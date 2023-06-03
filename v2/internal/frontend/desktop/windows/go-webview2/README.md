# go-webview2

This is a proof of concept for embedding Webview2 into Go without CGo. It is
based on [webview/webview](https://github.com/webview/webview) and provides a
compatible API.

## Notice

Because this version doesn't currently have an EdgeHTML fallback, it will not
work unless you have a Webview2 runtime installed. In addition, it requires the
Webview2Loader DLL in order to function. Adding an EdgeHTML fallback should be
technically possible but will likely require much worse hacks since the API is
not strictly COM to my knowledge.

## Demo

For now, you'll need to install the Webview2 runtime, as it does not ship with
Windows.

[WebView2 runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)

After that, you should be able to run go-webview2 directly:

```
go run go-webview2/cmd/demo
```

This will use go-winloader to load an embedded copy of WebView2Loader.dll.

If this does not work, please try running from a directory that has an
appropriate copy of `WebView2Loader.dll` for your GOARCH. If _that_ worked,
_please_ file a bug so we can figure out what's wrong with go-winloader :)
