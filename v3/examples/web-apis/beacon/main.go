package main

import (
	"embed"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

// ReceivedBeacon represents a beacon received by the server
type ReceivedBeacon struct {
	ID          int               `json:"id"`
	Timestamp   string            `json:"timestamp"`
	ContentType string            `json:"contentType"`
	Size        int               `json:"size"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
}

// BeaconService handles beacon reception and retrieval
type BeaconService struct {
	mu       sync.RWMutex
	beacons  []ReceivedBeacon
	counter  int
	notifier func()
}

func NewBeaconService() *BeaconService {
	return &BeaconService{
		beacons: make([]ReceivedBeacon, 0),
	}
}

// SetNotifier sets a callback to notify the frontend of new beacons
func (s *BeaconService) SetNotifier(fn func()) {
	s.notifier = fn
}

// GetBeacons returns all received beacons
func (s *BeaconService) GetBeacons() []ReceivedBeacon {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.beacons
}

// ClearBeacons removes all received beacons
func (s *BeaconService) ClearBeacons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.beacons = make([]ReceivedBeacon, 0)
	s.counter = 0
}

// addBeacon adds a new beacon (called from HTTP handler)
func (s *BeaconService) addBeacon(contentType string, body []byte, headers map[string]string) {
	s.mu.Lock()
	s.counter++
	beacon := ReceivedBeacon{
		ID:          s.counter,
		Timestamp:   time.Now().Format("15:04:05.000"),
		ContentType: contentType,
		Size:        len(body),
		Body:        string(body),
		Headers:     headers,
	}
	s.beacons = append(s.beacons, beacon)
	s.mu.Unlock()

	if s.notifier != nil {
		s.notifier()
	}
}

// GetServerPort returns the beacon server port
func (s *BeaconService) GetServerPort() int {
	return 9999
}

func main() {
	beaconService := NewBeaconService()

	// Start HTTP server for receiving beacons
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/beacon", func(w http.ResponseWriter, r *http.Request) {
			// Handle CORS preflight
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			headers := make(map[string]string)
			for k, v := range r.Header {
				if len(v) > 0 {
					headers[k] = v[0]
				}
			}

			beaconService.addBeacon(r.Header.Get("Content-Type"), body, headers)
			w.WriteHeader(http.StatusOK)
		})

		http.ListenAndServe(":9999", mux)
	}()

	app := application.New(application.Options{
		Name:        "Beacon API Demo",
		Description: "Demonstrates the Beacon API with a local server",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(beaconService),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Beacon API Demo",
		Width:  900,
		Height: 700,
		URL:    "/",
	})

	app.Run()
}
