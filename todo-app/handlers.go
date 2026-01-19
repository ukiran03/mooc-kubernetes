package main

import (
	"net/http"
	"unicode/utf8"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.PostFormValue("title")

	switch length := utf8.RuneCountInString(title); {
	case length == 0:
		http.Error(w, "Title cannot be empty",
			http.StatusBadRequest)
		return
	case length > 100:
		http.Error(w, "Title is too long (maximum 100 characters)",
			http.StatusBadRequest)
		return
	}

	_, err := app.taskdb.Insert(title, StateTodo)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
