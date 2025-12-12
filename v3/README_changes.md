Summary of Browser Mode Changes
I've implemented a Browser Mode feature for Wails v3 that allows you to use webview windows as browsers for external sites, with the ability to inject JavaScript and extract data like cookies, localStorage, and page content.

Files Modified/Created
1. v3/pkg/application/webview_window_options.go
   Added new types for browser mode configuration:
```go
type BrowserModeOptions struct {
    DisableWebSecurity           bool     // Disables CORS (Windows only)
    InjectScriptOnNavigation     string   // JS injected after page load
    InjectScriptAtDocumentStart  string   // JS injected before page scripts
    EnableDataExtraction         bool     // Enables data extraction APIs
    OnNavigationComplete         func(url string)
    OnDataExtracted              func(data *BrowserData)
}

type BrowserData struct {
    WindowName, URL, HTMLContent string
    Cookies, LocalStorage, SessionStorage, URLParams map[string]string
    CustomData map[string]interface{}
}
```

webview_window_options.go
v3/pkg/application
2. v3/pkg/application/webview_window.go
   Added methods to the WebviewWindow:

IsBrowserMode() - Check if browser mode is enabled
ExtractBrowserData() - Extract cookies, localStorage, etc. asynchronously
InjectBrowserScript(script) - Inject JS into current page
SetBrowserCustomData(key, value) - Store custom data for extraction
3. v3/pkg/application/webview_window_darwin.go
   macOS implementation:

Added windowAddUserScript() C function using WKUserScript
Added windowGetCurrentURL() C function
setupBrowserMode() injects scripts via WKUserContentController
Scripts run on every navigation, not just initial load
4. v3/pkg/application/webview_window_windows.go
   Windows implementation:

Added --disable-web-security flags when DisableWebSecurity is true
setupBrowserModeScripts() uses chromium.Init() for script injection
Added getCurrentURL() method
Navigation callback triggers OnNavigationComplete
5. v3/pkg/application/linux_cgo.go
   Linux implementation:

Added addUserScript() using webkit_user_script_new()
Added getCurrentURL() using webkit_web_view_get_uri()
setupBrowserMode() adds scripts to WebKitUserContentManager
6. v3/pkg/application/webview_window_linux.go
   Calls setupBrowserMode() before loading URL
   Triggers OnNavigationComplete callback on WindowLoadChanged event
7. v3/pkg/application/browser_data_store.go (NEW)
   Thread-safe store for extracted browser data:

GetBrowserDataStore() - Get global instance
Store(windowName, data) / Get(windowName) - Store/retrieve data
Helper methods: GetCookies(), GetLocalStorage(), etc.
8. v3/pkg/application/messageprocessor_browser.go
   Added BrowserSendData method to handle extracted data from JavaScript and store it in BrowserDataStore.

How It Works
Create a browser mode window:
```go
browserWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "browser",
    URL:   "https://example.com",
    BrowserMode: &application.BrowserModeOptions{
        DisableWebSecurity: true,  // Windows only - disables CORS
        EnableDataExtraction: true,
        InjectScriptOnNavigation: `
                console.log("Injected on every page!");
                window._wailsSetBrowserData("customKey", document.title);
            `,
        OnNavigationComplete: func(url string) {
			fmt.Println("Navigated to:", url)
        },
        OnDataExtracted: func(data *application.BrowserData) {
            fmt.Println("Cookies:", data.Cookies)
            fmt.Println("LocalStorage:", data.LocalStorage)
        },
    },
})
```
Extract data on demand:
```go
browserWindow.ExtractBrowserData()  // Async - triggers OnDataExtracted callback
```

Access data from any window:
```go
store := application.GetBrowserDataStore()
data := store.Get("browser")
cookies := store.GetCookies("browser")
```

Platform Differences

| Feature                                      | macOS                           | Windows               | Linux               |
|----------------------------------------------|---------------------------------|-----------------------|---------------------|
| Script Injection	                            | ✅ WKUserScript                  | 	✅ chromium.Init()    | 	✅ WebKitUserScript |
| Disable Web Security | 	⚠️ Limited                     | 	✅ Chrome flags	      | ⚠️ Limited          |
| Cookie Access | 	✅ JS only | 	✅ JS only | 	✅ JS only          |

**Note: DisableWebSecurity fully works on Windows via Chromium flags. On macOS/Linux, WKWebView and WebKitGTK don't have simple flags to disable CORS - you'd need a proxy server for true cross-origin requests.**