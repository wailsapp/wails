# Windows asset request cancellation regression probe

This probe exercises the Wails v3 application path for issue [#5963](https://github.com/wailsapp/wails/issues/5963), including the shared `webViewAssetRequest` wrapper. It launches hidden WebView2 windows and a loopback HTTP control server, runs for 15 seconds, writes `request-cancellation-results.json` and `native.log`, and exits nonzero on a failed assertion.

Build from `v3`:

```sh
go build -o cancellation-probe.exe ./pkg/application/testdata/request-cancellation
```

Run the executable in an interactive Windows desktop session. An SSH service session cannot initialise WebView2; use an interactive scheduled task or run from the desktop. No remote debugging port or external CDP client is required.

The probe checks fetch/XHR aborts, POST bodies and headers, redirects, identical URLs, iframe requests, worker aborts and termination, navigation, window closure, and 40 concurrent aborts. Controls verify normal completion, local and HTTP keepalive navigation, keepalive redirects and uploads, external redirects without added CORS preflights, and unchanged response URLs.

Additional automated checks:

```sh
go test ./pkg/application -run 'TestWindowsCancellation|TestAssetRequest'
go test ./internal/assetserver/...
node --test pkg/application/request_keepalive_windows_test.mjs
```

## Implementation notes

WebView2 does not expose an early cancellation event on `WebResourceRequested`. The Windows implementation uses its in-process DevTools API: root `Fetch.requestPaused` events associate an opaque marker with a network request ID, and `Network` terminal events cancel that request's Go context. Worker network events arrive on child sessions while their interceptions arrive on the root session. Window/application teardown also cancels pending contexts.

DevTools reports `ERR_ABORTED` when a page disappears even for keepalive fetches. A script installed before page and worker code gives local keepalive requests an independent identity and reports explicit aborts separately. This identity travels in the URL fragment, which HTTP does not send to handlers or external redirect targets; adding a header here would cause unwanted CORS preflights. Native request markers are stripped before application handlers run. Failed JavaScript response delivery lets an already dispatched keepalive handler finish, then releases its context.

The shared application wrapper forwards native contexts on every platform. Apple platforms already provide native cancellation contexts; their forwarding is covered by the application regression tests. Linux's native abort wiring is outside this Windows change.
