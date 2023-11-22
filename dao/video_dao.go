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

type VideoDAO interface {
	SaveVideo(video *model.Video) error
	GetVideoByID(id int) (model.Video, error)
	UpdateVideo(video *model.Video) (model.Video, error)
	DeleteVideo(id int) error
	ShowAllVideos() ([]model.Video, error)
}

type videoDAO struct {
	// Fields if needed
}

func NewVideoDAO() VideoDAO {
	return &videoDAO{}
}

func (dao *videoDAO) SaveVideo(video *model.Video) error {
	// Insert a new video entry into the database
	query := "INSERT INTO videos (user_firebase_uid, title, author, link, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, video.UserFirebaseUID, video.Title, video.Author, video.Link, video.ItemCategoriesID, video.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting video entry into the database: %v", dbErr)
		return dbErr
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return err
	}
	video.ID = int(lastInsertID) // Convert lastInsertID from int64 to int

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Retrieve the inserted video entry
	selectQuery := "SELECT created_at, updated_at FROM videos WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created video entry: %v", err)
		return err
	}
	video.CreatedAt = createdAt.Time
	video.UpdatedAt = updatedAt.Time

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range video.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return insertErr
		}
	}

	// Now, insert rows into the item_images table for each image
	for _, image := range video.Images {
		insertItemImageQuery := "INSERT INTO item_images (item_id, item_categories_id, images, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemImageQuery, video.ID, video.ItemCategoriesID, image, video.CreatedAt, video.UpdatedAt)
		if insertErr != nil {
			log.Printf("Error inserting into item_images table: %v", insertErr)
			return insertErr
		}
	}
	return nil
}

func (dao *videoDAO) GetVideoByID(id int) (model.Video, error) {
	var video model.Video
	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs, images sql.NullString

	// Query the database to retrieve the video entry by ID
	query := `
	SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids, GROUP_CONCAT(IFNULL(ii.images, '')) AS images
	FROM videos AS b
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
		&video.ID, &video.UserFirebaseUID, &video.Title, &video.Author, &video.Link, &video.Likes,
		&video.ItemCategoriesID, &video.Explanation, &createdAt, &updatedAt, &video.ItemCategoriesName, &curriculumIDs, &images,
	)
	if err != nil {
		return model.Video{}, err
	}
	video.CreatedAt = createdAt.Time
	video.UpdatedAt = updatedAt.Time

	// Process curriculum IDs
	if curriculumIDs.Valid {
		curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
		for _, idStr := range curriculumIDSlice {
			if id, err := strconv.Atoi(idStr); err == nil {
				video.CurriculumIDs = append(video.CurriculumIDs, id)
			}
		}
	}

	// Process images
	if images.Valid {
		imagesSlice := strings.Split(images.String, ",")
		video.Images = append(video.Images, imagesSlice...)
	}

	return video, nil
}

func (dao *videoDAO) UpdateVideo(video *model.Video) (model.Video, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return model.Video{}, err
	}

	// Check if video exists
	checkQuery := "SELECT id FROM videos WHERE id = ?"
	err = tx.QueryRow(checkQuery, video.ID).Scan(&video.ID)
	if err != nil {
		tx.Rollback() // Corrected to use tx instead of database.DB
		if err == sql.ErrNoRows {
			return model.Video{}, fmt.Errorf("video with ID %d does not exist", video.ID)
		}
		return model.Video{}, err
	}

	// Update the video
	updateQuery := "UPDATE videos SET title = ?, author = ?, link = ?, user_firebase_uid = ?, item_categories_id = ?, explanation = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery, video.Title, video.Author, video.Link, video.UserFirebaseUID, video.ItemCategoriesID, video.Explanation, video.ID)
	if err != nil {
		tx.Rollback()
		return model.Video{}, err
	}

	// Update the curriculum IDs for the video in the item_curriculums table
	// First, delete the existing entries for this video
	deleteQuery := "DELETE FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?"
	_, err = tx.Exec(deleteQuery, video.ID, video.ItemCategoriesID)
	if err != nil {
		tx.Rollback()
		return model.Video{}, err
	}

	// Now, insert the updated curriculum IDs
	for _, curriculumID := range video.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := tx.Exec(insertItemCurriculumQuery, video.ID, video.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			tx.Rollback()
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			return model.Video{}, insertErr
		}
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Query the database to retrieve the updated video entry
	selectQuery := "SELECT * FROM videos WHERE id = ?"
	err = tx.QueryRow(selectQuery, video.ID).Scan(&video.ID, &video.UserFirebaseUID, &video.Title, &video.Author, &video.Link, &video.Likes, &video.ItemCategoriesID, &video.Explanation, &createdAt, &updatedAt)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving the updated video entry: %v", err)
		return model.Video{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return model.Video{}, err
	}

	video.CreatedAt = createdAt.Time
	video.UpdatedAt = updatedAt.Time
	return *video, nil
}

func (dao *videoDAO) DeleteVideo(id int) error {
	// Query the database to delete the video entry by ID
	query := "DELETE FROM videos WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting video entry: %v", err)
		return err
	}
	return nil
}

func (dao *videoDAO) ShowAllVideos() ([]model.Video, error) {
	// Query the database to retrieve all video entries with their associated curriculum IDs
	query := `
		SELECT b.*, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
		FROM videos AS b
		LEFT JOIN item_curriculums AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
		GROUP BY b.id
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying video entries: %v", err)
		return nil, err
	}
	defer rows.Close()

	var videos []model.Video

	for rows.Next() {
		var video model.Video
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs sql.NullString

		if err := rows.Scan(
			&video.ID, &video.UserFirebaseUID, &video.Title, &video.Author, &video.Link, &video.Likes, &video.ItemCategoriesID,
			&video.Explanation, &createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning video entry: %v", err)
			return nil, err
		}
		video.CreatedAt = createdAt.Time
		video.UpdatedAt = updatedAt.Time

		// Process curriculum IDs
		if curriculumIDs.Valid {
			curriculumIDSlice := strings.Split(curriculumIDs.String, ",")
			for _, idStr := range curriculumIDSlice {
				if id, err := strconv.Atoi(idStr); err == nil {
					video.CurriculumIDs = append(video.CurriculumIDs, id)
				}
			}
		}

		videos = append(videos, video)
	}
	return videos, nil
}
