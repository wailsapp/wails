//
// Created by Lea Anthony on 6/1/21.
//

#ifndef MENU_DARWIN_H
#define MENU_DARWIN_H

#include "common.h"

typedef struct {

    /*** Internal ***/

    const char *menuAsJSON;

    struct hashmap_s menuItemMap;
    struct hashmap_s radioGroupMap;

} Menu;

// NewMenu creates a new Menu struct, saving the given menu structure as JSON
Menu* NewMenu(const char *menuAsJSON) {

    Menu *result = malloc(sizeof(Menu));

    // menuAsJSON is allocated and freed by Go
    result->menuAsJSON = menuAsJSON;

    // Allocate MenuItem Map
    if( 0 != hashmap_create((const unsigned)16, &result->menuItemMap)) {
        ABORT("[NewMenu] Not enough memory to allocate menuItemMap!");
    }
    // Allocate the Radio Group Map
    if( 0 != hashmap_create((const unsigned)4, &result->radioGroupMap)) {
        ABORT("[NewMenu] Not enough memory to allocate radioGroupMap!");
    }

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

    free(menu);
}


void Create() {

//
//    // Allocate the hashmaps we need
//    allocateMenuHashMaps(app);
//
//    // Create a new menu bar
//    id menubar = createMenu(str(""));
//
//    // Parse the processed menu json
//    app->processedMenu = json_decode(app->menuAsJSON);
//
//    if( app->processedMenu == NULL ) {
//        // Parse error!
//        Fatal(app, "Unable to parse Menu JSON: %s", app->menuAsJSON);
//        return;
//    }
//
//
//    // Pull out the Menu
//    JsonNode *menuData = json_find_member(app->processedMenu, "Menu");
//    if( menuData == NULL ) {
//        // Parse error!
//        Fatal(app, "Unable to find Menu data: %s", app->processedMenu);
//        return;
//    }
//
//
//    parseMenu(app, menubar, menuData, &menuItemMapForApplicationMenu,
//              "checkboxMenuCallbackForApplicationMenu:", "radioMenuCallbackForApplicationMenu:", "menuCallbackForApplicationMenu:");
//
//    // Create the radiogroup cache
//    JsonNode *radioGroups = json_find_member(app->processedMenu, "RadioGroups");
//    if( radioGroups == NULL ) {
//        // Parse error!
//        Fatal(app, "Unable to find RadioGroups data: %s", app->processedMenu);
//        return;
//    }
//
//    // Iterate radio groups
//    JsonNode *radioGroup;
//    json_foreach(radioGroup, radioGroups) {
//        // Get item label
//        processRadioGroup(radioGroup, &menuItemMapForApplicationMenu, &radioGroupMapForApplicationMenu);
//    }
//
//    // Apply the menu bar
//    msg(msg(c("NSApplication"), s("sharedApplication")), s("setMainMenu:"), menubar);

}

#endif //ASSETS_C_MENU_DARWIN_H
