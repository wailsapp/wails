//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// BenchmarkResult mirrors the struct in main.go
type BenchmarkResult struct {
	Name       string        `json:"name"`
	Iterations int           `json:"iterations"`
	TotalTime  time.Duration `json:"total_time_ns"`
	AvgTime    time.Duration `json:"avg_time_ns"`
	MinTime    time.Duration `json:"min_time_ns"`
	MaxTime    time.Duration `json:"max_time_ns"`
}

// BenchmarkReport mirrors the struct in main.go
type BenchmarkReport struct {
	GTKVersion string            `json:"gtk_version"`
	Platform   string            `json:"platform"`
	GoVersion  string            `json:"go_version"`
	Timestamp  time.Time         `json:"timestamp"`
	Results    []BenchmarkResult `json:"results"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run compare.go <gtk3-report.json> <gtk4-report.json>")
		os.Exit(1)
	}

	gtk3Report, err := loadReport(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading GTK3 report: %v\n", err)
		os.Exit(1)
	}

	gtk4Report, err := loadReport(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading GTK4 report: %v\n", err)
		os.Exit(1)
	}

	compareReports(gtk3Report, gtk4Report)
}

func loadReport(filename string) (*BenchmarkReport, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var report BenchmarkReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

func compareReports(gtk3, gtk4 *BenchmarkReport) {
	fmt.Println("=" + strings.Repeat("=", 89))
	fmt.Println("GTK3 vs GTK4 Benchmark Comparison")
	fmt.Println("=" + strings.Repeat("=", 89))
	fmt.Println()

	fmt.Println("Report Details:")
	fmt.Printf("  GTK3: %s (%s)\n", gtk3.GTKVersion, gtk3.Timestamp.Format(time.RFC3339))
	fmt.Printf("  GTK4: %s (%s)\n", gtk4.GTKVersion, gtk4.Timestamp.Format(time.RFC3339))
	fmt.Printf("  Platform: %s\n", gtk3.Platform)
	fmt.Println()

	// Build maps for easy lookup
	gtk3Results := make(map[string]BenchmarkResult)
	gtk4Results := make(map[string]BenchmarkResult)

	for _, r := range gtk3.Results {
		gtk3Results[r.Name] = r
	}
	for _, r := range gtk4.Results {
		gtk4Results[r.Name] = r
	}

	// Get all benchmark names
	names := make([]string, 0)
	for name := range gtk3Results {
		names = append(names, name)
	}
	sort.Strings(names)

	// Print comparison table
	fmt.Printf("%-35s %12s %12s %12s %10s\n", "Benchmark", "GTK3 Avg", "GTK4 Avg", "Difference", "Change")
	fmt.Println(strings.Repeat("-", 90))

	var totalGTK3, totalGTK4 time.Duration
	improvements := 0
	regressions := 0

	for _, name := range names {
		gtk3r, ok3 := gtk3Results[name]
		gtk4r, ok4 := gtk4Results[name]

		if !ok3 || !ok4 {
			continue
		}

		diff := gtk4r.AvgTime - gtk3r.AvgTime
		var pctChange float64
		if gtk3r.AvgTime > 0 {
			pctChange = float64(diff) / float64(gtk3r.AvgTime) * 100
		}

		totalGTK3 += gtk3r.AvgTime
		totalGTK4 += gtk4r.AvgTime

		changeSymbol := ""
		if pctChange < -5 {
			changeSymbol = "✓ FASTER"
			improvements++
		} else if pctChange > 5 {
			changeSymbol = "✗ SLOWER"
			regressions++
		} else {
			changeSymbol = "≈ SAME"
		}

		fmt.Printf("%-35s %12v %12v %+12v %+9.1f%% %s\n",
			name, gtk3r.AvgTime, gtk4r.AvgTime, diff, pctChange, changeSymbol)
	}

	fmt.Println(strings.Repeat("-", 90))

	// Summary
	totalDiff := totalGTK4 - totalGTK3
	var totalPctChange float64
	if totalGTK3 > 0 {
		totalPctChange = float64(totalDiff) / float64(totalGTK3) * 100
	}

	fmt.Printf("%-35s %12v %12v %+12v %+9.1f%%\n",
		"TOTAL", totalGTK3, totalGTK4, totalDiff, totalPctChange)

	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Improvements (>5%% faster): %d\n", improvements)
	fmt.Printf("  Regressions (>5%% slower):  %d\n", regressions)
	fmt.Printf("  No significant change:      %d\n", len(names)-improvements-regressions)
	fmt.Println()

	if totalPctChange < 0 {
		fmt.Printf("Overall: GTK4 is %.1f%% faster than GTK3\n", -totalPctChange)
	} else if totalPctChange > 0 {
		fmt.Printf("Overall: GTK4 is %.1f%% slower than GTK3\n", totalPctChange)
	} else {
		fmt.Println("Overall: No significant difference")
	}
}
