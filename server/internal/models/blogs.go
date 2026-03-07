package models

import (
	"database/sql"
	"errors"
	"time"
)

type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Expires   time.Time `json:"expires"`
}

type BlogModelInterface interface {
	Insert(title string, content string, author string, expires int) (int, error)
	Get(id int) (*Blog, error)
	Latest() ([]*Blog, error)
}

type BlogModel struct {
	DB *sql.DB
}

func (b *BlogModel) Insert(title string, content string, author string, expires int) (int, error) {
	query := `INSERT INTO blogs ( title, content, author,  expires)
			  VALUES ( $1, $2, $3, NOW() + ($4 * INTERVAL '1 day')) returning id`

	var id int

	err := b.DB.QueryRow(query, title, content, author, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (b *BlogModel) Get(id int) (*Blog, error) {
	query := `SELECT id, title, content, author, created_at, expires FROM blogs
			  WHERE expires > NOW() AND id = $1`

	blog := &Blog{}
	err := b.DB.QueryRow(query, id).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return blog, nil
}

// it will return 10 most recent blogs
func (b *BlogModel) Latest() ([]*Blog, error) {
	query := `SELECT id, title, content, author, created_at, expires FROM blogs
	          WHERE expires > NOW() 
			  ORDER BY created_at DESC`

	rows, err := b.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog

	for rows.Next() {
		var b Blog
		if err := rows.Scan(&b.ID, &b.Title, &b.Content, &b.Author, &b.CreatedAt, &b.Expires); err != nil {
			return nil, err
		}
		blogs = append(blogs, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}
