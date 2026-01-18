package main

import (
	"html/template"
	"log"
	"net/http"
)

func (p *PageData) homeHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := template.New("index.tmpl").Funcs(functions).ParseFiles(
		"./ui/html/index.tmpl",
	)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, p)
	if err != nil {
		log.Printf("Execution error: %v", err)
	}
}

func (p *PageData) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
}
