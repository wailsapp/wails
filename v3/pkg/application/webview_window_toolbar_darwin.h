//go:build darwin && !ios && !server

#ifndef WebviewWindowToolbarDarwin_h
#define WebviewWindowToolbarDarwin_h

#import <Cocoa/Cocoa.h>
#import <AppKit/NSSharingServicePickerToolbarItem.h>

extern void processToolbarItemClick(unsigned int itemID);
extern void processToolbarSearch(unsigned int itemID, char* query);
extern void processToolbarShareResult(unsigned int itemID, char* service, char* errorMessage);
// The caller owns the returned UTF-8 JSON string and must free it.
extern char* processToolbarShareData(unsigned int providerID, char* contentType);
extern void processToolbarShareProviderRelease(unsigned int providerID);

@interface WailsToolbarItem : NSToolbarItem
@property unsigned int itemID;
@end

@interface WailsToolbarSearchTarget : NSObject
@property unsigned int itemID;
- (void)handleSearch:(id)sender;
@end

@interface WailsToolbarGroupTarget : NSObject
@property (strong) NSArray<NSNumber*>* itemIDs;
- (void)handleClick:(id)sender;
@end

@interface WailsToolbarShareTarget : NSObject <NSSharingServicePickerToolbarItemDelegate, NSSharingServiceDelegate>
@property unsigned int itemID;
@property (strong) NSArray* items;
@property (copy) NSString* subject;
@property (assign) NSWindow* window;
@property (strong) NSSharingServicePicker* activePicker;
- (void)showSharePicker:(id)sender;
@end

@interface WailsToolbarShareProviderLifetime : NSObject
@property unsigned int providerID;
@end

// WailsToolbarDelegate owns the fully-built item list for one toolbar. Items
// are constructed before the delegate is installed so AppKit's first request
// for default identifiers observes a complete toolbar.
//
// orderedIdentifiers is the default layout; allowedIdentifiers is the full
// set offered by the customization palette (nil means "same as the default
// layout"); knownIdentifiers records every identifier ever synchronised into
// the live toolbar so a later structural sync can tell a newly added item
// from one the user removed through the palette.
@interface WailsToolbarDelegate : NSObject <NSToolbarDelegate>
@property (strong) NSMutableArray<NSToolbarItemIdentifier>* orderedIdentifiers;
@property (strong) NSArray<NSToolbarItemIdentifier>* allowedIdentifiers;
@property (strong) NSMutableSet<NSToolbarItemIdentifier>* knownIdentifiers;
@property (strong) NSMutableDictionary<NSToolbarItemIdentifier, NSToolbarItem*>* itemsByIdentifier;
@end

// Creates a detached toolbar handle. identifier becomes the NSToolbar
// identifier (the autosave key when customizable). The caller owns the handle
// until toolbarRelease.
void* toolbarCreate(const char* identifier, bool customizable);

// Toggles allowsUserCustomization and autosavesConfiguration on a live
// toolbar. The autosave identifier is fixed at toolbarCreate.
void toolbarSetCustomizable(void* handlePtr, bool customizable);

// Opens the standard "Customize Toolbar..." sheet.
void toolbarRunCustomizationPalette(void* handlePtr);

// Synchronises the toolbar layout with a JSON description produced by Go:
// {"customizable":bool,"default":[entry],"allowed":[entry],"centered":[entry],
//  "moved":"identifier"} where entry is {"id":"...","kind":"..."}. Kinds map
// standard AppKit identifiers ("space", "flexibleSpace", "sidebarToggle",
// "sidebarSeparator", "inspectorToggle", "inspectorSeparator"); "item" uses
// the id as-is. Before attachment only the delegate lists are updated. After
// attachment the live NSToolbar is diffed against the layout: a
// non-customizable toolbar mirrors the default list exactly, a customizable
// toolbar drops disallowed items, inserts newly added ones, replaces rebuilt
// items and honours "moved" while preserving the user's layout otherwise.
void toolbarSync(void* handlePtr, const char* layoutJSON);

// Whether NSToolbar.centeredItemIdentifiers is available on this system.
bool toolbarSupportsCenteredItems(void);

// Drops a Wails-owned item from the delegate's item table. The live toolbar
// is updated by the next toolbarSync.
void toolbarForgetItem(void* handlePtr, const char* identifier);

// Attaches a fully populated toolbar to the window.
void toolbarAttach(void* nsWindow, void* handlePtr, int style);

// Releases the caller-owned toolbar handle. The NSWindow independently retains
// an attached toolbar until it is replaced, detached, or destroyed.
void toolbarRelease(void* handlePtr);

// Removes the window's toolbar entirely.
void toolbarDetach(void* nsWindow);

// Sets NSToolbar.displayMode. Values map to NSToolbarDisplayMode.
void toolbarSetDisplayMode(void* handlePtr, int displayMode);

void* toolbarAddButtonItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool prominent, bool disabled, bool hidden,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount);

void* toolbarAddGroupItem(void* handlePtr, const char* identifier,
    const char* label, void** memberItems, int memberCount, int selectionMode, int selectedIndex);

void* toolbarBuildButtonItemStandalone(const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool disabled, bool hidden);

void* toolbarAddSearchItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* tooltip, bool disabled, bool hidden);

void* toolbarAddShareItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool disabled, bool hidden, const char* providerJSON);

void toolbarAddFlexibleSpaceIdentifier(void* handlePtr);
void toolbarAddSpaceIdentifier(void* handlePtr);

// Adds an NSMenuToolbarItem (macOS 10.15+) showing nsMenu. On earlier
// releases the item is omitted and a message is logged.
void* toolbarAddMenuItem(void* handlePtr, const char* identifier, const char* label,
    const char* symbolName, const char* tooltip, bool bordered, bool disabled, bool hidden,
    void* nsMenu, bool showsIndicator);
void toolbarMenuItemSetShowsIndicator(void* handlePtr, const char* identifier, bool showsIndicator);

void toolbarItemSetVisibilityPriority(void* handlePtr, const char* identifier, int priority);
// macOS 11+; no-op on earlier releases.
void toolbarItemSetNavigational(void* handlePtr, const char* identifier, bool navigational);

// Search-field configuration. Each targets the NSSearchToolbarItem's
// searchField (macOS 11+) or the NSSearchField view of the fallback item.
void toolbarSearchItemSetPlaceholder(void* handlePtr, const char* identifier, const char* placeholder);
void toolbarSearchItemSetIncremental(void* handlePtr, const char* identifier, bool incremental);
void toolbarSearchItemSetRecents(void* handlePtr, const char* identifier, const char* autosaveName, int maximumRecents);
// nsMenu may be NULL to clear the custom menu. includeRecents appends the
// standard recent-searches section (title, recents, clear) to the template.
void toolbarSearchItemSetMenu(void* handlePtr, const char* identifier, void* nsMenu, bool includeRecents);

// Standard AppKit sidebar items. AppKit creates and owns these; they never
// pass through the delegate's item factory.
void toolbarAddSidebarToggleIdentifier(void* handlePtr);
// No-op on macOS releases without NSToolbarSidebarTrackingSeparatorItemIdentifier.
void toolbarAddSidebarTrackingSeparatorIdentifier(void* handlePtr);
// Uses NSToolbarToggleInspectorItemIdentifier on macOS 14+. Earlier releases
// receive a Wails-owned native item routed through itemID.
void toolbarAddInspectorToggleItem(void* handlePtr, const char* identifier, unsigned int itemID);
// No-op on macOS releases without NSToolbarInspectorTrackingSeparatorItemIdentifier.
void toolbarAddInspectorTrackingSeparatorIdentifier(void* handlePtr);

void toolbarItemSetLabel(void* handlePtr, const char* identifier, const char* label);
void toolbarItemSetSymbol(void* handlePtr, const char* identifier, const char* symbolName);
void toolbarItemSetTooltip(void* handlePtr, const char* identifier, const char* tooltip);
void toolbarItemSetBordered(void* handlePtr, const char* identifier, bool bordered);
void toolbarItemSetProminent(void* handlePtr, const char* identifier, bool prominent);
void toolbarItemSetTintColor(void* handlePtr, const char* identifier, bool hasTint,
    double tintR, double tintG, double tintB, double tintA);
void toolbarItemSetEnabled(void* handlePtr, const char* identifier, bool enabled);
void toolbarItemSetHidden(void* handlePtr, const char* identifier, bool hidden);
void toolbarItemSetBadgeCount(void* handlePtr, const char* identifier, int badgeCount);
void toolbarGroupSetSelectedIndex(void* handlePtr, const char* identifier, int index);
void toolbarGroupSetSelectionMode(void* handlePtr, const char* identifier, int selectionMode);
void toolbarShareItemSetProvider(void* handlePtr, const char* identifier, const char* providerJSON);

#endif /* WebviewWindowToolbarDarwin_h */
