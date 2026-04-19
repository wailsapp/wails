// Package debug provides programmatic access to crash diagnostics for a Wails
// application: system info, build info, runtime process/memory/module state,
// and — on Windows — a minidump of the current process that can be opened in
// WinDbg or Visual Studio.
//
// Two entry points:
//
//   - Dump   — write a minidump of the current process (Windows only).
//   - Report — collect a rich crash/diagnostic snapshot, optionally including
//              a minidump.
//
// Both run under the calling user's token. No elevation is required to dump
// the current process: Windows' GetCurrentProcess pseudo-handle bypasses ACL
// checks for own-process access.
//
// Typical integration with application.PanicHandler:
//
//	app := application.New(application.Options{
//	    PanicHandler: func(pd *application.PanicDetails) {
//	        r, _ := debug.Report(debug.WithDump())
//	        // upload r / show dialog / log pd
//	    },
//	})
//
// Or manually, at any point in the app lifetime:
//
//	r, err := debug.Report(debug.WithDump())
//	log.Printf("dump=%s goroutines=%d", r.DumpPath, r.Crash.ProcessInfo.Goroutines)
package debug

import (
	"fmt"
	"runtime"
	rtdebug "runtime/debug"
	"time"

	doctor "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

// ---- Dump -----------------------------------------------------------------

type dumpConfig struct {
	path       string
	fullMemory bool
}

// DumpOption configures a Dump call.
type DumpOption func(*dumpConfig)

// WithPath sets an explicit output path for the dump file. When omitted,
// Dump writes to <TempDir>/wails-crash-<pid>-<unix>.dmp.
func WithPath(path string) DumpOption {
	return func(c *dumpConfig) { c.path = path }
}

// WithFullMemory requests a full-memory minidump. Produces much larger
// files (hundreds of MB to GB) but preserves the entire process address
// space for postmortem analysis.
func WithFullMemory() DumpOption {
	return func(c *dumpConfig) { c.fullMemory = true }
}

// Dump writes a minidump of the current process and returns the absolute
// path of the written file. On non-Windows platforms returns an error.
func Dump(opts ...DumpOption) (string, error) {
	cfg := &dumpConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return writeCoreDump(cfg.path, cfg.fullMemory)
}

// ---- Report ---------------------------------------------------------------

type reportConfig struct {
	withDump   bool
	dumpPath   string
	fullMemory bool
}

// ReportOption configures a Report call.
type ReportOption func(*reportConfig)

// WithDump instructs Report to also write a minidump of the current process
// (Windows only) and record its path in CrashReport.DumpPath.
func WithDump() ReportOption {
	return func(c *reportConfig) { c.withDump = true }
}

// WithDumpPath is like WithDump but writes the dump to an explicit location.
func WithDumpPath(path string) ReportOption {
	return func(c *reportConfig) {
		c.withDump = true
		c.dumpPath = path
	}
}

// WithDumpFullMemory is like WithDump but produces a full-memory minidump.
func WithDumpFullMemory() ReportOption {
	return func(c *reportConfig) {
		c.withDump = true
		c.fullMemory = true
	}
}

// Report collects a rich diagnostic snapshot of the current process:
// system info, build info, process/memory/module state, and a diagnostics
// pass from the doctor package. With WithDump(), a minidump is also written.
//
// Errors from the minidump step are attached to the returned error but do
// not prevent the other fields from being populated — callers get a partial
// report on failure.
func Report(opts ...ReportOption) (*CrashReport, error) {
	cfg := &reportConfig{}
	for _, o := range opts {
		o(cfg)
	}

	r := &CrashReport{Timestamp: time.Now()}

	if err := collectSystemInfo(r); err != nil {
		return r, fmt.Errorf("system info: %w", err)
	}
	collectBuildInfo(r)
	collectDiagnostics(r)
	r.Crash = collectCrashInfo(r)

	if cfg.withDump {
		path, err := writeCoreDump(cfg.dumpPath, cfg.fullMemory)
		if err != nil {
			return r, fmt.Errorf("minidump: %w", err)
		}
		r.DumpPath = path
	}
	return r, nil
}

// ---- Internal collectors --------------------------------------------------

func collectSystemInfo(r *CrashReport) error {
	d := doctor.New()
	rep, err := d.Run()
	if err != nil {
		return err
	}
	r.System = rep.System
	r.Diagnostics = rep.Diagnostics
	return nil
}

func collectBuildInfo(r *CrashReport) {
	buildInfo, ok := rtdebug.ReadBuildInfo()
	if !ok {
		return
	}
	settings := make(map[string]string, len(buildInfo.Settings))
	for _, s := range buildInfo.Settings {
		settings[s.Key] = s.Value
	}
	r.Build = BuildInfo{
		GoVersion:  runtime.Version(),
		BuildMode:  settings["-buildmode"],
		Compiler:   settings["-compiler"],
		CGOEnabled: settings["CGO_ENABLED"] == "1",
		Settings:   settings,
	}
}

// collectDiagnostics is currently folded into collectSystemInfo (same doctor
// pass covers both). Kept as a hook for future expansion — e.g. running a
// secondary diagnostic that is not part of the doctor report.
func collectDiagnostics(r *CrashReport) {}

func collectCrashInfo(r *CrashReport) *CrashInfo {
	info := &CrashInfo{
		Timestamp:   time.Now(),
		OS:          r.System.OS,
		Hardware:    r.System.Hardware,
		Build:       r.Build,
		Environment: collectEnvironmentVars(),
	}
	collectProcessInfo(info)
	collectMemoryInfo(info)
	collectThreadInfo(info)
	loadModules(info)
	return info
}
