package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var defaultOutFile = "data.txt"

type Ping struct {
	outFile string
	counter int
}

func (p *Ping) incrementAndSave() int {
	p.counter++
	data := []byte(strconv.Itoa(p.counter))
	err := os.WriteFile(p.outFile, data, 0o644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
	return p.counter
}

func main() {
	outFile := os.Getenv("OUTFILE")
	if outFile == "" {
		outFile = defaultOutFile
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3001")
		port = "3001"
	}

	f, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	srv := &Ping{outFile: outFile, counter: 0}

	http.HandleFunc("/pingpong", srv.pingHandler)

	addr := ":" + port
	log.Printf("Server starting on %s\n", addr)
	log.Fatal((http.ListenAndServe(addr, nil)))
}

func (p *Ping) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Fprintf(w, "ping %d", p.incrementAndSave())
}
