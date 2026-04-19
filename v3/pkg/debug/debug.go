package debug

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

type Collector struct {
	report *Report
}

func New() *Collector {
	return &Collector{
		report: NewReport(),
	}
}

func (c *Collector) Run() (*Report, error) {
	if err := c.collectSystemInfo(); err != nil {
		return nil, err
	}

	if err := c.collectBuildInfo(); err != nil {
		return nil, err
	}

	c.collectDiagnostics()

	return c.report, nil
}

func (c *Collector) collectSystemInfo() error {
	d := doctor.New()
	report, err := d.Run()
	if err != nil {
		return err
	}

	c.report.System = report.System

	return nil
}

func (c *Collector) collectBuildInfo() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	settings := make(map[string]string)
	for _, s := range buildInfo.Settings {
		settings[s.Key] = s.Value
	}

	c.report.Build = BuildInfo{
		GoVersion:  runtime.Version(),
		BuildMode:  settings["-buildmode"],
		Compiler:   settings["-compiler"],
		CGOEnabled: settings["CGO_ENABLED"] == "1",
		Settings:   settings,
	}

	return nil
}

func (c *Collector) collectDiagnostics() {
	d := doctor.New()
	report, err := d.Run()
	if err != nil {
		return
	}

	c.report.Diagnostics = report.Diagnostics

	c.report.CrashInfo = c.collectCrashInfo()
}

func (c *Collector) collectCrashInfo() *CrashInfo {
	info := &CrashInfo{
		Timestamp:   time.Now(),
		OS:         c.report.System.OS,
		Hardware:   c.report.System.Hardware,
		Build:      c.report.Build,
		Environment: collectEnvironmentVars(),
	}

	collectProcessInfo(info)
	collectMemoryInfo(info)
	collectThreadInfo(info)
	loadModules(info)

	return info
}