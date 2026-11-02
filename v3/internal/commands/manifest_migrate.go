package commands

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/version"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/migration"
	"github.com/wailsapp/wails/v3/internal/wake/migration/taskfile"
	wakeast "github.com/wailsapp/wails/v3/internal/wake/migration/taskfile/ast"
)

type MigrateOptions struct {
	DryRun   bool   `name:"dry-run" description:"Analyse without writing project files"`
	JSON     bool   `name:"json" description:"Print the machine-readable report"`
	Activate bool   `name:"activate" description:"Validate a reviewed migration draft and atomically activate it"`
	Output   string `name:"output" description:"Write or activate an inactive proposal at this project-relative path"`
}

type MigrationDiagnostic = migration.Diagnostic
type MigrationReport = migration.Report

func Migrate(options *MigrateOptions) error {
	return migrateWithOperations(options, productionMigrationCommandOperations())
}

type migrationCommandOperations struct {
	getwd                func() (string, error)
	activeManifestExists func(string) bool
	stat                 func(string) (fs.FileInfo, error)
	analyse              func(string) (MigrationReport, manifest.Document, error)
	writeDraft           func(string, string, manifest.Document, []string) error
	loadDraft            func(string, string, string) (*manifest.Loaded, error)
	publishDraft         func(string, string) error
	version              func() string
}

func productionMigrationCommandOperations() migrationCommandOperations {
	return migrationCommandOperations{
		getwd:                os.Getwd,
		activeManifestExists: manifest.Exists,
		stat:                 os.Stat,
		analyse:              analyseMigration,
		writeDraft:           manifest.WriteMigrationDraftAt,
		loadDraft:            manifest.LoadFile,
		publishDraft: func(draft, active string) error {
			return activateMigrationDraft(draft, active, migrationActivationOperations{link: os.Link, remove: os.Remove})
		},
		version: version.String,
	}
}

func migrateWithOperations(options *MigrateOptions, operations migrationCommandOperations) error {
	if options.JSON {
		// Keep stdout as one machine-readable document when the CLI's deferred
		// footer runs after this command returns.
		DisableFooter = true
	}
	root, err := operations.getwd()
	if err != nil {
		return err
	}
	if options.Activate {
		return activateMigrationWithOperations(root, options, operations)
	}
	if operations.activeManifestExists(root) {
		return fmt.Errorf("%s already exists; migration will not overwrite the active manifest", manifest.Filename)
	}
	draft, relativeDraft, err := migrationDraftPath(root, options.Output)
	if err != nil {
		return err
	}
	draftExists := false
	if _, err := operations.stat(draft); err == nil {
		draftExists = true
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	report, doc, err := operations.analyse(root)
	if err != nil {
		return err
	}
	if !options.DryRun && !draftExists {
		if err := operations.writeDraft(root, relativeDraft, doc, migrationBlockerComments(report)); err != nil {
			return err
		}
		report.Wrote = append(report.Wrote, relativeDraft)
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
	} else if draftExists {
		fmt.Printf("Left existing inactive %s unchanged; analysis only\n", relativeDraft)
	} else {
		fmt.Printf("Wrote inactive %s\n", relativeDraft)
	}
	return nil
}

func activateMigration(root string, options *MigrateOptions) error {
	return activateMigrationWithOperations(root, options, productionMigrationCommandOperations())
}

func activateMigrationWithOperations(root string, options *MigrateOptions, operations migrationCommandOperations) error {
	draft, relativeDraft, err := migrationDraftPath(root, options.Output)
	if err != nil {
		return err
	}
	if _, err := operations.loadDraft(root, draft, ""); err != nil {
		return fmt.Errorf("validate %s: %w", relativeDraft, err)
	}
	report, _, err := operations.analyse(root)
	if err != nil {
		return err
	}
	if !report.Complete {
		if options.JSON {
			printMigrationJSON(report)
		}
		return fmt.Errorf("migration has unresolved blockers; rerun `wails3 migrate` for current diagnostics")
	}
	if options.DryRun {
		if options.JSON {
			printMigrationJSON(report)
			return nil
		}
		fmt.Printf("Migration can be activated; no files written (--dry-run)\n")
		return nil
	}
	active := filepath.Join(root, manifest.Filename)
	if err := operations.publishDraft(draft, active); err != nil {
		return err
	}
	report.CompletedBy = operations.version()
	report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "activated", Message: "reviewed migration draft was activated; legacy Taskfiles remain untouched"})
	if options.JSON {
		printMigrationJSON(report)
		return nil
	}
	fmt.Printf("Activated %s; legacy Taskfiles are now ignored but were not changed\n", manifest.Filename)
	return nil
}

type migrationActivationOperations struct {
	link   func(string, string) error
	remove func(string) error
}

func activateMigrationDraft(draft, active string, operations migrationActivationOperations) error {
	if err := operations.link(draft, active); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%s already exists; refusing to overwrite it", manifest.Filename)
		}
		return fmt.Errorf("activate migration: %w", err)
	}
	if err := operations.remove(draft); err != nil {
		activationErr := fmt.Errorf("activate migration: remove inactive draft: %w", err)
		if rollbackErr := operations.remove(active); rollbackErr != nil {
			return errors.Join(activationErr, fmt.Errorf("activate migration: roll back active manifest: %w", rollbackErr))
		}
		return activationErr
	}
	return nil
}

func migrationDraftPath(root, configured string) (string, string, error) {
	if configured == "" {
		configured = manifest.MigratedFilename
	}
	clean := filepath.ToSlash(filepath.Clean(configured))
	if strings.EqualFold(clean, manifest.Filename) || strings.EqualFold(clean, manifest.EjectedFilename) || clean == "." || strings.HasPrefix(strings.ToLower(clean), ".wails/") {
		return "", "", fmt.Errorf("migration output %q must be an inactive project-owned YAML file", configured)
	}
	path, err := manifest.ResolveProjectPath(root, "migration output", clean, false)
	if err != nil {
		return "", "", err
	}
	if extension := strings.ToLower(filepath.Ext(path)); extension != ".yaml" && extension != ".yml" {
		return "", "", fmt.Errorf("migration output %q must use the .yaml or .yml extension", configured)
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return "", "", err
	}
	return path, filepath.ToSlash(relative), nil
}

func migrationBlockerComments(report MigrationReport) []string {
	comments := []string{"Generated by wails3 migrate."}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.Severity != "warning" {
			continue
		}
		location := diagnostic.File
		if diagnostic.Task != "" {
			location += " [" + diagnostic.Task + "]"
		}
		comments = append(comments, "BLOCKED: "+location+": "+diagnostic.Message)
	}
	return comments
}

func printMigrationJSON(report MigrationReport) {
	// MigrationReport contains only JSON-safe scalar, slice, and struct fields,
	// so encoding cannot fail. Keep this helper infallible instead of carrying
	// an unreachable error branch through every activation outcome.
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
}

type migrationTaskNode struct {
	key, prefix, path, name string
	task                    *wakeast.Task
	stockContext            bool
}

type migrationReachability struct {
	reachable map[string]map[string]bool
	issues    []MigrationDiagnostic
}

func (r migrationReachability) contains(path, task string) bool {
	return r.reachable[filepath.Clean(path)][task]
}

func analyseTaskReachability(rootTask string, files []string) (migrationReachability, error) {
	nodes := map[string][]migrationTaskNode{}
	loadedPaths := map[string]bool{}
	visited := map[string]bool{}
	ancestors := map[string]bool{}
	var includeIssues []MigrationDiagnostic
	var load func(string, string, bool) error
	load = func(path, prefix string, inheritedStockContext bool) error {
		path = filepath.Clean(path)
		if ancestors[path] {
			includeIssues = append(includeIssues, MigrationDiagnostic{Severity: "warning", Code: "cyclic-include", File: path, Message: "Taskfile include cycle cannot be proven equivalent"})
			return nil
		}
		visitKey := path + "\x00" + prefix + "\x00" + strconv.FormatBool(inheritedStockContext)
		if visited[visitKey] {
			return nil
		}
		visited[visitKey] = true
		ancestors[path] = true
		defer delete(ancestors, path)
		loadedPaths[path] = true
		tf, err := taskfile.Parse(path)
		if err != nil {
			return err
		}
		stockContext := inheritedStockContext && len(tf.Unsupported) == 0 && tf.Run == "" && tf.Interval == "" && stockMigrationVariables(rootTask, path, tf.Vars)
		for name, task := range tf.Tasks {
			key := prefix + name
			nodes[key] = append(nodes[key], migrationTaskNode{key: key, prefix: prefix, path: path, name: name, task: task, stockContext: stockContext})
		}
		includeNames := make([]string, 0, len(tf.Includes))
		for name := range tf.Includes {
			includeNames = append(includeNames, name)
		}
		sort.Strings(includeNames)
		for _, name := range includeNames {
			if tf.Includes[name].Taskfile == "" || strings.Contains(tf.Includes[name].Taskfile, "{{") {
				includeIssues = append(includeIssues, MigrationDiagnostic{Severity: "warning", Code: "dynamic-include", File: path, Message: "templated include " + name + " cannot be proven equivalent"})
				continue
			}
			included, err := resolveMigrationInclude(path, tf.Includes[name])
			if errors.Is(err, fs.ErrNotExist) && tf.Includes[name].Optional {
				continue
			}
			if errors.Is(err, fs.ErrNotExist) {
				includeIssues = append(includeIssues, MigrationDiagnostic{Severity: "warning", Code: "missing-include", File: path, Message: "required include " + name + " is missing"})
				continue
			}
			if err != nil {
				return err
			}
			include := tf.Includes[name]
			includeContext := stockContext && include.Dir == "" && len(include.Vars) == 0 && len(include.Unsupported) == 0
			if err := load(included, prefix+name+":", includeContext); err != nil {
				return err
			}
		}
		return nil
	}
	if err := load(rootTask, "", true); err != nil {
		return migrationReachability{}, err
	}
	// Conventional local/override files are merged at the root namespace even
	// though they are not reached through an explicit include.
	for _, path := range files {
		path = filepath.Clean(path)
		if !loadedPaths[path] {
			if err := load(path, "", true); err != nil {
				return migrationReachability{}, err
			}
		}
	}

	result := migrationReachability{reachable: map[string]map[string]bool{}, issues: includeIssues}
	queue := make([]string, 0, 30)
	for _, root := range []string{"build", "package", "sign", "dev", "run"} {
		if len(nodes[root]) > 0 {
			queue = append(queue, root)
		}
		for _, platform := range []string{"windows", "darwin", "linux", "ios", "android"} {
			candidate := platform + ":" + root
			if len(nodes[candidate]) > 0 {
				queue = append(queue, candidate)
			}
		}
	}
	seen := map[string]bool{}
	issueSeen := map[string]bool{}
	addIssue := func(node migrationTaskNode, code, message string) {
		key := node.path + "\x00" + node.name + "\x00" + code + "\x00" + message
		if issueSeen[key] {
			return
		}
		issueSeen[key] = true
		result.issues = append(result.issues, MigrationDiagnostic{Severity: "warning", Code: code, File: node.path, Task: node.name, Message: message})
	}
	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]
		if seen[key] {
			continue
		}
		seen[key] = true
		for _, node := range nodes[key] {
			if result.reachable[node.path] == nil {
				result.reachable[node.path] = map[string]bool{}
			}
			result.reachable[node.path][node.name] = true
			for _, reference := range migrationTaskReferences(node.task) {
				references := []string{reference}
				stockSelector := strings.Contains(reference, "{{")
				if stockSelector {
					if node.prefix == "" && node.stockContext && recognizedRootTask(node.name, node.task) {
						suffix := ":" + node.name
						for candidate := range nodes {
							if strings.HasSuffix(candidate, suffix) {
								queue = append(queue, candidate)
							}
						}
						continue
					}
					references = stockDynamicTaskTargets(rootTask, node, reference)
					if len(references) == 0 {
						addIssue(node, "dynamic-task-reference", "reachable task selects another task dynamically and cannot be proven equivalent")
						continue
					}
				}
				for _, reference := range references {
					candidate := reference
					if node.prefix != "" && (stockSelector || len(nodes[candidate]) == 0) {
						candidate = node.prefix + reference
					}
					if len(nodes[candidate]) == 0 {
						addIssue(node, "unresolved-task-reference", "reachable task refers to missing task "+reference)
						continue
					}
					queue = append(queue, candidate)
				}
			}
		}
	}
	sort.Slice(result.issues, func(i, j int) bool {
		if result.issues[i].File != result.issues[j].File {
			return result.issues[i].File < result.issues[j].File
		}
		if result.issues[i].Task != result.issues[j].Task {
			return result.issues[i].Task < result.issues[j].Task
		}
		return result.issues[i].Code < result.issues[j].Code
	})
	return result, nil
}

func resolveMigrationInclude(taskfile string, include *wakeast.Include) (string, error) {
	path := include.Taskfile
	if path == "" || strings.Contains(path, "{{") {
		return "", fmt.Errorf("dynamic include %q in %s cannot be analysed", path, taskfile)
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(taskfile), filepath.FromSlash(path))
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return filepath.Clean(path), nil
	}
	for _, name := range []string{"Taskfile.yml", "Taskfile.yaml"} {
		candidate := filepath.Join(path, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	return "", fs.ErrNotExist
}

// Only root variables translated by migration and unchanged platform defaults
// may supply inputs to a stock selector. Custom variable evaluation stays blocked.
func stockMigrationVariables(rootTask, path string, variables map[string]*wakeast.Var) bool {
	rel, err := filepath.Rel(filepath.Dir(rootTask), path)
	if err != nil {
		return false
	}
	role := taskfileRole(rel)
	if role != "root" {
		if len(variables) == 0 {
			return true
		}
		for _, stock := range stockTaskfileVariants(role) {
			if reflect.DeepEqual(variables, stock.Vars) {
				return true
			}
		}
		return false
	}
	for name, value := range variables {
		if value == nil || !value.Literal {
			return false
		}
		switch name {
		case "APP_NAME", "BIN_DIR", "VITE_PORT", "PACKAGE_MANAGER":
			if !strings.Contains(value.Static, "{{") {
				continue
			}
		}
		if name == "PACKAGE_MANAGER" && value.Static == `{{.PACKAGE_MANAGER | default "npm"}}` ||
			name == "VITE_PORT" && value.Static == `{{.WAILS_VITE_PORT | default 9245}}` ||
			name == "GOOS" && value.Static == `{{.GOOS | default OS}}` {
			continue
		}
		return false
	}
	return true
}

// Stock selectors have a finite set of branches. Traverse all of them without
// evaluating Taskfile expressions, and only after matching the entire stock task.
func stockDynamicTaskTargets(rootTask string, node migrationTaskNode, reference string) []string {
	rel, err := filepath.Rel(filepath.Dir(rootTask), node.path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil
	}
	if !node.stockContext {
		return nil
	}
	role := taskfileRole(rel)
	matched := false
	for _, stock := range stockTaskVariants(role) {
		if expected, ok := stock[node.name]; ok && reflect.DeepEqual(node.task, expected) {
			matched = true
			break
		}
	}
	if !matched {
		return nil
	}
	switch reference {
	case "install:frontend:deps:{{.PACKAGE_MANAGER}}":
		return []string{"install:frontend:deps:npm", "install:frontend:deps:pnpm", "install:frontend:deps:yarn", "install:frontend:deps:bun"}
	case "frontend:run:{{.PACKAGE_MANAGER}}":
		return []string{"frontend:run:npm", "frontend:run:pnpm", "frontend:run:yarn", "frontend:run:bun"}
	case "frontend:dev:{{.PACKAGE_MANAGER}}":
		return []string{"frontend:dev:npm", "frontend:dev:pnpm", "frontend:dev:yarn", "frontend:dev:bun"}
	case `{{if and (ne OS "windows") (eq .CGO_ENABLED "1")}}build:docker{{else}}build:native{{end}}`,
		`{{if or (eq OS "darwin") (has "purego" (splitList "," (.EXTRA_TAGS | default "")))}}build:native{{else}}build:docker{{end}}`,
		`{{if or (ne .CGO_ENABLED "1") (and (eq OS "linux") (eq .HAS_CC "true") (eq .TARGET_ARCH ARCH))}}build:native{{else}}build:docker{{end}}`:
		return []string{"build:native", "build:docker"}
	case `{{if eq (.FORMAT | default "nsis") "msix"}}create:msix:package{{else}}create:nsis:installer{{end}}`:
		return []string{"create:msix:package", "create:nsis:installer"}
	case `{{if eq OS "darwin"}}build:universal:lipo:native{{else}}build:universal:lipo:go{{end}}`:
		return []string{"build:universal:lipo:native", "build:universal:lipo:go"}
	case `{{if eq OS "darwin"}}codesign:adhoc{{else}}codesign:skip{{end}}`:
		return []string{"codesign:adhoc", "codesign:skip"}
	}
	return nil
}

func migrationTaskReferences(task *wakeast.Task) []string {
	if task == nil {
		return nil
	}
	result := make([]string, 0, len(task.Deps)+len(task.Cmds))
	for _, dep := range task.Deps {
		if dep != nil && dep.Task != "" {
			result = append(result, dep.Task)
		}
	}
	for _, command := range task.Cmds {
		if command == nil {
			continue
		}
		if command.Task != "" {
			result = append(result, command.Task)
		}
		if command.For != nil && command.For.Task != "" {
			result = append(result, command.For.Task)
		}
	}
	return result
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
	reachability, err := analyseTaskReachability(rootTask, files)
	if err != nil {
		return report, manifest.Document{}, err
	}
	for _, diagnostic := range reachability.issues {
		if relative, relErr := filepath.Rel(root, diagnostic.File); relErr == nil {
			diagnostic.File = filepath.ToSlash(relative)
		}
		report.Complete = false
		report.Diagnostics = append(report.Diagnostics, diagnostic)
	}
	packageManager := ""
	binaryName := ""
	buildOutput := ""
	devPort := 0
	typescriptBindings := false
	interfaceBindings := false
	migratedHooks := make(map[manifest.HookPhase]manifest.Hook)
	for _, path := range files {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		rel = filepath.ToSlash(rel)
		if rel == ".." || strings.HasPrefix(rel, "../") {
			if len(reachability.reachable[filepath.Clean(path)]) > 0 {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "external-taskfile", File: rel, Message: "reachable included Taskfile is outside the project and must be migrated manually"})
			} else {
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "unrelated-taskfile", File: rel, Message: "unreachable included Taskfile remains user-owned and does not block migration"})
			}
			continue
		}
		digest, err := digestFile(path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		report.Sources[rel] = digest
		tf, err := taskfile.Parse(path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		if len(reachability.reachable[filepath.Clean(path)]) > 0 && (len(tf.Env) > 0 || len(tf.Dotenv) > 0 || len(tf.Shopt) > 0 || tf.Requires != nil) {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "taskfile-execution-policy", File: rel, Message: "reachable Taskfile environment, dotenv, shell options or requirements are not represented in YAML; keep using Taskfiles until migrated manually"})
		}
		role := taskfileRole(rel)
		allowed := legacyTaskNames[role]
		knownStock := knownStockTaskfiles[role][digest]
		if !knownStock && len(knownStockTaskfiles[role]) > 0 {
			data, err := os.ReadFile(path)
			if err != nil {
				return report, manifest.Document{}, err
			}
			// Git commonly checks out generated Taskfiles as CRLF on Windows.
			// YAML normalises these line breaks; keep the source receipt's raw
			// digest, but compare historical defaults using their shipped LF form.
			sum := sha256.Sum256(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
			knownStock = knownStockTaskfiles[role][hex.EncodeToString(sum[:])]
		}
		classification := migration.Taskfile{File: rel, Role: role}
		var canonical canonicalDiff
		if knownStock {
			// Shipped defaults changed throughout the v3 previews. Exact
			// fingerprints from generated projects are stock, even when their
			// task AST differs from the defaults compiled into this CLI.
			classification.Classification = "historical-default"
		} else {
			switch role {
			case "root":
				classification.Classification = "current-default"
			case "unknown":
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
		}
		taskNames := make([]string, 0, len(tf.Tasks))
		for name := range tf.Tasks {
			taskNames = append(taskNames, name)
		}
		sort.Strings(taskNames)
		for _, name := range taskNames {
			typescript, interfaces := legacyBindingOptions(tf.Tasks[name])
			typescriptBindings = typescriptBindings || typescript
			interfaceBindings = interfaceBindings || interfaces
			reachable := reachability.contains(path, name)
			if phase, script, ok := legacyLifecycleScript(root, filepath.Dir(path), name, tf.Tasks[name]); ok {
				severity := "info"
				code := "unrelated-task"
				message := "unreachable custom " + phase + " script remains user-owned"
				if reachable {
					hookPhase := manifest.HookPhase(phase)
					if previous, exists := migratedHooks[hookPhase]; exists && previous.Script != script {
						report.Complete = false
						severity, code = "warning", "duplicate-hook"
						message = "multiple reachable scripts map to " + phase + "; choose one manually"
					} else {
						migratedHooks[hookPhase] = manifest.Hook{Script: script}
						code = "translated-hook"
						message = "translated custom " + phase + " script to hook " + script
					}
				}
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: severity, Code: code, File: rel, Task: name, Message: message})
				if classification.Classification == "current-default" {
					classification.Classification = "customised"
				}
				classification.ModifiedTasks = appendUnique(classification.ModifiedTasks, name)
				continue
			}
			if !allowed[name] {
				severity := "info"
				code := "unrelated-task"
				message := "custom utility task is outside Wails build, package, sign, and dev paths"
				if reachable {
					report.Complete = false
					severity, code = "warning", "unsupported-task"
					message = "reachable custom task is not representable in config-only YAML; keep using Taskfiles"
				}
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: severity, Code: code, File: rel, Task: name, Message: message})
				if classification.Classification == "current-default" {
					classification.Classification = "customised"
				}
				classification.ModifiedTasks = appendUnique(classification.ModifiedTasks, name)
				continue
			}
			if reachable && (containsString(canonical.changed, name) || containsString(canonical.added, name)) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "generated task was modified and is not proven representable in config-only YAML; keep using Taskfiles"})
			}
			if reachable && role == "root" && !recognizedRootTask(name, tf.Tasks[name]) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "root dispatch task was modified and is not proven representable in config-only YAML; keep using Taskfiles"})
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
				if key == "BIN_DIR" && value.Static != "" && !strings.Contains(value.Static, "{{") {
					candidate := filepath.ToSlash(filepath.Clean(value.Static))
					if _, pathErr := manifest.ResolveProjectPath(root, "BIN_DIR", candidate, false); pathErr != nil {
						report.Complete = false
						report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "external-output", File: rel, Message: "BIN_DIR points outside the project and must be migrated manually"})
					} else {
						buildOutput = candidate
					}
				}
				if key == "VITE_PORT" && value.Static != "" && !strings.Contains(value.Static, "{{") {
					port, convErr := strconv.Atoi(strings.TrimSpace(value.Static))
					if convErr != nil || port < 1 || port > 65535 {
						devPort = 0
						report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "invalid-port", File: rel, Message: "VITE_PORT=" + value.Static + " is not a valid port and was not migrated"})
					} else if port != defaultVitePort {
						report.Complete = false
						report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "unsupported-dev-port", File: rel, Message: "VITE_PORT is not representable in YAML; use wails3 dev --port explicitly and remove the legacy override before activation"})
					} else {
						devPort = port
					}
				}
			}
		}
	}
	project := manifest.Project{Name: filepath.Base(root), ProductName: filepath.Base(root), Identifier: "com.wails." + normaliseName(filepath.Base(root)), Version: "0.1.0"}
	project.BinaryName = binaryName
	if binaryName != "" {
		project.Name = binaryName
	}
	doc := manifest.NewDocument(project)
	doc.Hooks = migratedHooks
	if buildOutput != "" {
		doc.Build.OutputDirectory = buildOutput
	}
	if devPort != 0 {
		doc.Dev.Port = devPort
	}
	doc.Frontend.Bindings.TypeScript = typescriptBindings
	doc.Frontend.Bindings.Interfaces = interfaceBindings
	if err := applyLegacyBuildConfiguration(root, &report, &doc); err != nil {
		return report, manifest.Document{}, err
	}
	if err := applyLegacyFrontendConfiguration(root, packageManager, &report, &doc); err != nil {
		return report, manifest.Document{}, err
	}
	if err := applyConventionalLegacyAssets(root, &report, &doc); err != nil {
		return report, manifest.Document{}, err
	}
	return report, doc, nil
}

func legacyBindingOptions(task *wakeast.Task) (typescript, interfaces bool) {
	if task == nil {
		return false, false
	}
	for _, command := range task.Cmds {
		if command == nil || !strings.Contains(command.Cmd, "generate bindings") {
			continue
		}
		for _, argument := range strings.Fields(command.Cmd) {
			name, value, hasValue := strings.Cut(strings.TrimSpace(argument), "=")
			if hasValue && value != "true" {
				continue
			}
			switch name {
			case "-ts", "--ts", "-typescript", "--typescript":
				typescript = true
			case "-interfaces", "--interfaces":
				interfaces = true
			}
		}
	}
	return typescript, interfaces
}

func normalizeLegacyAppName(value string) string {
	// Older generated root Taskfiles escaped spaces because APP_NAME was
	// interpolated into unquoted shell commands. The backslash is command
	// syntax, not part of the executable's file name.
	return strings.ReplaceAll(value, `\ `, " ")
}

func legacyLifecycleScript(root, taskfileDir, name string, task *wakeast.Task) (string, string, bool) {
	phase := strings.ReplaceAll(strings.ToLower(name), "-", "_")
	if !containsString([]string{"before_build", "after_build", "before_package", "after_package", "before_sign", "after_sign"}, phase) {
		return "", "", false
	}
	if task == nil || len(task.Cmds) != 1 || len(task.Deps) != 0 || task.Cmds[0].Cmd == "" || task.Cmds[0].Task != "" {
		return "", "", false
	}
	command := strings.TrimSpace(task.Cmds[0].Cmd)
	if strings.ContainsAny(command, " \t\r\n;&|`$<>") || strings.Contains(command, "{{") {
		return "", "", false
	}
	base := taskfileDir
	if task.Dir != "" {
		base = filepath.Join(base, filepath.FromSlash(task.Dir))
	}
	script := filepath.Clean(filepath.Join(base, filepath.FromSlash(strings.TrimPrefix(command, "./"))))
	relative, err := filepath.Rel(root, script)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", "", false
	}
	resolved, err := manifest.ResolveProjectPath(root, "hook script", filepath.ToSlash(relative), true)
	if err != nil {
		return "", "", false
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() || runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		return "", "", false
	}
	if runtime.GOOS != "windows" {
		file, openErr := os.Open(resolved)
		if openErr != nil {
			return "", "", false
		}
		var prefix [2]byte
		_, readErr := io.ReadFull(file, prefix[:])
		closeErr := file.Close()
		if readErr != nil || closeErr != nil || string(prefix[:]) != "#!" {
			return "", "", false
		}
	} else if !validWindowsHookMigrationScript(relative) {
		return "", "", false
	}
	return phase, filepath.ToSlash(relative), true
}

func validWindowsHookMigrationScript(value string) bool {
	switch strings.ToLower(filepath.Ext(value)) {
	case ".cmd", ".bat", ".ps1":
		return true
	default:
		return false
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
	switch strings.TrimSpace(command) {
	case "go mod tidy", "wails3 build DEV=true", "wails3 task common:install:frontend:deps", "wails3 task common:dev:frontend", "wails3 task build", "wails3 task run":
		return true
	}
	return false
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
		tf, err := taskfile.Parse(abs)
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
			if strings.Contains(path, "{{") {
				continue
			}
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
				continue
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

var stockVariantCache = struct {
	sync.Mutex
	values map[string][]*wakeast.Taskfile
}{values: make(map[string][]*wakeast.Taskfile, 6)}

func stockTaskVariants(role string) []map[string]*wakeast.Task {
	var result []map[string]*wakeast.Task
	for _, variant := range stockTaskfileVariants(role) {
		result = append(result, variant.Tasks)
	}
	return result
}

func stockTaskfileVariants(role string) []*wakeast.Taskfile {
	stockVariantCache.Lock()
	variants, ok := stockVariantCache.values[role]
	stockVariantCache.Unlock()
	if ok {
		return variants
	}
	variants = loadStockTaskVariants(role)
	stockVariantCache.Lock()
	if existing, loaded := stockVariantCache.values[role]; loaded {
		variants = existing
	} else {
		stockVariantCache.values[role] = variants
	}
	stockVariantCache.Unlock()
	return variants
}

func loadStockTaskVariants(role string) []*wakeast.Taskfile {
	if role == "root" || role == "unknown" {
		return nil
	}
	if role == "common" {
		data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
		if err != nil {
			return nil
		}
		var result []*wakeast.Taskfile
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
				if tasks := parseStockTaskfile(rendered.Bytes()); tasks != nil {
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
	if tasks := parseStockTaskfile(data); tasks != nil {
		return []*wakeast.Taskfile{tasks}
	}
	return nil
}

func parseStockTasks(data []byte) map[string]*wakeast.Task {
	if tf := parseStockTaskfile(data); tf != nil {
		return tf.Tasks
	}
	return nil
}

func parseStockTaskfile(data []byte) *wakeast.Taskfile {
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
	tf, err := taskfile.Parse(name)
	if err != nil {
		return nil
	}
	return tf
}

func closestCanonical(actual map[string]*wakeast.Task, variants []map[string]*wakeast.Task) canonicalDiff {
	best := canonicalDiff{added: sortedTaskNames(actual)}
	if len(variants) == 0 {
		return best
	}
	bestScore := int(^uint(0) >> 1)
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
	switch rel {
	case "Taskfile.yml", "Taskfile.yaml":
		return "root"
	case "build/Taskfile.yml", "build/Taskfile.yaml":
		return "common"
	}
	switch {
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

// Exact SHA256 fingerprints of generated platform Taskfiles shipped in v3
// examples before the manifest cutover. These distinguish historical defaults
// from user edits without weakening comparison against the current templates.
var knownStockTaskfiles = map[string]map[string]bool{
	"common": set(
		"eab2e71045aa6a1bd6cf54fc1a9222184dc9fb6af00ba6e137d3797cdefe26cd",
		"d1b5e40d4c6dd9a2eff9fb4b16ffa829522c683fac49b473414b5f1690e75b6d",
		"ed59e465c150f605abf361a0d85f8f7b85183b21a659a1d7a73e4b40740e9198",
		"8340567cb86c254059cd2de5822c93d39d952648f46c10a337e993bf1fb63336",
		"33339f2554bc2422fc64a2c92bed7d35f1815ea2630733b83289bc340a7e5ddd",
		"3b1c9ac25d1491ce00a66dbd249604e4e01c6f420d53878125f20faa4b56cff3",
		"374ead652337eef98dc5fbf8bea6ed9a451e65f0bfcca1173736568c4999bee4",
	),
	"windows": set(
		"7142682189eaa323166330c87fe4741503db68f8d7d4a50a7e69d043f936dd2f",
		"bbc1c366f0d2dd14ad9bf4722e7b2e5977c49461b11f17d7789b8dd10ed3ba45",
		"6dbce0ac85b317d49b8c53acc67844d40e12eb625a63627d10c598bfcf25ca39",
	),
	"darwin": set(
		"6b75611710a8ab7aeb33e61b925c15fcb38aa7b22c1482f45e7b5e492909cec4",
		"a5e086690d6dbeb1506718b26516b8025631f48107f471e470a81904c7d9ab64",
		"1a8995dc7e3a0d710218203e16a914be5c775867289b58acd2dd367a97b7a14a",
	),
	"linux": set(
		"208c899a8871442bc67b41d95fcd70a9a671548a44e8cd995e747178d0d06357",
		"d7cc6d72e5cd370b1f6d2a65bb5602d0063cf1abbb0cd7482ae68090230c04f2",
	),
	"ios": set(
		"3fc2d20d261308e347bc59b4c85c2753d51e7374ea5c593d5f6a1d48238608ff",
		"f40e4938c8a4b85baf0a015608f8f9c92b57b36c514459a9c01abd93a43bc003",
	),
	"android": set(
		"a0c72b5349146dc8b59b8a4c424b1ee68080289c0d3ce2f5ad35fb3e96ae5610",
		"48b6e8f3ac3e87c1898f9aa018919d7b95a8cf4ee8011105c33b37410525ff7f",
	),
}

func recognizedRootTask(name string, task *wakeast.Task) bool {
	if task == nil || len(task.Cmds) != 1 || len(task.Deps) != 0 || len(task.Env) > 0 || len(task.Vars) > 0 || task.Dir != "" || len(task.Precondition) > 0 || len(task.Status) > 0 || len(task.Defer) > 0 {
		return false
	}
	if !rootTaskDefinitionSupported(task.Definition) {
		return false
	}
	command := task.Cmds[0]
	if command == nil || len(command.Vars) > 0 || command.IgnoreError || command.For != nil {
		return false
	}
	if name == "dev" {
		return command.Task == "" && (strings.TrimSpace(command.Cmd) == "wails3 dev -config ./build/config.yml" || strings.TrimSpace(command.Cmd) == "wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}")
	}
	refs := map[string]string{
		"setup:docker": "common:setup:docker", "build:server": "common:build:server",
		"run:server": "common:run:server", "build:docker": "common:build:docker", "run:docker": "common:run:docker",
	}
	if expected := refs[name]; expected != "" {
		return command.Task == expected && command.Cmd == ""
	}
	if name == "build" || name == "package" || name == "run" {
		return command.Cmd == "" && (command.Task == "{{OS}}:"+name || command.Task == "{{.GOOS}}:"+name)
	}
	return false
}

// Root dispatchers have a small contract; unknown execution fields must not
// disappear merely because the declaration reader does not interpret them.
func rootTaskDefinitionSupported(definition any) bool {
	if definition == nil {
		return true
	}
	fields, ok := definition.(map[string]any)
	if !ok {
		return false
	}
	for key := range fields {
		switch key {
		case "summary", "desc", "cmds":
		default:
			return false
		}
	}
	commands, ok := fields["cmds"].([]any)
	if !ok || len(commands) != 1 {
		return false
	}
	switch command := commands[0].(type) {
	case string:
		return true
	case map[string]any:
		if len(command) == 1 {
			_, task := command["task"]
			_, cmd := command["cmd"]
			return task || cmd
		}
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

func digestFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
func first(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
