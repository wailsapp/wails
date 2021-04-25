// Some code may be inspired by or directly used from Webview.
#include "ffenestri_windows.h"

int debug = 0;

struct Application{
    // Window specific
    HWND window;
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
};

struct Application *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel, int hideWindowOnClose) {

    // Create application
    struct Application *result = malloc(sizeof(struct Application));

    result->title = title;
    result->width = width;
    result->height = height;
    result->resizable = resizable;
    result->devtools = devtools;
    result->fullscreen = fullscreen;
    result->startHidden = startHidden;
    result->logLevel = logLevel;
    result->hideWindowOnClose = hideWindowOnClose;

    // Min/Max Width/Height
    result->minWidth = 0;
    result->minHeight = 0;
    result->maxWidth = 0;
    result->maxHeight = 0;

    return result;
}

void SetMinWindowSize(struct Application* app, int minWidth, int minHeight) {
    app->minWidth = (LONG)minWidth;
    app->minHeight = (LONG)minHeight;
}

void SetMaxWindowSize(struct Application* app, int maxWidth, int maxHeight) {
    app->maxWidth = (LONG)maxWidth;
    app->maxHeight = (LONG)maxHeight;
}

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {

    struct Application *app = (struct Application *)GetWindowLongPtr(hwnd, GWLP_USERDATA);

    switch(msg) {

        case WM_DESTROY: {
            DestroyApplication(app);
            break;
        }
        case WM_GETMINMAXINFO: {
            // Exit early if this is called before the window is created.
            if ( app == NULL ) {
                return 0;
            }
            LPMINMAXINFO mmi = (LPMINMAXINFO) lParam;
            if (app->minWidth > 0 && app->minHeight > 0) {
                mmi->ptMinTrackSize.x = app->minWidth;
                mmi->ptMinTrackSize.y = app->minHeight;
            }
            if (app->maxWidth > 0 && app->maxHeight > 0) {
                mmi->ptMaxSize.x = app->maxWidth;
                mmi->ptMaxSize.y = app->maxHeight;
                mmi->ptMaxTrackSize.x = app->maxWidth;
                mmi->ptMaxTrackSize.y = app->maxHeight;
            }
            return 0;
        }
        default:
            return DefWindowProc(hwnd, msg, wParam, lParam);
    }
}

void Run(struct Application* app, int argc, char **argv) {

    WNDCLASSEX wc;
    HINSTANCE hInstance = GetModuleHandle(NULL);
    ZeroMemory(&wc, sizeof(WNDCLASSEX));
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.hInstance = hInstance;
    wc.lpszClassName = "ffenestri";
    wc.lpfnWndProc   = WndProc;

    // Process window resizable
    DWORD dwStyle = WS_OVERLAPPEDWINDOW;
    if (app->resizable == 0) {
      dwStyle = WS_OVERLAPPED | WS_MINIMIZEBOX | WS_SYSMENU;
    }
    RegisterClassEx(&wc);
    app->window = CreateWindow("ffenestri", "", dwStyle, CW_USEDEFAULT,
                                      CW_USEDEFAULT, app->width, app->height, NULL, NULL,
                                      GetModuleHandle(NULL), NULL);
    // Set Title
    SetWindowText(app->window, app->title);

    // Store application pointer in window handle
    SetWindowLongPtr(app->window, GWLP_USERDATA, (LONG_PTR)app);

    // Process whether window should show by default
    int startVisibility = SW_SHOWNORMAL;
    if ( app->startHidden == 1 ) {
        startVisibility = SW_HIDE;
    }
    ShowWindow(app->window, startVisibility);
    UpdateWindow(app->window);
    SetFocus(app->window);

    // TODO: Add webview2

    // Main event loop
    MSG  msg;
    BOOL res;
    while ((res = GetMessage(&msg, NULL, 0, 0)) != -1) {
      if (msg.hwnd) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
        continue;
      }
      if (msg.message == WM_APP) {
      } else if (msg.message == WM_QUIT) {
        return;
      }
    }
}

void DestroyApplication(struct Application* app) {
    PostQuitMessage(0);
}
void SetDebug(struct Application* app, int flag) {
    debug = flag;
}

void SetBindings(struct Application* app, const char *bindings) {
}

void ExecJS(struct Application* app, const char *script) {
}

void Hide(struct Application* app) {
    ShowWindow(app->window, SW_HIDE);
}

void Show(struct Application* app) {
    ShowWindow(app->window, SW_SHOW);
}

void Center(struct Application* app) {
}

UINT getWindowPlacement(struct Application* app) {
    WINDOWPLACEMENT lpwndpl;
    lpwndpl.length = sizeof(WINDOWPLACEMENT);
    BOOL result = GetWindowPlacement(app->window, &lpwndpl);
    if( result == 0 ) {
        // TODO: Work out what this call failing means
        return -1;
    }
    return lpwndpl.showCmd;
}

int isMaximised(struct Application* app) {
    return getWindowPlacement(app) == SW_SHOWMAXIMIZED;
}

void Maximise(struct Application* app) {
    ShowWindow(app->window, SW_MAXIMIZE);
}

void Unmaximise(struct Application* app) {
    ShowWindow(app->window, SW_RESTORE);
}

void ToggleMaximise(struct Application* app) {
    if(isMaximised(app)) {
        return Unmaximise(app);
    }
    return Maximise(app);
}

int isMinimised(struct Application* app) {
    return getWindowPlacement(app) == SW_SHOWMINIMIZED;
}

void Minimise(struct Application* app) {
    ShowWindow(app->window, SW_MINIMIZE);
}

void Unminimise(struct Application* app) {
    ShowWindow(app->window, SW_RESTORE);
}

void ToggleMinimise(struct Application* app) {
    if(isMinimised(app)) {
        return Unminimise(app);
    }
    return Minimise(app);
}

void SetColour(struct Application* app, int red, int green, int blue, int alpha) {
}
void SetSize(struct Application* app, int width, int height) {
}
void SetPosition(struct Application* app, int x, int y) {
}
void Quit(struct Application* app) {
}
void SetTitle(struct Application* app, const char *title) {
    SetWindowText(app->window, title);
}
void Fullscreen(struct Application* app) {
}
void UnFullscreen(struct Application* app) {
}
void ToggleFullscreen(struct Application* app) {
}
void DisableFrame(struct Application* app) {
}
void OpenDialog(struct Application* app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolvesAliases, int treatPackagesAsDirectories) {
}
void SaveDialog(struct Application* app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories) {
}
void MessageDialog(struct Application* app, char *callbackID, char *type, char *title, char *message, char *icon, char *button1, char *button2, char *button3, char *button4, char *defaultButton, char *cancelButton) {
}
void DarkModeEnabled(struct Application* app, char *callbackID) {
}
void SetApplicationMenu(struct Application* app, const char *applicationMenuJSON) {
}
void AddTrayMenu(struct Application* app, const char *menuTrayJSON) {
}
void SetTrayMenu(struct Application* app, const char *menuTrayJSON) {
}
void DeleteTrayMenuByID(struct Application* app, const char *id) {
}
void UpdateTrayMenuLabel(struct Application* app, const char* JSON) {
}
void AddContextMenu(struct Application* app, char *contextMenuJSON) {
}
void UpdateContextMenu(struct Application* app, char *contextMenuJSON) {
}