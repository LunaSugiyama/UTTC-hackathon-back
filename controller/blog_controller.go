package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUsecase usecase.BlogUsecase
}

func NewBlogController(blogUsecase usecase.BlogUsecase) *BlogController {
	return &BlogController{
		blogUsecase: blogUsecase,
	}
}

func (bc *BlogController) CreateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.blogUsecase.CreateBlog(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (bc *BlogController) GetBlog(c *gin.Context) {
	blogID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	blog, err := bc.blogUsecase.GetBlog(blogID)
	if err != nil {
		// Handle different types of errors appropriately
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (bc *BlogController) UpdateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog, err := bc.blogUsecase.UpdateBlog(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (bc *BlogController) DeleteBlog(c *gin.Context) {
	blogID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := bc.blogUsecase.DeleteBlog(blogID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog entry deleted successfully"})
}

func (bc *BlogController) ShowAllBlogs(c *gin.Context) {
	blogs, err := bc.blogUsecase.ShowAllBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blogs)
}
