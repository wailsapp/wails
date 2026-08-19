// Package migration owns the ephemeral analysis types printed by `wails3
// migrate`. Migration has no durable state: routing depends only on wails.hcl.
package migration

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
	Sources     map[string]string `json:"sources"`
	Taskfiles   []Taskfile        `json:"taskfiles,omitempty"`
	Diagnostics []Diagnostic      `json:"diagnostics,omitempty"`
}
