//go:build openharmony

package assetserver

import "net/url"

// OpenHarmony serves the embedded frontend over a local HTTP server,
// loaded by ArkWeb via http://127.0.0.1:<port>/. Use a wails scheme base
// consistent with the linux/desktop default; the shell points ArkWeb at the
// real loopback URL, so this base is only used for relative-URL resolution.
var baseURL = url.URL{
	Scheme: "wails",
	Host:   "localhost",
}