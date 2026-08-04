//go:build darwin && !ios && !server

#ifndef WebviewWindowToolbarDarwin_h
#define WebviewWindowToolbarDarwin_h

#import <Cocoa/Cocoa.h>

extern void processToolbarItemClick(unsigned int itemID);
extern void processToolbarSearch(unsigned int itemID, char* query);

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

// WailsToolbarDelegate owns the fully-built item list for one toolbar. Items
// are constructed before the delegate is installed so AppKit's first request
// for default identifiers observes a complete toolbar.
@interface WailsToolbarDelegate : NSObject <NSToolbarDelegate>
@property (strong) NSMutableArray<NSToolbarItemIdentifier>* orderedIdentifiers;
@property (strong) NSMutableDictionary<NSToolbarItemIdentifier, NSToolbarItem*>* itemsByIdentifier;
@end

// Creates a detached toolbar handle with a unique internal identifier. The
// caller owns the handle until toolbarRelease.
void* toolbarCreate(const char* identifier);

// Attaches a fully populated toolbar to the window.
void toolbarAttach(void* nsWindow, void* handlePtr, int style);

// Releases the caller-owned toolbar handle. The NSWindow independently retains
// an attached toolbar until it is replaced, detached, or destroyed.
void toolbarRelease(void* handlePtr);

// Removes the window's toolbar entirely.
void toolbarDetach(void* nsWindow);

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

void toolbarAddFlexibleSpaceIdentifier(void* handlePtr);

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

#endif /* WebviewWindowToolbarDarwin_h */
