package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/term"
)

// migrator collects everything the migration produces so the summary and
// MIGRATION_REPORT.md can be rendered at the end.
type migrator struct {
	projectDir string
	cfg        *v2ProjectConfig
	rawConfig  map[string]any

	main         *v2MainInfo
	runtimeCalls []string

	created      []string    // files written by the migrator
	skipped      []string    // files preserved because they already exist
	configMapped [][2]string // v2 wails.json key -> v3 home
	configManual []string    // config items needing manual attention
	mainManual   []string    // main.go items needing manual attention
}

func (m *migrator) wasCreated(name string) bool {
	for _, created := range m.created {
		if created == name {
			return true
		}
	}
	return false
}

// Migrate assists with migrating a Wails v2 project to v3. It maps the v2
// wails.json onto the v3 project layout (Taskfile.yml + build assets),
// converts the declarative wails.Run options into a programmatic v3
// bootstrap (main_v3.go.example) with Bind entries converted to Services,
// and writes a MIGRATION_REPORT.md describing what was automated and what
// needs manual work. It never deletes or overwrites existing project files.
func Migrate(options *flags.Migrate) error {
	if options.Quiet {
		term.DisableOutput()
		defer term.EnableOutput()
	}

	projectDir, err := filepath.Abs(options.ProjectDir)
	if err != nil {
		return err
	}
	if err := validateV2Project(projectDir); err != nil {
		return err
	}

	term.Header("Wails v2 to v3 Migration Assistant")
	term.Section("Analysing project")

	cfg, raw, err := loadV2Config(projectDir)
	if err != nil {
		return err
	}
	m := &migrator{
		projectDir: projectDir,
		cfg:        cfg,
		rawConfig:  raw,
	}
	term.Infof("Project: %s (v2)\n", cfg.Name)

	term.Section("Generating v3 configuration")
	if err := m.migrateConfig(); err != nil {
		return err
	}

	term.Section("Generating v3 bootstrap")
	m.migrateMain()

	if err := m.writeReport(); err != nil {
		return err
	}
	m.created = append(m.created, migrationReportName)

	m.printSummary()
	return nil
}

// validateV2Project checks that the directory looks like a Wails v2 project:
// a wails.json plus a go.mod requiring github.com/wailsapp/wails/v2.
func validateV2Project(dir string) error {
	explain := func(reason string) error {
		return fmt.Errorf(
			"%s does not look like a Wails v2 project: %s.\n"+
				"`wails3 migrate` needs a project with a `wails.json` and a `go.mod` that requires %s.\n"+
				"To start a fresh v3 project instead, run `wails3 init`", dir, reason, v2ModulePath)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return fmt.Errorf("project directory %s not found", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "wails.json")); err != nil {
		return explain("no wails.json found")
	}
	gomod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return explain("no go.mod found")
	}
	if !strings.Contains(string(gomod), v2ModulePath) {
		return explain("go.mod does not require " + v2ModulePath)
	}
	return nil
}

func (m *migrator) printSummary() {
	term.Section("Summary")

	// Collapse the generated build assets into a single line; the full list
	// is in the report.
	buildFiles := 0
	var created []string
	for _, file := range m.created {
		if strings.HasPrefix(file, "build/") {
			buildFiles++
			continue
		}
		created = append(created, file)
	}
	if buildFiles > 0 {
		created = append([]string{fmt.Sprintf("build/ (v3 build system, %d files)", buildFiles)}, created...)
	}
	for _, file := range created {
		term.Println("  created    " + file)
	}
	if len(m.skipped) > 0 {
		term.Println(fmt.Sprintf("  preserved  %d existing file(s), listed in %s", len(m.skipped), migrationReportName))
	}
	term.Println("")

	if m.main != nil {
		term.Success(fmt.Sprintf("Generated %s with %d service(s) from your Bind list.", generatedMainName, len(m.main.BindExprs)))
	}
	manualItems := len(m.configManual) + len(m.mainManual)
	if manualItems > 0 {
		term.Warningf("%d item(s) need manual attention. See %s for details.\n", manualItems, migrationReportName)
	} else {
		term.Success("No manual follow-ups detected. See " + migrationReportName + " for next steps.")
	}
	term.Println("")
	term.Println("Migration guide: " + term.Hyperlink("https://v3.wails.io/migration/v2-to-v3/", "https://v3.wails.io/migration/v2-to-v3/"))
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
