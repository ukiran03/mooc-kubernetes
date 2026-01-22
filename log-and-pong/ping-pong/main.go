package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Ping struct {
	pingdb *PingModel
}

func main() {
	pingPort := os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}
	addr := ":" + pingPort

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = initializePingsTable(db)
	if err != nil {
		log.Printf("Error createTable: %v", err)
		return
	}

	srv := &Ping{
		pingdb: &PingModel{DB: db},
	}

	log.Printf("Server starting on %s\n", addr)
	log.Fatal((http.ListenAndServe(addr, srv.routes())))
}
