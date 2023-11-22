// main.go
package main

import (
	"uttc-hackathon/database"
	"uttc-hackathon/router"
)

func main() {
	database.InitializeDB()
	// err := firebaseinit.InitFirebase() // Use the updated package name
	// if err != nil {
	// 	log.Fatalf("Error initializing Firebase: %v\n", err)
	// }
	database.CreateItemCategoriesTable()
	// database.PopulateItemCategoriesTable()
	database.CreateCurriculumsTable()
	// database.PopulateCurriculumsTable()
	database.CreateItemCurriculumsTable()
	database.CreateUsersTable()
	database.CreateBlogsTable()
	database.CreateBooksTable()
	database.CreateVideosTable()
	database.CreateStarredItemsTable()
	database.CreateLikedItemsTable()
	database.CreateItemImagesTable()
	database.CreateItemCommentsTable()

	// // Register the AuthMiddleware for routes that require authentication.
	// r.Use(func(c *gin.Context) {
	// 	if c.FullPath() != "/users/register" {
	// 		middlewares.AuthMiddleware()(c)
	// 	}
	// })

	// r := gin.Default()
	// r.Use(middlewares.CORS())

	r := router.SetupRouter() // Pass the Gin engine instance as an argument

	r.Run(":8000")
}
