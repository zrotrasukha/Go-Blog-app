package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("server/ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", home)
	mux.HandleFunc("/post/view", postView)
	mux.HandleFunc("/post/create", postCreate)

	log.Println("server is starting at :8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
