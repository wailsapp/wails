#ifndef window_h
#define window_h

#include <JavaScriptCore/JavaScript.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

typedef struct DragOptions
{
    void *webview;
    GtkWindow *mainwindow;
} DragOptions;

typedef struct ResizeOptions
{
    void *webview;
    GtkWindow *mainwindow;
    GdkWindowEdge edge;
} ResizeOptions;

typedef struct JSCallback
{
    void *webview;
    char *script;
} JSCallback;

typedef struct MessageDialogOptions
{
    void *window;
    char *title;
    char *message;
    int messageType;
} MessageDialogOptions;

typedef struct OpenFileDialogOptions
{
    GtkWindow *window;
    char *title;
    char *defaultFilename;
    char *defaultDirectory;
    int createDirectories;
    int multipleFiles;
    int showHiddenFiles;
    GtkFileChooserAction action;
    GtkFileFilter **filters;
} OpenFileDialogOptions;

typedef struct RGBAOptions
{
    uint8_t r;
    uint8_t g;
    uint8_t b;
    uint8_t a;
    void *webview;
    void *webviewBox;
    gboolean windowIsTranslucent;
} RGBAOptions;

typedef struct SetTitleArgs
{
    GtkWindow *window;
    char *title;
} SetTitleArgs;

typedef struct SetPositionArgs
{
    int x;
    int y;
    void *window;
} SetPositionArgs;

void ExecuteOnMainThread(void *f, gpointer jscallback);

GtkWidget *GTKWIDGET(void *pointer);
GtkWindow *GTKWINDOW(void *pointer);
GtkContainer *GTKCONTAINER(void *pointer);
GtkBox *GTKBOX(void *pointer);

// window
ulong SetupInvokeSignal(void *contentManager);

void SetWindowIcon(GtkWindow *window, const guchar *buf, gsize len);
void SetWindowTransparency(GtkWidget *widget);
void SetBackgroundColour(void *data);
void SetTitle(GtkWindow *window, char *title);
void SetPosition(void *window, int x, int y);
void SetMinMaxSize(GtkWindow *window, int min_width, int min_height, int max_width, int max_height);
void DisableContextMenu(void *webview);
void ConnectButtons(void *webview);

int IsFullscreen(GtkWidget *widget);
int IsMaximised(GtkWidget *widget);
int IsMinimised(GtkWidget *widget);

gboolean Center(gpointer data);
gboolean Show(gpointer data);
gboolean Hide(gpointer data);
gboolean Maximise(gpointer data);
gboolean UnMaximise(gpointer data);
gboolean Minimise(gpointer data);
gboolean UnMinimise(gpointer data);
gboolean Fullscreen(gpointer data);
gboolean UnFullscreen(gpointer data);

// WebView
GtkWidget *SetupWebview(void *contentManager, GtkWindow *window, int hideWindowOnClose, int gpuPolicy);
void LoadIndex(void *webview, char *url);
void DevtoolsEnabled(void *webview, int enabled, bool showInspector);
void ExecuteJS(void *data);

// Drag
void StartDrag(void *webview, GtkWindow *mainwindow);
void StartResize(void *webview, GtkWindow *mainwindow, GdkWindowEdge edge);

// Dialog
void MessageDialog(void *data);
GtkFileFilter **AllocFileFilterArray(size_t ln);
void Opendialog(void *data);

// Inspector
void sendShowInspectorMessage();
void ShowInspector(void *webview);
void InstallF12Hotkey(void *window);

#endif /* window_h */
