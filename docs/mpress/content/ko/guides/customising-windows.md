---
title: "Wails에서 창 사용자 지정하기"
description: "Wails 애플리케이션에서 창의 모양과 동작을 사용자 지정합니다"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

관련 플랫폼: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails는 창 컨트롤의 모양과 기능을 제어하는 API를 제공합니다. 이 기능은 Windows와 macOS에서 사용할 수 있지만 Linux에서는 사용할 수 없습니다.

## 창 버튼 상태 설정하기

버튼 상태는 `ButtonState` 열거형으로 정의됩니다.

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`: 버튼이 활성화되어 있으며 표시됩니다.
- `ButtonDisabled`: 버튼이 표시되지만 비활성화되어 있습니다(회색으로 표시됨).
- `ButtonHidden`: 제목 표시줄에서 버튼을 숨깁니다.

창을 생성할 때 또는 런타임에 버튼 상태를 설정할 수 있습니다.

### 창 생성 시 버튼 상태 설정하기

새 창을 생성할 때 `WebviewWindowOptions` 구조체를 사용하여 버튼의 초기 상태를 설정할 수 있습니다.

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

위 예제에서는 최소화 버튼이 숨겨지고, 최대화 버튼이 비활성화되며(회색으로 표시됨), 닫기 버튼은 활성화됩니다.

### 런타임에 버튼 상태 설정하기

`Window` 인터페이스의 다음 메서드를 사용하여 런타임에 버튼 상태를 변경할 수도 있습니다.

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS: MaximiseButtonState와 FullscreenButtonState는 하나의 버튼을 공유합니다

macOS에서 녹색 신호등 버튼(`NSWindowZoomButton`)은 최대화와 전체 화면에 모두 사용되는 동일한 물리적 컨트롤입니다. 별도의 처리가 없다면, 창을 생성할 때 `MaximiseButtonState`과 `FullscreenButtonState`를 서로 다른 값으로 설정할 경우 별도의 알림 없이 마지막으로 설정한 값이 앞선 값을 덮어쓰게 됩니다.

이를 방지하기 위해 Wails는 초기화 시 `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden` 순서에 따라 두 상태 중 <strong>더 제한적인 상태</strong>를 적용합니다.

| `MaximiseButtonState` | `FullscreenButtonState` | macOS에서 적용되는 상태 |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

런타임에는 macOS에서 `SetMaximiseButtonState`과 `SetFullscreenButtonState`이 모두 `NSWindowZoomButton`을 대상으로 하므로 마지막 호출이 적용됩니다.

### 플랫폼별 차이점

버튼 상태 기능은 Windows와 macOS에서 약간 다르게 동작합니다.

|  | Windows | Mac |
| --- | --- | --- |
| 최소화/최대화/닫기 비활성화 | 최소화/최대화/닫기 비활성화 | 최소화/최대화/닫기 비활성화 |
| 최소화 숨기기 | 최소화 비활성화 | 최소화 버튼 숨기기 |
| 최대화 숨기기 | 최대화 비활성화 | 최대화 버튼 숨기기 |
| 닫기 숨기기 | 모든 컨트롤 숨기기 | 닫기 숨기기 |
| `FullscreenButtonState` | 아무 작업도 하지 않음 | 확대/축소(녹색) 버튼을 대상으로 함 |

참고: Windows에서는 최소화/최대화 버튼을 개별적으로 숨길 수 없습니다. 하지만 두 버튼을 모두 비활성화하면 두 컨트롤이 모두 숨겨지고 닫기 버튼만 표시됩니다. Windows의 표준 제목 표시줄에는 전용 전체 화면 버튼이 없으므로 `FullscreenButtonState`은(는) 아무 동작도 하지 않습니다.

### 창 스타일 제어(Windows)

Windows에서 제목 표시줄의 스타일을 제어하려면 `WebviewWindowOptions` 구조체의 `ExStyle` 필드를 사용할 수 있습니다:

예:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

창의 확장 스타일에 영향을 주는 다음 옵션은 이 설정으로 재정의됩니다:

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
