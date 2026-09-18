//go:build darwin && !ios && !server

#ifndef WailsNativeWindowDarwin_h
#define WailsNativeWindowDarwin_h

#import <Cocoa/Cocoa.h>

extern void processNativeWindowClosed(unsigned int windowID);

// WailsNativeWindowConfig carries the MacWindow options that must be known
// when the NSWindow (or NSPanel) is constructed. Options that AppKit accepts
// after construction are applied through mac_window_chrome_darwin.h.
typedef struct {
    unsigned int windowID;
    int width;
    int height;
    bool hideOnClose;
    // frameless is TitleBar.Hide. borderless and cornerRadius come from
    // resolveNativeMacFrame and follow WebviewWindow's frameless style mask.
    bool frameless;
    bool borderless;
    double cornerRadius;
    // isPanel selects an NSPanel subclass; the four preferences mirror
    // MacPanelPreferences.
    bool isPanel;
    bool floatingPanel;
    bool becomesKeyOnlyIfNeeded;
    bool nonActivating;
    bool utilityWindow;
    bool disableEscapeExitsFullscreen;
} WailsNativeWindowConfig;

// WailsNativeBackdropConfig describes the surface applied once the split
// layout is installed. backdrop is a MacBackdrop value; the remaining fields
// are the MacLiquidGlass options and are read only for the Liquid Glass backdrop.
typedef struct {
    int backdrop;
    int glassStyle;
    int glassMaterial;
    double glassCornerRadius;
    int tintR;
    int tintG;
    int tintB;
    int tintA;
    const char* groupID;
    double groupSpacing;
} WailsNativeBackdropConfig;

void* nativeWindowCreate(WailsNativeWindowConfig config);
void nativeWindowDestroy(void* nsWindow);
void nativeWindowSetTitle(void* nsWindow, const char* title);
void nativeWindowSetResizable(void* nsWindow, bool resizable);
void nativeWindowSetMinSize(void* nsWindow, int width, int height);
void nativeWindowSetMaxSize(void* nsWindow, int width, int height);
void nativeWindowConfigureTitlebar(void* nsWindow, bool appearsTransparent,
    bool fullSizeContent, bool hideTitle, bool hideToolbarSeparator, int toolbarStyle);
// nativeWindowApplyContentChrome applies the options that depend on the
// installed content hierarchy: the content layout (belowToolbar removes the
// full-size content style so the whole split layout sits under the toolbar),
// the custom rounded-corner mask, and the backdrop. It must run after the
// split view has replaced the window's content view.
void nativeWindowApplyContentChrome(void* nsWindow, bool belowToolbar,
    double cornerRadius, WailsNativeBackdropConfig backdrop);
void nativeWindowShow(void* nsWindow);
void nativeWindowHide(void* nsWindow);
void nativeWindowFocus(void* nsWindow);
bool nativeWindowIsVisible(void* nsWindow);
void nativeWindowCenter(void* nsWindow);
void nativeWindowSetPosition(void* nsWindow, int x, int y);

#endif
