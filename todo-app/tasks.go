package main

import (
	"time"
)

type State int

const (
	StateTodo State = iota
	// StateInProgress
	StateDone
)

type Task struct {
	Title     string
	State     State
	CreatedAt time.Time
	UpdatedAt time.Time
}

type List []Task

var demoTasks List = []Task{
	{"Complete Mooc Kubernetes before jan-31st", 0, time.Now(), time.Now()},
	{
		"Complete Mooc Docker before jan-11st", 1,
		time.Date(2025, time.December, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC),
	},
	{"Wash clothes", 0, time.Now(), time.Now()},
}
