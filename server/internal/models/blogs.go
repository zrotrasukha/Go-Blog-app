package models

import (
	"database/sql"
	"errors"
	"time"
)

type Blog struct {
	ID         int
	title      string
	content    string
	author     string
	created_at time.Time
	expires    time.Time
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
	err := b.DB.QueryRow(query, id).Scan(&blog.ID, &blog.title, &blog.content, &blog.author, &blog.created_at, &blog.expires)
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
		if err := rows.Scan(&b.ID, &b.title, &b.content, &b.author, &b.created_at, &b.expires); err != nil {
			return nil, err
		}
		blogs = append(blogs, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}
