---
title: "충돌 디버깅"
description: "Wails 애플리케이션의 진단 보고서와 미니덤프 수집"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

프로덕션에서 문제가 발생하면 로그만으로는 부족한 경우가 많습니다. 이 가이드는 구조화된 진단 보고서와 Windows에서 WinDbg 또는 Visual Studio로 열 수 있는 미니덤프를 생성하는 `github.com/wailsapp/wails/v3/pkg/debug` 패키지를 설명합니다.

이 패키지는 의도적으로 작게 설계되었으며 진입점은 두 가지입니다:

- `debug.Dump(...)` — 현재 프로세스의 Windows 미니덤프를 기록합니다.
- `debug.Report(...)` — 선택적으로 미니덤프를 포함하는 상세한 진단 스냅샷을 수집합니다.

둘 다 호출한 사용자의 토큰으로 실행됩니다. 권한 상승, `SeDebugPrivilege`, `OpenProcess`는 필요하지 않습니다. 현재 프로세스 덤프에는 자기 프로세스 접근에 대한 ACL 검사를 우회하는 `GetCurrentProcess` 의사 핸들을 사용합니다.

## 빠른 참조

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report`는 일부 수집에 실패해도 항상 `*CrashReport`를 반환합니다. 오류가 발생하기 전에 수집한 정보가 포함되므로 `err != nil`일 때도 보고서는 `nil`이 아닙니다.

## `CrashReport` 구조체

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

충돌 보고 서비스로 보내려면 `encoding/json`으로 직렬화하고, 일부 정보만 필요하면 필드를 직접 읽으세요.

## `PanicHandler`와 통합

가장 유용한 패턴은 Wails가 패닉을 감지한 순간 보고서와 선택적 미니덤프를 수집한 뒤, 결합한 데이터를 자체 보고 파이프라인에 전달하는 것입니다:

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails는 `debug.Report`를 자동으로 호출하지 **않습니다**. 충돌 보고에는 사용자 동의나 개인 식별 정보 삭제가 필요할 수 있으므로 핸들러에서 결정해야 합니다. Wails가 `PanicHandler`를 호출하는 시점과 직접 생성한 고루틴까지 처리하는 방법은 [패닉 처리 가이드](/guides/panic-handling/)를 참조하세요.

## 수동 디버깅

`debug.Dump`와 `debug.Report`는 패닉이 발생하지 않은 상황에서도 유용합니다:

- 멈춤 감지: 감시 타이머가 작동하면 프로세스를 덤프하여 WinDbg에서 각 스레드가 무엇을 기다리는지 확인할 수 있습니다.
- 비정상 상태 보고: 예를 들어 고루틴 수가 끝없이 증가하면 보고서를 수집하고 앱은 계속 실행한 뒤 나중에 조사합니다.
- 요청에 따른 지원 자료: «버그 신고» 메뉴를 `debug.Report(debug.WithDump())`에 연결하고 결과를 사용자의 지원 요청에 첨부합니다.

## Windows 미니덤프 활용

미니덤프는 기본적으로 정보가 풍부하면서도 크기가 작은 플래그 조합을 사용합니다:

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

이 조합은 모든 스레드 스택, 핸들 테이블, 모듈 목록을 보존하는 약 5–50 MB의 덤프를 만듭니다. 덤프 시점에 각 고루틴이 어디에 있었는지 재구성하기에 충분합니다.

`debug.WithFullMemory()` 또는 `Report`의 `debug.WithDumpFullMemory()`는 전체 프로세스 주소 공간을 수집하는 `MiniDumpWithFullMemory`를 추가합니다. 특정 포인터가 손상된 이유를 조사하는 데 유용하지만 파일 크기가 수백 MB 또는 GB에 달하므로 심층 조사에만 사용하세요.

생성된 `.dmp`는 WinDbg에서 `windbg -z C:\path\to\your.dmp`로 열거나 Visual Studio의 File → Open Crash Dump로 엽니다. 릴리스 빌드는 보통 심볼이 제거되므로 유용한 스택 추적을 얻으려면 같은 빌드의 `.pdb` 또는 디버그 정보가 포함된 Go 바이너리를 함께 사용하세요.

## Windows 이외의 플랫폼

Windows 이외에서는 `debug.Dump`가 구현되지 않았다는 오류를 반환합니다. Linux와 macOS는 `GOTRACEBACK=crash`를 통한 코어 덤프, `rr`, 플랫폼별 충돌 보고 도구 등 서로 다른 사후 분석 도구를 사용하며, 이는 `MiniDumpWriteDump` API에 그대로 대응하지 않습니다. `debug.Report`는 계속 작동하며 덤프 파일을 제외한 정보를 제공합니다.

`Dump`가 반환하는 오류를 확인하여 기능을 자연스럽게 축소할 수 있습니다:

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## 보안 참고 사항

- **권한 상승이 필요하지 않습니다.** 이 패키지는 자격 증명 수집 도구가 사용하는 `SeDebugPrivilege` 및 `OpenProcess` / 원격 PID 접근을 의도적으로 피합니다. 호출한 프로세스만 덤프합니다.
- **덤프 파일에는 메모리의 모든 내용이 포함됩니다.** 사용 중인 자격 증명, 인증 토큰, 사용자 데이터, 암호화 키도 포함됩니다. 전체 메모리 이미지처럼 취급하여 전송 중 보호하고, 저장된 데이터를 정리하며, 제삼자에게 업로드하기 전에 민감 정보를 삭제하는 것을 고려하세요.
- **기본 경로는 `os.TempDir()`입니다.** Windows에서는 `%LOCALAPPDATA%\Temp`로 해석됩니다. 사용자별 경로라 누구나 읽을 수는 없지만 파일은 남아 있습니다. 보관할 덤프는 안전한 위치로 옮기세요.
