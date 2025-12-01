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

---

## Priority 1: Core Functionality

### Message Dialogs (Info, Warning, Error, Question)
**Files:** `pkg/application/dialogs_android.go`

**Current State:** Stub implementation returning empty values.

**Implementation:**
1. Add JNI method to WailsBridge.java using AlertDialog.Builder
2. Handle async response via callback (dialogs are async on Android)
3. Support title, message, and button configuration

**Example Java:**
```java
public void showDialog(String type, String title, String message, String[] buttons, DialogCallback callback) {
    mainHandler.post(() -> {
        AlertDialog.Builder builder = new AlertDialog.Builder(context);
        builder.setTitle(title).setMessage(message);
        // Add buttons based on type...
        builder.show();
    });
}
```

---

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

## Priority 2: System Information

### Dark Mode Detection
**Files:** `pkg/application/application_android.go`, `pkg/application/application_android_nocgo.go`

**Current State:** Returns hardcoded `false`.

**Implementation:**
1. Add JNI method to query `Configuration.uiMode`:
   ```java
   public boolean isDarkMode() {
       int nightMode = context.getResources().getConfiguration().uiMode & Configuration.UI_MODE_NIGHT_MASK;
       return nightMode == Configuration.UI_MODE_NIGHT_YES;
   }
   ```
2. Wire up in Go via JNI

---

### Screen Information (DisplayMetrics)
**Files:** `pkg/application/screen_android.go`

**Current State:** Returns hardcoded 1080x2400.

**Implementation:**
1. Add JNI method to get DisplayMetrics:
   ```java
   public String getScreenInfo() {
       DisplayMetrics metrics = context.getResources().getDisplayMetrics();
       JSONObject info = new JSONObject();
       info.put("width", metrics.widthPixels);
       info.put("height", metrics.heightPixels);
       info.put("density", metrics.density);
       info.put("densityDpi", metrics.densityDpi);
       return info.toString();
   }
   ```
2. Parse JSON in Go and populate Screen struct

---

## Priority 3: Android-Specific Enhancements

### About Dialog
**Files:** `pkg/application/application_android.go:502`

**Current State:** TODO comment, no implementation.

**Implementation:**
Use AlertDialog with app icon, name, version from PackageInfo.

---

### WebView Background Color
**Files:** `pkg/application/webview_window_android.go:128-131`

**Current State:** Logs but doesn't apply.

**Implementation:**
1. Add JNI method: `setWebViewBackgroundColor(int color)`
2. Call `webView.setBackgroundColor(color)` on main thread

---

### setHTML / setURL
**Files:** `pkg/application/webview_window_android.go:336-356`

**Current State:** TODO comments, only logs.

**Implementation:**
1. Add JNI methods in WailsBridge
2. Call `webView.loadUrl()` or `webView.loadDataWithBaseURL()`

---

## Priority 4: Future Features

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

## Code Locations with TODO Comments

| File | Line | TODO |
|------|------|------|
| `clipboard_android.go` | 13, 19 | JNI to ClipboardManager |
| `dialogs_android.go` | 18-34 | AlertDialog implementation |
| `dialogs_android.go` | 37-58 | File picker Intent |
| `application_android.go` | 427, 502 | Dark mode, About dialog |
| `application_android_nocgo.go` | 52, 130 | Dark mode, About dialog |
| `screen_android.go` | 8 | DisplayManager support |
| `webview_window_android.go` | 129, 337, 354 | Background color, setHTML, setURL |

---

*Last Updated: December 2024*
