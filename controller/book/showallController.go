package book

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ShowAllBooks(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book entries"})
		return
	}
	defer rows.Close()

	var books []model.Book

	for rows.Next() {
		var book model.Book
		var createdAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var updatedAt mysql.NullTime // Use mysql.NullTime for MySQL TIMESTAMP columns
		var curriculumIDs string

		if err := rows.Scan(
			&book.ID, &book.UserFirebaseUID, &book.Title, &book.Author, &book.Link, &book.Likes, &book.ItemCategoriesID, &book.Explanation,
			&createdAt, &updatedAt, &curriculumIDs,
		); err != nil {
			log.Printf("Error scanning book entry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan book entries"})
			return
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

		books = append(books, book)
	}

	// Return the list of books as JSON
	c.JSON(http.StatusOK, books)
}
