package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

type TestService struct{}

func (t *TestService) TriggerPanic() string {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		var ptr *time.Time
		_ = ptr.Unix()
	}()
	return "Panic triggered in goroutine - check if app crashes or recovers"
}

func (t *TestService) TriggerImmediatePanic() string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from immediate panic: %v", err)
		}
	}()
	var ptr *time.Time
	_ = ptr.Unix()
	return "This should not be reached"
}

func (t *TestService) SafeMethod() string {
	return fmt.Sprintf("Safe method called at %s", time.Now().Format(time.RFC3339))
}

func main() {
	app := application.New(application.Options{
		Name:        "Signal Handler Test (#3965)",
		Description: "Test for signal handler crash on Linux",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(&TestService{}),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Signal Handler Test (#3965)",
		Width:  600,
		Height: 400,
		URL:    "/",
	})

	log.Println("Starting application...")
	log.Println("Click 'Trigger Panic' to test - app should NOT crash")
	log.Println("Before fix: app crashes with 'non-Go code set up signal handler without SA_ONSTACK flag'")
	log.Println("After fix: panic is recovered and logged")

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
