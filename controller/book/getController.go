package book

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func GetBook(c *gin.Context) {
	var book model.Book
	// Get the 'id' parameter from the query string
	bookID := c.Query("id")

	// Check if 'id' is empty or not provided
	if bookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	var curriculumIDs string
	var Images string

	// Convert the bookID to an integer
	id, err := strconv.Atoi(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the book entry by ID, including associated curriculum IDs
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
	err = database.DB.QueryRow(query, id).Scan(
		&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation,
		&createdAt, &updatedAt, &book.ItemCategoriesName, &curriculumIDs, &Images,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// The book entry doesn't exist, return a not found response
			c.JSON(http.StatusNotFound, gin.H{"error": "Book entry not found"})
			return
		} else {
			// Handle other database errors
			log.Printf("Error retrieving book entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book entry"})
			return
		}
	}
	book.CreatedAt = createdAt.Time
	book.UpdatedAt = updatedAt.Time

	// Split the curriculum IDs into a slice
	curriculumIDSlice := strings.Split(curriculumIDs, ",")
	for _, idStr := range curriculumIDSlice {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			book.CurriculumIDs = append(book.CurriculumIDs, id)
		}
	}

	// Split the images into a slice
	imageslice := strings.Split(Images, ",")
	book.Images = append(book.Images, imageslice...)

	// Book entry found, return it as JSON
	c.JSON(http.StatusOK, book)
}
