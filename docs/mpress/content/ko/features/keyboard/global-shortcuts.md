---
title: "전역 단축키"
description: "애플리케이션에 포커스가 없을 때도 작동하는 시스템 전역 키보드 단축키를 등록합니다"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

전역 단축키는 Wails 애플리케이션이 실행되는 동안 현재 어떤 애플리케이션에 포커스가 있는지와 관계없이 작동하는 시스템 전역 키보드 단축키입니다. 표시/숨기기 단축키, 빠른 캡처 도구, 미디어 제어 기능처럼 사용자가 어디서나 이용할 수 있기를 기대하는 기능에 적합합니다.

@note{type="info" title="전역 단축키와 키 바인딩 비교"}
[키 바인딩](/features/keyboard/shortcuts/)(`app.KeyBinding`)은 애플리케이션 창 중 하나에 포커스가 있을 때만 작동합니다. 전역 단축키(`app.GlobalShortcut`)는 애플리케이션이 백그라운드에 있을 때도 시스템 전역에서 작동합니다. 필요에 맞는 방식을 사용하세요.

@end

전역 단축키는 각 플랫폼의 네이티브 기능을 직접 기반으로 구현되며 서드 파티 종속성을 추가하지 않습니다.

## 전역 단축키 관리자에 접근하기

애플리케이션 인스턴스의 `GlobalShortcut` 속성을 통해 관리자에 접근할 수 있습니다.

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 단축키 등록하기

`Register`은 액셀러레이터와 콜백을 인수로 받습니다. 단축키를 누를 때마다 콜백이 별도의 고루틴에서 실행됩니다.

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

`app.Run`을 호출하기 전에 단축키를 등록할 수 있습니다. 그러면 애플리케이션이 시작될 때 운영 체제와의 바인딩이 자동으로 수행됩니다.

@note{type="tip" title="콜백에서 UI 조작하기"}
콜백은 메인 스레드 외부에서 실행됩니다. 콜백에서 창이나 다른 UI와 상호 작용해야 하는 경우에는 창 메서드가 이를 대신 처리하지만, 메인 스레드에서 사용자 지정 작업을 실행하려면 해당 작업을 `application.InvokeSync`으로 감싸세요.

@end

### 액셀러레이터 형식

전역 단축키는 메뉴 액셀러레이터 및 키 바인딩과 같은 액셀러레이터 형식을 사용합니다.

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl`은 macOS에서는 Command로, Windows와 Linux에서는 Control로 해석되므로 크로스 플랫폼 단축키를 편리하게 정의할 수 있습니다.

## 단축키 관리하기

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

등록된 모든 단축키는 애플리케이션이 종료될 때 자동으로 해제되므로 수동으로 정리할 필요가 없습니다.

## 같은 단축키를 두 번 등록하면 발생하는 동작

서로 다른 두 가지 경우가 있으며 Wails는 각 경우를 다르게 처리합니다.

### 같은 애플리케이션에서 단축키를 두 번 등록하는 경우

이 경우는 Wails 자체에서 처리하며 모든 플랫폼에서 동일하게 동작합니다. 두 번째 `Register` 호출은 오류를 반환하고 기존 바인딩은 그대로 유지됩니다("오류를 반환하고 유지"). 따라서 작동 중인 단축키를 조용히 대체하는 대신 실수를 드러내어 예측 가능한 동작을 유지합니다.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

단축키의 콜백을 변경하려면 먼저 해당 단축키를 `Unregister`한 다음 `Register`하세요.

### 다른 애플리케이션이 이미 단축키를 소유한 경우

이 경우는 운영 체제가 결정하므로 결과는 플랫폼별로 다릅니다.

| 플랫폼 | 다른 애플리케이션이 단축키를 소유한 경우의 동작 |
| --- | --- |
| **macOS** | 등록에 성공합니다. macOS에서는 여러 애플리케이션이 같은 단축키를 등록할 수 있으므로 등록이 거부되지 않고 기존 소유자의 콜백과 함께 애플리케이션의 콜백이 추가됩니다. |
| **Windows** | 등록에 실패하고 `Register`이 오류를 반환합니다. 단축키를 먼저 등록한 애플리케이션이 계속 소유합니다. |
| **Linux (X11)** | X 서버가 같은 키 조합의 두 번째 그랩을 거부하므로 등록에 실패하고 `Register`이 오류를 반환합니다. |
| **Linux (Wayland)** | 컴포지터가 조정합니다. 일반적으로 사용자에게 데스크톱의 전역 단축키 대화 상자를 통해 바인딩을 승인하거나 선택하라는 메시지가 표시됩니다. |

이러한 차이가 있으므로 항상 `Register`이 반환하는 오류를 확인하고, 단축키를 확보할 수 없을 때는 대체 단축키를 제공하거나 사용자에게 피드백을 표시하세요.

## 플랫폼별 고려 사항

@tabs
[macOS]
전역 단축키는 Carbon Event Manager의 단축키 API를 사용합니다. 이는 macOS에서 시스템 전역 단축키를 구현하는 표준 메커니즘이며 손쉬운 사용 권한이 필요하지 않습니다.

단축키는 키의 물리적 위치에 바인딩되므로 QWERTY가 아닌 레이아웃에서도 표준 ANSI/QWERTY 위치의 키에 매핑됩니다.

@note{type="caution" title="숨기기 단축키와 `ApplicationShouldTerminateAfterLastWindowClosed`"}
macOS에서 `window.Hide()`은 `orderOut:`을 사용하여 창을 보이지 않게 합니다. AppKit은 보이지 않는 마지막 창을 닫힌 것으로 처리하므로 `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`을 설정한 상태에서 전역 단축키를 사용해 유일한 창을 숨기면 애플리케이션이 백그라운드에 남지 않고 종료됩니다. 숨기기/표시 단축키를 사용하는 경우에는 이 옵션을 설정하지 않은 상태(기본값)로 두어야 창을 숨겼다가 나중에 다시 불러올 수 있습니다.

@end

[Windows]
전역 단축키는 Win32 `RegisterHotKey` API를 사용합니다. 자동 반복이 억제되므로 키를 계속 누르고 있어도 콜백이 반복해서 실행되지 않고 한 번만 실행됩니다.

다른 애플리케이션이 이미 해당 키 조합을 소유하고 있으면 등록에 실패하므로 기본값으로는 덜 흔한 키 조합을 사용하세요.

[Linux]
**X11** 세션에서 Wails는 X 서버로부터 바로 단축키를 가져오므로 요청한 가속 키가 지정한 그대로 바인딩됩니다.

**Wayland** 세션에서는 설계상 애플리케이션이 키를 직접 가져올 방법이 없습니다. 대신 Wails는 XDG Desktop Portal `org.freedesktop.portal.GlobalShortcuts` 인터페이스를 사용합니다. 포털을 사용할 때 전달하는 가속 키는 *선호* 트리거이며, 최종 키 조합은 컴포지터(궁극적으로는 사용자)가 결정합니다. 단축키가 활성화되면 콜백은 계속 호출되지만, 실제 키가 요청한 키와 정확히 일치한다고 보장할 수 없으며 `IsRegistered`/`GetAll`은 컴포지터가 바인딩한 조합이 아니라 애플리케이션이 API에 전달한 조합을 보고합니다.

포털을 사용하려면 전역 단축키 포털을 구현하는 데스크톱 환경(예: 최신 GNOME 또는 KDE Plasma)이 필요합니다.

@end

## 전체 예제

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="중요한 시스템 단축키 피하기"}
일부 키 조합은 운영 체제나 데스크톱 환경에서 예약되어 있어 애플리케이션이 가져올 수 없습니다. 충돌할 가능성이 낮은 기본값을 선택하고 `Register`에서 반환되는 오류를 항상 처리하세요.

@end
