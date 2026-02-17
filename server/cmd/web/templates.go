package main

import (
	"html/template"
	"path/filepath"

	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

type templateData struct {
	Blog  *models.Blog
	Blogs []*models.Blog
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("ui/html/pages/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles("ui/html/page/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
