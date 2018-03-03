package main

import "time"

type Todo struct {
	Isbn  int    `json:"isbn"`
	Title string `json:"title"`

	DeadLine    string `json:"deadline"`
	TimeCreated time.Time
}

type Todos []Todo
