//
// Created by Lea Anthony on 7/1/21.
//

#ifndef CONTEXTMENUSTORE_DARWIN_H
#define CONTEXTMENUSTORE_DARWIN_H

typedef struct {
    // This is our context menu store which keeps track
    // of all instances of ContextMenus
    struct hashmap_s contextMenuStore;

    // The raw JSON defining the context menus
    const char* contextMenusAsJSON;

    // The optional data that may be passed with a context menu selection
    const char* contextMenuData;

    // The processed context menus
    JsonNode* processedContextMenus;

} ContextMenuStore;

#endif //ASSETS_C_CONTEXTMENUSTORE_DARWIN_H
