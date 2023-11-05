package item

// func LikeItem(c *gin.Context) {
// 	// Get the 'id' parameter from the query string
// 	itemID := c.Query("id")

// 	// Check if 'id' is empty or not provided
// 	if itemID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing or empty"})
// 		return
// 	}

// 	// Convert the itemID to an integer
// 	id, err := strconv.Atoi(itemID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
// 		return
// 	}

// 	// Get the user's Firebase UID from the request context
// 	userFirebaseUID, ok := c.Get("user_firebase_uid")
// 	if !ok {
// 		// The user's Firebase UID is not present in the request context, return an error
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Firebase UID from request context"})
// 		return
// 	}

// 	// Query the database to check if the user has already liked the item
// 	query := `
// 		SELECT *
// 		FROM starred_items
// 		WHERE user_firebase_uid = ? AND item_id = ?
// 	`
// 	var starredItem model.StarredItem
// 	err = database.DB.QueryRow(query, userFirebaseUID, id).Scan(
// 		&starredItem.ID, &starredItem.UserFirebaseUID, &starredItem.ItemID, &starredItem.ItemCategoriesID, &starredItem.CreatedAt, &starredItem.UpdatedAt,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// The user has not liked the item yet, create a new entry in the starred_items table
// 			query := `
// 				INSERT INTO starred_items (user_firebase_uid, item_id, item_categories_id, created_at, updated_at)
// 				VALUES (?, ?, ?, NOW(), NOW())
// 			`
// 			_, err := database.DB.Exec(query, userFirebaseUID, id, item.ItemCategoriesID)
// 			if err != nil {
// 				log.Printf("Error creating starred item entry: %v", err)
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create starred item entry"})
// 				return
// 			}

// 			// Increment the item's likes count
// 			query = `
// 				UPDATE items
// 				SET likes = likes + 1
// 				WHERE id = ? AND item_categories_id = ?
