// Some code may be inspired by or directly used from Webview.
#include "ffenestri_windows.h"

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
    int minWidth;
    int minHeight;
    int maxWidth;
    int maxHeight;
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

      return result;
}

void SetMinWindowSize(struct Application* app, int minWidth, int minHeight) {
    app->minWidth = minWidth;
    app->minHeight = minHeight;
}

void SetMaxWindowSize(struct Application* app, int maxWidth, int maxHeight) {
    app->maxWidth = maxWidth;
    app->maxHeight = maxHeight;
}

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {

    switch(msg) {

      case WM_DESTROY:

          PostQuitMessage(0);
          break;
    }

    return DefWindowProcW(hwnd, msg, wParam, lParam);
}

void Run(struct Application* app, int argc, char **argv) {

      WNDCLASSEX wc;
      HINSTANCE hInstance = GetModuleHandle(NULL);
      ZeroMemory(&wc, sizeof(WNDCLASSEX));
      wc.cbSize = sizeof(WNDCLASSEX);
      wc.hInstance = hInstance;
      wc.lpszClassName = "ffenestri";
      wc.lpfnWndProc   = WndProc;

      // TODO: Trim title to 256 chars
      // https://stackoverflow.com/a/20458904
      wchar_t wchTitle[256];
      MultiByteToWideChar(CP_ACP, 0, app->title, -1, wchTitle, 256);

      RegisterClassEx(&wc);
      app->window = CreateWindow("ffenestri", wchTitle, WS_OVERLAPPEDWINDOW, CW_USEDEFAULT,
                                          CW_USEDEFAULT, app->width, app->height, NULL, NULL,
                                          GetModuleHandle(NULL), NULL);

    MSG  msg;
    ShowWindow(app->window, SW_SHOWNORMAL);
    UpdateWindow(app->window);
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
}
void SetDebug(struct Application* app, int flag) {
}
void SetBindings(struct Application* app, const char *bindings) {
}
void ExecJS(struct Application* app, const char *script) {
}
void Hide(struct Application* app) {
}
void Show(struct Application* app) {
}
void Center(struct Application* app) {
}
void Maximise(struct Application* app) {
}
void Unmaximise(struct Application* app) {
}
void ToggleMaximise(struct Application* app) {
}
void Minimise(struct Application* app) {
}
void Unminimise(struct Application* app) {
}
void ToggleMinimise(struct Application* app) {
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