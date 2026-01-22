package main

import (
	"fmt"
	"log"
	"net/http"
)

func (p *Ping) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", p.homeHandler)
	mux.HandleFunc("/pings", p.getPings)
	mux.HandleFunc("/pingpong", p.pingHandler)
	return mux
}

func (p *Ping) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "visit /pingpong")
}

func (p *Ping) getPings(w http.ResponseWriter, r *http.Request) {
	count, err := p.pingdb.Get()
	if err != nil {
		serverError(w, r, err, "")
		return
	}
	fmt.Fprintf(w, "ping %d", count)
}

func (p *Ping) pingHandler(w http.ResponseWriter, r *http.Request) {
	count, err := p.pingdb.Increment()
	if err != nil {
		serverError(w, r, err, "")
		return
	}
	fmt.Fprintf(w, "ping %d", count)
}

func serverError(w http.ResponseWriter, r *http.Request, err error, text string) {
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
