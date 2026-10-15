---
title: "iOS"
description: "iOS에서 Wails 애플리케이션 빌드 및 실행 — 설정, 시뮬레이터, 기기 빌드, 구성 및 네이티브 기능"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="실험적 기능"}
iOS 지원은 실험적 기능이며 향후 릴리스에서 변경될 수 있습니다.

@end

@note{type="tip"}
Wails로 모바일 개발을 처음 시작하시나요? 단계별 안내는 [첫 번째 모바일 앱 →](/guides/mobile/first-mobile-app/)에서 확인한 후, 전체 레퍼런스를 보려면 이 페이지로 돌아오세요.

@end

Wails v3 앱은 iOS에서 완전한 네이티브 앱으로 실행되며, 무엇보다 데스크톱 버전과  *완전히 동일하게* 작동합니다. 동일한 Go 백엔드와 프런트엔드를 사용하며, 서비스 바인딩, 이벤트, 대화 상자, 클립보드 등 동일한 `@wailsio/runtime`이 모두 똑같이 작동합니다. 모바일 전용 재구성은  **전혀** 필요하지 않습니다. 별도의 모바일 코드베이스나 포팅 계층이 없고 새로 익혀야 할 특수 API도 없습니다. 기존 Wails 앱이 iOS에서 그대로 실행됩니다. 포팅은 정말 매끄럽습니다. 앱을 있는 그대로 가져와 출시하면 됩니다.

동일한 `main.go`이 데스크톱과 iOS 모두에서 빌드되며, iOS 전용 세부 설정은  `application.Options.IOS`을 통해 구성합니다.

## 요구 사항

- <strong>전체 Xcode</strong>가 설치된 macOS가 필요합니다(명령줄 도구만으로는 충분하지 않음). `wails3 doctor`에는 찾을 수 있는 iOS SDK가 표시됩니다.
- Go 1.25 이상 및 npm

## 시뮬레이터

프로젝트 디렉터리에서 다음을 실행하세요:

```bash
wails3 task ios:run
```

이 명령은 앱을 빌드하고, 실행 중인 시뮬레이터가 없으면 시뮬레이터를 부팅한 다음 앱을 실행합니다.

함께 사용하면 유용한 명령:

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

디버그 빌드에서는 Safari의 개발자용 메뉴에서 WebView를 검사할 수 있습니다.

## 패키징

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

이 빌드는 최적화되고 디버그 정보가 제거된 프로덕션 빌드입니다.

## 기기 빌드

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device`은 실제 기기용으로 빌드합니다. 엔타이틀먼트는  `build/ios/entitlements.plist`에서 가져오며 기기 빌드에만 적용됩니다. 앱에 필요한 기능 키를 추가하세요.

@note{type="tip"}
자동 관리 방식으로 서명과 프로비저닝을 처리하고 App Store 아카이브를 만들려면, 생성된 Xcode 프로젝트를 `wails3 task ios:xcode`으로 열고 Xcode에서 빌드하세요.

@end

## 구성

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

시작 옵션(`application.Options.IOS`)에는 `DisableScroll`,  `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`,  `EnableBackForwardNavigationGestures`, `DisableLinkPreview`,  `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`,  `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`,  `BackgroundColour` 및 `EnableNativeTabs` +  `NativeTabsItems`을 통한 네이티브 하단 탭이 포함됩니다.

## 네이티브 기능

iOS 전용 기능은 `application.IOS`을 통해 사용할 수 있습니다. 공유 코드를 플랫폼에 종속되지 않게 유지하려면 `//go:build ios` 파일 내의 Go 코드에서 호출하세요. Android에서도 `application.Android`을 통해 동일한 기능 세트를 제공합니다.

일회성 작업은 즉시 반환됩니다:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

쿼리 헬퍼인 `SafeAreaJSON()`, `AppInfoJSON()`,  `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`,  `GetBrightness()`은 결과를 JSON으로 반환합니다. `StoragePath()`은 앱의 Application Support 디렉터리에 대한 절대 경로를 반환합니다. 이 디렉터리는 데이터베이스와 기타 영구 파일을 저장하기에 적합합니다(Android의 `getFilesDir()`에 해당하는 iOS 디렉터리). 디렉터리는 처음 접근할 때 생성됩니다. 생성할 수 없으면 `StoragePath()`이 빈 문자열을 반환하므로, 사용하기 전에 `""`인지 확인하세요.

### 이벤트

권한 요청, 센서 스트림, 카메라 캡처처럼 나중에 완료되는 작업은 반환 값 대신 <strong>이벤트</strong>로 결과를 전달하며, Go 또는 프런트엔드에서 수신할 수 있습니다. Android와 공유하는 기능의 이름에는 `common:` 접두사가 붙고, iOS 전용 기능에는 `ios:` 접두사가 붙습니다.

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| 이벤트 | 트리거 작업 | 페이로드 |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

`v3/examples/mobile` 아래의 종합 예제는 위의 모든 기능을 처음부터 끝까지 연결합니다.

## WebView 제어

일부 WebView 동작은 런타임에 Go에서 변경할 수도 있습니다:

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

번들로 제공되는 `@wailsio/runtime`은 소규모 **프런트엔드** iOS 네임스페이스도 노출합니다:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

네이티브 하단 탭 선택은 `window`에서 `nativeTabSelected` 이벤트로 전달됩니다.

## 지원 현황

| 영역 | 상태 |
| --- | --- |
| 프런트엔드 렌더링 및 에셋 | ✅ |
| 서비스 바인딩 및 이벤트(양방향) | ✅ |
| 메시지 대화 상자 | ✅ |
| 파일/여러 파일/디렉터리 열기 대화 상자 | ✅ 샌드박스 복사본으로 가져옴 |
| 파일 저장 대화 상자 | ❌ 대신 앱 샌드박스 내부에 쓰기 |
| 클립보드 | ✅ |
| 화면 API | ✅ 안전 영역을 반영한 작업 영역 포함 |
| 수명 주기 이벤트 | ✅ |
| 창 위치 및 크기, 메뉴, 시스템 트레이 | iOS에서는 아무 작업도 수행하지 않음 |
| 다중 창 | 첫 번째 창만 표시됨 |

## 포팅 참고 사항

- 데스크톱 코드는 변경 없이 iOS용으로 빌드되며, 창, 메뉴 및 시스템 트레이 호출은 아무 작업도 수행하지 않습니다.
- 파일 저장 대화 상자를 앱 샌드박스에 쓴 후 공유하는 방식으로 대체하세요.
- 프런트엔드를 반응형으로 설계하세요. 안전 영역은 자동으로 처리됩니다.
