package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Title   string
	Message string
	Image   string
}

const (
	pathname = "./ui/static/image.jpg"
	url      = "https://picsum.photos/1200"
)

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
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting Todo-App server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, img := GetImage(currentImg)
	data := PageData{
		Title:   "Hello from Todo-App",
		Message: "This is from Exercise: 1.12",
		Image:   strings.TrimPrefix(img.name, "./ui"),
	}

	tmpl, err := template.ParseFiles("./ui/index.tmpl")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
