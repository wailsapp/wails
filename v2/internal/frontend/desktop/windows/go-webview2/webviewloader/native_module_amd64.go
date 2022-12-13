//go:build windows && native_webview2loader

package webviewloader

import _ "embed"

//go:embed x64/WebView2Loader.dll
var WebView2Loader []byte
