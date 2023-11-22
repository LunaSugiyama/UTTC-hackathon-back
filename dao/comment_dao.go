package dao

import (
	"database/sql"
	"errors"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/go-sql-driver/mysql"
)

type CommentDAO interface {
	CreateComment(comment *model.Comment) error
	GetComments(itemID int, itemCategoriesID int) ([]model.Comment, error)
	DeleteComment(comment model.Comment) error
	UpdateComment(comment *model.Comment) error
}

type commentDAO struct {
	db *sql.DB
}

func NewCommentDAO(db *sql.DB) CommentDAO {
	return &commentDAO{
		db: db,
	}
}

func (dao *commentDAO) CreateComment(comment *model.Comment) error {
	result, err := dao.db.Exec("INSERT INTO item_comments (item_id, item_categories_id, user_firebase_uid, comment) VALUES (?, ?, ?, ?)",
		comment.ItemID, comment.ItemCategoriesID, comment.UserFirebaseUID, comment.Comment)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(lastInsertID)
	return nil
}

func (dao *commentDAO) GetComments(itemID int, itemCategoriesID int) ([]model.Comment, error) {
	selectQuery := "SELECT * FROM item_comments WHERE item_id = ? AND item_categories_id = ?"
	rows, err := database.DB.Query(selectQuery, itemID, itemCategoriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Close the rows after we're done with them

	createdAt := mysql.NullTime{}
	updatedAt := mysql.NullTime{}

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.ID, &comment.UserFirebaseUID, &comment.ItemID, &comment.ItemCategoriesID, &comment.Comment, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAt.Time
		comment.UpdatedAt = updatedAt.Time
		comments = append(comments, comment)
	}

	return comments, nil
}

func (dao *commentDAO) DeleteComment(comment model.Comment) error {
	deleteQuery := "DELETE FROM item_comments WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ? AND comment = ?"
	result, err := database.DB.Exec(deleteQuery, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID, comment.Comment)
	if err != nil {
		return err
	}

	// Check the number of rows affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found or not authorized to delete")
	}
	return nil
}

func (dao *commentDAO) UpdateComment(comment *model.Comment) error {
	var updatedAt mysql.NullTime

	updateQuery := "UPDATE item_comments SET comment = ?, updated_at = ? WHERE user_firebase_uid = ? AND item_id = ? AND item_categories_id = ?"
	_, err := database.DB.Exec(updateQuery, comment.Comment, updatedAt, comment.UserFirebaseUID, comment.ItemID, comment.ItemCategoriesID)
	if err != nil {
		return err
	}

	// Assign the converted time values to the comment struct
	comment.UpdatedAt = updatedAt.Time
	return nil
}
