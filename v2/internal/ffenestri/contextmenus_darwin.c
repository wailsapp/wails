////
//// Created by Lea Anthony on 6/1/21.
////
//

#include "ffenestri_darwin.h"
#include "common.h"
#include "contextmenus_darwin.h"
#include "menu_darwin.h"

ContextMenu* NewContextMenu(const char* contextMenuJSON) {
    ContextMenu* result = malloc(sizeof(ContextMenu));

    JsonNode* processedJSON = json_decode(contextMenuJSON);
    if( processedJSON == NULL ) {
        ABORT("[NewTrayMenu] Unable to parse TrayMenu JSON: %s", contextMenuJSON);
    }
    // Save reference to this json
    result->processedJSON = processedJSON;

    result->ID = mustJSONString(processedJSON, "ID");
    JsonNode* processedMenu = mustJSONObject(processedJSON, "ProcessedMenu");

    result->menu = NewMenu(processedMenu);
    result->nsmenu = NULL;
    result->menu->menuType = ContextMenuType;
    result->menu->parentData = result;
    result->contextMenuData = NULL;
    return result;
}

ContextMenu* GetContextMenuByID(ContextMenuStore* store, const char *contextMenuID) {
    return (ContextMenu*)hashmap_get(&store->contextMenuMap, (char*)contextMenuID, strlen(contextMenuID));
}

void DeleteContextMenu(ContextMenu* contextMenu) {
    // Free Menu
    DeleteMenu(contextMenu->menu);

    // Delete any context menu data we may have stored
    if( contextMenu->contextMenuData != NULL ) {
        MEMFREE(contextMenu->contextMenuData);
    }

    // Free JSON
    if (contextMenu->processedJSON != NULL ) {
        json_delete(contextMenu->processedJSON);
        contextMenu->processedJSON = NULL;
    }

    // Free context menu
    free(contextMenu);
}

int freeContextMenu(void *const context, struct hashmap_element_s *const e) {
    DeleteContextMenu(e->data);
    return -1;
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
	id contentView = msg_reg(mainWindow, s("contentView"));

	// Get the triggering event
	id menuEvent = msg_reg(mainWindow, s("currentEvent"));

	if( contextMenu->nsmenu == NULL ) {
        // GetMenu creates the NSMenu
        contextMenu->nsmenu = GetMenu(contextMenu->menu);
	}

	// Show popup
	((id(*)(id, SEL, id, id, id))objc_msgSend)(c("NSMenu"), s("popUpContextMenu:withEvent:forView:"), contextMenu->nsmenu, menuEvent, contentView);

}

