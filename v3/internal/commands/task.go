package commands

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/term"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/task/v3"
	"github.com/wailsapp/task/v3/taskfile/ast"
)

// BuildSettings contains the CLI build settings
var BuildSettings = map[string]string{}

func fatal(message string) {
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

	var index int
	var arg string
	for index, arg = range os.Args[2:] {
		if !strings.HasPrefix(arg, "-") {
			break
		}
	}

	var tasksAndVars []string
	for _, taskAndVar := range os.Args[index+2:] {
		if taskAndVar == "--" {
			break
		}
		tasksAndVars = append(tasksAndVars, taskAndVar)
	}

	if len(tasksAndVars) > 0 && len(otherArgs) > 0 {
		if tasksAndVars[0] == otherArgs[0] {
			otherArgs = otherArgs[1:]
		}
	}

	if err := e.RunTask(context.Background(), &ast.Call{Task: tasksAndVars[0]}); err != nil {
		fatal(err.Error())
	}
	return nil
}
