---
title: "Обработка паник"
description: "Как обрабатывать паники в приложении Wails"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

В приложениях на Go паника может возникнуть во время выполнения, если происходит что-то непредвиденное. В этом руководстве объясняется, как обрабатывать паники как в обычном коде Go, так и непосредственно в приложении Wails.

## Общие сведения о паниках в Go

Прежде чем переходить к обработке паник в Wails, важно понять, как они работают в Go:

1. Паники предназначены для неисправимых ошибок, которые не должны возникать при нормальной работе
2. Паника разворачивает стек текущей горутины; если она достигает вершины стека и не была перехвачена, процесс завершается вместе с остальными горутинами
3. Восстановление должно выполняться в отложенной функции в той же горутине, в которой возникла паника

Ниже приведён базовый пример обработки паники в Go:

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

Подробнее о panic и recover в Go см. в статье [блога Go «Defer, Panic, and Recover»](https://go.dev/blog/defer-panic-and-recover).

## Обработка паник в Wails

Wails автоматически обрабатывает паники, возникающие в методах Service при их вызове из фронтенда. Поэтому добавлять восстановление после паники в эти методы не требуется: Wails перехватит панику и передаст её настроенному обработчику паник.

Обработчик паник предназначен специально для перехвата:

- Паник в привязанных методах сервисов, вызванных из фронтенда
- Паник в путях выполнения среды исполнения, где явно установлено восстановление Wails; это не глобальный обработчик всех горутин процесса

В других сценариях, например в фоновых горутинах или автономном коде Go, следует самостоятельно обрабатывать паники с помощью стандартных механизмов восстановления Go.

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

Эта структура содержит исчерпывающую информацию о панике:

- `StackTrace`: форматированная строка со стеком вызовов, приведших к панике
- `Error`: фактическая ошибка или сообщение паники
- `Time`: точное время возникновения паники
- `FullStackTrace`: полная трассировка стека, включая кадры среды исполнения

@note{type="info" title="Паники в коде Service"}
Когда паника перехватывается в коде Service после его вызова из фронтенда, трассировка стека сокращается, чтобы точно показать место в вашем коде, где возникла паника. Чтобы просмотреть полную трассировку стека, используйте поле `FullStackTrace`.

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

Помните, что обработчик паник Wails предназначен специально для обработки паник в привязанных методах и внутренних ошибок среды исполнения. В остальных частях приложения следует при необходимости использовать стандартные шаблоны обработки ошибок и механизмы восстановления после паник Go. Как и во всех приложениях Go, по возможности лучше предотвращать паники с помощью надлежащей обработки ошибок.
