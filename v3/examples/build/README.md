# Build

Wails has adopted [Taskfile](https://taskfile.dev) as its build tool. This is optional
and any build tool can be used. However, Taskfile is a great tool, and we recommend it.

The Wails CLI has built-in integration with Taskfile so the standalone version is not a
requirement.

## Building

To build the example, run:

```bash
wails3 task build
```

> **Note:** Bare `go build .` outputs a binary named `build` (matching the directory
> name). Use `wails3 task build` or an explicit output path instead:
> `go build -o bin/buildtest .`

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | Working |
| Linux    |         |