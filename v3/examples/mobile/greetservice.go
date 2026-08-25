package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// SystemService is a small service the frontend calls over the Wails bindings
// to demonstrate JS -> Go calls returning values, structs and errors.
type SystemService struct{}

// Greet returns a greeting for the given name (empty -> "anonymous").
func (s *SystemService) Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "anonymous"
	}
	return "Hello " + name + " 👋"
}

// Add returns the sum of two integers (demonstrates typed args/returns).
func (s *SystemService) Add(a int, b int) int {
	return a + b
}

// SystemInfo is returned by Info() and serialised to the frontend as JSON.
type SystemInfo struct {
	GoVersion  string `json:"goVersion"`
	GOOS       string `json:"goos"`
	GOARCH     string `json:"goarch"`
	NumCPU     int    `json:"numCPU"`
	ServerTime string `json:"serverTime"`
}

// Info returns information about the Go runtime hosting the app.
func (s *SystemService) Info() SystemInfo {
	return SystemInfo{
		GoVersion:  runtime.Version(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
		NumCPU:     runtime.NumCPU(),
		ServerTime: time.Now().Format(time.RFC1123),
	}
}

// Divide demonstrates a Go method returning an error to the frontend.
func (s *SystemService) Divide(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %v by zero", a)
	}
	return a / b, nil
}
