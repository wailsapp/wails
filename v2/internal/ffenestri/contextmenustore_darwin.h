//
// Created by Lea Anthony on 7/1/21.
//

#ifndef CONTEXTMENUSTORE_DARWIN_H
#define CONTEXTMENUSTORE_DARWIN_H

#include "common.h"

typedef struct {

    int dummy;

    // This is our context menu store which keeps track
    // of all instances of ContextMenus
    struct hashmap_s contextMenuMap;

} ContextMenuStore;

ContextMenuStore* NewContextMenuStore();

void DeleteContextMenuStore(ContextMenuStore* store);
void UpdateContextMenuInStore(ContextMenuStore* store, const char* menuJSON);

void AddContextMenuToStore(ContextMenuStore* store, const char* contextMenuJSON);

#endif //CONTEXTMENUSTORE_DARWIN_H
