//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenustore_darwin.h"
#include "traymenu_darwin.h"

TrayMenuStore* NewTrayMenuStore() {

    TrayMenuStore* result = malloc(sizeof(TrayMenuStore));

    // Allocate Context Menu Store
    if( 0 != hashmap_create((const unsigned)4, &result->trayMenuMap)) {
        ABORT("[NewTrayMenuStore] Not enough memory to allocate trayMenuMap!");
    }

    return result;
}

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON) {
    TrayMenu* newMenu = NewTrayMenu(menuJSON);

    const char *ID = "TEST";

    hashmap_put(&store->trayMenuMap, ID, strlen(ID), newMenu);

}

int freeTrayMenu(void *const context, struct hashmap_element_s *const e) {
    DeleteTrayMenu(e->data);
    return -1;
}

void DeleteTrayMenuStore(TrayMenuStore *trayMenuStore) {

    // Delete context menus
    if( hashmap_num_entries(&trayMenuStore->trayMenuMap) > 0 ) {
        if (0 != hashmap_iterate_pairs(&trayMenuStore->trayMenuMap, freeTrayMenu, NULL)) {
            ABORT("[DeleteContextMenuStore] Failed to release contextMenuStore entries!");
        }
    }

    // Destroy tray menu map
    hashmap_destroy(&trayMenuStore->trayMenuMap);
}