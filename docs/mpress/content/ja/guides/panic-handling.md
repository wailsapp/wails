---
title: "パニックの処理"
description: "Wails アプリケーションでパニックを処理する方法"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Go アプリケーションでは、予期しない事態が発生すると、実行時にパニックが起こることがあります。このガイドでは、一般的な Go コードと Wails アプリケーション固有のコードの両方でパニックを処理する方法を説明します。

## Go のパニックについて

Wails 固有のパニック処理について説明する前に、Go でパニックがどのように機能するかを理解しておくことが重要です。

1. パニックは、通常の動作中には発生しないはずの回復不能なエラーに使用します
2. パニックが発生すると、現在の goroutine のスタックが巻き戻されます。回復されないまま最上位まで達すると、ほかの goroutine も含めてプロセスが終了します
3. 回復処理は、パニックが発生した goroutine と同じ goroutine 内の遅延実行関数で実行する必要があります

Go でパニックを処理する基本的な例を次に示します。

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

Go の panic と recover の詳細については、[Go ブログ：Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)を参照してください。

## Wails でのパニック処理

フロントエンドから Service メソッドが呼び出された際に発生したパニックは、Wails が自動的に処理します。そのため、これらのメソッドにパニック回復処理を追加する必要はありません。Wails がパニックを捕捉し、設定済みのパニックハンドラーで処理します。

パニックハンドラーは、特に次のパニックを捕捉するように設計されています。

- フロントエンドから呼び出された、バインド済みのサービスメソッド内のパニック
- Wails の回復処理を明示的に設定しているランタイムパス内のパニック。これは、すべての goroutine を対象とするプロセス全体のハンドラーではありません

バックグラウンドの goroutine や独立した Go コードなど、その他の場面では、Go 標準のパニック回復機構を使用して、パニックを自身で処理する必要があります。

## PanicDetails 構造体

パニックが発生すると、Wails はパニックに関する重要な情報を `PanicDetails` 構造体に記録します。

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

この構造体は、パニックに関する包括的な情報を提供します。

- `StackTrace`：パニックに至った呼び出しスタックを示す、整形済みの文字列
- `Error`：実際のエラーまたはパニックメッセージ
- `Time`：パニックが発生した正確な時刻
- `FullStackTrace`：ランタイムフレームを含む完全なスタックトレース

@note{type="info" title="Service コード内のパニック"}
フロントエンドから呼び出された後に Service コード内で発生したパニックが捕捉されると、スタックトレースは、コード内でパニックが発生した正確な位置に焦点を当てるように切り詰められます。完全なスタックトレースを確認するには、`FullStackTrace` フィールドを使用できます。

@end

## デフォルトのパニックハンドラー

カスタムパニックハンドラーを指定しない場合、Wails はデフォルトのハンドラーを使用します。このハンドラーは、エラー情報を整形されたログメッセージとして出力してから終了します。例：

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

## カスタムパニックハンドラー

アプリケーションの作成時に `PanicHandler` オプションを設定すると、独自のパニックハンドラーを実装できます。例：

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

Wails アプリケーションでのパニック処理について、完全に動作する例は、`v3/examples/panic-handling` にある panic-handling の例を参照してください。

## 最後に

Wails のパニックハンドラーは、バインド済みメソッド内のパニックと内部ランタイムエラーを処理するためのものであることに注意してください。アプリケーションのその他の部分では、必要に応じて Go 標準のエラー処理パターンとパニック回復機構を使用してください。すべての Go アプリケーションと同様に、可能な限り適切なエラー処理によってパニックを防ぐことを推奨します。
