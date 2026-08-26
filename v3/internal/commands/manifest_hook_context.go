package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

const hookContextFileEnvironment = "WAILS_HOOK_CONTEXT_FILE"

var legacyHookContextEnvironment = []string{
	"WAILS_HOOK_CONTEXT",
	hookContextFileEnvironment,
	"WAILS_PROJECT_DIR",
	"WAILS_TARGET_OS",
	"WAILS_TARGET_ARCH",
	"WAILS_PROFILE",
	"WAILS_OUTPUT",
	"WAILS_PIPELINE_VERSION",
}

type hookExecutionContext struct {
	Version          int                `json:"version"`
	Phase            manifest.HookPhase `json:"phase"`
	Command          string             `json:"command"`
	Scope            pipeline.Scope     `json:"scope"`
	ProjectDirectory string             `json:"project_dir"`
	WorkingDirectory string             `json:"working_directory"`
	Profile          string             `json:"profile"`
	Target           *hookContextTarget `json:"target"`
	Output           string             `json:"output"`
	DeclaredOutputs  []string           `json:"declared_outputs"`
}

type hookContextTarget struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

func (h *manifestHandler) writeHookContext(spec pipeline.HookSpec, directory, output string) (string, string, error) {
	version := spec.ContextVersion
	if version == 0 {
		version = 1
	}
	command := spec.Command
	if command == "" {
		command = "build"
	}
	scope := spec.Scope
	if scope == "" {
		scope = pipeline.ProjectScope
		if spec.TargetOS != "" || spec.TargetArch != "" {
			scope = pipeline.TargetScope
		}
	}
	var target *hookContextTarget
	if spec.TargetOS != "" || spec.TargetArch != "" {
		target = &hookContextTarget{OS: spec.TargetOS, Arch: spec.TargetArch}
	}
	declaredOutputs := make([]string, 0, len(spec.DeclaredOutputs))
	for _, declared := range spec.DeclaredOutputs {
		resolved, err := manifest.ResolveProjectPath(h.root, "hook declared output", declared, false)
		if err != nil {
			return "", "", err
		}
		declaredOutputs = append(declaredOutputs, resolved)
	}
	content, err := json.MarshalIndent(hookExecutionContext{
		Version:          version,
		Phase:            spec.Phase,
		Command:          command,
		Scope:            scope,
		ProjectDirectory: h.root,
		WorkingDirectory: directory,
		Profile:          spec.Profile,
		Target:           target,
		Output:           output,
		DeclaredOutputs:  declaredOutputs,
	}, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("encode hook %s context: %w", spec.Phase, err)
	}
	content = append(content, '\n')
	contextRoot, err := hookScopeOutputPath(h.root, filepath.ToSlash(filepath.Join(".wails", "hooks", "context")))
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(contextRoot, 0o700); err != nil {
		return "", "", fmt.Errorf("create hook context root: %w", err)
	}
	if err := protectHookContextPath(contextRoot, true); err != nil {
		return "", "", fmt.Errorf("protect hook context root: %w", err)
	}
	contextDirectory, err := os.MkdirTemp(contextRoot, string(spec.Phase)+"-")
	if err != nil {
		return "", "", fmt.Errorf("create hook %s context directory: %w", spec.Phase, err)
	}
	removeDirectory := true
	defer func() {
		if removeDirectory {
			_ = os.RemoveAll(contextDirectory)
		}
	}()
	if err := protectHookContextPath(contextDirectory, true); err != nil {
		return "", "", fmt.Errorf("protect hook %s context directory: %w", spec.Phase, err)
	}
	temporary, err := os.CreateTemp(contextDirectory, ".context-")
	if err != nil {
		return "", "", fmt.Errorf("create hook %s context: %w", spec.Phase, err)
	}
	temporaryName := temporary.Name()
	removeTemporary := true
	defer func() {
		_ = temporary.Close()
		if removeTemporary {
			_ = os.Remove(temporaryName)
		}
	}()
	if err := protectHookContextPath(temporaryName, false); err != nil {
		return "", "", fmt.Errorf("protect hook %s context: %w", spec.Phase, err)
	}
	if _, err := temporary.Write(content); err != nil {
		return "", "", fmt.Errorf("write hook %s context: %w", spec.Phase, err)
	}
	if err := temporary.Sync(); err != nil {
		return "", "", fmt.Errorf("sync hook %s context: %w", spec.Phase, err)
	}
	if err := temporary.Close(); err != nil {
		return "", "", fmt.Errorf("close hook %s context: %w", spec.Phase, err)
	}
	contextFile := filepath.Join(contextDirectory, "context.json")
	if err := os.Rename(temporaryName, contextFile); err != nil {
		return "", "", fmt.Errorf("publish hook %s context: %w", spec.Phase, err)
	}
	removeTemporary = false
	removeDirectory = false
	return contextFile, contextDirectory, nil
}

func removeEnvironmentKeys(environment []string, names ...string) []string {
	removed := make(map[string]bool, len(names))
	for _, name := range names {
		if runtime.GOOS == "windows" {
			name = strings.ToUpper(name)
		}
		removed[name] = true
	}
	result := make([]string, 0, len(environment))
	for _, item := range environment {
		name, _, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		if runtime.GOOS == "windows" {
			name = strings.ToUpper(name)
		}
		if !removed[name] {
			result = append(result, item)
		}
	}
	return result
}
