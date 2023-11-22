package controller

import (
	"net/http"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type ItemCategoryController struct {
	itemCategoryUsecase usecase.ItemCategoryUsecase
}

func NewItemCategoryController(itemCategoryUsecase usecase.ItemCategoryUsecase) *ItemCategoryController {
	return &ItemCategoryController{
		itemCategoryUsecase: itemCategoryUsecase,
	}
}

func (ic *ItemCategoryController) ShowAllItemCategories(c *gin.Context) {
	itemCategories, err := ic.itemCategoryUsecase.ShowAllItemCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemCategories)
}
