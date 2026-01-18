package main

import (
	"fmt"
	"html/template"
)

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
