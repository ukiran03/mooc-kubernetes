package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type PageData struct {
	Image    string
	TaskList []Task
}

func (app *application) render(
	w http.ResponseWriter, r *http.Request,
	status int, data PageData,
) {
	ts := app.tmplCache
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		app.serverError(w, r, err)
		log.Printf("Execution error: %v", err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func Print(t *Task) string {
	var state string
	switch t.State {
	case StateTodo:
		state = "[X]"
	case StateDone:
		state = "[ ]"
	}
	return fmt.Sprintf("%s %s\n", state, t.Title)
}

var functions = template.FuncMap{
	"Print": Print,
}

func newTemplateCache() (*template.Template, error) {
	index := "./ui/html/index.tmpl"
	name := filepath.Base(index)
	ts := template.New(name).Funcs(functions)
	ts, err := ts.ParseFiles(index)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (app *application) serverError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
