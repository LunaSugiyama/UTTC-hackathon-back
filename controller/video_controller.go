package controller

import (
	"net/http"
	"strconv"
	"uttc-hackathon/model"
	"uttc-hackathon/usecase"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoUsecase usecase.VideoUsecase
}

func NewVideoController(videoUsecase usecase.VideoUsecase) *VideoController {
	return &VideoController{
		videoUsecase: videoUsecase,
	}
}

func (bc *VideoController) CreateVideo(c *gin.Context) {
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.videoUsecase.CreateVideo(&video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (bc *VideoController) GetVideo(c *gin.Context) {
	videoID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	video, err := bc.videoUsecase.GetVideo(videoID)
	if err != nil {
		// Handle different types of errors appropriately
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (bc *VideoController) UpdateVideo(c *gin.Context) {
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	video, err := bc.videoUsecase.UpdateVideo(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (bc *VideoController) DeleteVideo(c *gin.Context) {
	videoID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := bc.videoUsecase.DeleteVideo(videoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video entry deleted successfully"})
}

func (bc *VideoController) ShowAllVideos(c *gin.Context) {
	videos, err := bc.videoUsecase.ShowAllVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}
