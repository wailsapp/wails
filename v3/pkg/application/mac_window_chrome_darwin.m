//go:build darwin && !ios && !server

#import "mac_window_chrome_darwin.h"

void windowSetToolbarStyle(void* pointer, int style) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil || window.toolbar == nil) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) window.toolbarStyle = (NSWindowToolbarStyle)style;
#endif
}

void windowSetHideToolbarSeparator(void* pointer, bool hideSeparator) {
    NSWindow* window = (NSWindow*)pointer;
    if (window == nil || window.toolbar == nil) return;
    window.toolbar.showsBaselineSeparator = !hideSeparator;
}

void windowSetShowToolbarWhenFullscreen(void* pointer, bool setting) {
    NSWindow* window = (NSWindow*)pointer;
    id delegate = window.delegate;
    if (delegate != nil && [delegate respondsToSelector:@selector(setShowToolbarWhenFullscreen:)]) {
        [delegate setValue:@(setting) forKey:@"showToolbarWhenFullscreen"];
    }
}
