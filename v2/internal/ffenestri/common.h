//
// Created by Lea Anthony on 6/1/21.
//

#ifndef COMMON_H
#define COMMON_H

#include "hashmap.h"

void ABORT(const char *message) {
    printf("%s\n", message);
    exit(1);
}

int freeHashmapItem(void *const context, struct hashmap_element_s *const e) {
    free(e->data);
    return -1;
}

#endif //ASSETS_C_COMMON_H
