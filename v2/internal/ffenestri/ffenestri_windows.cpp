// Some code may be inspired by or directly used from Webview.
#include "ffenestri_windows.h"
#include "wv2ComHandler_windows.h"
#include <functional>
#include <atomic>
#include <Shlwapi.h>
#include <locale>
#include <codecvt>
#include "windows/WebView2.h"

int debug = 0;
DWORD mainThread;

// dispatch will execute the given `func` pointer
void dispatch(dispatchFunction func) {
    PostThreadMessage(mainThread, WM_APP, 0, (LPARAM) new dispatchFunction(func));
}

struct Application *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel, int hideWindowOnClose) {

    // Create application
    struct Application *result = (struct Application*)malloc(sizeof(struct Application));

    result->window = nullptr;
    result->webview = nullptr;
    result->webviewController = nullptr;

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

    // Have a frame by default
    result->frame = 1;

    // Capture Main Thread
    mainThread = GetCurrentThreadId();

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
        case WM_SIZE: {
            if( app->webviewController != nullptr) {
                RECT bounds;
                GetClientRect(app->window, &bounds);
                app->webviewController->put_Bounds(bounds);
            }
            break;
        }
        case WM_GETMINMAXINFO: {
            // Exit early if this is called before the window is created.
            if ( app == NULL ) {
                return 0;
            }

            // get pixel density
            HDC hDC = GetDC(NULL);
            double DPIScaleX = GetDeviceCaps(hDC, 88)/96.0;
            double DPIScaleY = GetDeviceCaps(hDC, 90)/96.0;
            ReleaseDC(NULL, hDC);

            RECT rcClient, rcWind;
            POINT ptDiff;
            GetClientRect(hwnd, &rcClient);
            GetWindowRect(hwnd, &rcWind);

            int widthExtra = (rcWind.right - rcWind.left) - rcClient.right;
            int heightExtra = (rcWind.bottom - rcWind.top) - rcClient.bottom;

            LPMINMAXINFO mmi = (LPMINMAXINFO) lParam;
            if (app->minWidth > 0 && app->minHeight > 0) {
                mmi->ptMinTrackSize.x = app->minWidth * DPIScaleX + widthExtra;
                mmi->ptMinTrackSize.y = app->minHeight * DPIScaleY + heightExtra;
            }
            if (app->maxWidth > 0 && app->maxHeight > 0) {
                mmi->ptMaxSize.x = app->maxWidth * DPIScaleX + widthExtra;
                mmi->ptMaxSize.y = app->maxHeight * DPIScaleY + heightExtra;
                mmi->ptMaxTrackSize.x = app->maxWidth * DPIScaleX + widthExtra;
                mmi->ptMaxTrackSize.y = app->maxHeight * DPIScaleY + heightExtra;
            }
            return 0;
        }
        default:
            return DefWindowProc(hwnd, msg, wParam, lParam);
    }
    return 0;
}
//
bool initWebView2(struct Application *app, int debug, messageCallback cb) {
    CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);

    std::atomic_flag flag = ATOMIC_FLAG_INIT;
    flag.test_and_set();

    char currentExePath[MAX_PATH];
    GetModuleFileNameA(NULL, currentExePath, MAX_PATH);
    char *currentExeName = PathFindFileNameA(currentExePath);


//    std::wstring_convert<std::codecvt_utf8_utf16<wchar_t>> wideCharConverter;
//    std::wstring userDataFolder = wideCharConverter.from_bytes(std::getenv("APPDATA"));
//    std::wstring currentExeNameW = wideCharConverter.from_bytes(currentExeName);

//    printf("userdata folder = %s\n", userDataFolder.c_str());

    ICoreWebView2Controller *controller;
    ICoreWebView2* webview;

    HRESULT res = CreateCoreWebView2EnvironmentWithOptions(
            nullptr, /*(userDataFolder + L"/" + currentExeNameW).c_str()*/ nullptr, nullptr,
            new wv2ComHandler(app->window, cb,
                                     [&](ICoreWebView2Controller *webviewController) {
                                         controller = webviewController;
                                         controller->get_CoreWebView2(&webview);
                                         webview->AddRef();
                                         flag.clear();
                                     }));
    if (res != S_OK) {
        CoUninitialize();
        return false;
    }

    MSG msg = {};
    while (flag.test_and_set() && GetMessage(&msg, NULL, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }
    app->webviewController = controller;
    app->webview = webview;
    // Resize WebView to fit the bounds of the parent window
    RECT bounds;
    GetClientRect(app->window, &bounds);
    app->webviewController->put_Bounds(bounds);

    // Schedule an async task to navigate to Bing
    app->webview->Navigate(L"https://wails.app/");
////    init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}");
    return true;
}

void initialCallback(std::string message) {
    printf("MESSAGE=%s\n", message);
}

void Run(struct Application* app, int argc, char **argv) {

    WNDCLASSEX wc;
    HINSTANCE hInstance = GetModuleHandle(NULL);
    ZeroMemory(&wc, sizeof(WNDCLASSEX));
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.style = CS_HREDRAW | CS_VREDRAW;
    wc.hInstance = hInstance;
    wc.lpszClassName = (LPCWSTR)"ffenestri";
    wc.lpfnWndProc   = WndProc;

    // TODO: Menu
//    wc.lpszMenuName = nullptr;


    // Process window resizable
    DWORD windowStyle = WS_OVERLAPPEDWINDOW;
    if (app->resizable == 0) {
        windowStyle &= ~WS_MAXIMIZEBOX;
        windowStyle &= ~WS_THICKFRAME;
    }
    if ( app->frame == 0 ) {
        windowStyle = WS_POPUP;
    }

    RegisterClassEx(&wc);
    app->window = CreateWindow((LPCWSTR)"ffenestri", (LPCWSTR)"", windowStyle, CW_USEDEFAULT,
                                CW_USEDEFAULT, app->width, app->height, NULL, NULL,
                                hInstance, NULL);

    // Private setTitle as we're on the main thread
    setTitle(app, app->title);

    // Store application pointer in window handle
    SetWindowLongPtr(app->window, GWLP_USERDATA, (LONG_PTR)app);

    // Process whether window should show by default
    int startVisibility = SW_SHOWNORMAL;
    if ( app->startHidden == 1 ) {
        startVisibility = SW_HIDE;
    }

    // private center() as we are on main thread
    center(app);
//    if( debug == 1 ) {
//        BOOL supported = SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
//        if( !supported ) {
//            SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE);
//        }
//    }
    ShowWindow(app->window, startVisibility);
    UpdateWindow(app->window);
    SetFocus(app->window);

    // Add webview2
    initWebView2(app, 1, initialCallback);

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
          dispatchFunction *f = (dispatchFunction*) msg.lParam;
          (*f)();
          delete(f);
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

void hide(struct Application* app) {
    ShowWindow(app->window, SW_HIDE);
}

void Hide(struct Application* app) {
    ON_MAIN_THREAD(
        hide(app);
    );
}

void show(struct Application* app) {
    ShowWindow(app->window, SW_SHOW);
}

void Show(struct Application* app) {
    ON_MAIN_THREAD(
        show(app);
    );
}

void center(struct Application* app) {

    HMONITOR currentMonitor = MonitorFromWindow(app->window, MONITOR_DEFAULTTONEAREST);
    MONITORINFO info = {0};
    info.cbSize = sizeof(info);
    GetMonitorInfoA(currentMonitor, &info);
    RECT workRect = info.rcWork;
    LONG screenMiddleW = (workRect.right - workRect.left) / 2;
    LONG screenMiddleH = (workRect.bottom - workRect.top) / 2;
    RECT winRect;
    if (app->frame == 1) {
        GetWindowRect(app->window, &winRect);
    } else {
        GetClientRect(app->window, &winRect);
    }
    LONG winWidth = winRect.right - winRect.left;
    LONG winHeight = winRect.bottom - winRect.top;

    LONG windowX = screenMiddleW - (winWidth / 2);
    LONG windowY = screenMiddleH - (winHeight / 2);

    SetWindowPos(app->window, HWND_TOP, windowX, windowY, winWidth, winHeight, SWP_NOSIZE);
}

void Center(struct Application* app) {
    ON_MAIN_THREAD(
        center(app);
    );
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

void maximise(struct Application* app) {
    ShowWindow(app->window, SW_MAXIMIZE);
}

void Maximise(struct Application* app) {
    ON_MAIN_THREAD(
        maximise(app);
    );
}

void unmaximise(struct Application* app) {
    ShowWindow(app->window, SW_RESTORE);
}

void Unmaximise(struct Application* app) {
    ON_MAIN_THREAD(
        unmaximise(app);
    );
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

void minimise(struct Application* app) {
    ShowWindow(app->window, SW_MINIMIZE);
}

void Minimise(struct Application* app) {
    ON_MAIN_THREAD(
        minimise(app);
    );
}

void unminimise(struct Application* app) {
    ShowWindow(app->window, SW_RESTORE);
}

void Unminimise(struct Application* app) {
    ON_MAIN_THREAD(
        unminimise(app);
    );
}

void ToggleMinimise(struct Application* app) {
    if(isMinimised(app)) {
        return Unminimise(app);
    }
    return Minimise(app);
}

void SetColour(struct Application* app, int red, int green, int blue, int alpha) {
// TBD
}

void SetSize(struct Application* app, int width, int height) {
    // TBD
}

void setPosition(struct Application* app, int x, int y) {
    // TBD
}

void SetPosition(struct Application* app, int x, int y) {
    ON_MAIN_THREAD(
        setPosition(app, x, y);
    );
}

void Quit(struct Application* app) {
}

// Credit: https://stackoverflow.com/a/6693107
void setTitle(struct Application* app, const char *title) {
    int wchars_num = MultiByteToWideChar( CP_UTF8 , 0 , title , -1, NULL , 0 );
    wchar_t* wstr = new wchar_t[wchars_num];
    MultiByteToWideChar( CP_UTF8 , 0 , title , -1, wstr , wchars_num );
    SetWindowText(app->window, wstr);
    delete[] wstr;
}

void SetTitle(struct Application* app, const char *title) {
    ON_MAIN_THREAD(
        setTitle(app, title);
    );
}

void Fullscreen(struct Application* app) {
}

void UnFullscreen(struct Application* app) {
}

void ToggleFullscreen(struct Application* app) {
}

void DisableFrame(struct Application* app) {
    app->frame = 0;
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