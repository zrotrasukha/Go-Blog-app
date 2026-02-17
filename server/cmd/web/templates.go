package main

import "github.com/zrotrasukha/Go-Blog-app/internal/models"

type templateData struct {
	Blog  *models.Blog
	Blogs []*models.Blog
}
