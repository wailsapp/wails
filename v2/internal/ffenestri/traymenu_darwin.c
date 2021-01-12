//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenu_darwin.h"

extern struct hashmap_s trayIconCache;

TrayMenu* NewTrayMenu(const char* menuJSON) {
    TrayMenu* result = malloc(sizeof(TrayMenu));

/*
 {"ID":"0","Label":"Test Tray Label","Icon":"","ProcessedMenu":{"Menu":{"Items":[{"ID":"0","Label":"Show Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"1","Label":"Hide Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"2","Label":"Minimise Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0},{"ID":"3","Label":"Unminimise Window","Type":"Text","Disabled":false,"Hidden":false,"Checked":false,"Foreground":0,"Background":0}]},"RadioGroups":null}}
*/
    JsonNode* processedJSON = json_decode(menuJSON);
    if( processedJSON == NULL ) {
        ABORT("[NewTrayMenu] Unable to parse TrayMenu JSON: %s", menuJSON);
    }

    // TODO: Make this configurable
    result->trayIconPosition = NSImageLeft;

    result->ID = mustJSONString(processedJSON, "ID");
    result->label = mustJSONString(processedJSON, "Label");
    result->icon = mustJSONString(processedJSON, "Icon");
    JsonNode* processedMenu = mustJSONObject(processedJSON, "ProcessedMenu");

    // Create the menu
    result->menu = NewMenu(processedMenu);

    // Set the menu type and store the tray ID in the parent data
    result->menu->menuType = TrayMenuType;
    result->menu->parentData = (void*) result->ID;

    return result;
}

void DumpTrayMenu(TrayMenu* trayMenu) {
    printf("    ['%s':%p] = { label: '%s', icon: '%s', menu: %p }\n", trayMenu->ID, trayMenu, trayMenu->label, trayMenu->icon, trayMenu->menu );
}

void ShowTrayMenu(TrayMenu* trayMenu) {

    // Create a status bar item if we don't have one
    if( trayMenu->statusbaritem == NULL ) {
        id statusBar = msg( c("NSStatusBar"), s("systemStatusBar") );
        trayMenu->statusbaritem = msg(statusBar, s("statusItemWithLength:"), NSVariableStatusItemLength);
        msg(trayMenu->statusbaritem, s("retain"));
        id statusBarButton = msg(trayMenu->statusbaritem, s("button"));
        msg(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);

        // Update the icon if needed
        UpdateTrayMenuIcon(trayMenu);

        // Update the label if needed
        UpdateTrayMenuLabel(trayMenu);
    }

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

    // Exit early if NULL or emptystring
    if( trayMenu->icon == NULL || STREMPTY(trayMenu->icon ) ) {
        return;
    }
    id trayImage = hashmap_get(&trayIconCache, trayMenu->icon, strlen(trayMenu->icon));
    id statusBarButton = msg(trayMenu->statusbaritem, s("button"));
    msg(statusBarButton, s("setImagePosition:"), trayMenu->trayIconPosition);
    msg(statusBarButton, s("setImage:"), trayImage);
}

void DeleteTrayMenu(TrayMenu* trayMenu) {

//    printf("Freeing TrayMenu:\n");
//    DumpTrayMenu(trayMenu);

    // Delete the menu
    DeleteMenu(trayMenu->menu);

    // Free JSON
    json_delete(trayMenu->processedJSON);

    // Free the tray menu memory
    MEMFREE(trayMenu);
}