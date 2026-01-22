package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

func loadConfig() (*Config, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	if host == "" {
		host = "postgres-svc"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" || passwd == "" || dbName == "" {
		return nil, fmt.Errorf(
			`missing required environment variables (DB_USER, DB_PASSWD, or DB_NAME)`)
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName,
	)
	return &Config{DSN: dsn}, nil
}

func openDB() (*sql.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", cfg.DSN)
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

func initializePingsTable(db *sql.DB) error {
	// 1. Create table
	stmt := `CREATE TABLE IF NOT EXISTS pings (val INTEGER)`
	if _, err := db.Exec(stmt); err != nil {
		return err
	}

	// 2. Check if it's empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pings").Scan(&count)
	if err != nil {
		return err
	}

	// 3. If empty, insert the first "0"
	if count == 0 {
		_, err = db.Exec("INSERT INTO pings (val) VALUES (0)")
		if err != nil {
			return err
		}
		log.Println("Seeded database with initial row.")
	}

	return nil
}
