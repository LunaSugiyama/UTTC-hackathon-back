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

type BookDAO interface {
	SaveBook(book *model.Book) error
	GetBookByID(id int) (model.Book, error)
	UpdateBook(book *model.Book) (model.Book, error)
	DeleteBook(id int) error
	ShowAllBooks() ([]model.Book, error)
}

type bookDAO struct {
	// Fields if needed
}

func NewBookDAO() BookDAO {
	return &bookDAO{}
}

func (dao *bookDAO) SaveBook(book *model.Book) error {
	// Insert a new book entry into the database
	query := "INSERT INTO books (user_firebase_uid, title, author, link, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, book.UserFirebaseUID, book.Title, book.Author, book.Link, book.ItemCategoriesID, book.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting book entry into the database: %v", dbErr)
		return dbErr
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return err
	}
	book.ID = int(lastInsertID) // Convert lastInsertID from int64 to int

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Retrieve the inserted book entry
	selectQuery := "SELECT created_at, updated_at FROM books WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created book entry: %v", err)
		return err
	}
	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range book.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, book.ID, book.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return insertErr
		}
	}

	// Now, insert rows into the item_images table for each image
	for _, image := range book.Images {
		insertItemImageQuery := "INSERT INTO item_images (item_id, item_categories_id, images, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemImageQuery, book.ID, book.ItemCategoriesID, image, book.CreatedAt, book.UpdatedAt)
		if insertErr != nil {
			log.Printf("Error inserting into item_images table: %v", insertErr)
			return insertErr
		}
	}
	return nil
}

func (dao *bookDAO) GetBookByID(id int) (model.Book, error) {
	var book model.Book
	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs, images sql.NullString

	// Query the database to retrieve the book entry by ID
	query := `
	SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids, GROUP_CONCAT(IFNULL(ii.images, '')) AS images
	FROM books AS b
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
		&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes,
		&book.ItemCategoriesID, &book.Explanation, &createdAt, &updatedAt, &book.ItemCategoriesName, &curriculumIDs, &images,
	)
	if err != nil {
		return model.Book{}, err
	}
	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time

	// Process curriculum IDs
	if curriculumIDs.Valid {
		curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
		for _, idStr := range curriculumIDSlice {
			if id, err := strconv.Atoi(idStr); err == nil {
				book.CurriculumIDs = append(book.CurriculumIDs, id)
			}
		}
	}

	// Process images
	if images.Valid {
		imagesSlice := strings.Split(images.String, ",")
		book.Images = append(book.Images, imagesSlice...)
	}

	return book, nil
}

func (dao *bookDAO) UpdateBook(book *model.Book) (model.Book, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return model.Book{}, err
	}

	// Check if book exists
	checkQuery := "SELECT id FROM books WHERE id = ?"
	err = tx.QueryRow(checkQuery, book.ID).Scan(&book.ID)
	if err != nil {
		tx.Rollback() // Corrected to use tx instead of database.DB
		if err == sql.ErrNoRows {
			return model.Book{}, fmt.Errorf("book with ID %d does not exist", book.ID)
		}
		return model.Book{}, err
	}

	// Update the book
	updateQuery := "UPDATE books SET title = ?, author = ?, link = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery, book.Title, book.Author, book.Link, book.UserFirebaseUID, book.ItemCategoriesID, book.Explanation, book.ID)
	if err != nil {
		tx.Rollback()
		return model.Book{}, err
	}

	// Update the curriculum IDs for the book in the item_curriculums table
	// First, delete the existing entries for this book
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, err = tx.Exec(deleteQuery, book.ID, book.ItemCategoriesID)
	if err != nil {
		tx.Rollback()
		return model.Book{}, err
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range book.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := tx.Exec(insertItemCurriculumQuery, book.ID, book.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			tx.Rollback()
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return model.Book{}, insertErr
		}
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated book entry
	selectQuery := "SELECT * FROM books WHERE id = ?"
	err = tx.QueryRow(selectQuery, book.ID).Scan(&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation, &createdAt, &updatedAt)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving the updated book entry: %v", err)
		return model.Book{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return model.Book{}, err
	}

	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time
	return *book, nil
}

func (dao *bookDAO) DeleteBook(id int) error {
	// Query the database to delete the book entry by ID
	query := "DELETE FROM books WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting book entry: %v", err)
		return err
	}
	return nil
}

func (dao *bookDAO) ShowAllBooks() ([]model.Book, error) {
	// Query the database to retrieve all book entries with their associated curriculum IDs
	query := `
		SELECT b.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
		FROM books AS b
		LEFT JOIN item_curriculums AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
		GROUP BY b.id
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying book entries: %v", err)
		return nil, err
	}
	defer rows.Close()

	var books []model.Book

	for rows.Next() {
		var book model.Book
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs sql.NullString

		if err := rows.Scan(
			&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID,
			&book.Explanation, &createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning book entry: %v", err)
			return nil, err
		}
		book.CreatedAt = createdAt.Time
		book.UpdatedAt = updatedAt.Time

		// Split the curriculum IDs into a slice
		if curriculumIDs.Valid {
			curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
			for _, idStr := range curriculumIDSlice {
				if id, err := strconv.Atoi(idStr); err == nil {
					book.CurriculumIDs = append(book.CurriculumIDs, id)
				}
			}
		}

		books = append(books, book)
	}
	return books, nil
}
