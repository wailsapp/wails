---
title: "모바일 개요"
description: "데스크톱 앱과 동일한 Go 코드베이스로 iOS 및 Android 앱 빌드하기"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3는 데스크톱용으로 이미 작성한 것과 동일한 `main.go` 및 프런트엔드를 사용하여 **iOS와 Android**에서 실행됩니다. 별도의 모바일 프로젝트나 코드 공유 브리지가 필요하지 않으며 코드를 다시 작성할 필요도 없습니다. Go 바이너리는 모바일 대상을 위해 컴파일되고, 네이티브 WebView가 기존 프런트엔드를 렌더링합니다.

@cards{cols="2"}
iOS
WKWebView + UIKit 호스트. 사용자 지정 `wails://` 스킴을 통해 애셋을 제공하므로 포트를 열지 않습니다. 전체 Xcode가 설치된 **macOS**가 필요합니다.

[iOS 가이드 →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. Go는 NDK를 통해 `libwails.so`로 컴파일됩니다. macOS, Linux 및 Windows에서 작동합니다.

[Android 가이드 →](/guides/mobile/android/)

@end

## 실행 모습 보기: Kitchen Sink 예제

무엇이 가능한지 이해하는 가장 좋은 방법은 하나의 코드베이스로 iOS, Android 및 데스크톱에서 동일하게 실행되는 단일 Wails 앱인 **Kitchen Sink**를 살펴보는 것입니다.

@linkcard{title="모바일 Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="바인딩 · 이벤트 · 대화 상자 · 햅틱 · 위치 정보 · 생체 인식 · 알림 · 보안 저장소 · 그 밖의 기능 — 모두 하나의 main.go에서 제공"}
이 예제는 7개의 탭에서 모든 주요 모바일 API 영역을 보여 주며 데스크톱에서도 실행됩니다. 프런트엔드의 플랫폼 검사를 통해 데스크톱에서는 **Mobile** 및 **Hardware** 탭을 숨깁니다. 데스크톱용으로 빌드할 때 Go 측에서는 `common:*` 모바일 이벤트의 핸들러를 등록하지 않습니다. 하나의 코드베이스를 모든 플랫폼에 배포할 때 권장되는 패턴입니다.

| 탭 | 플랫폼 | 제공 기능 |
| --- | --- | --- |
| **바인딩** | 모두 | 값, 구조체 및 오류를 반환하는 JS → Go 서비스 호출 |
| **이벤트** | 모두 | Go → JS 시계, JS → Go → JS 핑퐁, OS 시스템 이벤트(배터리, 네트워크, 테마) |
| **대화 상자** | 모두 | 각 플랫폼의 네이티브 메시지 대화 상자 |
| **시스템** | 모두 | 클립보드, 화면 메트릭, 기기 정보 |
| **모바일** | iOS + Android | 공유 시트, 화면 켜짐 유지, 손전등, 밝기, 생체 인식, 로컬 알림, 보안 저장소 |
| **하드웨어** | iOS + Android | 햅틱, 위치 정보, 가속도계, 근접 센서, 텍스트 음성 변환 |
| **네이티브** | iOS + Android | iOS: 햅틱 + WKWebView 토글 · Android: 진동 + 토스트 |

직접 실행하려면 다음과 같이 하세요.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 작동 방식

모든 플랫폼에 동일한 애플리케이션 모델이 적용됩니다.

1. **Go 백엔드** — 서비스, 이벤트 핸들러 및 애플리케이션 로직을 변경하지 않고 `GOOS=ios` 및 `GOOS=android`용으로 컴파일합니다.
2. **프런트엔드** — 완전히 동일한 HTML/JS/CSS를 사용합니다. `@wailsio/runtime` 패키지도 동일하게 작동하며 서비스 바인딩, 이벤트, 대화 상자 및 클립보드는 모두 동일한 프로세스 내 전송 계층을 통해 라우팅됩니다.
3. **WebView 호스트** — iOS에서는 `UIViewController` 내부의 `WKWebView`를, Android에서는 `Activity` 내부의 `WebView`를 사용합니다. Wails가 메시지 브리지를 자동으로 연결합니다.
4. **프로세스 내 애셋 제공** — localhost 서버가 아니라 Go 메모리에서 애셋을 직접 제공합니다. 열린 포트나 루프백이 필요하지 않으며 추가 지연 시간도 없습니다.

플랫폼별 동작은 `//go:build ios` 또는 `//go:build android`로 보호된 파일에 두어 공유 코드를 깔끔하게 유지합니다.

## 필수 조건 한눈에 보기

| 요구 사항 | iOS | Android |
| --- | --- | --- |
| 운영 체제 | macOS만 지원 | macOS, Linux, Windows |
| 도구 체인 | 전체 Xcode(CLI 도구만으로는 안 됨) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 확인 명령 | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
도구 체인을 설정한 후 `wails3 doctor`을 실행하세요. 각 플랫폼에서 발견된 항목과 누락된 항목을 정확히 보여 줍니다.

@end

## 지원되는 기능

두 플랫폼은 동일한 핵심 기능 세트를 공유합니다.

| 기능 | iOS | Android |
| --- | --- | --- |
| 서비스 바인딩(JS → Go) | ✅ | ✅ |
| 이벤트(양방향) | ✅ | ✅ |
| 메시지 대화 상자 | ✅ UIAlertController | ✅ AlertDialog |
| 파일 열기 대화 상자 | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| 파일 저장 대화 상자 | ❌ 대신 샌드박스에 쓰기 | ❌ 대신 샌드박스에 쓰기 |
| 클립보드 | ✅ UIPasteboard | ✅ ClipboardManager |
| 화면/안전 영역 메트릭 | ✅ | ✅ |
| 수명 주기 이벤트 | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| 햅틱 | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| 기기 정보 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| 네이티브 탭(iOS) | ✅ UITabBar | — |
| 토스트 메시지(Android) | — | ✅ `Android.Toast.Show` |
| 다중 창 | ❌ 첫 번째 창만 지원 | ❌ 첫 번째 창만 지원 |
| 창 크기 및 위치 / 메뉴 / 트레이 | 의도적으로 아무 작업도 하지 않음 | 의도적으로 아무 작업도 하지 않음 |

## 빌드 태그 규칙

플랫폼 조건부 코드를 작성할 때 알아야 할 두 가지 중요한 규칙은 다음과 같습니다:

- <strong>`ios`는 `darwin`</strong>를 내포하므로, `//go:build darwin` 태그가 지정된 파일은 iOS에서도 컴파일됩니다. macOS만 대상으로 하려면 `//go:build darwin && !ios`를 사용하세요.
- <strong>`android`는 `linux`</strong>를 내포하므로, `//go:build linux` 태그가 지정된 파일은 Android에서도 컴파일됩니다. 데스크톱 Linux만 대상으로 하려면 `//go:build linux && !android`를 사용하세요.

런타임에 `runtime.GOOS`는 각각 `"ios"`와 `"android"`를 반환합니다.

## 런타임 플랫폼 감지

빌드 태그는 특정 플랫폼에서만 *컴파일*할 수 있는 코드에 사용합니다. 공유 코드에서 일반적인 조건 분기를 수행하려면 `application.System`를 사용하세요. 모든 빌드에서 사용할 수 있으므로(빌드 태그 불필요) 동일한 파일이 모든 플랫폼에서 작동합니다:

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

사용 가능: `IsMobile()`, `IsDesktop()`, `IsServer()`(`server` 빌드 태그), `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

프런트엔드에서는 `@wailsio/runtime`에 대응하는 헬퍼를 제공합니다:

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 다음 단계

@cards{cols="2"}
🚀 첫 번째 모바일 앱
데스크톱 Wails 앱을 몇 분 만에 iOS Simulator 또는 Android Emulator에서 실행해 보세요.

[시작하기 →](/guides/mobile/first-mobile-app/)

---
iOS 가이드
iOS 도구 체인의 전체 설정, 시뮬레이터, 기기 빌드, 서명, 구성 및 API 레퍼런스를 안내합니다.

[iOS 가이드 →](/guides/mobile/ios/)

---
Android 가이드
Android SDK/NDK의 전체 설정, 에뮬레이터, APK 서명, Play Store 패키징 및 API 레퍼런스를 안내합니다.

[Android 가이드 →](/guides/mobile/android/)

@end
