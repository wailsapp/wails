//
// Created by Lea Anthony on 7/1/21.
//

#ifndef CONTEXTMENUSTORE_DARWIN_H
#define CONTEXTMENUSTORE_DARWIN_H

#include "common.h"

typedef struct {
    // This is our context menu store which keeps track
    // of all instances of ContextMenus
    struct hashmap_s contextMenuStore;

    // The raw JSON defining the context menus
    const char* contextMenusAsJSON;

    // The processed context menus
    JsonNode* processedContextMenus;

} ContextMenuStore;

ContextMenuStore* NewContextMenuStore(const char* contextMenusAsJSON);
void DeleteContextMenuStore(ContextMenuStore* store);

#endif //CONTEXTMENUSTORE_DARWIN_H
