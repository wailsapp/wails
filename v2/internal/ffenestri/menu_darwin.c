//
// Created by Lea Anthony on 6/1/21.
//

#include "ffenestri_darwin.h"
#include "menu_darwin.h"
#include "contextmenus_darwin.h"
#include "common.h"

// NewMenu creates a new Menu struct, saving the given menu structure as JSON
Menu* NewMenu(JsonNode *menuData) {

    Menu *result = malloc(sizeof(Menu));

    result->processedMenu = menuData;

    // No title by default
    result->title = "";

    // Initialise menuCallbackDataCache
    vec_init(&result->callbackDataCache);

    // Allocate MenuItem Map
    if( 0 != hashmap_create((const unsigned)16, &result->menuItemMap)) {
        ABORT("[NewMenu] Not enough memory to allocate menuItemMap!");
    }
    // Allocate the Radio Group Map
    if( 0 != hashmap_create((const unsigned)4, &result->radioGroupMap)) {
        ABORT("[NewMenu] Not enough memory to allocate radioGroupMap!");
    }

    // Init other members
    result->menu = NULL;
    result->parentData = NULL;

    return result;
}

Menu* NewApplicationMenu(const char *menuAsJSON) {

    // Parse the menu json
    JsonNode *processedMenu = json_decode(menuAsJSON);
    if( processedMenu == NULL ) {
        // Parse error!
        ABORT("Unable to parse Menu JSON: %s", menuAsJSON);
    }

    Menu *result = NewMenu(processedMenu);
    result->menuType = ApplicationMenuType;
    return result;
}

MenuItemCallbackData* CreateMenuItemCallbackData(Menu *menu, id menuItem, const char *menuID, enum MenuItemType menuItemType) {
    MenuItemCallbackData* result = malloc(sizeof(MenuItemCallbackData));

    result->menu = menu;
    result->menuID = menuID;
    result->menuItem = menuItem;
    result->menuItemType = menuItemType;

    // Store reference to this so we can destroy later
    vec_push(&menu->callbackDataCache, result);

    return result;
}

void DeleteMenu(Menu *menu) {

    // Free menu item hashmap
    hashmap_destroy(&menu->menuItemMap);

    // Free radio group members
    if( hashmap_num_entries(&menu->radioGroupMap) > 0 ) {
        if (0 != hashmap_iterate_pairs(&menu->radioGroupMap, freeHashmapItem, NULL)) {
            ABORT("[DeleteMenu] Failed to release radioGroupMap entries!");
        }
    }

    // Free radio groups hashmap
    hashmap_destroy(&menu->radioGroupMap);

    // Free up the processed menu memory
    if (menu->processedMenu != NULL) {
        json_delete(menu->processedMenu);
        menu->processedMenu = NULL;
    }

    // Release the vector memory
    vec_deinit(&menu->callbackDataCache);

    // Free nsmenu if we have it
    if ( menu->menu != NULL ) {
        msg_reg(menu->menu, s("release"));
    }

    free(menu);
}

// Creates a JSON message for the given menuItemID and data
const char* createMenuClickedMessage(const char *menuItemID, const char *data, enum MenuType menuType, const char *parentID) {

    JsonNode *jsonObject = json_mkobject();
    if (menuItemID == NULL ) {
        ABORT("Item ID NULL for menu!!\n");
    }
    json_append_member(jsonObject, "menuItemID", json_mkstring(menuItemID));
    json_append_member(jsonObject, "menuType", json_mkstring(MenuTypeAsString[(int)menuType]));
    if (data != NULL) {
        json_append_member(jsonObject, "data", json_mkstring(data));
    }
    if (parentID != NULL) {
        json_append_member(jsonObject, "parentID", json_mkstring(parentID));
    }
    const char *payload = json_encode(jsonObject);
    json_delete(jsonObject);
    const char *result = concat("MC", payload);
    MEMFREE(payload);
    return result;
}

// Callback for text menu items
void menuItemCallback(id self, SEL cmd, id sender) {
    MenuItemCallbackData *callbackData = (MenuItemCallbackData *)msg_reg(msg_reg(sender, s("representedObject")), s("pointerValue"));
    const char *message;

    // Update checkbox / radio item
    if( callbackData->menuItemType == Checkbox) {
        // Toggle state
        bool state = msg_reg(callbackData->menuItem, s("state"));
        msg_int(callbackData->menuItem, s("setState:"), (state? NSControlStateValueOff : NSControlStateValueOn));
    } else if( callbackData->menuItemType == Radio ) {
        // Check the menu items' current state
        bool selected = (bool)msg_reg(callbackData->menuItem, s("state"));

        // If it's already selected, exit early
        if (selected) return;

        // Get this item's radio group members and turn them off
        id *members = (id*)hashmap_get(&(callbackData->menu->radioGroupMap), (char*)callbackData->menuID, strlen(callbackData->menuID));

        // Uncheck all members of the group
        id thisMember = members[0];
        int count = 0;
        while(thisMember != NULL) {
            msg_int(thisMember, s("setState:"), NSControlStateValueOff);
            count = count + 1;
            thisMember = members[count];
        }

        // check the selected menu item
        msg_int(callbackData->menuItem, s("setState:"), NSControlStateValueOn);
    }

    const char *menuID = callbackData->menuID;
    const char *data = NULL;
    enum MenuType menuType = callbackData->menu->menuType;
    const char *parentID = NULL;

    // Generate message to send to backend
    if( menuType == ContextMenuType ) {
        // Get the context menu data from the menu
        ContextMenu* contextMenu = (ContextMenu*) callbackData->menu->parentData;
        data = contextMenu->contextMenuData;
        parentID = contextMenu->ID;
    } else if ( menuType == TrayMenuType ) {
        parentID = (const char*) callbackData->menu->parentData;
    }

    message = createMenuClickedMessage(menuID, data, menuType, parentID);

    // Notify the backend
    messageFromWindowCallback(message);
    MEMFREE(message);
}

id processAcceleratorKey(const char *key) {

    // Guard against no accelerator key
    if( key == NULL ) {
        return str("");
    }

    if( STREQ(key, "backspace") ) {
        return strunicode(0x0008);
    }
    if( STREQ(key, "tab") ) {
        return strunicode(0x0009);
    }
    if( STREQ(key, "return") ) {
        return strunicode(0x000d);
    }
    if( STREQ(key, "enter") ) {
        return strunicode(0x000d);
    }
    if( STREQ(key, "escape") ) {
        return strunicode(0x001b);
    }
    if( STREQ(key, "left") ) {
        return strunicode(0x001c);
    }
    if( STREQ(key, "right") ) {
        return strunicode(0x001d);
    }
    if( STREQ(key, "up") ) {
        return strunicode(0x001e);
    }
    if( STREQ(key, "down") ) {
        return strunicode(0x001f);
    }
    if( STREQ(key, "space") ) {
        return strunicode(0x0020);
    }
    if( STREQ(key, "delete") ) {
        return strunicode(0x007f);
    }
    if( STREQ(key, "home") ) {
        return strunicode(0x2196);
    }
    if( STREQ(key, "end") ) {
        return strunicode(0x2198);
    }
    if( STREQ(key, "page up") ) {
        return strunicode(0x21de);
    }
    if( STREQ(key, "page down") ) {
        return strunicode(0x21df);
    }
    if( STREQ(key, "f1") ) {
        return strunicode(0xf704);
    }
    if( STREQ(key, "f2") ) {
        return strunicode(0xf705);
    }
    if( STREQ(key, "f3") ) {
        return strunicode(0xf706);
    }
    if( STREQ(key, "f4") ) {
        return strunicode(0xf707);
    }
    if( STREQ(key, "f5") ) {
        return strunicode(0xf708);
    }
    if( STREQ(key, "f6") ) {
        return strunicode(0xf709);
    }
    if( STREQ(key, "f7") ) {
        return strunicode(0xf70a);
    }
    if( STREQ(key, "f8") ) {
        return strunicode(0xf70b);
    }
    if( STREQ(key, "f9") ) {
        return strunicode(0xf70c);
    }
    if( STREQ(key, "f10") ) {
        return strunicode(0xf70d);
    }
    if( STREQ(key, "f11") ) {
        return strunicode(0xf70e);
    }
    if( STREQ(key, "f12") ) {
        return strunicode(0xf70f);
    }
    if( STREQ(key, "f13") ) {
        return strunicode(0xf710);
    }
    if( STREQ(key, "f14") ) {
        return strunicode(0xf711);
    }
    if( STREQ(key, "f15") ) {
        return strunicode(0xf712);
    }
    if( STREQ(key, "f16") ) {
        return strunicode(0xf713);
    }
    if( STREQ(key, "f17") ) {
        return strunicode(0xf714);
    }
    if( STREQ(key, "f18") ) {
        return strunicode(0xf715);
    }
    if( STREQ(key, "f19") ) {
        return strunicode(0xf716);
    }
    if( STREQ(key, "f20") ) {
        return strunicode(0xf717);
    }
    if( STREQ(key, "f21") ) {
        return strunicode(0xf718);
    }
    if( STREQ(key, "f22") ) {
        return strunicode(0xf719);
    }
    if( STREQ(key, "f23") ) {
        return strunicode(0xf71a);
    }
    if( STREQ(key, "f24") ) {
        return strunicode(0xf71b);
    }
    if( STREQ(key, "f25") ) {
        return strunicode(0xf71c);
    }
    if( STREQ(key, "f26") ) {
        return strunicode(0xf71d);
    }
    if( STREQ(key, "f27") ) {
        return strunicode(0xf71e);
    }
    if( STREQ(key, "f28") ) {
        return strunicode(0xf71f);
    }
    if( STREQ(key, "f29") ) {
        return strunicode(0xf720);
    }
    if( STREQ(key, "f30") ) {
        return strunicode(0xf721);
    }
    if( STREQ(key, "f31") ) {
        return strunicode(0xf722);
    }
    if( STREQ(key, "f32") ) {
        return strunicode(0xf723);
    }
    if( STREQ(key, "f33") ) {
        return strunicode(0xf724);
    }
    if( STREQ(key, "f34") ) {
        return strunicode(0xf725);
    }
    if( STREQ(key, "f35") ) {
        return strunicode(0xf726);
    }
//  if( STREQ(key, "Insert") ) {
//	return strunicode(0xf727);
//  }
//  if( STREQ(key, "PrintScreen") ) {
//	return strunicode(0xf72e);
//  }
//  if( STREQ(key, "ScrollLock") ) {
//	return strunicode(0xf72f);
//  }
    if( STREQ(key, "numLock") ) {
        return strunicode(0xf739);
    }

    return str(key);
}


void addSeparator(id menu) {
    id item = msg_reg(c("NSMenuItem"), s("separatorItem"));
    msg_id(menu, s("addItem:"), item);
}

id createMenuItemNoAutorelease( id title, const char *action, const char *key) {
    id item = ALLOC("NSMenuItem");
    ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), title, s(action), str(key));
    return item;
}

id createMenuItem(id title, const char *action, const char *key) {
    id item = ALLOC("NSMenuItem");
    ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), title, s(action), str(key));
    msg_reg(item, s("autorelease"));
    return item;
}

id addMenuItem(id menu, const char *title, const char *action, const char *key, bool disabled) {
    id item = createMenuItem(str(title), action, key);
    msg_bool(item, s("setEnabled:"), !disabled);
    msg_id(menu, s("addItem:"), item);
    return item;
}

id createMenu(id title) {
    id menu = ALLOC("NSMenu");
    msg_id(menu, s("initWithTitle:"), title);
    msg_bool(menu, s("setAutoenablesItems:"), NO);
//  msg(menu, s("autorelease"));
    return menu;
}

void createDefaultAppMenu(id parentMenu) {
// App Menu
    id appName = msg_reg(msg_reg(c("NSProcessInfo"), s("processInfo")), s("processName"));
    id appMenuItem = createMenuItemNoAutorelease(appName, NULL, "");
    id appMenu = createMenu(appName);

    msg_id(appMenuItem, s("setSubmenu:"), appMenu);
    msg_id(parentMenu, s("addItem:"), appMenuItem);

    id title = msg_id(str("Hide "), s("stringByAppendingString:"), appName);
    id item = createMenuItem(title, "hide:", "h");
    msg_id(appMenu, s("addItem:"), item);

    id hideOthers = addMenuItem(appMenu, "Hide Others", "hideOtherApplications:", "h", FALSE);
    msg_int(hideOthers, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagCommand));

    addMenuItem(appMenu, "Show All", "unhideAllApplications:", "", FALSE);

    addSeparator(appMenu);

    title = msg_id(str("Quit "), s("stringByAppendingString:"), appName);
    item = createMenuItem(title, "terminate:", "q");
    msg_id(appMenu, s("addItem:"), item);
}

void createDefaultEditMenu(id parentMenu) {
    // Edit Menu
    id editMenuItem = createMenuItemNoAutorelease(str("Edit"), NULL, "");
    id editMenu = createMenu(str("Edit"));

    msg_id(editMenuItem, s("setSubmenu:"), editMenu);
    msg_id(parentMenu, s("addItem:"), editMenuItem);

    addMenuItem(editMenu, "Undo", "undo:", "z", FALSE);
    addMenuItem(editMenu, "Redo", "redo:", "y", FALSE);
    addSeparator(editMenu);
    addMenuItem(editMenu, "Cut", "cut:", "x", FALSE);
    addMenuItem(editMenu, "Copy", "copy:", "c", FALSE);
    addMenuItem(editMenu, "Paste", "paste:", "v", FALSE);
    addMenuItem(editMenu, "Select All", "selectAll:", "a", FALSE);
}

void processMenuRole(Menu *menu, id parentMenu, JsonNode *item) {
    const char *roleName = item->string_;

    if ( STREQ(roleName, "appMenu") ) {
        createDefaultAppMenu(parentMenu);
        return;
    }
    if ( STREQ(roleName, "editMenu")) {
        createDefaultEditMenu(parentMenu);
        return;
    }
    if ( STREQ(roleName, "hide")) {
        addMenuItem(parentMenu, "Hide Window", "hide:", "h", FALSE);
        return;
    }
    if ( STREQ(roleName, "hideothers")) {
        id hideOthers = addMenuItem(parentMenu, "Hide Others", "hideOtherApplications:", "h", FALSE);
        msg_int(hideOthers, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagCommand));
        return;
    }
    if ( STREQ(roleName, "unhide")) {
        addMenuItem(parentMenu, "Show All", "unhideAllApplications:", "", FALSE);
        return;
    }
    if ( STREQ(roleName, "front")) {
        addMenuItem(parentMenu, "Bring All to Front", "arrangeInFront:", "", FALSE);
        return;
    }
    if ( STREQ(roleName, "undo")) {
        addMenuItem(parentMenu, "Undo", "undo:", "z", FALSE);
        return;
    }
    if ( STREQ(roleName, "redo")) {
        addMenuItem(parentMenu, "Redo", "redo:", "y", FALSE);
        return;
    }
    if ( STREQ(roleName, "cut")) {
        addMenuItem(parentMenu, "Cut", "cut:", "x", FALSE);
        return;
    }
    if ( STREQ(roleName, "copy")) {
        addMenuItem(parentMenu, "Copy", "copy:", "c", FALSE);
        return;
    }
    if ( STREQ(roleName, "paste")) {
        addMenuItem(parentMenu, "Paste", "paste:", "v", FALSE);
        return;
    }
    if ( STREQ(roleName, "delete")) {
        addMenuItem(parentMenu, "Delete", "delete:", "", FALSE);
        return;
    }
    if( STREQ(roleName, "pasteandmatchstyle")) {
        id pasteandmatchstyle = addMenuItem(parentMenu, "Paste and Match Style", "pasteandmatchstyle:", "v", FALSE);
        msg_int(pasteandmatchstyle, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagShift | NSEventModifierFlagCommand));
    }
    if ( STREQ(roleName, "selectall")) {
        addMenuItem(parentMenu, "Select All", "selectAll:", "a", FALSE);
        return;
    }
    if ( STREQ(roleName, "minimize")) {
        addMenuItem(parentMenu, "Minimize", "miniaturize:", "m", FALSE);
        return;
    }
    if ( STREQ(roleName, "zoom")) {
        addMenuItem(parentMenu, "Zoom", "performZoom:", "", FALSE);
        return;
    }
    if ( STREQ(roleName, "quit")) {
        addMenuItem(parentMenu, "Quit (More work TBD)", "terminate:", "q", FALSE);
        return;
    }
    if ( STREQ(roleName, "togglefullscreen")) {
        addMenuItem(parentMenu, "Toggle Full Screen", "toggleFullScreen:", "f", FALSE);
        return;
    }

}

// This converts a string array of modifiers into the
// equivalent MacOS Modifier Flags
unsigned long parseModifiers(const char **modifiers) {

    // Our result is a modifier flag list
    unsigned long result = 0;

    const char *thisModifier = modifiers[0];
    int count = 0;
    while( thisModifier != NULL ) {

        // Determine flags
        if( STREQ(thisModifier, "cmdorctrl") ) {
            result |= NSEventModifierFlagCommand;
        }
        if( STREQ(thisModifier, "optionoralt") ) {
            result |= NSEventModifierFlagOption;
        }
        if( STREQ(thisModifier, "shift") ) {
            result |= NSEventModifierFlagShift;
        }
        if( STREQ(thisModifier, "super") ) {
            result |= NSEventModifierFlagCommand;
        }
        if( STREQ(thisModifier, "ctrl") ) {
            result |= NSEventModifierFlagControl;
        }
        count++;
        thisModifier = modifiers[count];
    }
    return result;
}

id processRadioMenuItem(Menu *menu, id parentmenu, const char *title, const char *menuid, bool disabled, bool checked, const char *acceleratorkey) {
    id item = ALLOC("NSMenuItem");

    // Store the item in the menu item map
    hashmap_put(&menu->menuItemMap, (char*)menuid, strlen(menuid), item);

    // Create a MenuItemCallbackData
    MenuItemCallbackData *callback = CreateMenuItemCallbackData(menu, item, menuid, Radio);

    id wrappedId = msg_id(c("NSValue"), s("valueWithPointer:"), (id)callback);
    msg_id(item, s("setRepresentedObject:"), wrappedId);

    id key = processAcceleratorKey(acceleratorkey);

    ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), str(title), s("menuItemCallback:"), key);

    msg_bool(item, s("setEnabled:"), !disabled);
    msg_reg(item, s("autorelease"));
    msg_int(item, s("setState:"), (checked ? NSControlStateValueOn : NSControlStateValueOff));

    msg_id(parentmenu, s("addItem:"), item);
    return item;

}

id processCheckboxMenuItem(Menu *menu, id parentmenu, const char *title, const char *menuid, bool disabled, bool checked, const char *key) {

    id item = ALLOC("NSMenuItem");

    // Store the item in the menu item map
    hashmap_put(&menu->menuItemMap, (char*)menuid, strlen(menuid), item);

    // Create a MenuItemCallbackData
    MenuItemCallbackData *callback = CreateMenuItemCallbackData(menu, item, menuid, Checkbox);

    id wrappedId = msg_id(c("NSValue"), s("valueWithPointer:"), (id)callback);
    msg_id(item, s("setRepresentedObject:"), wrappedId);
    ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), str(title), s("menuItemCallback:"), str(key));
    msg_bool(item, s("setEnabled:"), !disabled);
    msg_reg(item, s("autorelease"));
    msg_int(item, s("setState:"), (checked ? NSControlStateValueOn : NSControlStateValueOff));
    msg_id(parentmenu, s("addItem:"), item);
    return item;
}

// getColour returns the colour from a styledLabel based on the key
const char* getColour(JsonNode *styledLabelEntry, const char* key) {
    JsonNode* colEntry = getJSONObject(styledLabelEntry, key);
    if( colEntry == NULL ) {
        return NULL;
    }
    return getJSONString(colEntry, "hex");
}

id createAttributedStringFromStyledLabel(JsonNode *styledLabel, const char* fontName, int fontSize) {

    // Create result
    id attributedString = ALLOC_INIT("NSMutableAttributedString");
    msg_reg(attributedString, s("autorelease"));

    // Create new Dictionary
    id dictionary = ALLOC_INIT("NSMutableDictionary");
    msg_reg(dictionary, s("autorelease"));

    // Use default font
    CGFloat fontSizeFloat = (CGFloat)fontSize;
    id font = ((id(*)(id, SEL, CGFloat))objc_msgSend)(c("NSFont"), s("menuBarFontOfSize:"), fontSizeFloat);

    // Check user supplied font
    if( STR_HAS_CHARS(fontName) ) {
        id fontNameAsNSString = str(fontName);
        id userFont = ((id(*)(id, SEL, id, CGFloat))objc_msgSend)(c("NSFont"), s("fontWithName:size:"), fontNameAsNSString, fontSizeFloat);
        if( userFont != NULL ) {
            font = userFont;
        }
    }

    id fan = lookupStringConstant(str("NSFontAttributeName"));
    id NSForegroundColorAttributeName = lookupStringConstant(str("NSForegroundColorAttributeName"));
    id NSBackgroundColorAttributeName = lookupStringConstant(str("NSBackgroundColorAttributeName"));

    // Loop over styled text creating NSAttributedText and appending to result
    JsonNode *styledLabelEntry;
    json_foreach(styledLabelEntry, styledLabel) {

        // Clear dictionary
        msg_reg(dictionary, s("removeAllObjects"));

        // Add font to dictionary
        msg_id_id(dictionary, s("setObject:forKey:"), font, fan);

        // Get Text
        const char* thisLabel = mustJSONString(styledLabelEntry, "Label");

        // Get foreground colour
        const char *hexColour = getColour(styledLabelEntry, "FgCol");
        if( hexColour != NULL) {
            unsigned short r, g, b, a;

            // white by default
            r = g = b = a = 255;
            int count = sscanf(hexColour, "#%02hx%02hx%02hx%02hx", &r, &g, &b, &a);
            if (count > 0) {
                id colour = ((id(*)(id, SEL, CGFloat, CGFloat, CGFloat, CGFloat))objc_msgSend)(c("NSColor"), s("colorWithCalibratedRed:green:blue:alpha:"),
                                    (CGFloat)r / (CGFloat)255.0,
                                    (CGFloat)g / (CGFloat)255.0,
                                    (CGFloat)b / (CGFloat)255.0,
                                    (CGFloat)a / (CGFloat)255.0);
                msg_id_id(dictionary, s("setObject:forKey:"), colour, NSForegroundColorAttributeName);
            }
        }

        // Get background colour
        hexColour = getColour(styledLabelEntry, "BgCol");
        if( hexColour != NULL) {
            unsigned short r, g, b, a;

            // white by default
            r = g = b = a = 255;
            int count = sscanf(hexColour, "#%02hx%02hx%02hx%02hx", &r, &g, &b, &a);
            if (count > 0) {
                id colour = ((id(*)(id, SEL, CGFloat, CGFloat, CGFloat, CGFloat))objc_msgSend)(c("NSColor"), s("colorWithCalibratedRed:green:blue:alpha:"),
                                    (CGFloat)r / (CGFloat)255.0,
                                    (CGFloat)g / (CGFloat)255.0,
                                    (CGFloat)b / (CGFloat)255.0,
                                    (CGFloat)a / (CGFloat)255.0);
                msg_id_id(dictionary, s("setObject:forKey:"), colour, NSForegroundColorAttributeName);
            }
        }

        // Create AttributedText
        id thisString = ALLOC("NSMutableAttributedString");
        msg_reg(thisString, s("autorelease"));
        msg_id_id(thisString, s("initWithString:attributes:"), str(thisLabel), dictionary);

        // Append text to result
        msg_id(attributedString, s("appendAttributedString:"), thisString);
    }

    return attributedString;

}


id createAttributedString(const char* title, const char* fontName, int fontSize, const char* RGBA) {

    // Create new Dictionary
    id dictionary = ALLOC_INIT("NSMutableDictionary");
    CGFloat fontSizeFloat = (CGFloat)fontSize;

    // Use default font
    id font = ((id(*)(id, SEL, CGFloat))objc_msgSend)(c("NSFont"), s("menuBarFontOfSize:"), fontSizeFloat);

    // Check user supplied font
    if( STR_HAS_CHARS(fontName) ) {
        id fontNameAsNSString = str(fontName);
        id userFont = ((id(*)(id, SEL, id, CGFloat))objc_msgSend)(c("NSFont"), s("fontWithName:size:"), fontNameAsNSString, fontSizeFloat);
        if( userFont != NULL ) {
            font = userFont;
        }
    }

    // Add font to dictionary
    id fan = lookupStringConstant(str("NSFontAttributeName"));
    msg_id_id(dictionary, s("setObject:forKey:"), font, fan);

    // RGBA
    if( RGBA != NULL && strlen(RGBA) > 0) {
        unsigned short r, g, b, a;

        // white by default
        r = g = b = a = 255;
        int count = sscanf(RGBA, "#%02hx%02hx%02hx%02hx", &r, &g, &b, &a);
        if (count > 0) {
			id colour = ((id(*)(id, SEL, CGFloat, CGFloat, CGFloat, CGFloat))objc_msgSend)(c("NSColor"), s("colorWithCalibratedRed:green:blue:alpha:"),
								(CGFloat)r / (CGFloat)255.0,
								(CGFloat)g / (CGFloat)255.0,
								(CGFloat)b / (CGFloat)255.0,
								(CGFloat)a / (CGFloat)255.0);
			id NSForegroundColorAttributeName = lookupStringConstant(str("NSForegroundColorAttributeName"));
            msg_id_id(dictionary, s("setObject:forKey:"), colour, NSForegroundColorAttributeName);
        }
    }

    id attributedString = ALLOC("NSMutableAttributedString");
    msg_id_id(attributedString, s("initWithString:attributes:"), str(title), dictionary);
    msg_reg(attributedString, s("autorelease"));
    msg_reg(dictionary, s("autorelease"));
    return attributedString;
}

id processTextMenuItem(Menu *menu, id parentMenu, const char *title, const char *menuid, bool disabled, const char *acceleratorkey, const char **modifiers, const char* tooltip, const char* image, const char* fontName, int fontSize, const char* RGBA, bool templateImage, bool alternate, JsonNode* styledLabel) {
    id item = ALLOC("NSMenuItem");

    // Create a MenuItemCallbackData
    MenuItemCallbackData *callback = CreateMenuItemCallbackData(menu, item, menuid, Text);

    id wrappedId = msg_id(c("NSValue"), s("valueWithPointer:"), (id)callback);
    msg_id(item, s("setRepresentedObject:"), wrappedId);

    if( !alternate ) {
        id key = processAcceleratorKey(acceleratorkey);
        ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), str(title),
            s("menuItemCallback:"), key);
    } else {
        ((id(*)(id, SEL, id, SEL, id))objc_msgSend)(item, s("initWithTitle:action:keyEquivalent:"), str(title), s("menuItemCallback:"), str(""));
    }

    if( tooltip != NULL ) {
        msg_id(item, s("setToolTip:"), str(tooltip));
    }

    // Process image
    if( image != NULL && strlen(image) > 0) {
        id nsimage = createImageFromBase64Data(image, templateImage);
        msg_id(item, s("setImage:"), nsimage);
    }

    id attributedString = NULL;
    if( styledLabel != NULL) {
        attributedString = createAttributedStringFromStyledLabel(styledLabel, fontName, fontSize);
    } else {
        attributedString = createAttributedString(title, fontName, fontSize, RGBA);
    }
    msg_id(item, s("setAttributedTitle:"), attributedString);

//msg_id(item, s("setTitle:"), str(title));

    msg_bool(item, s("setEnabled:"), !disabled);
    msg_reg(item, s("autorelease"));

    // Process modifiers
    if( modifiers != NULL && !alternate) {
        unsigned long modifierFlags = parseModifiers(modifiers);
        ((id(*)(id, SEL, unsigned long))objc_msgSend)(item, s("setKeyEquivalentModifierMask:"), modifierFlags);
    }

    // alternate
    if( alternate ) {
        msg_bool(item, s("setAlternate:"), true);
        msg_int(item, s("setKeyEquivalentModifierMask:"), NSEventModifierFlagOption);
    }
    msg_id(parentMenu, s("addItem:"), item);

    return item;
}

void processMenuItem(Menu *menu, id parentMenu, JsonNode *item) {

    // Check if this item is hidden and if so, exit early!
    bool hidden = false;
    getJSONBool(item, "Hidden", &hidden);
    if( hidden ) {
        return;
    }

    // Get the role
    JsonNode *role = json_find_member(item, "Role");
    if( role != NULL ) {
        processMenuRole(menu, parentMenu, role);
        return;
    }

    // This is a user menu. Get the common data
    // Get the label
    const char *label = getJSONString(item, "Label");
    if ( label == NULL) {
        label = "(empty)";
    }

    // Check for a styled label
    JsonNode *styledLabel = getJSONObject(item, "StyledLabel");

    // Is this an alternate menu item?
    bool alternate = false;
    getJSONBool(item, "MacAlternate", &alternate);

    const char *menuid = getJSONString(item, "ID");
    if ( menuid == NULL) {
        menuid = "";
    }

    bool disabled = false;
    getJSONBool(item, "Disabled", &disabled);

    // Get the Accelerator
    JsonNode *accelerator = json_find_member(item, "Accelerator");
    const char *acceleratorkey = NULL;
    const char **modifiers = NULL;

    const char *tooltip = getJSONString(item, "Tooltip");
    const char *image = getJSONString(item, "Image");
    const char *fontName = getJSONString(item, "FontName");
    const char *RGBA = getJSONString(item, "RGBA");
    bool templateImage = false;
    getJSONBool(item, "MacTemplateImage", &templateImage);

    int fontSize = 0;
    getJSONInt(item, "FontSize", &fontSize);

    // If we have an accelerator
    if( accelerator != NULL ) {
        // Get the key
        acceleratorkey = getJSONString(accelerator, "Key");
        // Check if there are modifiers
        JsonNode *modifiersList = json_find_member(accelerator, "Modifiers");
        if ( modifiersList != NULL ) {
            // Allocate an array of strings
            int noOfModifiers = json_array_length(modifiersList);

            // Do we have any?
            if (noOfModifiers > 0) {
                modifiers = malloc(sizeof(const char *) * (noOfModifiers + 1));
                JsonNode *modifier;
                int count = 0;
                // Iterate the modifiers and save a reference to them in our new array
                json_foreach(modifier, modifiersList) {
                    // Get modifier name
                    modifiers[count] = modifier->string_;
                    count++;
                }
                // Null terminate the modifier list
                modifiers[count] = NULL;
            }
        }
    }

    // Get the Type
    JsonNode *type = json_find_member(item, "Type");
    if( type != NULL ) {
        if( STREQ(type->string_, "Text") || STREQ(type->string_, "Submenu")) {
            id thisMenuItem = processTextMenuItem(menu, parentMenu, label, menuid, disabled, acceleratorkey, modifiers, tooltip, image, fontName, fontSize, RGBA, templateImage, alternate, styledLabel);

            // Check if this node has a submenu
            JsonNode *submenu = json_find_member(item, "SubMenu");
            if( submenu != NULL ) {
                // Get the label
                JsonNode *menuNameNode = json_find_member(item, "Label");
                const char *name = "";
                if ( menuNameNode != NULL) {
                    name = menuNameNode->string_;
                }

                id thisMenu = createMenu(str(name));

                msg_id(thisMenuItem, s("setSubmenu:"), thisMenu);

                JsonNode *submenuItems = json_find_member(submenu, "Items");
                // If we have no items, just return
                if ( submenuItems == NULL ) {
                    return;
                }

                // Loop over submenu items
                JsonNode *item;
                json_foreach(item, submenuItems) {
                    // Get item label
                    processMenuItem(menu, thisMenu, item);
                }
            }
        }
        else if ( STREQ(type->string_, "Separator")) {
            addSeparator(parentMenu);
        }
        else if ( STREQ(type->string_, "Checkbox")) {
            // Get checked state
            bool checked = false;
            getJSONBool(item, "Checked", &checked);

            processCheckboxMenuItem(menu, parentMenu, label, menuid, disabled, checked, "");
        }
        else if ( STREQ(type->string_, "Radio")) {
            // Get checked state
            bool checked = false;
            getJSONBool(item, "Checked", &checked);

            processRadioMenuItem(menu, parentMenu, label, menuid, disabled, checked, "");
        }
    }

    if ( modifiers != NULL ) {
        free(modifiers);
    }

    return;
}

void processMenuData(Menu *menu, JsonNode *menuData) {
    JsonNode *items = json_find_member(menuData, "Items");
    if( items == NULL ) {
        // Parse error!
        ABORT("Unable to find 'Items' in menu JSON!");
    }

    // Iterate items
    JsonNode *item;
    json_foreach(item, items) {
        // Process each menu item
        processMenuItem(menu, menu->menu, item);
    }
}

void processRadioGroupJSON(Menu *menu, JsonNode *radioGroup) {

    int groupLength;
    getJSONInt(radioGroup, "Length", &groupLength);
    JsonNode *members = json_find_member(radioGroup, "Members");
    JsonNode *member;

    // Allocate array
    size_t arrayLength = sizeof(id)*(groupLength+1);
    id memberList[arrayLength];

    // Build the radio group items
    int count=0;
    json_foreach(member, members) {
        // Get menu by id
        id menuItem = (id)hashmap_get(&menu->menuItemMap, (char*)member->string_, strlen(member->string_));
        // Save Member
        memberList[count] = menuItem;
        count = count + 1;
    }
    // Null terminate array
    memberList[groupLength] = 0;

    // Store the members
    json_foreach(member, members) {
        // Copy the memberList
        char *newMemberList = (char *)malloc(arrayLength);
        memcpy(newMemberList, memberList, arrayLength);
        // add group to each member of group
        hashmap_put(&menu->radioGroupMap, member->string_, strlen(member->string_), newMemberList);
    }

}

id GetMenu(Menu *menu) {

    // Pull out the menu data
    JsonNode *menuData = json_find_member(menu->processedMenu, "Menu");
    if( menuData == NULL ) {
        ABORT("Unable to find Menu data: %s", menu->processedMenu);
    }

    menu->menu = createMenu(str(""));

    // Process the menu data
    processMenuData(menu, menuData);

    // Create the radiogroup cache
    JsonNode *radioGroups = json_find_member(menu->processedMenu, "RadioGroups");
    if( radioGroups == NULL ) {
        // Parse error!
        ABORT("Unable to find RadioGroups data: %s", menu->processedMenu);
    }

    // Iterate radio groups
    JsonNode *radioGroup;
    json_foreach(radioGroup, radioGroups) {
        // Get item label
        processRadioGroupJSON(menu, radioGroup);
    }

    return menu->menu;
}

