package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorlog *log.Logger
	infolog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":8000", "HTTP networ address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	app := &application{
		errorlog: errorLog,
		infolog:  infoLog,
	}

	fileServer := http.FileServer(http.Dir("server/ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/view", app.postView)
	mux.HandleFunc("/post/create", app.postCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("server is starting at %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
