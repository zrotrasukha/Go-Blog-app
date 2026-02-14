package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Everything is alright"))
}
func postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Viewing post with ID: %d", id)
}

func postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	w.Write([]byte("creating a post"))
}

func home(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
