package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func generatedMigrationTasks(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	if err := GenerateBuildAssets(&BuildAssetsOptions{Name: "app", ProductName: "App", ProductVersion: "1.0.0", Dir: filepath.Join(root, "build"), Silent: true}); err != nil {
		t.Fatal(err)
	}
	rootTask := filepath.Join(root, "Taskfile.yml")
	data := `version: '3'
includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml
tasks:
  build:
    cmds:
      - task: windows:build
      - task: darwin:build
      - task: linux:build
      - task: darwin:build:universal
      - task: common:dev:frontend
`
	if err := os.WriteFile(rootTask, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, rootTask
}

func migrationTaskReachability(t *testing.T, rootTask string) migrationReachability {
	t.Helper()
	files, err := discoverTaskfiles(rootTask)
	if err != nil {
		t.Fatal(err)
	}
	reachable, err := analyseTaskReachability(rootTask, files)
	if err != nil {
		t.Fatal(err)
	}
	return reachable
}

func TestMigrationStockDynamicSelectorsVisitEveryBranch(t *testing.T) {
	root, rootTask := generatedMigrationTasks(t)
	reachable := migrationTaskReachability(t, rootTask)
	if len(reachable.issues) != 0 {
		t.Fatalf("stock tasks must be statically recognised: %+v", reachable.issues)
	}
	for file, names := range map[string][]string{
		"build/Taskfile.yml":         {"install:frontend:deps:npm", "install:frontend:deps:pnpm", "install:frontend:deps:yarn", "install:frontend:deps:bun", "frontend:run:npm", "frontend:run:pnpm", "frontend:run:yarn", "frontend:run:bun", "frontend:dev:npm", "frontend:dev:pnpm", "frontend:dev:yarn", "frontend:dev:bun"},
		"build/windows/Taskfile.yml": {"build:native", "build:docker", "create:msix:package", "create:nsis:installer"},
		"build/darwin/Taskfile.yml":  {"build:native", "build:docker", "build:universal:lipo:native", "build:universal:lipo:go", "codesign:adhoc", "codesign:skip"},
		"build/linux/Taskfile.yml":   {"build:native", "build:docker"},
	} {
		for _, name := range names {
			if !reachable.contains(filepath.Join(root, filepath.FromSlash(file)), name) {
				t.Errorf("stock selector hid branch %s [%s]", file, name)
			}
		}
	}
}

func TestMigrationChangedDynamicSelectorsRemainBlocked(t *testing.T) {
	for _, test := range []struct{ name, before, after string }{
		{"expression", `install:frontend:deps:{{.PACKAGE_MANAGER}}`, `install:frontend:deps:{{.CUSTOM_RUNNER}}`},
		{"extra command", `      - task: install:frontend:deps:{{.PACKAGE_MANAGER}}`, "      - cmd: echo extra behavior\n      - task: install:frontend:deps:{{.PACKAGE_MANAGER}}"},
		{"execution policy", "    run: once\n    cmds:\n      - task: install:frontend:deps:", "    run: always\n    cmds:\n      - task: install:frontend:deps:"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, rootTask := generatedMigrationTasks(t)
			path := filepath.Join(root, "build", "Taskfile.yml")
			replaceMigrationTaskText(t, path, test.before, test.after)
			reachable := migrationTaskReachability(t, rootTask)
			for _, diagnostic := range reachable.issues {
				if diagnostic.File == path && diagnostic.Task == "install:frontend:deps" && diagnostic.Code == "dynamic-task-reference" {
					return
				}
			}
			t.Fatalf("changed selector was accepted: %+v", reachable.issues)
		})
	}
}

func TestMigrationModifiedDynamicBranchRemainsReachable(t *testing.T) {
	for _, test := range []struct{ name, file, before, after string }{
		{"frontend:run:bun", "build/Taskfile.yml", "      - bun run {{.SCRIPT}}", "      - echo custom-bun {{.SCRIPT}}"},
		{"build:docker", "build/linux/Taskfile.yml", "  build:docker:\n", "  build:docker:\n    prompt: Confirm custom build\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, rootTask := generatedMigrationTasks(t)
			path := filepath.Join(root, filepath.FromSlash(test.file))
			replaceMigrationTaskText(t, path, test.before, test.after)
			reachable := migrationTaskReachability(t, rootTask)
			if !reachable.contains(path, test.name) {
				t.Fatal("modified dynamic branch was not visited")
			}
			report, _, err := analyseMigration(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, diagnostic := range report.Diagnostics {
				if diagnostic.File == test.file && diagnostic.Task == test.name && diagnostic.Code == "modified-task" {
					return
				}
			}
			t.Fatalf("modified dynamic branch was not blocked: %+v", report.Diagnostics)
		})
	}
}

func TestMigrationMissingDynamicBranchIsBlocked(t *testing.T) {
	root, rootTask := generatedMigrationTasks(t)
	path := filepath.Join(root, "build", "windows", "Taskfile.yml")
	replaceMigrationTaskText(t, path, "  create:msix:package:\n", "  renamed:msix:package:\n")
	reachable := migrationTaskReachability(t, rootTask)
	for _, diagnostic := range reachable.issues {
		if diagnostic.File == path && diagnostic.Task == "package" && diagnostic.Code == "unresolved-task-reference" && strings.Contains(diagnostic.Message, "create:msix:package") {
			return
		}
	}
	t.Fatalf("missing stock branch was accepted: %+v", reachable.issues)
}

func replaceMigrationTaskText(t *testing.T, path, before, after string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.ReplaceAll(string(data), "\r\n", "\n"))
	if !strings.Contains(string(data), before) {
		t.Fatalf("fixture %s does not contain %q", path, before)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), before, after, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMigrationStockSelectorsRejectAlteredContext(t *testing.T) {
	for _, test := range []struct{ name, file, before, after string }{
		{"root vars", "Taskfile.yml", "tasks:\n", "vars:\n  EXTRA_TAGS: custom\ntasks:\n"},
		{"platform vars", "build/linux/Taskfile.yml", "  CROSS_IMAGE: wails-cross", "  CROSS_IMAGE: custom-image"},
		{"common vars", "build/Taskfile.yml", "tasks:\n", "vars:\n  PACKAGE_MANAGER: custom-runner\ntasks:\n"},
		{"include vars", "Taskfile.yml", "  linux: ./build/linux/Taskfile.yml", "  linux:\n    taskfile: ./build/linux/Taskfile.yml\n    vars:\n      CGO_ENABLED: 1"},
		{"include directory", "Taskfile.yml", "  linux: ./build/linux/Taskfile.yml", "  linux:\n    taskfile: ./build/linux/Taskfile.yml\n    dir: other"},
		{"include unknown field", "Taskfile.yml", "  linux: ./build/linux/Taskfile.yml", "  linux:\n    taskfile: ./build/linux/Taskfile.yml\n    flatten: true"},
		{"task unknown field", "build/linux/Taskfile.yml", "  build:\n", "  build:\n    ignore_error: true\n"},
		{"file unknown field", "build/linux/Taskfile.yml", "tasks:\n", "method: none\ntasks:\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, rootTask := generatedMigrationTasks(t)
			replaceMigrationTaskText(t, filepath.Join(root, filepath.FromSlash(test.file)), test.before, test.after)
			reachable := migrationTaskReachability(t, rootTask)
			for _, diagnostic := range reachable.issues {
				if diagnostic.Code == "dynamic-task-reference" {
					return
				}
			}
			t.Fatalf("altered selector context was accepted: %+v", reachable.issues)
		})
	}
}

func TestMigrationStockDynamicBranchCannotBeMaskedByRootTask(t *testing.T) {
	root, rootTask := generatedMigrationTasks(t)
	f, err := os.OpenFile(rootTask, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("  build:docker:\n    cmds:\n      - task: common:build:docker\n")
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		t.Fatalf("append root task: %v, close: %v", err, closeErr)
	}
	path := filepath.Join(root, "build", "linux", "Taskfile.yml")
	replaceMigrationTaskText(t, path, "  build:docker:\n", "  build:docker:\n    ignore_error: true\n")
	reachable := migrationTaskReachability(t, rootTask)
	if !reachable.contains(path, "build:docker") {
		t.Fatal("root build:docker masked the platform's selected branch")
	}
	report, _, err := analyseMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.File == "build/linux/Taskfile.yml" && diagnostic.Task == "build:docker" && diagnostic.Code == "modified-task" {
			return
		}
	}
	t.Fatalf("platform mutation was not blocked: %+v", report.Diagnostics)
}

func TestMigrationRootDispatchRejectsUnrepresentedBehavior(t *testing.T) {
	for _, test := range []struct{ name, task, fields, command string }{
		{"unknown selector", "build", "", "      - task: '{{.CUSTOM}}:build'\n"},
		{"root execution field", "build", "    ignore_error: true\n", "      - task: '{{OS}}:build'\n"},
		{"command execution field", "build", "", "      - task: '{{OS}}:build'\n        defer: true\n"},
		{"dev flags", "dev", "", "      - wails3 dev -config ./build/config.yml -tags custom\n"},
		{"dev shell suffix", "dev", "", "      - wails3 dev -config ./build/config.yml && echo custom\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			text := "version: '3'\ntasks:\n  " + test.task + ":\n" + test.fields + "    cmds:\n" + test.command
			if err := os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
			report, _, err := analyseMigration(root)
			if err != nil {
				t.Fatal(err)
			}
			if report.Complete {
				t.Fatal("unsupported root dispatcher was accepted")
			}
			for _, diagnostic := range report.Diagnostics {
				if diagnostic.Task == test.task && diagnostic.Code == "modified-task" {
					return
				}
			}
			t.Fatalf("missing root dispatcher diagnostic: %+v", report.Diagnostics)
		})
	}
}
