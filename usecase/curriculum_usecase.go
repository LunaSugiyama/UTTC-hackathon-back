package usecase

import (
	"errors"
	"time"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type CurriculumUsecase interface {
	CreateCurriculum(curriculum *model.Curriculum) error
	GetCurriculum(id int) (model.Curriculum, error)
	UpdateCurriculum(curriculum *model.Curriculum) (model.Curriculum, error)
	DeleteCurriculum(id int) error
	ShowAllCurriculums() ([]model.Curriculum, error)
}

type curriculumUsecase struct {
	curriculumDAO dao.CurriculumDAO
}

func NewCurriculumUsecase(curriculumDAO dao.CurriculumDAO) CurriculumUsecase {
	return &curriculumUsecase{
		curriculumDAO: curriculumDAO,
	}
}

func (bu *curriculumUsecase) CreateCurriculum(curriculum *model.Curriculum) error {
	if curriculum.Name == "" {
		return errors.New("name is null")
	}
	// Business logic here
	curriculum.CreatedAt = time.Now()
	curriculum.UpdatedAt = time.Now()

	return bu.curriculumDAO.SaveCurriculum(curriculum)
}

func (bu *curriculumUsecase) GetCurriculum(id int) (model.Curriculum, error) {
	return bu.curriculumDAO.GetCurriculumByID(id)
}

func (bu *curriculumUsecase) UpdateCurriculum(curriculum *model.Curriculum) (model.Curriculum, error) {
	if curriculum.Name == "" {
		return model.Curriculum{}, errors.New("name is null")
	}

	// Business logic here
	curriculum.UpdatedAt = time.Now()

	return bu.curriculumDAO.UpdateCurriculum(curriculum)
}

func (bu *curriculumUsecase) DeleteCurriculum(id int) error {
	if id == 0 {
		return errors.New("id is null")
	}
	return bu.curriculumDAO.DeleteCurriculum(id)
}

func (bu *curriculumUsecase) ShowAllCurriculums() ([]model.Curriculum, error) {
	return bu.curriculumDAO.ShowAllCurriculums()
}
