package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	pathname = "./ui/static/image.jpg"
	imageUrl = "https://picsum.photos/1200"
)

type frontend struct {
	image     string
	tmplCache *template.Template
}

var currentImg = &Image{}

var backendUrl string

func init() {
	info, err := os.Stat(pathname)
	if err == nil {
		currentImg.name = pathname
		currentImg.modTime = info.ModTime()
	}
	backendEnvUrl := os.Getenv("BACKEND_URL")
	if backendEnvUrl == "" {
		backendUrl = "http://localhost:3000/tasks"
	} else {
		backendUrl = backendEnvUrl + "/tasks"
	}
}

func main() {
	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		fmt.Println("env FRONTEND_PORT was unset\nUsing Port 3001 as Frontend_Port")
		port = "3001"
	}
	addr := ":" + port

	tmplCache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	_, img := GetImage(currentImg)

	f := &frontend{
		image:     strings.TrimPrefix(img.name, "./ui"),
		tmplCache: tmplCache,
	}

	log.Printf("Starting Todo-App Frontend server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, f.routes()))
}

func (f *frontend) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", f.homeHandler)
	return mux
}

func (f *frontend) homeHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := fetchTasksFromBackend(backendUrl)
	if err != nil {
		f.backendError(w, r, err)
		return
	}
	data := PageData{
		Image:      f.image,
		TaskList:   *tasks,
		BackendURL: backendUrl,
	}
	f.render(w, r, http.StatusOK, data)
}
