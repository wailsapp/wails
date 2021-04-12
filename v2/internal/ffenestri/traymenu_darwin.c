//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenu_darwin.h"
#include "trayicons.h"

extern Class trayMenuDelegateClass;

// A cache for all our tray menu icons
// Global because it's a singleton
struct hashmap_s trayIconCache;

TrayMenu* NewTrayMenu(const char* menuJSON) {
    TrayMenu* result = malloc(sizeof(TrayMenu));

/*
 {"ID":"0","Label":"Test Tray Label","Icon":"","ProcessedMenu":{"Menu":{"Items":[{"ID":"0","Label":"Show Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"1","Label":"Hide Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"2","Label":"Minimise Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"3","Label":"Unminimise Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0}]},"RadioGroups":null}}
*/
    JsonNode* processedJSON = json_decode(menuJSON);
    if( processedJSON == NULL ) {
        ABORT("[NewTrayMenu] Unable to parse TrayMenu JSON: %s", menuJSON);
    }

    // Save reference to this json
    result->processedJSON = processedJSON;

    // TODO: Make this configurable
    result->trayIconPosition = NSImageLeft;

    result->ID = mustJSONString(processedJSON, "ID");
    result->label = mustJSONString(processedJSON, "Label");
    result->icon = mustJSONString(processedJSON, "Image");
    result->fontName = getJSONString(processedJSON, "FontName");
    result->RGBA = getJSONString(processedJSON, "RGBA");
    getJSONBool(processedJSON, "MacTemplateImage", &result->templateImage);
    result->fontSize = 0;
    getJSONInt(processedJSON, "FontSize", &result->fontSize);
    result->tooltip = NULL;
    result->tooltip = getJSONString(processedJSON, "Tooltip");
    result->disabled = false;
    getJSONBool(processedJSON, "Disabled", &result->disabled);

    result->styledLabel = getJSONObject(processedJSON, "StyledLabel");

    // Create the menu
    JsonNode* processedMenu = mustJSONObject(processedJSON, "ProcessedMenu");
    result->menu = NewMenu(processedMenu);

    result->delegate = NULL;

    // Init tray status bar item
    result->statusbaritem = NULL;

    // Set the menu type and store the tray ID in the parent data
    result->menu->menuType = TrayMenuType;
    result->menu->parentData = (void*) result->ID;

    return result;
}

void DumpTrayMenu(TrayMenu* trayMenu) {
    printf("    ['%s':%p] = { label: '%s', icon: '%s', menu: %p, statusbar: %p  }\n", trayMenu->ID, trayMenu, trayMenu->label, trayMenu->icon, trayMenu->menu, trayMenu->statusbaritem );
}


void UpdateTrayLabel(TrayMenu *trayMenu, const char *label, const char *fontName, int fontSize, const char *RGBA, const char *tooltip, bool disabled, JsonNode *styledLabel) {

    // Exit early if NULL
    if( trayMenu->label == NULL ) {
        return;
    }
    // Update button label
    id statusBarButton = msg_reg(trayMenu->statusbaritem, s("button"));
    id attributedString = NULL;
    if( styledLabel != NULL) {
        attributedString = createAttributedStringFromStyledLabel(styledLabel, fontName, fontSize);
    } else {
        attributedString = createAttributedString(label, fontName, fontSize, RGBA);
    }

    if( tooltip != NULL ) {
        msg_id(statusBarButton, s("setToolTip:"), str(tooltip));
    }

    msg_bool(statusBarButton, s("setEnabled:"), !disabled);

    msg_id(statusBarButton, s("setAttributedTitle:"), attributedString);
}

void UpdateTrayIcon(TrayMenu *trayMenu) {

    // Exit early if NULL
    if( trayMenu->icon == NULL ) {
        return;
    }

    id statusBarButton = msg_reg(trayMenu->statusbaritem, s("button"));

    // Empty icon means remove it
    if( STREMPTY(trayMenu->icon) ) {
        // Remove image
        msg_id(statusBarButton, s("setImage:"), NULL);
        return;
    }

    id trayImage = hashmap_get(&trayIconCache, trayMenu->icon, strlen(trayMenu->icon));

    // If we don't have the image in the icon cache then assume it's base64 encoded image data
    if (trayImage == NULL) {
        trayImage = createImageFromBase64Data(trayMenu->icon, trayMenu->templateImage);
    }

    msg_int(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);
    msg_id(statusBarButton, s("setImage:"), trayImage);

}

void ShowTrayMenu(TrayMenu* trayMenu) {

    // Create a status bar item if we don't have one
    if( trayMenu->statusbaritem == NULL ) {
        id statusBar = msg_reg( c("NSStatusBar"), s("systemStatusBar") );
        trayMenu->statusbaritem = ((id(*)(id, SEL, CGFloat))objc_msgSend)(statusBar, s("statusItemWithLength:"), NSVariableStatusItemLength);
        msg_reg(trayMenu->statusbaritem, s("retain"));
    }

    id statusBarButton = msg_reg(trayMenu->statusbaritem, s("button"));
    msg_uint(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);
    // Update the icon if needed
    UpdateTrayIcon(trayMenu);

    // Update the label if needed
    UpdateTrayLabel(trayMenu, trayMenu->label, trayMenu->fontName, trayMenu->fontSize, trayMenu->RGBA, trayMenu->tooltip, trayMenu->disabled, trayMenu->styledLabel);

    // Update the menu
    id menu = GetMenu(trayMenu->menu);
    objc_setAssociatedObject(menu, "trayMenuID", str(trayMenu->ID), OBJC_ASSOCIATION_ASSIGN);

	// Create delegate
	id trayMenuDelegate = msg_reg((id)trayMenuDelegateClass, s("new"));
	msg_id(menu, s("setDelegate:"), trayMenuDelegate);
    objc_setAssociatedObject(trayMenuDelegate, "menu", menu, OBJC_ASSOCIATION_ASSIGN);

    // Create menu delegate
    trayMenu->delegate = trayMenuDelegate;

    msg_id(trayMenu->statusbaritem, s("setMenu:"), menu);
}

// UpdateTrayMenuInPlace receives 2 menus. The current menu gets
// updated with the data from the new menu.
void UpdateTrayMenuInPlace(TrayMenu* currentMenu, TrayMenu* newMenu) {

    // Delete the old menu
    DeleteMenu(currentMenu->menu);

    // Set the new one
    currentMenu->menu = newMenu->menu;

    // Delete the old JSON
    json_delete(currentMenu->processedJSON);

    // Set the new JSON
    currentMenu->processedJSON = newMenu->processedJSON;

    // Copy the other data
    currentMenu->ID = newMenu->ID;
    currentMenu->label = newMenu->label;
    currentMenu->styledLabel = newMenu->styledLabel;
    currentMenu->trayIconPosition = newMenu->trayIconPosition;
    currentMenu->icon = newMenu->icon;

}

void DeleteTrayMenu(TrayMenu* trayMenu) {

//    printf("Freeing TrayMenu:\n");
//    DumpTrayMenu(trayMenu);

    // Delete the menu
    DeleteMenu(trayMenu->menu);

    // Free JSON
    if (trayMenu->processedJSON != NULL ) {
        json_delete(trayMenu->processedJSON);
    }

    // Free the status item
    if ( trayMenu->statusbaritem != NULL ) {
        id statusBar = msg_reg( c("NSStatusBar"), s("systemStatusBar") );
        msg_id(statusBar, s("removeStatusItem:"), trayMenu->statusbaritem);
        msg_reg(trayMenu->statusbaritem, s("release"));
        trayMenu->statusbaritem = NULL;
    }

    if ( trayMenu->delegate != NULL ) {
        msg_reg(trayMenu->delegate, s("release"));
    }

    // Free the tray menu memory
    MEMFREE(trayMenu);
}

void LoadTrayIcons() {

    // Allocate the Tray Icons
    if( 0 != hashmap_create((const unsigned)4, &trayIconCache)) {
        // Couldn't allocate map
        ABORT("Not enough memory to allocate trayIconCache!");
    }

    unsigned int count = 0;
    while( 1 ) {
        const unsigned char *name = trayIcons[count++];
        if( name == 0x00 ) {
            break;
        }
        const unsigned char *lengthAsString = trayIcons[count++];
        if( name == 0x00 ) {
            break;
        }
        const unsigned char *data = trayIcons[count++];
        if( data == 0x00 ) {
            break;
        }
        int length = atoi((const char *)lengthAsString);

        // Create the icon and add to the hashmap
        id imageData = ((id(*)(id, SEL, id, int))objc_msgSend)(c("NSData"), s("dataWithBytes:length:"), (id)data, length);
        id trayImage = ALLOC("NSImage");
        msg_id(trayImage, s("initWithData:"), imageData);
        hashmap_put(&trayIconCache, (const char *)name, strlen((const char *)name), trayImage);
    }
}

void UnloadTrayIcons() {
    // Release the tray cache images
    if( hashmap_num_entries(&trayIconCache) > 0 ) {
        if (0!=hashmap_iterate_pairs(&trayIconCache, releaseNSObject, NULL)) {
            ABORT("failed to release hashmap entries!");
        }
    }

    //Free radio groups hashmap
    hashmap_destroy(&trayIconCache);
}