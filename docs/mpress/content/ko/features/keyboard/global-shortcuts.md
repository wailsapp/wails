---
title: "전역 단축키"
description: "애플리케이션에 포커스가 없어도 실행되는 시스템 전역 키보드 단축키를 등록합니다"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

전역 단축키는 Wails 애플리케이션이 실행 중인 동안 현재 어떤 애플리케이션에 포커스가 있는지와 관계없이 작동하는 시스템 전역 키보드 단축키입니다. 표시/숨기기 단축키, 빠른 캡처 도구, 미디어 제어 및 사용자가 어디서든 이용할 수 있기를 기대하는 기타 기능에 적합합니다.

@note{type="info" title="전역 단축키와 키 바인딩 비교"}
[키 바인딩](/features/keyboard/shortcuts/)(`app.KeyBinding`)은 애플리케이션 창 중 하나에 포커스가 있을 때만 실행됩니다. 전역 단축키(`app.GlobalShortcut`)는 애플리케이션이 백그라운드에 있을 때도 시스템 전역에서 실행됩니다. 필요에 맞는 방식을 사용하세요.

@end

전역 단축키는 각 플랫폼의 네이티브 기능을 직접 기반으로 하며 타사 종속성을 추가하지 않습니다.

## 전역 단축키 관리자에 접근하기

관리자는 애플리케이션 인스턴스의 `GlobalShortcut` 속성에서 사용할 수 있습니다.

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## 단축키 등록하기

`Register`은 액셀러레이터와 콜백을 인수로 받습니다. 단축키를 누를 때마다 콜백이 별도의 goroutine에서 실행됩니다.

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
콜백은 메인 스레드 외부에서 실행됩니다. 콜백에서 창이나 다른 UI와 상호 작용해야 하는 경우 창 메서드가 이를 자동으로 처리하지만, 사용자 정의 작업을 메인 스레드에서 실행하려면 `application.InvokeSync`으로 감싸세요.

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

`CmdOrCtrl`은 macOS에서는 Command로, Windows와 Linux에서는 Control로 해석되므로 크로스 플랫폼 단축키를 편리하게 만들 수 있습니다.

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

## 같은 단축키를 두 번 등록하면 어떻게 되나요?

서로 다른 두 가지 경우가 있으며 Wails는 각 경우를 다르게 처리합니다.

### 같은 애플리케이션이 단축키를 두 번 등록하는 경우

이 경우 Wails 자체에서 처리하며 모든 플랫폼에서 동일하게 작동합니다. 두 번째 `Register` 호출은 오류를 반환하고 기존 바인딩은 그대로 유지됩니다("오류를 반환하고 유지"). 따라서 작동 중인 단축키를 아무 알림 없이 교체하지 않고 실수를 명확히 드러내어 동작을 예측 가능하게 유지합니다.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

단축키의 콜백을 변경하려면 먼저 `Unregister`한 다음 다시 `Register`하세요.

### 다른 애플리케이션이 이미 단축키를 사용 중인 경우

이는 운영 체제에서 결정하므로 결과가 플랫폼마다 다릅니다.

| 플랫폼 | 다른 애플리케이션이 단축키를 사용 중일 때의 동작 |
| --- | --- |
| **macOS** | 등록에 성공합니다. macOS에서는 여러 애플리케이션이 같은 단축키를 등록할 수 있으므로 요청을 거부하지 않고 기존 소유자의 콜백과 함께 해당 콜백을 추가합니다. |
| **Windows** | 등록에 실패하고 `Register`이 오류를 반환합니다. 단축키를 먼저 등록한 애플리케이션이 계속 사용합니다. |
| **Linux(X11)** | X 서버가 같은 조합에 대한 두 번째 grab을 거부하므로 등록에 실패하고 `Register`이 오류를 반환합니다. |
| **Linux(Wayland)** | 컴포지터가 조정합니다. 일반적으로 사용자에게 데스크톱의 전역 단축키 대화 상자를 통해 바인딩을 승인하거나 선택하라는 메시지가 표시됩니다. |

이러한 차이로 인해 항상 `Register`이 반환하는 오류를 확인하고, 단축키를 확보할 수 없을 때는 대체 단축키나 사용자 피드백을 제공하세요.

## 플랫폼별 고려 사항

@tabs
[macOS]
전역 단축키는 Carbon Event Manager의 단축키 API를 사용합니다. 이는 macOS에서 시스템 전역 단축키를 구현하는 표준 메커니즘이며 손쉬운 사용 권한이 필요하지 않습니다.

문자 단축키는 키보드 레이아웃을 인식합니다. Carbon API는 물리적인 키 위치에 바인딩되므로 Wails는 문자를 현재 활성 키보드 레이아웃에서 해당 문자를 입력하는 키로 변환합니다(`UCKeyTranslate` 사용). 따라서 `Cmd+A`은 AZERTY, QWERTZ, Dvorak 및 이와 유사한 레이아웃에서 QWERTY 위치가 아니라 **A**라고 표시된 키에 배치됩니다. 숫자, 문장 부호 및 이름이 지정된 키(화살표 키, 기능 키 등)는 고정된 물리적 위치를 사용합니다. 바인딩은 단축키를 등록할 때 결정됩니다. 애플리케이션이 실행되는 동안 키보드 레이아웃을 변경해도 기존 문자 단축키는 처음 바인딩된 물리적 키를 계속 사용합니다.

@note{type="caution" title="숨기기 단축키와 `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()`은 macOS에서 `orderOut:`을 사용하므로 창이 보이지 않게 됩니다. AppKit은 마지막으로 보이지 않게 된 창을 닫힌 것으로 처리합니다. 따라서 `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true`을 설정하고 전역 단축키로 유일한 창을 숨기면 애플리케이션이 백그라운드에 남아 있지 않고 종료됩니다. 숨기기/표시 단축키를 사용하는 경우에는 이 옵션을 설정하지 않은 상태(기본값)로 두어 창을 숨겼다가 나중에 다시 불러올 수 있게 하세요.

@end

[Windows]
전역 단축키는 Win32 `RegisterHotKey` API를 사용합니다. 자동 반복이 억제되므로 키를 길게 눌러도 콜백이 반복적으로 실행되지 않고 한 번만 실행됩니다.

다른 애플리케이션이 이미 해당 조합을 사용 중이면 등록에 실패하므로 기본값으로는 자주 사용되지 않는 조합을 선택하세요.

[Linux]
**X11** 세션에서 Wails는 X 서버로부터 단축키를 직접 grab하므로 요청한 액셀러레이터가 지정한 그대로 바인딩됩니다.

**Wayland** 세션에서는 설계상 애플리케이션이 키를 직접 grab할 수 없습니다. 대신 Wails는 XDG Desktop Portal의 `org.freedesktop.portal.GlobalShortcuts` 인터페이스를 사용합니다. 포털을 사용할 때 전달한 액셀러레이터는 *선호하는* 트리거이며 최종 키 조합은 컴포지터가 조정하고 궁극적으로 사용자가 결정합니다. 단축키가 활성화되면 콜백은 여전히 실행되지만 실제 키가 요청한 키와 일치한다고 보장할 수 없으며, `IsRegistered`/`GetAll`은 컴포지터가 바인딩한 값이 아니라 사용자가 요청한 값을 보고합니다.

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
일부 조합은 운영 체제나 데스크톱 환경에서 예약하므로 애플리케이션이 확보할 수 없습니다. 충돌할 가능성이 낮은 기본값을 선택하고 항상 `Register`의 오류를 처리하세요.

@end
