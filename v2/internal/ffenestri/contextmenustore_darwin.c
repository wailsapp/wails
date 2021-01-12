
#include "contextmenus_darwin.h"
#include "contextmenustore_darwin.h"

ContextMenuStore* NewContextMenuStore(const char* contextMenusAsJSON) {

    ContextMenuStore* result = malloc(sizeof(ContextMenuStore));

    // Init members
    result->contextMenusAsJSON = contextMenusAsJSON;
    result->processedContextMenus = NULL;
    result->contextMenuData = NULL;

    // Allocate Context Menu Store
    if( 0 != hashmap_create((const unsigned)4, &result->contextMenuStore)) {
        ABORT("[NewContextMenus] Not enough memory to allocate contextMenuStore!");
    }

    return result;
}


void DeleteContextMenuStore(ContextMenuStore* store) {

    // Guard against NULLs
    if( store == NULL ) {
        return;
    }

    // Delete context menus
    if( hashmap_num_entries(&store->contextMenuStore) > 0 ) {
        if (0 != hashmap_iterate_pairs(&store->contextMenuStore, freeContextMenu, NULL)) {
            ABORT("[DeleteContextMenuStore] Failed to release contextMenuStore entries!");
        }
    }

    // Free context menu hashmap
    hashmap_destroy(&store->contextMenuStore);

    // Destroy processed Context Menus
    if( store->processedContextMenus != NULL) {
        json_delete(store->processedContextMenus);
        store->processedContextMenus = NULL;
    }

    // Delete any context menu data we may have stored
    if( store->contextMenuData != NULL ) {
        MEMFREE(store->contextMenuData);
    }
}
