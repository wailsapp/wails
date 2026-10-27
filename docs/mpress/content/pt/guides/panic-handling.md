---
title: "Tratamento de panics"
description: "Como tratar panics em sua aplicação Wails"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Em aplicações Go, podem ocorrer panics durante a execução quando algo inesperado acontece. Este guia explica como tratar panics tanto no código Go em geral quanto especificamente em sua aplicação Wails.

## Entendendo panics em Go

Antes de abordar o tratamento de panics específico do Wails, é essencial entender como os panics funcionam em Go:

1. Panics são destinados a erros irrecuperáveis que não deveriam ocorrer durante a operação normal
2. Um panic desfaz a pilha da goroutine atual; se chegar ao topo sem recuperação, o processo será encerrado, incluindo as outras goroutines
3. A recuperação deve ser executada em uma função adiada na mesma goroutine em que ocorreu o panic

Veja um exemplo básico de tratamento de panic em Go:

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

Para obter informações mais detalhadas sobre panic e recover em Go, consulte [Blog do Go: Defer, Panic e Recover](https://go.dev/blog/defer-panic-and-recover).

## Tratamento de panics no Wails

O Wails trata automaticamente os panics que ocorrem nos métodos dos seus Services quando eles são chamados pelo frontend. Isso significa que você não precisa adicionar recuperação de panic a esses métodos — o Wails capturará o panic e o processará por meio do manipulador de panics configurado.

O manipulador de panics foi projetado especificamente para capturar:

- Panics em métodos de serviço vinculados chamados pelo frontend
- Panics em caminhos de execução que instalam explicitamente a recuperação do Wails; esse não é um manipulador para todas as goroutines do processo

Em outros cenários, como goroutines em segundo plano ou código Go independente, você deve tratar os panics por conta própria usando os mecanismos padrão de recuperação de panics do Go.

## A struct PanicDetails

Quando ocorre um panic, o Wails captura informações importantes sobre ele em uma struct `PanicDetails`:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Essa estrutura fornece informações abrangentes sobre o panic:

- `StackTrace`: uma string formatada que mostra a pilha de chamadas que levou ao panic
- `Error`: o erro ou a mensagem de panic propriamente dita
- `Time`: o momento exato em que ocorreu o panic
- `FullStackTrace`: o stack trace completo, incluindo os frames do runtime

@note{type="info" title="Panics no código do Service"}
Quando panics no código do seu Service são capturados após uma chamada feita pelo frontend, o stack trace é reduzido para destacar o ponto exato do seu código em que ocorreu o panic. Para ver o stack trace completo, use o campo `FullStackTrace`.

@end

## Manipulador de panics padrão

Se você não especificar um manipulador de panics personalizado, o Wails usará o manipulador padrão, que exibe as informações do erro em uma mensagem de log formatada e encerra a aplicação em seguida. Por exemplo:

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

## Manipulador de panics personalizado

Você pode implementar seu próprio manipulador de panics definindo a opção `PanicHandler` ao criar sua aplicação. Veja um exemplo:

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

Para ver um exemplo funcional completo de tratamento de panics em uma aplicação Wails, consulte o exemplo de tratamento de panics em `v3/examples/panic-handling`.

## Observações finais

Lembre-se de que o manipulador de panics do Wails se destina especificamente a gerenciar panics em métodos vinculados e erros internos do runtime. Para as outras partes da aplicação, use os padrões de tratamento de erros e os mecanismos de recuperação de panics padrão do Go quando apropriado. Como em todas as aplicações Go, é melhor evitar panics por meio do tratamento adequado de erros sempre que possível.
