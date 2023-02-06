package logger

import "time"

type Message struct {
	Level   string    `json:"log"`
	Message string    `json:"message"`
	Data    []any     `json:"data,omitempty"`
	Sender  string    `json:"-"`
	Time    time.Time `json:"-"`
}
