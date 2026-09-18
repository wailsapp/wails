---
title: "パニックの処理"
description: "Wails アプリケーションでパニックを処理する方法"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Go アプリケーションでは、予期しない事態が発生すると、実行時にパニックが起こることがあります。このガイドでは、一般的な Go コードと Wails アプリケーション固有のコードの両方でパニックを処理する方法について説明します。

## Go のパニックを理解する

Wails 固有のパニック処理について詳しく見る前に、Go でパニックがどのように機能するかを理解しておくことが重要です。

1. パニックは、通常の動作中には発生しないはずの回復不能なエラーに使用します
2. goroutine でパニックが発生した場合、影響を受けるのはその goroutine だけです
3. パニックからは、`defer`と`recover()`を使用して回復できます

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

Go の panic と recover の詳細については、[Go Blog：Defer、Panic、Recover](https://go.dev/blog/defer-panic-and-recover)を参照してください。

## Wails でのパニック処理

Wails は、フロントエンドから呼び出された Service メソッドで発生するパニックを自動的に処理します。そのため、これらのメソッドにパニックからの回復処理を追加する必要はありません。Wails がパニックを捕捉し、設定済みのパニックハンドラーで処理します。

パニックハンドラーは、特に次のパニックを捕捉するよう設計されています。

- フロントエンドから呼び出された、バインド済みサービスメソッド内のパニック
- Wails ランタイム内部のパニック

バックグラウンドの goroutine や独立した Go コードなど、その他の状況では、Go の標準的なパニック回復機構を使用して自分でパニックを処理してください。

## PanicDetails 構造体

パニックが発生すると、Wails はパニックに関する重要な情報を`PanicDetails`構造体に記録します。

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
フロントエンドから呼び出された Service コード内でパニックが捕捉されると、コード内でパニックが発生した正確な位置に焦点を当てるため、スタックトレースが短縮されます。 完全なスタックトレースを確認するには、`FullStackTrace`フィールドを使用できます。

@end

## デフォルトのパニックハンドラー

カスタムパニックハンドラーを指定しない場合、Wails はデフォルトのハンドラーを使用します。このハンドラーは、整形したログメッセージとしてエラー情報を出力した後、終了します。 例：

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

アプリケーションの作成時に`PanicHandler`オプションを設定することで、独自のパニックハンドラーを実装できます。例を次に示します。

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

## 自分で起動したゴルーチンのパニック処理 {#user-goroutines}

Wails が遅延実行によるリカバリーを設定できるのは、自身が制御する箇所だけです。これにはフロントエンドから呼ばれるバインド済みサービスメソッドや内部ランタイムのコールバックなどが含まれます。自分のコードで次のように書いた場合を考えます。

```go
go func() {
    // your work
}()
```

Wails はそのゴルーチンに `defer handlePanic()` を挿入できません。そこでパニックが発生すると、Go の既定の動作に従ってプロセス全体がクラッシュし、登録した `PanicHandler` は**呼び出されません**。

自分のゴルーチンのパニックを Wails が捕捉したパニックと同じハンドラーへ送るには、`PanicDetails` を手動で構築し、`PanicHandler` として登録した関数を呼び出す小さなヘルパーを追加します。

```go
import (
    "fmt"
    "runtime/debug"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

// reportPanic is the function you register as application.Options.PanicHandler.
func reportPanic(pd *application.PanicDetails) {
    // log to file / send to Sentry / show dialog / etc.
}

// recoverAndReport funnels goroutine panics to reportPanic. Defer it as the
// first statement of every goroutine you spawn in user code.
func recoverAndReport() {
    r := recover()
    if r == nil {
        return
    }
    err, ok := r.(error)
    if !ok {
        err = fmt.Errorf("%v", r)
    }
    stack := string(debug.Stack())
    reportPanic(&application.PanicDetails{
        Error:          err,
        Time:           time.Now(),
        StackTrace:     stack,
        FullStackTrace: stack,
    })
}
```

自分で起動する各ゴルーチンの先頭で使用してください。

```go
go func() {
    defer recoverAndReport()
    // your work
}()
```

Wails による捕捉経路と自分のゴルーチンの経路はどちらも `reportPanic` に到達し、報告処理を一元化できます。

Wails には 2 つの実行可能なサンプルがあります。

- `v3/examples/panic-handling` — 最小構成です。バインド済みメソッドのパニックだけを `PanicHandler` に送ります。
- `v3/examples/user-panic-handling` — バインド済みメソッドとバックグラウンドゴルーチンのパニックを同じハンドラーに送ります。

@note{type="caution" title="スタックトレースの忠実性"}

`PanicDetails.StackTrace` は Wails 自身のラッパーフレームを隠し、先頭に自分のコードが表示されるように短縮されます。`runtime/debug.Stack()` から自分で `PanicDetails` を構築した場合、この短縮は行われません。`StackTrace` と `FullStackTrace` は同一になり、ゴルーチンのスタック全体が含まれます。自分で起動したゴルーチンでは、通常これが望ましい動作です。

@end

## パニック時の診断情報の取得 {#panic-diagnostics}

`PanicHandler` はパニックそのものを受け取りますが、後から状況を再現するために、システム情報、ビルド情報、プロセス・メモリ・モジュールの状態、Windows ではミニダンプを含む、プロセス全体の状態も取得したい場合があります。

Wails はそのための専用パッケージを提供しています。[クラッシュのデバッグ](/guides/debugging-crashes/)を参照してください。一般的な統合例は次のとおりです。

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        report, _ := debug.Report(debug.WithDump())
        // pd holds the wails-side panic info; report adds system context
        // and (on Windows) a minidump at report.DumpPath.
        mycrashservice.Upload(pd, report)
    },
})
```

`debug.Report` は明示的に選択して使用します。クラッシュ報告にはユーザーの同意や個人情報の除去が必要なことが多いため、Wails がダンプやシステムスナップショットを自動保存することはありません。API 全体については[クラッシュのデバッグ](/guides/debugging-crashes/)を参照してください。

## 最後に

Wails のパニックハンドラーは、バインド済みメソッド内のパニックとランタイム内部のエラーを処理するためのものであることに注意してください。アプリケーションのその他の部分では、必要に応じて Go の標準的なエラー処理パターンとパニック回復機構を使用してください。すべての Go アプリケーションと同様に、可能な限り適切なエラー処理によってパニックを防ぐことが望まれます。
