package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.taskdb.GetTasks()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := PageData{
		Image:    app.image,
		TaskList: tasks,
	}
	app.render(w, r, http.StatusOK, data)
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	// TODO:
}
