package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookUsecase usecase.BookUsecase
}

func NewBookController(bookUsecase usecase.BookUsecase) *BookController {
	return &BookController{
		bookUsecase: bookUsecase,
	}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.bookUsecase.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) GetBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	book, err := bc.bookUsecase.GetBook(bookID)
	if err != nil {
		// Handle different types of errors appropriately
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	book, err := bc.bookUsecase.UpdateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := bc.bookUsecase.DeleteBook(bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book entry deleted successfully"})
}

func (bc *BookController) ShowAllBooks(c *gin.Context) {
	books, err := bc.bookUsecase.ShowAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
