package models

import "time"

type Log struct {
	Id        int64     `json:"id"`
	Module    string    `json:"module"`
	Operation string    `json:"operation"`
	Level     string    `json:"log_level"`
	Value     string    `json:"value"`
	Exploit   *string   `json:"exploit"`
	Attrs     string    `json:"attrs"`
	CreatedAt time.Time `json:"created_at"`
}

type Module struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type LogLevel struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
