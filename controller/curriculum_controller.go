package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type CurriculumController struct {
	curriculumUsecase usecase.CurriculumUsecase
}

func NewCurriculumController(curriculumUsecase usecase.CurriculumUsecase) *CurriculumController {
	return &CurriculumController{
		curriculumUsecase: curriculumUsecase,
	}
}

func (bc *CurriculumController) CreateCurriculum(c *gin.Context) {
	var curriculum model.Curriculum
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.curriculumUsecase.CreateCurriculum(&curriculum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (bc *CurriculumController) GetCurriculum(c *gin.Context) {
	curriculumID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	curriculum, err := bc.curriculumUsecase.GetCurriculum(curriculumID)
	if err != nil {
		// Handle different types of errors appropriately
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (bc *CurriculumController) UpdateCurriculum(c *gin.Context) {
	var curriculum model.Curriculum
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	curriculum, err := bc.curriculumUsecase.UpdateCurriculum(&curriculum)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (bc *CurriculumController) DeleteCurriculum(c *gin.Context) {
	curriculumID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := bc.curriculumUsecase.DeleteCurriculum(curriculumID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum deleted successfully"})
}

func (bc *CurriculumController) GetAllCurriculum(c *gin.Context) {
	curriculums, err := bc.curriculumUsecase.ShowAllCurriculums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curriculums)
}
