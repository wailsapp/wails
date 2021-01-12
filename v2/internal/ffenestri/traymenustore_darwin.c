//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenustore_darwin.h"

TrayMenuStore* NewTrayMenuStore() {

    TrayMenuStore* result = malloc(sizeof(TrayMenuStore));

    // Allocate Context Menu Store
    if( 0 != hashmap_create((const unsigned)4, &result->trayMenuMap)) {
        ABORT("[NewTrayMenuStore] Not enough memory to allocate trayMenuMap!");
    }

    return result;
}

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON) {

}

void DeleteTrayMenuStore(TrayMenuStore *trayMenuStore) {

    // Destroy tray menu map
    hashmap_destroy(&trayMenuStore->trayMenuMap);
}