# Android Implementation TODO

This document tracks outstanding implementation tasks for Android support in Wails v3.

## Completed ✅

### Clipboard Support ✅
**Files:** `WailsBridge.java`, `application_android.go`, `clipboard_android.go`
- Added `setClipboardText()` and `getClipboardText()` Java methods
- Added JNI wrappers and method caching
- Updated Go implementation to call JNI

### Dark Mode Detection ✅
**Files:** `WailsBridge.java`, `application_android.go`
- Added `isDarkMode()` Java method using Configuration.uiMode
- Added JNI wrapper for Go to call

### Screen Information ✅
**Files:** `WailsBridge.java`, `application_android.go`, `screen_android.go`
- Added `getScreenInfo()` Java method using DisplayMetrics
- Returns actual device dimensions instead of hardcoded values
- Calculates scale factor from DPI

### WebView Background Color ✅
**Files:** `WailsBridge.java`, `application_android.go`, `webview_window_android.go`
- Added `setWebViewBackgroundColor()` Java method
- Converts RGBA to Android ARGB format
- Applied via JNI on UI thread

### Message Dialogs ✅
**Files:** `WailsBridge.java`, `application_android.go`, `dialogs_android.go`
- Added `showMessageDialog()` Java method using AlertDialog.Builder
- Uses CountDownLatch for synchronous result handling
- Supports info, warning, error, and question dialog types
- Supports up to 3 buttons (positive, negative, neutral)

### WebView setHTML/setURL ✅
**Files:** `WailsBridge.java`, `application_android.go`, `webview_window_android.go`
- Added `setHTML()` and `setURL()` Java methods
- Uses `loadDataWithBaseURL()` for HTML content
- Uses `loadUrl()` for URL navigation
- Properly runs on main thread via Handler

---

## Priority 1: Core Functionality

### File Dialogs (Open, Save, Directory)
**Files:** `pkg/application/dialogs_android.go`

**Current State:** Stub implementation returning empty strings.

**Implementation:**
Use Android's Storage Access Framework (SAF):
1. Launch Intent with `ACTION_OPEN_DOCUMENT` / `ACTION_CREATE_DOCUMENT`
2. Handle result via Activity callback
3. Convert content:// URI to usable path or stream

**Complexity:** High - requires Activity result handling, SAF permissions

---

## Priority 2: Android-Specific Enhancements

### About Dialog
**Files:** `pkg/application/application_android.go`

**Current State:** TODO comment, no implementation.

**Implementation:**
Use AlertDialog with app icon, name, version from PackageInfo.

---

## Priority 3: Future Features

### Content Protection (FLAG_SECURE)
**Files:** `pkg/application/webview_window_android.go:332-334`

Prevent screenshots/screen recording:
```java
getWindow().setFlags(WindowManager.LayoutParams.FLAG_SECURE, WindowManager.LayoutParams.FLAG_SECURE);
```

### Multi-Display Support
**Files:** `pkg/application/screen_android.go`

Use DisplayManager to enumerate connected displays (tablets, ChromeOS, etc).

### Deep Linking
Handle custom URL schemes and app links.

### Push Notifications (FCM)
Firebase Cloud Messaging integration.

### Biometric Authentication
Use BiometricPrompt API.

---

*Last Updated: December 2024*
