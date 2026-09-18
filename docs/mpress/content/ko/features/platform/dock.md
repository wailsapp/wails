---
title: "Dock 및 작업 표시줄"
description: "macOS 및 Windows에서 Dock 아이콘 표시 여부를 관리하고 배지를 표시합니다"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## 소개

Wails는 데스크톱 애플리케이션을 위한 크로스 플랫폼 Dock 서비스를 제공합니다. 이 서비스를 사용하면 다음 작업을 수행할 수 있습니다.

- macOS Dock에서 애플리케이션 아이콘 숨기기 및 표시
- 애플리케이션 타일 또는 Dock/작업 표시줄 아이콘에 배지 표시(macOS 및 Windows)

## 기본 사용법

### 서비스 생성

먼저 Dock 서비스를 초기화합니다.

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### 사용자 지정 배지 옵션으로 서비스 생성(Windows 전용)

Windows에서는 다양한 옵션으로 배지 모양을 사용자 지정할 수 있습니다.

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Dock 작업

### Dock 앱 아이콘 숨기기

macOS Dock에서 앱 아이콘을 숨깁니다.

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Dock 앱 아이콘 표시하기

macOS Dock에 앱 아이콘을 표시합니다.

```go
// Show the app icon
dockService.ShowAppIcon()
```

## 배지 작업

### 배지 설정

애플리케이션 타일/Dock 아이콘에 배지를 설정합니다.

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### 사용자 지정 배지 설정(Windows 전용)

일회성 옵션을 적용하여 배지를 설정합니다.

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### 배지 제거

애플리케이션 아이콘에서 배지를 제거합니다.

```go
dockService.RemoveBadge()
```

### 설정된 배지 가져오기

```go
dockService.GetBadge()
```

## 플랫폼별 고려 사항

@tabs
[macOS]
macOS에서는 다음과 같이 동작합니다.

- Dock 아이콘을 **숨기거나** **표시할** 수 있습니다
- 배지가 Dock 아이콘에 직접 표시됩니다
- 배지 옵션은 **사용자 지정할 수 없습니다**(`NewWithOptions`/`SetCustomBadge`에 전달된 모든 옵션은 무시됩니다)
- 표시 모드에 자동으로 맞춰지는 표준 macOS Dock 배지 스타일을 사용합니다
- 레이블 오버플로는 시스템에서 처리합니다
- 빈 레이블을 지정하면 기본 배지인 "●"가 표시됩니다

[Windows]
Windows에서는 다음과 같이 동작합니다.

- 이 서비스는 현재 작업 표시줄 아이콘 숨기기/표시를 지원하지 않습니다
- 배지가 작업 표시줄의 오버레이 아이콘으로 표시됩니다
- 배지에서 텍스트 값을 지원합니다
- `BadgeOptions`을 통해 배지 모양을 사용자 지정할 수 있습니다
- 배지를 표시하려면 애플리케이션에 창이 있어야 합니다
- 여러 글자로 된 레이블에는 더 작은 글꼴 크기가 자동으로 적용됩니다
- 레이블 오버플로는 처리되지 않습니다
- 사용자 지정 옵션:
  - **TextColour**: 텍스트 색상(기본값: 흰색)
  - **BackgroundColour**: 배지 배경색(기본값: 빨간색)
  - **FontName**: 글꼴 파일 이름(기본값: "segoeuib.ttf")
  - **FontSize**: 한 글자용 글꼴 크기(기본값: 18)
  - **SmallFontSize**: 여러 글자용 글꼴 크기(기본값: 14)


[Linux]
Linux에서는 다음과 같이 동작합니다.

- Dock 아이콘 표시 여부 및 배지 기능을 사용할 수 없습니다

@end

## 권장 사항

1. **Dock 아이콘을 숨길 때(macOS):**
  - 사용자가 계속 앱에 접근할 수 있도록 해야 합니다(예: [시스템 트레이](/features/menus/systray/) 사용)
  - 대체 UI에 "종료" 옵션을 포함하세요
  - 앱이 Command+Tab 전환 화면에 표시되지 않습니다
  - 열려 있는 창은 계속 표시되며 정상적으로 작동합니다
  - 모든 창을 닫아도 앱이 종료되지 않을 수 있습니다(macOS 동작은 상황에 따라 다름)
  - 사용자는 Dock에서 마우스 오른쪽 버튼을 클릭하여 종료하는 표준 방법을 사용할 수 없게 됩니다


2. **배지는 필요한 경우에만 사용하세요.**
  - 배지를 너무 자주 업데이트하면 사용자의 주의를 분산시킬 수 있습니다
  - 중요한 알림에만 배지를 사용하세요


3. **배지 텍스트는 짧게 유지하세요.**
  - 숫자 배지가 가장 효과적입니다
  - macOS에서는 텍스트 배지를 간결하게 작성하는 것이 좋습니다


4. **Windows 배지를 사용자 지정할 때:**
  - 텍스트와 배경색의 대비를 높게 유지하세요
  - 텍스트가 길어질수록 글꼴 크기가 줄어드므로 다양한 텍스트 길이로 테스트하세요
  - 사용 가능성을 보장하려면 일반적인 시스템 글꼴을 사용하세요


## API 참조

### 서비스 관리

| 메서드 | 설명 |
| --- | --- |
| `New()` | 새 Dock 서비스를 생성합니다 |
| `NewWithOptions(options BadgeOptions)` | 사용자 지정 배지 옵션을 사용하여 새 Dock 서비스를 생성합니다(Windows 전용이며, macOS와 Linux에서는 옵션이 무시됨) |

### Dock 작업

| 메서드 | 설명 |
| --- | --- |
| `HideAppIcon()` | macOS Dock에서 앱 아이콘을 숨깁니다(macOS 전용) |
| `ShowAppIcon()` | macOS Dock에 앱 아이콘을 표시합니다(macOS 전용) |

### 배지 작업

| 메서드 | 설명 |
| --- | --- |
| `SetBadge(label string) error` | 지정된 레이블로 배지를 설정합니다 |
| `SetCustomBadge(label string, options BadgeOptions) error` | 지정된 레이블과 사용자 지정 스타일 옵션으로 배지를 설정합니다(Windows 전용) |
| `RemoveBadge() error` | 애플리케이션 아이콘에서 배지를 제거합니다 |
| `GetBadge() *string` | 현재 배지를 가져옵니다 |

### 구조체 및 형식

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
