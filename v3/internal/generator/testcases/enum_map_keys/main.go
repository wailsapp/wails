package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Status represents different status values
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Priority represents priority levels
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// Color represents color values
type Color uint8

const (
	Red   Color = 1
	Green Color = 2
	Blue  Color = 3
)

// EnumMapService tests various enum map key scenarios
type EnumMapService struct{}

// GetStatusMessages returns a map with string enum keys
func (*EnumMapService) GetStatusMessages() map[Status]string {
	return map[Status]string{
		StatusPending:   "Task is pending",
		StatusRunning:   "Task is running",
		StatusCompleted: "Task completed successfully",
		StatusFailed:    "Task failed",
	}
}

// GetPriorityWeights returns a map with integer enum keys
func (*EnumMapService) GetPriorityWeights() map[Priority]float64 {
	return map[Priority]float64{
		PriorityLow:    1.0,
		PriorityMedium: 2.5,
		PriorityHigh:   5.0,
	}
}

// GetColorCodes returns a map with uint8 enum keys
func (*EnumMapService) GetColorCodes() map[Color]string {
	return map[Color]string{
		Red:   "#FF0000",
		Green: "#00FF00",
		Blue:  "#0000FF",
	}
}

// GetNestedEnumMap returns a map with enum keys and complex values
func (*EnumMapService) GetNestedEnumMap() map[Status]map[Priority]string {
	return map[Status]map[Priority]string{
		StatusPending: {
			PriorityLow:    "Waiting in queue",
			PriorityMedium: "Scheduled soon",
			PriorityHigh:   "Next in line",
		},
		StatusRunning: {
			PriorityLow:    "Processing slowly",
			PriorityMedium: "Processing normally",
			PriorityHigh:   "Processing urgently",
		},
	}
}

// GetOptionalEnumMap returns a map with enum keys to optional values
func (*EnumMapService) GetOptionalEnumMap() map[Status]*string {
	running := "Currently running"
	return map[Status]*string{
		StatusPending:   nil,
		StatusRunning:   &running,
		StatusCompleted: nil,
		StatusFailed:    nil,
	}
}

// Person represents a person with status
type Person struct {
	Name   string
	Status Status
}

// GetPersonsByStatus returns a map with enum keys to struct values
func (*EnumMapService) GetPersonsByStatus() map[Status][]Person {
	return map[Status][]Person{
		StatusPending: {
			{Name: "Alice", Status: StatusPending},
			{Name: "Bob", Status: StatusPending},
		},
		StatusRunning: {
			{Name: "Charlie", Status: StatusRunning},
		},
		StatusCompleted: {
			{Name: "Dave", Status: StatusCompleted},
			{Name: "Eve", Status: StatusCompleted},
		},
	}
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&EnumMapService{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}