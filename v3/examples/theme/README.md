### Theme System for V3 (Alpha Feature)

This PR introduces the foundational **Theme system** for WINDOWS and MACOS theming in V3.

---

### Theme Architecture: Intent vs. Resolution

This system implements a strict separation between **Theme Intent** (what the developer/user wants) and **Theme Resolution** (the actual values sent to the operating system APIs especially in case of Windows OS). 

To achieve this cleanly, the architecture introduces three distinct theme enums across the application lifecycle:

#### 1. `AppTheme` (Application Level Intent)
*Defined in `theme_application.go`*

`AppTheme` represents the global theme preference for the entire application. It resolves to a string-based enum for maximum flexibility.
- **Available Values:** `AppSystemDefault` (Follow OS), `AppDark` (Force Dark), `AppLight` (Force Light).
- **Behavior:** This value is initialized via `App.Options.Theme` and stored in `App.theme`. It dictates the fallback behavior for all windows unless explicitly overridden. It's value can be read via `app.GetTheme()` and set via `app.SetTheme()`.

#### 2. `WinTheme` (Window Level Intent)
*Defined in `theme_window.go`*

`WinTheme` represents the theme preference for a specific, individual window cross-platform. It is also a string-based enum.
- **Available Values:**
  - `WinThemeApplication` - Inherit the `AppTheme` resolving logic (Default).
  - `WinThemeSystem` - Force the window to follow the OS theme, completely ignoring the `AppTheme`.
  - `WinThemeDark` - Force the window into Dark mode.
  - `WinThemeLight` - Force the window into Light mode.
- **Behavior:** Configured via window options at creation, and used via `window.GetTheme()` and `window.SetTheme()` methods. This state is purely intent-based and does not directly call OS bindings.

#### 3. `Theme` (Resolved Output State for Windows OS)
*Defined in `theme_window_windows.go`*

`Theme` is the final, resolved integer representation (`0` = System, `1` = Dark, `2` = Light) used exclusively by Windows OS to execute native OS layout changes.
- **Behavior:** This enum is completely invisible to the frontend and cross-platform APIs. When a Windows OS specific window is asked to update its appearance, it evaluates the `WinTheme` (and potentially the fallback `AppTheme`), runs the resolution function, and outputs this strongly-typed integer state to execute the native API calls.

---

### The `followApplicationTheme` Flag

Rather than storing a dedicated `Theme = WinThemeApplication` state persistently inside the unexported window struct, the system utilizes a simple boolean flag: `w.parent.followApplicationTheme`.

**Why use a flag instead of an enum?**
Because of how macOS handles themes natively (by querying the OS directly for what appearance is currently active), storing a resolved state enum on the window would be redundant and prone to desynchronization. Instead, we only need to know **one thing**: *Is this window currently slaved to the Application's global theme?*

1. During initialization (or passing `WinThemeApplication` to `SetTheme`), we set `followApplicationTheme = true`. 
2. Every time a theme needs to be calculated or fetched, we check this flag. 
3. If it is `true`, we infer the resolved theme from the Application theme. If it is `false` (because the user explicitly set the window to System, Dark, or Light), we resolve based on the window's specific settings.

---

### Windows Theme Handling (Windows OS)

For Windows, the UI framework relies on us explicitly telling the OS which theme to render. Therefore, Windows uses our 3-tier enum system to cascade state predictably:

1. **Check the Flag:** If `followApplicationTheme` is `true`, the window looks at the global application's `AppTheme` intent (`AppSystemDefault`, `AppDark`, `AppLight`) and maps that intent to the resolved integer `Theme` (`0`, `1`, or `2`).
2. **Check the Window Intent:** If `followApplicationTheme` is `false`, the window checks its explicit `WinTheme` intent (`WinThemeSystem`, `WinThemeDark`, or `WinThemeLight`) and maps that to the resolved `Theme`.

This strict unidirectional flow guarantees that the frontend developer only interacts with high-level strings (`AppTheme` and `WinTheme`), while the messy OS-specific fallback calculations are completely contained.

---

### macOS Theme Handling (Darwin)

MacOS handles light and dark Themes by managing `NSAppearance`. macOS has different Appearances for both light and dark themes, such as `NSAppearanceNameAqua` (light) vs. `NSAppearanceNameDarkAqua` (dark) or their accessibility high-contrast equivalents.

The macOS implementation handles light/dark/system swicthing by following the exact same Intent vs. Resolution flow as Windows, however rather than store the resolved theme, **it infers the resolved theme from Internal OS Mechanics.**

MacOS allows us to access the Appearance and Effective Appearance of the current window, and we use this to determine the current theme dynamically:

1. We utilize `NSAppearance` names to represent UI appearance.
2. The theme system divides different `NSAppearance` instances into light and dark variants.

**MacOS Initialization & Resolution:**
- If the user does *not* set an explicit appearance in the Mac window options, we assume the window should follow the app, and `followApplicationTheme` is set to `true`.
- If the Application's theme intent is System Default, we do nothing and keep appearance = `""` (which causes macOS to naturally inherit the OS appearance).
- If the Application's theme intent is Light or Dark, we assign the corresponding Light or Dark appearance to the window.
- Conversely, if the user explicitly set an appearance (e.g., `appearance = "NSAppearanceNameAqua"`) during window creation, `followApplicationTheme` is `false`, and we force that explicit appearance.

**Inferring State for `GetTheme()` and `SetTheme()`:**
Instead of saving what theme the window is, `win.GetTheme()` reverse-engineers the intent by querying macOS:
1. If `followApplicationTheme` is `true`, immediately return `WinThemeApplication`.
2. If it is `false`, we check the window's explicit macOS appearance. If it represents an empty string (`""`), it means we are following the system, returning `WinThemeSystem`.
3. If the appearance is explicitly populated, we check if that specific macOS Appearance is a dark variant or light variant, and return `WinThemeDark` or `WinThemeLight` accordingly.

This reverse-inference means our Window struct stays incredibly lightweight, relying directly on the native OS source-of-truth.

---

### Notes

* This PR focuses on introducing the **core theme abstraction and API surface**.
* Platform-specific behavior will continue to be implemented or extended in follow-up changes where necessary.
* **Refactor Note:** The previous integer-based `Theme` enum in `webview_window_options.go` has been relocated to `theme_window_windows.go` to demonstrate its new use as internal record keeping of applied Theme.

---

### Files Updated
1. `webview_window_options_test.go`
2. `webview_window_options.go`
3. `application.go`
4. `application_darwin.go`
5. `application_windows.go`
6. `application_android.go`
7. `application_ios.go`
8. `application_linux.go`
9. `application_linus_gtk4.go`
10. `application_server.go`

### Files Created
1. `theme_application.go`
2. `theme_window.go`
4. `theme_webview_window_darwin.go`
5. `theme_webview_window_windows.go`
