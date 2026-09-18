---
title: "Gestion des paniques"
description: "Comment gérer les paniques dans votre application Wails"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Dans les applications Go, des paniques peuvent survenir à l’exécution lorsqu’un événement inattendu se produit. Ce guide explique comment gérer les paniques dans le code Go en général, puis plus particulièrement dans votre application Wails.

## Comprendre les paniques dans Go

Avant d’aborder la gestion des paniques propre à Wails, il est essentiel de comprendre leur fonctionnement dans Go :

1. Les paniques sont réservées aux erreurs irrécupérables qui ne devraient pas se produire en fonctionnement normal
2. Lorsqu’une panique survient dans une goroutine, seule cette goroutine est affectée
3. Il est possible de récupérer après une panique à l’aide de `defer` et de `recover()`

Voici un exemple simple de gestion des paniques dans Go :

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

Pour en savoir plus sur panic et recover dans Go, consultez l’article [Blog Go : Defer, Panic et Recover](https://go.dev/blog/defer-panic-and-recover).

## Gestion des paniques dans Wails

Wails gère automatiquement les paniques qui surviennent dans vos méthodes de service lorsqu’elles sont appelées depuis le frontend. Vous n’avez donc pas besoin d’ajouter de mécanisme de récupération après une panique à ces méthodes : Wails intercepte la panique et la transmet au gestionnaire de paniques configuré.

Le gestionnaire de paniques est spécialement conçu pour intercepter :

- Les paniques dans les méthodes de service liées appelées depuis le frontend
- Les paniques internes provenant de l’environnement d’exécution de Wails

Dans les autres cas, par exemple dans des goroutines en arrière-plan ou du code Go autonome, vous devriez gérer vous-même les paniques à l’aide des mécanismes standard de récupération après une panique de Go.

## La structure PanicDetails

Lorsqu’une panique survient, Wails enregistre les informations importantes la concernant dans une structure `PanicDetails` :

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Cette structure fournit des informations détaillées sur la panique :

- `StackTrace` : chaîne mise en forme indiquant la pile d’appels ayant conduit à la panique
- `Error` : erreur ou message de panique proprement dit
- `Time` : heure exacte à laquelle la panique est survenue
- `FullStackTrace` : trace complète de la pile, y compris les cadres de pile de l’environnement d’exécution Go

@note{type="info" title="Paniques dans le code de service"}
Lorsque des paniques sont interceptées dans votre code de service après un appel depuis le frontend, la trace de la pile est raccourcie afin d’indiquer précisément où la panique est survenue dans votre code. Pour afficher la trace complète de la pile, vous pouvez utiliser le champ `FullStackTrace`.

@end

## Gestionnaire de paniques par défaut

Si vous ne spécifiez pas de gestionnaire de paniques personnalisé, Wails utilise son gestionnaire par défaut, qui affiche les informations sur l’erreur dans un message de journal mis en forme, puis quitte l’application. Par exemple :

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

## Gestionnaire de paniques personnalisé

Vous pouvez implémenter votre propre gestionnaire de paniques en définissant l’option `PanicHandler` lors de la création de votre application. Voici un exemple :

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

## Gérer les paniques de vos propres goroutines {#user-goroutines}

Wails ne peut installer sa récupération différée qu’aux endroits qu’il contrôle : méthodes liées appelées depuis le frontend, fonctions de rappel internes du runtime, etc. Si votre propre code fait ceci :

```go
go func() {
    // your work
}()
```

Wails ne peut pas injecter un `defer handlePanic()` dans cette goroutine. Si elle panique, tout le processus plante conformément au comportement par défaut de Go : votre `PanicHandler` enregistré **n’est pas** appelé.

Pour faire passer les paniques de vos goroutines par le même gestionnaire que celles interceptées par Wails, ajoutez une petite fonction qui construit manuellement un `PanicDetails` et appelle la fonction que vous avez enregistrée comme `PanicHandler` :

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

Utilisez-la au début de chaque goroutine dont vous êtes responsable :

```go
go func() {
    defer recoverAndReport()
    // your work
}()
```

Les paniques interceptées par Wails et celles de vos goroutines aboutissent maintenant à `reportPanic`, ce qui centralise le signalement.

Wails fournit deux exemples exécutables :

- `v3/examples/panic-handling` — exemple minimal : uniquement une panique dans une méthode liée, transmise à `PanicHandler`.
- `v3/examples/user-panic-handling` — une panique dans une méthode liée et une panique dans une goroutine d’arrière-plan passent toutes deux par le même gestionnaire.

@note{type="caution" title="Fidélité des traces de pile"}

`PanicDetails.StackTrace` est raccourci par Wails pour masquer ses propres fonctions d’enveloppe et placer votre code en haut de la trace. Lorsque vous construisez vous-même un `PanicDetails` à partir de `runtime/debug.Stack()`, aucun raccourcissement n’est appliqué : `StackTrace` et `FullStackTrace` sont identiques et contiennent toute la pile de votre goroutine. C’est généralement ce que vous souhaitez pour les goroutines que vous créez.

@end

## Capturer des diagnostics lors d’une panique {#panic-diagnostics}

Le `PanicHandler` reçoit la panique elle-même, mais vous souhaiterez souvent capturer l’état complet du processus : informations système et de build, état du processus, de la mémoire et des modules, ainsi qu’un minidump sous Windows, pour reconstituer les événements plus tard.

Wails fournit un package dédié à cela : [Déboguer les plantages](/guides/debugging-crashes/). Intégration habituelle :

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

`debug.Report` est facultatif : Wails n’enregistre jamais automatiquement un dump ou un instantané système, car le signalement de plantages implique souvent le consentement de l’utilisateur ou la suppression de données personnelles. Consultez [Déboguer les plantages](/guides/debugging-crashes/) pour l’API complète.

## Remarques finales

N’oubliez pas que le gestionnaire de paniques de Wails sert spécifiquement à gérer les paniques dans les méthodes liées et les erreurs internes de l’environnement d’exécution. Dans les autres parties de votre application, utilisez les modèles standard de gestion des erreurs et, lorsque cela est approprié, les mécanismes de récupération après une panique de Go. Comme pour toute application Go, il est préférable, dans la mesure du possible, de prévenir les paniques grâce à une gestion adéquate des erreurs.
