package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zrotrasukha/Go-Blog-app/internal/models"
)

type application struct {
	errorlog      *log.Logger
	infolog       *log.Logger
	blog          *models.BlogModel
	templateCache map[string]*template.Template
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	neon_dsn := os.Getenv("NEON_USER_DSN")
	if neon_dsn == "" {
		log.Fatal("NEON_DSN environment variable is not set")
	}

	addr := flag.String("addr", ":8000", "HTTP network address")
	dsn := flag.String("dsn", neon_dsn, "PostgreSQL data source name")
	// dsn := flag.String("dsn", "postgres://web:webp@localhost/blogs?sslmode=disable", "PostgreSQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorlog:      errorLog,
		infolog:       infoLog,
		blog:          &models.BlogModel{DB: db},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("server is starting at %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
