//
// Created by Lea Anthony on 12/1/21.
//

#include "common.h"
#include "traymenu_darwin.h"


TrayMenu* NewTrayMenu(const char* menuJSON) {
    TrayMenu* result = malloc(sizeof(TrayMenu));




    return result;
}

void DeleteTrayMenu(TrayMenu* trayMenu) {
    // Free the tray menu memory
    MEMFREE(trayMenu);
}