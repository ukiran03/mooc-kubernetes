package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Image    string
	TaskList []Task
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

	_, img := GetImage(currentImg)
	data := PageData{
		Image:    strings.TrimPrefix(img.name, "./ui"),
		TaskList: demoTasks,
	}

	mux.HandleFunc("/{$}", data.homeHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting Todo-App server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
