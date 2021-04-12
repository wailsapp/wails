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

    if (pthread_mutex_init(&result->lock, NULL) != 0) {
        printf("\n mutex init has failed\n");
        exit(1);
    }

    return result;
}

int dumpTrayMenu(void *const context, struct hashmap_element_s *const e) {
    DumpTrayMenu(e->data);
    return 0;
}

void DumpTrayMenuStore(TrayMenuStore* store) {
    pthread_mutex_lock(&store->lock);
    hashmap_iterate_pairs(&store->trayMenuMap, dumpTrayMenu, NULL);
    pthread_mutex_unlock(&store->lock);
}

void AddTrayMenuToStore(TrayMenuStore* store, const char* menuJSON) {

    TrayMenu* newMenu = NewTrayMenu(menuJSON);

    pthread_mutex_lock(&store->lock);
    //TODO: check if there is already an entry for this menu
    hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);
    pthread_mutex_unlock(&store->lock);
}

int showTrayMenu(void *const context, struct hashmap_element_s *const e) {
    ShowTrayMenu(e->data);
    // 0 to retain element, -1 to delete.
    return 0;
}

void ShowTrayMenusInStore(TrayMenuStore* store) {
    pthread_mutex_lock(&store->lock);
    if( hashmap_num_entries(&store->trayMenuMap) > 0 ) {
        hashmap_iterate_pairs(&store->trayMenuMap, showTrayMenu, NULL);
    }
    pthread_mutex_unlock(&store->lock);
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

    pthread_mutex_destroy(&store->lock);
}

TrayMenu* GetTrayMenuFromStore(TrayMenuStore* store, const char* menuID) {
    // Get the current menu
    pthread_mutex_lock(&store->lock);
    TrayMenu* result = hashmap_get(&store->trayMenuMap, menuID, strlen(menuID));
    pthread_mutex_unlock(&store->lock);
    return result;
}

TrayMenu* MustGetTrayMenuFromStore(TrayMenuStore* store, const char* menuID) {
    // Get the current menu
    pthread_mutex_lock(&store->lock);
    TrayMenu* result = hashmap_get(&store->trayMenuMap, menuID, strlen(menuID));
    pthread_mutex_unlock(&store->lock);

    if (result == NULL ) {
        ABORT("Unable to find TrayMenu with ID '%s' in the TrayMenuStore!", menuID);
    }
    return result;
}

void DeleteTrayMenuInStore(TrayMenuStore* store, const char* ID) {

    TrayMenu *menu = MustGetTrayMenuFromStore(store, ID);
    pthread_mutex_lock(&store->lock);
    hashmap_remove(&store->trayMenuMap, ID, strlen(ID));
    pthread_mutex_unlock(&store->lock);
    DeleteTrayMenu(menu);
}

void UpdateTrayMenuLabelInStore(TrayMenuStore* store, const char* JSON) {
    // Parse the JSON
    JsonNode *parsedUpdate = mustParseJSON(JSON);

    // Get the data out
    const char* ID = mustJSONString(parsedUpdate, "ID");
    const char* Label = mustJSONString(parsedUpdate, "Label");

    // Check we have this menu
    TrayMenu *menu = MustGetTrayMenuFromStore(store, ID);

    const char *fontName = getJSONString(parsedUpdate, "FontName");
    const char *RGBA = getJSONString(parsedUpdate, "RGBA");
    int fontSize = 0;
    getJSONInt(parsedUpdate, "FontSize", &fontSize);
    const char *tooltip = getJSONString(parsedUpdate, "Tooltip");
    bool disabled = false;
    getJSONBool(parsedUpdate, "Disabled", &disabled);

    JsonNode *styledLabel = getJSONObject(parsedUpdate, "StyledLabel");

    UpdateTrayLabel(menu, Label, fontName, fontSize, RGBA, tooltip, disabled, styledLabel);


}

void UpdateTrayMenuInStore(TrayMenuStore* store, const char* menuJSON) {
    TrayMenu* newMenu = NewTrayMenu(menuJSON);
//    DumpTrayMenu(newMenu);

    // Get the current menu
    TrayMenu *currentMenu = GetTrayMenuFromStore(store, newMenu->ID);

    // If we don't have a menu, we create one
    if ( currentMenu == NULL ) {
        // Store the new menu
        pthread_mutex_lock(&store->lock);
        hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);
        pthread_mutex_unlock(&store->lock);

        // Show it
        ShowTrayMenu(newMenu);
        return;
    }
//    DumpTrayMenu(currentMenu);

    // Save the status bar reference
    newMenu->statusbaritem = currentMenu->statusbaritem;

    pthread_mutex_lock(&store->lock);
    hashmap_remove(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID));
    pthread_mutex_unlock(&store->lock);

    // Delete the current menu
    DeleteMenu(currentMenu->menu);
    currentMenu->menu = NULL;

    // Free the tray menu memory
    MEMFREE(currentMenu);

    pthread_mutex_lock(&store->lock);
    hashmap_put(&store->trayMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);
    pthread_mutex_unlock(&store->lock);

    // Show the updated menu
    ShowTrayMenu(newMenu);
}
