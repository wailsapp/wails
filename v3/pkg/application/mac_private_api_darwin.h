//go:build darwin && !ios && !server

#ifndef WailsMacPrivateAPI_h
#define WailsMacPrivateAPI_h

// Internal bridge only; the public Go API is identical in both builds.
// mac_private_api_darwin.go implements all private macOS calls and requires
// -tags private_mac_apis. mac_public_api_darwin.go supplies public alternatives
// or no-ops by default. Do not put private selectors or keys in this header.

void wailsPrivateSetWebviewTransparent(void *webView);
void wailsPrivateSetWebviewBackgroundColour(void *webView, int r, int g, int b, int alpha);
void wailsPrivateSetGlassStyle(void *glassView, int style);
void wailsPrivateSetGlassGrouping(void *glassView, const char *groupID, double groupSpacing);

#ifdef WAILS_MAC_DEVTOOLS
void wailsPrivateOpenWebInspector(void *window);
void wailsPrivateEnableWebInspector(void *window);
#endif

#endif /* WailsMacPrivateAPI_h */
