////
//// Created by Lea Anthony on 6/1/21.
////
//
#ifndef CONTEXTMENU_DARWIN_H
#define CONTEXTMENU_DARWIN_H

#include "common.h"
#include "menu_darwin.h"
#include "contextmenustore_darwin.h"

typedef struct {
    const char* ID;
    id nsmenu;
    Menu* menu;
} ContextMenu;


ContextMenu* NewContextMenu(JsonNode* menuData, ContextMenuStore *store) {
    ContextMenu* result = malloc(sizeof(ContextMenu));
    result->menu = NewMenu(menuData);
    result->nsmenu = NULL;
    result->menu->menuType = ContextMenuType;
    result->menu->parentData = store;
    return result;
}


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

ContextMenu* GetContextMenuByID(ContextMenuStore* store, const char *contextMenuID) {
    return (ContextMenu*)hashmap_get(&store->contextMenuStore, (char*)contextMenuID, strlen(contextMenuID));
}

void DeleteContextMenu(ContextMenu* contextMenu) {
    // Free Menu
    DeleteMenu(contextMenu->menu);

    // Free context menu
    free(contextMenu);
}

int freeContextMenu(void *const context, struct hashmap_element_s *const e) {
    DeleteContextMenu(e->data);
    return -1;
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

void ProcessContextMenus(ContextMenuStore* store) {

    // Decode the context menus JSON
    store->processedContextMenus = json_decode(store->contextMenusAsJSON);
	if( store->processedContextMenus == NULL ) {
		ABORT("[ProcessContextMenus] Unable to parse Context Menus JSON: %s", store->contextMenusAsJSON);
	}

//	// Get the context menu items
//	JsonNode *contextMenuItems = json_find_member(store->processedContextMenus, "Items");
//	if( contextMenuItems == NULL ) {
//		ABORT("[ProcessContextMenus] Unable to find Items in processedContextMenus!");
//	}

	// Iterate context menus
	JsonNode *contextMenu;
	json_foreach(contextMenu, store->processedContextMenus) {

	    const char* ID = getJSONString(contextMenu, "ID");
	    if ( ID == NULL ) {
	        ABORT("Unable to read ID of contextMenu\n");
	    }

	    JsonNode* processedMenu = json_find_member(contextMenu, "ProcessedMenu");
        if ( processedMenu == NULL ) {
            ABORT("Unable to read ProcessedMenu of contextMenu\n");
        }
        // Create a new context menu instance
        ContextMenu *thisContextMenu = NewContextMenu(processedMenu, store);

		// Store the item in the context menu map
		hashmap_put(&store->contextMenuStore, (char*)ID, strlen(ID), thisContextMenu);
	}

}

//
//
//bool ContextMenuExists(ContextMenus *contextMenus, const char* contextMenuID) {
//    return hashmap_get(&contextMenus->contextMenuStore, contextMenuID, strlen(contextMenuID)) != NULL;
//}
//
//bool AddContextMenu(ContextMenu* contextMenu) {
//
//    // Check if we already have this
//    if( ContextMenuExists(contextMenu->ID) ) {
//        return false;
//    }
//
//    // Store the context menu
//    if (0 != hashmap_put(&contextMenus->contextMenuStore, contextMenu->ID, strlen(contextMenu->ID), contextMenu)) {
//        ABORT("Unable to add context menu with ID '%s'", contextMenu->ID);
//    }
//
//    return true;
//}
//
//ContextMenus* NewContextMenus(const char* contextMenusAsJSON) {
//
//    ContextMenus* result = malloc(sizeof(ContextMenus));
//
//    // Allocate Context Menu Store
//    if( 0 != hashmap_create((const unsigned)4, &result->contextMenuStore)) {
//        ABORT("[NewContextMenus] Not enough memory to allocate contextMenuStore!");
//    }
//
//    //
//
//    return result;
//}
//
//void ProcessContextMenus() {
//    // Parse the context menu json
//    processedContextMenus = json_decode(contextMenusAsJSON);
//    if( processedContextMenus == NULL ) {
//        // Parse error!
//        ABORT("Unable to parse Context Menus JSON: %s", contextMenusAsJSON);
//    }
//
//    JsonNode *contextMenuItems = json_find_member(processedContextMenus, "Items");
//    if( contextMenuItems == NULL ) {
//        // Parse error!
//        ABORT("Unable to find Items in Context menus");
//    }
//    // Iterate context menus
//    JsonNode *contextMenu;
//    json_foreach(contextMenu, contextMenuItems) {
//        Menu *contextMenu = NewMenu()
//
//        // Store the item in the context menu map
//        hashmap_put(&contextMenuMap, (char*)contextMenu->key, strlen(contextMenu->key), menu);
//    }
//
//}
//
//ContextMenu* NewContextMenu() {
//    ContextMenu* result = malloc(sizeof(ContextMenu));
//
//    result->menu = NewMenu(contextMenuAsJSON);
//
//    return result;
//}
//

//
//void InitContextMenuStore() {
//
//}
//

//
//void DeleteContextMenuStore() {
//    // Free radio group members
//    if( hashmap_num_entries(&contextMenuStore) > 0 ) {
//        if (0 != hashmap_iterate_pairs(&contextMenuStore, freeContextMenu, NULL)) {
//            ABORT("[DeleteContextMenuStore] Failed to release contextMenuStore entries!");
//        }
//    }
//}
//

void ShowContextMenu(ContextMenuStore* store, id mainWindow, const char *contextMenuID, const char *contextMenuData) {

	// If no context menu ID was given, abort
	if( contextMenuID == NULL ) {
		return;
	}

	ContextMenu* contextMenu = GetContextMenuByID(store, contextMenuID);

	// We don't need the ID now
    MEMFREE(contextMenuID);

    if( contextMenu == NULL ) {
        // Free context menu data
        if( contextMenuData != NULL ) {}
        MEMFREE(contextMenuData);
		return;
	}

    // We need to store the context menu data. Free existing data if we have it
    // and set to the new value.
    FREE_AND_SET(store->contextMenuData, contextMenuData);

	// Grab the content view and show the menu
	id contentView = msg(mainWindow, s("contentView"));

	// Get the triggering event
	id menuEvent = msg(mainWindow, s("currentEvent"));

	if( contextMenu->nsmenu == NULL ) {
        // GetMenu creates the NSMenu
        contextMenu->nsmenu = GetMenu(contextMenu->menu);
	}

	// Show popup
	msg(c("NSMenu"), s("popUpContextMenu:withEvent:forView:"), contextMenu->nsmenu, menuEvent, contentView);

}

#endif //CONTEXTMENU_DARWIN_H
