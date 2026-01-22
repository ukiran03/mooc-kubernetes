package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

func loadConfig() (*Config, error) {
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	if user == "" || passwd == "" || dbName == "" {
		return nil, fmt.Errorf(
			"missing required environment variables (DB_USER, DB_PASSWD, or DB_NAME)",
		)
	}
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		user, passwd, dbName,
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

func createTable(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS pings (val INTEGER)`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
