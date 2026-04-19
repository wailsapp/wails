# debug

The `debug` package provides runtime diagnostics and crash analysis utilities for Wails applications.

## Usage

```go
collector := debug.New()
report, err := collector.Run()
```

## Features

- System info from doctor-ng
- Process diagnostics (memory, threads, modules)
- Windows-first crash analysis