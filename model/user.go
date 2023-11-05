// user.go

package model

import (
	"time"
)

// User represents the structure of the 'users' table in the database.
type User struct {
	ID          int       `json:"id"`
	FirebaseUID string    `json:"firebase_uid"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Age         int       `json:"age"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Authority   int       `json:"authority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
