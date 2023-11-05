package router

import (
	"uttc-hackathon/controller"

	"github.com/gin-gonic/gin"
)

func SetupBlogRoutes(router *gin.Engine) {
	// Create a new blog
	router.POST("/create", controller.CreateBlog)

	// // Get a blog by ID
	// router.GET("/blogs/:id", controller.GetBlog)

	// // Update a blog by ID
	// router.PUT("/blogs/:id", controller.UpdateBlog)

	// // Delete a blog by ID
	// router.DELETE("/blogs/:id", controller.DeleteBlog)
}
