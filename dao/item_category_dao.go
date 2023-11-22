package dao

import (
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/go-sql-driver/mysql"
)

type ItemCategoryDAO interface {
	ShowAllItemCategories() ([]model.ItemCategory, error)
}

type itemCategoryDAO struct {
}

func NewItemCategoryDAO() ItemCategoryDAO {
	return &itemCategoryDAO{}
}

func (ic *itemCategoryDAO) ShowAllItemCategories() ([]model.ItemCategory, error) {
	query := "SELECT * FROM item_categories"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item_categories []model.ItemCategory

	for rows.Next() {
		var item_category model.ItemCategory
		var createdAt, updatedAt mysql.NullTime
		if err := rows.Scan(&item_category.ID, &item_category.Name, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		item_category.CreatedAt = createdAt.Time
		item_category.UpdatedAt = updatedAt.Time
		item_categories = append(item_categories, item_category)
	}
	return item_categories, nil
}
