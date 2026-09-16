//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowAccessory_h
#define WailsWebviewWindowAccessory_h

#import <Cocoa/Cocoa.h>

enum {
    WailsMacAccessoryStyleApplied = 0,
    WailsMacAccessoryStyleUnavailable = -1,
    WailsMacAccessoryStyleInvalidController = -2,
};

// Returns 1 for NSTitlebarAccessoryViewController, 2 for
// NSSplitViewItemAccessoryViewController, and 0 for any other object.
int macAccessoryViewControllerKind(void* controller);

// Returns whether preferredScrollEdgeEffectStyle can be used on this
// controller on the running macOS version.
bool macAccessoryViewControllerSupportsScrollEdgeEffectStyle(void* controller);

// Applies a MacScrollEdgeEffectStyle value. Automatic is a successful no-op
// on systems older than macOS 26.1; explicit styles report unavailable.
int macAccessoryViewControllerSetScrollEdgeEffectStyle(void* controller, int style);

// Returns a MacScrollEdgeEffectStyle value, or one of the negative result
// constants above.
int macAccessoryViewControllerScrollEdgeEffectStyle(void* controller);

// Native control strips (MacAccessory). Every function below must run on the
// application thread. Controller kinds and layouts mirror the Go constants.
enum {
    WailsMacAccessoryKindTitlebar = 1,
    WailsMacAccessoryKindSplitItem = 2,
};

enum {
    WailsMacAccessoryLayoutLeading = 0,
    WailsMacAccessoryLayoutTrailing = 1,
    WailsMacAccessoryLayoutBottom = 2,
    WailsMacAccessoryLayoutTop = 3,
};

extern void processMacAccessoryClicked(unsigned long long controlID);
extern void processMacAccessorySearched(unsigned long long controlID, char* query);
extern void processMacAccessorySegmentChanged(unsigned long long controlID, int selectedIndex);

// Reports whether NSSplitViewItem accessories exist on the running system.
bool macAccessorySplitItemAccessoriesSupported(void);

// macAccessoryCreate builds an accessory view controller of the given kind
// holding an empty horizontal control strip of the given height. The caller
// owns one reference and releases it with macAccessoryRelease after
// detaching. Returns NULL when the kind is unavailable.
void* macAccessoryCreate(int kind, int layout, double height);
void macAccessoryRelease(void* controller);

bool macAccessoryAttachToWindow(void* nsWindow, void* controller);
bool macAccessoryAttachToSplitItem(void* splitItem, void* controller, bool top);
// Detaches from whichever window currently hosts the controller.
void macAccessoryDetachFromWindow(void* controller);
void macAccessoryDetachFromSplitItem(void* splitItem, void* controller);

void macAccessorySetHidden(void* controller, bool hidden);
void macAccessorySetHeight(void* controller, double height);
void macAccessorySetFullScreenMinHeight(void* controller, double height);
void macAccessorySetAutomaticallyAdjustsSize(void* controller, bool adjusts);
void macAccessorySetAppliesContentInsets(void* controller, bool applies);

// Controls are appended in call order. Strings are borrowed for the call;
// labelsJSON and symbolsJSON are UTF-8 JSON string arrays. width 0 fits.
void macAccessoryAddSearch(void* controller, unsigned long long controlID,
    const char* placeholder, const char* text, bool incremental,
    const char* tooltip, bool disabled, bool hidden, double width);
void macAccessoryAddSegmented(void* controller, unsigned long long controlID,
    const char* labelsJSON, const char* symbolsJSON, int selected,
    const char* tooltip, bool disabled, bool hidden, double width);
// nsMenu non-NULL makes a menu button that pops the menu up when clicked.
void macAccessoryAddButton(void* controller, unsigned long long controlID,
    const char* title, const char* symbol, void* nsMenu,
    const char* tooltip, bool disabled, bool hidden, double width);
void macAccessoryAddLabel(void* controller, unsigned long long controlID,
    const char* text, const char* symbol, const char* tooltip, bool hidden, double width);
void macAccessoryAddFlexibleSpace(void* controller, unsigned long long controlID);
void macAccessoryAddNativeView(void* controller, unsigned long long controlID,
    void* nsView, bool hidden, double width);

void macAccessoryControlSetEnabled(void* controller, unsigned long long controlID, bool enabled);
void macAccessoryControlSetHidden(void* controller, unsigned long long controlID, bool hidden);
void macAccessoryControlSetTooltip(void* controller, unsigned long long controlID, const char* tooltip);
void macAccessoryControlSetWidth(void* controller, unsigned long long controlID, double width);
void macAccessoryControlSetText(void* controller, unsigned long long controlID, const char* text);
void macAccessoryControlSetSymbol(void* controller, unsigned long long controlID, const char* symbol);
void macAccessoryControlSetPlaceholder(void* controller, unsigned long long controlID, const char* placeholder);
void macAccessoryControlSetIncremental(void* controller, unsigned long long controlID, bool incremental);
void macAccessoryControlSetSegments(void* controller, unsigned long long controlID,
    const char* labelsJSON, const char* symbolsJSON, int selected);
void macAccessoryControlSetSelectedSegment(void* controller, unsigned long long controlID, int selected);
void macAccessoryControlSetMenu(void* controller, unsigned long long controlID, void* nsMenu);

#endif
