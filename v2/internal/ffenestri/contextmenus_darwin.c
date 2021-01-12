////
//// Created by Lea Anthony on 6/1/21.
////
//

#include "ffenestri_darwin.h"
#include "common.h"
#include "contextmenus_darwin.h"
#include "menu_darwin.h"

ContextMenu* NewContextMenu(JsonNode* menuData, ContextMenuStore *store) {
    ContextMenu* result = malloc(sizeof(ContextMenu));
    result->menu = NewMenu(menuData);
    result->nsmenu = NULL;
    result->menu->menuType = ContextMenuType;
    result->menu->parentData = result;
    result->contextMenuData = NULL;
    return result;
}


ContextMenu* GetContextMenuByID(ContextMenuStore* store, const char *contextMenuID) {
    return (ContextMenu*)hashmap_get(&store->contextMenuStore, (char*)contextMenuID, strlen(contextMenuID));
}

void DeleteContextMenu(ContextMenu* contextMenu) {
    // Free Menu
    DeleteMenu(contextMenu->menu);

    // Delete any context menu data we may have stored
    if( contextMenu->contextMenuData != NULL ) {
        MEMFREE(contextMenu->contextMenuData);
    }

    // Free context menu
    free(contextMenu);
}

int freeContextMenu(void *const context, struct hashmap_element_s *const e) {
    DeleteContextMenu(e->data);
    return -1;
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
        thisContextMenu->ID = ID;

		// Store the item in the context menu map
		hashmap_put(&store->contextMenuStore, (char*)ID, strlen(ID), thisContextMenu);
	}

}

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
        if( contextMenuData != NULL ) {
	        MEMFREE(contextMenuData);
			return;
		}
	}

    // We need to store the context menu data. Free existing data if we have it
    // and set to the new value.
    FREE_AND_SET(contextMenu->contextMenuData, contextMenuData);

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

