package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Everything is alright"))
}
func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
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

	app.render(w, http.StatusOK, "view.html", &templateData{
		Blog: b,
	})

}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail Climb Mount Fuji, But slowly, slowly! O snail Climb Mount Fuji, And be happy!"
	author := "Yosa Buson"
	expires := 7

	id, err := app.blog.Insert(title, content, author, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	blogs, err := app.blog.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.html", &templateData{
		Blogs: blogs,
	})
}
