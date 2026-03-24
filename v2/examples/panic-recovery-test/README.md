# Panic Recovery Test

This example demonstrates the Linux signal handler issue (#3965) and verifies the fix using `runtime.ResetSignalHandlers()`.

## The Problem

On Linux, WebKit installs signal handlers without the `SA_ONSTACK` flag, which prevents Go from recovering panics caused by nil pointer dereferences (SIGSEGV). Without the fix, the application crashes with:

```
signal 11 received but handler not on signal stack
fatal error: non-Go code set up signal handler without SA_ONSTACK flag
```

## The Solution

Call `runtime.ResetSignalHandlers()` immediately before code that might panic:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

go func() {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Recovered: %v", err)
        }
    }()
    runtime.ResetSignalHandlers()
    // Code that might panic...
}()
```

## How to Reproduce

### Prerequisites

- Linux with WebKit2GTK 4.1 installed
- Go 1.21+
- Wails CLI

### Steps

1. Build the example:
   ```bash
   cd v2/examples/panic-recovery-test
   wails build -tags webkit2_41
   ```

2. Run the application:
   ```bash
   ./build/bin/panic-recovery-test
   ```

3. Wait ~10 seconds (the app auto-calls `Greet` after 5s, then waits another 5s before the nil pointer dereference)

### Expected Result (with fix)

The panic is recovered and you see:
```
------------------------------"invalid memory address or nil pointer dereference"
```

The application continues running.

### Without the fix

Comment out the `runtime.ResetSignalHandlers()` call in `app.go` and rebuild. The application will crash with a fatal signal 11 error.

## Files

- `app.go` - Contains the `Greet` function that demonstrates panic recovery
- `frontend/src/main.js` - Auto-calls `Greet` after 5 seconds to trigger the test

## Related

- Issue: https://github.com/wailsapp/wails/issues/3965
- Original fix PR: https://github.com/wailsapp/wails/pull/2152
