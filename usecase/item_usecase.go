package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type ItemUsecase interface {
	SearchItems(searchRequest SearchRequest) ([]model.Item, error)
	LikeItem(likedItem *model.LikedItem) error
	UnlikeItem(likedItem *model.LikedItem) error
	StarItem(starredItem *model.StarredItem) error
	UnstarItem(starredItem *model.StarredItem) error
	CheckLiked(likedItem *model.LikedItem) (bool, error)
	CheckStarred(starredItem *model.StarredItem) (bool, error)
	CountLikes(itemID int, itemCategoriesID int) (int, error)
}

type itemUsecase struct {
	itemDAO dao.ItemDAO
}

func NewItemUsecase(itemDAO dao.ItemDAO) ItemUsecase {
	return &itemUsecase{
		itemDAO: itemDAO,
	}
}

type SearchRequest struct {
	Words          string `json:"words"`
	Sort           string `json:"sorting"`
	Order          string `json:"order"`
	ItemCategories string `json:"item_categories"`
	CurriculumIDs  string `json:"curriculum_ids"`
}

func (uc *itemUsecase) SearchItems(searchRequest SearchRequest) ([]model.Item, error) {
	var ItemCategories []int
	var CurriculumIDs []int
	var itemCategoriesSQL string
	var curriculumIDsSQL string
	var searchResults []model.Item

	if searchRequest.ItemCategories != "" {
		ItemCategories = parseIntList(searchRequest.ItemCategories)
		itemCategoriesSQL = intListToSQL(ItemCategories)
	}

	if searchRequest.CurriculumIDs != "" {
		CurriculumIDs = parseIntList(searchRequest.CurriculumIDs)
		curriculumIDsSQL = intListToSQL(CurriculumIDs)
	}

	terms := strings.Split(searchRequest.Words, " ")

	items, err := uc.itemDAO.GetItems(ItemCategories, CurriculumIDs, itemCategoriesSQL, curriculumIDsSQL)
	if err != nil {
		return nil, err
	}
	fmt.Println("item: ", items)

	for _, item := range items {
		for _, term := range terms {
			nameContainsTerm := strings.Contains(strings.ToLower(item.Title), strings.ToLower(term))
			explanationContainsTerm := strings.Contains(strings.ToLower(item.Explanation), strings.ToLower(term))

			if nameContainsTerm || explanationContainsTerm {
				searchResults = append(searchResults, item)
			}
		}
	}

	return searchResults, nil
}

func parseIntList(input string) []int {
	parts := strings.Split(input, ",")
	var result []int
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err == nil {
			result = append(result, val)
		}
	}
	return result
}

func intListToSQL(list []int) string {
	var parts []string
	for _, val := range list {
		parts = append(parts, strconv.Itoa(val))
	}
	return strings.Join(parts, ",")
}

func (uc *itemUsecase) LikeItem(likedItem *model.LikedItem) error {
	return uc.itemDAO.LikeItem(likedItem)
}

func (uc *itemUsecase) UnlikeItem(likedItem *model.LikedItem) error {
	return uc.itemDAO.UnlikeItem(likedItem)
}

func (uc *itemUsecase) StarItem(starredItem *model.StarredItem) error {
	return uc.itemDAO.StarItem(starredItem)
}

func (uc *itemUsecase) UnstarItem(starredItem *model.StarredItem) error {
	return uc.itemDAO.UnstarItem(starredItem)
}

func (uc *itemUsecase) CheckLiked(likedItem *model.LikedItem) (bool, error) {
	return uc.itemDAO.CheckLiked(likedItem)
}

func (uc *itemUsecase) CheckStarred(starredItem *model.StarredItem) (bool, error) {
	return uc.itemDAO.CheckStarred(starredItem)
}

func (uc *itemUsecase) CountLikes(itemID int, itemCategoriesID int) (int, error) {
	return uc.itemDAO.CountLikes(itemID, itemCategoriesID)
}
