package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	homeTemplate := resolveExistingPath("ui/html/home.html", "server/ui/html/home.html")
	if homeTemplate == "" {
		log.Println("home template not found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles(homeTemplate)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
