# Next Steps

Now that you have Wails installed, you can start exploring the alpha version.

The best place to start is the `examples` directory in the Wails repository.
This contains a number of examples that you can run and play with.

## Running an example

To run an example, you can simply use:

```shell
go run .
```

in the example directory.

## Creating a new project

To create a new project, you can use the `wails3 init` command. This will create
a new project in the current directory.

Wails3 uses [Task](https://taskfile.dev) as its build system by default,
although there is no reason why you can't use your own build system, or use
`go build` directly. Wails has the task build system built in and can be run
using `wails3 task`.

If you look through the `Taskfile.yaml` file, you will see that there are a
number of tasks defined. The most important one is the `build` task. This is the
task that is run when you use `wails3 build`.

The task file is unlikely to be complete and is subject to change over time.

## Building a project

To build a project, you can use the `wails3 build` command. This is a shortcut
for `wails3 task build`.
