---
title: "Panics behandeln"
description: "So behandeln Sie Panics in Ihrer Wails-Anwendung"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

In Go-Anwendungen können zur Laufzeit Panics auftreten, wenn etwas Unerwartetes geschieht. Dieser Leitfaden erläutert, wie Sie Panics sowohl in allgemeinem Go-Code als auch speziell in Ihrer Wails-Anwendung behandeln.

## Panics in Go verstehen

Bevor Sie sich mit der Wails-spezifischen Panic-Behandlung befassen, sollten Sie verstehen, wie Panics in Go funktionieren:

1. Panics sind für nicht behebbare Fehler vorgesehen, die im Normalbetrieb nicht auftreten sollten
2. Eine Panic wickelt die aktuelle Goroutine ab. Erreicht sie ohne Wiederherstellung die oberste Ebene, wird der Prozess einschließlich der anderen Goroutinen beendet
3. Die Wiederherstellung muss in einer verzögert ausgeführten Funktion in derselben Goroutine wie die Panic erfolgen

Hier ist ein einfaches Beispiel für die Behandlung von Panics in Go:

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

Ausführlichere Informationen zu panic und recover in Go finden Sie im [Go-Blog: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover).

## Panic-Behandlung in Wails

Wails behandelt automatisch Panics, die in Ihren Service-Methoden auftreten, wenn diese vom Frontend aufgerufen werden. Sie müssen diesen Methoden daher keine Panic-Wiederherstellung hinzufügen. Wails fängt die Panic ab und verarbeitet sie mit dem von Ihnen konfigurierten Panic-Handler.

Der Panic-Handler ist speziell dafür ausgelegt, Folgendes abzufangen:

- Panics in gebundenen Service-Methoden, die vom Frontend aufgerufen werden
- Panics in Laufzeitpfaden, die ausdrücklich die Wails-Wiederherstellung einrichten; dies ist kein prozessweiter Handler für jede Goroutine

In anderen Szenarien, etwa bei Hintergrund-Goroutinen oder eigenständigem Go-Code, sollten Sie Panics selbst mit den standardmäßigen Mechanismen von Go zur Panic-Wiederherstellung behandeln.

## Die PanicDetails-Struktur

Wenn eine Panic auftritt, erfasst Wails wichtige Informationen dazu in einer `PanicDetails`-Struktur:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Diese Struktur stellt umfassende Informationen über die Panic bereit:

- `StackTrace`: Eine formatierte Zeichenfolge mit dem Aufrufstapel, der zur Panic geführt hat
- `Error`: Die eigentliche Fehler- oder Panic-Meldung
- `Time`: Der genaue Zeitpunkt, zu dem die Panic aufgetreten ist
- `FullStackTrace`: Der vollständige Stacktrace einschließlich der Laufzeit-Frames

@note{type="info" title="Panics im Service-Code"}
Wenn Panics nach einem Aufruf aus dem Frontend in Ihrem Service-Code abgefangen werden, wird der Stacktrace so gekürzt, dass genau die Stelle in Ihrem Code im Mittelpunkt steht, an der die Panic aufgetreten ist. Den vollständigen Stacktrace können Sie über das Feld `FullStackTrace` anzeigen.

@end

## Standardmäßiger Panic-Handler

Wenn Sie keinen benutzerdefinierten Panic-Handler angeben, verwendet Wails seinen Standard-Handler. Dieser gibt Fehlerinformationen in einer formatierten Protokollmeldung aus und beendet anschließend die Anwendung. Beispiel:

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

## Benutzerdefinierter Panic-Handler

Sie können einen eigenen Panic-Handler implementieren, indem Sie beim Erstellen Ihrer Anwendung die Option `PanicHandler` festlegen. Hier ist ein Beispiel:

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

Ein vollständiges, funktionsfähiges Beispiel für die Panic-Behandlung in einer Wails-Anwendung finden Sie im Beispiel zur Panic-Behandlung unter `v3/examples/panic-handling`.

## Abschließende Hinweise

Beachten Sie, dass der Wails-Panic-Handler speziell für die Behandlung von Panics in gebundenen Methoden und von internen Laufzeitfehlern vorgesehen ist. In anderen Teilen Ihrer Anwendung sollten Sie gegebenenfalls die standardmäßigen Fehlerbehandlungsmuster und Mechanismen zur Panic-Wiederherstellung von Go verwenden. Wie bei allen Go-Anwendungen ist es nach Möglichkeit besser, Panics durch eine ordnungsgemäße Fehlerbehandlung zu verhindern.
