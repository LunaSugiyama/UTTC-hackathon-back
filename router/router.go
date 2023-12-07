// router.go
package router

import (
	"uttc-hackathon/controller"
	"uttc-hackathon/dao"
	"uttc-hackathon/database"
	"uttc-hackathon/middlewares"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORS())

	// Setup for each group of routes
	setupUserRoutes(r)
	setupItemRoutes(r)
	setupItemCategoryRoutes(r)
	setupBlogRoutes(r)
	setupBookRoutes(r)
	setupVideoRoutes(r)
	setupCurriculumRoutes(r)

	return r
}

func setupUserRoutes(r *gin.Engine) {
	userDao := dao.NewUserDao()
	userUsecase := usecase.NewUserUsecase(userDao)
	userController := controller.NewUserController(userUsecase)

	usersGroup := r.Group("/users")
	{
		usersGroup.POST("/register", userController.RegisterUser)
		usersGroup.POST("/login", userController.LoginUser)
		usersGroup.GET("/show", userController.ShowUser)
		// usersGroup.PUT("/update", userController.Update)
	}
}

func setupItemRoutes(r *gin.Engine) {
	itemDao := dao.NewItemDAO()
	itemUsecase := usecase.NewItemUsecase(itemDao)
	itemController := controller.NewItemController(itemUsecase)

	relatedItemDao := dao.NewRelatedItemDAO(database.DB)
	relatedItemUsecase := usecase.NewRelatedItemUsecase(relatedItemDao)
	relatedItemController := controller.NewRelatedItemController(relatedItemUsecase)

	commentDao := dao.NewCommentDAO(database.DB)
	commentUsecase := usecase.NewCommentUsecase(commentDao)
	commentController := controller.NewCommentController(commentUsecase)

	itemsGroup := r.Group("/items")
	{
		commentsGroup := itemsGroup.Group("/comments")
		{
			commentsGroup.POST("/create", commentController.CreateComment)
			commentsGroup.GET("/get", commentController.GetComments)
			commentsGroup.PUT("/update", commentController.UpdateComment)
			commentsGroup.DELETE("/delete", commentController.DeleteComment)
		}
		itemsGroup.GET("/checkstarred", itemController.CheckStarred)
		itemsGroup.POST("/star", itemController.StarItem)
		itemsGroup.POST("/unstar", itemController.UnstarItem)
		itemsGroup.GET("/countlikes", itemController.CountLikes)
		itemsGroup.GET("/checkliked", itemController.CheckLiked)
		itemsGroup.POST("/like", itemController.LikeItem)
		itemsGroup.POST("/unlike", itemController.UnlikeItem)
		itemsGroup.POST("/search", itemController.SearchItems)
		itemsGroup.GET("/related", relatedItemController.GetItemsSimilarity)
	}
}

func setupItemCategoryRoutes(r *gin.Engine) {
	item_categoryDao := dao.NewItemCategoryDAO()
	item_categoryUsecase := usecase.NewItemCategoryUsecase(item_categoryDao)
	item_categoryController := controller.NewItemCategoryController(item_categoryUsecase)
	itemcategoryGroup := r.Group("/item_categories")
	{
		itemcategoryGroup.GET("/showall", item_categoryController.ShowAllItemCategories)
	}
}

func setupBlogRoutes(r *gin.Engine) {
	blogDao := dao.NewBlogDAO()
	blogUsecase := usecase.NewBlogUsecase(blogDao)
	blogController := controller.NewBlogController(blogUsecase)

	blogsGroup := r.Group("/blogs")
	{
		blogsGroup.POST("/create", blogController.CreateBlog)
		blogsGroup.PUT("/update", blogController.UpdateBlog)
		blogsGroup.GET("/get", blogController.GetBlog)
		blogsGroup.GET("/showall", blogController.ShowAllBlogs)
		blogsGroup.DELETE("/delete", blogController.DeleteBlog)
	}
}

func setupBookRoutes(r *gin.Engine) {
	bookDAO := dao.NewBookDAO()
	bookUsecase := usecase.NewBookUsecase(bookDAO)
	bookController := controller.NewBookController(bookUsecase)

	booksGroup := r.Group("/books")
	{
		booksGroup.POST("/create", bookController.CreateBook)
		booksGroup.PUT("/update", bookController.UpdateBook)
		booksGroup.GET("/get", bookController.GetBook)
		booksGroup.GET("/showall", bookController.ShowAllBooks)
		booksGroup.DELETE("/delete", bookController.DeleteBook)
	}
}

func setupVideoRoutes(r *gin.Engine) {
	videoDAO := dao.NewVideoDAO()
	videoUsecase := usecase.NewVideoUsecase(videoDAO)
	videoController := controller.NewVideoController(videoUsecase)

	videosGroup := r.Group("/videos")
	{
		videosGroup.POST("/create", videoController.CreateVideo)
		videosGroup.PUT("/update", videoController.UpdateVideo)
		videosGroup.GET("/get", videoController.GetVideo)
		videosGroup.GET("/showall", videoController.ShowAllVideos)
		videosGroup.DELETE("/delete", videoController.DeleteVideo)
	}
}

func setupCurriculumRoutes(r *gin.Engine) {
	curriculumDAO := dao.NewCurriculumDAO(database.DB)
	curriculumUsecase := usecase.NewCurriculumUsecase(curriculumDAO)
	curriculumController := controller.NewCurriculumController(curriculumUsecase)

	curriculumsGroup := r.Group("/curriculums")
	{
		curriculumsGroup.POST("/create", curriculumController.CreateCurriculum)
		curriculumsGroup.PUT("/update", curriculumController.UpdateCurriculum)
		curriculumsGroup.GET("/get", curriculumController.GetCurriculum)
		curriculumsGroup.GET("/showall", curriculumController.GetAllCurriculum)
		curriculumsGroup.DELETE("/delete", curriculumController.DeleteCurriculum)
	}
}
