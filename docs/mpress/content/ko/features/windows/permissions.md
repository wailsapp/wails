---
title: "권한"
description: "웹 콘텐츠의 카메라, 마이크, 위치 정보 및 기타 기능 요청 제어"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

`navigator.mediaDevices.getUserMedia()`, Geolocation API 또는 Notifications API를 호출하는 웹 콘텐츠의 요청은 호스트 애플리케이션이 허용하거나 거부해야 합니다. Wails는 이를 선언적으로 제어할 수 있도록 `WebviewWindowOptions`에 크로스 플랫폼 `Permissions` 맵을 제공합니다. 플랫폼별 코드는 필요하지 않습니다.

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

해당 창의 웹 콘텐츠에서 보내는 카메라 및 마이크 요청은 브라우저 프롬프트 없이 허용됩니다.

## 권한 유형

`PermissionType`(uint8)은 웹 콘텐츠가 요청할 수 있는 기능을 식별합니다.

| 상수 | 기능 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 권한 값

`Permission`(uint8)은 지정된 유형에 적용되는 정책입니다.

| 상수 | 값 | 의미 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 플랫폼의 기본 처리 방식 사용(아래 참조) |
| `PermissionAllow` | 1 | 프롬프트 없이 허용 |
| `PermissionDeny` | 2 | 프롬프트 없이 거부 |

`PermissionDefault`은 0 값이므로 맵에서 설정되지 않은 항목은 기본값으로 동작합니다.

## 플랫폼별 동작

기반 웹뷰의 기본 동작이 서로 다르므로 각 플랫폼은 `PermissionDefault`을 다르게 처리합니다.

### Linux(WebKitGTK)

WebKitGTK에는 **기본 권한 프롬프트가 없습니다**. 핸들러가 연결되지 않으면 모든 요청을 아무 알림 없이 거부합니다. 이 때문에 이 기능이 추가되기 전에는 `getUserMedia`이 항상 `NotAllowedError`을 반환했습니다.

현재 Wails는 Linux에서 **카메라 및 마이크** 요청을 처리합니다. 위치 정보, 알림 및 클립보드 읽기는 아직 연결되지 않았으므로 설정한 정책과 관계없이 계속 거부됩니다.

| 정책 | 카메라/마이크 | 위치 정보, 알림, 클립보드 |
| --- | --- | --- |
| `PermissionDefault` | **허용됨**(getUserMedia 복원) | 항상 거부됨 |
| `PermissionAllow` | 허용됨 | 항상 거부됨(아직 구현되지 않음) |
| `PermissionDeny` | 거부됨 | 항상 거부됨 |

### Windows(WebView2)

WebView2에는 기본 권한 프롬프트와 유형별 권한 API가 있습니다. 다섯 가지 기능 유형이 모두 완전히 지원됩니다.

| 정책 | 동작 |
| --- | --- |
| `PermissionDefault` | WebView2가 운영 체제/브라우저의 기본 권한 프롬프트를 표시함 |
| `PermissionAllow` | 아무 알림 없이 허용됨 |
| `PermissionDeny` | 아무 알림 없이 거부됨 |

**중요:** 이 기능이 도입되기 전에는 Wails가 조건 없이 `SetGlobalPermission(Allow)`을 호출하여 모든 기능을 아무 알림 없이 허용했습니다. 이제 `Permissions`에 항목이 하나라도 있으면 이러한 일괄 허용이 설정되지 **않습니다**. 설정되지 않은 기능은 자동으로 허용되지 않고 WebView2의 기본 프롬프트로 넘어갑니다.

즉, Windows에서 `Permissions`을 하나라도 구성하면 명시적으로 나열하지 않은 기능은 아무 알림 없이 허용되는 대신 프롬프트를 표시합니다. 필요한 기능을 명시적으로 설정하세요.

### macOS(WKWebView + TCC)

macOS에서는 두 계층 모두 요청을 허용해야 합니다. WKWebView는 캡처 세션을 시작하기 전에 애플리케이션에 요청하며, `Permissions` 맵이 이 요청에 응답합니다. 그 아래에서는 시스템 개인정보 보호 프레임워크(TCC)가 기기 자체에 대한 접근을 제어합니다. 앱이 실제로 카메라나 마이크에 처음 접근할 때 운영 체제 프롬프트가 나타나며, 사용자의 선택은 시스템 설정 → 개인정보 보호 및 보안에서 앱별로 저장됩니다.

Wails는 macOS 12 이상에서 **카메라 및 마이크** 요청을 처리합니다. 위치 정보, 알림 및 클립보드 읽기에는 `WKUIDelegate`에 해당하는 항목이 없으므로 연결되지 않으며 TCC에서만 처리됩니다. Linux와 마찬가지로 이러한 기능에 설정한 정책은 아무 효과가 없습니다.

| 정책 | 카메라/마이크 | 위치 정보, 알림, 클립보드 |
| --- | --- | --- |
| `PermissionDefault` | WebKit이 자체 권한 프롬프트를 표시함 | TCC에서만 처리 |
| `PermissionAllow` | WebKit 프롬프트를 건너뜀 — **TCC는 계속 적용됨** | TCC에서만 처리 |
| `PermissionDeny` | WebKit에서 거부됨. 시스템 승인은 별도로 처리됨 | TCC에서만 처리 |

`PermissionAllow`은 기기 접근이 아니라 웹뷰의 요청을 허용합니다. 시스템 승인은 여전히 독립적으로 적용되며, 시스템 설정에서 거부된 앱은 계속 거부됩니다. WebKit 자체 프롬프트만 건너뜁니다.

macOS 12 미만에서는 해당 델리게이트 메서드가 존재하지 않으므로 맵이 무시되고 모든 요청이 WebKit 프롬프트로 대체됩니다.

`Info.plist`에 적절한 사용 목적 설명 키가 포함되어 있는지 확인하세요.

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="사용 목적 설명이 필요합니다"}

WebKit은 Wails 권한 델리게이트를 호출하기 전에 시스템 승인을 요청할 수 있습니다. `PermissionDeny`은 WebKit 요청을 차단하지만 TCC 프롬프트나 AVFoundation 승인 요청이 발생하지 않는다고 보장하지는 않습니다.

앱이 요청하는 기기에 적절한 사용 목적 설명을 포함하고, 샌드박스 환경에서는 그에 해당하는 기기 권한도 포함하세요. 이러한 선언을 `PermissionDeny`으로 대체하지 마세요.

@end

## 일반적인 패턴

### 미디어 캡처 앱

모든 플랫폼에서 카메라와 마이크를 허용합니다.

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Linux**에서는 두 장치를 모두 명시적으로 허용하며, 다른 기능은 계속 거부됩니다. **Windows**에서는 두 기능을 모두 허용하며, 목록에 지정하지 않은 다른 기능에는 네이티브 프롬프트가 표시됩니다. **macOS**에서는 WebKit 계층에서 두 기능을 모두 허용하므로 브라우저 프롬프트가 표시되지 않습니다. 하지만 처음 사용할 때는 TCC가 장치 자체에 대한 권한을 계속 요청하며, 두 사용 목적 설명 키가 모두 `Info.plist`에 있어야 합니다.

### Linux에서 미디어 캡처 거부

Linux에서는 기본적으로 카메라와 마이크를 허용합니다. 이를 허용하지 않으려면 다음과 같이 설정하세요.

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

창별 `Windows.Permissions` 필드(`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`)는 계속 작동하며, 크로스 플랫폼 맵을 적용한 후 개별 기능의 설정을 재정의할 수 있습니다. 크로스 플랫폼에 대응하는 항목이 없는 WebView2 권한 종류(예: `CoreWebView2PermissionKindOtherSensors`)에 접근해야 할 때 사용하세요.

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows에서 평가 순서는 다음과 같습니다.

1. 크로스 플랫폼 `Permissions` 맵(`SetPermission`을 통해 종류별 상태 설정)
2. `Windows.Permissions` 맵(개별 종류 재정의)
3. 두 맵 중 어느 쪽에도 포함되지 않은 종류: 정책이 구성되어 있으면 WebView2의 네이티브 프롬프트를 표시하고, 정책이 구성되어 있지 않으면 자동으로 허용합니다(기존 동작).

## 플랫폼 지원 표

| 기능 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 마이크 | ✅ | ✅ | ✅ (macOS 12 이상) |
| 카메라 | ✅ | ✅ | ✅ (macOS 12 이상) |
| 위치 정보 | ❌ 아직 지원되지 않음 | ✅ | ❌ 아직 지원되지 않음 |
| 알림 | ❌ 아직 지원되지 않음 | ✅ | ❌ 아직 지원되지 않음 |
| 클립보드 읽기 | ❌ 아직 지원되지 않음 | ✅ | ❌ 아직 지원되지 않음 |

macOS에서 ✅인 경우 정책이 WebKit의 요청에 응답하지만, 그와 별개로 TCC가 장치 접근을 계속 제어합니다. ❌인 경우 해당 기능은 TCC에서만 처리합니다.

## 문제 해결

**`getUserMedia` 업그레이드 후에도 Linux에서 계속 실패함**

`PermissionMicrophone: PermissionDeny` 또는 `PermissionCamera: PermissionDeny`을 명시적으로 설정하지 않았는지 확인하세요. 기본값(미설정)은 Linux에서 미디어 캡처를 허용합니다.

**Windows에서 구성하지 않은 권한을 요청함**

`Permissions`에 항목이 하나라도 있으면 Wails는 더 이상 포괄적인 `Allow` 허용을 설정하지 않습니다. 목록에 지정하지 않은 기능에는 WebView2의 네이티브 프롬프트가 표시됩니다. 앱에서 사용하는 모든 기능에 명시적인 `PermissionAllow` 항목을 추가하세요.

**macOS 권한이 작동하지 않음**

`Permissions`은 macOS 12 이상에서 카메라와 마이크를 지원합니다. 위치 정보, 알림, 클립보드 읽기는 아직 연동되지 않아 정책을 무시합니다. 정책이 적용되는 경우에도 그와 별개로 TCC가 장치 접근을 계속 제어합니다. `PermissionAllow`은 WebKit 프롬프트를 없애지만 시스템 프롬프트는 없애지 않습니다. `Info.plist`에 올바른 사용 목적 설명 키(`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`)가 포함되어 있는지, 그리고 사용자가 시스템 설정 → 개인정보 보호 및 보안에서 접근을 허용했는지 확인하세요.

**웹 콘텐츠에서 카메라나 마이크를 요청하면 macOS 앱이 종료됨**

앱에서 요청하는 각 장치에 대해 `Info.plist`에 `NSCameraUsageDescription` 또는 `NSMicrophoneUsageDescription`이 포함되어 있는지 확인하세요. 설명이 누락되면 요청이 거부될 수 있으며, 필수 설명 없이 AVFoundation에 접근하면 앱이 종료될 수 있습니다. `PermissionDeny`은 WebKit의 결정을 제어하며 이러한 키를 대체하지 않습니다.

**Linux에서 위치 정보/알림/클립보드 설정이 아무 효과가 없음**

현재 Linux에서는 카메라와 마이크만 처리합니다. 다른 기능 유형은 아직 지원되지 않으므로 설정한 정책과 관계없이 계속 거부됩니다.
