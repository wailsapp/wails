// +build windows

package x64

import _ "embed"

//go:embed webview.dll
var WebView2 []byte

//go:embed WebView2Loader.dll
var WebView2Loader []byte
