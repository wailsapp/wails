---
title: "Gestion des paniques"
description: "Comment gérer les paniques dans votre application Wails"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Dans les applications Go, des paniques peuvent survenir lors de l’exécution lorsqu’un événement inattendu se produit. Ce guide explique comment gérer les paniques aussi bien dans le code Go en général que dans votre application Wails en particulier.

## Comprendre les paniques en Go

Avant d’aborder la gestion des paniques propre à Wails, il est essentiel de comprendre leur fonctionnement en Go :

1. Les paniques sont destinées aux erreurs irrécupérables qui ne devraient pas se produire en fonctionnement normal
2. Une panique déroule la pile de la goroutine courante ; si elle atteint le sommet sans récupération, le processus se termine, ainsi que les autres goroutines
3. La récupération doit s’exécuter dans une fonction différée au sein de la même goroutine que celle où la panique s’est produite

Voici un exemple simple de gestion des paniques en Go :

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

Pour en savoir plus sur panic et recover en Go, consultez l’article du [blog Go : Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover).

## Gestion des paniques dans Wails

Wails gère automatiquement les paniques qui surviennent dans vos méthodes de service lorsqu’elles sont appelées depuis le frontend. Vous n’avez donc pas besoin d’ajouter de mécanisme de récupération des paniques à ces méthodes : Wails intercepte la panique et la traite à l’aide du gestionnaire de paniques configuré.

Le gestionnaire de paniques est spécialement conçu pour intercepter :

- Les paniques dans les méthodes de service liées appelées depuis le frontend
- Les paniques dans les chemins d’exécution qui installent explicitement le mécanisme de récupération de Wails ; il ne s’agit pas d’un gestionnaire à l’échelle du processus pour toutes les goroutines

Dans les autres cas, par exemple pour les goroutines en arrière-plan ou le code Go autonome, vous devez gérer vous-même les paniques à l’aide des mécanismes standard de récupération après une panique de Go.

## La structure PanicDetails

Lorsqu’une panique survient, Wails enregistre des informations importantes à son sujet dans une structure `PanicDetails` :

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Cette structure fournit des informations détaillées sur la panique :

- `StackTrace` : chaîne mise en forme qui affiche la pile d’appels ayant conduit à la panique
- `Error` : erreur réelle ou message de panique
- `Time` : heure exacte à laquelle la panique s’est produite
- `FullStackTrace` : trace de pile complète, y compris les frames du runtime

@note{type="info" title="Paniques dans le code des services"}
Lorsque des paniques sont interceptées dans votre code de service après un appel depuis le frontend, la trace de pile est tronquée afin de cibler précisément l’endroit de votre code où la panique s’est produite. Pour afficher la trace de pile complète, utilisez le champ `FullStackTrace`.

@end

## Gestionnaire de paniques par défaut

Si vous ne spécifiez pas de gestionnaire de paniques personnalisé, Wails utilise son gestionnaire par défaut, qui affiche les informations sur l’erreur dans un message de journal mis en forme, puis ferme l’application. Par exemple :

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

Gardez à l’esprit que le gestionnaire de paniques de Wails est spécifiquement destiné à gérer les paniques dans les méthodes liées et les erreurs internes du runtime. Pour les autres parties de votre application, utilisez les modèles standard de gestion des erreurs de Go et, lorsque cela est approprié, ses mécanismes de récupération après une panique. Comme pour toute application Go, il est préférable d’empêcher les paniques grâce à une gestion correcte des erreurs chaque fois que cela est possible.
