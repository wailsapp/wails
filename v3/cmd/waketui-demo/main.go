// Command waketui-demo drives the pulse Reporter with a scripted fake build
// so the visual design can be inspected without spinning up a real wake run.
//
// Subcommands:
//
//	serial      one worker, mixed cached/run/ok steps, success verdict
//	parallel    four workers, overlapping in-flight steps, success verdict
//	fail        one step fails partway through, panel rendered at end
//	multifail   three steps fail across a parallel build, banner + 3 panels
//	verbose     parallel with -v to show streamed commands and output
//	debug       parallel with -vv showing dag/dep/var debug lines
//
// Pass --no-anim to force the static non-TTY fallback path even from a real
// terminal — useful for screenshotting both modes side by side.
//
// Run from the repo root:
//
//	go run ./cmd/waketui-demo serial
//	go run ./cmd/waketui-demo parallel
//	go run ./cmd/waketui-demo fail
//	go run ./cmd/waketui-demo multifail
package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse"
)

func main() {
	verbose := flag.Bool("v", false, "verbose (show commands and output)")
	debug := flag.Bool("vv", false, "debug (show DAG/dep/var trace lines)")
	speed := flag.Float64("speed", 1.0, "time multiplier (smaller = faster demo)")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: waketui-demo [serial|parallel|fail|multifail|verbose|debug] [-v|-vv] [-speed=N]")
		os.Exit(2)
	}

	level := report.Normal
	if *verbose {
		level = report.Verbose
	}
	if *debug {
		level = report.Debug
	}

	r := pulse.New(os.Stdout, level)

	switch args[0] {
	case "serial":
		runSerial(r, *speed)
	case "parallel":
		runParallel(r, *speed, false, 0)
	case "fail":
		runSerialFail(r, *speed)
	case "multifail":
		runParallel(r, *speed, true, 3)
	case "verbose":
		level = report.Verbose
		r = pulse.New(os.Stdout, level)
		runParallel(r, *speed, false, 0)
	case "debug":
		level = report.Debug
		r = pulse.New(os.Stdout, level)
		runWithDebug(r, *speed)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		os.Exit(2)
	}
}

// --- scenarios -----------------------------------------------------------

func runSerial(r *pulse.Reporter, speed float64) {
	steps := []scriptStep{
		{name: "darwin:common:go:mod:tidy", label: "", dur: 350, status: report.StatusCached},
		{name: "darwin:common:generate:bindings", label: "", dur: 480, status: report.StatusOK,
			infos: []string{"loading packages…", "found 12 services", "writing bindings/"}},
		{name: "darwin:common:build:frontend", label: "", dur: 1300, status: report.StatusOK,
			infos: []string{"vite v5.1.0", "transforming…", "rendering chunks…", "built in 1.1s"}},
		{name: "darwin:build:amd64", label: "build:amd64", dur: 2100, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:build:arm64", label: "build:arm64", dur: 2050, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:lipo", label: "", dur: 180, status: report.StatusOK},
		{name: "darwin:codesign", label: "", dur: 920, status: report.StatusOK,
			infos: []string{"signing bundle", "verifying signature"}},
		{name: "darwin:notarize", label: "notarize (skipped — dev build)", dur: 0, status: report.StatusSkipped},
	}
	r.BuildStart("build", "darwin:build", len(steps))
	start := time.Now()
	for _, s := range steps {
		runOne(r, s, speed)
	}
	// Register the build artifacts the wake executor would discover from the
	// task's `generates:` declarations. Sizes are made-up for the demo; in a
	// real wake build, Artifact() is called with the actual file size from
	// the wake executor (or zero, in which case Artifact() stats the path).
	r.Artifact(pulse.Artifact{Path: "bin/myapp", Size: 12_734_512, Kind: "binary"})
	r.Artifact(pulse.Artifact{Path: "bin/myapp.dSYM", Size: 3_445_120, Kind: "debug"})
	r.BuildEnd(time.Since(start), true)
}

func runSerialFail(r *pulse.Reporter, speed float64) {
	steps := []scriptStep{
		{name: "go:mod:tidy", dur: 240, status: report.StatusCached},
		{name: "generate:bindings", dur: 380, status: report.StatusOK,
			infos: []string{"loading packages…", "found 12 services"}},
		{name: "test:unit", dur: 1800, status: report.StatusFailed,
			infos: []string{"running 184 tests", "FAIL TestAnalyser/Service13"},
			failure: &report.Failure{
				Command:  "go test -count=1 ./internal/generator/...",
				ExitCode: 1,
				Output: `=== RUN   TestAnalyser
=== RUN   TestAnalyser/Service13
    analyser_test.go:142: missing service "Service13" in generated bindings
    analyser_test.go:143:   want: ["Service11" "Service12" "Service13"]
    analyser_test.go:144:   got:  ["Service11" "Service12"]
--- FAIL: TestAnalyser (0.42s)
    --- FAIL: TestAnalyser/Service13 (0.18s)
FAIL
exit status 1`,
			}},
	}
	r.BuildStart("test", "ci", len(steps))
	start := time.Now()
	for _, s := range steps {
		runOne(r, s, speed)
		if s.status == report.StatusFailed {
			break
		}
	}
	r.BuildEnd(time.Since(start), false)
}

func runParallel(r *pulse.Reporter, speed float64, withFailures bool, nfail int) {
	steps := []scriptStep{
		{name: "darwin:common:go:mod:tidy", dur: 200, status: report.StatusCached},
		{name: "darwin:common:generate:bindings", dur: 460, status: report.StatusOK,
			infos: []string{"loading packages…", "found 12 services"}},
		{name: "darwin:common:build:frontend", dur: 1200, status: report.StatusOK,
			infos: []string{"vite v5.1.0", "transforming…", "built in 1.0s"}},
		{name: "darwin:build:amd64", label: "build:amd64", dur: 2100, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:build:arm64", label: "build:arm64", dur: 2050, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:test:unit", label: "test:unit", dur: 1700, status: report.StatusOK,
			infos: []string{"running 184 tests", "184 passed in 1.6s"}},
		{name: "darwin:test:integration", label: "test:integration", dur: 1450, status: report.StatusOK,
			infos: []string{"running 24 tests"}},
		{name: "darwin:lipo", dur: 180, status: report.StatusOK},
		{name: "darwin:vet", label: "vet", dur: 620, status: report.StatusOK,
			infos: []string{"vetting 248 packages"}},
		{name: "darwin:lint", label: "lint", dur: 540, status: report.StatusOK,
			infos: []string{"40 packages clean"}},
		{name: "darwin:codesign", dur: 900, status: report.StatusOK},
		{name: "darwin:package:dmg", label: "package:dmg", dur: 320, status: report.StatusOK,
			infos: []string{"building dmg"}},
	}
	if withFailures {
		// Mark the last N steps as failures so the user sees the aggregation banner.
		failedIdxs := []int{5, 8, 11} // test:unit, vet, package:dmg
		if nfail > len(failedIdxs) {
			nfail = len(failedIdxs)
		}
		for _, i := range failedIdxs[:nfail] {
			steps[i].status = report.StatusFailed
			steps[i].failure = sampleFailure(steps[i].name)
		}
	}

	r.BuildStart("build", "darwin:release", len(steps))
	start := time.Now()

	const workers = 4
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for _, s := range steps {
		s := s
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			runOneParallel(r, s, speed)
		}()
		// Stagger task launches so the parallel pane has a chance to fill up
		// visually instead of all four spinners launching on the same frame.
		time.Sleep(sleep(120, speed))
	}
	wg.Wait()
	// Register expected artifacts before BuildEnd so the summary's "Output"
	// section can render them. In a real wake build the executor calls
	// Artifact() per `generates:` declaration as each task completes.
	if !withFailures {
		r.Artifact(pulse.Artifact{Path: "bin/myapp.amd64", Size: 12_734_512, Kind: "binary"})
		r.Artifact(pulse.Artifact{Path: "bin/myapp.arm64", Size: 12_842_376, Kind: "binary"})
		r.Artifact(pulse.Artifact{Path: "bin/myapp", Size: 25_576_888, Kind: "universal"})
		r.Artifact(pulse.Artifact{Path: "bin/myapp.dmg", Size: 26_104_320, Kind: "installer"})
	}
	r.BuildEnd(time.Since(start), !withFailures)
}

func runWithDebug(r *pulse.Reporter, speed float64) {
	r.BuildStart("build", "darwin:release", 6)
	r.Debug(report.DebugLine{Category: "dag", Subject: "darwin:release",
		Fields: []report.DebugField{{Key: "tasks", Val: "6"}, {Key: "leaf", Val: "go:mod:tidy"}}})
	r.Debug(report.DebugLine{Category: "dep", Subject: "build:amd64", Arrow: "go:mod:tidy"})
	r.Debug(report.DebugLine{Category: "dep", Subject: "build:arm64", Arrow: "go:mod:tidy"})
	r.Debug(report.DebugLine{Category: "var", Subject: "VERSION", Arrow: "v3.0.0-pre.42",
		Fields: []report.DebugField{{Key: "src", Val: "git describe"}}})
	r.Debug(report.DebugLine{Category: "exec", Subject: "spawn worker pool",
		Fields: []report.DebugField{{Key: "workers", Val: "4"}}})

	steps := []scriptStep{
		{name: "darwin:common:go:mod:tidy", dur: 200, status: report.StatusCached},
		{name: "darwin:build:amd64", dur: 1500, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:build:arm64", dur: 1480, status: report.StatusOK,
			infos: []string{"compiling 248 packages", "linking…"}},
		{name: "darwin:test:unit", dur: 1200, status: report.StatusOK,
			infos: []string{"running 184 tests"}},
		{name: "darwin:lipo", dur: 180, status: report.StatusOK},
		{name: "darwin:codesign", dur: 750, status: report.StatusOK},
	}
	start := time.Now()
	for _, s := range steps {
		runOne(r, s, speed)
	}
	r.BuildEnd(time.Since(start), true)
}

// --- step driver ---------------------------------------------------------

type scriptStep struct {
	name    string
	label   string
	dur     int // ms
	status  report.Status
	infos   []string
	failure *report.Failure
}

// runOne drives one step through the serial Reporter interface.
func runOne(r *pulse.Reporter, s scriptStep, speed float64) {
	r.StepStart(s.name, s.label)
	if s.status == report.StatusCached || s.status == report.StatusSkipped {
		time.Sleep(sleep(40, speed))
		r.StepEnd(s.status, 0)
		return
	}
	emitInfos(r, s.infos, s.dur, speed, 0)
	if s.status == report.StatusFailed && s.failure != nil {
		r.StepFailed(*s.failure)
		return
	}
	r.StepEnd(s.status, time.Duration(s.dur)*time.Millisecond)
}

// runOneParallel drives one step through the parallel Reporter extension.
func runOneParallel(r *pulse.Reporter, s scriptStep, speed float64) {
	id := r.ParallelStepStart(s.name, s.label)
	if s.status == report.StatusCached || s.status == report.StatusSkipped {
		time.Sleep(sleep(80, speed))
		r.ParallelStepEnd(id, s.status, 0)
		return
	}
	emitInfosParallel(r, id, s.infos, s.dur, speed)
	if s.status == report.StatusFailed && s.failure != nil {
		r.ParallelStepFailed(id, *s.failure)
		return
	}
	r.ParallelStepEnd(id, s.status, time.Duration(s.dur)*time.Millisecond)
}

func emitInfos(r *pulse.Reporter, infos []string, totalMs int, speed float64, jitter int) {
	if len(infos) == 0 {
		time.Sleep(sleep(totalMs, speed))
		return
	}
	step := totalMs / len(infos)
	for _, msg := range infos {
		r.StepInfo(msg)
		time.Sleep(sleep(step+jitter, speed))
	}
}

func emitInfosParallel(r *pulse.Reporter, id pulse.StepID, infos []string, totalMs int, speed float64) {
	if len(infos) == 0 {
		time.Sleep(sleep(totalMs, speed))
		return
	}
	step := totalMs / len(infos)
	for _, msg := range infos {
		r.ParallelStepInfo(id, msg)
		// Each parallel worker jitters slightly so the spinners don't all
		// transition between detail lines on the same frame.
		j := rand.IntN(60) - 30
		time.Sleep(sleep(step+j, speed))
	}
}

func sleep(ms int, speed float64) time.Duration {
	if ms <= 0 {
		return 0
	}
	d := time.Duration(float64(ms)*speed) * time.Millisecond
	if d < 10*time.Millisecond {
		d = 10 * time.Millisecond
	}
	return d
}

func sampleFailure(taskName string) *report.Failure {
	switch taskName {
	case "darwin:test:unit":
		return &report.Failure{
			Command:  "go test -count=1 ./...",
			ExitCode: 1,
			Output: `=== RUN   TestAnalyser/Service13
    analyser_test.go:142: missing service "Service13" in generated bindings
    analyser_test.go:143:   want: ["Service11" "Service12" "Service13"]
    analyser_test.go:144:   got:  ["Service11" "Service12"]
--- FAIL: TestAnalyser (0.42s)
    --- FAIL: TestAnalyser/Service13 (0.18s)
FAIL
exit status 1`,
		}
	case "darwin:vet":
		return &report.Failure{
			Command:  "go vet ./...",
			ExitCode: 2,
			Output: `internal/wake/exec/runner.go:142:2: unreachable code
internal/wake/cmds/native.go:88:6: struct field tag ` + "`json:foo`" + ` not compatible with reflect.StructTag.Get`,
		}
	case "darwin:package:dmg":
		return &report.Failure{
			Command:  "hdiutil create -fs HFS+ -volname MyApp -srcfolder bin/ bin/MyApp.dmg",
			ExitCode: 1,
			Output: `hdiutil: create failed - No space left on device
DEBUG: hdiutil: detached image after failure`,
		}
	}
	return &report.Failure{Command: "true", ExitCode: 1, Output: "boom"}
}
