//
// Created by Lea Anthony on 6/1/21.
//

#ifndef MENU_DARWIN_H
#define MENU_DARWIN_H

#include "common.h"
#include "ffenestri_darwin.h"

enum MenuItemType {Text = 0, Checkbox = 1, Radio = 2};
enum MenuType {ApplicationMenuType = 0, ContextMenuType = 1, TrayMenuType = 2};
static const char *MenuTypeAsString[] = {
        "ApplicationMenu", "ContextMenu", "TrayMenu",
};

typedef struct _NSRange {
    unsigned long location;
    unsigned long length;
} NSRange;

#define NSFontWeightUltraLight -0.8
#define NSFontWeightThin -0.6
#define NSFontWeightLight -0.4
#define NSFontWeightRegular 0.0
#define NSFontWeightMedium 0.23
#define NSFontWeightSemibold 0.3
#define NSFontWeightBold 0.4
#define NSFontWeightHeavy 0.56
#define NSFontWeightBlack 0.62

extern void messageFromWindowCallback(const char *);

typedef struct {

    const char *title;

    /*** Internal ***/

    // The decoded version of the Menu JSON
    JsonNode *processedMenu;

    struct hashmap_s menuItemMap;
    struct hashmap_s radioGroupMap;

    // Vector to keep track of callback data memory
    vec_void_t callbackDataCache;

    // The NSMenu for this menu
    id menu;

    // The parent data, eg ContextMenuStore or Tray
    void *parentData;

    // The commands for the menu callbacks
    const char *callbackCommand;

    // This indicates if we are an Application Menu, tray menu or context menu
    enum MenuType menuType;


} Menu;


typedef struct {
    id menuItem;
    Menu *menu;
    const char *menuID;
    enum MenuItemType menuItemType;
} MenuItemCallbackData;



// NewMenu creates a new Menu struct, saving the given menu structure as JSON
Menu* NewMenu(JsonNode *menuData);

Menu* NewApplicationMenu(const char *menuAsJSON);
MenuItemCallbackData* CreateMenuItemCallbackData(Menu *menu, id menuItem, const char *menuID, enum MenuItemType menuItemType);

void DeleteMenu(Menu *menu);

// Creates a JSON message for the given menuItemID and data
const char* createMenuClickedMessage(const char *menuItemID, const char *data, enum MenuType menuType, const char *parentID);
// Callback for text menu items
void menuItemCallback(id self, SEL cmd, id sender);
id processAcceleratorKey(const char *key);


void addSeparator(id menu);
id createMenuItemNoAutorelease( id title, const char *action, const char *key);

id createMenuItem(id title, const char *action, const char *key);

id addMenuItem(id menu, const char *title, const char *action, const char *key, bool disabled);

id createMenu(id title);
void createDefaultAppMenu(id parentMenu);
void createDefaultEditMenu(id parentMenu);

void processMenuRole(Menu *menu, id parentMenu, JsonNode *item);
// This converts a string array of modifiers into the
// equivalent MacOS Modifier Flags
unsigned long parseModifiers(const char **modifiers);
id processRadioMenuItem(Menu *menu, id parentmenu, const char *title, const char *menuid, bool disabled, bool checked, const char *acceleratorkey);

id processCheckboxMenuItem(Menu *menu, id parentmenu, const char *title, const char *menuid, bool disabled, bool checked, const char *key);

id processTextMenuItem(Menu *menu, id parentMenu, const char *title, const char *menuid, bool disabled, const char *acceleratorkey, const char **modifiers, const char* tooltip, const char* image, const char* fontName, int fontSize, const char* RGBA, bool templateImage, bool alternate, JsonNode* styledLabel);
void processMenuItem(Menu *menu, id parentMenu, JsonNode *item);
void processMenuData(Menu *menu, JsonNode *menuData);

void processRadioGroupJSON(Menu *menu, JsonNode *radioGroup);
id GetMenu(Menu *menu);
id createAttributedString(const char* title, const char* fontName, int fontSize, const char* RGBA);
id createAttributedStringFromStyledLabel(JsonNode *styledLabel, const char* fontName, int fontSize);

#endif //ASSETS_C_MENU_DARWIN_H
