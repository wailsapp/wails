# WebView API Compatibility Check

This example application tests and reports which Web APIs are available in the current WebView engine.

## Purpose

Different platforms use different WebView engines:
- **Linux GTK4**: WebKitGTK 6.0 (WebKit-based)
- **Linux GTK3**: WebKit2GTK 4.1 (WebKit-based)
- **Windows**: WebView2 (Chromium-based)
- **macOS**: WKWebView (WebKit-based)

Each engine supports different Web APIs. This tool helps you understand what APIs are available for your Wails application on each platform.

## Building

```bash
# Linux GTK4
go build -tags gtk4 -o webview-api-check .

# Linux GTK3
go build -tags gtk3 -o webview-api-check .

# Windows/macOS
go build -o webview-api-check .
```

## Usage

1. Run the application
2. Click "Run API Tests" to test all Web APIs
3. View results organized by category
4. Use filters to find specific APIs
5. Export report as JSON for comparison

## API Categories Tested

| Category | APIs Tested |
|----------|-------------|
| Storage | localStorage, IndexedDB, Cache API, File System Access |
| Network | Fetch, WebSocket, EventSource, WebTransport |
| Media | Web Audio, MediaRecorder, MediaDevices, Speech |
| Graphics | Canvas, WebGL, WebGL2, WebGPU |
| Device | Geolocation, Sensors, Battery, Bluetooth, USB |
| Workers | Web Workers, Service Workers, Shared Workers |
| Performance | Observers, Timing APIs |
| Security | Web Crypto, Credentials, WebAuthn |
| UI & DOM | Custom Elements, Shadow DOM, Pointer Events |
| CSS | CSSOM, Container Queries, Modern Selectors |
| JavaScript | ES Modules, BigInt, Private Fields, etc. |

## Understanding Results

- **Supported** (green): API is fully available
- **Partial** (yellow): API exists but may have limitations
- **Unsupported** (red): API is not available

Some APIs are marked with notes:
- "Chromium only" - Available in WebView2 but not WebKit
- "Experimental" - May not be stable
- "Requires secure context" - Needs HTTPS
- "PWA context" - Only available in installed PWAs

## Comparing Platforms

Run the app on different platforms and export JSON reports. Compare them to understand API availability differences:

```bash
# On Linux GTK4
./webview-api-check
# Export: webview-api-report-linux-20240115-143052.json

# On Windows
./webview-api-check.exe
# Export: webview-api-report-windows-20240115-143052.json
```

## Known Differences

### WebKitGTK vs WebView2 (Chromium)

WebView2 (Windows) typically supports more APIs because Chromium is updated more frequently:
- File System Access API (Chromium only)
- Web Serial, WebHID, WebUSB (Chromium only)
- Various experimental features

WebKitGTK may have better support for:
- Standard DOM APIs
- CSS features (varies by version)

### GTK3 vs GTK4 WebKitGTK

GTK4 uses WebKitGTK 6.0, GTK3 uses WebKit2GTK 4.1. The WebKit version determines API support, not GTK version.
