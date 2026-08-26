package commands

import (
	"bytes"
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
	wakeast "github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/migration"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/zeebo/blake3"
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
		return "", "", fmt.Errorf("migration output %q must be an inactive project-owned HCL file", configured)
	}
	path, err := manifest.ResolveProjectPath(root, "migration output", clean, false)
	if err != nil {
		return "", "", err
	}
	if !strings.EqualFold(filepath.Ext(path), ".hcl") {
		return "", "", fmt.Errorf("migration output %q must use the .hcl extension", configured)
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
	var load func(string, string) error
	load = func(path, prefix string) error {
		path = filepath.Clean(path)
		if ancestors[path] {
			includeIssues = append(includeIssues, MigrationDiagnostic{Severity: "warning", Code: "cyclic-include", File: path, Message: "Taskfile include cycle cannot be proven equivalent"})
			return nil
		}
		visitKey := path + "\x00" + prefix
		if visited[visitKey] {
			return nil
		}
		visited[visitKey] = true
		ancestors[path] = true
		defer delete(ancestors, path)
		loadedPaths[path] = true
		tf, err := parse.Parse(path)
		if err != nil {
			return err
		}
		for name, task := range tf.Tasks {
			key := prefix + name
			nodes[key] = append(nodes[key], migrationTaskNode{key: key, prefix: prefix, path: path, name: name, task: task})
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
			if err := load(included, prefix+name+":"); err != nil {
				return err
			}
		}
		return nil
	}
	if err := load(rootTask, ""); err != nil {
		return migrationReachability{}, err
	}
	// Conventional local/override files are merged at the root namespace even
	// though they are not reached through an explicit include.
	for _, path := range files {
		path = filepath.Clean(path)
		if !loadedPaths[path] {
			if err := load(path, ""); err != nil {
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
				if strings.Contains(reference, "{{") {
					if node.prefix == "" && recognizedRootTask(node.name, node.task) {
						suffix := ":" + node.name
						for candidate := range nodes {
							if strings.HasSuffix(candidate, suffix) {
								queue = append(queue, candidate)
							}
						}
						continue
					}
					addIssue(node, "dynamic-task-reference", "reachable task selects another task dynamically and cannot be proven equivalent")
					continue
				}
				candidate := reference
				if len(nodes[candidate]) == 0 && node.prefix != "" {
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
		tf, err := parse.Parse(path)
		if err != nil {
			return report, manifest.Document{}, err
		}
		role := taskfileRole(rel)
		allowed := legacyTaskNames[role]
		knownStock := knownStockTaskfiles[role][digest]
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
					message = "reachable custom task is not representable in config-only HCL; keep using Taskfiles"
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
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "generated task was modified and is not proven representable in config-only HCL; keep using Taskfiles"})
			}
			if reachable && role == "root" && !recognizedRootTask(name, tf.Tasks[name]) {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: rel, Task: name, Message: "root dispatch task was modified and is not proven representable in config-only HCL; keep using Taskfiles"})
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
	values map[string][]map[string]*wakeast.Task
}{values: make(map[string][]map[string]*wakeast.Task, 6)}

func stockTaskVariants(role string) []map[string]*wakeast.Task {
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

func loadStockTaskVariants(role string) []map[string]*wakeast.Task {
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
