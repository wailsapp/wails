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

void UpdateTrayMenuInStore(TrayMenuStore* store, const char* menuJSON) {
    TrayMenu* newMenu = NewTrayMenu(menuJSON);

    // Get the current menu
    TrayMenu *currentMenu = GetTrayMenuFromStore(store, newMenu->ID);
    if ( currentMenu == NULL ) {
        ABORT("Attempted to update unknown tray menu with ID '%s'.", newMenu->ID);
    }

    // Save the status bar reference
    newMenu->statusbaritem = currentMenu->statusbaritem;

    hashmap_remove(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID));

    // Delete the current menu
    DeleteMenu(currentMenu->menu);
    currentMenu->menu = NULL;

    // Free JSON
    if (currentMenu->processedJSON != NULL ) {
        json_delete(currentMenu->processedJSON);
        currentMenu->processedJSON = NULL;
    }

    // Free the tray menu memory
    MEMFREE(currentMenu);

    hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);

    // Show the updated menu
    ShowTrayMenu(newMenu);

}
