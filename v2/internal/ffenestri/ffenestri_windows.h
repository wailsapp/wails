
#ifndef _FFENESTRI_WINDOWS_H
#define _FFENESTRI_WINDOWS_H

#define UNICODE 1

#include "ffenestri.h"
#include <windows.h>
#include <wingdi.h>

void center(struct Application*);
void setTitle(struct Application* app, const char *title);

#endif