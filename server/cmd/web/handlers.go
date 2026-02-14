package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Everything is alright"))
}
func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Viewing post with ID: %d", id)
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("creating a post"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"server/ui/html/pages/base.html",
		"server/ui/html/pages/partials/nav.html",
		"server/ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
