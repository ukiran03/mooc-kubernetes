package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

func (b *backend) getTasks(w http.ResponseWriter, r *http.Request) {
	data, err := b.taskdb.GetTasks()
	if err != nil {
		log.Fatal(err)
	}
	writeJSON(w, http.StatusOK, data)
}

func (b *backend) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var t Task
	if err := readJSON(w, r, &t); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch length := utf8.RuneCountInString(t.Title); {
	case length == 0:
		b.serverError(w, r, nil, "Title cannot be empty")
		return
	case length > 100:
		b.serverError(w, r, nil, "Title is too long (maximum 100 characters)")
		return
	}
	_, err := b.taskdb.Insert(t.Title, t.State)
	if err != nil {
		b.serverError(w, r, err, "")
		return
	}
}

func (b *backend) serverError(
	w http.ResponseWriter, r *http.Request, err error, text string,
) {
	errString := text
	if err != nil {
		if errString != "" {
			errString = fmt.Sprintf("%s: %v", text, err)
		} else {
			errString = err.Error()
		}
	}
	log.Printf("ERROR: %s", errString)
	http.Error(w, errString, http.StatusInternalServerError)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Disallow unknown fields

	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON: %v", err) // to debug
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
