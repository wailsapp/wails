---
title: "處理 Panic"
description: "如何處理 Wails 應用程式中的 panic"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

在 Go 應用程式中，執行階段發生非預期狀況時，可能會出現 panic。本指南說明如何在一般 Go 程式碼中，以及特別是在 Wails 應用程式中處理 panic。

## 瞭解 Go 中的 Panic

在深入瞭解 Wails 特有的 panic 處理方式之前，請務必先瞭解 panic 在 Go 中的運作方式：

1. Panic 用於處理正常運作期間不應發生且無法復原的錯誤
2. 當 goroutine 中發生 panic 時，只有該 goroutine 會受到影響
3. 可以使用`defer`和`recover()`從 panic 中復原

以下是在 Go 中處理 panic 的基本範例：

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

如需進一步瞭解 Go 中的 panic 和 recover，請參閱[Go 部落格：Defer、Panic 與 Recover](https://go.dev/blog/defer-panic-and-recover)。

## Wails 中的 Panic 處理

從前端呼叫 Service 方法時，如果方法中發生 panic，Wails 會自動處理。因此，你不需要在這些方法中加入 panic 復原機制；Wails 會攔截 panic，並交由你設定的 panic 處理常式處理。

Panic 處理常式專門用來攔截：

- 從前端呼叫的已繫結服務方法中發生的 panic
- Wails 執行階段內部發生的 panic

對於其他情境，例如背景 goroutine 或獨立的 Go 程式碼，你應使用 Go 的標準 panic 復原機制自行處理 panic。

## PanicDetails 結構

發生 panic 時，Wails 會將其重要資訊擷取至`PanicDetails`結構中：

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

此結構提供下列完整的 panic 資訊：

- `StackTrace`：顯示導致 panic 之呼叫堆疊的格式化字串
- `Error`：實際的錯誤或 panic 訊息
- `Time`：panic 發生的確切時間
- `FullStackTrace`：包含執行階段堆疊框架的完整堆疊追蹤

@note{type="info" title="Service 程式碼中的 Panic"}
從前端呼叫 Service 程式碼後，如果其中的 panic 被攔截，系統會裁剪堆疊追蹤，使其聚焦於程式碼中發生 panic 的確切位置。 若要查看完整的堆疊追蹤，可以使用`FullStackTrace`欄位。

@end

## 預設 Panic 處理常式

如果未指定自訂 panic 處理常式，Wails 會使用預設處理常式，以格式化的記錄訊息輸出錯誤資訊，然後結束應用程式。 例如：

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

## 自訂 Panic 處理常式

建立應用程式時，可以設定`PanicHandler`選項來實作自己的 panic 處理常式。範例如下：

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

如需 Wails 應用程式中 panic 處理方式的完整可運作範例，請參閱`v3/examples/panic-handling`中的 panic-handling 範例。

## 最後注意事項

請記住，Wails panic 處理常式專門用於管理已繫結方法中的 panic 和執行階段內部錯誤。對於應用程式的其他部分，你應在適當時使用 Go 的標準錯誤處理模式和 panic 復原機制。與所有 Go 應用程式一樣，最好盡可能透過妥善的錯誤處理來避免發生 panic。
