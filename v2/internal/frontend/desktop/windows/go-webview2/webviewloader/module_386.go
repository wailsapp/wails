package webviewloader

import _ "embed"

//go:embed x86/WebView2Loader.dll
var WebView2Loader []byte
