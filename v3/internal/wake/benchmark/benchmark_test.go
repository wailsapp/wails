package benchmark

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBuildOutputCapturesWorkAndCriticalPathDuration(t *testing.T) {
	sample := Sample{}
	parseBuildOutput("\x1b[31m✓  build succeeded    5ms    0 ran  •  4 cached\x1b[0m", &sample)
	assert.Equal(t, 0, sample.ExecutedSteps)
	assert.Equal(t, 4, sample.CachedSteps)
	assert.Equal(t, 5.0, sample.ReportedBuildMS)
}

func TestCheckUsesAbsoluteAndRelativeMedianBudgets(t *testing.T) {
	baseline := &Result{Scenario: "badge-noop", MedianWallMS: 70}
	require.NoError(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 80}, baseline, 100, 20))
	assert.ErrorContains(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 101}, baseline, 100, 20), "exceeds 100.00ms")
	assert.ErrorContains(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 85}, baseline, 100, 20), "baseline budget")
	assert.ErrorContains(t, Check(Result{Scenario: "other", MedianWallMS: 80}, baseline, 100, 20), "does not match")
}

func TestCheckBudgetEnforcesSampleVarianceWorkAndOverhead(t *testing.T) {
	zero, four := 0, 4
	result := Result{
		Scenario: "badge-noop", MedianWallMS: 84, MedianReportedBuildMS: 81,
		MedianAbsoluteDeviationMS: 3, MADPercent: 3.57, OrchestrationOverheadPercent: 3.70,
		Samples: []Sample{
			{WallMS: 80, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
			{WallMS: 84, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
			{WallMS: 87, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
		},
	}
	budget := Budget{MinSamples: 3, MaxMADPercent: 15, MaxOrchestrationOverheadPercent: 5, ExpectedExecutedSteps: &zero, ExpectedCachedSteps: &four}
	require.NoError(t, CheckBudget(result, nil, budget))

	bad := result
	bad.MADPercent = 16
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "variance budget")
	bad = result
	bad.OrchestrationOverheadPercent = 6
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "orchestration overhead")
	bad = result
	bad.Samples[1].ExecutedSteps = 1
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "sample 2 executed 1 steps")
	bad = result
	bad.Samples = bad.Samples[:2]
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "at least 3")
}

func TestCheckBudgetRequiresEvidenceForRequestedChecks(t *testing.T) {
	zero := 0
	result := Result{Scenario: "badge-noop", MedianWallMS: 80, Samples: []Sample{{WallMS: 80}}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{MaxRegressionPercent: 20}), "requires a baseline")
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{MaxOrchestrationOverheadPercent: 5}), "no reported build duration")
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{ExpectedExecutedSteps: &zero}), "did not report")
}

func TestRunMeasuresCommandSamples(t *testing.T) {
	command := []string{"sh", "-c", "printf 'build succeeded 2ms  1 ran • 3 cached\\n'"}
	if runtime.GOOS == "windows" {
		command = []string{"cmd.exe", "/c", "echo build succeeded 2ms  1 ran • 3 cached"}
	}
	result, err := Run(context.Background(), Config{Scenario: "test", Command: command, Warmups: 1, Samples: 3})
	require.NoError(t, err)
	assert.Len(t, result.Samples, 3)
	assert.Equal(t, 1, result.Samples[0].ExecutedSteps)
	assert.Equal(t, 3, result.Samples[0].CachedSteps)
	assert.True(t, result.Samples[0].StepsReported)
	assert.Positive(t, result.MedianWallMS)
	assert.Equal(t, 2.0, result.MedianReportedBuildMS)
	assert.GreaterOrEqual(t, result.MedianAbsoluteDeviationMS, 0.0)
	assert.GreaterOrEqual(t, result.MADPercent, 0.0)
}

func TestResultRoundTrip(t *testing.T) {
	path := t.TempDir() + string(os.PathSeparator) + "result.json"
	want := Result{Version: 1, Scenario: "badge", MedianWallMS: 70, Samples: []Sample{{WallMS: 70}}}
	require.NoError(t, Write(path, want))
	got, err := Read(path)
	require.NoError(t, err)
	assert.Equal(t, want, *got)
}

func TestRunPreparesEveryWarmupAndSampleOutsideTheMeasurement(t *testing.T) {
	root := t.TempDir()
	counter := filepath.Join(root, "counter")
	result, err := Run(context.Background(), Config{
		Scenario:         "prepared",
		Command:          benchmarkHelperCommand(t, "write", "observed-cwd", "yes"),
		BeforeEach:       benchmarkHelperCommand(t, "increment", counter),
		WorkingDirectory: root,
		Environment:      benchmarkHelperEnvironment(),
		Warmups:          2,
		Samples:          3,
	})
	require.NoError(t, err)
	assert.Len(t, result.Samples, 3)
	value, err := os.ReadFile(counter)
	require.NoError(t, err)
	assert.Equal(t, "5", string(value))
	value, err = os.ReadFile(filepath.Join(root, "observed-cwd"))
	require.NoError(t, err)
	assert.Equal(t, "yes", string(value))
}

func TestRunReportsPreparationFailure(t *testing.T) {
	_, err := Run(context.Background(), Config{
		Scenario:    "prepare-failure",
		Command:     benchmarkHelperCommand(t, "ok"),
		BeforeEach:  benchmarkHelperCommand(t, "fail"),
		Environment: benchmarkHelperEnvironment(),
		Samples:     1,
	})
	assert.ErrorContains(t, err, "prepare sample 1")
}

func TestRunCapturesArtifactDigestSizeAndPeakRSS(t *testing.T) {
	root := t.TempDir()
	artifact := filepath.Join(root, "artifact.bin")
	result, err := Run(context.Background(), Config{
		Scenario:         "artifact",
		Command:          benchmarkHelperCommand(t, "write", artifact, "stable bytes"),
		WorkingDirectory: root,
		Environment:      benchmarkHelperEnvironment(),
		Artifacts:        []string{"artifact.bin"},
		Samples:          2,
	})
	require.NoError(t, err)
	require.Len(t, result.Samples, 2)
	for _, sample := range result.Samples {
		require.Contains(t, sample.Artifacts, "artifact.bin")
		assert.Equal(t, int64(12), sample.Artifacts["artifact.bin"].SizeBytes)
		assert.Equal(t, "sha256:3821461753e58afa7abe81ccec8ea5ac178ea27ee92ede53771a95a101928e40", sample.Artifacts["artifact.bin"].Digest)
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			assert.Positive(t, sample.PeakRSSBytes)
		}
	}
	require.NoError(t, CheckBudget(result, nil, Budget{RequireStableArtifacts: true}))
}

func TestCheckBudgetRejectsArtifactDrift(t *testing.T) {
	result := Result{Scenario: "drift", Samples: []Sample{
		{Artifacts: map[string]Artifact{"app": {Digest: "sha256:one", SizeBytes: 1}}},
		{Artifacts: map[string]Artifact{"app": {Digest: "sha256:two", SizeBytes: 1}}},
	}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{RequireStableArtifacts: true}), "artifact app changed")
}

func TestCheckBudgetUsesMeasuredArtifactsWhenBaselineHasNoArtifactEvidence(t *testing.T) {
	artifacts := map[string]Artifact{"app": {Digest: "sha256:stable", SizeBytes: 12}}
	result := Result{Scenario: "stable", Samples: []Sample{
		{Artifacts: artifacts},
		{Artifacts: artifacts},
	}}
	baseline := &Result{Scenario: "stable", Samples: []Sample{{}}}

	require.NoError(t, CheckBudget(result, baseline, Budget{RequireStableArtifacts: true}))
}

func TestRunRejectsInvalidConfigurationAndFailedCommands(t *testing.T) {
	_, err := Run(context.Background(), Config{Samples: 1})
	assert.ErrorContains(t, err, "command is required")

	_, err = Run(context.Background(), Config{Command: benchmarkHelperCommand(t, "ok")})
	assert.ErrorContains(t, err, "samples must be positive")

	environment := benchmarkHelperEnvironment()
	_, err = Run(context.Background(), Config{Command: benchmarkHelperCommand(t, "ok"), BeforeEach: benchmarkHelperCommand(t, "fail"), Environment: environment, Warmups: 1, Samples: 1})
	assert.ErrorContains(t, err, "prepare warmup 1")

	_, err = Run(context.Background(), Config{Command: benchmarkHelperCommand(t, "fail"), Environment: environment, Warmups: 1, Samples: 1})
	assert.ErrorContains(t, err, "warmup 1")

	_, err = Run(context.Background(), Config{Command: benchmarkHelperCommand(t, "fail-output"), Environment: environment, Samples: 1})
	assert.ErrorContains(t, err, "sample 1")
	assert.ErrorContains(t, err, "deliberate output")

	_, err = Run(context.Background(), Config{Command: benchmarkHelperCommand(t, "ok"), Environment: environment, Artifacts: []string{"missing"}, Samples: 1})
	assert.ErrorContains(t, err, "snapshot artifact missing")
}

func TestRunComputesSecondsDurationOverheadAndPeakMemoryMedian(t *testing.T) {
	result, err := Run(context.Background(), Config{
		Scenario:    "reported",
		Command:     benchmarkHelperCommand(t, "report", "0.001s"),
		Environment: benchmarkHelperEnvironment(),
		Samples:     2,
	})
	require.NoError(t, err)
	assert.Equal(t, 1.0, result.MedianReportedBuildMS)
	assert.Positive(t, result.OrchestrationOverheadPercent)
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		assert.Positive(t, result.MedianPeakRSSBytes)
	}
}

func TestCheckBudgetCoversRemainingWorkArtifactAndBaselineFailures(t *testing.T) {
	zero, four := 0, 4
	result := Result{Scenario: "same", MedianWallMS: 10, Samples: []Sample{{StepsReported: true, ExecutedSteps: 0, CachedSteps: 3}}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{ExpectedExecutedSteps: &zero, ExpectedCachedSteps: &four}), "cached 3 steps")

	result.Samples = []Sample{{Artifacts: map[string]Artifact{"one": Artifact{Digest: "1"}}}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{RequireStableArtifacts: true, MinSamples: 2}), "at least 2")
	result.Samples = append(result.Samples, Sample{Artifacts: map[string]Artifact{}})
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{RequireStableArtifacts: true}), "artifact set changed")

	result.Samples[1].Artifacts = map[string]Artifact{"two": Artifact{Digest: "1"}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{RequireStableArtifacts: true}), "artifact one disappeared")

	baseline := &Result{Scenario: "same", MedianWallMS: 10, Samples: []Sample{{Artifacts: map[string]Artifact{"one": Artifact{Digest: "1"}}}}}
	result.Samples = []Sample{{Artifacts: map[string]Artifact{"one": Artifact{Digest: "1"}}}}
	require.NoError(t, CheckBudget(result, baseline, Budget{RequireStableArtifacts: true}))

	baseline.MedianWallMS = 0
	assert.ErrorContains(t, CheckBudget(result, baseline, Budget{MaxRegressionPercent: 1}), "baseline median must be positive")
}

func TestReadReportsMissingAndMalformedResults(t *testing.T) {
	root := t.TempDir()
	_, err := Read(filepath.Join(root, "missing.json"))
	require.Error(t, err)
	malformed := filepath.Join(root, "malformed.json")
	require.NoError(t, os.WriteFile(malformed, []byte("{"), 0o644))
	_, err = Read(malformed)
	assert.ErrorContains(t, err, "unexpected end")
}

func TestWriteReportsDirectoryAndRenameFailures(t *testing.T) {
	root := t.TempDir()
	parentFile := filepath.Join(root, "parent")
	require.NoError(t, os.WriteFile(parentFile, []byte("file"), 0o644))
	assert.Error(t, Write(filepath.Join(parentFile, "result.json"), Result{}))

	destinationDirectory := filepath.Join(root, "destination")
	require.NoError(t, os.Mkdir(destinationDirectory, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destinationDirectory, "keep"), []byte("keep"), 0o644))
	assert.Error(t, Write(destinationDirectory, Result{}))
}

func TestWriteReportsEveryFilesystemFailure(t *testing.T) {
	sentinel := errors.New("filesystem failure")
	base := benchmarkFilesystem{
		mkdirAll: func(string, fs.FileMode) error { return nil },
		createTemp: func(string, string) (namedWriteCloser, error) {
			return &faultWriteCloser{name: "temporary"}, nil
		},
		remove: func(string) error { return nil },
		rename: func(string, string) error { return nil },
	}

	filesystem := base
	filesystem.mkdirAll = func(string, fs.FileMode) error { return sentinel }
	assert.ErrorIs(t, writeResult(filesystem, "result.json", Result{}), sentinel)

	filesystem = base
	filesystem.createTemp = func(string, string) (namedWriteCloser, error) { return nil, sentinel }
	assert.ErrorIs(t, writeResult(filesystem, "result.json", Result{}), sentinel)

	filesystem = base
	filesystem.createTemp = func(string, string) (namedWriteCloser, error) {
		return &faultWriteCloser{name: "temporary", writeErr: sentinel}, nil
	}
	assert.ErrorIs(t, writeResult(filesystem, "result.json", Result{}), sentinel)

	filesystem = base
	filesystem.createTemp = func(string, string) (namedWriteCloser, error) {
		return &faultWriteCloser{name: "temporary", closeErr: sentinel}, nil
	}
	assert.ErrorIs(t, writeResult(filesystem, "result.json", Result{}), sentinel)

	filesystem = base
	filesystem.rename = func(string, string) error { return sentinel }
	assert.ErrorIs(t, writeResult(filesystem, "result.json", Result{}), sentinel)
}

func TestSnapshotArtifactCoversFilesDirectoriesAndRejectedEntries(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated rights on Windows")
	}
	root := t.TempDir()
	file := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(file, []byte("file"), 0o755))
	artifact, err := snapshotArtifact(file)
	require.NoError(t, err)
	assert.Equal(t, int64(4), artifact.SizeBytes)

	directory := filepath.Join(root, "tree")
	require.NoError(t, os.MkdirAll(filepath.Join(directory, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(directory, "nested", "data"), []byte("data"), 0o644))
	artifact, err = snapshotArtifact(directory)
	require.NoError(t, err)
	assert.Equal(t, int64(4), artifact.SizeBytes)

	_, err = snapshotArtifact(filepath.Join(root, "missing"))
	require.Error(t, err)
	require.NoError(t, os.Symlink(file, filepath.Join(root, "link")))
	_, err = snapshotArtifact(filepath.Join(root, "link"))
	assert.ErrorContains(t, err, "not a regular file or directory")

	require.NoError(t, os.Symlink(file, filepath.Join(directory, "nested", "link")))
	_, err = snapshotArtifact(directory)
	assert.ErrorContains(t, err, "contains symlink")
}

func TestSnapshotArtifactReportsEveryFilesystemFailure(t *testing.T) {
	sentinel := errors.New("filesystem failure")
	root := t.TempDir()
	file := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(file, []byte("file"), 0o644))
	directory := filepath.Join(root, "directory")
	require.NoError(t, os.Mkdir(directory, 0o755))

	filesystem := operatingSystem
	filesystem.open = func(string) (io.ReadCloser, error) { return nil, sentinel }
	_, err := snapshotArtifactWithFilesystem(filesystem, file)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.open = func(string) (io.ReadCloser, error) {
		return &faultReadCloser{readErr: sentinel}, nil
	}
	_, err = snapshotArtifactWithFilesystem(filesystem, file)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.open = func(string) (io.ReadCloser, error) {
		return &faultReadCloser{reader: bytes.NewBufferString("file"), closeErr: sentinel}, nil
	}
	_, err = snapshotArtifactWithFilesystem(filesystem, file)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.walkDir = func(path string, visit fs.WalkDirFunc) error {
		return visit(path, nil, sentinel)
	}
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.walkDir = func(path string, visit fs.WalkDirFunc) error {
		if err := visit(path, filesystemDirEntry{name: filepath.Base(path), directory: true}, nil); err != nil {
			return err
		}
		return visit(filepath.Join(path, "bad"), filesystemDirEntry{name: "bad", infoErr: sentinel}, nil)
	}
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.walkDir = func(string, fs.WalkDirFunc) error { return sentinel }
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)

	walkSingleFile := func(path string, visit fs.WalkDirFunc) error {
		if err := visit(path, filesystemDirEntry{name: filepath.Base(path), directory: true}, nil); err != nil {
			return err
		}
		return visit(filepath.Join(path, "file"), filesystemDirEntry{name: "file"}, nil)
	}
	filesystem = operatingSystem
	filesystem.walkDir = walkSingleFile
	filesystem.open = func(string) (io.ReadCloser, error) { return nil, sentinel }
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.walkDir = walkSingleFile
	filesystem.open = func(string) (io.ReadCloser, error) { return &faultReadCloser{readErr: sentinel}, nil }
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)

	filesystem = operatingSystem
	filesystem.walkDir = walkSingleFile
	filesystem.open = func(string) (io.ReadCloser, error) {
		return &faultReadCloser{reader: bytes.NewBufferString("file"), closeErr: sentinel}, nil
	}
	_, err = snapshotArtifactWithFilesystem(filesystem, directory)
	assert.ErrorIs(t, err, sentinel)
}

type faultWriteCloser struct {
	name     string
	writeErr error
	closeErr error
}

func (file *faultWriteCloser) Name() string { return file.name }
func (file *faultWriteCloser) Write(data []byte) (int, error) {
	if file.writeErr != nil {
		return 0, file.writeErr
	}
	return len(data), nil
}
func (file *faultWriteCloser) Close() error { return file.closeErr }

type faultReadCloser struct {
	reader   io.Reader
	readErr  error
	closeErr error
}

func (file *faultReadCloser) Read(data []byte) (int, error) {
	if file.readErr != nil {
		return 0, file.readErr
	}
	return file.reader.Read(data)
}
func (file *faultReadCloser) Close() error { return file.closeErr }

type filesystemDirEntry struct {
	name      string
	directory bool
	infoErr   error
}

func (entry filesystemDirEntry) Name() string      { return entry.name }
func (entry filesystemDirEntry) IsDir() bool       { return entry.directory }
func (entry filesystemDirEntry) Type() fs.FileMode { return 0 }
func (entry filesystemDirEntry) Info() (fs.FileInfo, error) {
	if entry.infoErr != nil {
		return nil, entry.infoErr
	}
	return filesystemFileInfo{name: entry.name, directory: entry.directory}, nil
}

type filesystemFileInfo struct {
	name      string
	directory bool
}

func (info filesystemFileInfo) Name() string       { return info.name }
func (info filesystemFileInfo) Size() int64        { return 0 }
func (info filesystemFileInfo) Mode() fs.FileMode  { return 0o644 }
func (info filesystemFileInfo) ModTime() time.Time { return time.Time{} }
func (info filesystemFileInfo) IsDir() bool        { return info.directory }
func (info filesystemFileInfo) Sys() any           { return nil }

func benchmarkHelperCommand(t *testing.T, args ...string) []string {
	t.Helper()
	executable, err := os.Executable()
	require.NoError(t, err)
	return append([]string{executable, "-test.run=TestBenchmarkHelperProcess", "--"}, args...)
}

func benchmarkHelperEnvironment() []string {
	return append(os.Environ(), "GO_WANT_BENCHMARK_HELPER=1")
}

func TestBenchmarkHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_BENCHMARK_HELPER") != "1" {
		return
	}
	separator := -1
	for index, argument := range os.Args {
		if argument == "--" {
			separator = index
			break
		}
	}
	if separator < 0 || separator+1 >= len(os.Args) {
		os.Exit(97)
	}
	arguments := os.Args[separator+1:]
	switch arguments[0] {
	case "ok":
		os.Exit(0)
	case "fail":
		os.Exit(23)
	case "fail-output":
		_, _ = os.Stderr.WriteString("deliberate output\n")
		os.Exit(28)
	case "report":
		time.Sleep(10 * time.Millisecond)
		_, _ = os.Stdout.WriteString("build succeeded " + arguments[1] + " 0 ran • 0 cached\n")
	case "write":
		if err := os.WriteFile(arguments[1], []byte(arguments[2]), 0o644); err != nil {
			os.Exit(25)
		}
	case "increment":
		value := 0
		if data, err := os.ReadFile(arguments[1]); err == nil {
			value, _ = strconv.Atoi(strings.TrimSpace(string(data)))
		} else if !os.IsNotExist(err) {
			os.Exit(26)
		}
		if err := os.WriteFile(arguments[1], []byte(strconv.Itoa(value+1)), 0o644); err != nil {
			os.Exit(27)
		}
	default:
		os.Exit(98)
	}
	os.Exit(0)
}
