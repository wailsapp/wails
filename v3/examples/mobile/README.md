# Wails Mobile Kitchen Sink

A single Wails v3 app that runs on **iOS, Android and desktop** from one
`main.go` and one frontend, exercising the cross-platform runtime surface:

| Tab | Demonstrates |
|---|---|
| **Bindings** | JS → Go service calls returning values, structs and an error |
| **Events** | Go → JS (`time` clock), JS → Go → JS (`ping`/`pong`), and OS **system events** (battery, network, theme, screen-lock, low-memory) |
| **Dialogs** | Native message dialogs (Info / Warning / Error / Question + callback) |
| **System** | Clipboard round-trip, screen metrics, device info |
| **Mobile** | Share sheet, open URL, keep-awake, torch, safe-area insets, brightness, app info, orientation lock, status bar, biometrics, local notifications, secure storage |
| **Hardware** | Haptics, geolocation, accelerometer, proximity, text-to-speech, storage, power/battery, network, keyboard insets, screen-capture |
| **Native** | Platform-specific: iOS haptics + WKWebView toggles, Android vibrate + toast |

The UI feature-detects the platform (`window.wails` on Android, the WKWebView
message handler on iOS) and shows only the controls that platform supports —
the **Mobile** and **Hardware** tabs appear on both iOS and Android.

## Run it

```bash
# iOS Simulator (requires full Xcode)
wails3 task ios:run

# Android emulator (requires the Android SDK + NDK + a JDK)
wails3 task android:run

# Android physical device (USB debugging enabled)
adb devices
GOWORK=off DEVICE_ID=<serial> wails3 task android:run:device
GOWORK=off DEVICE_ID=<serial> wails3 task android:deploy-device

# Desktop
wails3 task run
```

`wails3 task ios:package` / `android:package` produce release builds. See
[`../../IOS.md`](../../IOS.md) and [`../../ANDROID.md`](../../ANDROID.md) for the
toolchain requirements and device/signing details. `GOWORK=off` is only needed
when running this checked-in example from inside the Wails repository so Go uses
the example module instead of the repository workspace.

## How it works

- `main.go` is shared across all platforms. On Android the app is built as a
  c-shared library, so `main_android.go` registers `main` via
  `application.RegisterAndroidMain`; on iOS the generated build overlay invokes
  it; on desktop it runs directly.
- `SystemService` (`greetservice.go`) is the bound Go service.
- The frontend imports `@wailsio/runtime` and calls `Runtime.IOS.*` /
  `Runtime.Android.*` for the platform-specific features.
- The **Mobile** and **Hardware** tabs use an event bridge: the frontend emits
  a `common:*` event, a per-platform listener (`native_features_ios.go` /
  `native_features_android.go`, registered by `registerNativeFeatures`) calls
  the matching `application.IOS.*` / `application.Android.*` manager method, and
  asynchronous results (a biometric prompt, a GPS fix, a torch toggle) come back
  as `common:*` events. See [`../../IOS.md`](../../IOS.md) /
  [`../../ANDROID.md`](../../ANDROID.md) for the full list of these APIs.
- **System events** (Events tab) arrive in Go as `events.Common.*` application
  events (battery, network, theme, screen-lock, low-memory — mapped from the
  per-platform `ios:` / `android:` events); `main.go` forwards them to the
  frontend as `sys:*` custom events.
