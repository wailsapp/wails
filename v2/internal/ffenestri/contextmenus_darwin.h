////
//// Created by Lea Anthony on 6/1/21.
////
//
#ifndef CONTEXTMENU_DARWIN_H
#define CONTEXTMENU_DARWIN_H

#include "json.h"
#include "menu_darwin.h"
#include "contextmenustore_darwin.h"

typedef struct {
    const char* ID;
    id nsmenu;
    Menu* menu;

    JsonNode* processedJSON;

    // Context menu data is given by the frontend when clicking a context menu.
    // We send this to the backend when an item is selected
	const char* contextMenuData;
} ContextMenu;


ContextMenu* NewContextMenu(const char* contextMenuJSON);

ContextMenu* GetContextMenuByID( ContextMenuStore* store, const char *contextMenuID);
void DeleteContextMenu(ContextMenu* contextMenu);
int freeContextMenu(void *const context, struct hashmap_element_s *const e);

void ShowContextMenu(ContextMenuStore* store, id mainWindow, const char *contextMenuID, const char *contextMenuData);

#endif //CONTEXTMENU_DARWIN_H
