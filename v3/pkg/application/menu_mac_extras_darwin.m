//go:build darwin && !ios && !server

#import "menu_mac_extras_darwin.h"
#import "menuitem_darwin.h"
#import "application_darwin_delegate.h"

// Runs block on the main thread synchronously. Menu state must be applied
// synchronously so a reopen reflects it (see setMenuItemChecked in
// menuitem_darwin.go / wailsapp/wails#5002); dispatch_sync from the main
// thread would deadlock, so short-circuit there.
static void menuExtrasOnMain(dispatch_block_t block) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

// Mirrors toolbarSymbolImage in webview_window_toolbar_darwin.m.
static NSImage* menuSymbolImage(NSString* symbolName, NSString* accessibilityLabel) {
    if (symbolName == nil || symbolName.length == 0) return nil;
    if (@available(macOS 11.0, *)) {
        return [NSImage imageWithSystemSymbolName:symbolName
                         accessibilityDescription:accessibilityLabel];
    }
    return nil;
}

void setMenuItemSymbol(void* nsMenuItem, char* symbolName) {
    NSMenuItem* menuItem = (NSMenuItem*)nsMenuItem;
    NSString* name = symbolName != NULL ? [NSString stringWithUTF8String:symbolName] : @"";
    if (symbolName != NULL) free(symbolName);
    menuExtrasOnMain(^{
        if (name.length == 0) {
            menuItem.image = nil;
            return;
        }
        NSImage* image = menuSymbolImage(name, menuItem.title);
        if (image == nil) {
            NSLog(@"[wails] SetSymbol(%@) needs macOS 11 or later; ignored", name);
            return;
        }
        menuItem.image = image;
    });
}

void setMenuItemBadge(void* nsMenuItem, char* text, int count, bool present) {
    NSMenuItem* menuItem = (NSMenuItem*)nsMenuItem;
    NSString* badgeText = text != NULL ? [NSString stringWithUTF8String:text] : @"";
    if (text != NULL) free(text);
    menuExtrasOnMain(^{
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
        if (@available(macOS 14.0, *)) {
            if (!present) {
                menuItem.badge = nil;
                return;
            }
            NSMenuItemBadge* badge = badgeText.length > 0
                ? [[NSMenuItemBadge alloc] initWithString:badgeText]
                : [[NSMenuItemBadge alloc] initWithCount:count];
            menuItem.badge = badge;
            [badge release];
            return;
        }
#endif
        if (present) {
            NSLog(@"[wails] menu item badges need macOS 14 or later; ignored");
        }
    });
}

void setMenuItemMixed(void* nsMenuItem, bool mixed) {
    NSMenuItem* menuItem = (NSMenuItem*)nsMenuItem;
    menuExtrasOnMain(^{
        menuItem.state = mixed ? NSControlStateValueMixed : NSControlStateValueOff;
    });
}

void setMenuItemAlternate(void* nsMenuItem, bool alternate) {
    NSMenuItem* menuItem = (NSMenuItem*)nsMenuItem;
    menuExtrasOnMain(^{
        menuItem.alternate = alternate;
    });
}

void setMenuItemIndentationLevel(void* nsMenuItem, int level) {
    NSMenuItem* menuItem = (NSMenuItem*)nsMenuItem;
    menuExtrasOnMain(^{
        menuItem.indentationLevel = level;
    });
}

void* newMenuItemSectionHeader(unsigned int menuItemID, char* title) {
    NSString* nsTitle = title != NULL ? [NSString stringWithUTF8String:title] : @"";
    if (title != NULL) free(title);
    NSMenuItem* item = nil;
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
    if (@available(macOS 14.0, *)) {
        item = [[NSMenuItem sectionHeaderWithTitle:nsTitle] retain];
    }
#endif
    if (item == nil) {
        NSLog(@"[wails] section headers need macOS 14 or later; showing a disabled item");
        item = [[NSMenuItem alloc] initWithTitle:nsTitle action:nil keyEquivalent:@""];
        item.enabled = NO;
    }
    item.tag = menuItemID;
    return (void*)item;
}

void* newMenuItemPalette(unsigned int menuItemID, char* label, char* symbols,
    double* rgba, int colourCount, int selected) {
    NSString* nsLabel = label != NULL ? [NSString stringWithUTF8String:label] : @"";
    NSString* symbolList = symbols != NULL ? [NSString stringWithUTF8String:symbols] : @"";
    if (label != NULL) free(label);
    if (symbols != NULL) free(symbols);
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
    if (@available(macOS 14.0, *)) {
        NSArray<NSString*>* symbolNames = symbolList.length > 0
            ? [symbolList componentsSeparatedByString:@"\n"]
            : @[];
        NSMutableArray<NSColor*>* colours = [NSMutableArray arrayWithCapacity:colourCount];
        NSMutableArray<NSString*>* titles = [NSMutableArray arrayWithCapacity:colourCount];
        for (int i = 0; i < colourCount; i++) {
            [colours addObject:[NSColor colorWithSRGBRed:rgba[i * 4]
                                                   green:rgba[i * 4 + 1]
                                                    blue:rgba[i * 4 + 2]
                                                   alpha:rgba[i * 4 + 3]]];
            NSString* title = (NSInteger)symbolNames.count == colourCount
                ? symbolNames[i]
                : [NSString stringWithFormat:@"Swatch %d", i + 1];
            [titles addObject:title];
        }
        void (^onSelection)(NSMenu*) = ^(NSMenu* menu) {
            NSMenuItem* chosen = menu.selectedItems.firstObject;
            int index = chosen != nil ? (int)[menu indexOfItem:chosen] : -1;
            processMenuPaletteSelection(menuItemID, index);
        };
        NSMenu* paletteMenu = nil;
        if (symbolNames.count == 1) {
            NSImage* templateImage = menuSymbolImage(symbolNames[0], nsLabel);
            paletteMenu = [NSMenu paletteMenuWithColors:colours
                                                 titles:titles
                                          templateImage:templateImage
                                       selectionHandler:onSelection];
        } else {
            paletteMenu = [NSMenu paletteMenuWithColors:colours
                                                 titles:titles
                                       selectionHandler:onSelection];
            if ((NSInteger)symbolNames.count == colourCount) {
                for (int i = 0; i < colourCount && i < (int)paletteMenu.numberOfItems; i++) {
                    NSImage* image = menuSymbolImage(symbolNames[i], titles[i]);
                    if (image == nil) continue;
                    NSImageSymbolConfiguration* config =
                        [NSImageSymbolConfiguration configurationWithPaletteColors:@[colours[i]]];
                    [paletteMenu itemAtIndex:i].image = [image imageWithSymbolConfiguration:config];
                }
            }
        }
        paletteMenu.selectionMode = NSMenuSelectionModeSelectOne;
        if (selected >= 0 && selected < (int)paletteMenu.numberOfItems) {
            [paletteMenu itemAtIndex:selected].state = NSControlStateValueOn;
        }
        MenuItem* item = [MenuItem new];
        item.menuItemID = menuItemID;
        item.tag = menuItemID;
        item.title = nsLabel;
        item.submenu = paletteMenu;
        return (void*)item;
    }
#endif
    NSLog(@"[wails] palette menus need macOS 14 or later; item skipped");
    return NULL;
}

// Recent documents

void noteRecentDocument(char* path) {
    if (path == NULL) return;
    NSString* nsPath = [NSString stringWithUTF8String:path];
    free(path);
    NSURL* url = [NSURL fileURLWithPath:nsPath];
    if (url != nil) {
        [[NSDocumentController sharedDocumentController] noteNewRecentDocumentURL:url];
    }
}

void clearRecentDocumentsList(void) {
    [[NSDocumentController sharedDocumentController] clearRecentDocuments:nil];
}

int recentDocumentCount(void) {
    return (int)[[NSDocumentController sharedDocumentController] recentDocumentURLs].count;
}

char* recentDocumentPathAt(int index) {
    NSArray<NSURL*>* urls = [[NSDocumentController sharedDocumentController] recentDocumentURLs];
    if (index < 0 || index >= (int)urls.count) return NULL;
    NSString* path = urls[index].path;
    if (path == nil) return NULL;
    return strdup(path.UTF8String);
}

// WailsOpenRecentMenu populates an "Open Recent" submenu from
// NSDocumentController each time it is about to open. AppKit only fills the
// standard Open Recent menu for NSDocument-based apps, so this builds the
// same shape: one item per recent file with its Finder icon, a separator and
// "Clear Menu". Opening an item routes through HandleOpenFile, the same path
// Finder and file associations use.
@interface WailsOpenRecentMenu : NSObject <NSMenuDelegate>
+ (instancetype)shared;
- (void)populate:(NSMenu*)menu;
@end

@implementation WailsOpenRecentMenu

+ (instancetype)shared {
    static WailsOpenRecentMenu* instance = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        instance = [[WailsOpenRecentMenu alloc] init];
    });
    return instance;
}

- (void)menuNeedsUpdate:(NSMenu*)menu {
    [self populate:menu];
}

- (void)populate:(NSMenu*)menu {
    [menu removeAllItems];
    NSArray<NSURL*>* urls = [[NSDocumentController sharedDocumentController] recentDocumentURLs];
    for (NSURL* url in urls) {
        NSString* path = url.path;
        if (path == nil) continue;
        NSString* title = [[NSFileManager defaultManager] displayNameAtPath:path];
        NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:title
                                                      action:@selector(openRecent:)
                                               keyEquivalent:@""];
        item.target = self;
        item.representedObject = url;
        item.toolTip = path;
        NSImage* icon = [[NSWorkspace sharedWorkspace] iconForFile:path];
        if (icon != nil) {
            icon = [icon copy];
            icon.size = NSMakeSize(16, 16);
            item.image = icon;
            [icon release];
        }
        [menu addItem:item];
        [item release];
    }
    if (urls.count > 0) {
        [menu addItem:[NSMenuItem separatorItem]];
    }
    NSMenuItem* clear = [[NSMenuItem alloc] initWithTitle:@"Clear Menu"
                                                   action:@selector(clearRecentDocuments:)
                                            keyEquivalent:@""];
    clear.target = [NSDocumentController sharedDocumentController];
    clear.enabled = urls.count > 0;
    [menu addItem:clear];
    [clear release];
}

- (void)openRecent:(NSMenuItem*)sender {
    NSURL* url = sender.representedObject;
    if (![url isKindOfClass:[NSURL class]] || url.path == nil) return;
    [[NSDocumentController sharedDocumentController] noteNewRecentDocumentURL:url];
    HandleOpenFile((char*)url.path.UTF8String);
}

@end

void installOpenRecentMenu(void* nsMenu) {
    NSMenu* menu = (NSMenu*)nsMenu;
    menu.delegate = [WailsOpenRecentMenu shared];
    [[WailsOpenRecentMenu shared] populate:menu];
}
