package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	staticDir := resolveExistingPath("ui/static", "server/ui/static")
	if staticDir == "" {
		log.Fatal("static directory not found")
	}
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/post/view", postView)
	mux.HandleFunc("/post/create", postCreate)
	mux.HandleFunc("/heaelth", health)

	log.Println("server is starting at :8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)

}
