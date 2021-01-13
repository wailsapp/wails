//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenu_darwin.h"
#include "trayicons.h"

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
    result->icon = mustJSONString(processedJSON, "Icon");
    JsonNode* processedMenu = mustJSONObject(processedJSON, "ProcessedMenu");

    // Create the menu
    result->menu = NewMenu(processedMenu);

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

void ShowTrayMenu(TrayMenu* trayMenu) {

    // Create a status bar item if we don't have one
    if( trayMenu->statusbaritem == NULL ) {
        id statusBar = msg( c("NSStatusBar"), s("systemStatusBar") );
        trayMenu->statusbaritem = msg(statusBar, s("statusItemWithLength:"), NSVariableStatusItemLength);
        msg(trayMenu->statusbaritem, s("retain"));

    }

    id statusBarButton = msg(trayMenu->statusbaritem, s("button"));
    msg(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);

    // Update the icon if needed
    UpdateTrayMenuIcon(trayMenu);

    // Update the label if needed
    UpdateTrayMenuLabel(trayMenu);

	// Update the menu
    id menu = GetMenu(trayMenu->menu);
    msg(trayMenu->statusbaritem, s("setMenu:"), menu);
}

void UpdateTrayMenuLabel(TrayMenu *trayMenu) {

    // Exit early if NULL
    if( trayMenu->label == NULL ) {
        return;
    }
    // We don't check for a
    id statusBarButton = msg(trayMenu->statusbaritem, s("button"));
    msg(statusBarButton, s("setTitle:"), str(trayMenu->label));
}

void UpdateTrayMenuIcon(TrayMenu *trayMenu) {

    // Exit early if NULL or empty string
    if( trayMenu->icon == NULL || STREMPTY(trayMenu->icon ) ) {
        return;
    }
    id trayImage = hashmap_get(&trayIconCache, trayMenu->icon, strlen(trayMenu->icon));
    id statusBarButton = msg(trayMenu->statusbaritem, s("button"));
    msg(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);
    msg(statusBarButton, s("setImage:"), trayImage);
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
        id statusBar = msg( c("NSStatusBar"), s("systemStatusBar") );
        msg(statusBar, s("removeStatusItem:"), trayMenu->statusbaritem);
        msg(trayMenu->statusbaritem, s("release"));
        trayMenu->statusbaritem = NULL;
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
        id imageData = msg(c("NSData"), s("dataWithBytes:length:"), data, length);
        id trayImage = ALLOC("NSImage");
        msg(trayImage, s("initWithData:"), imageData);
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