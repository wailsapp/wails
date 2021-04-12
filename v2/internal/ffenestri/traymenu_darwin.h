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
    const char *tooltip;

    bool templateImage;
    const char *fontName;
    int fontSize;
    const char *RGBA;

    bool disabled;

    Menu* menu;

    id statusbaritem;
    unsigned int trayIconPosition;

    JsonNode* processedJSON;

    JsonNode* styledLabel;

    id delegate;

} TrayMenu;

TrayMenu* NewTrayMenu(const char *trayJSON);
void DumpTrayMenu(TrayMenu* trayMenu);
void ShowTrayMenu(TrayMenu* trayMenu);
void UpdateTrayMenuInPlace(TrayMenu* currentMenu, TrayMenu* newMenu);
void UpdateTrayIcon(TrayMenu *trayMenu);
void UpdateTrayLabel(TrayMenu *trayMenu, const char *label, const char *fontName, int fontSize, const char *RGBA, const char *tooltip, bool disabled, JsonNode *styledLabel);

void LoadTrayIcons();
void UnloadTrayIcons();

void DeleteTrayMenu(TrayMenu* trayMenu);

#endif //TRAYMENU_DARWIN_H
