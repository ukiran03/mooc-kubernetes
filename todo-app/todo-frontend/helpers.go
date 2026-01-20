package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func (f *frontend) backendError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	http.Error(
		w,
		fmt.Sprintf("Backed Error: %v", err),
		http.StatusInternalServerError,
	)
}

func (f *frontend) serverError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (f *frontend) render(
	w http.ResponseWriter, r *http.Request, status int, data PageData,
) {
	ts := f.tmplCache
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		f.serverError(w, r, err)
		log.Printf("Execution error: %v", err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
