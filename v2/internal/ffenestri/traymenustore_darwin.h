//
// Created by Lea Anthony on 7/1/21.
//

#ifndef TRAYMENUSTORE_DARWIN_H
#define TRAYMENUSTORE_DARWIN_H

#include "traymenu_darwin.h"

#include <pthread.h>

typedef struct {

	int dummy;

    // This is our tray menu map
    // It maps tray IDs to TrayMenu*
    struct hashmap_s trayMenuMap;

    pthread_mutex_t lock;

} TrayMenuStore;

TrayMenuStore* NewTrayMenuStore();

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON);
void UpdateTrayMenuInStore(TrayMenuStore* store, const char* menuJSON);
void ShowTrayMenusInStore(TrayMenuStore* store);
void DeleteTrayMenuStore(TrayMenuStore* store);

TrayMenu* GetTrayMenuByID(TrayMenuStore* store, const char* menuID);

void UpdateTrayMenuLabelInStore(TrayMenuStore* store, const char* JSON);
void DeleteTrayMenuInStore(TrayMenuStore* store, const char* id);

#endif //TRAYMENUSTORE_DARWIN_H
