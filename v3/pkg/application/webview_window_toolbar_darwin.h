//go:build darwin && !ios && !server

#ifndef WebviewWindowToolbarDarwin_h
#define WebviewWindowToolbarDarwin_h

#import <Cocoa/Cocoa.h>

extern void processToolbarItemClick(unsigned int itemID);
extern void processToolbarSearch(unsigned int itemID, char* query);

@interface WailsToolbarItem : NSToolbarItem
@property unsigned int itemID;
@end

@interface WailsSearchToolbarItem : NSSearchToolbarItem
@property unsigned int itemID;
@end

// WailsToolbarDelegate owns the fully-built item list for one toolbar.
// Items are constructed once, up front, by the Go-side tree walk (see
// webview_window_toolbar_darwin.go); the delegate just serves them back.
@interface WailsToolbarDelegate : NSObject <NSToolbarDelegate>
@property (strong) NSMutableArray<NSToolbarItemIdentifier>* orderedIdentifiers;
@property (strong) NSMutableDictionary<NSToolbarItemIdentifier, NSToolbarItem*>* itemsByIdentifier;
@end

// Creates a new, empty toolbar + delegate pair and attaches it to the
// window, replacing any existing toolbar. Returns the delegate (retained
// by the window's toolbar association; the caller does not need to keep it
// alive beyond populating it).
void* toolbarNewAndAttach(void* nsWindow);

// Removes the window's toolbar entirely.
void toolbarDetach(void* nsWindow);

// Appends a plain button item. tintR/G/B/A are ignored unless hasTint is
// true. Returns the constructed NSToolbarItem (also retained by the
// delegate's itemsByIdentifier).
void* toolbarAddButtonItem(void* delegatePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, bool bordered, bool prominent,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount);

// Appends a segmented group. memberItems must be items built via
// toolbarBuildButtonItemStandalone: each already carries its own
// target/action, so clicking a segment independently dispatches that
// segment's own itemID, the group itself has no target/action of its own.
void* toolbarAddGroupItem(void* delegatePtr, const char* identifier,
    const char* label, void** memberItems, int memberCount);

// Builds a standalone button item not yet attached to any delegate's
// identifier list, for use as a group member.
void* toolbarBuildButtonItemStandalone(const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, bool bordered);

// Appends a search field item.
void* toolbarAddSearchItem(void* delegatePtr, const char* identifier, unsigned int itemID, const char* label);

// Appends the system flexible-space identifier (no item construction needed).
void toolbarAddFlexibleSpaceIdentifier(void* delegatePtr);

// Appends a sidebar-toggle button. Wired via the standard AppKit responder
// chain (target nil, action toggleSidebar:), so it works against whichever
// NSSplitViewController is the window's contentViewController without a
// direct reference. Tracking-separator alignment with a specific pane's
// divider is added separately by the split-window implementation.
void toolbarAddSidebarToggleIdentifier(void* delegatePtr, const char* identifier);

#endif /* WebviewWindowToolbarDarwin_h */
