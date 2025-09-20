//go:build ios

#ifndef APPLICATION_IOS_H
#define APPLICATION_IOS_H

#include <stdbool.h>
#import <UIKit/UIKit.h>

// Forward declarations
@class WailsViewController;
@class WailsAppDelegate;

// Global references
extern WailsAppDelegate *appDelegate;
extern unsigned int nextWindowID;

// Initialize the iOS application
void ios_app_init(void);

// Run the iOS application main loop
void ios_app_run(void);

// Quit the iOS application
void ios_app_quit(void);

// Check if dark mode is enabled
bool ios_is_dark_mode(void);

// Configure/show state for iOS WKWebView input accessory view (keyboard toolbar)
// If disabled, the accessory view will be hidden.
void ios_set_disable_input_accessory(bool disabled);
bool ios_is_input_accessory_disabled(void);

// Scrolling & bounce & indicators
void ios_set_disable_scroll(bool disabled);
bool ios_is_scroll_disabled(void);
void ios_set_disable_bounce(bool disabled);
bool ios_is_bounce_disabled(void);
void ios_set_disable_scroll_indicators(bool disabled);
bool ios_is_scroll_indicators_disabled(void);

// Navigation gestures
void ios_set_enable_back_forward_gestures(bool enabled);
bool ios_is_back_forward_gestures_enabled(void);

// Link previews
void ios_set_disable_link_preview(bool disabled);
bool ios_is_link_preview_disabled(void);

// Media playback
void ios_set_enable_inline_media_playback(bool enabled);
bool ios_is_inline_media_playback_enabled(void);
void ios_set_enable_autoplay_without_user_action(bool enabled);
bool ios_is_autoplay_without_user_action_enabled(void);

// Inspector
void ios_set_disable_inspectable(bool disabled);
bool ios_is_inspectable_disabled(void);

// User agent customization
void ios_set_user_agent(const char* ua);
const char* ios_get_user_agent(void);
void ios_set_app_name_for_user_agent(const char* name);
const char* ios_get_app_name_for_user_agent(void);

// Live runtime mutations (apply to existing WKWebView instances)
// These functions iterate current view controllers and update the active webviews on the main thread.
void ios_runtime_set_scroll_enabled(bool enabled);
void ios_runtime_set_bounce_enabled(bool enabled);
void ios_runtime_set_scroll_indicators_enabled(bool enabled);
void ios_runtime_set_back_forward_gestures_enabled(bool enabled);
void ios_runtime_set_link_preview_enabled(bool enabled);
void ios_runtime_set_inspectable_enabled(bool enabled);
void ios_runtime_set_custom_user_agent(const char* ua);

// Native bottom tab bar (UITabBar) controls
void ios_native_tabs_set_enabled(bool enabled);
bool ios_native_tabs_is_enabled(void);
void ios_native_tabs_select_index(int index);
// Configure native tabs items as a JSON array: [{"Title":"...","SystemImage":"..."}]
void ios_native_tabs_set_items_json(const char* json);
const char* ios_native_tabs_get_items_json(void);

// App-wide background colour control
// Setter accepts RGBA (0-255) and a flag indicating whether the colour is intentionally set by the app.
void ios_set_app_background_color(unsigned char r, unsigned char g, unsigned char b, unsigned char a, bool isSet);
// Getter returns true if a colour was set; outputs RGBA components via pointers when non-null.
bool ios_get_app_background_color(unsigned char* r, unsigned char* g, unsigned char* b, unsigned char* a);

// Create a WebView window and return its ID
unsigned int ios_create_webview(void);

// Create a WebView window with specified Wails ID and return its native handle
void* ios_create_webview_with_id(unsigned int wailsID);

// Execute JavaScript in a WebView by ID (legacy)
void ios_execute_javascript(unsigned int windowID, const char* js);

// Direct JavaScript execution on a specific window handle
void ios_window_exec_js(void* viewController, const char* js);

// Loaders
void ios_window_load_url(void* viewController, const char* url);
void ios_window_set_html(void* viewController, const char* html);

// Get the window ID from a native handle
unsigned int ios_window_get_id(void* viewController);

// Release a native handle when window is destroyed
void ios_window_release_handle(void* viewController);

// Go callbacks
extern void ServeAssetRequest(unsigned int windowID, void* urlSchemeTask);
extern void HandleJSMessage(unsigned int windowID, char* message);
extern bool hasListeners(unsigned int eventId);

// iOS Runtime bridges
// Trigger haptic impact with a style: "light"|"medium"|"heavy"|"soft"|"rigid"
void ios_haptics_impact(const char* style);
// Returns a JSON string with basic device info. Caller is responsible for freeing with free().
const char* ios_device_info_json(void);

#endif // APPLICATION_IOS_H