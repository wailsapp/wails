//go:build darwin && !ios && !server

package application

import _ "embed"

// WKWebView passes a wails:// scheme handler an empty body for requests whose
// body is an ArrayBuffer or typed-array view (Blob bodies arrive intact), with
// both the cgo and the purego backend. Both install this script at document
// start so those bodies are wrapped in a Blob first. It is the macOS
// counterpart of the Linux Blob body shim. See BUGS_FOUND.md #9.
//
//go:embed darwin_body_fetch_shim.js
var darwinBodyFetchShimJS string
