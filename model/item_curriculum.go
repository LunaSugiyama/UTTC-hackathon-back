package model

// ItemCurriculum represents the 'item_curriculums' table structure
type ItemCurriculum struct {
	ID               int
	ItemID           int
	ItemCategoriesID int
	CurriculumID     int
}

// You can add methods to this struct for CRUD operations or any other functionality you need.
