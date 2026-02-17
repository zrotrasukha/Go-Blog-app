package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("server/ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/blog/view", app.blogView)
	mux.HandleFunc("/blog/create", app.blogCreate)

	return mux
}
