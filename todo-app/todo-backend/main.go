package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type backend struct {
	logger *slog.Logger
	taskdb *TaskModel
}

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		fmt.Println("env BACKEND_PORT was unset\nUsing Port 3000 as Backend_Port")
		port = "3000"
	}
	addr := ":" + port
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createInitialseTable(db)
	if err != nil {
		log.Printf("Error createInitialseTable: %v", err)
		return
	}

	b := &backend{
		logger: logger,
		taskdb: &TaskModel{DB: db},
	}
	logger.Info("Starting Todo-App Backend", "address", addr)
	err = http.ListenAndServe(addr, b.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func (b *backend) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", b.getTasks)
	mux.HandleFunc("POST /tasks", b.createTask)

	// Wrap the entire mux with our CORS logic
	return b.enableCORS(mux)
}

func (b *backend) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Set the headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 2. Handle the Pre-flight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 3. Pass the request to the mux (where it will match GET or POST)
		next.ServeHTTP(w, r)
	})
}
