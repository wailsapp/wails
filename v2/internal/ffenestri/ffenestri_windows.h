
#ifndef _FFENESTRI_WINDOWS_H
#define _FFENESTRI_WINDOWS_H

#define UNICODE 1

#include "ffenestri.h"
#include <windows.h>
#include <wingdi.h>
#include <functional>

#define ON_MAIN_THREAD(code) dispatch( [=]{ code; } )

typedef std::function<void()> dispatchFunction;
//typedef std::function<void(const std::string)> messageCallback;
//typedef std::function<void(ICoreWebView2Controller *)> comHandlerCallback;

void center(struct Application*);
void setTitle(struct Application* app, const char *title);

#endif