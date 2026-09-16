//go:build darwin && !ios && !server

#import "webview_window_tabs_darwin.h"
#include <string.h>

static NSWindowTabGroup* tabGroupForWindow(void* nsWindow) {
    if (nsWindow == NULL) return nil;
    return [(NSWindow*)nsWindow tabGroup];
}

// selectNextTab: and friends act relative to the receiver, so route them
// through the group's selected window to get "next after the current tab".
static NSWindow* actionWindowForGroup(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group != nil && group.selectedWindow != nil) return group.selectedWindow;
    return (NSWindow*)nsWindow;
}

int windowTabsAdd(void* nsWindow, void* other, int ordered) {
    if (nsWindow == NULL || other == NULL || nsWindow == other) {
        return WailsWindowTabsInvalidWindow;
    }
    if (ordered != NSWindowAbove && ordered != NSWindowBelow) {
        return WailsWindowTabsInvalidOrder;
    }
    [(NSWindow*)nsWindow addTabbedWindow:(NSWindow*)other ordered:(NSWindowOrderingMode)ordered];
    return WailsWindowTabsAdded;
}

bool windowTabsHasGroup(void* nsWindow) {
    return tabGroupForWindow(nsWindow) != nil;
}

void* windowTabsGroupPointer(void* nsWindow) {
    return (void*)tabGroupForWindow(nsWindow);
}

int windowTabsGroupCount(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group == nil) return 0;
    return (int)group.windows.count;
}

void* windowTabsGroupWindowAt(void* nsWindow, int index) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group == nil || index < 0 || (NSUInteger)index >= group.windows.count) return NULL;
    return (void*)group.windows[(NSUInteger)index];
}

void* windowTabsSelectedWindow(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group == nil) return NULL;
    return (void*)group.selectedWindow;
}

void windowTabsSelectNext(void* nsWindow) {
    if (nsWindow == NULL) return;
    [actionWindowForGroup(nsWindow) selectNextTab:nil];
}

void windowTabsSelectPrevious(void* nsWindow) {
    if (nsWindow == NULL) return;
    [actionWindowForGroup(nsWindow) selectPreviousTab:nil];
}

void windowTabsSelect(void* nsWindow, void* target) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group == nil || target == NULL) return;
    NSWindow* window = (NSWindow*)target;
    if (![group.windows containsObject:window]) return;
    group.selectedWindow = window;
}

bool windowTabsIsTabBarVisible(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    return group != nil && group.tabBarVisible;
}

void windowTabsToggleTabBar(void* nsWindow) {
    if (nsWindow == NULL) return;
    [(NSWindow*)nsWindow toggleTabBar:nil];
}

bool windowTabsIsOverviewVisible(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    return group != nil && group.overviewVisible;
}

void windowTabsToggleOverview(void* nsWindow) {
    if (nsWindow == NULL) return;
    [(NSWindow*)nsWindow toggleTabOverview:nil];
}

char* windowTabsIdentifier(void* nsWindow) {
    NSWindowTabGroup* group = tabGroupForWindow(nsWindow);
    if (group == nil || group.identifier == nil) return strdup("");
    return strdup([group.identifier UTF8String]);
}

void windowTabsMoveToNewWindow(void* nsWindow) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    // AppKit only detaches windows that share a group with someone else.
    if (window.tabGroup == nil || window.tabGroup.windows.count < 2) return;
    [window moveTabToNewWindow:nil];
}

void windowTabsMergeAllWindows(void* nsWindow) {
    if (nsWindow == NULL) return;
    [(NSWindow*)nsWindow mergeAllWindows:nil];
}

void windowTabsSetTitle(void* nsWindow, const char* title) {
    if (nsWindow == NULL) return;
    NSWindowTab* tab = [(NSWindow*)nsWindow tab];
    if (tab == nil) return;
    if (title == NULL || title[0] == '\0') {
        tab.title = nil;
        return;
    }
    tab.title = [NSString stringWithUTF8String:title];
}

void windowTabsSetToolTip(void* nsWindow, const char* tooltip) {
    if (nsWindow == NULL) return;
    NSWindowTab* tab = [(NSWindow*)nsWindow tab];
    if (tab == nil) return;
    if (tooltip == NULL || tooltip[0] == '\0') {
        tab.toolTip = nil;
        return;
    }
    tab.toolTip = [NSString stringWithUTF8String:tooltip];
}
