package mocks

import (
	"time"

	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

var mockBlog = &models.Blog{
	ID:        1,
	Title:     "Test Blog",
	Content:   "This is a test blog content.",
	CreatedAt: time.Now(),
	Expires:   time.Now(),
}

type BlogModel struct{}

func (m *BlogModel) Insert(title, author, content string, expires int) (int, error) {
	return 2, nil
}

func (m *BlogModel) Get(id int) (*models.Blog, error) {
	switch id {
	case 1:
		return mockBlog, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *BlogModel) Latest() ([]*models.Blog, error) {
	return []*models.Blog{mockBlog}, nil
}
