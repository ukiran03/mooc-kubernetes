package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	pathname = "./ui/static/image.jpg"
	url      = "https://picsum.photos/1200"
)

type application struct {
	image     string
	taskdb    *TaskModel
	tmplCache *template.Template
}

var currentImg = &Image{}

func init() {
	info, err := os.Stat(pathname)
	if err == nil {
		currentImg.name = pathname
		currentImg.modTime = info.ModTime()
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000")
		port = "3000"
	}
	addr := ":" + port

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmplCache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	_, img := GetImage(currentImg)
	app := &application{
		image:     strings.TrimPrefix(img.name, "./ui"),
		taskdb:    &TaskModel{DB: db},
		tmplCache: tmplCache,
	}
	log.Printf("Starting Todo-App server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("POST /create", app.createTask)
	return mux
}
