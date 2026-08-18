package commands

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/version"
	wakeast "github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/migration"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/zeebo/blake3"
	"gopkg.in/yaml.v3"
)

type MigrateOptions struct {
	DryRun   bool `name:"dry-run" description:"Analyse without writing project files"`
	JSON     bool `name:"json" description:"Print the machine-readable report"`
	Activate bool `name:"activate" description:"Validate a reviewed migration draft and atomically activate it"`
	Backup   bool `name:"backup" description:"Deprecated: migration never changes legacy files"`
	// Complete is retained as an undocumented compatibility spelling for one
	// transition. It has the same non-destructive semantics as --activate.
	Complete bool `name:"complete" description:"Deprecated: use --activate"`
}

type MigrationDiagnostic = migration.Diagnostic
type MigrationReport = migration.Report

func Migrate(options *MigrateOptions) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	if options.Activate || options.Complete {
		return activateMigration(root, options)
	}
	if manifest.Exists(root) {
		return fmt.Errorf("%s already exists; migration will not overwrite the active manifest", manifest.Filename)
	}
	draft := filepath.Join(root, manifest.MigratedFilename)
	if _, err := os.Stat(draft); err == nil {
		return fmt.Errorf("%s already exists; review it or remove it before running migration again", manifest.MigratedFilename)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	report, doc, err := analyseMigration(root)
	if err != nil {
		return err
	}
	if !options.DryRun {
		if err := manifest.WriteMigrationDraft(root, doc); err != nil {
			return err
		}
		report.Wrote = append(report.Wrote, manifest.MigratedFilename)
		report.Wrote = append(report.Wrote, migration.RelativeReportPath)
		if err := migration.Write(root, report); err != nil {
			return err
		}
	}
	if options.JSON {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	status := "complete"
	if !report.Complete {
		status = "needs manual changes"
	}
	fmt.Printf("Migration analysis: %s (%d source files, %d diagnostics)\n", status, len(report.Sources), len(report.Diagnostics))
	for _, d := range report.Diagnostics {
		location := d.File
		if d.Task != "" {
			location += " [" + d.Task + "]"
		}
		fmt.Printf("  %s %s: %s\n", strings.ToUpper(d.Severity), location, d.Message)
	}
	if options.DryRun {
		fmt.Println("No files written (--dry-run)")
	} else {
		fmt.Printf("Wrote inactive %s and %s\n", manifest.MigratedFilename, migration.RelativeReportPath)
	}
	return nil
}

func activateMigration(root string, options *MigrateOptions) error {
	draft := filepath.Join(root, manifest.MigratedFilename)
	if _, err := manifest.LoadFile(root, draft, ""); err != nil {
		return fmt.Errorf("validate %s: %w", manifest.MigratedFilename, err)
	}
	report, _, err := analyseMigration(root)
	if err != nil {
		return err
	}
	if !report.Complete {
		if options.JSON {
			if err := printMigrationJSON(report); err != nil {
				return err
			}
		}
		return fmt.Errorf("migration has unresolved blockers; review %s and rerun after resolving them", migration.RelativeReportPath)
	}
	if options.DryRun {
		if options.JSON {
			return printMigrationJSON(report)
		}
		fmt.Printf("Migration can be activated; no files written (--dry-run)\n")
		return nil
	}
	active := filepath.Join(root, manifest.Filename)
	if err := os.Link(draft, active); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%s already exists; refusing to overwrite it", manifest.Filename)
		}
		return fmt.Errorf("activate migration: %w", err)
	}
	if err := os.Remove(draft); err != nil {
		_ = os.Remove(active)
		return fmt.Errorf("activate migration: %w", err)
	}
	report.CompletedBy = version.String()
	report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "activated", Message: "reviewed migration draft was activated; legacy Taskfiles remain untouched"})
	if err := migration.Write(root, report); err != nil {
		return err
	}
	if options.JSON {
		return printMigrationJSON(report)
	}
	fmt.Printf("Activated %s; legacy Taskfiles are now ignored but were not changed\n", manifest.Filename)
	return nil
}

func printMigrationJSON(report MigrationReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func analyseMigration(root string) (MigrationReport, manifest.Document, error) {
	report := MigrationReport{Version: 1, CompletedBy: version.String(), Complete: true, Sources: map[string]string{}}
	root, err := filepath.Abs(root)
	if err != nil {
		return report, manifest.Document{}, err
	}
	rootTask, err := findTaskfile(root)
	if err != nil {
		return report, manifest.Document{}, err
	}
	files, err := discoverTaskfiles(rootTask)
	if err != nil {
		return report, manifest.Document{}, err
	}
	packageManager := ""
	binaryName := ""
	for _, path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		rel = filepath.ToSlash(rel)
		if rel == ".." || strings.HasPrefix(rel, "../") {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "external-taskfile", File: rel, Message: "included Taskfile is outside the project and must be migrated manually"})
			continue
		}
		digest, err := digestFile(path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		report.Sources[rel] = digest
		tf, err := parse.Parse(path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		role := taskfileRole(rel)
		allowed := legacyTaskNames[role]
		knownStock := knownStockTaskfiles[role][digest]
		classification := migration.Taskfile{File: rel, Role: role}
		var canonical canonicalDiff
		switch {
		case knownStock:
			// Shipped defaults changed throughout the v3 previews. Exact
			// fingerprints from generated projects are stock, even when their
			// task AST differs from the defaults compiled into this CLI.
			classification.Classification = "historical-default"
		case role == "root":
			classification.Classification = "current-default"
		case role == "unknown":
			classification.Classification = "custom"
		default:
			canonical = closestCanonical(tf.Tasks, stockTaskVariants(role))
			if canonical.exact {
				classification.Classification = "current-default"
			} else {
				classification.Classification = "customised"
				classification.ModifiedTasks = append(classification.ModifiedTasks, canonical.changed...)
				classification.ModifiedTasks = append(classification.ModifiedTasks, canonical.added...)
				classification.MissingTasks = append(classification.MissingTasks, canonical.missing...)
			}
		}
		taskNames := make([]string, 0, len(tf.Tasks))
		for name := range tf.Tasks {
			taskNames = append(taskNames, name)
		}
		sort.Strings(taskNames)
		for _, name := range taskNames {
			if phase, _, ok := migrateScriptHook(root, filepath.Dir(path), name, tf.Tasks[name]); ok {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "deferred-hook", File: rel, Task: name, Message: "custom " + phase + " script is not representable in config-only HCL and blocks activation"})
				if classification.Classification == "current-default" {
					classification.Classification = "customised"
				}
				classification.ModifiedTasks = appendUnique(classification.ModifiedTasks, name)
				continue
			}
			if !allowed[name] {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "unsupported-task", File: rel, Task: name, Message: "custom task requires a user-owned hook script or manual migration"})
				if classification.Classification == "current-default" {
					classification.Classification = "customised"
				}
				classification.ModifiedTasks = appendUnique(classification.ModifiedTasks, name)
				continue
			}
			if containsString(canonical.changed, name) || containsString(canonical.added, name) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "generated task was modified and requires a manifest field or user-owned hook script"})
			}
			if role == "root" && !recognizedRootTask(name, tf.Tasks[name]) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "root dispatch task was modified and requires a manifest field or user-owned hook script"})
				classification.Classification = "customised"
				classification.ModifiedTasks = appendUnique(classification.ModifiedTasks, name)
			}
		}
		if len(canonical.missing) > 0 {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "missing-generated-tasks", File: rel, Message: fmt.Sprintf("%d generated tasks are missing; review this Taskfile as customised", len(canonical.missing))})
		}
		sort.Strings(classification.ModifiedTasks)
		sort.Strings(classification.MissingTasks)
		report.Taskfiles = append(report.Taskfiles, classification)
		if role == "root" {
			for key, value := range tf.Vars {
				if key == "PACKAGE_MANAGER" && value.Static != "" && !strings.Contains(value.Static, "{{") {
					packageManager = value.Static
					report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "package-manager", File: rel, Message: "translated PACKAGE_MANAGER=" + value.Static})
				}
				if key == "APP_NAME" && value.Static != "" && !strings.Contains(value.Static, "{{") {
					binaryName = normalizeLegacyAppName(value.Static)
				}
			}
		}
	}
	project := manifest.Project{Name: filepath.Base(root), ProductName: filepath.Base(root), Identifier: "com.wails." + normaliseName(filepath.Base(root)), Version: "0.1.0"}
	project.BinaryName = binaryName
	if binaryName != "" {
		project.Name = binaryName
	}
	var associations []manifest.Association
	var protocols []manifest.Protocol
	configPath := filepath.Join(root, "build", "config.yml")
	if data, err := os.ReadFile(configPath); err == nil {
		var legacy struct {
			WailsConfig `yaml:",inline"`
			DevMode     struct {
				Executes []struct {
					Cmd  string `yaml:"cmd"`
					Type string `yaml:"type"`
				} `yaml:"executes"`
			} `yaml:"dev_mode"`
		}
		if err := yaml.Unmarshal(data, &legacy); err != nil {
			return report, manifest.Document{}, fmt.Errorf("parse build/config.yml: %w", err)
		}
		project.ProductName = first(legacy.Info.ProductName, project.ProductName)
		project.Identifier = first(legacy.Info.ProductIdentifier, project.Identifier)
		project.Version = first(legacy.Info.Version, project.Version)
		project.CompanyName = legacy.Info.CompanyName
		project.Description = legacy.Info.Description
		project.Copyright = legacy.Info.Copyright
		project.Comments = legacy.Info.Comments
		for _, file := range legacy.FileAssociations {
			associations = append(associations, manifest.Association{Extensions: []string{file.Ext}, Name: file.Name, Description: file.Description, Icon: file.IconName, Role: file.Role, MIMEType: file.MimeType})
		}
		for _, protocol := range legacy.Protocols {
			protocols = append(protocols, manifest.Protocol{Scheme: protocol.Scheme, Description: protocol.Description})
		}
		for _, execute := range legacy.DevMode.Executes {
			if !legacyDevCommand(execute.Cmd) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "unsupported-dev-command", File: "build/config.yml", Message: "dev command requires a user-owned hook or manual migration: " + execute.Cmd})
			}
		}
		digest, _ := digestFile(configPath)
		report.Sources["build/config.yml"] = digest
	} else if !errors.Is(err, fs.ErrNotExist) {
		return report, manifest.Document{}, err
	}
	doc := manifest.NewDocument(project)
	doc.Associations = associations
	doc.Protocols = protocols
	if packageManager != "" {
		doc.Frontend.PackageManager = packageManager
		doc.Frontend.Install = []string{packageManager, "install"}
		doc.Frontend.Build = []string{packageManager, "run", "build"}
		doc.Frontend.Dev = []string{packageManager, "run", "dev"}
	}
	return report, doc, nil
}

func normalizeLegacyAppName(value string) string {
	// Older generated root Taskfiles escaped spaces because APP_NAME was
	// interpolated into unquoted shell commands. The backslash is command
	// syntax, not part of the executable's file name.
	return strings.ReplaceAll(value, `\ `, " ")
}

func migrateScriptHook(root, taskfileDir, name string, task *wakeast.Task) (string, manifest.Hook, bool) {
	phase := strings.ReplaceAll(strings.ToLower(name), "-", "_")
	if !containsString([]string{"before_build", "after_build", "before_package", "after_package", "before_sign", "after_sign"}, phase) {
		return "", manifest.Hook{}, false
	}
	if task == nil || len(task.Cmds) != 1 || len(task.Deps) != 0 || task.Cmds[0].Cmd == "" || task.Cmds[0].Task != "" {
		return "", manifest.Hook{}, false
	}
	command := strings.TrimSpace(task.Cmds[0].Cmd)
	if strings.ContainsAny(command, " \t\r\n;&|`$<>") || strings.Contains(command, "{{") {
		return "", manifest.Hook{}, false
	}
	base := taskfileDir
	if task.Dir != "" {
		base = filepath.Join(base, filepath.FromSlash(task.Dir))
	}
	script := filepath.Clean(filepath.Join(base, filepath.FromSlash(strings.TrimPrefix(command, "./"))))
	relative, err := filepath.Rel(root, script)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", manifest.Hook{}, false
	}
	info, err := os.Stat(script)
	if err != nil || info.IsDir() {
		return "", manifest.Hook{}, false
	}
	return phase, manifest.Hook{Script: filepath.ToSlash(relative)}, true
}

func setMigratedHook(hooks *manifest.Hooks, phase string, hook manifest.Hook) {
	switch phase {
	case "before_build":
		hooks.BeforeBuild = hook
	case "after_build":
		hooks.AfterBuild = hook
	case "before_package":
		hooks.BeforePackage = hook
	case "after_package":
		hooks.AfterPackage = hook
	case "before_sign":
		hooks.BeforeSign = hook
	case "after_sign":
		hooks.AfterSign = hook
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func legacyDevCommand(command string) bool {
	command = strings.TrimSpace(command)
	if command == "go mod tidy" || command == "wails3 build DEV=true" {
		return true
	}
	return strings.HasPrefix(command, "wails3 task common:install:frontend:deps") || strings.HasPrefix(command, "wails3 task common:dev:frontend") || strings.HasPrefix(command, "wails3 task build") || strings.HasPrefix(command, "wails3 task run")
}

func discoverTaskfiles(root string) ([]string, error) {
	seen := map[string]bool{}
	var result []string
	var visit func(string) error
	visit = func(path string) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if seen[abs] {
			return nil
		}
		seen[abs] = true
		tf, err := parse.Parse(abs)
		if err != nil {
			return err
		}
		result = append(result, abs)
		base := filepath.Dir(abs)
		names := make([]string, 0, len(tf.Includes))
		for name := range tf.Includes {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			inc := tf.Includes[name]
			path := inc.Taskfile
			if !filepath.IsAbs(path) {
				path = filepath.Join(base, path)
			}
			info, err := os.Stat(path)
			if err == nil && info.IsDir() {
				yml := filepath.Join(path, "Taskfile.yml")
				yaml := filepath.Join(path, "Taskfile.yaml")
				if _, ymlErr := os.Stat(yml); ymlErr == nil {
					path = yml
				} else {
					path = yaml
				}
				_, err = os.Stat(path)
			}
			if errors.Is(err, fs.ErrNotExist) {
				if inc.Optional {
					continue
				}
				return fmt.Errorf("include %q points at a missing Taskfile: %s", name, path)
			}
			if err != nil {
				return err
			}
			if err := visit(path); err != nil {
				return err
			}
		}
		return nil
	}
	if err := visit(root); err != nil {
		return nil, err
	}
	projectRoot := filepath.Dir(root)
	for _, name := range []string{"Taskfile.override.yml", "Taskfile.override.yaml", "Taskfile.local.yml", "Taskfile.local.yaml"} {
		path := filepath.Join(projectRoot, name)
		if _, err := os.Stat(path); err == nil {
			if err := visit(path); err != nil {
				return nil, err
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

type canonicalDiff struct {
	exact   bool
	changed []string
	missing []string
	added   []string
}

func stockTaskVariants(role string) []map[string]*wakeast.Task {
	if role == "root" || role == "unknown" {
		return nil
	}
	if role == "common" {
		data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
		if err != nil {
			return nil
		}
		var result []map[string]*wakeast.Task
		for _, typescript := range []bool{false, true} {
			for _, useInterfaces := range []bool{false, true} {
				config := BuildConfig{
					BuildAssetsOptions: BuildAssetsOptions{Typescript: typescript, UseInterfaces: useInterfaces},
					TemplateEnrichment: TemplateEnrichment{Opn: "{{", Cls: "}}"},
				}
				tmpl, err := template.New("Taskfile.yml").Parse(string(data))
				if err != nil {
					return nil
				}
				var rendered bytes.Buffer
				if err := tmpl.Execute(&rendered, config); err != nil {
					return nil
				}
				if tasks := parseStockTasks(rendered.Bytes()); tasks != nil {
					result = append(result, tasks)
				}
			}
		}
		return result
	}
	data, err := buildAssets.ReadFile(filepath.ToSlash(filepath.Join("build_assets", role, "Taskfile.yml")))
	if err != nil {
		return nil
	}
	if tasks := parseStockTasks(data); tasks != nil {
		return []map[string]*wakeast.Task{tasks}
	}
	return nil
}

func parseStockTasks(data []byte) map[string]*wakeast.Task {
	tmp, err := os.CreateTemp("", "wails-stock-taskfile-*.yml")
	if err != nil {
		return nil
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return nil
	}
	if err := tmp.Close(); err != nil {
		return nil
	}
	tf, err := parse.Parse(name)
	if err != nil {
		return nil
	}
	return tf.Tasks
}

func closestCanonical(actual map[string]*wakeast.Task, variants []map[string]*wakeast.Task) canonicalDiff {
	best := canonicalDiff{added: sortedTaskNames(actual)}
	bestScore := len(best.added)
	if len(variants) == 0 {
		return best
	}
	bestScore = int(^uint(0) >> 1)
	for _, expected := range variants {
		var candidate canonicalDiff
		for name, task := range actual {
			expectedTask, ok := expected[name]
			switch {
			case !ok:
				candidate.added = append(candidate.added, name)
			case !reflect.DeepEqual(task, expectedTask):
				candidate.changed = append(candidate.changed, name)
			}
		}
		for name := range expected {
			if _, ok := actual[name]; !ok {
				candidate.missing = append(candidate.missing, name)
			}
		}
		sort.Strings(candidate.changed)
		sort.Strings(candidate.missing)
		sort.Strings(candidate.added)
		score := len(candidate.changed) + len(candidate.missing) + len(candidate.added)
		candidate.exact = score == 0
		if score < bestScore {
			best = candidate
			bestScore = score
		}
	}
	return best
}

func sortedTaskNames(tasks map[string]*wakeast.Task) []string {
	result := make([]string, 0, len(tasks))
	for name := range tasks {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func appendUnique(values []string, value string) []string {
	if containsString(values, value) {
		return values
	}
	return append(values, value)
}
func findTaskfile(root string) (string, error) {
	for _, name := range []string{"Taskfile.yml", "Taskfile.yaml"} {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no legacy Taskfile.yml or Taskfile.yaml found")
}
func taskfileRole(rel string) string {
	rel = filepath.ToSlash(rel)
	switch {
	case rel == "Taskfile.yml" || rel == "Taskfile.yaml":
		return "root"
	case rel == "build/Taskfile.yml" || rel == "build/Taskfile.yaml":
		return "common"
	case strings.Contains(rel, "/windows/"):
		return "windows"
	case strings.Contains(rel, "/darwin/"):
		return "darwin"
	case strings.Contains(rel, "/linux/"):
		return "linux"
	case strings.Contains(rel, "/ios/"):
		return "ios"
	case strings.Contains(rel, "/android/"):
		return "android"
	default:
		return "unknown"
	}
}

var legacyTaskNames = map[string]map[string]bool{
	"root":    set("build", "package", "run", "dev", "setup:docker", "build:server", "run:server", "build:docker", "run:docker"),
	"common":  set("go:mod:tidy", "install:frontend:deps", "install:frontend:deps:npm", "install:frontend:deps:bun", "install:frontend:deps:pnpm", "install:frontend:deps:yarn", "build:frontend", "frontend:run", "frontend:run:npm", "frontend:run:yarn", "frontend:run:pnpm", "frontend:run:bun", "frontend:vendor:puppertino", "generate:bindings", "generate:icons", "dev:frontend", "frontend:dev:npm", "frontend:dev:yarn", "frontend:dev:pnpm", "frontend:dev:bun", "update:build-assets", "build:server", "run:server", "build:docker", "run:docker", "setup:docker", "ios:device:list", "ios:run:device"),
	"windows": set("build", "build:native", "build:docker", "package", "generate:syso", "create:nsis:installer", "create:msix:package", "install:msix:tools", "run", "sign", "sign:installer"),
	"darwin":  set("build", "build:native", "build:docker", "build:universal", "build:universal:lipo:native", "build:universal:lipo:go", "package", "package:universal", "package:dmg", "create:dmg", "create:app:bundle", "codesign:adhoc", "codesign:skip", "run", "sign", "sign:notarize"),
	"linux":   set("build", "build:native", "build:docker", "package", "create:appimage", "create:deb", "create:rpm", "create:aur", "generate:deb", "generate:rpm", "generate:aur", "generate:dotdesktop", "run", "sign:deb", "sign:rpm", "sign:packages"),
}

// Exact BLAKE3 fingerprints of generated platform Taskfiles shipped in v3
// examples before the manifest cutover. These distinguish historical defaults
// from user edits without weakening comparison against the current templates.
var knownStockTaskfiles = map[string]map[string]bool{
	"common": set(
		"4ab65b0363866b550ef897db8f7aae12794789056dd4f2e82407a738a64f6819",
		"64cef372c5d2ae21526043751993f57b6d5b469287da0adfc8a78dde0d7ad86a",
		"9a85cde94f9c2562240676b4306b1176158f0df8956d010e7b9b0d15b08a76b9",
		"371fafce90a2648d89e705e377ae06ade62292271a4baabf1806e838e5431a7d",
		"5e14a693d4e63dd8c548998502bea0d66798a2d8dbbbb95a251056051ff63636",
		"c7413c1fd1dcae671502a94715b6309f56def27f644e1921f4ca982431a24e42",
		"3ee740d29d727fd11b098cc4ee343cf2d9bc40c8940d6bce3bc9df181f74f549",
	),
	"windows": set(
		"2036534db056f76202df0ade442e697a4572643b985ed5bc3c0042f97de773ea",
		"ca9fb9e4d61c1816b6cc6919598feb54adc69451c444afe3af68588ed1d7d441",
		"ab1b5cec4dc3f7f42a8c222b3b87744405e9b1b20f520ab7d9ab44bfcca88013",
	),
	"darwin": set(
		"ed7dd8950f9422105e7d4349c223e0aeda1eca364186a621b4083052c9e4ca68",
		"9c030aea973144d4041859e6a82fe791a5a3b919d1fbe3549830b19f93cb25b2",
	),
	"linux": set(
		"5fed27274feb32c231426a77edd1f66fdf8fb77a0f074ab7db862dc0d3f31d6e",
		"0fde154fbc54fc44c64e6ff4cc9e55474f118c80d750a7f479dd7c8a45c5cfa5",
	),
	"ios": set(
		"b0839989d6583756f6a1d71886a1462b9983e70d6255672ead8dcca952b42268",
		"676f5bb6bee3a521907a8cd213c9a5b43653d2cd16631debb71328f16c5d03c4",
	),
	"android": set(
		"c36c2687b8ef10117ef7ea2218d26b64cea9d5ef1a8cdca25e5c3cf30964be2e",
		"1e17b34e61333f3037986d27dcec870fe64432c72b22e70dad7603fbde0b4ba3",
	),
}

func recognizedRootTask(name string, task *wakeast.Task) bool {
	if task == nil || len(task.Cmds) != 1 || len(task.Deps) != 0 {
		return false
	}
	command := task.Cmds[0]
	if command == nil {
		return false
	}
	if name == "dev" {
		return strings.HasPrefix(strings.TrimSpace(command.Cmd), "wails3 dev ") && command.Task == ""
	}
	refs := map[string]string{
		"setup:docker": "common:setup:docker", "build:server": "common:build:server",
		"run:server": "common:run:server", "build:docker": "common:build:docker", "run:docker": "common:run:docker",
	}
	if expected := refs[name]; expected != "" {
		return command.Task == expected && command.Cmd == ""
	}
	if name == "build" || name == "package" || name == "run" {
		return command.Cmd == "" && strings.HasSuffix(command.Task, ":"+name) && strings.Contains(command.Task, "{{")
	}
	return false
}

func init() {
	legacyTaskNames["ios"] = set("install:deps", "build", "package", "package:ipa", "compile:ios:link", "create:app:bundle", "deploy-simulator", "deploy-device", "compile:ios", "generate:ios:bindings", "ensure-simulator", "generate:ios:overlay", "generate:ios:xcode", "run", "xcode", "logs", "logs:dev", "logs:wide")
	legacyTaskNames["android"] = set("install:deps", "build", "compile:go:shared", "compile:go:all-archs", "package", "package:fat", "bundle", "bundle:fat", "assemble:apk", "assemble:apk:release", "assemble:aab", "assemble:aab:release", "generate:android:overlay", "generate:android:bindings", "ensure-emulator", "deploy-emulator", "run", "device:list", "run:device", "deploy-device", "studio", "logs", "logs:all", "clean")
}
func set(values ...string) map[string]bool {
	result := map[string]bool{}
	for _, v := range values {
		result[v] = true
	}
	return result
}

func removeLegacySources(root string, sources map[string]string, complete bool) ([]string, []MigrationDiagnostic) {
	if !complete {
		return nil, []MigrationDiagnostic{{Severity: "warning", Code: "remove-blocked", Message: "legacy files were retained because migration is incomplete"}}
	}
	var removed []string
	var diagnostics []MigrationDiagnostic
	paths := make([]string, 0, len(sources))
	for path := range sources {
		paths = append(paths, path)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, rel := range paths {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
		if filepath.IsAbs(rel) || clean == ".." || strings.HasPrefix(clean, "../") {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "remove-outside-project", File: rel, Message: "retained because the source is outside the project"})
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(rel))
		digest, err := digestFile(path)
		if err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "remove-failed", File: rel, Message: err.Error()})
			continue
		}
		if digest != sources[rel] {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-source", File: rel, Message: "retained because it changed after analysis"})
			continue
		}
		if err := os.Remove(path); err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "remove-failed", File: rel, Message: err.Error()})
			continue
		}
		removed = append(removed, rel)
	}
	sort.Strings(removed)
	return removed, diagnostics
}

func verifyLegacySources(root string, sources map[string]string) []MigrationDiagnostic {
	var diagnostics []MigrationDiagnostic
	paths := make([]string, 0, len(sources))
	for path := range sources {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, rel := range paths {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
		if filepath.IsAbs(rel) || clean == ".." || strings.HasPrefix(clean, "../") {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "source-outside-project", File: rel, Message: "source is outside the project"})
			continue
		}
		digest, err := digestFile(filepath.Join(root, filepath.FromSlash(clean)))
		if err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "source-unavailable", File: rel, Message: err.Error()})
			continue
		}
		if digest != sources[rel] {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-source", File: rel, Message: "source changed after the original migration analysis"})
		}
	}
	return diagnostics
}

func backupLegacySources(root string, sources map[string]string) ([]string, []MigrationDiagnostic) {
	var backedUp []string
	var diagnostics []MigrationDiagnostic
	paths := make([]string, 0, len(sources))
	for path := range sources {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, rel := range paths {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
		if filepath.IsAbs(rel) || clean == ".." || strings.HasPrefix(clean, "../") {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-outside-project", File: rel, Message: "not backed up because the source is outside the project"})
			continue
		}
		source := filepath.Join(root, filepath.FromSlash(clean))
		digest, err := digestFile(source)
		if err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-failed", File: rel, Message: err.Error()})
			continue
		}
		if digest != sources[rel] {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-source", File: rel, Message: "not backed up because it changed after analysis"})
			continue
		}
		data, err := os.ReadFile(source)
		if err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-failed", File: rel, Message: err.Error()})
			continue
		}
		info, err := os.Stat(source)
		if err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-failed", File: rel, Message: err.Error()})
			continue
		}
		destination := filepath.Join(root, ".wails", "migration-backup", filepath.FromSlash(clean))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-failed", File: rel, Message: err.Error()})
			continue
		}
		if err := os.WriteFile(destination, data, info.Mode().Perm()); err != nil {
			diagnostics = append(diagnostics, MigrationDiagnostic{Severity: "warning", Code: "backup-failed", File: rel, Message: err.Error()})
			continue
		}
		backedUp = append(backedUp, filepath.ToSlash(filepath.Join(".wails", "migration-backup", clean)))
	}
	return backedUp, diagnostics
}

func digestFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := blake3.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
func first(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
