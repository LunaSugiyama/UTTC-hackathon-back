package usecase

import (
	"errors"
	"time"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type VideoUsecase interface {
	CreateVideo(video *model.Video) error
	GetVideo(id int) (model.Video, error)
	UpdateVideo(video *model.Video) (model.Video, error)
	DeleteVideo(id int) error
	ShowAllVideos() ([]model.Video, error)
}

type videoUsecase struct {
	videoDAO dao.VideoDAO
}

func NewVideoUsecase(videoDAO dao.VideoDAO) VideoUsecase {
	return &videoUsecase{
		videoDAO: videoDAO,
	}
}

func (bu *videoUsecase) CreateVideo(video *model.Video) error {
	if video.UserFirebaseUID == "" {
		return errors.New("missing required parameters: user_firebase_uid")
	}
	if video.Title == "" {
		return errors.New("missing required parameters: title")
	}
	if video.Author == "" {
		return errors.New("missing required parameters: author")
	}
	if video.Link == "" {
		return errors.New("missing required parameters: link")
	}
	if video.ItemCategoriesID == 0 {
		return errors.New("missing required parameters: item_categories_id")
	}
	if video.Explanation == "" {
		return errors.New("missing required parameters: explanation")
	}
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()

	return bu.videoDAO.SaveVideo(video)
}

func (bu *videoUsecase) GetVideo(id int) (model.Video, error) {
	if id == 0 {
		return model.Video{}, errors.New("id is null")
	}
	return bu.videoDAO.GetVideoByID(id)
}

func (bu *videoUsecase) UpdateVideo(video *model.Video) (model.Video, error) {
	if video.ID == 0 {
		return model.Video{}, errors.New("id is null")
	}
	if video.UserFirebaseUID == "" {
		return model.Video{}, errors.New("missing required parameters: user_firebase_uid")
	}
	if video.Title == "" {
		return model.Video{}, errors.New("missing required parameters: title")
	}
	if video.Author == "" {
		return model.Video{}, errors.New("missing required parameters: author")
	}
	if video.Link == "" {
		return model.Video{}, errors.New("missing required parameters: link")
	}
	if video.ItemCategoriesID == 0 {
		return model.Video{}, errors.New("missing required parameters: item_categories_id")
	}
	if video.Explanation == "" {
		return model.Video{}, errors.New("missing required parameters: explanation")
	}
	video.UpdatedAt = time.Now()

	return bu.videoDAO.UpdateVideo(video)
}

func (bu *videoUsecase) DeleteVideo(id int) error {
	if id == 0 {
		return errors.New("id is null")
	}
	return bu.videoDAO.DeleteVideo(id)
}

func (bu *videoUsecase) ShowAllVideos() ([]model.Video, error) {
	return bu.videoDAO.ShowAllVideos()
}
