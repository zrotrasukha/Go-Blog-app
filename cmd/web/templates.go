package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/zrotrasukha/Go-Blog-app/internal/models"
	"github.com/zrotrasukha/Go-Blog-app/ui"
)

type templateData struct {
	Blog         *models.Blog
	Blogs        []*models.Blog
	CurrentYear  int
	Form         any
	Flash        string
	Autheticated bool
	CSRFToken    string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/pages/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
