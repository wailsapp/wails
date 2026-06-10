# Wails Mobile Kitchen Sink

A single Wails v3 app that runs on **iOS, Android and desktop** from one
`main.go` and one frontend, exercising the cross-platform runtime surface:

| Tab | Demonstrates |
|---|---|
| **Bindings** | JS → Go service calls returning values, structs and an error |
| **Events** | Go → JS (`time` clock) and JS → Go → JS (`ping`/`pong`) |
| **Dialogs** | Native message dialogs (Info / Warning / Error / Question + callback) |
| **System** | Clipboard round-trip, screen metrics, device info |
| **Native** | Platform-specific: iOS haptics + WKWebView toggles, Android vibrate + toast |

The UI feature-detects the platform (`window.wails` on Android, the WKWebView
message handler on iOS) and shows only the controls that platform supports.

## Run it

```bash
# iOS Simulator (requires full Xcode)
wails3 task ios:run

# Android emulator (requires the Android SDK + NDK + a JDK)
wails3 task android:run

# Desktop
wails3 task run
```

`wails3 task ios:package` / `android:package` produce release builds. See
[`../../IOS.md`](../../IOS.md) and [`../../ANDROID.md`](../../ANDROID.md) for the
toolchain requirements and device/signing details.

## How it works

- `main.go` is shared across all platforms. On Android the app is built as a
  c-shared library, so `main_android.go` registers `main` via
  `application.RegisterAndroidMain`; on iOS the generated build overlay invokes
  it; on desktop it runs directly.
- `SystemService` (`greetservice.go`) is the bound Go service.
- The frontend imports `@wailsio/runtime` and calls `Runtime.IOS.*` /
  `Runtime.Android.*` for the platform-specific features.
