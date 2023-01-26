package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/go-task/task/v3/args"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

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
		if options.OutputGroupBegin != "" {
			return fmt.Errorf("task: You can't set --output-group-end without --output=group")
		}
	}

	e := task.Executor{
		Force:       options.Force,
		Watch:       options.Watch,
		Verbose:     options.Verbose,
		Silent:      options.Silent,
		Dir:         options.Dir,
		Dry:         options.Dry,
		Entrypoint:  options.EntryPoint,
		Summary:     options.Summary,
		Parallel:    options.Parallel,
		Color:       options.Color,
		Concurrency: options.Concurrency,
		Interval:    time.Duration(options.Interval) * time.Second,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		OutputStyle: taskfile.Output{
			Name: options.OutputName,
			Group: taskfile.OutputGroup{
				Begin: options.OutputGroupBegin,
				End:   options.OutputGroupEnd,
			},
		},
	}

	var listOptions = task.NewListOptions(options.List, options.ListAll, options.ListJSON)
	if err := listOptions.Validate(); err != nil {
		log.Fatal(err)
	}

	if (listOptions.ShouldListTasks()) && options.Silent {
		e.ListTaskNames(options.ListAll)
		return nil
	}

	if err := e.Setup(); err != nil {
		log.Fatal(err)
	}
	v, err := e.Taskfile.ParsedVersion()
	if err != nil {
		return err
	}

	if listOptions.ShouldListTasks() {
		if foundTasks, err := e.ListTasks(listOptions); !foundTasks || err != nil {
			os.Exit(1)
		}
		return nil
	}

	var (
		calls   []taskfile.Call
		globals *taskfile.Vars
	)

	var taskAndVars []string
	for _, taskAndVar := range os.Args[2:] {
		if taskAndVar == "--" {
			break
		}
		taskAndVars = append(taskAndVars, taskAndVar)
	}

	if len(taskAndVars) > 0 && len(otherArgs) > 0 {
		if taskAndVars[0] == otherArgs[0] {
			otherArgs = otherArgs[1:]
		}
	}

	if v >= 3.0 {
		calls, globals = args.ParseV3(taskAndVars...)
	} else {
		calls, globals = args.ParseV2(taskAndVars...)
	}

	globals.Set("CLI_ARGS", taskfile.Var{Static: strings.Join(otherArgs, " ")})
	e.Taskfile.Vars.Merge(globals)

	if !options.Watch {
		e.InterceptInterruptSignals()
	}

	ctx := context.Background()

	if options.Status {
		return e.Status(ctx, calls...)
	}

	if err := e.Run(ctx, calls...); err != nil {
		pterm.Error.Println(err.Error())

		if options.ExitCode {
			if err, ok := err.(*task.TaskRunError); ok {
				os.Exit(err.ExitCode())
			}
		}
		os.Exit(1)
	}
	return nil
}
