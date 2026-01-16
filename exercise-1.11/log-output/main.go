package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var defaultPongFile = "data.txt"

type LogEnt struct {
	timeStamp string
	randStr   string
}

func (l LogEnt) String() string {
	return fmt.Sprintf("%s %s\n", l.timeStamp, l.randStr)
}

type Data struct {
	file string
}

func randomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func main() {
	inputFile := os.Getenv("INFILE")
	if inputFile == "" {
		inputFile = defaultPongFile
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000")
		port = "3000"
	}
	data := &Data{file: inputFile}
	http.HandleFunc("/log", data.homeHandler)
	addr := ":" + port
	log.Printf("Server starting on %s\n", addr)
	log.Fatal((http.ListenAndServe(addr, nil)))
}

func (d *Data) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	count, err := os.ReadFile(d.file)
	if err != nil {
		log.Printf("Error reading file %s: %v", d.file, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	entry := LogEnt{
		timeStamp: time.Now().Format("2006-01-02 15:04:05"),
		randStr:   randomString(),
	}
	fmt.Fprintf(w, "%sPing / Pongs: %s", entry, count)
}
