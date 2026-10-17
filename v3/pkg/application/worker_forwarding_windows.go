//go:build windows && !server

package application

import _ "embed"

// workerForwardingScript is injected into every attached worker DevTools
// session, before the worker's first script runs.
//
// WebView2 does not deliver the response for a worker script that another worker
// requests from the asset server, so a Worker created from inside a Worker never
// starts and reports a silent error event (wails#6128). The script replaces the
// worker's Worker constructor so that a same-origin URL is forwarded through a
// blob URL that imports the original absolute URL, which WebView2 does serve.
//
//go:embed worker_forwarding_windows.js
var workerForwardingScript string
