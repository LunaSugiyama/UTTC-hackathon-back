package usecase

import (
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type ItemCategoryUsecase interface {
	ShowAllItemCategories() ([]model.ItemCategory, error)
}

type itemCategoryUsecase struct {
	itemCategoryDAO dao.ItemCategoryDAO
}

func (ic *itemCategoryUsecase) ShowAllItemCategories() ([]model.ItemCategory, error) {
	return ic.itemCategoryDAO.ShowAllItemCategories()
}
