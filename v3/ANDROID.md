# Wails v3 on Android

Wails v3 apps run on Android as native applications: an Android `WebView`
renders the frontend, assets are served **in-process** through a
`WebViewAssetLoader` backed by the Go asset server (no localhost server, no
open ports), and the standard Wails runtime (`@wailsio/runtime`) works
unchanged — bindings, events, dialogs and clipboard all route through the Go
message processor.

The same `main.go` builds for desktop and Android. Android-specific behaviour
is configured through `application.Options.Android` and per-platform files
guarded by `//go:build android`. The Go code is compiled as a C shared
library (`libwails.so`, `-buildmode=c-shared`, `GOOS=android` + the NDK
toolchain) and loaded by a small Java host (`MainActivity` + `WailsBridge`).

## Status

| Area | Status |
|---|---|
| WebView + in-process asset serving (`WebViewAssetLoader`) | ✅ Working |
| Service bindings (JS → Go calls) | ✅ Working (JavascriptInterface transport) |
| Events (Go → JS and JS → Go) | ✅ Working |
| Message dialogs (Info/Question/Warning/Error) | ✅ AlertDialog, with button callbacks |
| Open file / files dialogs | ✅ Storage Access Framework (files imported as cache copies) |
| Open directory dialogs | ❌ Returns an error — SAF yields tree URIs, not filesystem paths |
| Save file dialogs | ❌ Returns an error — write inside the app sandbox instead |
| Clipboard | ✅ ClipboardManager |
| Screens API | ✅ WindowMetrics/DisplayMetrics (dp, pixels, scale, system-bar work area) |
| App lifecycle events (`events.Android.*`, `Common.ApplicationStarted`) | ✅ Working |
| Haptics (vibrate), device info, toast | ✅ `Android.*` runtime API |
| System events (battery, network, theme, screen-lock, low-memory) | ✅ `events.Android.*` → `events.Common.*` application events |
| Native mobile features (share, torch, biometrics, geolocation, sensors, …) | ✅ Exported `application.Android*` functions — see [Native mobile features](#native-mobile-features) |
| Window geometry (SetSize/SetPosition/Minimize/...) | Intentional no-ops (apps are fullscreen) |
| Menus, system tray | Intentional no-ops |
| Multiple windows | ⚠️ Only the first window is displayed |
| Emulator + device builds | ✅ emulator: `wails3 task android:run` / `android:deploy-emulator`; physical device: `android:run:device` / `android:deploy-device` |
| Release signing | ✅ Debug keystore by default; real keystore via env vars |
| Play Store submission pipeline | ⚠️ Sign a release APK/AAB with your own keystore (see below) |

## Requirements

- The **Android SDK** with platform-tools, an SDK platform (API 35),
  build-tools and the **NDK** (r26+/26.3.x). `wails3 doctor` reports what it
  can see.
- A **JDK** (e.g. OpenJDK 21) for Gradle. Set `JAVA_HOME` if `java` is not on
  your `PATH`.
- Go 1.24+, npm.
- `ANDROID_HOME` (or `ANDROID_SDK_ROOT`) pointing at the SDK; optionally
  `ANDROID_NDK_HOME` (otherwise the newest installed NDK is used).

Install the SDK pieces with the command-line tools:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Quickstart (Emulator)

From a Wails v3 project:

```bash
wails3 task android:run
```

This boots an emulator if none is running, generates the TypeScript bindings,
builds the frontend, compiles the Go code to `libwails.so` for the host
architecture (the emulator's ABI), assembles a debug APK with Gradle, then
installs and launches it. Stream logs with `wails3 task android:logs`. In
debug builds the WebView is inspectable from Chrome at `chrome://inspect`.

`wails3 task android:package` produces a production release APK
(`-tags production,android`, stripped, framework diagnostics compiled out).
`wails3 task android:run` and `wails3 task android:deploy-emulator` are
emulator-oriented tasks; they use `ensure-emulator` and do not select a
physical device.

## Device & release builds

```bash
# Debug build/install in the Android Emulator
wails3 task android:run

# Debug build/install on the first connected physical Android device (arm64)
wails3 task android:run:device

# Production release APK (signed with the debug keystore by default)
wails3 task android:package

# Production release APK installed and launched on the first connected
# physical Android device (arm64)
wails3 task android:deploy-device

# Target a specific physical device serial
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device

# Production release APK signed with your own keystore
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

- `wails3 task android:package:fat` builds both `arm64-v8a` and `x86_64`
  into one APK (useful for distributing a single artifact that runs on
  devices and emulators).
- `wails3 task android:run:device` builds a debug APK for `arm64`, selects the
  first connected non-emulator device from `adb devices`, installs it with that
  device serial, and launches `com.wails.app.MainActivity`.
- `wails3 task android:deploy-device` does the same for a production APK.
  If no physical device is connected, pass `DEVICE_ID=<serial>` to either task;
  find serials with `adb devices`.
- `wails3 task android:deploy-emulator` installs and launches a production APK
  in the Android Emulator after ensuring an emulator is running.
- Without keystore env vars, release builds are signed with the Android
  **debug** keystore so they install for testing. **Play Store uploads
  require your own keystore** (set the `ANDROID_KEYSTORE_*` variables, or
  open `build/android/` in Android Studio and use *Build → Generate Signed
  Bundle / APK* to produce an `.aab`).

## Configuration

`build/config.yml`:

```yaml
# Android options are read by the build tasks; APP_ID controls the package name
APP_ID: com.example.myapp
```

App-level options (`application.Options.Android`) are a placeholder today; the
Android surface is driven mostly from the frontend through the `Android`
runtime object: `Android.Haptics.Vibrate(durationMs)`,
`Android.Device.Info()`, `Android.Toast.Show(message)`.

## Native mobile features

Android exposes a set of "genuinely mobile" capabilities as exported
`application.Android*` functions (guarded by `//go:build android`), each
forwarded to a matching method on the Java `WailsBridge` via JNI. They mirror
the `application.IOS*` surface, so one event-driven layer drives both platforms
(see the `mobile` example's `registerNativeFeatures`).

For the subset of capabilities whose signature is identical on both platforms,
`application.Mobile` provides one build-guarded entry point: it dispatches to
`Android` on Android, `IOS` on iOS, and a no-op stub on desktop — so
cross-platform code can call e.g. `application.Mobile.StoragePath()` without its
own `//go:build` split. Platform-specific calls stay on `Android` / `IOS`.

| Capability | API | Notes |
|---|---|---|
| Share sheet | `Android.Share(json)` | `Intent.ACTION_SEND` |
| Open URL externally | `Android.OpenURL(url)` | `Intent.ACTION_VIEW` |
| Keep screen awake | `Android.SetKeepAwake(bool)` | `FLAG_KEEP_SCREEN_ON` |
| Torch / flashlight | `Android.SetTorch(bool)` | `CameraManager` → `common:torch` |
| Safe-area insets | `Android.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Brightness | `Android.SetBrightness(0-100)` / `Android.BrightnessJSON()` | |
| App info | `Android.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Orientation lock | `Android.SetOrientation(...)` / `Android.OrientationJSON()` | |
| Status bar | `Android.SetStatusBar(json)` | style + visibility |
| Biometrics | `Android.BiometricAuthenticate(reason)` | `BiometricPrompt` → `common:biometric` |
| Local notification | `Android.Notify(json)` | `NotificationManager` (POST_NOTIFICATIONS) |
| Secure storage | `Android.SecureSet/Get/Delete` | `EncryptedSharedPreferences` |
| Haptics | `Android.Haptic(type)` | `VibrationEffect` |
| Geolocation | `Android.GetLocation()` | one-shot → `common:location` (ACCESS_FINE_LOCATION) |
| Accelerometer | `Android.SetMotion(bool)` | `SensorManager` stream → `common:motion` |
| Proximity | `Android.SetProximity(bool)` | → `common:proximity` |
| Text-to-speech | `Android.Speak(text)` / `Android.StopSpeak()` | `TextToSpeech` |
| Storage info | `Android.StorageJSON()` | `{free,total}` bytes (`StatFs`) |
| Storage path | `Android.StoragePath()` | `getFilesDir()` (for databases & persistent files); `""` if unavailable |
| Power / battery | `Android.PowerJSON()` | `{level,charging,lowPower}` |
| Network status | `Android.NetworkJSON()` | `{connected,type}` (`ConnectivityManager`) |
| Keyboard insets | `Android.SetKeyboardWatch(bool)` | → `common:keyboard {visible,height}` |
| Screen-capture | `Android.SetScreenProtect(bool)` | `FLAG_SECURE` → `common:screenCapture` |

Asynchronous results flow back to the frontend through the bridge's
`nativeEmitEvent` JNI export → `globalApp.Event.Emit`. Geolocation, biometrics
and notifications need their permissions in `AndroidManifest.xml`
(`ACCESS_FINE_LOCATION`, `USE_BIOMETRIC`, `POST_NOTIFICATIONS`); location and
notifications are requested at runtime on first use.

## System events

OS signals surface as typed application events: `events.Android.BatteryChanged`,
`NetworkChanged`, `ThemeChanged`, `ScreenLocked`, `ScreenUnlocked` and
`ApplicationLowMemory`, each also mapped to a platform-neutral `events.Common.*`
(`BatteryChanged`, `NetworkChanged`, `ThemeChanged`, `ScreenLocked`,
`ScreenUnlocked`, `LowMemory`). Subscribe with
`app.Event.OnApplicationEvent(events.Common.BatteryChanged, …)` and read the
payload (battery level, network type, dark-mode flag, …) from the event
context. The `mobile` example forwards these to its frontend as `sys:*` events.

## Porting an existing desktop app

- Everything compiles unchanged under `GOOS=android`; Android is fullscreen,
  so window-geometry calls, menus and the system tray become no-ops.
- Save-file and choose-directory dialogs return an error on Android — write
  into the app sandbox (the app's files/cache directory) and share via an
  intent instead. Open-file dialogs work and import the chosen documents as
  cache-directory copies so you get real filesystem paths.
- **`android` implies the `linux` build tag** (Android is a Linux kernel):
  desktop-Linux-only files need `//go:build linux && !android`, and
  `runtime.GOOS == "android"` (not `"linux"`) at runtime. This is the Android
  analogue of `ios`/`darwin`.
- A real app is always built with `CGO_ENABLED=1` and the NDK (the JNI
  bridge needs cgo). The non-cgo path exists only so tooling such as
  `wails3 generate bindings` can load the package.
- Design the frontend responsively; the WebView fills the display and the
  `Screens` work area excludes the status and navigation bars.

## Architecture notes

- `MainActivity` creates the `WebView`, wires a `WebViewAssetLoader` to the
  Go asset server, and exposes `window.wails` (a JavascriptInterface). It
  calls `WailsBridge.initialize()`, which loads `libwails.so` and calls
  `nativeInit`; Go then starts `main()` on a goroutine and blocks the Android
  lifecycle in `platformRun`.
- Assets and runtime calls flow through JNI: the WebView's
  `shouldInterceptRequest` → `WailsBridge.serveAsset` → Go asset server.
  Because the Android WebView cannot deliver `fetch()` POST bodies to
  `shouldInterceptRequest`, runtime calls use a dedicated transport — the
  bundled runtime detects `window.wails.invokeAsync` and routes calls through
  it to `nativeHandleRuntimeCall`, with responses delivered back via
  `window._wailsAndroidCallback`.
- Go → JS uses `WebView.evaluateJavascript` on the main looper; lifecycle
  callbacks (`onResume`/`onPause`/...) become `events.Android.*` events,
  with `ActivityCreated` mapped to `Common.ApplicationStarted`.
- Native facilities (dialogs, clipboard, screen/device info, toast, vibrate,
  main-thread dispatch) are methods on `WailsBridge` called from Go over JNI.
- Framework diagnostics are compiled out of production builds
  (`-tags production`); debug builds log to logcat under the `Wails` tag.

## Layout of the generated Android project

`build/android/` is a standard Gradle project:

```
build/android/
  app/
    build.gradle                     # SDK levels, ABIs, signing config
    src/main/
      AndroidManifest.xml
      java/com/wails/app/
        MainActivity.java            # WebView host + lifecycle + file picker
        WailsBridge.java             # JNI bridge + native facilities
        WailsJSBridge.java           # window.wails JavascriptInterface
        WailsPathHandler.java        # asset loader → Go
      jniLibs/<abi>/libwails.so      # compiled Go (produced by the build)
  gradlew, settings.gradle, ...
```
