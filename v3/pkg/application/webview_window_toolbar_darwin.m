#import "webview_window_toolbar_darwin.h"
#import <objc/runtime.h>
#import <string.h>

@implementation WailsToolbarItem

- (void)handleClick {
    processToolbarItemClick(self.itemID);
}

@end

@implementation WailsSearchToolbarItem

- (void)handleSearch:(id)sender {
    NSSearchField* field = (NSSearchField*)sender;
    char* query = (char*)[field.stringValue UTF8String];
    processToolbarSearch(self.itemID, query);
}

@end

@implementation WailsToolbarDelegate

- (instancetype)init {
    self = [super init];
    if (self) {
        _orderedIdentifiers = [[NSMutableArray alloc] init];
        _itemsByIdentifier = [[NSMutableDictionary alloc] init];
    }
    return self;
}

- (void)dealloc {
    [_orderedIdentifiers release];
    [_itemsByIdentifier release];
    [super dealloc];
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarDefaultItemIdentifiers:(NSToolbar*)toolbar {
    return self.orderedIdentifiers;
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarAllowedItemIdentifiers:(NSToolbar*)toolbar {
    return self.orderedIdentifiers;
}

- (NSToolbarItem*)toolbar:(NSToolbar*)toolbar
        itemForItemIdentifier:(NSToolbarItemIdentifier)itemIdentifier
    willBeInsertedIntoToolbar:(BOOL)flag {
    return self.itemsByIdentifier[itemIdentifier];
}

@end

static void applyCommonItemStyle(NSToolbarItem* item, bool bordered, bool prominent,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    if (@available(macOS 10.15, *)) {
        item.bordered = bordered;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        if (prominent) {
            item.style = NSToolbarItemStyleProminent;
        }
        if (hasTint) {
            item.backgroundTintColor = [NSColor colorWithRed:tintR green:tintG blue:tintB alpha:tintA];
        }
        if (badgeCount > 0) {
            item.badge = [NSItemBadge badgeWithCount:badgeCount];
        }
    }
#endif
}

void* toolbarNewAndAttach(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    WailsToolbarDelegate* delegate = [[[WailsToolbarDelegate alloc] init] autorelease];
    NSToolbar* toolbar = [[[NSToolbar alloc] initWithIdentifier:@"wails.toolbar"] autorelease];
    toolbar.delegate = delegate;
    toolbar.displayMode = NSToolbarDisplayModeIconAndLabel;
    // NSToolbar.delegate is a weak/assign reference: without this, the
    // delegate has no owner once this function returns and gets
    // deallocated, silently breaking every item lookup. Associate it with
    // RETAIN so it lives exactly as long as the toolbar does.
    objc_setAssociatedObject(toolbar, "wailsToolbarDelegate", delegate, OBJC_ASSOCIATION_RETAIN);
    window.toolbar = toolbar;
    return (void*)delegate;
}

void toolbarDetach(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    window.toolbar = nil;
}

// Returns a +1 retained item: the caller (toolbarAddButtonItem, or
// toolbarAddGroupItem via the member array) is responsible for handing it
// to something that retains it (an NSDictionary/NSArray) and then
// releasing this +1.
void* toolbarBuildButtonItemStandalone(const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, bool bordered) {
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = itemID;
    item.label = [NSString stringWithUTF8String:label];
    item.target = item;
    item.action = @selector(handleClick);
    if (@available(macOS 10.15, *)) {
        item.bordered = bordered;
    }
    if (symbolName != NULL && strlen(symbolName) > 0) {
        if (@available(macOS 11.0, *)) {
            NSString* symbolStr = [NSString stringWithUTF8String:symbolName];
            item.image = [NSImage imageWithSystemSymbolName:symbolStr accessibilityDescription:item.label];
        }
    }
    return (void*)item; // +1, intentionally not autoreleased, see callers
}

void* toolbarAddButtonItem(void* delegatePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, bool bordered, bool prominent,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    WailsToolbarItem* item = (WailsToolbarItem*)toolbarBuildButtonItemStandalone(
        identifier, itemID, label, symbolName, bordered); // +1
    applyCommonItemStyle(item, bordered, prominent, hasTint, tintR, tintG, tintB, tintA, badgeCount);

    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary takes its own retain
    [item release]; // release our +1; the dictionary now owns the only reference
    return (void*)item; // safe: dictionary keeps it alive; caller must not release
}

void* toolbarAddGroupItem(void* delegatePtr, const char* identifier,
    const char* label, void** memberItems, int memberCount) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;

    NSMutableArray<NSToolbarItem*>* subitems = [NSMutableArray arrayWithCapacity:memberCount];
    for (int i = 0; i < memberCount; i++) {
        NSToolbarItem* memberItem = (NSToolbarItem*)memberItems[i]; // +1 from the standalone builder
        [subitems addObject:memberItem]; // array takes its own retain
        [memberItem release]; // release our +1
    }

    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    NSToolbarItemGroup* group = [[NSToolbarItemGroup alloc] initWithItemIdentifier:identifierStr];
    // subitems is a plain property (there's no subitems-taking initializer);
    // each subitem already carries its own target/action from
    // toolbarBuildButtonItemStandalone, so clicking a segment dispatches
    // that segment's own itemID directly, no group-level target/action.
    group.subitems = subitems;
    group.label = [NSString stringWithUTF8String:label];
    if (@available(macOS 10.15, *)) {
        group.selectionMode = NSToolbarItemGroupSelectionModeMomentary;
        group.controlRepresentation = NSToolbarItemGroupControlRepresentationAutomatic;
    }

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = group; // dictionary retains
    [group release]; // release our +1 from alloc
    return (void*)group;
}

void* toolbarAddSearchItem(void* delegatePtr, const char* identifier, unsigned int itemID, const char* label) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];

    WailsSearchToolbarItem* item = [[WailsSearchToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = itemID;
    item.label = [NSString stringWithUTF8String:label];
    item.searchField.target = item;
    item.searchField.action = @selector(handleSearch:);
    item.searchField.sendsWholeSearchString = YES;

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary retains
    [item release]; // release our +1 from alloc
    return (void*)item;
}

void toolbarAddFlexibleSpaceIdentifier(void* delegatePtr) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    [delegate.orderedIdentifiers addObject:NSToolbarFlexibleSpaceItemIdentifier];
}

void toolbarAddSidebarToggleIdentifier(void* delegatePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];

    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = 0; // not user-dispatched; standard AppKit action below
    item.label = @"Sidebar";
    if (@available(macOS 11.0, *)) {
        item.image = [NSImage imageWithSystemSymbolName:@"sidebar.left" accessibilityDescription:@"Toggle Sidebar"];
    }
    item.target = nil; // let the responder chain find toggleSidebar: on the NSSplitViewController
    item.action = @selector(toggleSidebar:);

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary retains
    [item release]; // release our +1 from alloc
}
