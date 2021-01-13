
#include "contextmenus_darwin.h"
#include "contextmenustore_darwin.h"

ContextMenuStore* NewContextMenuStore() {

    ContextMenuStore* result = malloc(sizeof(ContextMenuStore));

    // Allocate Context Menu Store
    if( 0 != hashmap_create((const unsigned)4, &result->contextMenuMap)) {
        ABORT("[NewContextMenus] Not enough memory to allocate contextMenuStore!");
    }

    return result;
}

void AddContextMenuToStore(ContextMenuStore* store, const char* contextMenuJSON) {
    ContextMenu* newMenu = NewContextMenu(contextMenuJSON);

    //TODO: check if there is already an entry for this menu
    hashmap_put(&store->contextMenuMap, newMenu->ID, strlen(newMenu->ID), newMenu);
}

ContextMenu* GetContextMenuFromStore(ContextMenuStore* store, const char* menuID) {
    // Get the current menu
    return hashmap_get(&store->contextMenuMap, menuID, strlen(menuID));
}

void UpdateContextMenuInStore(ContextMenuStore* store, const char* menuJSON) {
    ContextMenu* newContextMenu = NewContextMenu(menuJSON);

    // Get the current menu
    ContextMenu *currentMenu = GetContextMenuFromStore(store, newContextMenu->ID);
    if ( currentMenu == NULL ) {
        ABORT("Attempted to update unknown context menu with ID '%s'.", newContextMenu->ID);
    }

    hashmap_remove(&store->contextMenuMap, newContextMenu->ID, strlen(newContextMenu->ID));

    // Save the status bar reference
    DeleteContextMenu(currentMenu);

    hashmap_put(&store->contextMenuMap, newContextMenu->ID, strlen(newContextMenu->ID), newContextMenu);

}


void DeleteContextMenuStore(ContextMenuStore* store) {

    // Guard against NULLs
    if( store == NULL ) {
        return;
    }

    // Delete context menus
    if( hashmap_num_entries(&store->contextMenuMap) > 0 ) {
        if (0 != hashmap_iterate_pairs(&store->contextMenuMap, freeContextMenu, NULL)) {
            ABORT("[DeleteContextMenuStore] Failed to release contextMenuStore entries!");
        }
    }

    // Free context menu hashmap
    hashmap_destroy(&store->contextMenuMap);

}
