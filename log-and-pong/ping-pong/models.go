package main

import (
	"database/sql"
)

type PingModel struct {
	DB *sql.DB
}

func (m *PingModel) Get() (int, error) {
	stmt := `SELECT val FROM pings LIMIT 1`
	var count int
	err := m.DB.QueryRow(stmt).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *PingModel) Increment() (int, error) {
	stmt := `UPDATE pings SET val = val + 1 RETURNING val`
	var count int
	err := m.DB.QueryRow(stmt).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
