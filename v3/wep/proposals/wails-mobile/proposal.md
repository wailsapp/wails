# WEP: Wails v3 Mobile Support (Android)

**Author:** Ernest ([@its-ernest](https://github.com/its-ernest))  
**Created:** 2026-05-26

### Summary
This proposal introduces native Android support for Wails v3, enabling developers to build cross-platform mobile applications using the same Go-based backend logic and web-based frontend. It leverages `gomobile` for the Go-Java bridge and provides a plugin-based architecture for native platform integration.

### Motivation
Currently, Wails is primarily focused on desktop (Windows, macOS, Linux). As mobile development is a significant part of the ecosystem, providing a path for Wails developers to target mobile devices with minimal changes to their Go code is essential. The goal is to provide a "Wails-way" of handling mobile: Go for logic, Web for UI, and a native bridge designed in Java for Android, Swift for iOS, for native hardware access.

### Detailed Design
The proposed architecture follows a **Tri-Bridge** model:

1.  **Core Go Runtime (Mobile Core):**
    *   A mobile-compatible subset of the Wails v3 runtime.
    *   Uses `gomobile bind` to generate an AAR for Android.
    *   Implements an `EventBus` that supports polling from the native side to push events to JavaScript.

2.  **Native Container (Android/Java):**
    *   A custom `Activity` (e.g., `WailsWebViewActivity`) that hosts the `WebView`.
    *   Intercepts URL requests (e.g., `https://wails.local`) to serve assets directly from Go's `embed.FS`.
    *   Implements a `NativeCallHandler` interface to receive and route calls from Go.

3.  **Plugin Architecture:**
    *   Native features (Camera, Permissions, Biometrics) are implemented as "Plugins" in Java.
    *   Go orchestrates these plugins. JS calls Go -> Go calls Native Plugin -> Result returns to Go -> Go returns to JS.
    *   This ensures business logic stays in Go, while platform-specific UI or hardware logic stays in Java.

4.  **Frontend Bridge (`wails.js`):**
    *   A simplified contract providing `Wails.CallGo(method, ...args)` and `Wails.on(event, callback)`.

### Pros/Cons
**Pros:**
*   **Go-Centric Logic:** Developers write their business logic once in Go for both desktop and mobile.
*   **Modular Plugins:** Native features are decoupled, making the framework easier to maintain and extend.
*   **Performant Asset Loading:** Assets are served from memory (embedded FS) rather than a local web server.
*   **Consistency:** Follows the Wails v3 programming model (Service binding, Events).

**Cons:**
*   **Build Toolchain Complexity:** Requires Android SDK/NDK(Ultimately Android Studio) and `gomobile`(`wailsm` CLI helper exists to automate `gomobile` complexity)
*   **Platform-Specific Code:** While Go logic is shared, complex native features still require some Java/Kotlin knowledge to write new plugins.

### Alternatives Considered
*   **Local Web Server:** Running a Go-based web server on the mobile device. *Rejected due to OS restrictions, battery drain, and security concerns.*
*   **Capacitor/Cordova Approach:** Using a JS-to-Native bridge directly. *Rejected because it sidelines Go, which is the core value proposition of Wails.*

### Backwards Compatibility
*   This is an additive feature for Wails v3.
*   The `Wails.CallGo` and `Wails.on` APIs match the existing v3 desktop patterns to ensure a high degree of code reuse in the frontend.

### Test Plan
*   **Unit Tests:** Go-side logic for binding and event dispatching.
*   **Integration Tests:** A "Hello World" sample app testing the full loop: JS -> Go -> Java (Permissions) -> Result back to JS.
*   **Manual Testing:** Testing on various Android API levels (Min SDK 23+) and hardware configurations.

### Reference Implementation
A working prototype is available at: [https://github.com/its-ernest/wails-mobile](https://github.com/its-ernest/wails-mobile)
*Example implementation includes:*
*   Asset serving from Go.
*   Bidirectional Go-Java bridge.
*   Native Permissions Plugin.
*   Background Event Polling system.

### Maintenance Plan
The author (@its-ernest) will act as the primary maintainer for the Android mobile implementation. We will follow Wails v3 coding standards and ensure documentation is kept up to date in the official Wails docs.

### Conclusion
Mobile support is the next logical step for Wails v3. By adopting a Go-centric orchestration model and a plugin-based native bridge, we can provide a powerful, flexible, and familiar development experience for Gophers targeting mobile platforms.
