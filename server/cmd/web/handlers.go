package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/zrotrasukha/Go-Blog-app/internal/models"
	"github.com/zrotrasukha/Go-Blog-app/internal/validator"
)

type blogCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Author              string `form:"author"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData()
	data.Blog = b
	data.Flash = flash

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
	var form blogCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.NotBlank(form.Author), "author", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Author, 50), "author", "This field cannot be more than 50 characters long")
	form.CheckField(form.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be equal to 1, 7 or 365")

	if !form.Valid() {
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
	app.sessionManager.Put(r.Context(), "flash", "Blog post successfully created!")

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
