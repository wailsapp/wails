// Some code may be inspired by or directly used from Webview (c) zserge.
// License included in README.md

#include "ffenestri_windows.h"
#include "wv2ComHandler_windows.h"
#include <functional>
#include <atomic>
#include <Shlwapi.h>
#include <locale>
#include <codecvt>
#include "windows/WebView2.h"
#include <winuser.h>
#include "effectstructs_windows.h"
#include <Shlobj.h>

int debug = 0;
DWORD mainThread;

#define WS_EX_NOREDIRECTIONBITMAP 0x00200000L

// --- Assets
extern const unsigned char runtime;
extern const unsigned char *defaultDialogIcons[];

// dispatch will execute the given `func` pointer
void dispatch(dispatchFunction func) {
    PostThreadMessage(mainThread, WM_APP, 0, (LPARAM) new dispatchFunction(func));
}

LPWSTR cstrToLPWSTR(const char *cstr) {
    int wchars_num = MultiByteToWideChar( CP_UTF8 , 0 , cstr , -1, NULL , 0 );
    wchar_t* wstr = new wchar_t[wchars_num+1];
    MultiByteToWideChar( CP_UTF8 , 0 , cstr , -1, wstr , wchars_num );
    return wstr;
}

// Credit: https://stackoverflow.com/a/9842450
char* LPWSTRToCstr(LPWSTR input) {
    int length = WideCharToMultiByte(CP_UTF8, 0, input, -1, 0, 0, NULL, NULL);
    char* output = new char[length];
    WideCharToMultiByte(CP_UTF8, 0, input, -1, output , length, NULL, NULL);
    return output;
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
    result->webviewIsTranparent = false;
    result->windowBackgroundIsTranslucent = false;
    result->disableWindowIcon = false;

    // Min/Max Width/Height
    result->minWidth = 0;
    result->minHeight = 0;
    result->maxWidth = 0;
    result->maxHeight = 0;

    // Default colour
    result->backgroundColour.R = 255;
    result->backgroundColour.G = 255;
    result->backgroundColour.B = 255;
    result->backgroundColour.A = 255;

    // Have a frame by default
    result->frame = 1;

    // Capture Main Thread
    mainThread = GetCurrentThreadId();

    // Startup url
    result->startupURL = nullptr;

    // Used to remember the window location when going fullscreen
    result->previousPlacement = { sizeof(result->previousPlacement) };

    return result;
}

void* GetWindowHandle(struct Application *app) {
    return (void*)app->window;
}

void SetMinWindowSize(struct Application* app, int minWidth, int minHeight) {
    app->minWidth = (LONG)minWidth;
    app->minHeight = (LONG)minHeight;
}

void SetMaxWindowSize(struct Application* app, int maxWidth, int maxHeight) {
    app->maxWidth = (LONG)maxWidth;
    app->maxHeight = (LONG)maxHeight;
}

void SetBindings(struct Application *app, const char *bindings) {
    std::string temp = std::string("window.wailsbindings = \"") + std::string(bindings) + std::string("\";");
    app->bindings = new char[temp.length()+1];
	memcpy(app->bindings, temp.c_str(), temp.length()+1);
}

void performShutdown(struct Application *app) {
    if( app->startupURL != nullptr ) {
        delete[] app->startupURL;
    }
    messageFromWindowCallback("WC");
}

// Credit: https://gist.github.com/ysc3839/b08d2bff1c7dacde529bed1d37e85ccf
void enableTranslucentBackground(struct Application *app) {
    HMODULE hUser = GetModuleHandleA("user32.dll");
    if (hUser)
    {
        pfnSetWindowCompositionAttribute setWindowCompositionAttribute = (pfnSetWindowCompositionAttribute)GetProcAddress(hUser, "SetWindowCompositionAttribute");
        if (setWindowCompositionAttribute)
        {
            ACCENT_POLICY accent = { ACCENT_ENABLE_BLURBEHIND, 0, 0, 0 };
            WINDOWCOMPOSITIONATTRIBDATA data;
            data.Attrib = WCA_ACCENT_POLICY;
            data.pvData = &accent;
            data.cbData = sizeof(accent);
            setWindowCompositionAttribute(app->window, &data);
        }
    }
}

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {

    struct Application *app = (struct Application *)GetWindowLongPtr(hwnd, GWLP_USERDATA);

    switch(msg) {

        case WM_CREATE: {
            createApplicationMenu(hwnd);
            break;
        }
        case WM_COMMAND:
            menuClicked(LOWORD(wParam));
        break;

        case WM_CLOSE: {
            DestroyWindow( app->window );
            break;
        }
        case WM_DESTROY: {
            if( app->hideWindowOnClose ) {
                Hide(app);
            } else {
                PostQuitMessage(0);
            }
            break;
        }
        case WM_SIZE: {
            if ( app == NULL ) {
                return 0;
            }
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

void init(struct Application *app, const char* js) {
    LPCWSTR wjs = cstrToLPWSTR(js);
    app->webview->AddScriptToExecuteOnDocumentCreated(wjs, nullptr);
    delete[] wjs;
}

void execJS(struct Application* app, const char *script) {
    LPWSTR s = cstrToLPWSTR(script);
    app->webview->ExecuteScript(s, nullptr);
    delete[] s;
}

void loadAssets(struct Application* app) {

    // setup window.wailsInvoke
    std::string initialCode = std::string("window.wailsInvoke=function(m){window.chrome.webview.postMessage(m)};");

    // Load bindings
    initialCode += std::string(app->bindings);
    delete[] app->bindings;

    // Load runtime
    initialCode += std::string((const char*)&runtime);

    int index = 1;
    while(1) {
        // Get next asset pointer
        const unsigned char *asset = assets[index];

        // If we have no more assets, break
        if (asset == 0x00) {
            break;
        }

        initialCode += std::string((const char*)asset);
        index++;
    };

    // Disable context menu if not in debug mode
    if( debug != 1 ) {
        initialCode += std::string("wails._.DisableDefaultContextMenu();");
    }

    initialCode += std::string("window.wailsInvoke('completed');");

    // Keep a copy of the code
    app->initialCode = new char[initialCode.length()+1];
	memcpy(app->initialCode, initialCode.c_str(), initialCode.length()+1);

    execJS(app, app->initialCode);

    // Show app if we need to
    if( app->startHidden == false ) {
        Show(app);
    }
}

// This is called when all our assets are loaded into the DOM
void completed(struct Application* app) {
    delete[] app->initialCode;
    app->initialCode = nullptr;

    // Process whether window should show by default
    int startVisibility = SW_SHOWNORMAL;
    if ( app->startHidden == 1 ) {
        startVisibility = SW_HIDE;
    }

    // Fix for webview2 bug: https://github.com/MicrosoftEdge/WebView2Feedback/issues/1077
    // Will be fixed in next stable release
    app->webviewController->put_IsVisible(false);
    app->webviewController->put_IsVisible(true);

    // Private setTitle as we're on the main thread
    if( app->frame == 1) {
        setTitle(app, app->title);
    }

    ShowWindow(app->window, startVisibility);
    UpdateWindow(app->window);
    SetFocus(app->window);

    if( app->startupURL == nullptr ) {
        messageFromWindowCallback("SS");
        return;
    }
    std::string readyMessage = std::string("SS") + std::string(app->startupURL);
    messageFromWindowCallback(readyMessage.c_str());
}


//
bool initWebView2(struct Application *app, int debugEnabled, messageCallback cb) {

    debug = debugEnabled;

    CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);

    std::atomic_flag flag = ATOMIC_FLAG_INIT;
    flag.test_and_set();

    char currentExePath[MAX_PATH];
    GetModuleFileNameA(NULL, currentExePath, MAX_PATH);
    char *currentExeName = PathFindFileNameA(currentExePath);

    std::wstring_convert<std::codecvt_utf8_utf16<wchar_t>> wideCharConverter;
    std::wstring userDataFolder =
        wideCharConverter.from_bytes(std::getenv("APPDATA"));
    std::wstring currentExeNameW = wideCharConverter.from_bytes(currentExeName);

    ICoreWebView2Controller *controller;
    ICoreWebView2* webview;

    HRESULT res = CreateCoreWebView2EnvironmentWithOptions(
            nullptr, (userDataFolder + L"/" + currentExeNameW).c_str(), nullptr,
            new wv2ComHandler(app, app->window, cb,
                                     [&](ICoreWebView2Controller *webviewController) {
                                         controller = webviewController;
                                         controller->get_CoreWebView2(&webview);
                                         webview->AddRef();
                                         ICoreWebView2Settings* settings;
                                         webview->get_Settings(&settings);
                                         if ( debugEnabled == 0 ) {
                                            settings->put_AreDefaultContextMenusEnabled(FALSE);
                                         }
                                         // Fix for invisible webview
                                         if( app->startHidden ) {}
                                         flag.clear();
                                     }));
    if (!SUCCEEDED(res))
    {
        switch (res)
        {
            case HRESULT_FROM_WIN32(ERROR_FILE_NOT_FOUND):
            {
                MessageBox(
                    app->window,
                    L"Couldn't find Edge installation. "
                    "Do you have a version installed that's compatible with this "
                    "WebView2 SDK version?",
                    nullptr, MB_OK);
            }
            break;
            case HRESULT_FROM_WIN32(ERROR_FILE_EXISTS):
            {
                MessageBox(
                    app->window, L"User data folder cannot be created because a file with the same name already exists.", nullptr, MB_OK);
            }
            break;
            case E_ACCESSDENIED:
            {
                MessageBox(
                    app->window, L"Unable to create user data folder, Access Denied.", nullptr, MB_OK);
            }
            break;
            case E_FAIL:
            {
                MessageBox(
                    app->window, L"Edge runtime unable to start", nullptr, MB_OK);
            }
            break;
            default:
            {
                 MessageBox(app->window, L"Failed to create WebView2 environment", nullptr, MB_OK);
            }
        }
    }

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

    // Let the backend know we have initialised
    app->webview->AddScriptToExecuteOnDocumentCreated(L"window.chrome.webview.postMessage('initialised');", nullptr);
    // Load the HTML
    LPCWSTR html = (LPCWSTR) cstrToLPWSTR((char*)assets[0]);
    app->webview->Navigate(html);

	messageFromWindowCallback("Ej{\"name\":\"wails:launched\",\"data\":[]}");
    return true;
}

void initialCallback(std::string message) {
    printf("MESSAGE=%s\n", message);
}

void Run(struct Application* app, int argc, char **argv) {

    // Register the window class.
    const wchar_t CLASS_NAME[]  = L"Ffenestri";

    WNDCLASSEX wc = { };

    wc.cbSize = sizeof(WNDCLASSEX);
    wc.lpfnWndProc   = WndProc;
    wc.hInstance     = GetModuleHandle(NULL);
    wc.lpszClassName = CLASS_NAME;

    if( app->disableWindowIcon == false ) {
        wc.hIcon = LoadIcon(wc.hInstance, MAKEINTRESOURCE(100));
        wc.hIconSm = LoadIcon(wc.hInstance, MAKEINTRESOURCE(100));
    }

    // Configure translucency
    DWORD dwExStyle = 0;
    if ( app->windowBackgroundIsTranslucent) {
        dwExStyle = WS_EX_NOREDIRECTIONBITMAP;
        wc.hbrBackground = CreateSolidBrush(RGB(255,255,255));
    }

    RegisterClassEx(&wc);

    // Process window style
    DWORD windowStyle = WS_OVERLAPPEDWINDOW | WS_THICKFRAME | WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX | WS_MAXIMIZEBOX;

    if (app->resizable == 0) {
        windowStyle &= ~WS_MAXIMIZEBOX;
        windowStyle &= ~WS_THICKFRAME;
    }
    if ( app->frame == 0 ) {
        windowStyle &= ~WS_OVERLAPPEDWINDOW;
        windowStyle &= ~WS_CAPTION;
        windowStyle |= WS_POPUP;
    }

    // Create the window.
    app->window = CreateWindowEx(
        dwExStyle,      // Optional window styles.
        CLASS_NAME,     // Window class
        L"",            // Window text
        windowStyle,    // Window style

        // Size and position
        CW_USEDEFAULT, CW_USEDEFAULT, app->width, app->height,

        NULL,       // Parent window
        NULL,       // Menu
        wc.hInstance,  // Instance handle
        NULL        // Additional application data
    );

    if (app->window == NULL)
    {
        return;
    }

    if ( app->fullscreen ) {
        fullscreen(app);
    }

    // Credit: https://stackoverflow.com/a/35482689
    if( app->disableWindowIcon && app->frame == 1 ) {
        int extendedStyle = GetWindowLong(app->window, GWL_EXSTYLE);
        SetWindowLong(app->window, GWL_EXSTYLE, extendedStyle | WS_EX_DLGMODALFRAME);
        SetWindowPos(nullptr, nullptr, 0, 0, 0, 0, SWP_FRAMECHANGED | SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER);
    }

    if ( app->windowBackgroundIsTranslucent ) {

        // Enable the translucent background effect
        enableTranslucentBackground(app);

        // Setup transparency of main window. This allows the blur to show through.
        SetLayeredWindowAttributes(app->window,RGB(255,255,255),0,LWA_COLORKEY);
    }

    // Store application pointer in window handle
    SetWindowLongPtr(app->window, GWLP_USERDATA, (LONG_PTR)app);

    // private center() as we are on main thread
    center(app);

    // Add webview2
    initWebView2(app, debug, initialCallback);

    if( app->webviewIsTranparent ) {
        wchar_t szBuff[64];
        ICoreWebView2Controller2 *wc2;
        wc2 = nullptr;
        app->webviewController->QueryInterface(IID_ICoreWebView2Controller2, (void**)&wc2);

        COREWEBVIEW2_COLOR wvColor;
        wvColor.R = app->backgroundColour.R;
        wvColor.G = app->backgroundColour.G;
        wvColor.B = app->backgroundColour.B;
        wvColor.A = app->backgroundColour.A == 0 ? 0 : 255;
        if( app->windowBackgroundIsTranslucent ) {
            wvColor.A = 0;
        }
        HRESULT result = wc2->put_DefaultBackgroundColor(wvColor);
        if (!SUCCEEDED(result))
        {
            switch (result)
            {
                case HRESULT_FROM_WIN32(ERROR_FILE_NOT_FOUND):
                {
                    MessageBox(
                        app->window,
                        L"Couldn't find Edge installation. "
                        "Do you have a version installed that's compatible with this "
                        "WebView2 SDK version?",
                        nullptr, MB_OK);
                }
                break;
                case HRESULT_FROM_WIN32(ERROR_FILE_EXISTS):
                {
                    MessageBox(
                        app->window, L"User data folder cannot be created because a file with the same name already exists.", nullptr, MB_OK);
                }
                break;
                case E_ACCESSDENIED:
                {
                    MessageBox(
                        app->window, L"Unable to create user data folder, Access Denied.", nullptr, MB_OK);
                }
                break;
                case E_FAIL:
                {
                    MessageBox(
                        app->window, L"Edge runtime unable to start", nullptr, MB_OK);
                }
                break;
                default:
                {
                     MessageBox(app->window, L"Failed to create WebView2 environment", nullptr, MB_OK);
                }
            }
        }

    }

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
        performShutdown(app);
        return;
      }
    }
}

void SetDebug(struct Application* app, int flag) {
    debug = flag;
}

void ExecJS(struct Application* app, const char *script) {
    ON_MAIN_THREAD(
        execJS(app, script);
    );
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

void DisableWindowIcon(struct Application* app) {
    app->disableWindowIcon = true;
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
    app->backgroundColour.R = red;
    app->backgroundColour.G = green;
    app->backgroundColour.B = blue;
    app->backgroundColour.A = alpha;
}

void SetSize(struct Application* app, int width, int height) {
    if( app->maxWidth > 0 && width > app->maxWidth ) {
        width = app->maxWidth;
    }
    if ( app->maxHeight > 0 && height > app->maxHeight ) {
        height = app->maxHeight;
    }
    SetWindowPos(app->window, nullptr, 0, 0, width, height, SWP_NOMOVE);
}

void setPosition(struct Application* app, int x, int y) {
    HMONITOR currentMonitor = MonitorFromWindow(app->window, MONITOR_DEFAULTTONEAREST);
    MONITORINFO info = {0};
    info.cbSize = sizeof(info);
    GetMonitorInfoA(currentMonitor, &info);
    RECT workRect = info.rcWork;
    LONG newX = workRect.left + x;
    LONG newY = workRect.top + y;

    SetWindowPos(app->window, HWND_TOP, newX, newY, 0, 0, SWP_NOSIZE);
}

void SetPosition(struct Application* app, int x, int y) {
    ON_MAIN_THREAD(
        setPosition(app, x, y);
    );
}

void Quit(struct Application* app) {
    // Override the hide window on close flag
    app->hideWindowOnClose = 0;
    ON_MAIN_THREAD(
        DestroyWindow(app->window);
    );
}


// Credit: https://stackoverflow.com/a/6693107
void setTitle(struct Application* app, const char *title) {
    LPCTSTR text = cstrToLPWSTR(title);
    SetWindowText(app->window, text);
    delete[] text;
}

void SetTitle(struct Application* app, const char *title) {
    ON_MAIN_THREAD(
        setTitle(app, title);
    );
}

void fullscreen(struct Application* app) {

    // Ensure we aren't in fullscreen
    if (app->isFullscreen) return;

    app->isFullscreen = true;
    app->previousWindowStyle = GetWindowLong(app->window, GWL_STYLE);
    MONITORINFO mi = { sizeof(mi) };
    if (GetWindowPlacement(app->window, &(app->previousPlacement)) && GetMonitorInfo(MonitorFromWindow(app->window, MONITOR_DEFAULTTOPRIMARY), &mi)) {
        SetWindowLong(app->window, GWL_STYLE, app->previousWindowStyle & ~WS_OVERLAPPEDWINDOW);
        SetWindowPos(app->window, HWND_TOP,
            mi.rcMonitor.left,
            mi.rcMonitor.top,
            mi.rcMonitor.right - mi.rcMonitor.left,
            mi.rcMonitor.bottom - mi.rcMonitor.top,
            SWP_NOOWNERZORDER | SWP_FRAMECHANGED);
    }
}

void Fullscreen(struct Application* app) {
    ON_MAIN_THREAD(
        fullscreen(app);
        show(app);
    );
}

void unfullscreen(struct Application* app) {
    if (app->isFullscreen) {
        SetWindowLong(app->window, GWL_STYLE, app->previousWindowStyle);
        SetWindowPlacement(app->window, &(app->previousPlacement));
        SetWindowPos(app->window, NULL, 0, 0, 0, 0,
                         SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER |
                         SWP_NOOWNERZORDER | SWP_FRAMECHANGED);
        app->isFullscreen = false;
    }
}

void UnFullscreen(struct Application* app) {
    ON_MAIN_THREAD(
        unfullscreen(app);
    );
}

void DisableFrame(struct Application* app) {
    app->frame = 0;
}

// WebviewIsTransparent will make the webview transparent
// revealing the window underneath
void WebviewIsTransparent(struct Application *app) {
	app->webviewIsTranparent = true;
}

void WindowBackgroundIsTranslucent(struct Application *app) {
	app->windowBackgroundIsTranslucent = true;
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