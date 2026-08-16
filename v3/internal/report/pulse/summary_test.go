package pulse

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
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
}
