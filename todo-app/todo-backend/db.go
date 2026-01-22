package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/rqlite/gorqlite/stdlib"
)

type Config struct {
	DSN string
}

func loadConfig() (*Config, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	// Format: "http://admin:supersecret@sqlite-service:4001"
	if host == "" || port == "" || user == "" || passwd == "" {
		return nil, fmt.Errorf(
			`missing required environment variables like
             (DB_HOST, DB_USER, DB_PASSWD, DB_PATH)`)
	}
	dsn := fmt.Sprintf("http://%s:%s@%s:%s", user, passwd, host, port)
	return &Config{DSN: dsn}, nil
}

func openDB() (*sql.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("rqlite", cfg.DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
