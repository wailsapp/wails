
#ifndef _FFENESTRI_WINDOWS_H
#define _FFENESTRI_WINDOWS_H

#define WIN32_LEAN_AND_MEAN
#define UNICODE 1

#include "ffenestri.h"
#include <windows.h>
#include <wingdi.h>
#include <functional>
#include <codecvt>
#include "windows/WebView2.h"

#include "assets.h"

// TODO:
//#include "userdialogicons.h"


struct Application{
    // Window specific
    HWND window;
    ICoreWebView2 *webview;
    ICoreWebView2Controller* webviewController;

    // Application
    const char *title;
    int width;
    int height;
    int resizable;
    int devtools;
    int fullscreen;
    int startHidden;
    int logLevel;
    int hideWindowOnClose;
    int minSizeSet;
    LONG minWidth;
    LONG minHeight;
    int maxSizeSet;
    LONG maxWidth;
    LONG maxHeight;
    int frame;

    char* bindings;
};

#define ON_MAIN_THREAD(code) dispatch( [=]{ code; } )

typedef std::function<void()> dispatchFunction;
typedef std::function<void(const std::string)> messageCallback;
typedef std::function<void(ICoreWebView2Controller *)> comHandlerCallback;

void center(struct Application*);
void setTitle(struct Application* app, const char *title);
char* LPWSTRToCstr(LPWSTR input);
void loadAssets(struct Application* app);

// Callback
extern "C" {
    void messageFromWindowCallback(const char *);
}

#endif