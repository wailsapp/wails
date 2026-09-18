//go:build darwin && !ios && !server

#ifndef menu_mac_extras_darwin_h
#define menu_mac_extras_darwin_h

#import <Cocoa/Cocoa.h>
#include <stdbool.h>

// Go callbacks (menu_mac_extras_darwin.go)
extern void processMenuPaletteSelection(unsigned int menuItemID, int index);

// SF Symbol image on a menu item (macOS 11+). Frees symbolName. An empty
// name clears the image.
void setMenuItemSymbol(void* nsMenuItem, char* symbolName);

// NSMenuItemBadge (macOS 14+). Frees text. present == false clears the
// badge; a non-empty text wins over count.
void setMenuItemBadge(void* nsMenuItem, char* text, int count, bool present);

// NSControlStateValueMixed on/off.
void setMenuItemMixed(void* nsMenuItem, bool mixed);

// isAlternate flag.
void setMenuItemAlternate(void* nsMenuItem, bool alternate);

// indentationLevel (0..15).
void setMenuItemIndentationLevel(void* nsMenuItem, int level);

// Section header item (macOS 14+, disabled plain item below). Frees title.
// Returns a retained NSMenuItem.
void* newMenuItemSectionHeader(unsigned int menuItemID, char* title);

// Palette submenu item (macOS 14+). symbols is a newline-joined list of SF
// Symbol names (may be empty), rgba holds 4 doubles (0..1) per colour.
// Frees label and symbols. Returns a retained NSMenuItem whose submenu is
// the palette, or NULL below macOS 14.
void* newMenuItemPalette(unsigned int menuItemID, char* label, char* symbols,
    double* rgba, int colourCount, int selected);

// Recent documents (NSDocumentController). Main thread only.
void noteRecentDocument(char* path);
void clearRecentDocumentsList(void);
int recentDocumentCount(void);
char* recentDocumentPathAt(int index);

// Installs the "Open Recent" delegate on an NSMenu created by createNSMenu
// and fills it. The delegate repopulates it each time it opens.
void installOpenRecentMenu(void* nsMenu);

#endif /* menu_mac_extras_darwin_h */
