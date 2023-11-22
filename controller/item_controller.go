package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	itemUsecase usecase.ItemUsecase
}

func NewItemController(itemUsecase usecase.ItemUsecase) *ItemController {
	return &ItemController{
		itemUsecase: itemUsecase,
	}
}

func (ctrl *ItemController) SearchItems(c *gin.Context) {
	var searchRequest usecase.SearchRequest
	if err := c.ShouldBindJSON(&searchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := ctrl.itemUsecase.SearchItems(searchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (ctrl *ItemController) LikeItem(c *gin.Context) {
	var likedItem model.LikedItem
	if err := c.ShouldBindJSON(&likedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.itemUsecase.LikeItem(&likedItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked item"})
}

func (ctrl *ItemController) UnlikeItem(c *gin.Context) {
	var likedItem model.LikedItem
	if err := c.ShouldBindJSON(&likedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.itemUsecase.UnlikeItem(&likedItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unliked item"})
}

func (ctrl *ItemController) StarItem(c *gin.Context) {
	var starredItem model.StarredItem
	if err := c.ShouldBindJSON(&starredItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.itemUsecase.StarItem(&starredItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully starred item"})
}

func (ctrl *ItemController) UnstarItem(c *gin.Context) {
	var starredItem model.StarredItem
	if err := c.ShouldBindJSON(&starredItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.itemUsecase.UnstarItem(&starredItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unstar item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unstarred item"})
}

func (ctrl *ItemController) CheckLiked(c *gin.Context) {
	var likedItem model.LikedItem

	itemIDStr := c.Query("item_id")
	fmt.Println(itemIDStr)
	itemID, err := strconv.Atoi(itemIDStr)
	fmt.Println(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
		return
	}

	itemCategoriesIDStr := c.Query("item_categories_id")
	itemCategoriesID, err := strconv.Atoi(itemCategoriesIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_categories_id"})
		return
	}

	likedItem.ItemID = itemID
	likedItem.ItemCategoriesID = itemCategoriesID
	likedItem.UserFirebaseUID = c.Query("user_firebase_uid")

	isLiked, err := ctrl.itemUsecase.CheckLiked(&likedItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}

func (ctrl *ItemController) CheckStarred(c *gin.Context) {
	var starredItem model.StarredItem
	starredItem.UserFirebaseUID = c.Query("user_firebase_uid")

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
		return
	}

	itemCategoriesIDStr := c.Query("item_categories_id")
	itemCategoriesID, err := strconv.Atoi(itemCategoriesIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_categories_id"})
		return
	}

	starredItem.ItemID = itemID
	starredItem.ItemCategoriesID = itemCategoriesID

	isStarred, err := ctrl.itemUsecase.CheckStarred(&starredItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isStarred": isStarred})
}

func (ctrl *ItemController) CountLikes(c *gin.Context) {
	var itemID int
	var itemCategoriesID int

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
		return
	}

	itemCategoriesIDStr := c.Query("item_categories_id")
	itemCategoriesID, err = strconv.Atoi(itemCategoriesIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_categories_id"})
		return
	}

	count, err := ctrl.itemUsecase.CountLikes(itemID, itemCategoriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
