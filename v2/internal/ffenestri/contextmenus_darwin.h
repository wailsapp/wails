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
} ContextMenu;


ContextMenu* NewContextMenu(JsonNode* menuData, ContextMenuStore* store);

ContextMenu* GetContextMenuByID( ContextMenuStore* store, const char *contextMenuID);
void DeleteContextMenu(ContextMenu* contextMenu);
int freeContextMenu(void *const context, struct hashmap_element_s *const e);
void ProcessContextMenus( ContextMenuStore* store);

void ShowContextMenu(ContextMenuStore* store, id mainWindow, const char *contextMenuID, const char *contextMenuData);

#endif //CONTEXTMENU_DARWIN_H
