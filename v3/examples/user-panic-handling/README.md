# User Panic Handling Example

Wails' built-in `PanicHandler` catches panics in two places:

- Bound service methods invoked from the frontend
- Internal wails runtime callbacks

Panics in goroutines **you** spawn are not covered — they follow standard Go
behaviour and crash the process. This example shows how to funnel those
panics into the same handler wails uses, so every panic in the app — no
matter where it originates — ends up in one place.

## What it demonstrates

1. **Bound-method path.** The button in the UI calls
   `WindowService.GeneratePanic()`. Wails' deferred `handlePanic()` wrapper
   recovers, builds a `PanicDetails`, and calls the `PanicHandler` registered
   in `application.Options`.

2. **User goroutine path.** `BackgroundWorker.tickLoop` runs in a goroutine
   spawned by user code. It defers `recoverAndReport()`, which recovers, builds
   a `PanicDetails` manually (using `runtime/debug.Stack()` for the trace),
   and calls the same handler.

Both paths end up at `reportPanic(*application.PanicDetails)`.

## Running

```bash
go run .
```

Expected output within ~3 seconds:

```
*** PANIC ***
Source: user goroutine (recovered via recoverAndReport)
Time:   ...
Error:  ticker exploded after 3 ticks
Stack:
  ...
```

Then click the button in the window and you'll see a second panic logged,
this time marked `Source: wails runtime (...)`.

## The user-side helper

```go
func recoverAndReport() {
    r := recover()
    if r == nil { return }
    err, ok := r.(error)
    if !ok { err = fmt.Errorf("%v", r) }
    stack := string(debug.Stack())
    reportPanic(&application.PanicDetails{
        Error:          err,
        Time:           time.Now(),
        StackTrace:     stack,
        FullStackTrace: stack,
    })
}
```

Put `defer recoverAndReport()` as the first line of every goroutine you
spawn in user code.

## Why wails can't do this for you

Wails can only defer `handlePanic()` inside code paths it owns. When you
write `go func() { ... }()`, wails never sees that goroutine, so it cannot
inject a deferred recover into it. The user-side helper is the smallest
pattern that fixes that — and it routes through your existing
`PanicHandler` so reporting stays centralised.
