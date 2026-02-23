package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

type blogCreateForm struct {
	Title      string
	Content    string
	Author     string
	Expires    int
	FieldError map[string]string
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Everything is alright"))
}
func (app *application) blogView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	b, err := app.blog.Get(id)
	if err != nil {
		if err == models.ErrNoRecord {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Blog = b

	app.render(w, http.StatusOK, "view.html", data)
}

// for displaying the form to create a new blog
func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = blogCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

// for creating the post
func (app *application) blogCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.Form.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validation
	form := blogCreateForm{
		Title:      r.PostForm.Get("title"),
		Content:    r.PostForm.Get("content"),
		Author:     r.PostForm.Get("author"),
		Expires:    expires,
		FieldError: make(map[string]string),
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldError["title"] = "Title is required"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldError["title"] = "Title must not exceed 100 characters"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldError["content"] = "Content is required"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldError["author"] = "Author is required"
	} else if utf8.RuneCountInString(form.Author) > 50 {
		form.FieldError["author"] = "Author must not exceed 50 characters"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldError["expires"] = "Expires must be 1, 7, or 365"
	}

	if len(form.FieldError) > 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.blog.Insert(form.Title, form.Content, form.Author, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	blogs, err := app.blog.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData()
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.html", data)
}
