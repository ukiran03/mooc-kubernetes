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
		b.logger.Error("failed to decode task JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	titleLen := utf8.RuneCountInString(t.Title)
	if titleLen == 0 || titleLen > 140 {
		b.logger.Warn("invalid task title length",
			"length", titleLen,
			"remote_addr", r.RemoteAddr,
		)
		http.Error(
			w,
			"Title must be between 1 and 140 characters",
			http.StatusBadRequest,
		)
		return
	}

	id, err := b.taskdb.Insert(t.Title, t.State)
	if err != nil {
		b.logger.Error("database insertion failed",
			"error", err,
			"title", t.Title,
		)
		b.serverError(w, r, err, "Could not save task")
		return
	}
	b.logger.Info("task created successfully", "id", id)
	w.WriteHeader(http.StatusCreated)
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
