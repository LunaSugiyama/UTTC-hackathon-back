package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"uttc-hackathon/model"

	"github.com/go-sql-driver/mysql"
)

type CurriculumDAO interface {
	SaveCurriculum(curriculum *model.Curriculum) error
	GetCurriculumByID(id int) (model.Curriculum, error)
	UpdateCurriculum(curriculum *model.Curriculum) (model.Curriculum, error)
	DeleteCurriculum(id int) error
	ShowAllCurriculums() ([]model.Curriculum, error)
}

type curriculumDAO struct {
	db *sql.DB
}

func NewCurriculumDAO(db *sql.DB) CurriculumDAO {
	return &curriculumDAO{
		db: db,
	}
}

func (cd *curriculumDAO) SaveCurriculum(curriculum *model.Curriculum) error {
	query := fmt.Sprintf("INSERT INTO curriculums ( name) VALUES ('%s')", curriculum.Name)
	result, err := cd.db.Exec(query)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	curriculum.ID = int(lastInsertID) // Convert lastInsertID from int64 to int
	return nil
}

func (cd *curriculumDAO) GetCurriculumByID(id int) (model.Curriculum, error) {
	var curriculum model.Curriculum
	var createdAt, updatedAt mysql.NullTime
	query := fmt.Sprintf("SELECT * FROM curriculums WHERE id = %d", id)
	err := cd.db.QueryRow(query).Scan(&curriculum.ID, &curriculum.Name, &createdAt, &updatedAt)
	if err != nil {
		return curriculum, err
	}

	curriculum.CreatedAt = createdAt.Time
	curriculum.UpdatedAt = updatedAt.Time
	return curriculum, nil
}

func (cd *curriculumDAO) UpdateCurriculum(curriculum *model.Curriculum) (model.Curriculum, error) {
	checkQuery := "SELECT id FROM curriculums WHERE id = ?"
	err := cd.db.QueryRow(checkQuery, curriculum.ID).Scan(&curriculum.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Curriculum{}, fmt.Errorf("curriculum not found")
		}
		return model.Curriculum{}, err
	}

	updateQuery := "UPDATE curriculums SET name = ? WHERE id = ?"
	result, err := cd.db.Exec(updateQuery, curriculum.Name, curriculum.ID)
	if err != nil {
		return model.Curriculum{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Curriculum{}, err
	}

	if rowsAffected == 0 {
		return model.Curriculum{}, errors.New("no rows affected")
	}

	selectQuery := "SELECT * FROM curriculums WHERE id = ?"
	var createdAt, updatedAt mysql.NullTime
	err = cd.db.QueryRow(selectQuery, curriculum.ID).Scan(&curriculum.ID, &curriculum.Name, &createdAt, &updatedAt)
	if err != nil {
		return model.Curriculum{}, err
	}

	curriculum.CreatedAt = createdAt.Time
	curriculum.UpdatedAt = updatedAt.Time

	return *curriculum, nil
}

func (cd *curriculumDAO) DeleteCurriculum(id int) error {
	checkQuery := "SELECT id FROM curriculums WHERE id = ?"
	err := cd.db.QueryRow(checkQuery, id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("curriculum not found")
		}
		return err
	}

	query := fmt.Sprintf("DELETE FROM curriculums WHERE id = %d", id)
	_, err = cd.db.Exec(query)
	return err
}

func (cd *curriculumDAO) ShowAllCurriculums() ([]model.Curriculum, error) {
	var curriculums []model.Curriculum
	query := "SELECT * FROM curriculums"
	rows, err := cd.db.Query(query)
	if err != nil {
		return curriculums, err
	}
	defer rows.Close()

	for rows.Next() {
		var curriculum model.Curriculum
		var createdAt, updatedAt mysql.NullTime
		err := rows.Scan(&curriculum.ID, &curriculum.Name, &createdAt, &updatedAt)
		if err != nil {
			return curriculums, err
		}
		curriculum.CreatedAt = createdAt.Time
		curriculum.UpdatedAt = updatedAt.Time
		curriculums = append(curriculums, curriculum)
	}
	return curriculums, nil
}
