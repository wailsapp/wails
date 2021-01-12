//
// Created by Lea Anthony on 7/1/21.
//

#ifndef TRAYMENUSTORE_DARWIN_H
#define TRAYMENUSTORE_DARWIN_H


typedef struct {

	int dummy;

    // This is our tray menu map
    // It maps tray IDs to TrayMenu*
    struct hashmap_s trayMenuMap;

} TrayMenuStore;

TrayMenuStore* NewTrayMenuStore();

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON);
void ShowTrayMenusInStore(TrayMenuStore* store);
void DeleteTrayMenuStore(TrayMenuStore* store);

#endif //TRAYMENUSTORE_DARWIN_H
