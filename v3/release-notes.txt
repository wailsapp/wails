## Added
- Added `AdditionalLaunchArgs` to `WindowsWindow` options to allow for additional command line arguments to be passed to the WebView2 browser. in [PR](https://github.com/wailsapp/wails/pull/4467)
- Added Run go mod tidy automatically after wails init [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows Snapassist feature by @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)

## Fixed
- Add support for `allowsBackForwardNavigationGestures` in macOS WKWebView to enable two-finger swipe navigation gestures (#1857)
- Fixes issue where onClick didn't work for menu items initially set as disabled by @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Thanks to @IanVS for the initial investigation.
- Fix Vite server not being cleaned up when build fails (#4403)