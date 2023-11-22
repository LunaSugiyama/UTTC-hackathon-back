// dao/user_dao.go

package dao

import (
	"fmt"
	"log"
	"uttc-hackathon/database"
	"uttc-hackathon/model"
	// SQLドライバやデータベース関連のインポート
)

type UserDao interface {
	SaveUser(user model.User) error
	UpdateUser(user model.User) error
	ShowUser(user *model.User) error
}

type userDao struct {
	// Struct fields, if any
}

func NewUserDao() UserDao {
	return &userDao{}
}

func (dao *userDao) SaveUser(user model.User) error {
	insertUserSQL := `
        INSERT INTO users (firebase_uid, username, email, name, age, created_at) VALUES (?, ?, ?, ?, ?, NOW())`

	_, err := database.DB.Exec(insertUserSQL, user.FirebaseUID, user.Email, user.Email, user.Name, user.Age)
	if err != nil {
		log.Printf("Error saving user to SQL database: %v", err)
		// You might want to return an error response here, or handle the error according to your application's logic.
		return err
	}

	// User successfully inserted into the database
	log.Printf("User inserted into the database.")
	return nil
}

func (dao *userDao) UpdateUser(user model.User) error {
	updateUserSQL := `
		UPDATE users SET name=?, email=?, age=?, password=? WHERE firebase_uid=?`

	_, err := database.DB.Exec(updateUserSQL, user.Name, user.Email, user.Age, user.Password, user.FirebaseUID)
	if err != nil {
		log.Printf("Error updating user in SQL database: %v", err)
		return err
	}

	log.Printf("User updated in the database.")
	return nil
}

func (dao *userDao) ShowUser(user *model.User) error {
	fmt.Println(user.FirebaseUID)
	showUserSQL := `
		SELECT firebase_uid, name, email, age FROM users WHERE firebase_uid = ?`

	err := database.DB.QueryRow(showUserSQL, user.FirebaseUID).Scan(&user.FirebaseUID, &user.Name, &user.Email, &user.Age)
	if err != nil {
		log.Printf("Error showing user in SQL database: %v", err)
		return err
	}

	log.Printf("User showed in the database.")
	return nil
}
