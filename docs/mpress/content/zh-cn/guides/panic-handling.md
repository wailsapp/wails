---
title: "处理 panic"
description: "如何处理 Wails 应用程序中的 panic"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

在 Go 应用程序中，如果运行时发生意外情况，可能会出现 panic。本指南介绍如何在常规 Go 代码中以及专门在 Wails 应用程序中处理 panic。

## 了解 Go 中的 panic

在深入了解 Wails 特有的 panic 处理方式之前，必须先了解 panic 在 Go 中的工作原理：

1. panic 用于表示正常运行期间不应发生且无法恢复的错误
2. 当某个 goroutine 中发生 panic 时，只有该 goroutine 会受到影响
3. 可以使用 `defer` 和 `recover()` 从 panic 中恢复

下面是一个在 Go 中处理 panic 的基本示例：

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

有关 Go 中 panic 和 recover 的更多详细信息，请参阅[Go 博客：Defer、Panic 和 Recover](https://go.dev/blog/defer-panic-and-recover)。

## Wails 中的 panic 处理

从前端调用 Service 方法时，Wails 会自动处理这些方法中发生的 panic。这意味着你无需为这些方法添加 panic 恢复逻辑——Wails 会捕获 panic，并通过你配置的 panic 处理程序进行处理。

panic 处理程序专门用于捕获：

- 从前端调用的已绑定服务方法中发生的 panic
- Wails 运行时内部的 panic

对于后台 goroutine 或独立 Go 代码等其他场景，你应使用 Go 的标准 panic 恢复机制自行处理 panic。

## PanicDetails 结构体

发生 panic 时，Wails 会将有关该 panic 的重要信息捕获到 `PanicDetails` 结构体中：

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

此结构体提供关于 panic 的完整信息：

- `StackTrace`：显示导致 panic 的调用栈的格式化字符串
- `Error`：实际的错误或 panic 消息
- `Time`：panic 发生的确切时间
- `FullStackTrace`：包含运行时栈帧的完整堆栈跟踪

@note{type="info" title="Service 代码中的 panic"}
从前端调用 Service 代码并捕获其中的 panic 后，系统会裁剪堆栈跟踪，使其重点显示代码中发生 panic 的确切位置。 如果你想查看完整的堆栈跟踪，可以使用 `FullStackTrace` 字段。

@end

## 默认 panic 处理程序

如果你未指定自定义 panic 处理程序，Wails 将使用默认处理程序，以格式化日志消息的形式输出错误信息，然后退出。 例如：

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

## 自定义 panic 处理程序

创建应用程序时，可以通过设置 `PanicHandler` 选项来实现自己的 panic 处理程序。示例如下：

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

有关在 Wails 应用程序中处理 panic 的完整可运行示例，请参阅 `v3/examples/panic-handling` 中的 panic-handling 示例。

## 最后说明

请记住，Wails panic 处理程序专门用于处理已绑定方法中的 panic 和运行时内部错误。对于应用程序的其他部分，应在适当情况下使用 Go 的标准错误处理模式和 panic 恢复机制。与所有 Go 应用程序一样，应尽可能通过正确的错误处理来防止 panic。
