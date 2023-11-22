package dao

import (
	"fmt"
	"log"
	"net/http"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type ItemDAO interface {
	GetItems(itemCategories []int, curriculumIDs []int, itemCategoriesSQL string, curriculumIDsSQL string) ([]model.Item, error)
	LikeItem(likedItem *model.LikedItem) error
	UnlikeItem(likedItem *model.LikedItem) error
	StarItem(starredItem *model.StarredItem) error
	UnstarItem(starredItem *model.StarredItem) error
	CheckLiked(likedItem *model.LikedItem) (bool, error)
	CheckStarred(starredItem *model.StarredItem) (bool, error)
	CountLikes(itemID int, itemCategoriesID int) (int, error)
}

type itemDAO struct {
	// db connection
}

func NewItemDAO() ItemDAO {
	return &itemDAO{
		// db connection
	}
}

func (dao *itemDAO) GetItems(itemCategories []int, curriculumIDs []int, itemCategoriesSQL string, curriculumIDsSQL string) ([]model.Item, error) {
	var items []model.Item

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
				query += "AND i.item_categories_id IN (" + itemCategoriesSQL + ") "
			}
			if len(curriculumIDs) > 0 {
				query += "AND ic.curriculum_id IN (" + curriculumIDsSQL + ") "
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
				var item model.Item

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

func (dao *itemDAO) LikeItem(likedItem *model.LikedItem) error {
	// Retrieve the name from the "item_categories" table
	itemName, err := getItemCategoryName(likedItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	// Check if the name exists
	if itemName == "" {
		return fmt.Errorf("item category with ID %d does not exist", likedItem.ItemCategoriesID)
	}

	// Check if the item exists in the table with the retrieved name
	if recordLikedExists(likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID) {
		// Item already exists, return an "already starred" message
		return &gin.Error{
			Err:  fmt.Errorf("item is already liked"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"status": http.StatusConflict},
		}
	}

	getItemCategoriesQuery := "SELECT name FROM item_categories WHERE id = ?"
	err = database.DB.QueryRow(getItemCategoriesQuery, likedItem.ItemCategoriesID).Scan(&itemName)
	if err != nil {
		return err
	}
	fmt.Println(itemName)

	// Update the "likes" column in the "items" table for the liked item
	query := fmt.Sprintf("UPDATE %s SET likes = likes + 1 WHERE id = ?", itemName)
	_, err = database.DB.Exec(query, likedItem.ItemID)
	if err != nil {
		return err
	}

	// Insert a new record into the "starred_items" table
	_, err = database.DB.Exec("INSERT INTO liked_items (user_firebase_uid, item_id, item_categories_id) VALUES (?, ?, ?)", likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve the name from the "item_categories" table
func getItemCategoryName(itemCategoriesID int) (string, error) {
	var itemName string
	err := database.DB.QueryRow("SELECT name FROM item_categories WHERE id = ?", itemCategoriesID).Scan(&itemName)
	if err != nil {
		return "", err
	}
	return itemName, nil
}

// Check if the item with a specific name and itemID exists in a table with a dynamic name
func recordLikedExists(userFirebaseUID string, itemID int, itemCategoriesID int) bool {
	var count int
	query := "SELECT COUNT(*) FROM liked_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, userFirebaseUID, itemID, itemCategoriesID).Scan(&count)
	if err != nil {
		// Handle any errors here, such as database connection issues
		return false
	}
	return count > 0
}

func (dao *itemDAO) UnlikeItem(likedItem *model.LikedItem) error {
	// Check if the item to unlike exists in the "liked_items" table
	if !recordLikedExists(likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID) {
		// Item does not exist, return an error
		return &gin.Error{
			Err:  fmt.Errorf("item is not liked"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"status": http.StatusNotFound},
		}
	}

	// Remove the record from the "liked_items" table
	_, err := database.DB.Exec("DELETE FROM liked_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?", likedItem.UserFirebaseUID, likedItem.ItemID, likedItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	// Retrieve the name from the "item_categories" table
	itemName, err := getItemCategoryName(likedItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	getItemCategoriesQuery := "SELECT name FROM item_categories WHERE id = ?"
	err = database.DB.QueryRow(getItemCategoriesQuery, likedItem.ItemCategoriesID).Scan(&itemName)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET likes = likes - 1 WHERE id = ?", itemName)
	_, err = database.DB.Exec(query, likedItem.ItemID)
	if err != nil {
		return err
	}

	return nil
}

func (dao *itemDAO) StarItem(starredItem *model.StarredItem) error {
	// Retrieve the name from the "item_categories" table
	itemName, err := getItemCategoryName(starredItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	// Check if the name exists
	if itemName == "" {
		return fmt.Errorf("item category with ID %d does not exist", starredItem.ItemCategoriesID)
	}

	// Check if the item exists in the table with the retrieved name
	if recordExists(starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID) {
		// Item already exists, return an "already starred" message
		return &gin.Error{
			Err:  fmt.Errorf("item is already starred"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"status": http.StatusConflict},
		}
	}

	// Insert a new record into the "starred_items" table
	_, err = database.DB.Exec("INSERT INTO starred_items (user_firebase_uid, item_id, item_categories_id) VALUES (?, ?, ?)", starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	return nil
}

// Check if the item with a specific name and itemID exists in a table with a dynamic name
func recordExists(UserFirebaseUID string, itemID, itemCategoriesID int) bool {
	var count int
	query := "SELECT COUNT(*) FROM starred_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, UserFirebaseUID, itemID, itemCategoriesID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (dao *itemDAO) UnstarItem(starredItem *model.StarredItem) error {
	// Check if the item to unstar exists in the "starred_items" table
	if !recordExists(starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID) {
		// Item does not exist, return an error
		return &gin.Error{
			Err:  fmt.Errorf("item is not starred"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"status": http.StatusNotFound},
		}
	}

	// Remove the record from the "starred_items" table
	_, err := database.DB.Exec("DELETE FROM starred_items WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?", starredItem.UserFirebaseUID, starredItem.ItemID, starredItem.ItemCategoriesID)
	if err != nil {
		return err
	}

	return nil
}

func (dao *itemDAO) CheckLiked(likedItem *model.LikedItem) (bool, error) {
	var isLiked int
	query := "SELECT COUNT(*) FROM liked_items WHERE item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, likedItem.ItemID, likedItem.ItemCategoriesID).Scan(&isLiked)
	if err != nil {
		return false, err
	}
	if isLiked > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (dao *itemDAO) CheckStarred(starredItem *model.StarredItem) (bool, error) {
	var isStarred int
	query := "SELECT COUNT(*) FROM starred_items WHERE item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, starredItem.ItemID, starredItem.ItemCategoriesID).Scan(&isStarred)
	if err != nil {
		return false, err
	}
	if isStarred > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (dao *itemDAO) CountLikes(itemID int, itemCategoriesID int) (int, error) {
	var likes int
	query := "SELECT COUNT(*) FROM liked_items WHERE item_id = ? AND item_categories_id = ?"
	err := database.DB.QueryRow(query, itemID, itemCategoriesID).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}
