package usecase

import (
	"errors"
	"fmt"
	"time"
	"uttc-hackathon/dao"
	"uttc-hackathon/model"
)

type BlogUsecase interface {
	CreateBlog(blog *model.Blog) error
	GetBlog(id int) (model.Blog, error)
	UpdateBlog(blog *model.Blog) (model.Blog, error)
	DeleteBlog(id int) error
	ShowAllBlogs() ([]model.Blog, error)
}

type blogUsecase struct {
	blogDAO dao.BlogDAO
}

func NewBlogUsecase(blogDAO dao.BlogDAO) BlogUsecase {
	return &blogUsecase{
		blogDAO: blogDAO,
	}
}

func (bu *blogUsecase) CreateBlog(blog *model.Blog) error {
	if blog.UserFirebaseUID == "" {
		return errors.New("missing required parameters: user_firebase_uid")
	}
	if blog.Title == "" {
		return errors.New("missing required parameters: title")
	}
	if blog.Author == "" {
		return errors.New("missing required parameters: author")
	}
	if blog.Link == "" {
		return errors.New("missing required parameters: link")
	}
	if blog.ItemCategoriesID == 0 {
		return errors.New("missing required parameters: item_categories_id")
	}
	if blog.Explanation == "" {
		return errors.New("missing required parameters: explanation")
	}
	// Business logic here
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	return bu.blogDAO.SaveBlog(blog)
}

func (bu *blogUsecase) GetBlog(id int) (model.Blog, error) {
	if id == 0 {
		return model.Blog{}, errors.New("id is null")
	}
	return bu.blogDAO.GetBlogByID(id)
}

func (bu *blogUsecase) UpdateBlog(blog *model.Blog) (model.Blog, error) {
	if blog.ID == 0 {
		return model.Blog{}, errors.New("id is null")
	}
	if blog.UserFirebaseUID == "" {
		return model.Blog{}, errors.New("missing required parameters: user_firebase_uid")
	}
	if blog.Title == "" {
		return model.Blog{}, errors.New("missing required parameters: title")
	}
	if blog.Author == "" {
		return model.Blog{}, errors.New("missing required parameters: author")
	}
	if blog.Link == "" {
		return model.Blog{}, errors.New("missing required parameters: link")
	}
	if blog.ItemCategoriesID == 0 {
		return model.Blog{}, errors.New("missing required parameters: item_categories_id")
	}
	if blog.Explanation == "" {
		return model.Blog{}, errors.New("missing required parameters: explanation")
	}
	fmt.Println(blog)

	// Business logic here
	blog.UpdatedAt = time.Now()

	return bu.blogDAO.UpdateBlog(blog)
}

func (bu *blogUsecase) DeleteBlog(id int) error {
	if id == 0 {
		return errors.New("id is null")
	}
	return bu.blogDAO.DeleteBlog(id)
}

func (bu *blogUsecase) ShowAllBlogs() ([]model.Blog, error) {
	return bu.blogDAO.ShowAllBlogs()
}
