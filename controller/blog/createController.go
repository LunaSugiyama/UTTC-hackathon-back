package blog

import (
	"log"
	"net/http"
	"time"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Claims struct {
	UserID string `json:"user_firebase_uid"`
	jwt.StandardClaims
}

func CreateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert a new blog entry into the database
	query := "INSERT INTO blogs (user_firebase_uid, title, author, link, item_categories_id, explanation, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, dbErr := database.DB.Exec(query, blog.UserFirebaseUID, blog.Title, blog.Author, blog.Link, blog.ItemCategoriesID, blog.Explanation, time.Now())
	if dbErr != nil {
		log.Printf("Error inserting blog entry into the database: %v", dbErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert blog entry"})
		return
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last inserted ID"})
		return
	}

	var createdAt mysql.NullTime
	var updatedAt mysql.NullTime
	// Retrieve the inserted blog entry
	selectQuery := "SELECT * FROM blogs WHERE id = ?"
	err = database.DB.QueryRow(selectQuery, lastInsertID).Scan(&blog.ID, &blog.UserFirebaseUID, &blog.Title, &blog.Author, &blog.Link, &blog.Likes, &blog.ItemCategoriesID, &blog.Explanation, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving the created blog entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created blog entry"})
		return
	}
	blog.CreatedAt = createdAt.Time
	blog.UpdatedAt = updatedAt.Time

	// Now, insert rows into the item_curriculums table for each curriculum ID
	for _, curriculumID := range blog.CurriculumIDs {
		insertItemCurriculumQuery := "INSERT INTO item_curriculums (item_id, item_categories_id, curriculum_id) VALUES (?, ?, ?)"
		_, insertErr := database.DB.Exec(insertItemCurriculumQuery, blog.ID, blog.ItemCategoriesID, curriculumID)
		if insertErr != nil {
			log.Printf("Error inserting into item_curriculums table: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into item_curriculums table"})
			return
		}
	}

	c.JSON(http.StatusOK, blog)
}
