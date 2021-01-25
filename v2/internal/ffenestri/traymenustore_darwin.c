//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenustore_darwin.h"
#include "traymenu_darwin.h"
#include <stdlib.h>

TrayMenuStore* NewTrayMenuStore() {

    TrayMenuStore* result = malloc(sizeof(TrayMenuStore));

    // Allocate Tray Menu Store
    if( 0 != hashmap_create((const unsigned)4, &result->trayMenuMap)) {
        ABORT("[NewTrayMenuStore] Not enough memory to allocate trayMenuMap!");
    }

    return result;
}

int dumpTrayMenu(void *const context, struct hashmap_element_s *const e) {
    DumpTrayMenu(e->data);
    return 0;
}

void DumpTrayMenuStore(TrayMenuStore* store) {
    hashmap_iterate_pairs(&store->trayMenuMap, dumpTrayMenu, NULL);
}

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON) {
    TrayMenu* newMenu = NewTrayMenu(menuJSON);

    //TODO: check if there is already an entry for this menu
    hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);

}

int showTrayMenu(void *const context, struct hashmap_element_s *const e) {
    ShowTrayMenu(e->data);
    // 0 to retain element, -1 to delete.
    return 0;
}

void ShowTrayMenusInStore(TrayMenuStore* store) {
    if( hashmap_num_entries(&store->trayMenuMap) > 0 ) {
        hashmap_iterate_pairs(&store->trayMenuMap, showTrayMenu, NULL);
    }
}


int freeTrayMenu(void *const context, struct hashmap_element_s *const e) {
    DeleteTrayMenu(e->data);
    return -1;
}

void DeleteTrayMenuStore(TrayMenuStore *store) {

    // Delete context menus
    if (hashmap_num_entries(&store->trayMenuMap) > 0) {
        if (0 != hashmap_iterate_pairs(&store->trayMenuMap, freeTrayMenu, NULL)) {
            ABORT("[DeleteContextMenuStore] Failed to release contextMenuStore entries!");
        }
    }

    // Destroy tray menu map
    hashmap_destroy(&store->trayMenuMap);
}

TrayMenu* GetTrayMenuFromStore(TrayMenuStore* store, const char* menuID) {
    // Get the current menu
    return hashmap_get(&store->trayMenuMap, menuID, strlen(menuID));
}

TrayMenu* MustGetTrayMenuFromStore(TrayMenuStore* store, const char* menuID) {
    // Get the current menu
    TrayMenu* result = hashmap_get(&store->trayMenuMap, menuID, strlen(menuID));
    if (result == NULL ) {
        ABORT("Unable to find TrayMenu with ID '%s' in the TrayMenuStore!", menuID);
    }
    return result;
}

void UpdateTrayMenuLabelInStore(TrayMenuStore* store, const char* JSON) {
    // Parse the JSON
    JsonNode *parsedUpdate = mustParseJSON(JSON);

    // Get the data out
    const char* ID = mustJSONString(parsedUpdate, "ID");
    const char* Label = mustJSONString(parsedUpdate, "Label");

    // Check we have this menu
    TrayMenu *menu = MustGetTrayMenuFromStore(store, ID);
    UpdateTrayLabel(menu, Label);

}

void UpdateTrayMenuInStore(TrayMenuStore* store, const char* menuJSON) {
    TrayMenu* newMenu = NewTrayMenu(menuJSON);
//    DumpTrayMenu(newMenu);

    // Get the current menu
    TrayMenu *currentMenu = GetTrayMenuFromStore(store, newMenu->ID);

    // If we don't have a menu, we create one
    if ( currentMenu == NULL ) {
        // Store the new menu
        hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);

        // Show it
        ShowTrayMenu(newMenu);
        return;
    }
//    DumpTrayMenu(currentMenu);

    // Save the status bar reference
    newMenu->statusbaritem = currentMenu->statusbaritem;

    hashmap_remove(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID));

    // Delete the current menu
    DeleteMenu(currentMenu->menu);
    currentMenu->menu = NULL;

    // Free the tray menu memory
    MEMFREE(currentMenu);

    hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);

    // Show the updated menu
    ShowTrayMenu(newMenu);

}
