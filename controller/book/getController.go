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

	// Convert the bookID to an integer
	id, err := strconv.Atoi(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Query the database to retrieve the book entry by ID, including associated curriculum IDs
	query := `
        SELECT b.*, icat.name, GROUP_CONCAT(ic.curriculum_id) AS curriculum_ids
        FROM books AS b
        LEFT JOIN item_curriculums AS ic ON b.id = ic.item_id AND b.item_categories_id = ic.item_categories_id
		LEFT JOIN item_categories AS icat ON b.item_categories_id = icat.id
        WHERE b.id = ?
        GROUP BY b.id
    `
	err = database.DB.QueryRow(query, id).Scan(
		&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation,
		&createdAt, &updatedAt, &book.ItemCategoriesName, &curriculumIDs,
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

	// Book entry found, return it as JSON
	c.JSON(http.StatusOK, book)
}
