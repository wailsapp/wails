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

Полный рабочий пример обработки паник в приложении Wails см. в примере panic-handling в `v3/examples/panic-handling`.

## Заключительные замечания

Помните, что обработчик паник Wails предназначен непосредственно для обработки паник в привязанных методах и внутренних ошибок среды выполнения. В других частях приложения следует по мере необходимости использовать стандартные шаблоны обработки ошибок и механизмы восстановления после паники в Go. Как и в любых приложениях Go, по возможности лучше предотвращать паники посредством правильной обработки ошибок.
