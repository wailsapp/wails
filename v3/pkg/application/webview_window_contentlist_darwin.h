//go:build darwin && !ios && !server

#ifndef WebviewWindowContentListDarwin_h
#define WebviewWindowContentListDarwin_h

#import <Cocoa/Cocoa.h>

// Content-list callbacks (implemented in webview_window_contentlist_darwin.go).
// rowIDs lists the selection in row order.
extern void processMacContentListSelectionChanged(unsigned long long paneID, unsigned long long* rowIDs, int count);
extern void processMacContentListRowActivated(unsigned long long paneID, unsigned long long rowID);
extern void processMacContentListSortChanged(unsigned long long paneID, int column, bool ascending);
// Returns the NSMenu to show for a right-click on rowID (0 for the empty
// area), or NULL. Runs synchronously on the application thread.
extern void* processMacContentListContextMenu(unsigned long long paneID, unsigned long long rowID);

// Hooks used by webview_window_split_darwin.m to host the table in a
// content-list split item. The controller is created lazily so rows can be
// staged before the window exists.
NSViewController* wailsContentListCreateController(unsigned long long paneID);
// Loads the view with the pane surface colour. Returns false on failure.
bool wailsContentListPrepareForInstall(NSViewController* controller, NSColor* surfaceColor);
// Builds the NSSplitViewItem: contentListWithViewController: on macOS 11 and
// newer, a regular item otherwise.
NSSplitViewItem* wailsContentListSplitItem(NSViewController* controller);
// Reloads the table once the split item is in the window.
void wailsContentListDidInstall(NSViewController* controller);

// splitViewContentListController returns the pane's (lazily created) table
// controller, or nil when the pane is not a content list. Implemented in
// webview_window_split_darwin.m because the pane records live there.
NSViewController* splitViewContentListController(void* handlePtr, unsigned long long paneID);

// Presentation options for one content list. style mirrors
// MacContentListStyle; sortColumn is -1 when unsorted.
typedef struct WailsContentListOptions {
    bool allowsMultipleSelection;
    bool sortable;
    int sortColumn;
    bool sortAscending;
    double rowHeight;
    int style;
    bool alternatingRows;
    const char* emptyText;
    bool headerVisible;
} WailsContentListOptions;

// One row's presentation. cellsJSON is a UTF-8 JSON string array for table
// mode. Strings are borrowed for the duration of the call.
typedef struct WailsContentListRowSpec {
    const char* title;
    const char* subtitle;
    const char* detail;
    const char* symbol;
    const char* tooltip;
    int badge;
    const char* cellsJSON;
    bool disabled;
    bool hidden;
} WailsContentListRowSpec;

// Native contents for a content-list pane. Reset drops columns and rows;
// AddColumn/AddRow stage the model in order; Reload rebuilds the table with
// the selection preserved. UpdateRow refreshes one row in place.
void splitViewContentListReset(void* handlePtr, unsigned long long paneID);
void splitViewContentListSetOptions(void* handlePtr, unsigned long long paneID, WailsContentListOptions options);
void splitViewContentListAddColumn(void* handlePtr, unsigned long long paneID,
    const char* title, double width, double minWidth, double maxWidth, bool sortable, int alignment);
void splitViewContentListAddRow(void* handlePtr, unsigned long long paneID,
    unsigned long long rowID, WailsContentListRowSpec spec);
void splitViewContentListSetSelection(void* handlePtr, unsigned long long paneID,
    unsigned long long* rowIDs, int count);
void splitViewContentListReload(void* handlePtr, unsigned long long paneID);
void splitViewContentListUpdateRow(void* handlePtr, unsigned long long paneID,
    unsigned long long rowID, WailsContentListRowSpec spec);

#endif /* WebviewWindowContentListDarwin_h */
