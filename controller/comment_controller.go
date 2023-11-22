package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentUsecase usecase.CommentUsecase
}

func NewCommentController(commentUsecase usecase.CommentUsecase) *CommentController {
	return &CommentController{
		commentUsecase: commentUsecase,
	}
}

func (ctrl *CommentController) CreateComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.commentUsecase.CreateComment(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (ctrl *CommentController) GetComments(c *gin.Context) {
	item_id := c.Query("item_id")
	item_categories_id := c.Query("item_categories_id")

	item_id_int, err := strconv.Atoi(item_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ItemID parameter"})
		return
	}

	item_categories_id_int, err := strconv.Atoi(item_categories_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ItemCategoriesID parameter"})
		return
	}

	comments, err := ctrl.commentUsecase.GetComments(item_id_int, item_categories_id_int)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.commentUsecase.DeleteComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted comment"})
}

func (ctrl *CommentController) UpdateComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := ctrl.commentUsecase.UpdateComment(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
