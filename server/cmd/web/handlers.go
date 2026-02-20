package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

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

	app.render(w, http.StatusOK, "create.html", data)
}

// for creating the post
func (app *application) blogCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	author := r.Form.Get("author")
	expires, err := strconv.Atoi(r.Form.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.blog.Insert(title, content, author, expires)
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
