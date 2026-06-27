# Wails v3 on iOS

Wails v3 apps run on iOS as native UIKit applications: a `WKWebView` renders
the frontend, assets are served **in-process** over a custom `wails://` URL
scheme (no localhost server, no open ports), and the standard Wails runtime
(`@wailsio/runtime`) works unchanged — bindings, events, dialogs and
clipboard all route through the same `/wails/runtime` transport as on
desktop.

The same `main.go` builds for desktop and iOS. iOS-specific behaviour is
configured through `application.Options.IOS` and per-platform files guarded
by `//go:build ios`.

## Status

| Area | Status |
|---|---|
| WKWebView + in-process asset serving (`wails://`) | ✅ Working |
| Service bindings (JS → Go calls) | ✅ Working |
| Events (Go → JS and JS → Go) | ✅ Working |
| Message dialogs (Info/Question/Warning/Error) | ✅ UIAlertController, with button callbacks |
| Open file / files / directory dialogs | ✅ UIDocumentPickerViewController (files are imported as sandbox copies) |
| Save file dialogs | ❌ Returns an error — write inside the app sandbox instead |
| Clipboard | ✅ UIPasteboard |
| Screens API | ✅ UIScreen (points, pixels, scale, safe-area work area) |
| App lifecycle events (`events.IOS.*`, `Common.ApplicationStarted`) | ✅ Working |
| Haptics, device info, native UITabBar, scroll/bounce/UA options | ✅ `IOSOptions` + `IOS.*` runtime API |
| System events (battery, network, theme, screen-lock, low-memory) | ✅ `events.IOS.*` → `events.Common.*` application events |
| Native mobile features (share, torch, biometrics, geolocation, sensors, …) | ✅ Exported `application.IOS*` functions — see [Native mobile features](#native-mobile-features) |
| Window geometry (SetSize/SetPosition/Minimize/...) | Intentional no-ops (apps are fullscreen) |
| Menus, system tray | Intentional no-ops |
| Multiple windows | ⚠️ Only the first window is displayed |
| Simulator builds | ✅ `wails3 task ios:run` / `ios:package` |
| Device builds | ✅ `ios:package IOS_PLATFORM=device` (manual signing) or the generated Xcode project (automatic signing) |
| App Store submission pipeline | ⚠️ Use the generated Xcode project (`ios:xcode`) for archive/upload |

## Requirements

- macOS with **full Xcode** installed (the CLI tools alone are not enough);
  `wails3 doctor` reports the iOS SDKs it can see
- Go 1.24+, npm

## Quickstart (Simulator)

From a Wails v3 project:

```bash
wails3 task ios:run
```

This generates the Go build overlay and Xcode scaffolding (`build/ios/`),
builds the Go code as a C archive (`GOOS=ios`), links it with the UIKit
bootstrap, bundles, ad-hoc signs, boots a simulator if needed, and launches
the app. Use `wails3 task ios:logs:dev` to stream its logs. In debug builds
the webview is inspectable from Safari's Develop menu.

`wails3 task ios:package` produces a production `.app`
(`-tags production,ios`, stripped, framework diagnostics compiled out);
`ios:deploy-simulator` installs and launches it.

## Device builds

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision
wails3 task ios:deploy-device [DEVICE_ID=<udid>]      # xcrun devicectl
wails3 task ios:package:ipa IOS_PLATFORM=device ...   # distribution .ipa
```

- `IOS_PLATFORM=device` switches to the `iphoneos` SDK and the
  `arm64-apple-ios<min>` target.
- Entitlements: `build/ios/entitlements.plist` is applied to device builds
  only (`get-task-allow` by default; add capability keys as needed).
- For **automatically managed** signing, provisioning and App Store
  archives, open the generated project instead: `wails3 task ios:xcode`.

## Configuration

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"   # templated into Info.plist and used by the toolchain
```

App-level options (`application.Options.IOS`): `DisableScroll`,
`DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`,
`EnableBackForwardNavigationGestures`, `DisableLinkPreview`,
`EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`,
`DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`,
`BackgroundColour`, `EnableNativeTabs` + `NativeTabsItems`.

The frontend can drive iOS features at runtime through the `IOS` runtime
object: `IOS.Haptics.Impact(style)`, `IOS.Device.Info()`,
`IOS.Scroll.SetEnabled(...)`, etc. Native tab selections arrive as a
`nativeTabSelected` `CustomEvent` on `window`.

## Native mobile features

Beyond the cross-platform runtime, iOS exposes a set of "genuinely mobile"
capabilities as exported `application.IOS*` functions (guarded by
`//go:build ios`, implemented by a small Objective-C bridge). Each has an
`application.Android*` counterpart, so a single event-driven layer can drive
both platforms (see the `mobile` example's `registerNativeFeatures`).

For the subset of capabilities whose signature is identical on both platforms,
`application.Mobile` provides one build-guarded entry point: it dispatches to
`IOS` on iOS, `Android` on Android, and a no-op stub on desktop — so
cross-platform code can call e.g. `application.Mobile.StoragePath()` without its
own `//go:build` split. Platform-specific calls stay on `IOS` / `Android`.

| Capability | API | Notes |
|---|---|---|
| Share sheet | `IOS.Share(json)` | `UIActivityViewController` |
| Open URL externally | `IOS.OpenURL(url)` | Opens in Safari |
| Keep screen awake | `IOS.SetKeepAwake(bool)` | Idle-timer toggle |
| Torch / flashlight | `IOS.SetTorch(bool)` | → `common:torch` event |
| Safe-area insets | `IOS.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Brightness | `IOS.SetBrightness(0-1)` / `IOS.GetBrightness()` | |
| App info | `IOS.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Orientation lock | `IOS.SetOrientation("portrait\|landscape\|auto")` / `IOS.GetOrientation()` | |
| Status bar | `IOS.SetStatusBar(json)` | style + visibility |
| Biometrics | `IOS.BiometricAuthenticate(reason)` | Face ID / Touch ID, passcode fallback → `common:biometric` |
| Local notification | `IOS.PostNotification(json)` | `UNUserNotificationCenter` |
| Secure storage | `IOS.SecureSet/Get/Delete` | Keychain |
| Haptics | `IOS.Haptic(type)` | impact / notification / selection |
| Geolocation | `IOS.GetLocation()` | one-shot → `common:location` (needs `NSLocationWhenInUseUsageDescription`) |
| Accelerometer | `IOS.SetMotion(bool)` | Core Motion stream → `common:motion` (needs `NSMotionUsageDescription`) |
| Proximity | `IOS.SetProximity(bool)` | → `common:proximity` |
| Text-to-speech | `IOS.Speak(text)` / `IOS.StopSpeak()` | `AVSpeechSynthesizer` |
| Storage info | `IOS.StorageJSON()` | `{free,total}` bytes |
| Storage path | `IOS.StoragePath()` | Application Support dir, created on first access (for databases & persistent files); `""` if it can't be created |
| Power / battery | `IOS.PowerJSON()` | `{level,charging,lowPower}` |
| Network status | `IOS.NetworkJSON()` | `{connected,type}` |
| Keyboard insets | `IOS.SetKeyboardWatch(bool)` | → `common:keyboard {visible,height}` |
| Screen-capture | `IOS.SetScreenProtect(bool)` | Detects screenshots & recording (iOS can't block them) → `common:screenCapture` |

Asynchronous results are delivered to the frontend as custom events
(`iosEmitNativeEvent` → `globalApplication.Event.Emit`). Geolocation and motion
require the matching `NS*UsageDescription` keys in `Info.plist`; the linker
pulls in `CoreLocation`, `CoreMotion` and `SystemConfiguration` (already wired
in `build/ios/Taskfile.yml`).

## System events

OS signals surface as typed application events: `events.IOS.BatteryChanged`,
`NetworkChanged`, `ThemeChanged`, `ScreenLocked`, `ScreenUnlocked` and
`ApplicationDidReceiveMemoryWarning`, each also mapped to a platform-neutral
`events.Common.*` (`BatteryChanged`, `NetworkChanged`, `ThemeChanged`,
`ScreenLocked`, `ScreenUnlocked`, `LowMemory`). Subscribe with
`app.Event.OnApplicationEvent(events.Common.BatteryChanged, …)` and read the
payload (battery level, network type, dark-mode flag, …) from the event
context. The `mobile` example forwards these to its frontend as `sys:*` events.

## Porting an existing desktop app

- Everything compiles unchanged under `GOOS=ios`; iOS is fullscreen, so
  window-geometry calls, menus and the system tray become no-ops.
- Save-file dialogs return an error on iOS — write into the app sandbox
  (e.g. `os.UserHomeDir()` + `Documents`) and offer a share flow instead.
- Note that `ios` implies the `darwin` build tag: desktop-mac-only files
  need `//go:build darwin && !ios`, and `runtime.GOOS == "ios"` (not
  `"darwin"`) at runtime.
- Design the frontend responsively (safe areas are handled natively; the
  webview is laid out inside them).

## Architecture notes

- `main.m` starts the Go runtime on a background queue and runs
  `UIApplicationMain` on the main thread; the Go side waits for the app
  delegate's launch signal, then drives window/webview creation on the main
  queue.
- `WKURLSchemeHandler` forwards `wails://` requests to the Go asset server;
  the first `/wails/runtime` call also marks the window's JS runtime ready
  (belt-and-braces for the `wails:runtime:ready` handshake).
- Go → JS uses `evaluateJavaScript` on the main queue; JS → Go uses
  `window.webkit.messageHandlers.external.postMessage` and the HTTP-style
  runtime transport over the custom scheme.
- Framework diagnostics are compiled out of production builds
  (`-tags production`); debug builds log through the unified log and the
  webview console.
