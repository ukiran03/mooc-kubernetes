package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Title   string
	Message string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: env PORT was unset")
	}

	addr := ":" + port
	http.HandleFunc("/", helloHandler)

	log.Printf("Starting Todo-App server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Hello from Todo-App",
		Message: "This is from Exercise: 1.5",
	}

	tmpl, err := template.ParseFiles("./ui/index.tmpl")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
