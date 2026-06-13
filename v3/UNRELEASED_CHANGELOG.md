# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->
- iOS: native message dialogs (UIAlertController) and open file/files/directory dialogs (UIDocumentPickerViewController); save dialogs return an explicit error
- iOS: clipboard support via UIPasteboard
- iOS: real screen metrics via UIScreen (points, pixels, scale, safe-area work area)
- iOS: device builds (`IOS_PLATFORM=device`), code-signing identity / provisioning profile / entitlements support, `.ipa` packaging, and `deploy-device` via devicectl
- iOS: configurable minimum iOS version (`ios.minIOSVersion` in build/config.yml)
- iOS: `wails3 doctor` reports Xcode and iOS SDK availability on macOS
- iOS: system events — battery, network, theme, screen-lock and low-memory surface as `events.IOS.*` and platform-neutral `events.Common.*` application events
- iOS: native mobile feature bridge (exported `application.IOS*`) — share sheet, open URL, keep-awake, torch, safe-area insets, brightness, app info, orientation lock, status bar, biometrics (Face ID/Touch ID), local notifications and Keychain secure storage
- iOS: sensors & hardware — haptics, one-shot geolocation, accelerometer, proximity, text-to-speech, storage info, power/battery state, network status, keyboard insets and screen-capture detection
- iOS: documentation (IOS.md and a docs-site guide)
- Android: native message dialogs (AlertDialog) and open file/files dialogs (Storage Access Framework, imported as cache copies); open-directory and save dialogs return an explicit error
- Android: clipboard support via ClipboardManager
- Android: real screen metrics via WindowMetrics/DisplayMetrics (dp, pixels, scale, system-bar work area)
- Android: haptics (`Android.Haptics.Vibrate`), device info (`Android.Device.Info`) and toast (`Android.Toast.Show`) runtime methods
- Android: typed lifecycle events (`events.Android.*`, generated from events.txt) with `ActivityCreated` mapped to `Common.ApplicationStarted`
- Android: build pipeline produces installable debug and release APKs (`android:run`, `android:package`, `android:package:fat`); release signing via the debug keystore by default or a real keystore through `ANDROID_KEYSTORE_*` env vars
- Android: `wails3 doctor` reports the Android SDK, NDK and JDK
- Android: system events — battery, network, theme, screen-lock and low-memory surface as `events.Android.*` and platform-neutral `events.Common.*` application events
- Android: native mobile feature bridge (exported `application.Android*`) — share, open URL, keep-awake, torch, safe-area insets, brightness, app info, orientation lock, status bar, biometrics (BiometricPrompt), local notifications and EncryptedSharedPreferences secure storage
- Android: sensors & hardware — haptics, one-shot geolocation, accelerometer, proximity, text-to-speech, storage info, power/battery state, network status, keyboard insets and FLAG_SECURE screen-capture blocking
- Android: documentation (ANDROID.md and a docs-site guide)
- Example: the `mobile` kitchen sink gains Mobile and Hardware tabs demonstrating the native feature bridge across iOS and Android (pill tabs wrap to multiple rows)
- Mobile: battery — the accelerometer, proximity sensor, torch and the example's periodic clock are paused when the app is backgrounded and restored on return (Android keeps the process running in the background, and the torch is hardware state that persists on iOS), and Android system-event receivers are only registered while the app is in the foreground
- iOS: camera capture — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → a `native:capture` event with a base64 thumbnail)
- iOS: background execution — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (a UIApplication background-task window) and a configurable `ios.backgroundModes` (build/config.yml) that templates `UIBackgroundModes` into the generated Info.plist
- Android: camera capture — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (system camera via FileProvider → a `native:capture` event)
- Android: foreground service — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (a `WailsForegroundService` with an ongoing notification keeps the process alive for long-running background work)
- Example: a Camera tab demonstrating photo/video capture and background execution (foreground service on Android, background-task window on iOS)

## Changed
<!-- Changes in existing functionality -->
- Replace `github.com/go-git/go-git/v5` direct dependency with calls to the system `git` CLI (`internal/git` package). **Note: `git` must be installed on the system.** Graceful errors are returned when `git` is not found in `PATH`.


## Fixed
<!-- Bug fixes -->
- Fix `getUserMedia` always failing with `NotAllowedError` on Linux: WebKitGTK denies permission requests nobody handles, and the `permission-request` signal was not connected. Camera/microphone are now handled per a new cross-platform `WebviewWindowOptions.Permissions` map (`map[PermissionType]Permission`), honored on both Linux (WebKitGTK) and Windows (WebView2). On Linux, which has no native prompt, camera/microphone default to allowed (restoring `getUserMedia`) and can be turned off with `PermissionDeny` (#5552)
- iOS: `GOOS=ios` compiles again (exported `events.IOS`, mobile method-name stubs) and production-tagged builds compile (build-tag fixes in pkg/application and several services)
- iOS: Go→JS events and ExecJS now work — the page no longer loads twice at startup and the `wails:runtime:ready` handshake can no longer be lost
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` no longer race app startup; removed the fixed 2-second startup sleep
- iOS: fixed a C-string leak on every Go→JS JavaScript execution
- iOS: `hasListeners` now reflects real listener registration
- iOS: framework debug logging is compiled out of production builds
- Android: `GOOS=android` compiles again — defined `events.Android`, removed the out-of-bounds `events_android.go` listener array, added the mobile method-name stub, and stopped desktop-Linux files (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) leaking into Android builds
- Android: JS→Go bindings now work — the WebView cannot deliver `fetch()` POST bodies to `shouldInterceptRequest`, so runtime calls route through a JavascriptInterface transport (`nativeHandleRuntimeCall`) instead of crashing on a nil request body
- Android: `Screens.*` runtime calls return real data — the ScreenManager is now populated at startup (it was never wired, so `GetAll` returned nil)
- Android: framework debug logging is compiled out of production builds and routes through logcat under the `Wails` tag in debug builds
- Android: real `hasListeners` registry, JNI reference/exception handling, and a single-load page lifecycle (no double navigation)
- Fix `wails3 generate bindings` failing with "Access is denied" on Windows when the Vite dev server is running, by syncing generated files into the output directory instead of renaming over it (#5515)
- Fix intermittent fatal crash on macOS when reading screen information after a display change: the screen id and name stored pointers to autoreleased `UTF8String` buffers that could be freed before Go copied them (use-after-free). The strings are now `strdup`'d and freed after conversion, and screen enumeration runs in an explicit autorelease pool so it no longer leaks when called from Go goroutines (#5556)
- Fix intermittent SIGSEGV on Linux when the assetserver closes a `WebKitURISchemeRequest`: the final `g_object_unref` ran on the assetserver goroutine, finalizing a WebKit GObject off the GTK main thread. The unref is now marshalled onto the GTK main context via `g_main_context_invoke` (#5557)

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new `SetWindowIcon()` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
