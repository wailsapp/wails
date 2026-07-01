package commands

import (
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/buildwarnings"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/wake"
)

// runTaskFunc is a variable to allow mocking in tests
var runTaskFunc = RunTask

// validPlatforms for GOOS
var validPlatforms = map[string]bool{
	"windows": true,
	"darwin":  true,
	"linux":   true,
}

// rootDispatchTasks are the verbs that this wrapper routes through the
// top-level task in the generated root Taskfile, which dispatches to the
// platform-specific task via the GOOS variable (e.g. `build` -> `{{.GOOS}}:build`).
// For these we run the root task and pass GOOS as a variable, so user
// customisations in the root Taskfile are honoured for both native and
// cross-compilation builds. Only `build` and `package` go through this wrapper:
// the root Taskfile also defines a `run` dispatch task, but `run` is invoked by
// `wails3 dev` directly rather than through here. Verbs without a root task
// (e.g. `sign`, which only exists per-platform) always target the
// platform-specific task directly.
var rootDispatchTasks = map[string]bool{
	"build":   true,
	"package": true,
}

const (
	// mcpEnvVar enables the MCP service build tag when set to a truthy value.
	mcpEnvVar = "WAILS_MCP"
	// mcpBuildTag is the Go build tag that compiles in the MCP service.
	mcpBuildTag = "mcp"
)

// envTags returns the build tags implied by environment variables.
func envTags() []string {
	var tags []string
	switch strings.ToLower(strings.TrimSpace(os.Getenv(mcpEnvVar))) {
	case "1", "true", "on", "yes":
		tags = append(tags, mcpBuildTag)
	}
	return tags
}

// mergeTags appends extra tags to a comma-separated tag list, skipping duplicates.
func mergeTags(tags string, extra ...string) string {
	existing := strings.Split(tags, ",")
	for _, tag := range extra {
		if !slices.Contains(existing, tag) {
			existing = append(existing, tag)
		}
	}
	return strings.Trim(strings.Join(existing, ","), ",")
}

func Build(buildFlags *flags.Build, otherArgs []string) error {
	buildFlags.Tags = mergeTags(buildFlags.Tags, envTags()...)
	if buildFlags.Tags != "" {
		otherArgs = append(otherArgs, "EXTRA_TAGS="+buildFlags.Tags)
	}
	if buildFlags.Obfuscated {
		otherArgs = append(otherArgs, "OBFUSCATED=true")
	}
	if buildFlags.GarbleArgs != "" {
		otherArgs = append(otherArgs, "GARBLE_ARGS="+buildFlags.GarbleArgs)
	}
	return wrapTask("build", otherArgs)
}

func Package(_ *flags.Package, otherArgs []string) error {
	return wrapTask("package", otherArgs)
}

func SignWrapper(_ *flags.SignWrapper, otherArgs []string) error {
	return wrapTask("sign", otherArgs)
}

func wrapTask(action string, otherArgs []string) error {
	// Match the banner other wails3 commands print; the footer is restored by
	// leaving DisableFooter at its default so printFooter runs on exit.
	term.Header(title(action))

	// Create a per-invocation warnings file so subprocess commands (e.g.
	// wails3 tool has-cc) can append deprecation notices that we collect
	// and print after the task finishes.
	if f, err := os.CreateTemp("", "wails-build-warnings-*"); err == nil {
		f.Close()
		prev, hadPrev := os.LookupEnv(buildwarnings.EnvVar)
		os.Setenv(buildwarnings.EnvVar, f.Name())
		defer func() {
			buildwarnings.FlushAndPrint()
			if hadPrev {
				os.Setenv(buildwarnings.EnvVar, prev)
			} else {
				os.Unsetenv(buildwarnings.EnvVar)
			}
		}()
	}

	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	var remainingArgs []string

	for _, arg := range otherArgs {
		switch {
		case strings.HasPrefix(arg, "GOOS="):
			goos = strings.TrimPrefix(arg, "GOOS=")
		case strings.HasPrefix(arg, "GOARCH="):
			goarch = strings.TrimPrefix(arg, "GOARCH=")
		default:
			remainingArgs = append(remainingArgs, arg)
		}
	}

	// platformTaskName always targets the platform-specific task (e.g.
	// "linux:build"). The experimental wake runner is built against this concrete
	// name, so it keeps using it unconditionally.
	platformTaskName := action
	if validPlatforms[goos] {
		platformTaskName = goos + ":" + action
	}

	// Pass GOOS/ARCH through as Taskfile variables. The root build/package/run
	// tasks dispatch on {{.GOOS}}, so running the root task (rather than the
	// platform-prefixed one) means any customisations in the root Taskfile are
	// honoured, for both native and cross-compilation builds. Verbs without a
	// root task (e.g. `sign`) still target the platform task directly. See #5615.
	remainingArgs = append(remainingArgs, "GOOS="+goos, "ARCH="+goarch)

	taskName := platformTaskName
	if rootDispatchTasks[action] {
		taskName = action
	}

	if useWake() {
		return runWakeTask(action, platformTaskName, goos, goarch, remainingArgs)
	}

	newArgs := []string{"wails3", "task", taskName}
	newArgs = append(newArgs, remainingArgs...)
	os.Args = newArgs
	return runTaskFunc(&RunTaskOptions{Name: taskName}, remainingArgs)
}

func useWake() bool {
	return os.Getenv("WAILS_USE_WAKE") == "true"
}

// title capitalises an action ("build" -> "Build") for the command banner.
func title(action string) string {
	if action == "" {
		return action
	}
	return strings.ToUpper(action[:1]) + action[1:]
}

func runWakeTask(verb, taskName, goos, goarch string, cliVars []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	vars := make(map[string]string)
	for _, v := range cliVars {
		if strings.Contains(v, "=") {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				vars[parts[0]] = parts[1]
			}
		}
	}

	opts := wake.ExecuteOptions{
		Dir:      dir,
		Platform: goos,
		Arch:     goarch,
		Verb:     verb,
		Vars:     vars,
		Verbose:  os.Getenv("WAKE_VERBOSE") != "",
		Silent:   os.Getenv("WAKE_SILENT") != "",
		Debug:    os.Getenv("WAKE_DEBUG") != "",
		// Parallel execution is the default. Set WAKE_SERIAL=true to opt out
		// (useful when debugging task ordering or when stdout interleaving
		// from sibling steps would muddle a specific investigation).
		Parallel: os.Getenv("WAKE_SERIAL") == "",
		// WAKE_FORCE=true skips every cache lookup, both the Taskfile
		// sources/generates/status check and the implicit native-Go cache.
		// Use when you want a true "clean" build without rm -rf .wake/.
		Force: os.Getenv("WAKE_FORCE") != "",
	}

	return wake.Execute(taskName, opts)
}
