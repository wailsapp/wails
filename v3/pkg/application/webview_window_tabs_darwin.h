//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowTabs_h
#define WailsWebviewWindowTabs_h

#import <Cocoa/Cocoa.h>

// All functions expect to run on the main thread; the Go side dispatches
// through InvokeSync. Every function tolerates a NULL window.

enum {
    WailsWindowTabsAdded = 0,
    WailsWindowTabsInvalidWindow = -1,
    WailsWindowTabsInvalidOrder = -2,
};

// Calls addTabbedWindow:ordered: with an NSWindowOrderingMode value
// (NSWindowAbove = 1, NSWindowBelow = -1).
int windowTabsAdd(void* nsWindow, void* other, int ordered);

// Whether the window currently belongs to an NSWindowTabGroup.
bool windowTabsHasGroup(void* nsWindow);

// The unretained NSWindowTabGroup pointer, used only for identity checks.
void* windowTabsGroupPointer(void* nsWindow);

// Number of windows in the group, and the unretained NSWindow at an index.
int windowTabsGroupCount(void* nsWindow);
void* windowTabsGroupWindowAt(void* nsWindow, int index);

void* windowTabsSelectedWindow(void* nsWindow);
void windowTabsSelectNext(void* nsWindow);
void windowTabsSelectPrevious(void* nsWindow);
void windowTabsSelect(void* nsWindow, void* target);

bool windowTabsIsTabBarVisible(void* nsWindow);
void windowTabsToggleTabBar(void* nsWindow);
bool windowTabsIsOverviewVisible(void* nsWindow);
void windowTabsToggleOverview(void* nsWindow);

// Returns a malloc'd copy of the tabbing identifier; the caller frees it.
char* windowTabsIdentifier(void* nsWindow);

void windowTabsMoveToNewWindow(void* nsWindow);
void windowTabsMergeAllWindows(void* nsWindow);

// An empty string restores the default tab title (the window title) or
// removes the tooltip.
void windowTabsSetTitle(void* nsWindow, const char* title);
void windowTabsSetToolTip(void* nsWindow, const char* tooltip);

#endif
