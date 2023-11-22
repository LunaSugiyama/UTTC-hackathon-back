package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type RelatedItemController struct {
	relatedItemUsecase usecase.RelatedItemUsecase
}

func NewRelatedItemController(relatedItemUsecase usecase.RelatedItemUsecase) *RelatedItemController {
	return &RelatedItemController{
		relatedItemUsecase: relatedItemUsecase,
	}
}

func (ctrl *RelatedItemController) GetItemsSimilarity(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Query("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	itemCategoriesName := c.Query("item_categories_name")

	items, err := ctrl.relatedItemUsecase.GetItemsSimilarity(itemID, itemCategoriesName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
		return
	}

	c.JSON(http.StatusOK, items)
}
