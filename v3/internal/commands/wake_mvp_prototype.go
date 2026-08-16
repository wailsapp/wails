package commands

// This file is a throwaway prototype. It exists to measure and review the
// manifest-driven Wake fast path before the production architecture is built.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/wake/mvpprototype"
)

type mvpNodeKind string

const (
	mvpInstall  mvpNodeKind = "InstallFrontendDependencies"
	mvpBindings mvpNodeKind = "GenerateBindings"
	mvpFrontend mvpNodeKind = "BuildFrontend"
	mvpCompile  mvpNodeKind = "CompileApplication"
)

type mvpNode struct {
	Kind   mvpNodeKind
	Label  string
	Output string
	Spec   any
}

type mvpPlan struct {
	Target string
	Nodes  []mvpNode
}

type mvpBuild struct {
	root        string
	manifest    mvpprototype.Manifest
	plan        mvpPlan
	cache       *mvpprototype.Cache
	reporter    report.Reporter
	buildFlags  *flags.Build
	npm         mvpTool
	goTool      mvpTool
	goInput     string
	moduleInput string
	frontInput  string
	installIn   string
}

type mvpTool struct {
	Path   string
	Runner string
	Script string
	ID     string
}

func hasWakeMVPManifest(root string) bool {
	_, err := os.Stat(filepath.Join(root, "wails.toml"))
	return err == nil
}

func runWakeMVP(buildFlags *flags.Build) error {
	overallStarted := time.Now()
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	manifest, err := mvpprototype.LoadManifest(root)
	if err != nil {
		return err
	}
	cache, err := mvpprototype.OpenCache(root)
	if err != nil {
		return err
	}
	npm, err := resolveNPM()
	if err != nil {
		return err
	}
	goTool, err := resolveTool("go")
	if err != nil {
		return err
	}

	tags := []string{"production"}
	if buildFlags.Tags != "" {
		tags = append(tags, strings.Split(buildFlags.Tags, ",")...)
	}
	plan := mvpPlan{
		Target: runtime.GOOS + "/" + runtime.GOARCH,
		Nodes: []mvpNode{
			{Kind: mvpInstall, Label: "Install frontend dependencies", Spec: map[string]any{"manager": "npm"}},
			{Kind: mvpBindings, Label: "Generate bindings", Output: "frontend/bindings", Spec: map[string]any{"tags": tags, "typescript": true, "interfaces": true}},
			{Kind: mvpFrontend, Label: "Build frontend", Output: "frontend/dist", Spec: map[string]any{"script": "build", "mode": "production"}},
			{Kind: mvpCompile, Label: "Compile application", Output: filepath.ToSlash(filepath.Join("bin", executableName(manifest.Project.BinaryName))), Spec: map[string]any{"target": runtime.GOOS + "/" + runtime.GOARCH, "tags": tags}},
		},
	}

	rep := pulse.New(os.Stdout, report.Normal)
	b := &mvpBuild{root: root, manifest: manifest, plan: plan, cache: cache, reporter: rep, buildFlags: buildFlags, npm: npm, goTool: goTool}
	if err := b.discoverInputs(); err != nil {
		return err
	}

	term.Header("Build")
	fmt.Fprintln(os.Stderr, "  wake manifest MVP · local prototype · no Taskfiles")
	report.SetActive(rep)
	defer report.SetActive(nil)
	rep.BuildStart("build", b.plan.Target, len(b.plan.Nodes))
	err = b.execute(context.Background())
	if err != nil {
		rep.BuildEnd(time.Since(overallStarted), false)
		return err
	}
	output := filepath.Join(root, filepath.FromSlash(plan.Nodes[len(plan.Nodes)-1].Output))
	rep.Artifact(report.Artifact{Path: output, Kind: "binary"})
	rep.BuildEnd(time.Since(overallStarted), true)
	return cache.Save()
}

func (b *mvpBuild) discoverInputs() error {
	var err error
	b.installIn, err = b.cache.SnapshotFiles("frontend-install",
		"frontend/package.json", "frontend/package-lock.json", "frontend/npm-shrinkwrap.json", "frontend/.npmrc")
	if err != nil {
		return err
	}
	moduleRoot := nearestModuleRoot(b.root)
	b.goInput, err = b.cache.Snapshot(mvpprototype.SnapshotOptions{
		Label: "go-module-sources", Root: moduleRoot,
		IncludeNames:      []string{"go.mod", "go.sum", "go.work", "go.work.sum"},
		IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"},
		ExcludeDirs:       []string{".git", ".wails", "bin", "build", "dist", "frontend", "node_modules"},
	})
	if err != nil {
		return err
	}
	moduleFiles := nearestModuleFiles(b.root)
	b.moduleInput, err = b.cache.SnapshotFiles("go-module", moduleFiles...)
	if err != nil {
		return err
	}
	b.frontInput, err = b.snapshotFrontend()
	return err
}

func (b *mvpBuild) execute(ctx context.Context) error {
	installKey, err := b.runInstall(ctx, b.plan.Nodes[0])
	if err != nil {
		return err
	}
	// Installation may create a lockfile. Refresh the frontend Snapshot before
	// assigning its Action Key so a successful cold build is already stable on
	// the immediately following no-op build.
	b.frontInput, err = b.snapshotFrontend()
	if err != nil {
		return err
	}
	bindingsArtifact, err := b.runArtifactNode(ctx, b.plan.Nodes[1], []string{b.goInput, b.moduleInput}, nil)
	if err != nil {
		return err
	}
	frontendArtifact, err := b.runArtifactNode(ctx, b.plan.Nodes[2], []string{b.frontInput}, []string{installKey, bindingsArtifact})
	if err != nil {
		return err
	}
	_, err = b.runArtifactNode(ctx, b.plan.Nodes[3], []string{b.goInput, b.moduleInput}, []string{frontendArtifact})
	return err
}

func (b *mvpBuild) snapshotFrontend() (string, error) {
	return b.cache.Snapshot(mvpprototype.SnapshotOptions{
		Label: "frontend-source", Root: filepath.Join(b.root, "frontend"), IncludeAll: true,
		ExcludeDirs: []string{".git", ".wails", "bindings", "dist", "node_modules"},
	})
}

func (b *mvpBuild) runInstall(ctx context.Context, node mvpNode) (string, error) {
	key, err := mvpprototype.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": b.npm.ID}, []string{b.installIn}, nil)
	if err != nil {
		return "", err
	}
	id := b.reporter.StepStart(string(node.Kind), node.Label)
	started := time.Now()
	if b.cache.HasReceipt(key, "frontend/node_modules") {
		b.reporter.StepInfo(id, "receipt matched; dependency tree not scanned")
		b.reporter.StepEnd(id, report.StatusCached, time.Since(started))
		return key, nil
	}
	b.reporter.StepCommand(id, "npm install --no-audit --no-fund")
	output, runErr := b.npm.run(ctx, filepath.Join(b.root, "frontend"), "install", "--no-audit", "--no-fund")
	if runErr != nil {
		b.reporter.StepFailed(id, report.Failure{Task: string(node.Kind), Command: "npm install", Output: output, Err: runErr, ExitCode: exitCode(runErr)})
		return "", runErr
	}
	// npm may create a lockfile. Record the Receipt under the resulting input
	// identity so the very next build is already a hit.
	postInput, err := b.cache.SnapshotFiles("frontend-install",
		"frontend/package.json", "frontend/package-lock.json", "frontend/npm-shrinkwrap.json", "frontend/.npmrc")
	if err != nil {
		return "", err
	}
	key, err = mvpprototype.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": b.npm.ID}, []string{postInput}, nil)
	if err != nil {
		return "", err
	}
	if err := b.cache.RecordReceipt(key); err != nil {
		return "", err
	}
	b.reporter.StepEnd(id, report.StatusOK, time.Since(started))
	return key, nil
}

func (b *mvpBuild) runArtifactNode(ctx context.Context, node mvpNode, inputs, deps []string) (string, error) {
	tool := b.goTool.ID
	if node.Kind == mvpFrontend {
		tool = b.npm.ID
	}
	key, err := mvpprototype.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": tool}, inputs, deps)
	if err != nil {
		return "", err
	}
	id := b.reporter.StepStart(string(node.Kind), node.Label)
	started := time.Now()
	b.cache.ResetStats()
	status, artifact, err := b.cache.Lookup(key, node.Output)
	if err != nil {
		return "", err
	}
	if status == mvpprototype.LookupHit || status == mvpprototype.LookupRestored {
		stats := b.cache.Stats()
		detail := fmt.Sprintf("%s · %d file digests reused · %s read", status, stats.DigestsReused, formatBytes(stats.BytesRead))
		b.reporter.StepInfo(id, detail)
		b.reporter.StepEnd(id, report.StatusCached, time.Since(started))
		return artifact, nil
	}
	if status == mvpprototype.LookupDirty {
		b.reporter.StepInfo(id, "generated output was modified; rebuilding")
	}
	command, output, runErr := b.executeNode(ctx, node)
	if command != "" {
		b.reporter.StepCommand(id, command)
	}
	if runErr != nil {
		b.reporter.StepFailed(id, report.Failure{Task: string(node.Kind), Command: command, Output: output, Err: runErr, ExitCode: exitCode(runErr)})
		return "", runErr
	}
	artifact, err = b.cache.RecordAction(key, node.Output)
	if err != nil {
		b.reporter.StepFailed(id, report.Failure{Task: string(node.Kind), Err: err, ExitCode: -1})
		return "", err
	}
	b.reporter.StepEnd(id, report.StatusOK, time.Since(started))
	return artifact, nil
}

func (b *mvpBuild) executeNode(ctx context.Context, node mvpNode) (command, output string, err error) {
	switch node.Kind {
	case mvpBindings:
		tags := []string{"production"}
		if b.buildFlags.Tags != "" {
			tags = append(tags, strings.Split(b.buildFlags.Tags, ",")...)
		}
		command = "generate bindings (in process)"
		previousFooter := DisableFooter
		err = GenerateBindings(&flags.GenerateBindingsOptions{
			BuildFlagsString: "-tags " + strings.Join(tags, ","),
			OutputDirectory:  filepath.Join(b.root, filepath.FromSlash(node.Output)),
			ModelsFilename:   "models", IndexFilename: "index", TimeType: "string",
			TS: true, UseInterfaces: true, Clean: true, Silent: true,
		}, nil)
		DisableFooter = previousFooter
		return command, "", err
	case mvpFrontend:
		command = "npm run build --silent"
		output, err = b.npm.run(ctx, filepath.Join(b.root, "frontend"), "run", "build", "--silent")
		return command, output, err
	case mvpCompile:
		args := []string{"build", "-tags", "production", "-trimpath", "-buildvcs=false", "-ldflags", "-w -s", "-o", filepath.Join(b.root, filepath.FromSlash(node.Output)), "."}
		if b.buildFlags.Tags != "" {
			args[2] += "," + b.buildFlags.Tags
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(b.root, filepath.FromSlash(node.Output))), 0o755); err != nil {
			return "", "", err
		}
		command = "go " + strings.Join(args, " ")
		output, err = runCommand(ctx, b.root, b.goTool.Path, args...)
		return command, output, err
	default:
		return "", "", fmt.Errorf("wake mvp: unsupported node %s", node.Kind)
	}
}

func resolveNPM() (mvpTool, error) {
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return mvpTool{}, fmt.Errorf("wake mvp: npm not found: %w", err)
	}
	script, err := filepath.EvalSymlinks(npmPath)
	if err != nil {
		script = npmPath
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return mvpTool{}, fmt.Errorf("wake mvp: node not found: %w", err)
	}
	id, err := toolIdentity(nodePath, script)
	if err != nil {
		return mvpTool{}, err
	}
	return mvpTool{Path: npmPath, Runner: nodePath, Script: script, ID: id}, nil
}

func resolveTool(name string) (mvpTool, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return mvpTool{}, fmt.Errorf("wake mvp: %s not found: %w", name, err)
	}
	id, err := toolIdentity(path)
	if err != nil {
		return mvpTool{}, err
	}
	return mvpTool{Path: path, ID: id}, nil
}

func toolIdentity(paths ...string) (string, error) {
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%s:%d:%d:%o", path, info.Size(), info.ModTime().UnixNano(), info.Mode()))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|"), nil
}

func (t mvpTool) run(ctx context.Context, dir string, args ...string) (string, error) {
	if t.Runner != "" {
		return runCommand(ctx, dir, t.Runner, append([]string{t.Script}, args...)...)
	}
	return runCommand(ctx, dir, t.Path, args...)
}

func runCommand(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PRODUCTION=true")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return output.String(), err
}

func executableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func exitCode(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

func formatBytes(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d B", value)
	}
	return fmt.Sprintf("%.1f KiB", float64(value)/1024)
}

func nearestModuleFiles(root string) []string {
	dir := nearestModuleRoot(root)
	var files []string
	for _, name := range []string{"go.mod", "go.sum", "go.work", "go.work.sum"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}
	return files
}

func nearestModuleRoot(root string) string {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return root
		}
		dir = parent
	}
}
