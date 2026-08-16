// Package migration owns the private workflow state used while replacing
// legacy Taskfiles with wails.toml. None of these fields are build intent, so
// they deliberately do not form part of the manifest schema.
package migration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const RelativeReportPath = ".wails/migration-report.json"

type Diagnostic struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	File     string `json:"file,omitempty"`
	Task     string `json:"task,omitempty"`
	Message  string `json:"message"`
}

// Taskfile records how a project Taskfile relates to defaults known by the
// migrating CLI. Classification is one of current-default,
// historical-default, customised, or custom.
type Taskfile struct {
	File           string   `json:"file"`
	Role           string   `json:"role"`
	Classification string   `json:"classification"`
	ModifiedTasks  []string `json:"modified_tasks,omitempty"`
	MissingTasks   []string `json:"missing_tasks,omitempty"`
}

type Report struct {
	Version     int               `json:"version"`
	CompletedBy string            `json:"completed_by"`
	Complete    bool              `json:"complete"`
	Wrote       []string          `json:"wrote,omitempty"`
	Removed     []string          `json:"removed,omitempty"`
	BackedUp    []string          `json:"backed_up,omitempty"`
	Sources     map[string]string `json:"sources"`
	Taskfiles   []Taskfile        `json:"taskfiles,omitempty"`
	Diagnostics []Diagnostic      `json:"diagnostics,omitempty"`
}

func Path(root string) string {
	return filepath.Join(root, filepath.FromSlash(RelativeReportPath))
}

// Read returns the report and whether it exists. An invalid report is an
// error rather than an implicit activation decision.
func Read(root string) (Report, bool, error) {
	data, err := os.ReadFile(Path(root))
	if errors.Is(err, fs.ErrNotExist) {
		return Report{}, false, nil
	}
	if err != nil {
		return Report{}, false, fmt.Errorf("read %s: %w", RelativeReportPath, err)
	}
	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return Report{}, true, fmt.Errorf("parse %s: %w", RelativeReportPath, err)
	}
	if report.Version < 1 {
		return Report{}, true, fmt.Errorf("parse %s: unsupported report version %d", RelativeReportPath, report.Version)
	}
	return report, true, nil
}

func Write(root string, report Report) error {
	path := Path(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".migration-report-*.json")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}
