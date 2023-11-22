package usecase

import (
	"errors"
	"time"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type CommentUsecase interface {
	CreateComment(comment *model.Comment) error
	GetComments(itemID int, itemCategoriesID int) ([]model.Comment, error)
	DeleteComment(comment model.Comment) error
	UpdateComment(comment *model.Comment) error
}

type commentUsecase struct {
	commentDAO dao.CommentDAO
}

func NewCommentUsecase(commentDAO dao.CommentDAO) CommentUsecase {
	return &commentUsecase{
		commentDAO: commentDAO,
	}
}

func (uc *commentUsecase) CreateComment(comment *model.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return uc.commentDAO.CreateComment(comment)
}

func (uc *commentUsecase) GetComments(itemID int, itemCategoriesID int) ([]model.Comment, error) {
	return uc.commentDAO.GetComments(itemID, itemCategoriesID)
}

func (uc *commentUsecase) DeleteComment(comment model.Comment) error {
	if comment.UserFirebaseUID == "" || comment.ItemID == 0 || comment.ItemCategoriesID == 0 || comment.Comment == "" {
		return errors.New("missing required parameters")
	}

	return uc.commentDAO.DeleteComment(comment)
}

func (uc *commentUsecase) UpdateComment(comment *model.Comment) error {
	comment.UpdatedAt = time.Now()
	return uc.commentDAO.UpdateComment(comment)
}
