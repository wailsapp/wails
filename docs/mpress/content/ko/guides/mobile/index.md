---
title: "모바일 개요"
description: "데스크톱 앱과 동일한 Go 코드베이스로 iOS 및 Android 앱 빌드"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3는 데스크톱용으로 이미 작성한 것과 동일한 `main.go` 및 프런트엔드를 사용하여 <strong>iOS와 Android</strong>에서 실행됩니다. 별도의 모바일 프로젝트도, 코드 공유 브리지도, 재작성도 필요하지 않습니다. Go 바이너리가 모바일 대상으로 컴파일되고 네이티브 WebView가 기존 프런트엔드를 렌더링합니다.

@cards{cols="2"}
iOS
WKWebView + UIKit 호스트. 사용자 지정 `wails://` 스킴을 통해 애셋을 제공하므로 열린 포트가 없습니다. 전체 Xcode가 설치된 <strong>macOS</strong>가 필요합니다.

[iOS 가이드 →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. Go는 NDK를 통해 `libwails.so`(으)로 컴파일됩니다. macOS, Linux 및 Windows에서 작동합니다.

[Android 가이드 →](/guides/mobile/android/)

@end

## 실행 모습 살펴보기: Kitchen Sink 예제

무엇이 가능한지 이해하는 가장 좋은 방법은 하나의 코드베이스로 iOS, Android 및 데스크톱에서 동일하게 실행되는 단일 Wails 앱인 <strong>Kitchen Sink</strong>를 살펴보는 것입니다.

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="바인딩 · 이벤트 · 대화 상자 · 햅틱 · 위치 정보 · 생체 인증 · 알림 · 보안 저장소 등 — 모두 하나의 main.go에서 제공"}
7개의 탭에서 모든 주요 모바일 API 영역을 보여 주며 데스크톱에서도 실행됩니다. 프런트엔드의 플랫폼 검사를 통해 데스크톱에서는 **모바일** 및 **하드웨어** 탭을 숨깁니다. 데스크톱용으로 빌드할 때 Go 측에서는 `common:*` 모바일 이벤트의 핸들러를 등록하지 않습니다. 하나의 코드베이스를 모든 플랫폼에 배포할 때 권장하는 패턴입니다.

| 탭 | 플랫폼 | 표시하는 기능 |
| --- | --- | --- |
| **바인딩** | 전체 | 값, 구조체 및 오류를 반환하는 JS → Go 서비스 호출 |
| **이벤트** | 전체 | Go → JS 시계, JS → Go → JS 핑퐁, OS 시스템 이벤트(배터리, 네트워크, 테마) |
| **대화 상자** | 전체 | 각 플랫폼의 네이티브 메시지 대화 상자 |
| **시스템** | 전체 | 클립보드, 화면 메트릭, 기기 정보 |
| **모바일** | iOS + Android | 공유 시트, 절전 방지, 손전등, 밝기, 생체 인증, 로컬 알림, 보안 저장소 |
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

1. **Go 백엔드** — 서비스, 이벤트 핸들러 및 애플리케이션 로직이 변경 없이 `GOOS=ios` 및 `GOOS=android`용으로 컴파일됩니다.
2. **프런트엔드** — 완전히 동일한 HTML/JS/CSS입니다. `@wailsio/runtime` 패키지도 동일하게 작동하며 서비스 바인딩, 이벤트, 대화 상자 및 클립보드가 모두 동일한 프로세스 내 전송 계층을 통해 라우팅됩니다.
3. **WebView 호스트** — iOS에서는 `UIViewController` 내부의 `WKWebView`이고, Android에서는 `Activity` 내부의 `WebView`입니다. Wails가 메시지 브리지를 자동으로 연결합니다.
4. **프로세스 내 애셋 제공** — 애셋은 localhost 서버가 아니라 Go 메모리에서 직접 제공됩니다. 열린 포트도, 루프백도, 추가 지연 시간도 없습니다.

플랫폼별 동작은 `//go:build ios` 또는 `//go:build android`의 보호를 받는 파일에 두므로 공유 코드를 깔끔하게 유지할 수 있습니다.

## 한눈에 보는 필수 구성 요소

| 요구 사항 | iOS | Android |
| --- | --- | --- |
| 운영 체제 | macOS만 지원 | macOS, Linux, Windows |
| 도구 체인 | 전체 Xcode(CLI 도구만으로는 안 됨) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 확인 명령 | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
도구 체인을 설정한 후 `wails3 doctor`을 실행하세요. 각 플랫폼에서 감지된 항목과 누락된 항목을 정확히 보여 줍니다.

@end

## 지원되는 기능

두 플랫폼은 동일한 핵심 기능 세트를 공유합니다.

| 기능 | iOS | Android |
| --- | --- | --- |
| 서비스 바인딩(JS → Go) | ✅ | ✅ |
| 이벤트(양방향) | ✅ | ✅ |
| 메시지 대화상자 | ✅ UIAlertController | ✅ AlertDialog |
| 파일 열기 대화상자 | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| 파일 저장 대화상자 | ❌ 대신 샌드박스에 쓰기 | ❌ 대신 샌드박스에 쓰기 |
| 클립보드 | ✅ UIPasteboard | ✅ ClipboardManager |
| 화면 / 안전 영역 측정값 | ✅ | ✅ |
| 수명 주기 이벤트 | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| 햅틱 피드백 | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| 기기 정보 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| 네이티브 탭(iOS) | ✅ UITabBar | — |
| 토스트 메시지(Android) | — | ✅ `Android.Toast.Show` |
| 다중 창 | ❌ 첫 번째 창만 지원 | ❌ 첫 번째 창만 지원 |
| 창 크기 및 위치 / 메뉴 / 트레이 | 의도적으로 아무 작업도 수행하지 않음 | 의도적으로 아무 작업도 수행하지 않음 |

## 빌드 태그 규칙

플랫폼 조건부 코드를 작성할 때 알아야 할 두 가지 중요한 규칙은 다음과 같습니다.

- <strong>`ios`는 `darwin`</strong>를 포함합니다. 따라서 `//go:build darwin` 태그가 지정된 파일은 iOS용으로도 컴파일됩니다. macOS만 대상으로 하려면 `//go:build darwin && !ios`을 사용하세요.
- <strong>`android`는 `linux`</strong>를 포함합니다. 따라서 `//go:build linux` 태그가 지정된 파일은 Android용으로도 컴파일됩니다. 데스크톱 Linux만 대상으로 하려면 `//go:build linux && !android`을 사용하세요.

런타임에 `runtime.GOOS`는 각각 `"ios"`와 `"android"`를 반환합니다.

## 런타임 플랫폼 감지

빌드 태그는 지정된 플랫폼에서만 <em>컴파일</em>할 수 있는 코드에 사용합니다. 공유 코드에서 일반적인 분기 처리를 하려면 `application.System`을 사용하세요. 이 기능은 모든 빌드에서 사용할 수 있으므로(빌드 태그 불필요) 같은 파일이 모든 플랫폼에서 작동합니다.

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

사용 가능 항목: `IsMobile()`, `IsDesktop()`, `IsServer()`(`server` 빌드 태그), `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

프런트엔드의 `@wailsio/runtime`에는 이에 대응하는 헬퍼가 있습니다.

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 다음 단계

@cards{cols="2"}
🚀 첫 번째 모바일 앱
데스크톱 Wails 앱을 단 몇 분 만에 iOS Simulator 또는 Android Emulator에서 실행하세요.

[시작하기 →](/guides/mobile/first-mobile-app/)

---
iOS 가이드
전체 iOS 도구 체인 설정, 시뮬레이터, 기기 빌드, 서명, 구성 및 API 레퍼런스입니다.

[iOS 가이드 →](/guides/mobile/ios/)

---
Android 가이드
전체 Android SDK/NDK 설정, 에뮬레이터, APK 서명, Play Store 패키징 및 API 레퍼런스입니다.

[Android 가이드 →](/guides/mobile/android/)

@end
