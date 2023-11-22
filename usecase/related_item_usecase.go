package usecase

import (
	"errors"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type RelatedItemUsecase interface {
	GetItemsSimilarity(itemID int, itemCategoriesName string) ([]model.Item, error)
}

type relatedItemUsecase struct {
	relatedItemDAO dao.RelatedItemDAO
}

func NewRelatedItemUsecase(relatedItemDAO dao.RelatedItemDAO) RelatedItemUsecase {
	return &relatedItemUsecase{
		relatedItemDAO: relatedItemDAO,
	}
}

func (riu *relatedItemUsecase) GetItemsSimilarity(itemID int, itemCategoriesName string) ([]model.Item, error) {
	if itemID == 0 {
		return nil, errors.New("item_id is null")
	}
	if itemCategoriesName == "" {
		return nil, errors.New("item_categories_name is null")
	}

	return riu.relatedItemDAO.GetItemsSimilarity(itemID, itemCategoriesName)
}
