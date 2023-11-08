package item

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"uttc-hackathon/database"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type ItemWithExplanation struct {
	ID                 int
	UserFirebaseUID    string
	Title              string
	Author             string
	Link               string
	Likes              int
	Explanation        string
	ItemCategoriesID   int
	ItemCategoriesName string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CurriculumIDs      []int // New field to store curriculum_ids
}

type seachRequest struct {
	Words          string `json:"words"`
	Sort           string `json:"sorting"`
	Order          string `json:"order"`
	ItemCategories string `json:"item_categories"`
	CurriculumIDs  string `json:"curriculum_ids"`
}

func SearchItems(c *gin.Context) {
	// Define a struct to bind the request data
	var searchRequest seachRequest

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&searchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(searchRequest)

	// Split the query string into terms based on whitespace
	terms := strings.Split(searchRequest.Words, " ")
	fmt.Println(terms)

	items, err := getAllItemsInfo(searchRequest.Sort, searchRequest.Order, searchRequest.ItemCategories, searchRequest.CurriculumIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
		return
	}
	fmt.Println("item: ", items)

	// Perform the search logic
	results := performSearch(terms, items, searchRequest.Sort, searchRequest.Order)

	// Return the results
	c.JSON(http.StatusOK, results)
}

func performSearch(terms []string, items []ItemWithExplanation, sorting string, order string) []ItemWithExplanation {
	// Perform your search logic here
	var searchResults []ItemWithExplanation

	for _, item := range items {
		for _, term := range terms {
			nameContainsTerm := strings.Contains(strings.ToLower(item.Title), strings.ToLower(term))
			explanationContainsTerm := strings.Contains(strings.ToLower(item.Explanation), strings.ToLower(term))

			if nameContainsTerm || explanationContainsTerm {
				searchResults = append(searchResults, item)
			}
		}
	}

	// Check if sortinging criteria are provided
	if sorting != "" {
		// Implement sortinging based on the 'sorting' and 'order' parameters
		switch sorting {
		case "created_at":
			// sorting by the 'CreatedAt' field
			if order == "asc" {
				sort.Slice(searchResults, func(i, j int) bool {
					return searchResults[i].CreatedAt.Before(searchResults[j].CreatedAt)
				})
			} else if order == "desc" {
				sort.Slice(searchResults, func(i, j int) bool {
					return searchResults[i].CreatedAt.After(searchResults[j].CreatedAt)
				})
			}
		// Add cases for other sortinging criteria if needed
		case "likes":
			sort.Slice(items, func(i, j int) bool {
				if order == "asc" {
					return items[i].Likes < items[j].Likes
				} else {
					return items[i].Likes > items[j].Likes
				}
			})
		case "updated_at":
			// sorting by the 'UpdatedAt' field
			if order == "asc" {
				sort.Slice(searchResults, func(i, j int) bool {
					return searchResults[i].UpdatedAt.Before(searchResults[j].UpdatedAt)
				})
			} else if order == "desc" {
				sort.Slice(searchResults, func(i, j int) bool {
					return searchResults[i].UpdatedAt.After(searchResults[j].UpdatedAt)
				})
			}
		default:
			// Handle unsupported sortinging criteria
		}
	}

	return searchResults
}

func getAllItemsInfo(sortingField string, sortingOrder string, itemCategoriesStr string, curriculumIDsStr string) ([]ItemWithExplanation, error) {

	var items []ItemWithExplanation

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
		return nil, err
	}

	// Create a map to track unique item categories for each item
	uniqueItemCategories := make(map[string]bool)

	// Iterate through curriculum IDs and table names to retrieve filtered items
	for _, curriculumID := range curriculumIDsAll {
		// Iterate through the tableNames (or categories) and retrieve items for each curriculum
		tableNames, err := getTableNames()
		if err != nil {
			return nil, err
		}

		// Iterate through table names to retrieve items for each category
		for _, tableName := range tableNames {
			// Build the query based on curriculum and category
			query := fmt.Sprintf("SELECT i.id, i.user_firebase_uid, i.title, i.author, i.link, i.explanation, i.likes, i.item_categories_id, icat.name AS item_category_name, i.created_at, i.updated_at FROM %s AS i "+
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
				return nil, err
			}
			defer rows.Close()

			// Retrieve and append data from the table
			for rows.Next() {
				var item ItemWithExplanation

				var CreatedAt, UpdatedAt mysql.NullTime

				// Scan the data into the item struct
				if err := rows.Scan(
					&item.ID, &item.UserFirebaseUID, &item.Title, &item.Author, &item.Link, &item.Explanation,
					&item.Likes, &item.ItemCategoriesID, &item.ItemCategoriesName, &CreatedAt, &UpdatedAt); err != nil {
					fmt.Println("error in scan")
					return nil, err
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
					fmt.Println("error: Failed to retrieve curriculum IDs for item")
					return nil, err
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
	return items, nil
}
