package item

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"uttc-hackathon/database"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var curriculums map[int]string
var itemToCurriculum map[int]int // Map item ID to curriculum ID

func getAllCurriculumIDs() ([]int, error) {
	var curriculumIDs []int

	// Query the database to fetch all curriculum IDs
	rows, err := database.DB.Query("SELECT id FROM curriculums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var curriculumID int

		if err := rows.Scan(&curriculumID); err != nil {
			return nil, err
		}

		curriculumIDs = append(curriculumIDs, curriculumID)
	}

	return curriculumIDs, nil
}

// Your existing code for GetAllItems and other functions here

func getTableNames() ([]string, error) {
	var tableNames []string

	query := "SELECT name FROM item_categories"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, name)
	}

	return tableNames, nil
}

type ItemWithCurriculum struct {
	ID                 int
	UserFirebaseUID    string
	Title              string
	Author             string
	Link               string
	Likes              int
	ItemCategoriesID   int
	ItemCategoriesName string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CurriculumIDs      []int // New field to store curriculum_ids
}

func GetAllItems(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size parameter"})
		return
	}

	var items []ItemWithCurriculum

	// Get the sorting criteria from the query parameters
	sortField := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "asc")

	// Get item_categories and curriculum_ids from query parameters
	itemCategoriesStr := c.Query("item_categories")
	curriculumIDsStr := c.Query("curriculum_ids")

	// Parse item_categories and curriculum_ids into slices of integers
	var itemCategories []int
	var curriculumIDs []int

	if itemCategoriesStr != "" {
		itemCategories = parseIntList(itemCategoriesStr)
	}

	if curriculumIDsStr != "" {
		curriculumIDs = parseIntList(curriculumIDsStr)
	}

	// Get all curriculum IDs from the categories table
	curriculumIDsAll, err := getAllCurriculumIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve curriculum IDs"})
		return
	}

	// Create a map to track unique item categories for each item
	uniqueItemCategories := make(map[string]bool)

	// Iterate through curriculum IDs and table names to retrieve filtered items
	for _, curriculumID := range curriculumIDsAll {
		// Iterate through the tableNames (or categories) and retrieve items for each curriculum
		tableNames, err := getTableNames()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve table names"})
			return
		}

		// Iterate through table names to retrieve items for each category
		for _, tableName := range tableNames {
			// Build the query based on curriculum and category
			query := fmt.Sprintf("SELECT i.id, i.user_firebase_uid, i.title, i.author, i.link, i.likes, i.item_categories_id, icat.name AS item_category_name, i.created_at, i.updated_at FROM %s AS i "+
				"INNER JOIN item_curriculums AS ic ON i.id = ic.item_id AND i.item_categories_id = ic.item_categories_id "+
				"INNER JOIN item_categories AS icat ON ic.item_categories_id = icat.id "+
				"WHERE ic.curriculum_id = ? ", tableName)

			// If item_categories or curriculum_ids are provided, modify the query
			if len(itemCategories) > 0 {
				query += "AND i.item_categories_id IN (" + intListToSQL(itemCategories) + ") "
			}
			if len(curriculumIDs) > 0 {
				query += "AND ic.curriculum_id IN (" + intListToSQL(curriculumIDs) + ") "
			}

			rows, err := database.DB.Query(query, curriculumID) // Provide both curriculumID and tableName as arguments

			if err != nil {
				// Log the error for debugging
				log.Printf("Error executing query: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from " + tableName})
				return
			}
			defer rows.Close()

			// Retrieve and append data from the table
			for rows.Next() {
				var item ItemWithCurriculum

				var CreatedAt, UpdatedAt mysql.NullTime

				// Scan the data into the item struct
				if err := rows.Scan(
					&item.ID, &item.UserFirebaseUID, &item.Title, &item.Author, &item.Link,
					&item.Likes, &item.ItemCategoriesID, &item.ItemCategoriesName, &CreatedAt, &UpdatedAt); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan data from " + tableName})
					return
				}

				// Convert NullTime to time.Time if not null
				if CreatedAt.Valid {
					item.CreatedAt = CreatedAt.Time
				}
				if UpdatedAt.Valid {
					item.UpdatedAt = UpdatedAt.Time
				}

				// Retrieve and add the curriculum IDs for this item
				curriculumIDs, err := getCurriculumIDsForItem(item.ID, item.ItemCategoriesID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve curriculum IDs for item"})
					return
				}
				item.CurriculumIDs = curriculumIDs // Assign the curriculum IDs to the item

				// Create a unique identifier for the item based on item_id and item_categories_id
				itemIdentifier := fmt.Sprintf("%d-%d", item.ID, item.ItemCategoriesID)

				// Check if this item is unique based on the item identifier
				if !uniqueItemCategories[itemIdentifier] {
					// If it's unique, append it to the items slice and mark it as displayed
					items = append(items, item)
					uniqueItemCategories[itemIdentifier] = true
				}

			}
		}
	}
	// Calculate the offset based on the page and page size
	offset := (page - 1) * pageSize

	start := offset
	end := offset + pageSize
	if start < 0 {
		start = 0
	}
	if end > len(items) {
		end = len(items)
	}
	pagedItems := items[start:end]
	// Implement sorting logic based on the specified sortField and sortOrder
	sortItems(pagedItems, sortField, sortOrder)

	// Return the combined and sorted data from all tables as JSON
	c.JSON(http.StatusOK, pagedItems)
}

// Helper function to parse a comma-separated list of integers
func parseIntList(input string) []int {
	parts := strings.Split(input, ",")
	var result []int
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err == nil {
			result = append(result, val)
		}
	}
	return result
}

// Helper function to convert a list of integers to a SQL-friendly string
func intListToSQL(list []int) string {
	var parts []string
	for _, val := range list {
		parts = append(parts, strconv.Itoa(val))
	}
	return strings.Join(parts, ",")
}

// Helper function to sort items by various fields
func sortItems(items []ItemWithCurriculum, sortField string, order string) {
	switch sortField {
	case "created_at", "updated_at":
		sortItemsByTime(items, sortField, order)
	case "item_categories_id":
		sortItemsByInt(items, sortField, order)
	case "curriculum_id":
		sortItemsByCurriculum(items, order, curriculums, itemToCurriculum)
	}
}

// Helper function to sort items by time (created_at or updated_at)
func sortItemsByTime(items []ItemWithCurriculum, sortField string, order string) {
	// Implement sorting logic based on the specified sortField and order
	switch sortField {
	case "created_at":
		sort.Slice(items, func(i, j int) bool {
			timeI := items[i].CreatedAt
			timeJ := items[j].CreatedAt

			if order == "asc" {
				return timeI.Before(timeJ)
			} else {
				return timeI.After(timeJ)
			}
		})
	case "updated_at":
		sort.Slice(items, func(i, j int) bool {
			timeI := items[i].UpdatedAt
			timeJ := items[j].UpdatedAt

			if order == "asc" {
				return timeI.Before(timeJ)
			} else {
				return timeI.After(timeJ)
			}
		})
	case "likes":
		sort.Slice(items, func(i, j int) bool {
			if order == "asc" {
				return items[i].Likes < items[j].Likes
			} else {
				return items[i].Likes > items[j].Likes
			}
		})
	}

}

// Helper function to sort items by integer field (e.g., item_categories_id)
func sortItemsByInt(items []ItemWithCurriculum, sortField string, order string) {
	switch order {
	case "asc":
		sort.Slice(items, func(i, j int) bool {
			switch sortField {
			case "item_categories_id":
				return items[i].ItemCategoriesID < items[j].ItemCategoriesID
			default:
				return false
			}
		})
	case "desc":
		sort.Slice(items, func(i, j int) bool {
			switch sortField {
			case "item_categories_id":
				return items[i].ItemCategoriesID > items[j].ItemCategoriesID
			default:
				return false
			}
		})
	}
}

// Helper function to sort items by curriculum
func sortItemsByCurriculum(items []ItemWithCurriculum, order string, curriculums map[int]string, itemToCurriculum map[int]int) {
	switch order {
	case "asc":
		sort.Slice(items, func(i, j int) bool {
			itemID1 := items[i].ID
			itemID2 := items[j].ID

			curriculumID1, exists1 := itemToCurriculum[itemID1]
			curriculumID2, exists2 := itemToCurriculum[itemID2]

			if exists1 && exists2 {
				return curriculumID1 < curriculumID2
			}
			return exists1
		})
	case "desc":
		sort.Slice(items, func(i, j int) bool {
			itemID1 := items[i].ID
			itemID2 := items[j].ID

			curriculumID1, exists1 := itemToCurriculum[itemID1]
			curriculumID2, exists2 := itemToCurriculum[itemID2]

			if exists1 && exists2 {
				return curriculumID1 > curriculumID2
			}
			return exists2
		})
	}
}

func getCurriculumIDsForItem(itemID int, ItemCategoriesID int) ([]int, error) {
	var curriculumIDs []int

	// Query the database to fetch curriculum IDs related to the item ID
	rows, err := database.DB.Query("SELECT curriculum_id FROM item_curriculums WHERE item_id = ? AND item_categories_id = ?", itemID, ItemCategoriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var curriculumID int

		if err := rows.Scan(&curriculumID); err != nil {
			return nil, err
		}

		curriculumIDs = append(curriculumIDs, curriculumID)
	}

	return curriculumIDs, nil
}
