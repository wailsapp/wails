# Wails v3 iOS Architecture

## Executive Summary

This document provides a comprehensive technical architecture for iOS support in Wails v3. The implementation enables Go applications to run natively on iOS with a WKWebView frontend, maintaining the Wails philosophy of using web technologies for UI while leveraging Go for business logic.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [Layer Architecture](#layer-architecture)
4. [Implementation Details](#implementation-details)
5. [Battery Optimization](#battery-optimization)
6. [Build System](#build-system)
7. [Security Considerations](#security-considerations)
8. [API Reference](#api-reference)

## Architecture Overview

### Design Principles

1. **Battery Efficiency First**: All architectural decisions prioritize battery life
2. **No Network Ports**: Asset serving happens in-process via native APIs
3. **Minimal WebView Instances**: Maximum 2 concurrent WebViews (1 primary, 1 for transitions)
4. **Native Integration**: Deep iOS integration using Objective-C runtime
5. **Wails v3 Compatibility**: Maintain API compatibility with existing Wails v3 applications

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     iOS Application                          │
├─────────────────────────────────────────────────────────────┤
│                    UIKit Framework                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              WailsViewController                     │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │            WKWebView Instance                  │  │   │
│  │  │  ┌─────────────────────────────────────────┐  │  │   │
│  │  │  │         Web Application (HTML/JS)        │  │  │   │
│  │  │  └─────────────────────────────────────────┘  │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                  Bridge Layer (CGO)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │URL Handler   │  │JS Bridge     │  │Message Handler│     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
├─────────────────────────────────────────────────────────────┤
│                    Go Runtime                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                 Wails Application                     │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │  │
│  │  │App Logic │  │Services  │  │Asset Server      │  │  │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Platform Layer (`application_ios.go`)

**Purpose**: Go interface for iOS platform operations

**Key Functions**:
- `platformRun()`: Initialize and run the iOS application
- `platformQuit()`: Gracefully shutdown the application
- `isDarkMode()`: Detect iOS dark mode state
- `ExecuteJavaScript(windowID uint, js string)`: Execute JS in WebView

**Exported Go Functions (Called from Objective-C)**:
- `ServeAssetRequest(windowID C.uint, urlStr *C.char, callbackID C.uint)`
- `HandleJSMessage(windowID C.uint, message *C.char)`

### 2. Native iOS Layer (`application_ios.m`)

**Components**:

#### WailsSchemeHandler
```objc
@interface WailsSchemeHandler : NSObject <WKURLSchemeHandler>
```
- Implements `WKURLSchemeHandler` protocol
- Intercepts `wails://` URL requests
- Bridges to Go for asset serving
- Manages pending requests with callback IDs

**Methods**:
- `startURLSchemeTask:`: Intercept request, call Go handler
- `stopURLSchemeTask:`: Cancel pending request
- `completeRequest:withData:mimeType:`: Complete request with data from Go

#### WailsMessageHandler
```objc
@interface WailsMessageHandler : NSObject <WKScriptMessageHandler>
```
- Implements JavaScript to Go communication
- Handles `window.webkit.messageHandlers.external.postMessage()`
- Serializes messages to JSON for Go processing

**Methods**:
- `userContentController:didReceiveScriptMessage:`: Process JS messages

#### WailsViewController
```objc
@interface WailsViewController : UIViewController
```
- Main view controller containing WKWebView
- Manages WebView lifecycle
- Handles JavaScript execution requests

**Properties**:
- `webView`: WKWebView instance
- `schemeHandler`: Custom URL scheme handler
- `messageHandler`: JS message handler
- `windowID`: Unique window identifier

**Methods**:
- `viewDidLoad`: Initialize WebView with configuration
- `executeJavaScript:`: Run JS code in WebView

### 3. Bridge Layer (CGO)

**C Interface Functions**:
```c
void ios_app_init(void);                    // Initialize iOS app
void ios_app_run(void);                     // Run main loop
void ios_app_quit(void);                    // Quit application
bool ios_is_dark_mode(void);                // Check dark mode
unsigned int ios_create_webview(void);      // Create WebView
void ios_execute_javascript(unsigned int windowID, const char* js);
void ios_complete_request(unsigned int callbackID, const char* data, const char* mimeType);
```

## Layer Architecture

### Layer 1: Presentation Layer (WebView)

**Responsibilities**:
- Render HTML/CSS/JavaScript UI
- Handle user interactions
- Communicate with native layer

**Key Features**:
- WKWebView for modern web standards
- Hardware-accelerated rendering
- Efficient memory management

### Layer 2: Communication Layer

**Request Interception**:
```
WebView Request → WKURLSchemeHandler → Go ServeAssetRequest → AssetServer → Response
```

**JavaScript Bridge**:
```
JS postMessage → WKScriptMessageHandler → Go HandleJSMessage → Process → ExecuteJavaScript
```

### Layer 3: Application Layer (Go)

**Components**:
- Application lifecycle management
- Service binding and method calls
- Asset serving from embedded fs.FS
- Business logic execution

### Layer 4: Platform Integration Layer

**iOS-Specific Features**:
- Dark mode detection
- System appearance integration
- iOS-specific optimizations

## Implementation Details

### Request Handling Flow

1. **WebView makes request** to `wails://localhost/path`
2. **WKURLSchemeHandler intercepts** request
3. **Creates callback ID** and stores `WKURLSchemeTask`
4. **Calls Go function** `ServeAssetRequest` with URL and callback ID
5. **Go processes request** through AssetServer
6. **Go calls** `ios_complete_request` with response data
7. **Objective-C completes** the `WKURLSchemeTask` with response

### JavaScript Execution Flow

1. **Go calls** `ios_execute_javascript` with JS code
2. **Bridge dispatches** to main thread
3. **WKWebView evaluates** JavaScript
4. **Completion handler** logs any errors

### Message Passing Flow

1. **JavaScript calls** `window.webkit.messageHandlers.wails.postMessage(data)`
2. **WKScriptMessageHandler receives** message
3. **Serializes to JSON** and passes to Go
4. **Go processes** message in `HandleJSMessage`
5. **Go can respond** via `ExecuteJavaScript`

## Battery Optimization

### WebView Configuration

```objc
// Disable unnecessary features
config.suppressesIncrementalRendering = NO;
config.allowsInlineMediaPlayback = YES;
config.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone;
```

### Memory Management

1. **Single WebView Instance**: Reuse instead of creating new instances
2. **Automatic Reference Counting**: Use ARC for Objective-C objects
3. **Lazy Loading**: Initialize components only when needed
4. **Resource Cleanup**: Properly release resources when done

### Request Optimization

1. **In-Process Serving**: No network overhead
2. **Direct Memory Transfer**: Pass data directly without serialization
3. **Efficient Caching**: Leverage WKWebView's built-in cache
4. **Minimal Wake Locks**: No background network activity

## Build System

### Build Tags

```go
//go:build ios
```

### CGO Configuration

```go
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework WebKit
```

### Build Script (`build_ios.sh`)

**Steps**:
1. Check dependencies (go, xcodebuild, xcrun)
2. Set up iOS cross-compilation environment
3. Build Go binary with iOS tags
4. Create app bundle structure
5. Generate Info.plist
6. Sign for simulator
7. Create launch script

**Environment Variables**:
```bash
export CGO_ENABLED=1
export GOOS=ios
export GOARCH=arm64
export SDK_PATH=$(xcrun --sdk iphonesimulator --show-sdk-path)
```

### Simulator Deployment

```bash
xcrun simctl install "$DEVICE_ID" "WailsIOSDemo.app"
xcrun simctl launch "$DEVICE_ID" "com.wails.iosdemo"
```

## Security Considerations

### URL Scheme Security

1. **Custom Scheme**: Use `wails://` to avoid conflicts
2. **Origin Validation**: Only serve to authorized WebViews
3. **No External Access**: Scheme handler only responds to app's WebView

### JavaScript Execution

1. **Input Validation**: Sanitize JS code before execution
2. **Sandboxed Execution**: WKWebView provides isolation
3. **No eval()**: Avoid dynamic code evaluation

### Data Protection

1. **In-Memory Only**: No temporary files on disk
2. **ATS Compliance**: App Transport Security enabled
3. **Secure Communication**: All data stays within app process

## API Reference

### Go API

#### Application Functions

```go
// Create new iOS application
app := application.New(application.Options{
    Name: "App Name",
    Description: "App Description",
})

// Run the application
app.Run()

// Execute JavaScript
app.ExecuteJavaScript(windowID, "console.log('Hello')")
```

#### Service Binding

```go
type MyService struct{}

func (s *MyService) Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

### JavaScript API

#### Send Message to Go

```javascript
window.webkit.messageHandlers.wails.postMessage({
    type: 'methodCall',
    service: 'MyService',
    method: 'Greet',
    args: ['World']
});
```

#### Receive from Go

```javascript
window.wailsCallback = function(data) {
    console.log('Received:', data);
};
```

### Objective-C Bridge API

#### From Go to Objective-C

```c
// Execute JavaScript
ios_execute_javascript(windowID, "alert('Hello')");

// Complete asset request
ios_complete_request(callbackID, htmlData, "text/html");
```

#### From Objective-C to Go

```c
// Serve asset request
ServeAssetRequest(windowID, urlString, callbackID);

// Handle JavaScript message
HandleJSMessage(windowID, jsonMessage);
```

## Performance Metrics

### Target Metrics

- **WebView Creation**: < 100ms
- **Asset Request**: < 10ms for cached, < 50ms for first load
- **JS Execution**: < 5ms for simple scripts
- **Message Passing**: < 2ms round trip
- **Memory Usage**: < 50MB baseline
- **Battery Impact**: < 2% per hour active use

### Monitoring

1. **Xcode Instruments**: CPU, Memory, Energy profiling
2. **WebView Inspector**: JavaScript performance
3. **Go Profiling**: pprof for Go code analysis

## Future Enhancements

### Phase 1: Core Stability
- [ ] Production-ready error handling
- [ ] Comprehensive test suite
- [ ] Performance optimization

### Phase 2: Feature Parity
- [ ] Multiple window support
- [ ] System tray integration
- [ ] Native menu implementation

### Phase 3: iOS-Specific Features
- [ ] Widget extension support
- [ ] App Clip support
- [ ] ShareSheet integration
- [ ] Siri Shortcuts

### Phase 4: Advanced Features
- [ ] Background task support
- [ ] Push notifications
- [ ] CloudKit integration
- [ ] Apple Watch companion app

## Conclusion

This architecture provides a solid foundation for iOS support in Wails v3. The design prioritizes battery efficiency, native performance, and seamless integration with the existing Wails ecosystem. The proof of concept demonstrates all four required capabilities:

1. ✅ **WebView Creation**: Native WKWebView with optimized configuration
2. ✅ **Request Interception**: Custom scheme handler without network ports
3. ✅ **JavaScript Execution**: Bidirectional communication bridge
4. ✅ **iOS Simulator Support**: Complete build and deployment pipeline

The architecture is designed to scale from this proof of concept to a full production implementation while maintaining the simplicity and elegance that Wails developers expect.