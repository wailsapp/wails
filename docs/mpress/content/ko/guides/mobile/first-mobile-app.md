---
title: "첫 번째 모바일 앱"
description: "몇 분 만에 Wails 앱을 iOS Simulator 또는 Android Emulator에서 실행하기"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

이 가이드에서는 표준 Wails 데스크톱 앱을 iOS Simulator 또는 Android Emulator에서 실행합니다. **Go 코드는 변경할 필요가 없습니다.** 동일한 `main.go`에서 모든 대상을 빌드할 수 있습니다.

**예상 소요 시간:** 15~30분(대부분은 첫 실행 시 도구 체인을 설치하는 데 걸리는 시간입니다)

## 데스크톱 프로젝트에서 시작하기

아직 프로젝트가 없다면 새 프로젝트를 만드세요:

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

먼저 데스크톱 앱이 제대로 작동하는지 확인하세요:

```bash
wails3 dev
```

앱이 열리면 종료한 후 다음 단계로 진행하세요. 데스크톱에서 실행되는 모든 기능은 모바일에서도 실행되므로, 이 가이드에서는 `main.go`이나 Go 코드를 전혀 수정할 필요가 없습니다.

---

## 플랫폼 선택하기

@tabs{sync-key="mobile-platform"}
[iOS Simulator]
### 요구 사항

- **macOS**(iOS 빌드는 macOS에서만 지원됨)
- **전체 Xcode** — 명령줄 도구만으로는 충분하지 않습니다. App Store에서 설치한 후 다음을 실행하세요:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.25 이상** 및 **npm**(`wails3 init`을 실행했다면 이미 설치되어 있음)

`wails3 doctor`을 실행하여 확인하세요. 이 명령은 검색된 iOS SDK를 나열합니다.

### Simulator에서 실행하기

@steps
### 앱 실행하기
```bash
wails3 task ios:run
```

이것으로 끝입니다. 앱을 빌드하고, 실행 중인 Simulator가 없으면 부팅한 다음 앱을 실행합니다.

@note{type="tip"}
첫 실행에는 몇 분이 걸립니다(iOS용 Wails 프레임워크를 컴파일하고 캐시하기 때문입니다). 이후 실행부터는 훨씬 빨라집니다.

@end

실행되면 수정하지 않은 데스크톱 앱이 iOS Simulator에서 작동합니다. `main.go`도 같고 프런트엔드도 같습니다:

![iOS Simulator에서 실행 중인 기본 Wails 앱](/assets/ios-simulator-first-app.png)

### 로그 실시간 보기
별도의 터미널에서 다음을 실행하세요:

```bash
wails3 task ios:logs:dev
```

이 명령은 앱에 해당하는 항목만 필터링하여 Simulator 로그를 실시간으로 표시합니다. `fmt.Println` 및 `log.Println` 출력이 여기에 나타납니다.

### WebView 검사하기
Safari에서 <strong>개발자용 → Simulator → 해당 앱</strong>으로 이동하세요. 콘솔, 디버거, 네트워크 패널 등 Web Inspector의 모든 기능을 사용할 수 있습니다.

### 변경 사항 적용하기
프런트엔드 파일(`frontend/src/main.js`, `index.html` 등)을 수정하고 `wails3 task ios:run`을 다시 실행하세요. Wails가 프런트엔드를 다시 빌드하고 앱을 재실행합니다.

Go 코드를 변경한 경우에도 `wails3 task ios:run`을 다시 실행하세요. Go 재컴파일은 증분 방식이므로 변경된 패키지만 다시 빌드됩니다.

@end

### Xcode에서 열기(선택 사항)

```bash
wails3 task ios:xcode
```

이 명령은 `build/ios/`을 Xcode에서 엽니다. Xcode를 사용하여 기기에 배포하거나, 고급 프로파일링을 수행하거나, 프로비저닝 프로파일을 관리할 수 있습니다. Wails는 빌드할 때마다 Xcode 프로젝트를 다시 생성하므로 생성된 파일을 직접 수정하지 마세요.

[Android Emulator]
### 요구 사항

**Android SDK**, **NDK** 및 <strong>JDK</strong>가 필요합니다. 가장 간편한 방법은 Android Studio를 사용하는 것이며, 명령줄 도구를 사용할 수도 있습니다:

@steps
### Android 명령줄 도구 설치하기
[developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only)에서 다운로드한 후 `~/android-sdk/cmdline-tools/latest/`에 압축을 푸세요.

### SDK 구성 요소 설치하기
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Emulator 만들기
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### 환경 변수 설정하기
`~/.zshrc` 또는 `~/.bashrc`에 다음 내용을 추가하세요:

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

다음 명령으로 다시 불러오세요: `source ~/.zshrc`

### JDK 설치하기
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

`wails3 doctor`을 실행하여 필요한 항목이 모두 검색되는지 확인하세요.

### Emulator에서 실행하기

@steps
### 앱 실행하기
```bash
wails3 task android:run
```

처음 실행하면 다음 작업을 수행합니다:

- 실행 중인 Emulator가 없으면 부팅합니다
- 바인딩을 생성하고 프런트엔드를 빌드합니다
- NDK 크로스 컴파일러를 사용하여 Go 코드를 `libwails.so`로 컴파일합니다
- Gradle로 디버그 APK를 빌드합니다
- Emulator에 APK를 설치하고 실행합니다

@note{type="tip"}
첫 빌드에서는 Gradle을 다운로드하고 NDK 도구 체인을 컴파일하므로 5~10분 정도 걸립니다. 이후 빌드는 증분 방식으로 진행되며 1분 이내에 완료됩니다.

@end

### 로그 실시간 보기
별도의 터미널에서 다음을 실행하세요:

```bash
wails3 task android:logs
```

이 명령은 `adb logcat`을 실행하고 앱에 해당하는 항목만 필터링합니다. `fmt.Println` 출력이 여기에 나타납니다.

### WebView 검사하기
Chrome을 열고 `chrome://inspect`으로 이동하세요. 앱의 WebView가 **원격 대상** 아래에 표시됩니다. <strong>검사</strong>를 클릭하여 DevTools를 여세요.

### 변경 사항 적용하기
파일을 수정하고 `wails3 task android:run`을 다시 실행하세요. Gradle의 증분 빌드 기능으로 변경된 코드만 다시 컴파일됩니다.

@end

@end

---

## 수행된 작업 이해하기

`main.go`은 전혀 변경되지 않았습니다. 모든 작업은 Wails가 처리했습니다:

- **빌드 시스템** — 프로젝트의 `Taskfile.yml`에는 플랫폼별 도구 체인을 구동하는 `ios:*` 및 `android:*` 작업이 포함되어 있습니다.
- **Go 크로스 컴파일** — 적절한 `GOARCH` 및 sysroot와 함께 `GOOS=ios` 또는 `GOOS=android`을 사용합니다.
- **네이티브 호스트** — 컴파일된 Go 코드를 포함하고 WebView를 호스팅하는, 생성된 Xcode 프로젝트(iOS) 또는 Gradle 프로젝트(Android)입니다.
- **애셋 제공** — `frontend/dist/`은 Go 바이너리에 포함되며 프로세스 내부에서 제공됩니다. localhost 서버는 필요하지 않습니다.

---

## 앱이 모바일 환경을 인식하도록 만들기

앱은 이미 작동하지만 휴대전화 화면에서 데스크톱 앱처럼 보입니다. 몇 가지 작은 변경만으로도 큰 차이를 만들 수 있습니다.

### 반응형 CSS

모바일 화면은 더 좁고 입력 방식도 다릅니다. `frontend/public/style.css`(또는 이에 해당하는 파일)에서 다음과 같이 설정하세요.

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Go에서 플랫폼 감지하기

공유 코드를 복잡하게 만들지 않고 플랫폼별 동작을 추가하려면 빌드 태그를 사용하세요.

iOS 전용 코드용 `mobile_ios.go`을 만드세요.

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Android 전용 코드용 `mobile_android.go`을 만드세요.

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

공유 코드가 데스크톱에서도 컴파일되도록 스텁인 `mobile_desktop.go`을 만드세요.

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### JavaScript에서 플랫폼을 감지하고 모바일 전용 UI 제한하기

`IOS.*` 및 `Android.*` 런타임 객체는 각각 해당 플랫폼에만 존재합니다. 데스크톱에서 이 객체를 호출하면 예외가 발생합니다. Kitchen Sink에서 사용하는 올바른 패턴은 플랫폼을 한 번 감지한 후 모바일 전용 컨트롤을 완전히 숨기는 것입니다.

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

그런 다음 HTML에서 다음과 같이 작성하세요.

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

이렇게 하면 모바일 전용 버튼이 데스크톱에서 전혀 렌더링되지 않으며, 개별 호출마다 `if (isMobile)` 검사를 추가할 필요도 없습니다.

Go 측에서는 빌드 태그가 지정된 스텁과 함께 사용하여, 필요한 플랫폼에서만 이벤트 핸들러가 등록되도록 하세요.

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

이는 [Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)에서 사용하는 것과 정확히 같은 패턴입니다. `native_features_stub.go`, `native_features_ios.go` 및 `native_features_android.go`을 참조하세요.

@note{type="note" title="네이티브 기능 API 및 이벤트 명명 규칙"}
알아 두어야 할 규칙은 두 가지입니다.

- **Go 측 네이티브 기능은 플랫폼 관리자를 사용합니다.** `application.IOS.*` 및 `application.Android.*` 싱글턴을 통해 호출하세요. 예를 들면 `application.IOS.Haptic("medium")` 또는 `application.Android.Share(payload)`입니다. 각 관리자는 자체 플랫폼에만 존재하므로 관련 호출은 `//go:build ios` / `//go:build android` 파일에 둡니다.
- **이벤트의 네임스페이스는 적용 범위에 따라 구분됩니다.** 두 플랫폼 모두에서 처리할 수 있는 이벤트에는 `common:*` 접두사를 사용하고(`common:haptic`, `common:location`, …), 한 플랫폼에서만 생성하거나 처리할 수 있는 이벤트에는 `ios:*` 또는 `android:*`을 사용합니다(예: `ios:backgroundTask`, `android:foregroundService`). 거의 모든 모바일 기능이 공유되므로 프런트엔드에서는 `common:*` 아래에서 이벤트마다 리스너 하나만 유지하면 됩니다.

@end

### 햅틱 피드백 추가하기(iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### 진동 추가하기(Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## 프로덕션용으로 빌드하기

@tabs{sync-key="mobile-platform"}
[iOS]
**시뮬레이터 빌드**(시뮬레이터 테스트용이며 서명이 필요하지 않음):

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**기기 빌드**(서명 ID와 프로비저닝 프로파일이 필요함):

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**배포용 IPA**(App Store 또는 TestFlight용):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
App Store Connect에 업로드할 때는 `wails3 task ios:xcode`을 사용하고 Xcode에서 서명과 아카이브를 관리하도록 하세요. Xcode가 인증서, 프로파일 및 공증의 복잡한 과정을 자동으로 처리합니다.

@end

[Android]
**디버그 APK**(Android 디버그 키 저장소로 서명되며 직접 설치 가능):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**릴리스 APK**(자체 키 저장소로 서명):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**유니버설 APK**(단일 파일에 arm64 + x86_64 포함):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Play Store에 업로드하려면 APK 대신 `.aab`(Android App Bundle)를 생성하세요. Android Studio에서 `build/android/`을 열고 <strong>Build → Generate Signed Bundle / APK</strong>를 사용하세요.

@end

@end

---

## 문제 해결

### `wails3 task ios:run` 실행 시 "no iOS SDKs found" 오류가 발생하는 경우

전체 Xcode가 설치되어 있고 선택되어 있어야 합니다.

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` 실행 시 "SDK not found" 오류가 발생하는 경우

`ANDROID_HOME`이 설정되고 내보내졌는지 확인하세요. 다음 명령으로 검증할 수 있습니다.

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### 시뮬레이터가 부팅되지 않는 경우

사용 가능한 시뮬레이터를 나열하고 하나를 수동으로 부팅하세요.

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect`에 대상이 표시되지 않는 경우

WebView가 디버그 모드여야 합니다(`android:run`의 기본값). 프로덕션 빌드가 아닌 디버그 빌드를 실행하고 있는지 확인하세요. 또한 `adb devices`에 에뮬레이터가 연결된 것으로 표시되는지도 확인하세요.

#### 안전 영역 인셋이 적용되지 않는 경우

HTML에 viewport 메타 태그가 포함되어 있는지 확인하세요:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Kitchen Sink 살펴보기

첫 번째 앱이 실행되면, 그 밖에 어떤 기능을 구현할 수 있는지 가장 빠르게 알아보는 방법은 **Kitchen Sink** 예제를 살펴보는 것입니다. 이 예제는 단일 코드베이스로 iOS, Android 및 데스크톱에서 실행되는 완전한 Wails 앱이며, 햅틱 피드백, 위치 정보, 생체 인증, 로컬 알림, 보안 저장소 등을 다룹니다:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

[`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)에서 소스 코드를 살펴보세요. 특히 `native_features_ios.go` 및 `native_features_android.go` 파일은 플랫폼별 기능을 복사하여 활용하기 위한 시작점으로 유용합니다.

## 다음 단계

@cards{cols="2"}
iOS 가이드
전체 참조 문서: 구성 옵션, 네이티브 탭, WKWebView 전환 옵션, 기기용 빌드, 서명.

[iOS 가이드 →](/guides/mobile/ios/)

---
Android 가이드
전체 참조 문서: 구성, 토스트 알림, Play Store 패키징, NDK 세부 정보.

[Android 가이드 →](/guides/mobile/android/)

---
📖 Kitchen Sink 소스 코드
햅틱 피드백, 위치 정보, 생체 인증, 알림, 보안 저장소를 하나의 실행 가능한 앱에서 모두 확인할 수 있습니다.

[GitHub에서 보기 →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
