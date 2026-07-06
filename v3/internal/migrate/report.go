package migrate

import (
	"fmt"
	"sort"
	"strings"
)

// Report collects everything the migrator did (or could not do) so the user
// gets a single MIGRATION.md summarising the state of the new project.
type Report struct {
	mapped []string          // options that were migrated automatically
	manual map[string]string // v2 option/feature -> what the user must do
	notes  []string          // informational notes
}

func NewReport() *Report {
	return &Report{manual: map[string]string{}}
}

// Mapped records a successfully migrated v2 option: e.g. Mapped("Title",
// "WebviewWindowOptions.Title").
func (r *Report) Mapped(v2Option, target string) {
	r.mapped = append(r.mapped, fmt.Sprintf("| `%s` | `%s` |", v2Option, target))
}

// Manual records something that needs the user's attention.
func (r *Report) Manual(what, instructions string) {
	r.manual[what] = instructions
}

// Note records an informational message.
func (r *Report) Note(note string) {
	r.notes = append(r.notes, note)
}

// HasManualSteps reports whether any manual step was recorded.
func (r *Report) HasManualSteps() bool {
	return len(r.manual) > 0
}

// Markdown renders the report as the MIGRATION.md contents.
func (r *Report) Markdown() string {
	var sb strings.Builder
	sb.WriteString("# Migration Report\n\n")
	sb.WriteString("This project was migrated from Wails v2 by `wails3 migrate`.\n\n")
	sb.WriteString("## Next steps\n\n")
	sb.WriteString("1. Run `wails3 doctor` to check your environment.\n")
	sb.WriteString("2. Run `wails3 dev` to build and run the migrated app.\n")
	sb.WriteString("3. Work through the *Manual steps* below, if any.\n")
	sb.WriteString("4. Migrate incrementally from the `v2compat/runtime` bridge to the v3 API.\n")
	sb.WriteString("   Every function in `github.com/wailsapp/wails/v3/pkg/v2compat/runtime` documents its v3 replacement.\n")
	sb.WriteString("   See https://v3.wails.io/migration/v2-to-v3/ for the full guide.\n\n")

	if len(r.manual) > 0 {
		sb.WriteString("## Manual steps\n\n")
		keys := make([]string, 0, len(r.manual))
		for k := range r.manual {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", k, r.manual[k]))
		}
		sb.WriteString("\n")
	}

	if len(r.notes) > 0 {
		sb.WriteString("## Notes\n\n")
		for _, n := range r.notes {
			sb.WriteString("- " + n + "\n")
		}
		sb.WriteString("\n")
	}

	if len(r.mapped) > 0 {
		sb.WriteString("## Migrated options\n\n")
		sb.WriteString("| Wails v2 | Wails v3 |\n|---|---|\n")
		mapped := append([]string(nil), r.mapped...)
		sort.Strings(mapped)
		for _, m := range mapped {
			sb.WriteString(m + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
