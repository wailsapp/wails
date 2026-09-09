package setupwizard

import (
	"encoding/json"
	"net/http"
)

// InitTemplate is a selectable project template, sourced from the templates
// package by the caller (the setupwizard package must not import templates).
type InitTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// InitData is the payload exchanged with the init wizard frontend. The caller
// (commands.Init) seeds it with defaults + the available templates, and gets it
// back with the user's choices after Create.
type InitData struct {
	// Mode lets the frontend distinguish the init wizard from the setup wizard.
	Mode string `json:"mode"`

	ProjectName       string `json:"projectName"`
	TemplateName      string `json:"templateName"`
	ProductName       string `json:"productName"`
	ProductCompany    string `json:"productCompany"`
	ProductIdentifier string `json:"productIdentifier"`
	ProductDescription string `json:"productDescription"`
	ProductVersion    string `json:"productVersion"`
	ProductCopyright  string `json:"productCopyright"`
	ProductComments   string `json:"productComments"`
	// UseInterfaces selects interface vs class bindings for TypeScript projects.
	UseInterfaces bool `json:"useInterfaces"`

	// BaseDir is the absolute directory the project folder will be created in
	// (display only). The project lands at BaseDir/<ProjectName>.
	BaseDir string `json:"baseDir"`

	// Catalogue for the picker (not modified by the user).
	Templates       []InitTemplate `json:"templates"`
	DefaultTemplate string         `json:"defaultTemplate"`
}

// NewInitWizard creates a wizard that runs in project-init mode, pre-seeded with
// the given data (defaults + template list).
func NewInitWizard(data InitData) *Wizard {
	d := data
	d.Mode = "init"
	w := New()
	w.initData = &d
	return w
}

// RunInit launches the init wizard and blocks until the user clicks Create (or
// closes the window). On Create it returns the user-edited data; on close it
// returns (nil, nil).
func (w *Wizard) RunInit() (*InitData, error) {
	if err := w.Run(); err != nil {
		return nil, err
	}
	return w.initResult, nil
}

// handleInit returns the seeded init data (and signals init mode). In setup mode
// it returns null so the frontend falls back to the standard OOBE flow.
func (w *Wizard) handleInit(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if w.initData == nil {
		_, _ = rw.Write([]byte("null"))
		return
	}
	_ = json.NewEncoder(rw).Encode(w.initData)
}

// handleInitCreate records the user's final choices and unblocks RunInit.
func (w *Wizard) handleInitCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if w.initData == nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "not in init mode"})
		return
	}

	var result InitData
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Invalid request: " + err.Error()})
		return
	}
	result.Mode = "init"

	// Single-shot: only the first create request records the result and unblocks
	// RunInit. Concurrent requests are no-ops (no result race, no double close).
	w.doneOnce.Do(func() {
		w.initResult = &result
		close(w.done)
	})

	json.NewEncoder(rw).Encode(map[string]any{"success": true})
}
