package main

import (
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

type Event struct {
	Name string
	When time.Time
}

func (*Service) GetTime() time.Time {
	return time.Now()
}

func (*Service) GetEvent() Event {
	return Event{
		Name: "test",
		When: time.Now(),
	}
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&Service{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
