//
// Created by Lea Anthony on 12/1/21.
//

#ifndef TRAYMENU_DARWIN_H
#define TRAYMENU_DARWIN_H

#include "menu_darwin.h"

typedef struct {

    const char *label;
    const char *icon;
    const char *trayID;

    Menu* menu;

} TrayMenu;

#endif //TRAYMENU_DARWIN_H
