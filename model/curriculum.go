package model

import (
	"time"
)

// Curriculum represents the 'curriculums' table structure
type Curriculum struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// You can add methods to this struct for CRUD operations or any other functionality you need.
