//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/benchmark"
)

func main() {
	name := flag.String("name", "manifest-build", "benchmark scenario name")
	samples := flag.Int("samples", 7, "number of measured samples")
	warmups := flag.Int("warmups", 2, "number of unmeasured warmups")
	output := flag.String("output", "", "optional JSON output path")
	baselinePath := flag.String("baseline", "", "optional baseline JSON")
	maxMS := flag.Float64("max-ms", 0, "absolute median wall-time budget")
	maxRegression := flag.Float64("max-regression", 0, "allowed median regression percentage")
	maxMAD := flag.Float64("max-mad-percent", 15, "maximum median absolute deviation as a percentage of median wall time")
	maxOverhead := flag.Float64("max-overhead-percent", 0, "maximum process orchestration overhead over reported graph duration")
	minSamples := flag.Int("min-samples", 5, "minimum sample count required for acceptance")
	expectedRan := flag.Int("expect-ran", -1, "expected executed step count in every sample; negative disables")
	expectedCached := flag.Int("expect-cached", -1, "expected cached step count in every sample; negative disables")
	timeout := flag.Duration("timeout", 10*time.Minute, "timeout for the complete sample set")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./scripts/benchmark-manifest-build.go [flags] -- command [args...]")
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, err := benchmark.Run(ctx, benchmark.Config{Scenario: *name, Command: flag.Args(), Warmups: *warmups, Samples: *samples})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *output != "" {
		if err := benchmark.Write(*output, result); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
	var baseline *benchmark.Result
	if *baselinePath != "" {
		baseline, err = benchmark.Read(*baselinePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	budget := benchmark.Budget{MaxMedianWallMS: *maxMS, MaxRegressionPercent: *maxRegression, MaxMADPercent: *maxMAD, MaxOrchestrationOverheadPercent: *maxOverhead, MinSamples: *minSamples}
	if *expectedRan >= 0 {
		budget.ExpectedExecutedSteps = expectedRan
	}
	if *expectedCached >= 0 {
		budget.ExpectedCachedSteps = expectedCached
	}
	if err := benchmark.CheckBudget(result, baseline, budget); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
