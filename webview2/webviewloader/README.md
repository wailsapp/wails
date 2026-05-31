# Webviewloader

Webviewloader is a port of [OpenWebView2Loader](https://github.com/jchv/OpenWebView2Loader) to Go.

It is intended to be feature-complete with the original WebView2Loader distributed with
the WebView2 NuGet package, but some features are intentionally not implemented.

## Status

- [x] CompareBrowserVersions
- [x] CreateCoreWebView2Environment
- [x] CreateCoreWebView2EnvironmentWithOptions
- [x] GetAvailableCoreWebView2BrowserVersionString

## Not implemented features

- Registry Overrides of Parameters
- Env Variable Overrides of Parameters
- Does not incorporate `GetCurrentPackageInfo` to search for an installed runtime
