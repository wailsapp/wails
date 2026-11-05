---
title: "Android"
description: "Android에서 Wails 애플리케이션 빌드 및 실행 — 툴체인 설정, 에뮬레이터, APK 서명, Play Store 패키징 및 API 참조"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="실험적 기능"}
Android 지원은 실험적 기능이며 향후 릴리스에서 변경될 수 있습니다.

@end

@note{type="tip"}
Wails로 모바일 개발을 처음 시작하시나요? 단계별 안내인 [첫 번째 모바일 앱 →](/guides/mobile/first-mobile-app/)부터 살펴본 후 전체 참조 문서를 확인하려면 여기로 돌아오세요.

@end

Wails v3 애플리케이션은 Android에서 네이티브 앱으로 실행됩니다. Android `WebView`가 프런트엔드를 렌더링하고, 애셋은 Go 애셋 서버가 지원하는 `WebViewAssetLoader`를 통해 **프로세스 내에서** 제공됩니다(로컬 호스트 서버와 열린 포트 없음). 또한 표준 `@wailsio/runtime`는 변경 없이 작동하며 서비스 바인딩, 이벤트, 대화 상자 및 클립보드는 Go 메시지 프로세서를 통해 처리됩니다.

동일한 `main.go`로 데스크톱 및 Android용 빌드를 생성합니다. Go 코드는 C 공유 라이브러리(`libwails.so`, `GOOS=android` + NDK 툴체인)로 컴파일되어 작은 Java 호스트에서 로드됩니다. Android 전용 동작은 `//go:build android`로 보호되는 플랫폼별 Go 파일에 구현되어 있습니다.

## 요구 사항

- platform-tools, SDK 플랫폼(API 35), build-tools 및 **NDK**(26.3.x)가 포함된 **Android SDK** — `wails3 doctor`에서 감지된 항목을 확인할 수 있습니다
- Gradle용 **JDK**(예: OpenJDK 21). `java`가 `PATH`에 없으면 `JAVA_HOME`를 설정하세요
- Go 1.26 이상 및 npm
- SDK를 가리키는 `ANDROID_HOME`(또는 `ANDROID_SDK_ROOT`)

명령줄 도구를 사용하여 SDK 구성 요소를 설치하세요:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## 에뮬레이터에서 실행

프로젝트 디렉터리에서 다음을 실행하세요:

```bash
wails3 task android:run
```

실행 중인 에뮬레이터가 없으면 에뮬레이터를 부팅하고, 바인딩을 생성하며, 프런트엔드를 빌드한 다음 Go 코드를 에뮬레이터의 ABI에 맞는 `libwails.so`로 컴파일합니다. 이어서 Gradle로 디버그 APK를 어셈블하고 설치한 후 실행합니다.

함께 사용하면 유용한 명령:

```bash
wails3 task android:logs    # stream the app's logcat output
```

디버그 빌드에서는 Chrome의 `chrome://inspect`에서 WebView를 검사할 수 있습니다.

## 패키징

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

프로덕션 빌드는 `-tags production,android`을 사용하고 심벌을 제거하며 프레임워크의 내부 진단 기능을 컴파일에서 제외합니다. `wails3 task android:package:fat`은 `arm64-v8a`와 `x86_64`를 모두 하나의 APK로 빌드합니다.

Google Play에 새 앱을 제출하려면 Android App Bundle(`.aab`) 형식을 사용해야 하며, 새 제출 앱은 Android 15(API 35) 이상을 대상으로 해야 합니다. 프로젝트 템플릿은 `build/android/app/build.gradle`에서 `compileSdk`와 `targetSdk`를 35로 설정합니다. `wails3 task android:bundle:fat`은 두 ABI를 모두 포함한 `bin/<AppName>.aab`을 생성합니다. Google Play는 이 번들로부터 기기별 최적화 APK를 생성하므로, 모든 ABI를 포함한 번들이 스토어 업로드에 적합한 아티팩트입니다. `.aab`은 `adb`로 직접 설치할 수 없으므로 로컬 및 에뮬레이터 테스트에는 여전히 APK가 가장 빠릅니다.

`android:run`과 `android:deploy-emulator`은 에뮬레이터용 작업입니다. 실제 Android 기기에서는 디버그 APK의 경우 `android:run:device`를, 릴리스 APK의 경우 `android:deploy-device`을 사용하세요. 두 작업 모두 `arm64`용으로 빌드하고, `adb devices`에서 연결된 비에뮬레이터 항목 중 첫 번째 항목을 선택하여 설치한 후 `com.wails.app.MainActivity`을 실행합니다. 특정 기기를 대상으로 하려면 `DEVICE_ID=<serial>`을 전달하세요.

## 서명 및 릴리스 빌드

키 저장소가 없으면 릴리스 빌드는 테스트용으로 설치할 수 있도록 Android **디버그** 키 저장소로 서명됩니다. 자체 키 저장소로 서명하려면 다음을 설정하세요:

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

동일한 변수가 App Bundle에도 서명합니다. 이 변수들을 설정한 상태에서 `wails3 task android:bundle:fat`을 실행하면 Play에 제출할 수 있는 `.aab`이 생성됩니다. 변수를 설정하지 않으면 번들이 디버그 키 저장소로 서명되어 Google Play에서 거부되므로 작업에서 경고를 출력합니다.

@note{type="tip"}
[Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756)을 사용하면 로컬에서 서명할 때 사용하는 키 저장소가 **업로드 키**입니다. Google은 이 키로 업로드를 확인한 다음 관리 중인 앱 서명 키로 앱에 다시 서명합니다. 또한 Google Play는 업로드할 때마다 더 높은 `versionCode`를 요구하므로 `build/android/app/build.gradle`에서 값을 올리세요.

@end

## 구성

프런트엔드는 런타임에 `Android` 런타임 객체를 통해 Android 기능(`Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`)을 제어합니다. 패키지 이름은 빌드 작업의 `APP_ID`에서 제어합니다.

## 지원되는 기능과 지원되지 않는 기능

| 영역 | 상태 |
| --- | --- |
| WebView + 프로세스 내 애셋(`WebViewAssetLoader`) | ✅ |
| 서비스 바인딩, 이벤트(양방향) | ✅ |
| 메시지 대화 상자 | ✅ 버튼 콜백을 지원하는 AlertDialog |
| 파일 열기/여러 파일 열기 대화 상자 | ✅ Storage Access Framework(파일을 캐시 복사본으로 가져옴) |
| 디렉터리 열기/파일 저장 대화 상자 | ❌ 오류를 반환합니다. 대신 앱 샌드박스 내부에 쓰세요 |
| 클립보드 | ✅ ClipboardManager |
| Screens API | ✅ 시스템 표시줄 작업 영역을 포함한 WindowMetrics |
| 수명 주기 이벤트(`events.Android.*`) | ✅ |
| 햅틱, 기기 정보, 토스트 | ✅ `Android.*` 런타임 API |
| 에뮬레이터 + 실제 기기 빌드 | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| 창 위치 및 크기, 메뉴, 시스템 트레이 | 의도적으로 아무 작업도 하지 않음 |
| 여러 창 | 첫 번째 창만 표시됨 |

## 포팅 참고 사항

- 데스크톱 코드는 `GOOS=android`에서 변경 없이 컴파일됩니다. Android 앱은 전체 화면으로 실행되므로 창 위치 및 크기, 메뉴, 트레이 호출은 아무 작업도 하지 않습니다.
- `android` **에는 `linux` 빌드 태그**가 암시적으로 적용됩니다(Android는 Linux 커널을 사용함). 데스크톱 Linux 전용 파일에는 `//go:build linux && !android`가 필요하며, 런타임에서 `runtime.GOOS`는 `"android"`입니다.
- 파일 저장 및 디렉터리 선택 대화 상자를 앱 샌드박스에 쓰고 intent 공유 흐름을 사용하는 방식으로 대체하세요. 파일 열기 대화 상자는 작동하며 선택한 문서를 캐시 디렉터리의 복사본으로 가져오므로 실제 파일 시스템 경로를 얻을 수 있습니다.
- 실제 앱은 항상 `CGO_ENABLED=1` 및 NDK로 빌드됩니다. non-cgo 경로는 `wails3 generate bindings` 같은 도구가 패키지를 로드할 수 있도록 하기 위해서만 존재합니다.
- 프런트엔드를 반응형으로 설계하세요. `Screens` 작업 영역에는 상태 표시줄과 탐색 표시줄이 포함되지 않습니다.
