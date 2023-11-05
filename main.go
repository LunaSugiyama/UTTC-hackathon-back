// main.go
package main

import (
	"uttc-hackathon/controller/blog"
	"uttc-hackathon/controller/book"
	"uttc-hackathon/controller/curriculum"
	"uttc-hackathon/controller/item"
	itemcategory "uttc-hackathon/controller/item_category"
	"uttc-hackathon/controller/user"
	"uttc-hackathon/controller/video"
	"uttc-hackathon/middlewares"

	"uttc-hackathon/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitializeDB()
	// err := firebaseinit.InitFirebase() // Use the updated package name
	// if err != nil {
	// 	log.Fatalf("Error initializing Firebase: %v\n", err)
	// }
	// database.CreateItemCategoriesTable()
	// database.CreateCurriculumsTable()
	// database.CreateItemCurriculumsTable()
	// database.DropUsersTable()
	// database.CreateUsersTable()
	// database.AddFirebaseUID()
	// database.PasswordColumnToNull()
	// database.DropBlogsTable()
	// database.CreateBlogsTable()
	// database.DropBooksTable()
	// database.CreateBooksTable()
	// database.DropVideosTable()
	// database.CreateVideosTable()
	// database.CreateStarredItemsTable()
	// ginMode := os.Getenv("GIN_MODE")
	// gin.SetMode(ginMode)

	r := gin.Default()
	r.Use(middlewares.CORS())
	// Register the AuthMiddleware for routes that require authentication.
	r.Use(func(c *gin.Context) {
		if c.FullPath() != "/users/register" {
			middlewares.AuthMiddleware()(c)
		}
	})

	usersGroup := r.Group("/users")
	{
		// usersGroup.POST("/register", user.Register)
		usersGroup.POST("/login", user.Login)
	}

	itemsGroup := r.Group("/items")
	{
		itemsGroup.GET("/showall", item.GetAllItems)
		itemsGroup.GET("/checkstarred", item.IsItemStarred)
		itemsGroup.POST("/star", item.StarItem)
		itemsGroup.POST("/unstar", item.UnstarItem)
	}

	itemcategoryGroup := r.Group("/item_categories")
	{
		itemcategoryGroup.GET("/showall", itemcategory.ShowAllItemCategories)
	}

	blogsGroup := r.Group("/blogs")
	{
		blogsGroup.POST("/create", blog.CreateBlog)
		blogsGroup.PUT("/update", blog.UpdateBlog)
		blogsGroup.GET("/get", blog.GetBlog)
		blogsGroup.GET("/showall", blog.ShowAllBlogs)
		blogsGroup.DELETE("/delete", blog.DeleteBlog)
	}

	booksGroup := r.Group("/books")
	{
		booksGroup.POST("/create", book.CreateBook)
		booksGroup.PUT("/update", book.UpdateBook)
		booksGroup.GET("/get", book.GetBook)
		booksGroup.GET("/showall", book.ShowAllBooks)
		booksGroup.DELETE("/delete", book.DeleteBook)
	}

	videosGroup := r.Group("/videos")
	{
		videosGroup.POST("/create", video.CreateVideo)
		videosGroup.PUT("/update", video.UpdateVideo)
		videosGroup.GET("/get", video.GetVideo)
		videosGroup.GET("/showall", video.ShowAllVideos)
		videosGroup.DELETE("/delete", video.DeleteVideo)
	}

	curriculumsGroup := r.Group("/curriculums")
	{
		curriculumsGroup.POST("/create", curriculum.CreateCurriculum)
		curriculumsGroup.PUT("/update", curriculum.UpdateCurriculum)
		curriculumsGroup.GET("/get", curriculum.GetCurriculum)
		curriculumsGroup.GET("/showall", curriculum.ShowAllCurriculums)
		curriculumsGroup.DELETE("/delete", curriculum.DeleteCurriculum)
	}

	r.Run(":8000")
}
