// +build windows

package x64

import _ "embed"

//go:embed WebView2Loader.dll
var WebView2Loader []byte
