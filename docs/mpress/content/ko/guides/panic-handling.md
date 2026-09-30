---
title: "패닉 처리"
description: "Wails 애플리케이션에서 패닉을 처리하는 방법"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Go 애플리케이션에서는 런타임 중 예기치 않은 상황이 발생하면 패닉이 일어날 수 있습니다. 이 가이드에서는 일반적인 Go 코드와 Wails 애플리케이션에서 패닉을 처리하는 방법을 설명합니다.

## Go의 패닉 이해하기

Wails에 특화된 패닉 처리 방법을 살펴보기 전에 Go에서 패닉이 어떻게 작동하는지 이해해야 합니다:

1. 패닉은 정상적인 작동 중에는 발생해서는 안 되는 복구 불가능한 오류에 사용됩니다
2. goroutine에서 패닉이 발생하면 해당 goroutine만 영향을 받습니다
3. `defer`와 `recover()`를 사용해 패닉에서 복구할 수 있습니다

다음은 Go에서 패닉을 처리하는 기본 예제입니다:

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

Go의 panic과 recover에 관한 자세한 내용은 [Go 블로그: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)를 참조하세요.

## Wails의 패닉 처리

Wails는 프런트엔드에서 호출된 Service 메서드에서 발생하는 패닉을 자동으로 처리합니다. 따라서 이러한 메서드에 패닉 복구 코드를 추가할 필요가 없습니다. Wails가 패닉을 포착하여 구성된 패닉 핸들러로 처리합니다.

패닉 핸들러는 다음 패닉을 포착하도록 특별히 설계되었습니다:

- 프런트엔드에서 호출된 바인딩된 서비스 메서드의 패닉
- Wails 런타임에서 발생하는 내부 패닉

백그라운드 goroutine이나 독립 실행형 Go 코드 같은 그 밖의 시나리오에서는 Go의 표준 패닉 복구 메커니즘을 사용해 직접 패닉을 처리하는 것이 좋습니다.

## PanicDetails 구조체

패닉이 발생하면 Wails는 패닉에 관한 중요한 정보를 `PanicDetails` 구조체에 기록합니다:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

이 구조체는 패닉에 관한 포괄적인 정보를 제공합니다:

- `StackTrace`: 패닉으로 이어진 호출 스택을 보여 주는 형식화된 문자열
- `Error`: 실제 오류 또는 패닉 메시지
- `Time`: 패닉이 발생한 정확한 시각
- `FullStackTrace`: 런타임 프레임을 포함한 전체 스택 트레이스

@note{type="info" title="Service 코드의 패닉"}
프런트엔드에서 호출된 Service 코드의 패닉이 포착되면 코드에서 패닉이 발생한 정확한 위치에 초점을 맞추도록 스택 트레이스가 잘립니다. 전체 스택 트레이스를 확인하려면 `FullStackTrace` 필드를 사용할 수 있습니다.

@end

## 기본 패닉 핸들러

사용자 지정 패닉 핸들러를 지정하지 않으면 Wails는 오류 정보를 형식화된 로그 메시지로 출력한 다음 종료하는 기본 핸들러를 사용합니다. 예:

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

## 사용자 지정 패닉 핸들러

애플리케이션을 생성할 때 `PanicHandler` 옵션을 설정하여 자체 패닉 핸들러를 구현할 수 있습니다. 다음은 예제입니다:

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

Wails 애플리케이션에서 패닉을 처리하는 완전한 실행 예제는 `v3/examples/panic-handling`의 panic-handling 예제를 참조하세요.

## 마무리 참고 사항

Wails 패닉 핸들러는 바인딩된 메서드의 패닉과 내부 런타임 오류를 처리하기 위한 것임을 기억하세요. 애플리케이션의 다른 부분에서는 적절한 경우 Go의 표준 오류 처리 패턴과 패닉 복구 메커니즘을 사용하는 것이 좋습니다. 모든 Go 애플리케이션이 그렇듯이 가능하면 적절한 오류 처리를 통해 패닉을 방지하는 것이 좋습니다.
