//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/benchmark"
)

func main() {
	name := flag.String("name", "manifest-build", "benchmark scenario name")
	samples := flag.Int("samples", 7, "number of measured samples")
	warmups := flag.Int("warmups", 2, "number of unmeasured warmups")
	output := flag.String("output", "", "optional JSON output path")
	directory := flag.String("dir", "", "working directory for the measured and preparation commands")
	beforeJSON := flag.String("before-json", "", "optional JSON command array run before every warmup and sample outside measurement")
	artifactsCSV := flag.String("artifacts", "", "comma-separated artifact paths to hash and size after every sample")
	baselinePath := flag.String("baseline", "", "optional baseline JSON")
	maxMS := flag.Float64("max-ms", 0, "absolute median wall-time budget")
	maxRegression := flag.Float64("max-regression", 0, "allowed median regression percentage")
	maxMAD := flag.Float64("max-mad-percent", 15, "maximum median absolute deviation as a percentage of median wall time")
	maxOverhead := flag.Float64("max-overhead-percent", 0, "maximum process orchestration overhead over reported graph duration")
	minSamples := flag.Int("min-samples", 5, "minimum sample count required for acceptance")
	expectedRan := flag.Int("expect-ran", -1, "expected executed step count in every sample; negative disables")
	expectedCached := flag.Int("expect-cached", -1, "expected cached step count in every sample; negative disables")
	requireStableArtifacts := flag.Bool("require-stable-artifacts", false, "require measured artifact digests and sizes to remain identical")
	timeout := flag.Duration("timeout", 10*time.Minute, "timeout for the complete sample set")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./scripts/benchmark-manifest-build.go [flags] -- command [args...]")
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var beforeEach []string
	if *beforeJSON != "" {
		if err := json.Unmarshal([]byte(*beforeJSON), &beforeEach); err != nil {
			fmt.Fprintln(os.Stderr, "invalid -before-json command:", err)
			os.Exit(2)
		}
		if len(beforeEach) == 0 {
			fmt.Fprintln(os.Stderr, "-before-json command must not be empty")
			os.Exit(2)
		}
	}
	var artifacts []string
	for _, path := range strings.Split(*artifactsCSV, ",") {
		if path = strings.TrimSpace(path); path != "" {
			artifacts = append(artifacts, path)
		}
	}
	if *requireStableArtifacts && len(artifacts) == 0 {
		fmt.Fprintln(os.Stderr, "-require-stable-artifacts requires -artifacts")
		os.Exit(2)
	}
	result, err := benchmark.Run(ctx, benchmark.Config{Scenario: *name, Command: flag.Args(), BeforeEach: beforeEach, WorkingDirectory: *directory, Artifacts: artifacts, Warmups: *warmups, Samples: *samples})
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
	budget := benchmark.Budget{MaxMedianWallMS: *maxMS, MaxRegressionPercent: *maxRegression, MaxMADPercent: *maxMAD, MaxOrchestrationOverheadPercent: *maxOverhead, MinSamples: *minSamples, RequireStableArtifacts: *requireStableArtifacts}
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
