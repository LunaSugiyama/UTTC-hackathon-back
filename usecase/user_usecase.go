package usecase

import (
	"errors"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type UserUsecase interface {
	RegisterUser(user model.User) error
	LoginUser(data model.LoginData) (string, error)
	UpdateUser(user model.User) error
	ShowUser(user *model.User) error
}

type userUsecase struct {
	userDao dao.UserDao
}

func NewUserUsecase(userDao dao.UserDao) UserUsecase {
	return &userUsecase{
		userDao: userDao,
	}
}

func (u *userUsecase) RegisterUser(user model.User) error {
	// ユーザー情報のバリデーション
	if user.FirebaseUID == "" {
		return errors.New("missing required parameters: firebase_uid")
	}
	if user.Name == "" {
		return errors.New("missing required parameters: name")
	}
	if user.Email == "" {
		return errors.New("missing required parameters: email")
	}
	if user.Age == 0 {
		return errors.New("missing required parameters: age")
	}

	// データベースにユーザー情報を保存
	return u.userDao.SaveUser(user)
}

func (u *userUsecase) LoginUser(data model.LoginData) (string, error) {
	// ユーザー情報のバリデーション
	if data.FirebaseUID == "" {
		return "", errors.New("missing required parameters: uid")
	}
	if data.IDToken == "" {
		return "", errors.New("missing required parameters: idToken")
	}

	return data.IDToken, nil
}

func (u *userUsecase) UpdateUser(user model.User) error {
	// ユーザー情報のバリデーション
	if user.FirebaseUID == "" {
		return errors.New("missing required parameters: firebase_uid")
	}
	if user.Name == "" {
		return errors.New("missing required parameters: name")
	}
	if user.Email == "" {
		return errors.New("missing required parameters: email")
	}
	if user.Age == 0 {
		return errors.New("missing required parameters: age")
	}

	// データベースのユーザー情報を更新
	return u.userDao.UpdateUser(user)
}

func (u *userUsecase) ShowUser(user *model.User) error {
	// ユーザー情報のバリデーション
	if user.FirebaseUID == "" {
		return errors.New("missing required parameters: firebase_uid")
	}

	// データベースのユーザー情報を取得
	return u.userDao.ShowUser(user)
}
