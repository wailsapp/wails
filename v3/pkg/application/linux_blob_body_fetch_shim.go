//go:build linux && cgo && !android

package application

import _ "embed"

// WebKit does not support resolving some request bodies for custom URI schemes.
// On WebKitGTK, asking for a Blob- or FormData-backed body crashes in
// webkit_uri_scheme_request_get_http_body instead of returning an empty stream.
// Install this script at document start so fetch serialises Blob and FormData
// bodies, plus body streams from preconstructed Request objects, to ordinary
// bytes before WebKit hands a wails:// request to the native scheme handler.
//
//go:embed linux_blob_body_fetch_shim.js
var linuxBlobBodyFetchShimJS string
