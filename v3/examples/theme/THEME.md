### Theme System for V3 (Alpha Feature)

This PR introduces the foundational **Theme system** for WINDOWS and MACOS theming in V3.

---

### Theme Scope

This proposed system divides Theme handling to two levels:
1. Application Level
2. Window Level
Each level defines how theme values are resolved and inherited.

---

#### 1. Application Level

At the application level, the theme applies to the entire Go application. The application theme can be configured to:
- Follow the operating system theme (Default Behaviour)
- Force Light theme
- Force Dark theme

When set at this level, all windows inherit the application theme unless explicitly overridden.

---

#### 2. Window Level

At the window level, the theme can be configured per individual window. The window level theme can be configured to:
- Inherit the application-level theme (default behavior)
- Follow the operating system theme
- Force Light theme
- Force Dark theme

If a window explicitly sets its own theme, it overrides the application-level setting for that window only.

---

### Theme Enum

A new THEME string enum is introduced to define available theme modes.

Available values:
* **System** – Follow the OS theme.
* **Light** – Force the light theme.
* **Dark** – Force the dark theme.
* **Application** – Follow the application-level theme setting *(applicable for individual windows only)*.

This enum provides a unified abstraction for theme control across platforms while allowing platform-specific implementations where required.

Note:- The old Int based Enum defined in *webview_window_options.go* has been removed. The Enum is defined in *theme.go*.

---

### Theme Enum Update

Theme Representation Update has been updated to use string based Enum instead of Int based Enum used previously.

The previous theme configuration used integer constants:

- SystemDefault = 0
- Dark          = 1
- Light         = 2

This PR updates the theme representation to use string values instead. The motivation for this change is due to following reasons:
1. Support the new Application theme value which is used for Window Level Theme (without adding separate Enum for Window Level Theme)
2. Support the new Application theme value
3. Improve usability for SetTheme / GetTheme APIs

Note:- The test case for the Enum has been updated in *webview_window_options_test*.

---

### Application-Level Theme

For Application level theme, the configuration follows the same architectural pattern used throughout the codebase: Options are used for initialization, while the App struct stores runtime state.

Options.Theme which is already defined in *webview_window_options.go* supplies the initial theme value when creating the application. 

This options value is parsed during initialization and then applied to the internal application state defined in *application.go*.

Application-level theme control has been added through two new methods:

* `GetTheme()`
* `SetTheme(theme Theme)`

These methods are defined in *application.go* allow retrieving and updating the current application theme preference.

---

### macOS Theme Handling (Darwin)


MacOS allows light and dark Theme by manging NSAppearance. macOS has different Appearances for both light and dark themes. Such as NSAppearanceNameAqua (light) - NSAppearanceNameDarkAqua (dark) or NSAppearanceNameAccessibilityHighContrastAqua (light) - NSAppearanceNameAccessibilityHighContrastDarkAqua (dark)



The macOS implementation handles light/dark appearance switching while respecting system behaviour, accessibility appearances, and application-level theme inheritance.

1. We utilize NSAppearance names to represent UI appearance.
2. The theme system divides different NSAppearance into light and dark variants. Such as: NSAppearanceNameAqua - NSAppearanceNameAqua, NSAppearanceNameDarkAqua - NSAppearanceNameVibrantLight, etc.

The macOS Theme handles system/light/dark switch in the following manner:
- Get Resolved Theme value.
- If user wants to follow System then we clear appearance, which defaults to following OS
- If user wants light or dark we look for current appearance and get its light/dark version and apply it





Note:- the functions defined for macOS Theme hadling are defined in *theme_darwin.go*




### Notes

* This PR focuses on introducing the **core theme abstraction and API surface**.
* Platform-specific behavior will be implemented or extended in follow-up changes where necessary.
* The goal is to establish a consistent interface for theme control across platforms while keeping the implementation minimal for the alpha stage.

---

### Files Updated
1. webview_window_options_test.go
2. webview_window_options.go
3. application.go

### Files Created
1. theme.go
