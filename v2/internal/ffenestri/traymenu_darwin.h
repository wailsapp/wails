//
// Created by Lea Anthony on 12/1/21.
//

#ifndef TRAYMENU_DARWIN_H
#define TRAYMENU_DARWIN_H

#include "common.h"
#include "menu_darwin.h"

typedef struct {

    const char *label;
    const char *icon;
    const char *ID;

    Menu* menu;

    id statusbaritem;
    int trayIconPosition;

    JsonNode* processedJSON;

} TrayMenu;

TrayMenu* NewTrayMenu(const char *trayJSON);
void DumpTrayMenu(TrayMenu* trayMenu);
void ShowTrayMenu(TrayMenu* trayMenu);
void UpdateTrayMenuInPlace(TrayMenu* currentMenu, TrayMenu* newMenu);
void UpdateTrayIcon(TrayMenu *trayMenu);
void UpdateTrayLabel(TrayMenu *trayMenu, const char*);

void LoadTrayIcons();
void UnloadTrayIcons();

void DeleteTrayMenu(TrayMenu* trayMenu);

#endif //TRAYMENU_DARWIN_H
