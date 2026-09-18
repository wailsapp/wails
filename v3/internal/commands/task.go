package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/buildwarnings"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/wake"

	"github.com/wailsapp/task/v3"
	"github.com/wailsapp/task/v3/taskfile/ast"
)

// wakeRoutableInvocation reports whether `wails3 task <name>` can be safely
// dispatched to wake. wake doesn't model the introspection flags (list,
// dry-run, summary, watch) — those bypass wake and fall through to the
// embedded task runtime below.
//
// `--dir` / `--taskfile` also bypass wake: they reshape where the Taskfile
// is loaded from, and the wake fast path currently always loads from the
// process cwd. Routing those through wake would silently run the task from
// the wrong directory or ignore the user's chosen entrypoint, so let the
// embedded runtime keep its semantics for that case.
func wakeRoutableInvocation(o *RunTaskOptions) bool {
	if o.List || o.ListAll || o.ListJSON || o.Status || o.Watch ||
		o.Dry || o.Summary {
		return false
	}
	if o.Dir != "" || o.EntryPoint != "" {
		return false
	}
	return true
}

// parseCLIVars splits the trailing "KEY=VALUE" arguments off `wails3 task`
// invocations into a map for wake.ExecuteOptions.Vars. Args without an "="
// are silently skipped (matching the embedded task runtime's behaviour).
func parseCLIVars(args []string) map[string]string {
	vars := make(map[string]string)
	for _, a := range args {
		if k, v, ok := strings.Cut(a, "="); ok {
			vars[k] = v
		}
	}
	return vars
}

// BuildSettings contains the CLI build settings
var BuildSettings = map[string]string{}

func fatal(message string) {
	buildwarnings.FlushAndPrint()
	term.Error(message)
	os.Exit(1)
}

type RunTaskOptions struct {
	Name             string `pos:"1"`
	Help             bool   `name:"h" description:"shows Task usage"`
	Init             bool   `name:"i" description:"creates a new Taskfile.yml"`
	List             bool   `name:"list" description:"tasks with description of current Taskfile"`
	ListAll          bool   `name:"list-all" description:"lists tasks with or without a description"`
	ListJSON         bool   `name:"json" description:"formats task list as json"`
	Status           bool   `name:"status" description:"exits with non-zero exit code if any of the given tasks is not up-to-date"`
	Force            bool   `name:"f" description:"forces execution even when the task is up-to-date"`
	Watch            bool   `name:"w" description:"enables watch of the given task"`
	Verbose          bool   `name:"v" description:"enables verbose mode"`
	Version          bool   `name:"version" description:"prints version"`
	Silent           bool   `name:"s" description:"disables echoing"`
	Parallel         bool   `name:"p" description:"executes tasks provided on command line in parallel"`
	Dry              bool   `name:"dry" description:"compiles and prints tasks in the order that they would be run, without executing them"`
	Summary          bool   `name:"summary" description:"show summary about a task"`
	ExitCode         bool   `name:"x" description:"pass-through the exit code of the task command"`
	Dir              string `name:"dir" description:"sets directory of execution"`
	EntryPoint       string `name:"taskfile" description:"choose which Taskfile to run."`
	OutputName       string `name:"output" description:"sets output style: [interleaved|group|prefixed]"`
	OutputGroupBegin string `name:"output-group-begin" description:"message template to print before a task's grouped output"`
	OutputGroupEnd   string `name:"output-group-end" description:"message template to print after a task's grouped output"`
	Color            bool   `name:"c" description:"colored output. Enabled by default. Set flag to false or use NO_COLOR=1 to disable" default:"true"`
	Concurrency      int    `name:"C" description:"limit number tasks to run concurrently"`
	Interval         int64  `name:"interval" description:"interval to watch for changes"`
}

func RunTask(options *RunTaskOptions, otherArgs []string) error {

	if options.Version {
		ver := BuildSettings["mod.github.com/wailsapp/task/v3"]
		fmt.Println("Task Version:", ver)
		return nil
	}

	if options.Init {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		return task.InitTaskfile(os.Stdout, wd)
	}

	// When WAILS_USE_WAKE=true, route arbitrary `wails3 task <name>`
	// invocations through wake too, for consistency with `wails3 build` /
	// `wails3 package` / `wails3 sign` (which all dispatch via wrapTask).
	// Operations wake doesn't support (--list, --watch, --dry, etc.) fall
	// through to the embedded task runtime below. wake itself will fall
	// back to the embedded runtime if the Taskfile uses an unsupported
	// feature (dotenv, requires, ...), so this path is always safe.
	if useWake() && options.Name != "" && wakeRoutableInvocation(options) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		vars := parseCLIVars(otherArgs)
		opts := wake.ExecuteOptions{
			Dir:      dir,
			Verb:     options.Name,
			Vars:     vars,
			Verbose:  options.Verbose || os.Getenv("WAKE_VERBOSE") != "",
			Silent:   options.Silent || os.Getenv("WAKE_SILENT") != "",
			Debug:    os.Getenv("WAKE_DEBUG") != "",
			Parallel: options.Parallel || os.Getenv("WAKE_SERIAL") == "",
			Force:    options.Force || os.Getenv("WAKE_FORCE") != "",
		}
		return wake.Execute(options.Name, opts)
	}


	if options.Dir != "" && options.EntryPoint != "" {
		return fmt.Errorf("task: You can't set both --dir and --taskfile")
	}

	if options.EntryPoint != "" {
		options.Dir = filepath.Dir(options.EntryPoint)
		options.EntryPoint = filepath.Base(options.EntryPoint)
	}

	if options.OutputName != "group" {
		if options.OutputGroupBegin != "" {
			return fmt.Errorf("task: You can't set --output-group-begin without --output=group")
		}
		if options.OutputGroupEnd != "" {
			return fmt.Errorf("task: You can't set --output-group-end without --output=group")
		}
	}

	e := task.Executor{
		Force:               options.Force,
		Watch:               options.Watch,
		Verbose:             options.Verbose,
		Silent:              options.Silent,
		Dir:                 options.Dir,
		Dry:                 options.Dry,
		Entrypoint:          options.EntryPoint,
		Summary:             options.Summary,
		Parallel:            options.Parallel,
		Color:               options.Color,
		Concurrency:         options.Concurrency,
		Interval:            time.Duration(options.Interval) * time.Second,
		DisableVersionCheck: true,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	listOptions := task.NewListOptions(options.List, options.ListAll, options.ListJSON, false)
	if err := listOptions.Validate(); err != nil {
		fatal(err.Error())
	}

	if listOptions.ShouldListTasks() && options.Silent {
		e.ListTaskNames(options.ListAll)
		return nil
	}

	if err := e.Setup(); err != nil {
		fatal(err.Error())
	}

	if listOptions.ShouldListTasks() {
		if foundTasks, err := e.ListTasks(listOptions); !foundTasks || err != nil {
			os.Exit(1)
		}
		return nil
	}

	if e.Taskfile != nil && e.Taskfile.Vars == nil {
		e.Taskfile.Vars = &ast.Vars{}
	}

	// Parse task name and CLI variables from otherArgs or os.Args
	var tasksAndVars []string

	// Check if we have a task name specified in options
	if options.Name != "" {
		// If task name is provided via options, use it and treat otherArgs as CLI variables
		tasksAndVars = append([]string{options.Name}, otherArgs...)
	} else if len(otherArgs) > 0 {
		// Use otherArgs directly if provided
		tasksAndVars = otherArgs
	} else {
		// Fall back to parsing os.Args for backward compatibility
		var index int
		var arg string
		for index, arg = range os.Args[2:] {
			if !strings.HasPrefix(arg, "-") {
				break
			}
		}

		for _, taskAndVar := range os.Args[index+2:] {
			if taskAndVar == "--" {
				break
			}
			tasksAndVars = append(tasksAndVars, taskAndVar)
		}
	}

	// Default task
	if len(tasksAndVars) == 0 {
		tasksAndVars = []string{"default"}
	}

	// Parse task name and CLI variables
	taskName := tasksAndVars[0]
	cliVars := tasksAndVars[1:]

	// Create call with CLI variables
	call := &ast.Call{
		Task: taskName,
		Vars: &ast.Vars{},
	}

	// Parse CLI variables (format: KEY=VALUE)
	for _, v := range cliVars {
		if strings.Contains(v, "=") {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				call.Vars.Set(parts[0], ast.Var{
					Value: parts[1],
				})
				if e.Taskfile != nil {
					e.Taskfile.Vars.Set(parts[0], ast.Var{
						Value: parts[1],
					})
				}
			}
		}
	}

	if err := e.RunTask(context.Background(), call); err != nil {
		fatal(err.Error())
	}
	return nil
}
