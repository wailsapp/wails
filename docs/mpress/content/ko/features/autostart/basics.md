---
title: "자동 시작"
description: "macOS, Windows 및 Linux에서 사용자 로그인 시 애플리케이션이 실행되도록 등록합니다"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## 자동 시작

`app.Autostart`는 사용자가 로그인할 때 애플리케이션이 자동으로 실행되도록 등록합니다. 플랫폼별로 적절한 네이티브 메커니즘을 선택하고 심볼릭 링크가 적용된 설치 경로(Homebrew, Scoop)를 해석하므로, 바이너리를 업그레이드해도 등록이 손상되지 않습니다.

등록은 즉시 적용되지 않고 <strong>다음 로그인</strong>부터 적용됩니다.

## 빠른 시작

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

기본 옵션으로 로그인 시 애플리케이션이 실행되도록 등록합니다.

```go
func (m *AutostartManager) Enable() error
```

`Enable`를 반복해서 호출해도 안전합니다. 호출할 때마다 등록을 덮어쓰므로, 사용자의 기본 설정을 저장해 두었다면 시작할 때마다 호출해도 됩니다.

### `EnableWithOptions`

사용자 지정 옵션으로 등록합니다.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| 필드 | 타입 | 설명 |
| --- | --- | --- |
| `Identifier` | `string` | 자동으로 파생된 등록 ID를 재정의합니다. 아래의 "식별자"를 참조하세요. |
| `Arguments` | `[]string` | 로그인 시 실행할 때 실행 파일 경로 뒤에 추가되는 인수입니다(예: `--hidden`). |

### `Disable`

자동 시작 등록을 제거합니다. 애플리케이션이 등록되어 있지 않으면 `nil`를 반환합니다. 비활성화 작업은 멱등성을 갖습니다.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

등록이 존재하는지 알려 줍니다. 등록된 경로를 검증하지 않으므로 빠릅니다.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

전체 등록 상태를 반환합니다.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| 필드 | 타입 | 설명 |
| --- | --- | --- |
| `Enabled` | `bool` | 등록의 존재 여부입니다. |
| `Path` | `string` | 등록 아티팩트의 디스크상 위치입니다(plist 경로, `.desktop` 경로 또는 레지스트리 하위 키). `Enabled`가 false이면 비어 있습니다. |
| `Strategy` | `AutostartStrategy` | 앱을 등록한 메커니즘입니다([플랫폼별 동작](#--2) 참조). |

## 플랫폼별 동작

@tabs{sync-key="platform"}
[macOS]
앱 패키징 방식에 따라 다음 두 메커니즘을 사용합니다.

- **macOS 13 이상, 번들된 `.app`**: `SMAppService.mainAppService`. 샌드박스 앱과 Mac App Store 빌드에서 작동합니다. TCC 자동화 프롬프트가 표시되지 않습니다(기존 AppleScript 방식에서는 이 프롬프트가 표시되었습니다).
- **13 이전 버전의 macOS 또는 번들되지 않은 바이너리**: `RunAtLoad=true`를 사용하여 `~/Library/LaunchAgents/<identifier>.plist`에 LaunchAgent plist를 작성합니다.

호출자가 어떤 경로를 사용했는지 알 수 있도록 `Status()`는 `AutostartStrategySMAppService` 또는 `AutostartStrategyLaunchAgent`를 반환합니다.

앱이 번들되지 않은 형태에서 번들된 형태로 업그레이드되면 `Status()`에서 두 경로를 모두 확인하고 `Disable()`에서 정리하므로, 남겨진 LaunchAgent가 이전 빌드를 계속 실행하지 않습니다.

[Windows]
`HKCU\Software\Microsoft\Windows\CurrentVersion\Run` 아래에 레지스트리 값을 추가합니다. 값 이름은 자동 시작 `Identifier`로 설정하고, 데이터는 따옴표로 묶은 실행 파일 경로에 모든 `Arguments`를 더한 값으로 설정합니다.

공백이나 따옴표가 포함된 경로도 올바르게 왕복 변환되도록, 인수 인용은 `CommandLineToArgvW` 규칙을 따릅니다(따옴표 앞의 백슬래시는 두 개로 만듭니다).

`Status().Strategy`는 `AutostartStrategyRegistryRun`를 반환합니다.

[Linux]
다음 설정으로 XDG 자동 시작 항목을 `$XDG_CONFIG_HOME/autostart/<identifier>.desktop`에 작성합니다(기본값은 `~/.config/autostart/`).

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

`Exec` 필드는 [freedesktop.org Desktop Entry 사양](https://specifications.freedesktop.org/desktop-entry-spec/)에 따라 이스케이프합니다. 예약 문자(`"`, `` ` ``, `$`, `\\`)는 백슬래시로 이스케이프하고, 값에 공백 문자가 포함되어 있으면 큰따옴표로 묶습니다.

`Status().Strategy`는 `AutostartStrategyXDGAutostart`를 반환합니다.

[iOS / Android / 서버]
지원되지 않습니다. 모든 메서드는 `ErrAutostartNotSupported`를 반환합니다. 이를 명확하게 감지하려면 `errors.Is(err, application.ErrAutostartNotSupported)`를 사용하세요.

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## 식별자

`Options.Identifier`가 비어 있으면 앱 이름에서 기본값을 파생합니다.

| 플랫폼 | 기본 식별자 |
| --- | --- |
| macOS(번들됨) | 앱의 번들 식별자(예: `com.example.MyApp`) |
| macOS(번들되지 않음) | `application.Options.Name`에서 `<slug>`를 파생한 `wails.autostart.<slug>` |
| Windows | `application.Options.Name`의 슬러그(소문자로 변환하고 `A-Za-z0-9._-`가 아닌 문자를 제거하며 공백을 대시로 변경) |
| Linux | Windows와 동일한 슬러그 |

식별자는 `^[A-Za-z0-9._-]+$`와 일치해야 하며 200자를 초과해서는 안 됩니다. macOS에서는 역방향 DNS 형식을 권장합니다(launchd Label을 관례적으로 작성하는 방식과 일치합니다).

`AutostartOptions.Identifier`을 재정의하면 Windows에서는 레지스트리 값 이름으로, Linux에서는 `.desktop` 파일 이름으로 같은 식별자를 재사용하므로 하나의 문자열로 플랫폼 전반의 등록을 식별할 수 있습니다.

## 오래된 등록 감지

`Disable()`과 `Status()`은 식별자를 조회하는 대신, <strong>등록된 실행 파일 경로를 `os.Executable()`(모든 심볼릭 링크를 해석한 경로)</strong>와 대조하여 등록을 찾습니다. 따라서 다음과 같이 동작합니다.

- **릴리스 간에 식별자를 변경해도 안전합니다.** 실행 파일 경로가 동일하다면 `Status()`에서 이전 등록을 계속 찾을 수 있으며 `Disable()`에서 이를 정리합니다.
- **다른 경로에 있는 앱의 두 번째 복사본이 첫 번째 복사본의 등록을 덮어쓰지 않습니다.** 각 바이너리 위치는 독립적으로 추적됩니다.
- **심볼릭 링크를 사용하는 설치(Homebrew, Scoop)는 안정적으로 유지됩니다.** 대조하기 전에 `os.Executable()`에 `filepath.EvalSymlinks`을 적용하므로, Homebrew 업그레이드로 링크 대상이 바뀌어도 항목이 고립되지 않습니다.

여기서 *처리하지 않는* 경우는 사용자가 바이너리를 관련 없는 경로로 이동하거나 이름을 변경하는 경우입니다. 이때 이전 등록은 현재 존재하지 않는 파일을 가리키므로 고립됩니다. 안정적인 설치 경로에서 배포되는 앱은 이를 걱정할 필요가 없습니다. 이식 가능한 단일 파일 바이너리로 배포되는 앱은 자체 파일을 이동하기 전에 `Disable()`을 호출하거나 항상 안정적인 심볼릭 링크를 통해 실행하는 것이 좋습니다.

## 예제

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

상태 확인, 활성화 및 비활성화 버튼을 갖춘 실행 가능한 전체 예제는 [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart)에 있습니다.
