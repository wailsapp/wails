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

Pour consulter un exemple fonctionnel complet de gestion des paniques dans une application Wails, reportez-vous à l’exemple panic-handling dans `v3/examples/panic-handling`.

## Remarques finales

N’oubliez pas que le gestionnaire de paniques de Wails sert spécifiquement à gérer les paniques dans les méthodes liées et les erreurs internes de l’environnement d’exécution. Dans les autres parties de votre application, utilisez les modèles standard de gestion des erreurs et, lorsque cela est approprié, les mécanismes de récupération après une panique de Go. Comme pour toute application Go, il est préférable, dans la mesure du possible, de prévenir les paniques grâce à une gestion adéquate des erreurs.
