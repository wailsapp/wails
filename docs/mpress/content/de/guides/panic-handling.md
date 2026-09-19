---
title: "Panics behandeln"
description: "So behandeln Sie Panics in Ihrer Wails-Anwendung"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

In Go-Anwendungen können zur Laufzeit Panics auftreten, wenn etwas Unerwartetes geschieht. Dieser Leitfaden erläutert, wie Sie Panics sowohl in allgemeinem Go-Code als auch speziell in Ihrer Wails-Anwendung behandeln.

## Panics in Go verstehen

Bevor Sie sich mit der Wails-spezifischen Panic-Behandlung befassen, müssen Sie verstehen, wie Panics in Go funktionieren:

1. Panics sind für nicht behebbare Fehler vorgesehen, die im normalen Betrieb nicht auftreten sollten
2. Wenn in einer Goroutine eine Panic auftritt, ist nur diese Goroutine betroffen
3. Panics können mit `defer` und `recover()` abgefangen werden

Hier sehen Sie ein einfaches Beispiel für die Panic-Behandlung in Go:

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

Ausführlichere Informationen zu panic und recover in Go finden Sie im [Go-Blogbeitrag „Defer, Panic, and Recover“](https://go.dev/blog/defer-panic-and-recover).

## Panic-Behandlung in Wails

Wails behandelt automatisch Panics, die in Ihren Service-Methoden auftreten, wenn diese vom Frontend aufgerufen werden. Sie müssen diesen Methoden daher keine Panic-Wiederherstellung hinzufügen: Wails fängt die Panic ab und verarbeitet sie mit Ihrem konfigurierten Panic-Handler.

Der Panic-Handler ist speziell dafür ausgelegt, Folgendes abzufangen:

- Panics in gebundenen Service-Methoden, die vom Frontend aufgerufen werden
- Interne Panics der Wails-Laufzeitumgebung

In anderen Szenarien, etwa bei Hintergrund-Goroutinen oder eigenständigem Go-Code, sollten Sie Panics selbst mit den standardmäßigen Mechanismen zur Panic-Wiederherstellung von Go behandeln.

## Die Struktur PanicDetails

Wenn eine Panic auftritt, erfasst Wails wichtige Informationen dazu in einer `PanicDetails`-Struktur:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Diese Struktur stellt umfassende Informationen zur Panic bereit:

- `StackTrace`: Eine formatierte Zeichenfolge mit dem Aufruf-Stack, der zur Panic geführt hat
- `Error`: Die eigentliche Fehler- oder Panic-Meldung
- `Time`: Der genaue Zeitpunkt, zu dem die Panic aufgetreten ist
- `FullStackTrace`: Der vollständige Stacktrace einschließlich Laufzeit-Frames

@note{type="info" title="Panics im Service-Code"}
Wenn Panics in Ihrem Service-Code nach einem Aufruf vom Frontend abgefangen werden, wird der Stacktrace so gekürzt, dass er genau die Stelle in Ihrem Code hervorhebt, an der die Panic aufgetreten ist. Wenn Sie den vollständigen Stacktrace sehen möchten, können Sie das Feld `FullStackTrace` verwenden.

@end

## Standardmäßiger Panic-Handler

Wenn Sie keinen benutzerdefinierten Panic-Handler angeben, verwendet Wails seinen Standard-Handler. Dieser gibt Fehlerinformationen als formatierte Protokollmeldung aus und beendet anschließend die Anwendung. Beispiel:

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

Sie können einen eigenen Panic-Handler implementieren, indem Sie beim Erstellen Ihrer Anwendung die Option `PanicHandler` festlegen. Beispiel:

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

## Panics in eigenen Goroutinen behandeln {#user-goroutines}

Wails kann seine verzögerte Wiederherstellung nur an Stellen einrichten, die es kontrolliert: gebundene, vom Frontend aufgerufene Servicemethoden, interne Laufzeit-Callbacks und Ähnliches. Wenn dein eigener Code Folgendes tut:

```go
go func() {
    // your work
}()
```

Wails kann kein `defer handlePanic()` in diese Goroutine einfügen. Bei einer Panic stürzt gemäß dem Standardverhalten von Go der gesamte Prozess ab. Dein registrierter `PanicHandler` wird **nicht** aufgerufen.

Um Panics eigener Goroutinen durch denselben Handler wie von Wails abgefangene Panics zu leiten, ergänze eine kleine Hilfsfunktion, die manuell ein `PanicDetails` erstellt und die als `PanicHandler` registrierte Funktion aufruft:

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

Verwende sie am Anfang jeder von dir gestarteten Goroutine:

```go
go func() {
    defer recoverAndReport()
    // your work
}()
```

Sowohl von Wails abgefangene Panics als auch Panics eigener Goroutinen gelangen nun zu `reportPanic`; die Berichterstattung bleibt zentralisiert.

Wails enthält zwei ausführbare Beispiele:

- `v3/examples/panic-handling` — minimal: nur eine Panic in einer gebundenen Methode, weitergeleitet an `PanicHandler`.
- `v3/examples/user-panic-handling` — eine Panic in einer gebundenen Methode und eine in einer Hintergrund-Goroutine gelangen beide zum selben Handler.

@note{type="caution" title="Genauigkeit des Stacktraces"}

`PanicDetails.StackTrace` wird von Wails gekürzt, um eigene Wrapper-Frames auszublenden und deinen Code am Anfang anzuzeigen. Wenn du ein `PanicDetails` selbst aus `runtime/debug.Stack()` erstellst, findet keine Kürzung statt: `StackTrace` und `FullStackTrace` sind identisch und enthalten den vollständigen Stack deiner Goroutine. Für selbst gestartete Goroutinen ist das normalerweise erwünscht.

@end

## Diagnosedaten bei einer Panic erfassen {#panic-diagnostics}

Der `PanicHandler` erhält die Panic selbst. Häufig möchtest du jedoch den vollständigen Prozesszustand erfassen: System- und Build-Informationen, Prozess-, Speicher- und Modulzustand sowie unter Windows einen Minidump, um den Ablauf später rekonstruieren zu können.

Wails bietet dafür ein eigenes Paket: [Abstürze debuggen](/guides/debugging-crashes/). Typische Integration:

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

`debug.Report` muss ausdrücklich aufgerufen werden. Wails speichert niemals automatisch einen Dump oder eine Systemaufnahme, da Absturzberichte oft die Zustimmung des Benutzers oder das Entfernen personenbezogener Daten erfordern. Die vollständige API beschreibt [Abstürze debuggen](/guides/debugging-crashes/).

## Abschließende Hinweise

Beachten Sie, dass der Wails-Panic-Handler speziell für Panics in gebundenen Methoden und für interne Laufzeitfehler vorgesehen ist. In anderen Teilen Ihrer Anwendung sollten Sie gegebenenfalls die standardmäßigen Fehlerbehandlungsmuster und Mechanismen zur Panic-Wiederherstellung von Go verwenden. Wie bei allen Go-Anwendungen ist es nach Möglichkeit besser, Panics durch eine ordnungsgemäße Fehlerbehandlung zu verhindern.
