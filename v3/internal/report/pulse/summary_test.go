package pulse

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

func TestUnknownExitCodeShowsUnderlyingError(t *testing.T) {
	var output bytes.Buffer
	reporter := New(&output, report.Normal)
	reporter.BuildStart("build", "linux/amd64", 1)
	step := reporter.StepStart("frontend:install", "Install frontend dependencies")
	reporter.StepFailed(step, report.Failure{
		Task:     "frontend:install",
		ExitCode: -1,
		Err:      errors.New(`npm not found: exec: "npm": executable file not found in $PATH`),
	})
	reporter.BuildEnd(time.Second, false)

	rendered := output.String()
	if !strings.Contains(rendered, `npm not found: exec: "npm": executable file not found in $PATH`) {
		t.Fatalf("failure panel omitted the underlying error:\n%s", rendered)
	}
	if strings.Contains(rendered, "exit -1") {
		t.Fatalf("failure panel rendered an unknown exit code as a real exit status:\n%s", rendered)
	}
	if strings.Contains(rendered, "│ error     npm not found") {
		t.Fatalf("failure panel repeated the error badge inside the body:\n%s", rendered)
	}
}

func TestFailurePanelHeaderKeepsBorderColourAfterStyledLabels(t *testing.T) {
	var output bytes.Buffer
	reporter := &Reporter{w: &output, s: newStyler(ProfileANSI)}
	reporter.writePanelLocked("frontend:install", report.Failure{ExitCode: -1, Err: errors.New("npm not found")})

	top := strings.SplitN(output.String(), "\n", 2)[0]
	if strings.Contains(top, ansi.Reset+" ") {
		t.Fatalf("failure header returned to the terminal's default colour between red labels:\n%q", top)
	}
}
