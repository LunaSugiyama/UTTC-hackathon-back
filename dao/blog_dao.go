package dao

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/go-sql-driver/mysql"
)

type BlogDAO interface {
	SaveBlog(blog *model.Blog) error
	GetBlogByID(id int) (model.Blog, error)
	UpdateBlog(blog *model.Blog) (model.Blog, error)
	DeleteBlog(id int) error
	ShowAllBlogs() ([]model.Blog, error)
}

type blogDAO struct {
	// Fields if needed
}

func NewBlogDAO() BlogDAO {
	return &blogDAO{}
}

func (dao *blogDAO) SaveBlog(blog *model.Blog) error {
	// Insert a new blog entry into the database
	query := "INSERT INTO blogs (user_firebase_uid, title, author, link, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, blog.UserFirebaseUID, blog.Title, blog.Author, blog.Link, blog.ItemCategoriesID, blog.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting blog entry into the database: %v", dbErr)
		return dbErr
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return err
	}
	blog.ID = int(lastInsertID) // Convert lastInsertID from int64 to int

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Retrieve the inserted blog entry
	selectQuery := "SELECT created_at, updated_at FROM blogs WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created blog entry: %v", err)
		return err
	}
	blog.CreatedAt = createdAt.Time
	blog.UpdatedAt = updatedAt.Time

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range blog.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, blog.ID, blog.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return insertErr
		}
	}

	// Now, insert rows into the item_images table for each image
	for _, image := range blog.Images {
		insertItemImageQuery := "INSERT INTO item_images (item_id, item_categories_id, images, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemImageQuery, blog.ID, blog.ItemCategoriesID, image, blog.CreatedAt, blog.UpdatedAt)
		if insertErr != nil {
			log.Printf("Error inserting into item_images table: %v", insertErr)
			return insertErr
		}
	}
	return nil
}

func (dao *blogDAO) GetBlogByID(id int) (model.Blog, error) {
	var blog model.Blog
	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs, images sql.NullString

	// Query the database to retrieve the blog entry by ID
	query := `
	SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids, GROUP_CONCAT(IFNULL(ii.images, '')) AS images
	FROM blogs AS b
	LEFT JOIN item_categories AS icat ON b.item_categories_id = icat.id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(images) AS images
		FROM item_images
		GROUP BY item_id, item_categories_id
	) AS ii ON b.id = ii.item_id AND b.item_categories_id = ii.item_categories_id
	LEFT JOIN (
		SELECT item_id, item_categories_id, GROUP_CONCAT(curriculum_id) AS curriculum_id
		FROM item_curriculums
		GROUP BY item_id, item_categories_id
	) AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
	WHERE b.id = ?
	GROUP BY b.id
	`
	err := database.DB.QueryRow(query, id).Scan(
		&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes,
		&blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt, &blog.ItemCategoriesName, &curriculumIDs, &images,
	)
	if err != nil {
		return model.Blog{}, err
	}
	blog.CreatedAt = createdAt.Time
	blog.UpdatedAt = updatedAt.Time

	// Process curriculum IDs
	if curriculumIDs.Valid {
		curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
		for _, idStr := range curriculumIDSlice {
			if id, err := strconv.Atoi(idStr); err == nil {
				blog.CurriculumIDs = append(blog.CurriculumIDs, id)
			}
		}
	}

	// Process images
	if images.Valid {
		imagesSlice := strings.Split(images.String, ",")
		blog.Images = append(blog.Images, imagesSlice...)
	}

	return blog, nil
}

func (dao *blogDAO) UpdateBlog(blog *model.Blog) (model.Blog, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return model.Blog{}, err
	}

	// Check if blog exists
	checkQuery := "SELECT id FROM blogs WHERE id = ?"
	err = tx.QueryRow(checkQuery, blog.ID).Scan(&blog.ID)
	if err != nil {
		tx.Rollback() // Corrected to use tx instead of database.DB
		if err == sql.ErrNoRows {
			return model.Blog{}, fmt.Errorf("blog with ID %d does not exist", blog.ID)
		}
		return model.Blog{}, err
	}

	// Update the blog
	updateQuery := "UPDATE blogs SET title = ?, author = ?, link = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery, blog.Title, blog.Author, blog.Link, blog.UserFirebaseUID, blog.ItemCategoriesID, blog.Explanation, blog.ID)
	if err != nil {
		tx.Rollback()
		return model.Blog{}, err
	}

	// Update the curriculum IDs for the blog in the item_curriculums table
	// First, delete the existing entries for this blog
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, err = tx.Exec(deleteQuery, blog.ID, blog.ItemCategoriesID)
	if err != nil {
		tx.Rollback()
		return model.Blog{}, err
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range blog.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := tx.Exec(insertItemCurriculumQuery, blog.ID, blog.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			tx.Rollback()
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return model.Blog{}, insertErr
		}
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated blog entry
	selectQuery := "SELECT * FROM blogs WHERE id = ?"
	err = tx.QueryRow(selectQuery, blog.ID).Scan(&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes, &blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving the updated blog entry: %v", err)
		return model.Blog{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return model.Blog{}, err
	}

	blog.CreatedAt = createdAt.Time
	blog.UpdatedAt = updatedAt.Time
	return *blog, nil
}

func (dao *blogDAO) DeleteBlog(id int) error {
	// Query the database to delete the blog entry by ID
	query := "DELETE FROM blogs WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting blog entry: %v", err)
		return err
	}
	return nil
}

func (dao *blogDAO) ShowAllBlogs() ([]model.Blog, error) {
	// Query the database to retrieve all blog entries with their associated curriculum IDs
	query := `
		SELECT b.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
		FROM blogs AS b
		LEFT JOIN item_curriculums AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
		GROUP BY b.id
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying blog entries: %v", err)
		return nil, err
	}
	defer rows.Close()

	var blogs []model.Blog

	for rows.Next() {
		var blog model.Blog
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs sql.NullString

		if err := rows.Scan(
			&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes, &blog.ItemCategoriesID,
			&blog.Explanation, &createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning blog entry: %v", err)
			return nil, err
		}
		blog.CreatedAt = createdAt.Time
		blog.UpdatedAt = updatedAt.Time

		// Split the curriculum IDs into a slice
		if curriculumIDs.Valid {
			curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
			for _, idStr := range curriculumIDSlice {
				if id, err := strconv.Atoi(idStr); err == nil {
					blog.CurriculumIDs = append(blog.CurriculumIDs, id)
				}
			}
		}

		blogs = append(blogs, blog)
	}
	return blogs, nil
}
