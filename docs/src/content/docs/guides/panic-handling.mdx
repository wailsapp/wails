---
title: Handling Panics
description: How to handle panics in your Wails application
---

In Go applications, panics can occur during runtime when something unexpected happens. This guide explains how to handle panics both in general Go code and specifically in your Wails application.

## Understanding Panics in Go

Before diving into Wails-specific panic handling, it's essential to understand how panics work in Go:

1. Panics are for unrecoverable errors that shouldn't happen during normal operation
2. When a panic occurs in a goroutine, only that goroutine is affected
3. Panics can be recovered using `defer` and `recover()`

Here's a basic example of panic handling in Go:

```go
func doSomething() {
    // Deferred functions run even when a panic occurs
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    // Your code that might panic
    panic("something went wrong")
}
```

For more detailed information about panic and recover in Go, see the [Go Blog: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover).

## Panic Handling in Wails

Wails automatically handles panics that occur in your Service methods when they are called from the frontend. This means you don't need to add panic recovery to these methods - Wails will catch the panic and process it through your configured panic handler.

The panic handler is specifically designed to catch:
- Panics in bound service methods called from the frontend
- Internal panics from the Wails runtime

For other scenarios, such as background goroutines or standalone Go code, you should handle panics yourself using Go's standard panic recovery mechanisms.

## The PanicDetails Struct

When a panic occurs, Wails captures important information about the panic in a `PanicDetails` struct:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

This structure provides comprehensive information about the panic:
- `StackTrace`: A formatted string showing the call stack that led to the panic
- `Error`: The actual error or panic message
- `Time`: The exact time when the panic occurred
- `FullStackTrace`: The complete stack trace including runtime frames

:::note[Panics in Service Code]

When panics are caught in your Service code after being called from the frontend, the stack trace is trimmed to focus on exactly where in your code the panic occurred.
If you want to see the full stack trace, you can use the `FullStackTrace` field.

:::

## Default Panic Handler

If you don't specify a custom panic handler, Wails will use its default handler which outputs error information in a formatted log message and then quits.
For example:

```
************************ FATAL ******************************
* There has been a catastrophic failure in your application *
********************* Error Details *************************
panic error: oh no! something went wrong deep in my service! :(
main.(*WindowService).call2
	at E:/wails/v3/examples/panic-handling/main.go:23
main.(*WindowService).call1
	at E:/wails/v3/examples/panic-handling/main.go:19
main.(*WindowService).GeneratePanic
	at E:/wails/v3/examples/panic-handling/main.go:15
*************************************************************
```

## Custom Panic Handler

You can implement your own panic handler by setting the `PanicHandler` option when creating your application. Here's an example:

```go
app := application.New(application.Options{
    Name: "My App",
    PanicHandler: func(panicDetails *application.PanicDetails) {
        fmt.Printf("*** Custom Panic Handler ***\n")
        fmt.Printf("Time: %s\n", panicDetails.Time)
        fmt.Printf("Error: %s\n", panicDetails.Error)
        fmt.Printf("Stacktrace: %s\n", panicDetails.StackTrace)
        fmt.Printf("Full Stacktrace: %s\n", panicDetails.FullStackTrace)
        
        // You could also:
        // - Log to a file
        // - Send to a crash reporting service
        // - Show a user-friendly error dialog
        // - Attempt to recover or restart the application
    },
})
```

For a complete working example of panic handling in a Wails application, see the panic-handling example in `v3/examples/panic-handling`.

## Final Notes

Remember that the Wails panic handler is specifically for managing panics in bound methods and internal runtime errors. For other parts of your application, you should use Go's standard error handling patterns and panic recovery mechanisms where appropriate. As with all Go applications, it's better to prevent panics through proper error handling where possible.
