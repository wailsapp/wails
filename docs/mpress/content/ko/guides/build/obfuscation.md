---
title: "난독화 빌드"
description: "소스 코드를 리버스 엔지니어링으로부터 보호하도록 Garble로 Wails 애플리케이션 빌드하기"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble)은 `go build`을 대체하여 심볼 이름을 변경하고, 상수를 난독화하며, 생성된 바이너리에서 디버그 정보를 제거하는 Go 빌드 도구입니다. Wails v3에서는 두 가지 새로운 명령을 통해 Garble을 기본 지원합니다.

## 사전 요구 사항

- **Go 1.26.2 이상** — Garble v0.16.0에 필요합니다.
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Garble에 필요한 최소 Go 버전은 릴리스마다 달라집니다. 이전 Go 툴체인을 사용 중이라면 설치하기 전에 [Garble 릴리스 페이지](https://github.com/burrowers/garble/releases)에서 사용 중인 툴체인과 호환되는 버전을 확인하세요.

@end

## 필수: 서비스 타입에 JSON 태그 추가하기

바인딩된 서비스 메서드가 반환하거나 인수로 받는 모든 구조체에는 내보낸 각 필드에 명시적인 JSON 태그가 있어야 합니다.

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble은 내보낸 구조체 필드의 이름을 변경하며, Wails는 Garble이 정적으로 추적할 수 없는 `interface{}` 매개변수를 통해 이러한 구조체를 `json.Marshal`에 전달합니다. JSON 태그 없이 난독화 빌드를 실행하면 컴파일에는 성공하지만 런타임에 프런트엔드가 난독화되거나 비어 있는 필드 이름을 받게 됩니다. 난독화 빌드를 실행하기 전에 태그를 추가하세요.

@end

Wails 자체 타입인 `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities`에는 이미 태그가 지정되어 있습니다. 직접 정의한 타입에만 태그를 지정하면 됩니다.

## 난독화하여 빌드하기

@steps
### 안정적인 ID 파일 생성
바인딩된 서비스 메서드를 추가하거나, 이름을 변경하거나, 제거할 때마다 다음을 실행하세요.

```bash
wails3 generate bindings -obfuscated
```

그러면 main 패키지 디렉터리에 `wails_obfuscated.gen.go`이 생성됩니다. 이 파일을 커밋하세요.

### Garble로 빌드
```bash
wails3 build --obfuscated
```

난독화된 바인딩을 사용하여 애플리케이션을 빌드합니다.

@end

## Garble에 추가 플래그 전달하기

옵션을 `garble`에 직접 전달하려면 `--garbleargs`을 사용하세요.

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

지원되는 플래그의 전체 목록은 [Garble 문서](https://github.com/burrowers/garble#flags)를 참조하세요.

## 고급: ID 파일을 다른 패키지에 쓰기

기본적으로 `wails_obfuscated.gen.go`은 `main` 패키지와 같은 위치에 기록됩니다. 프로젝트에서 서비스를 `main`이 가져오는 하위 패키지에 두는 경우에는 `-obfuscated-output`을 사용하여 파일을 해당 위치에 기록할 수 있습니다.

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
대상 패키지의 `init()`이 시작 시 실행되도록 해당 패키지를 `main` 패키지에서 직접 또는 전이적으로 가져와야 합니다. 해당 패키지에 도달할 수 없으면 안정적인 ID가 등록되지 않으며, 바인딩 호출이 실패합니다(예: 런타임에 `binding not found` 오류 발생).

@end

## 문제 해결

### `garble: command not found`

Garble이 설치되어 있지 않거나 `$(go env GOPATH)/bin`이 `PATH`에 없습니다.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 프런트엔드가 잘못되거나 비어 있는 필드 값을 받음

서비스 반환 타입에 `json:"..."` 태그가 없습니다. 바인딩된 메서드가 반환하는 모든 구조체를 확인하고 내보낸 각 필드에 명시적인 태그를 추가하세요.

### 브라우저 콘솔의 `binding not found` 오류

안정적인 ID 파일이 없거나 컴파일에 포함되지 않았습니다. 다음을 확인하세요.

- main 패키지 디렉터리(또는 `-obfuscated-output`에 전달한 디렉터리)에 `wails_obfuscated.gen.go`이 있는지 확인하세요.
- `wails_obfuscated` 빌드 태그를 추가하는 `wails3 build --obfuscated`을 실행했는지 확인하세요.
- `-obfuscated-output`을 사용했다면 대상 패키지를 `main`에서 가져오는지 확인하세요.

### Windows Defender가 빌드를 바이러스로 탐지함

Garble로 난독화된 Go 바이너리는 디버그 심볼이 없고 패킹된 실행 파일과 유사하므로 빌드 중 Windows Defender가 휴리스틱 방식으로 탐지합니다. 빌드는 다음 오류와 함께 실패합니다.

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

임시 디렉터리(Go가 중간 빌드 아티팩트를 기록하는 위치)와 프로젝트 디렉터리를 Defender의 제외 목록에 추가하세요.

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

이러한 제외 설정은 지정된 경로에만 적용되며 Defender를 시스템 전체에서 비활성화하지 않습니다.

### `unsupported Go version` 오류로 빌드 실패

Garble v0.16.0에는 Go 1.26.2 이상이 필요합니다. Go를 업그레이드하거나 [Garble 릴리스 페이지](https://github.com/burrowers/garble/releases)에서 사용 중인 툴체인과 호환되는 버전을 확인하세요.
