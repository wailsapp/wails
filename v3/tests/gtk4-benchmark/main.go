//go:build linux

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

// BenchmarkResult holds the result of a single benchmark
type BenchmarkResult struct {
	Name       string        `json:"name"`
	Iterations int           `json:"iterations"`
	TotalTime  time.Duration `json:"total_time_ns"`
	AvgTime    time.Duration `json:"avg_time_ns"`
	MinTime    time.Duration `json:"min_time_ns"`
	MaxTime    time.Duration `json:"max_time_ns"`
}

// BenchmarkReport holds all benchmark results
type BenchmarkReport struct {
	GTKVersion string            `json:"gtk_version"`
	Platform   string            `json:"platform"`
	GoVersion  string            `json:"go_version"`
	Timestamp  time.Time         `json:"timestamp"`
	Results    []BenchmarkResult `json:"results"`
}

var (
	app    *application.App
	report BenchmarkReport
)

func main() {
	report = BenchmarkReport{
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		GoVersion: runtime.Version(),
		Timestamp: time.Now(),
		Results:   []BenchmarkResult{},
	}

	app = application.New(application.Options{
		Name:        "GTK Benchmark",
		Description: "Benchmark comparing GTK3 vs GTK4 performance",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "GTK Benchmark",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	// Run benchmarks after a short delay to ensure app is initialized
	go func() {
		time.Sleep(1 * time.Second)
		runBenchmarks()
	}()

	err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runBenchmarks() {
	// Give the app a moment to fully initialize
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println("GTK Benchmark Suite")
	fmt.Println("=" + strings.Repeat("=", 59))

	// Detect GTK version
	report.GTKVersion = getGTKVersionString()
	fmt.Printf("GTK Version: %s\n", report.GTKVersion)
	fmt.Printf("Platform: %s\n", report.Platform)
	fmt.Printf("Go Version: %s\n", report.GoVersion)
	fmt.Println()

	// Run all benchmarks
	benchmarkScreenEnumeration()
	benchmarkWindowCreation()
	benchmarkWindowOperations()
	benchmarkMenuCreation()
	benchmarkEventDispatch()
	benchmarkDialogSetup()

	// Print and save report
	printReport()
	saveReport()

	// Exit after benchmarks complete
	time.Sleep(500 * time.Millisecond)
	app.Quit()
}

func benchmark(name string, iterations int, fn func()) BenchmarkResult {
	fmt.Printf("Running: %s (%d iterations)...\n", name, iterations)

	var times []time.Duration
	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		fn()
		elapsed := time.Since(start)
		times = append(times, elapsed)
		totalTime += elapsed
	}

	minTime := times[0]
	maxTime := times[0]
	for _, t := range times {
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
	}

	result := BenchmarkResult{
		Name:       name,
		Iterations: iterations,
		TotalTime:  totalTime,
		AvgTime:    totalTime / time.Duration(iterations),
		MinTime:    minTime,
		MaxTime:    maxTime,
	}

	report.Results = append(report.Results, result)
	fmt.Printf("  Average: %v\n", result.AvgTime)

	return result
}

func benchmarkScreenEnumeration() {
	benchmark("Screen Enumeration", 100, func() {
		_ = app.Screen.GetAll()
	})

	benchmark("Primary Screen Query", 100, func() {
		_ = app.Screen.GetPrimary()
	})
}

func benchmarkWindowCreation() {
	benchmark("Window Create/Destroy", 20, func() {
		w := app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title:  "Benchmark Window",
			Width:  400,
			Height: 300,
			Hidden: true,
		})
		// Small delay to ensure window is created
		time.Sleep(10 * time.Millisecond)
		w.Close()
		time.Sleep(10 * time.Millisecond)
	})
}

func benchmarkWindowOperations() {
	// Create a test window
	testWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Operations Test",
		Width:  400,
		Height: 300,
	})
	time.Sleep(100 * time.Millisecond)

	benchmark("Window SetSize", 50, func() {
		testWindow.SetSize(500, 400)
		testWindow.SetSize(400, 300)
	})

	benchmark("Window SetTitle", 100, func() {
		testWindow.SetTitle("Test Title " + time.Now().String())
	})

	benchmark("Window Size Query", 100, func() {
		_, _ = testWindow.Size()
	})

	benchmark("Window Position Query", 100, func() {
		_, _ = testWindow.Position()
	})

	benchmark("Window Center", 50, func() {
		testWindow.Center()
	})

	benchmark("Window Show/Hide", 20, func() {
		testWindow.Hide()
		time.Sleep(5 * time.Millisecond)
		testWindow.Show()
		time.Sleep(5 * time.Millisecond)
	})

	testWindow.Close()
}

func benchmarkMenuCreation() {
	benchmark("Menu Creation (Simple)", 50, func() {
		menu := app.Menu.New()
		menu.Add("Item 1")
		menu.Add("Item 2")
		menu.Add("Item 3")
	})

	benchmark("Menu Creation (Complex)", 20, func() {
		menu := app.Menu.New()
		for i := 0; i < 5; i++ {
			submenu := menu.AddSubmenu(fmt.Sprintf("Menu %d", i))
			for j := 0; j < 10; j++ {
				submenu.Add(fmt.Sprintf("Item %d-%d", i, j))
			}
		}
	})

	benchmark("Menu with Accelerators", 50, func() {
		menu := app.Menu.New()
		menu.Add("Cut").SetAccelerator("CmdOrCtrl+X")
		menu.Add("Copy").SetAccelerator("CmdOrCtrl+C")
		menu.Add("Paste").SetAccelerator("CmdOrCtrl+V")
	})
}

func benchmarkEventDispatch() {
	received := make(chan struct{}, 1000)

	app.Event.On("benchmark-event", func(event *application.CustomEvent) {
		select {
		case received <- struct{}{}:
		default:
		}
	})

	benchmark("Event Emit", 100, func() {
		app.Event.Emit("benchmark-event", map[string]interface{}{
			"timestamp": time.Now().UnixNano(),
			"data":      "test payload",
		})
	})

	benchmark("Event Emit+Receive", 50, func() {
		app.Event.Emit("benchmark-event", nil)
		select {
		case <-received:
		case <-time.After(100 * time.Millisecond):
		}
	})
}

func benchmarkDialogSetup() {
	// Dialog benchmarks - measure setup time only (Show() would block)
	benchmark("Dialog Setup (Info)", 20, func() {
		_ = app.Dialog.Info().
			SetTitle("Benchmark").
			SetMessage("Test message")
	})

	benchmark("Dialog Setup (Question)", 20, func() {
		_ = app.Dialog.Question().
			SetTitle("Benchmark").
			SetMessage("Test question?")
	})
}

func printReport() {
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println("Benchmark Results")
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Printf("GTK Version: %s\n", report.GTKVersion)
	fmt.Printf("Platform: %s\n", report.Platform)
	fmt.Printf("Timestamp: %s\n", report.Timestamp.Format(time.RFC3339))
	fmt.Println()

	fmt.Printf("%-35s %10s %12s %12s\n", "Benchmark", "Iterations", "Avg Time", "Total Time")
	fmt.Println(strings.Repeat("-", 75))

	for _, r := range report.Results {
		fmt.Printf("%-35s %10d %12v %12v\n",
			r.Name, r.Iterations, r.AvgTime, r.TotalTime)
	}
	fmt.Println()
}

func saveReport() {
	filename := fmt.Sprintf("benchmark-%s-%s.json",
		strings.ReplaceAll(report.GTKVersion, " ", "-"),
		report.Timestamp.Format("20060102-150405"))
	filename = strings.ReplaceAll(filename, "/", "-")
	filename = strings.ReplaceAll(filename, "(", "")
	filename = strings.ReplaceAll(filename, ")", "")

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling report: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
		return
	}

	fmt.Printf("Report saved to: %s\n", filename)
}
