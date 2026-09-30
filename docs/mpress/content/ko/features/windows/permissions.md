---
title: "권한"
description: "웹 콘텐츠의 카메라, 마이크, 위치 정보 및 기타 기능 요청 제어"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

`navigator.mediaDevices.getUserMedia()`, Geolocation API 또는 Notifications API를 호출하는 웹 콘텐츠의 요청은 호스트 애플리케이션에서 허용하거나 거부해야 합니다. Wails는 `WebviewWindowOptions`에 크로스 플랫폼 `Permissions` 맵을 제공하므로 플랫폼별 코드 없이 이를 선언적으로 제어할 수 있습니다.

## 빠른 시작

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

해당 창의 웹 콘텐츠에서 요청한 카메라 및 마이크 접근은 브라우저 메시지를 표시하지 않고 허용됩니다.

## 권한 유형

`PermissionType`(uint8)은 웹 콘텐츠에서 요청할 수 있는 기능을 식별합니다.

| 상수 | 기능 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 권한 값

`Permission`(uint8)은 지정된 유형에 적용할 정책입니다.

| 상수 | 값 | 의미 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 플랫폼의 기본 처리 방식 사용(아래 참조) |
| `PermissionAllow` | 1 | 메시지를 표시하지 않고 허용 |
| `PermissionDeny` | 2 | 메시지를 표시하지 않고 거부 |

`PermissionDefault`은 0 값이므로 맵에 설정되지 않은 항목은 기본값으로 동작합니다.

## 플랫폼별 동작

기반 웹뷰의 기본 동작이 서로 다르므로 각 플랫폼은 `PermissionDefault`을 다르게 처리합니다.

### Linux (WebKitGTK)

WebKitGTK에는 **기본 권한 요청 메시지가 없습니다**. 핸들러가 연결되어 있지 않으면 모든 요청을 사용자에게 알리지 않고 거부합니다. 이 때문에 이 기능이 추가되기 전에는 `getUserMedia`이 항상 `NotAllowedError`을 반환했습니다.

현재 Wails는 Linux에서 **카메라 및 마이크** 요청을 처리합니다. 위치 정보, 알림 및 클립보드 읽기는 아직 연결되지 않았으므로 설정한 정책과 관계없이 계속 거부됩니다.

| 정책 | 카메라/마이크 | 위치 정보, 알림, 클립보드 |
| --- | --- | --- |
| `PermissionDefault` | **허용됨**(getUserMedia 복원) | 항상 거부됨 |
| `PermissionAllow` | 허용됨 | 항상 거부됨(아직 구현되지 않음) |
| `PermissionDeny` | 거부됨 | 항상 거부됨 |

### Windows (WebView2)

WebView2에는 기본 권한 요청 메시지와 종류별 권한 API가 있습니다. 다섯 가지 기능 유형을 모두 완전히 지원합니다.

| 정책 | 동작 |
| --- | --- |
| `PermissionDefault` | WebView2가 자체 OS/브라우저 권한 요청 메시지를 표시함 |
| `PermissionAllow` | 사용자에게 알리지 않고 허용됨 |
| `PermissionDeny` | 사용자에게 알리지 않고 거부됨 |

**중요:** 이 기능이 도입되기 전에는 Wails가 조건 없이 `SetGlobalPermission(Allow)`을 호출하여 모든 기능을 사용자에게 알리지 않고 허용했습니다. 이제 `Permissions`에 항목이 하나라도 있으면 이러한 일괄 허용을 설정하지 **않습니다**. 설정되지 않은 기능은 자동으로 허용되는 대신 WebView2의 기본 권한 요청 메시지로 넘어갑니다.

즉, Windows에서 `Permissions`을 하나라도 구성하면 명시적으로 나열하지 않은 모든 기능은 사용자에게 알리지 않고 허용되는 대신 권한 요청 메시지를 표시합니다. 필요한 기능을 명시적으로 설정하십시오.

### macOS (TCC)

macOS는 시스템 개인정보 보호 프레임워크를 통해 카메라, 마이크, 위치 정보 및 알림 접근을 관리합니다. 웹 콘텐츠에서 기능을 처음 요청하면 OS 메시지가 자동으로 표시되며, 사용자의 선택은 시스템 설정 → 개인정보 보호 및 보안에서 애플리케이션별로 기억됩니다.

이는 `Permissions`을 구성하지 않아도 올바르게 작동합니다. 현재 macOS에서는 이 맵이 **무시됩니다**. 무엇을 설정하든 모든 요청이 TCC를 거칩니다. 실질적인 제약은 macOS에서 `PermissionDeny`이 아무 효과도 없다는 것입니다. TCC가 시스템 수준에서 이미 허용한 기능을 웹뷰에서 사용하지 못하도록 차단할 수 없습니다.

`Info.plist`에 적절한 사용 목적 설명 키가 포함되어 있는지 확인하십시오.

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## 일반적인 패턴

### 미디어 캡처 애플리케이션

모든 플랫폼에서 카메라와 마이크를 허용합니다:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

<strong>Linux</strong>에서는 두 장치가 모두 명시적으로 허용되며, 다른 기능은 계속 거부됩니다. <strong>Windows</strong>에서는 두 기능이 모두 허용되며, 목록에 지정하지 않은 다른 기능에는 네이티브 프롬프트가 표시됩니다. <strong>macOS</strong>에서는 아무 효과가 없으며, 모든 권한은 TCC에서 처리합니다.

### Linux에서 미디어 캡처 거부

Linux에서는 기본적으로 카메라와 마이크가 허용됩니다. 이를 허용하지 않으려면 다음과 같이 설정하세요.

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Windows에서 모든 기능 허용

프롬프트를 표시하지 않고 모든 기능을 허용하려면 다음과 같이 설정하세요.

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### 창별 정책

창마다 서로 다른 정책을 적용할 수 있습니다.

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Windows 전용 재정의

창별 `Windows.Permissions` 필드(`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`)도 계속 작동하며, 크로스 플랫폼 맵이 적용된 후 개별 기능의 설정을 재정의할 수 있습니다. 크로스 플랫폼에서 이에 대응하는 기능이 없는 WebView2 권한 종류(예: `CoreWebView2PermissionKindOtherSensors`)에 접근해야 할 때 사용하세요.

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows에서의 평가 순서는 다음과 같습니다.

1. 크로스 플랫폼 `Permissions` 맵(`SetPermission`을 통해 종류별 상태 설정)
2. `Windows.Permissions` 맵(개별 종류 재정의)
3. 두 맵 어디에도 포함되지 않은 종류: 정책이 구성된 경우 WebView2의 네이티브 프롬프트를 표시하고, 정책이 구성되지 않은 경우 자동으로 허용합니다(기존 동작).

## 플랫폼 지원 표

| 기능 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 마이크 | ✅ | ✅ | TCC에서만 처리 |
| 카메라 | ✅ | ✅ | TCC에서만 처리 |
| 위치 정보 | ❌ 아직 지원되지 않음 | ✅ | TCC에서만 처리 |
| 알림 | ❌ 아직 지원되지 않음 | ✅ | TCC에서만 처리 |
| 클립보드 읽기 | ❌ 아직 지원되지 않음 | ✅ | TCC에서만 처리 |

## 문제 해결

**`getUserMedia`가 업그레이드 후에도 Linux에서 계속 실패함**

`PermissionMicrophone: PermissionDeny` 또는 `PermissionCamera: PermissionDeny`을 명시적으로 설정하지 않았는지 확인하세요. 기본값(설정되지 않음)은 Linux에서 미디어 캡처를 허용합니다.

**Windows에서 구성하지 않은 권한을 묻는 프롬프트가 표시됨**

`Permissions`에 항목이 하나라도 있으면 Wails는 더 이상 일괄 `Allow` 허용을 설정하지 않습니다. 목록에 지정하지 않은 기능에는 WebView2의 네이티브 프롬프트가 표시됩니다. 앱에서 사용하는 모든 기능에 명시적인 `PermissionAllow` 항목을 추가하세요.

**macOS 권한이 작동하지 않음**

`Permissions` 맵은 macOS에서 아무 효과가 없습니다. `Info.plist`에 올바른 사용 목적 설명 키(`NSMicrophoneUsageDescription`, `NSCameraUsageDescription` 등)가 포함되어 있는지, 그리고 사용자가 시스템 설정 → 개인정보 보호 및 보안에서 접근을 허용했는지 확인하세요.

**Linux에서 위치 정보/알림/클립보드 설정이 적용되지 않음**

현재 Linux에서는 카메라와 마이크만 처리됩니다. 다른 기능 유형에 대한 지원은 아직 구현되지 않았으므로 설정한 정책과 관계없이 계속 거부됩니다.
