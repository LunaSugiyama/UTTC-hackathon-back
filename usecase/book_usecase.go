package usecase

import (
	"errors"
	"fmt"
	"time"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type BookUsecase interface {
	CreateBook(book *model.Book) error
	GetBook(id int) (model.Book, error)
	UpdateBook(book *model.Book) (model.Book, error)
	DeleteBook(id int) error
	ShowAllBooks() ([]model.Book, error)
}

type bookUsecase struct {
	bookDAO dao.BookDAO
}

func NewBookUsecase(bookDAO dao.BookDAO) BookUsecase {
	return &bookUsecase{
		bookDAO: bookDAO,
	}
}

func (bu *bookUsecase) CreateBook(book *model.Book) error {
	if book.UserFirebaseUID == "" {
		return errors.New("missing required parameters: user_firebase_uid")
	}
	if book.Title == "" {
		return errors.New("missing required parameters: title")
	}
	if book.Author == "" {
		return errors.New("missing required parameters: author")
	}
	if book.Link == "" {
		return errors.New("missing required parameters: link")
	}
	if book.ItemCategoriesID == 0 {
		return errors.New("missing required parameters: item_categories_id")
	}
	if book.Explanation == "" {
		return errors.New("missing required parameters: explanation")
	}
	// Business logic here
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	return bu.bookDAO.SaveBook(book)
}

func (bu *bookUsecase) GetBook(id int) (model.Book, error) {
	if id == 0 {
		return model.Book{}, errors.New("id is null")
	}
	fmt.Println("id", id)
	return bu.bookDAO.GetBookByID(id)
}

func (bu *bookUsecase) UpdateBook(book *model.Book) (model.Book, error) {
	if book.ID == 0 {
		return model.Book{}, errors.New("id is null")
	}
	if book.UserFirebaseUID == "" {
		return model.Book{}, errors.New("missing required parameters: user_firebase_uid")
	}
	if book.Title == "" {
		return model.Book{}, errors.New("missing required parameters: title")
	}
	if book.Author == "" {
		return model.Book{}, errors.New("missing required parameters: author")
	}
	if book.Link == "" {
		return model.Book{}, errors.New("missing required parameters: link")
	}
	if book.ItemCategoriesID == 0 {
		return model.Book{}, errors.New("missing required parameters: item_categories_id")
	}
	if book.Explanation == "" {
		return model.Book{}, errors.New("missing required parameters: explanation")
	}
	// Business logic here
	book.UpdatedAt = time.Now()

	return bu.bookDAO.UpdateBook(book)
}

func (bu *bookUsecase) DeleteBook(id int) error {
	if id == 0 {
		return errors.New("id is null")
	}
	return bu.bookDAO.DeleteBook(id)
}

func (bu *bookUsecase) ShowAllBooks() ([]model.Book, error) {
	return bu.bookDAO.ShowAllBooks()
}
