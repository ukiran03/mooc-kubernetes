package main

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type PageData struct {
	Image      string
	TaskList   []Task
	BackendURL string
}

func Print(t *Task) string {
	var state string
	switch t.State {
	case StateTodo:
		state = "[ ]"
	case StateDone:
		state = "[X]"
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
