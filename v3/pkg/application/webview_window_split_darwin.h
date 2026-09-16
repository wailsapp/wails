//go:build darwin && !ios && !server

#ifndef WebviewWindowSplitDarwin_h
#define WebviewWindowSplitDarwin_h

#import <Cocoa/Cocoa.h>
#ifndef WAILS_NATIVE_ONLY
#import <WebKit/WebKit.h>
#import "webview_window_darwin.h"
#endif

extern void processMacSplitPaneLoaded(unsigned long long paneID);
extern void processMacSplitPaneNavigationStarted(unsigned long long paneID);
extern void processMacSplitPaneCollapsed(unsigned long long paneID, bool collapsed);
extern void processMacSidebarItemSelected(unsigned long long itemID);
extern void processMacInspectorTextChanged(unsigned long long controlID, char* value);
extern void processMacInspectorToggleChanged(unsigned long long controlID, bool checked);
extern void processMacInspectorSelectionChanged(unsigned long long controlID, int selectedIndex);
extern void processMacTextEditorChanged(unsigned long long editorID);

// Sidebar callbacks (implemented in webview_window_sidebar_darwin.go).
// itemIDs lists the clicked row first, then the rest of the selection.
extern void processMacSidebarSelectionChanged(unsigned long long paneID, unsigned long long* itemIDs, int count);
extern void processMacSidebarItemExpanded(unsigned long long itemID, bool expanded);
extern void processMacSidebarItemRenamed(unsigned long long itemID, char* label);
extern void processMacSidebarItemMoved(unsigned long long paneID, unsigned long long itemID,
    unsigned long long parentID, int index);
// Returns the NSMenu to show for a right-click on itemID (0 for the empty
// area), or NULL. Runs synchronously on the application thread.
extern void* processMacSidebarContextMenu(unsigned long long paneID, unsigned long long itemID);

// Inspector callbacks (implemented in webview_window_inspector_darwin.go).
extern void processMacInspectorNumberChanged(unsigned long long controlID, int kind, double value);
extern void processMacInspectorSegmentChanged(unsigned long long controlID, int selectedIndex);
extern void processMacInspectorColorChanged(unsigned long long controlID, int red, int green, int blue, int alpha);
extern void processMacInspectorDateChanged(unsigned long long controlID, double unixSeconds);
extern void processMacInspectorButtonClicked(unsigned long long controlID);
extern void processMacInspectorSectionCollapsed(unsigned long long sectionID, bool collapsed);

// splitPrimaryPaneIDForWebView returns the primary split-pane ID associated
// with the window's preserved WebView, or 0 when no split is installed.
#ifndef WAILS_NATIVE_ONLY
unsigned long long splitPrimaryPaneIDForWebView(WKWebView* webView);
#else
unsigned long long splitPrimaryPaneIDForWebView(void* webView);
#endif

// Pane roles mirror macSplitPaneRole on the Go side.
enum {
    WailsSplitPaneRoleSidebar = 0,
    WailsSplitPaneRolePrimary = 1,
    WailsSplitPaneRoleInspector = 2,
};

// splitViewCreate allocates an owning handle for one split layout. The caller
// owns the handle until splitViewRelease.
void* splitViewCreate(const char* autosaveName);

// splitViewAddPane appends one pane description in leading-to-trailing order.
// Values guarded by a has* flag are applied only when explicitly configured
// so AppKit's semantic role defaults are preserved otherwise.
void splitViewAddPane(void* handlePtr, unsigned long long paneID, int role, bool primary,
    double minThickness, double maxThickness,
    double preferredFraction, bool hasPreferredFraction,
    double holdingPriority, bool hasHoldingPriority,
    bool collapsible, bool hasCollapsible,
    bool canCollapseFromResize, bool hasCanCollapseFromResize,
    bool startCollapsed, int contentLayout);

// splitViewInstall builds the NSSplitViewController, reparents the window's
// primary WKWebView into the pane marked primary, creates native AppKit
// sidebar controllers, and installs the controller as the window's content
// view controller. A normal backdrop is made opaque; explicit transparent,
// translucent, and Liquid Glass backdrops retain their configured behavior.
// Must run on the application thread before any toolbar.
bool splitViewInstall(void* handlePtr, void* nsWindow, bool normalBackdrop);

// splitViewInstallNative installs the same native sidebar/inspector layout
// into an ordinary NSWindow and uses the configured NSTextView pane as its
// primary content. No WKWebView is created.
bool splitViewInstallNative(void* handlePtr, void* nsWindow, bool normalBackdrop);

void splitViewConfigureTextEditor(void* handlePtr, unsigned long long paneID,
    unsigned long long editorID, const char* text, bool editable);
void splitViewTextEditorSetText(void* handlePtr, unsigned long long paneID, const char* text);
// The caller owns the returned UTF-8 string and must free it.
char* splitViewTextEditorCopyText(void* handlePtr, unsigned long long paneID);
void splitViewTextEditorSetEditable(void* handlePtr, unsigned long long paneID, bool editable);
void splitViewTextEditorFocus(void* handlePtr, unsigned long long paneID);

// splitViewTeardown removes observers and drops the owner's native references. It never
// detaches the window's content view controller: the window keeps the view
// hierarchy (including the preserved primary WebView) alive until normal
// window teardown. Idempotent.
void splitViewTeardown(void* handlePtr);

// splitViewRelease frees the caller-owned handle after teardown.
void splitViewRelease(void* handlePtr);

void splitViewPaneSetMinimumThickness(void* handlePtr, unsigned long long paneID, double value);
void splitViewPaneSetMaximumThickness(void* handlePtr, unsigned long long paneID, double value);
void splitViewPaneSetPreferredFraction(void* handlePtr, unsigned long long paneID, double value);
void splitViewPaneSetHoldingPriority(void* handlePtr, unsigned long long paneID, double value);
void splitViewPaneSetCollapsible(void* handlePtr, unsigned long long paneID, bool collapsible);
void splitViewPaneSetCanCollapseFromWindowResize(void* handlePtr, unsigned long long paneID, bool allowed);
void splitViewPaneSetContentLayout(void* handlePtr, unsigned long long paneID, int layout);
void splitViewPaneSetCollapsed(void* handlePtr, unsigned long long paneID, bool collapsed);
void splitViewPaneToggleCollapsed(void* handlePtr, unsigned long long paneID);

// Native source-list contents for a sidebar pane. Sections are nonselectable;
// item identifiers are generated by Go and used only to route callbacks.
// Rows are added beneath parentID: a section, another row, or 0 for the root.
void splitViewSidebarReset(void* handlePtr, unsigned long long paneID);
void splitViewSidebarSetOptions(void* handlePtr, unsigned long long paneID,
    bool allowsMultipleSelection, bool reorderable);
void splitViewSidebarAddSection(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, const char* label);
void splitViewSidebarAddItem(void* handlePtr, unsigned long long paneID,
    unsigned long long parentID, unsigned long long itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool disabled, bool hidden, bool expanded, bool editable, int badge,
    const char* accessorySymbol, bool hasTint, int red, int green, int blue, int alpha);
void splitViewSidebarSetSelectedItem(void* handlePtr, unsigned long long paneID,
    unsigned long long itemID);

// WailsInspectorControlSpec carries one control's presentation. kind mirrors
// MacInspectorControlKind. Options are a UTF-8 JSON string array. Strings are
// borrowed for the duration of the call.
typedef struct WailsInspectorControlSpec {
    int kind;
    const char* label;
    const char* value;
    bool checked;
    const char* optionsJSON;
    int selectedIndex;
    const char* tooltip;
    bool disabled;
    bool hidden;
    double number;
    double minimum;
    double maximum;
    double step;
    int red;
    int green;
    int blue;
    int alpha;
    double dateSeconds;
} WailsInspectorControlSpec;

// Native property-inspector contents.
void splitViewInspectorReset(void* handlePtr, unsigned long long paneID);
void splitViewInspectorAddSection(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, const char* label, bool collapsible, bool collapsed);
void splitViewInspectorAddControl(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, unsigned long long controlID,
    WailsInspectorControlSpec spec);
void splitViewInspectorReload(void* handlePtr, unsigned long long paneID);
void splitViewInspectorUpdateControl(void* handlePtr, unsigned long long paneID,
    unsigned long long controlID, WailsInspectorControlSpec spec);

#endif /* WebviewWindowSplitDarwin_h */
