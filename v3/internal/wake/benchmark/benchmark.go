// Package benchmark measures complete Wails CLI invocations and applies
// explicit median-based performance budgets.
package benchmark

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Scenario         string
	Command          []string
	BeforeEach       []string
	WorkingDirectory string
	Environment      []string
	Artifacts        []string
	Warmups          int
	Samples          int
}

type Artifact struct {
	Digest    string `json:"digest"`
	SizeBytes int64  `json:"size_bytes"`
}

type Sample struct {
	WallMS          float64             `json:"wall_ms"`
	CPUMS           float64             `json:"cpu_ms"`
	PeakRSSBytes    int64               `json:"peak_rss_bytes,omitempty"`
	ReportedBuildMS float64             `json:"reported_build_ms,omitempty"`
	ExecutedSteps   int                 `json:"executed_steps,omitempty"`
	CachedSteps     int                 `json:"cached_steps,omitempty"`
	StepsReported   bool                `json:"steps_reported,omitempty"`
	Artifacts       map[string]Artifact `json:"artifacts,omitempty"`
}

type Result struct {
	Version            int      `json:"version"`
	Scenario           string   `json:"scenario"`
	Command            []string `json:"command"`
	Samples            []Sample `json:"samples"`
	MedianWallMS       float64  `json:"median_wall_ms"`
	MinWallMS          float64  `json:"min_wall_ms"`
	MaxWallMS          float64  `json:"max_wall_ms"`
	MedianCPUMS        float64  `json:"median_cpu_ms"`
	MedianPeakRSSBytes int64    `json:"median_peak_rss_bytes,omitempty"`
	// MedianReportedBuildMS is the CLI's own graph duration. WallMS includes
	// process startup and reporting, so keeping both separates orchestration
	// critical-path work from complete command latency.
	MedianReportedBuildMS float64 `json:"median_reported_build_ms,omitempty"`
	// MedianAbsoluteDeviationMS and MADPercent describe run-to-run wall-time
	// variance without allowing one noisy outlier to decide the result.
	MedianAbsoluteDeviationMS float64 `json:"median_absolute_deviation_ms"`
	MADPercent                float64 `json:"mad_percent"`
	// OrchestrationOverheadPercent is complete process wall time minus the
	// duration reported by the build graph, as a percentage of graph time.
	OrchestrationOverheadPercent float64 `json:"orchestration_overhead_percent,omitempty"`
}

type Budget struct {
	MaxMedianWallMS                 float64
	MaxRegressionPercent            float64
	MaxMADPercent                   float64
	MaxOrchestrationOverheadPercent float64
	MinSamples                      int
	ExpectedExecutedSteps           *int
	ExpectedCachedSteps             *int
	RequireStableArtifacts          bool
}

type namedWriteCloser interface {
	io.Writer
	io.Closer
	Name() string
}

type benchmarkFilesystem struct {
	mkdirAll   func(string, fs.FileMode) error
	createTemp func(string, string) (namedWriteCloser, error)
	remove     func(string) error
	rename     func(string, string) error
	lstat      func(string) (fs.FileInfo, error)
	open       func(string) (io.ReadCloser, error)
	walkDir    func(string, fs.WalkDirFunc) error
}

var operatingSystem = benchmarkFilesystem{
	mkdirAll: os.MkdirAll,
	createTemp: func(directory, pattern string) (namedWriteCloser, error) {
		return os.CreateTemp(directory, pattern)
	},
	remove: os.Remove,
	rename: os.Rename,
	lstat:  os.Lstat,
	open: func(path string) (io.ReadCloser, error) {
		return os.Open(path)
	},
	walkDir: filepath.WalkDir,
}

func Run(ctx context.Context, config Config) (Result, error) {
	if len(config.Command) == 0 {
		return Result{}, fmt.Errorf("benchmark command is required")
	}
	if config.Samples <= 0 {
		return Result{}, fmt.Errorf("samples must be positive")
	}
	for index := 0; index < config.Warmups; index++ {
		if err := prepare(ctx, config); err != nil {
			return Result{}, fmt.Errorf("prepare warmup %d: %w", index+1, err)
		}
		if _, _, err := runOnce(ctx, config); err != nil {
			return Result{}, fmt.Errorf("warmup %d: %w", index+1, err)
		}
	}
	result := Result{Version: 3, Scenario: config.Scenario, Command: append([]string(nil), config.Command...)}
	for index := 0; index < config.Samples; index++ {
		if err := prepare(ctx, config); err != nil {
			return Result{}, fmt.Errorf("prepare sample %d: %w", index+1, err)
		}
		sample, output, err := runOnce(ctx, config)
		if err != nil {
			return Result{}, fmt.Errorf("sample %d: %w\n%s", index+1, err, strings.TrimSpace(output))
		}
		result.Samples = append(result.Samples, sample)
	}
	result.MedianWallMS = median(sampleValues(result.Samples, func(sample Sample) float64 { return sample.WallMS }))
	result.MedianCPUMS = median(sampleValues(result.Samples, func(sample Sample) float64 { return sample.CPUMS }))
	result.MedianPeakRSSBytes = int64(median(sampleValues(result.Samples, func(sample Sample) float64 { return float64(sample.PeakRSSBytes) })))
	result.MedianReportedBuildMS = median(nonzeroSampleValues(result.Samples, func(sample Sample) float64 { return sample.ReportedBuildMS }))
	deviations := sampleValues(result.Samples, func(sample Sample) float64 {
		return absolute(sample.WallMS - result.MedianWallMS)
	})
	result.MedianAbsoluteDeviationMS = median(deviations)
	if result.MedianWallMS > 0 {
		result.MADPercent = result.MedianAbsoluteDeviationMS / result.MedianWallMS * 100
	}
	if result.MedianReportedBuildMS > 0 && result.MedianWallMS > result.MedianReportedBuildMS {
		result.OrchestrationOverheadPercent = (result.MedianWallMS - result.MedianReportedBuildMS) / result.MedianReportedBuildMS * 100
	}
	result.MinWallMS, result.MaxWallMS = result.Samples[0].WallMS, result.Samples[0].WallMS
	for _, sample := range result.Samples[1:] {
		if sample.WallMS < result.MinWallMS {
			result.MinWallMS = sample.WallMS
		}
		if sample.WallMS > result.MaxWallMS {
			result.MaxWallMS = sample.WallMS
		}
	}
	return result, nil
}

func Check(result Result, baseline *Result, maxMS, maxRegressionPercent float64) error {
	return CheckBudget(result, baseline, Budget{MaxMedianWallMS: maxMS, MaxRegressionPercent: maxRegressionPercent})
}

func CheckBudget(result Result, baseline *Result, budget Budget) error {
	if budget.MinSamples > 0 && len(result.Samples) < budget.MinSamples {
		return fmt.Errorf("%s has %d samples; acceptance requires at least %d", result.Scenario, len(result.Samples), budget.MinSamples)
	}
	if budget.MaxMedianWallMS > 0 && result.MedianWallMS > budget.MaxMedianWallMS {
		return fmt.Errorf("%s median %.2fms exceeds %.2fms budget", result.Scenario, result.MedianWallMS, budget.MaxMedianWallMS)
	}
	if budget.MaxMADPercent > 0 && result.MADPercent > budget.MaxMADPercent {
		return fmt.Errorf("%s wall-time MAD %.2f%% exceeds %.2f%% variance budget", result.Scenario, result.MADPercent, budget.MaxMADPercent)
	}
	if budget.MaxOrchestrationOverheadPercent > 0 {
		if result.MedianReportedBuildMS <= 0 {
			return fmt.Errorf("%s has no reported build duration for orchestration-overhead acceptance", result.Scenario)
		}
		if result.OrchestrationOverheadPercent > budget.MaxOrchestrationOverheadPercent {
			return fmt.Errorf("%s orchestration overhead %.2f%% exceeds %.2f%% budget", result.Scenario, result.OrchestrationOverheadPercent, budget.MaxOrchestrationOverheadPercent)
		}
	}
	for index, sample := range result.Samples {
		if budget.ExpectedExecutedSteps == nil && budget.ExpectedCachedSteps == nil {
			break
		}
		if !sample.StepsReported {
			return fmt.Errorf("%s sample %d did not report executed and cached work", result.Scenario, index+1)
		}
		if budget.ExpectedExecutedSteps != nil && sample.ExecutedSteps != *budget.ExpectedExecutedSteps {
			return fmt.Errorf("%s sample %d executed %d steps; expected %d", result.Scenario, index+1, sample.ExecutedSteps, *budget.ExpectedExecutedSteps)
		}
		if budget.ExpectedCachedSteps != nil && sample.CachedSteps != *budget.ExpectedCachedSteps {
			return fmt.Errorf("%s sample %d cached %d steps; expected %d", result.Scenario, index+1, sample.CachedSteps, *budget.ExpectedCachedSteps)
		}
	}
	if budget.RequireStableArtifacts && len(result.Samples) > 0 {
		reference := result.Samples[0].Artifacts
		if baseline != nil && len(baseline.Samples) > 0 {
			reference = baseline.Samples[0].Artifacts
		}
		for index, sample := range result.Samples {
			if err := compareArtifacts(reference, sample.Artifacts); err != nil {
				return fmt.Errorf("%s sample %d: %w", result.Scenario, index+1, err)
			}
		}
	}
	if budget.MaxRegressionPercent > 0 {
		if baseline == nil {
			return fmt.Errorf("%s regression budget requires a baseline", result.Scenario)
		}
		if baseline.Scenario != result.Scenario {
			return fmt.Errorf("baseline scenario %q does not match %q", baseline.Scenario, result.Scenario)
		}
		if baseline.MedianWallMS <= 0 {
			return fmt.Errorf("baseline median must be positive")
		}
		limit := baseline.MedianWallMS * (1 + budget.MaxRegressionPercent/100)
		if result.MedianWallMS > limit {
			return fmt.Errorf("%s median %.2fms exceeds %.2fms baseline budget (%.2fms + %.1f%%)", result.Scenario, result.MedianWallMS, limit, baseline.MedianWallMS, budget.MaxRegressionPercent)
		}
	}
	return nil
}

func Read(path string) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Write(path string, result Result) error {
	return writeResult(operatingSystem, path, result)
}

func writeResult(filesystem benchmarkFilesystem, path string, result Result) error {
	data, _ := json.MarshalIndent(result, "", "  ")
	data = append(data, '\n')
	directory := filepath.Dir(path)
	if err := filesystem.mkdirAll(directory, 0o755); err != nil {
		return err
	}
	temporary, err := filesystem.createTemp(directory, ".benchmark-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer filesystem.remove(temporaryPath)
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return filesystem.rename(temporaryPath, path)
}

func runOnce(ctx context.Context, config Config) (Sample, string, error) {
	cmd := exec.CommandContext(ctx, config.Command[0], config.Command[1:]...)
	cmd.Dir = config.WorkingDirectory
	if config.Environment != nil {
		cmd.Env = config.Environment
	}
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	started := time.Now()
	err := cmd.Run()
	sample := Sample{WallMS: durationMS(time.Since(started))}
	if cmd.ProcessState != nil {
		sample.CPUMS = durationMS(cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime())
		sample.PeakRSSBytes = peakRSSBytes(cmd.ProcessState)
	}
	parseBuildOutput(output.String(), &sample)
	if err == nil {
		sample.Artifacts, err = snapshotArtifacts(config)
	}
	return sample, output.String(), err
}

func prepare(ctx context.Context, config Config) error {
	if len(config.BeforeEach) == 0 {
		return nil
	}
	cmd := exec.CommandContext(ctx, config.BeforeEach[0], config.BeforeEach[1:]...)
	cmd.Dir = config.WorkingDirectory
	if config.Environment != nil {
		cmd.Env = config.Environment
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func snapshotArtifacts(config Config) (map[string]Artifact, error) {
	if len(config.Artifacts) == 0 {
		return nil, nil
	}
	result := make(map[string]Artifact, len(config.Artifacts))
	for _, configuredPath := range config.Artifacts {
		path := configuredPath
		if !filepath.IsAbs(path) && config.WorkingDirectory != "" {
			path = filepath.Join(config.WorkingDirectory, path)
		}
		artifact, err := snapshotArtifact(path)
		if err != nil {
			return nil, fmt.Errorf("snapshot artifact %s: %w", configuredPath, err)
		}
		result[configuredPath] = artifact
	}
	return result, nil
}

func snapshotArtifact(path string) (Artifact, error) {
	return snapshotArtifactWithFilesystem(operatingSystem, path)
}

func snapshotArtifactWithFilesystem(filesystem benchmarkFilesystem, path string) (Artifact, error) {
	info, err := filesystem.lstat(path)
	if err != nil {
		return Artifact{}, err
	}
	hash := sha256.New()
	var size int64
	if !info.IsDir() {
		if !info.Mode().IsRegular() {
			return Artifact{}, fmt.Errorf("not a regular file or directory")
		}
		file, err := filesystem.open(path)
		if err != nil {
			return Artifact{}, err
		}
		written, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return Artifact{}, copyErr
		}
		if closeErr != nil {
			return Artifact{}, closeErr
		}
		size = written
		return Artifact{Digest: fmt.Sprintf("sha256:%x", hash.Sum(nil)), SizeBytes: size}, nil
	}
	root := filepath.Clean(path)
	err = filesystem.walkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, _ := filepath.Rel(root, current)
		if relative == "." {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("artifact contains symlink %s", relative)
		}
		_, _ = fmt.Fprintf(hash, "%s\x00%o\x00", filepath.ToSlash(relative), info.Mode().Perm())
		if entry.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("artifact contains non-regular file %s", relative)
		}
		file, err := filesystem.open(current)
		if err != nil {
			return err
		}
		written, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		size += written
		return nil
	})
	if err != nil {
		return Artifact{}, err
	}
	return Artifact{Digest: fmt.Sprintf("sha256:%x", hash.Sum(nil)), SizeBytes: size}, nil
}

func compareArtifacts(expected, actual map[string]Artifact) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("artifact set changed: expected %d artifacts, got %d", len(expected), len(actual))
	}
	for path, expectedArtifact := range expected {
		actualArtifact, exists := actual[path]
		if !exists {
			return fmt.Errorf("artifact %s disappeared", path)
		}
		if actualArtifact != expectedArtifact {
			return fmt.Errorf("artifact %s changed: expected %s/%d bytes, got %s/%d bytes", path, expectedArtifact.Digest, expectedArtifact.SizeBytes, actualArtifact.Digest, actualArtifact.SizeBytes)
		}
	}
	return nil
}

var (
	csiPattern      = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
	oscPattern      = regexp.MustCompile(`\x1b\][^\x07]*(\x07|\x1b\\)`)
	stepsPattern    = regexp.MustCompile(`([0-9]+) ran(?:\s+•\s+([0-9]+) cached)?`)
	durationPattern = regexp.MustCompile(`build succeeded\s+([0-9.]+)(ms|s)`)
)

func parseBuildOutput(output string, sample *Sample) {
	plain := csiPattern.ReplaceAllString(oscPattern.ReplaceAllString(output, ""), "")
	if matches := stepsPattern.FindStringSubmatch(plain); matches != nil {
		sample.StepsReported = true
		sample.ExecutedSteps, _ = strconv.Atoi(matches[1])
		if matches[2] != "" {
			sample.CachedSteps, _ = strconv.Atoi(matches[2])
		}
	}
	if matches := durationPattern.FindStringSubmatch(plain); matches != nil {
		value, _ := strconv.ParseFloat(matches[1], 64)
		if matches[2] == "s" {
			value *= 1000
		}
		sample.ReportedBuildMS = value
	}
}

func sampleValues(samples []Sample, value func(Sample) float64) []float64 {
	result := make([]float64, len(samples))
	for index, sample := range samples {
		result[index] = value(sample)
	}
	return result
}

func nonzeroSampleValues(samples []Sample, value func(Sample) float64) []float64 {
	result := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if measured := value(sample); measured > 0 {
			result = append(result, measured)
		}
	}
	return result
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	values = append([]float64(nil), values...)
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}
	return values[middle]
}

func durationMS(duration time.Duration) float64 {
	return float64(duration.Microseconds()) / 1000
}

func absolute(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
