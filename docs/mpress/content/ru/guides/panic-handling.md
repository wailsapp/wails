---
title: "Обработка паник"
description: "Как обрабатывать паники в приложении Wails"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

В приложениях на Go паники могут возникать во время выполнения, когда происходит что-то непредвиденное. В этом руководстве объясняется, как обрабатывать паники как в коде Go в целом, так и непосредственно в приложении Wails.

## Основные сведения о паниках в Go

Прежде чем переходить к обработке паник в Wails, важно понять, как паники работают в Go:

1. Паники предназначены для неисправимых ошибок, которые не должны возникать при нормальной работе
2. Когда в горутине возникает паника, она затрагивает только эту горутину
3. После паники выполнение можно восстановить с помощью `defer` и `recover()`

Ниже приведён простой пример обработки паники в Go:

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

Подробные сведения о panic и recover в Go см. в статье [«Блог Go: Defer, Panic и Recover»](https://go.dev/blog/defer-panic-and-recover).

## Обработка паник в Wails

Wails автоматически обрабатывает паники, возникающие в методах Service при их вызове из фронтенда. Поэтому добавлять в эти методы восстановление после паники не нужно — Wails перехватит панику и передаст её настроенному обработчику паник.

Обработчик паник предназначен непосредственно для перехвата следующих паник:

- Паники в привязанных методах сервисов, вызванных из фронтенда
- Внутренние паники среды выполнения Wails

В других сценариях, например в фоновых горутинах или автономном коде Go, следует самостоятельно обрабатывать паники с помощью стандартных механизмов восстановления после паники в Go.

## Структура PanicDetails

При возникновении паники Wails сохраняет важную информацию о ней в структуре `PanicDetails`:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Эта структура предоставляет исчерпывающую информацию о панике:

- `StackTrace`: форматированная строка со стеком вызовов, приведших к панике
- `Error`: фактическая ошибка или сообщение о панике
- `Time`: точное время возникновения паники
- `FullStackTrace`: полная трассировка стека, включая кадры среды выполнения

@note{type="info" title="Паники в коде Service"}
Когда паника перехватывается в коде Service после вызова из фронтенда, трассировка стека сокращается, чтобы точно показать место в вашем коде, где возникла паника. Чтобы просмотреть полную трассировку стека, используйте поле `FullStackTrace`.

@end

## Обработчик паник по умолчанию

Если пользовательский обработчик паник не указан, Wails использует обработчик по умолчанию, который выводит сведения об ошибке в виде форматированного сообщения журнала, а затем завершает работу. Например:

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

## Пользовательский обработчик паник

Чтобы реализовать собственный обработчик паник, задайте параметр `PanicHandler` при создании приложения. Ниже приведён пример:

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

## Обработка паник в собственных горутинах {#user-goroutines}

Wails может установить отложенное восстановление только там, где контролирует выполнение: в связанных методах сервисов, вызываемых из фронтенда, внутренних обработчиках среды выполнения и подобных местах. Если ваш код делает следующее:

```go
go func() {
    // your work
}()
```

Wails не может внедрить `defer handlePanic()` в эту горутину. При её панике весь процесс аварийно завершается согласно стандартному поведению Go; зарегистрированный `PanicHandler` **не** вызывается.

Чтобы направить паники собственных горутин в тот же обработчик, что и паники, перехваченные Wails, добавьте небольшую функцию, которая вручную создаёт `PanicDetails` и вызывает функцию, зарегистрированную как `PanicHandler`:

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

Используйте её в начале каждой горутины, которую запускаете сами:

```go
go func() {
    defer recoverAndReport()
    // your work
}()
```

Теперь и путь перехвата Wails, и путь пользовательской горутины приводят к `reportPanic`, сохраняя централизованную обработку отчётов.

С Wails поставляются два запускаемых примера:

- `v3/examples/panic-handling` — минимальный: только паника связанного метода, направляемая в `PanicHandler`.
- `v3/examples/user-panic-handling` — паника связанного метода и паника фоновой горутины проходят через один обработчик.

@note{type="caution" title="Точность трассировки стека"}

`PanicDetails.StackTrace` сокращается Wails, чтобы скрыть собственные кадры обёрток и показать ваш код в начале трассировки. Если вы создаёте `PanicDetails` самостоятельно из `runtime/debug.Stack()`, сокращения нет: `StackTrace` и `FullStackTrace` будут одинаковыми и будут содержать полный стек вашей горутины. Для самостоятельно запущенных горутин обычно это и требуется.

@end

## Сбор диагностики при панике {#panic-diagnostics}

Обработчик `PanicHandler` получает саму панику, но часто нужно сохранить полное состояние процесса: сведения о системе и сборке, состояние процесса, памяти и модулей, а в Windows ещё и минидамп, чтобы позднее восстановить произошедшее.

Для этого Wails предоставляет отдельный пакет: [Отладка сбоев](/guides/debugging-crashes/). Типичная интеграция:

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

`debug.Report` вызывается только по явному выбору: Wails никогда автоматически не сохраняет дамп или системный снимок, поскольку отчёты о сбоях часто требуют согласия пользователя или удаления персональных данных. Полный API описан в разделе [Отладка сбоев](/guides/debugging-crashes/).

## Заключительные замечания

Помните, что обработчик паник Wails предназначен непосредственно для обработки паник в привязанных методах и внутренних ошибок среды выполнения. В других частях приложения следует по мере необходимости использовать стандартные шаблоны обработки ошибок и механизмы восстановления после паники в Go. Как и в любых приложениях Go, по возможности лучше предотвращать паники посредством правильной обработки ошибок.
