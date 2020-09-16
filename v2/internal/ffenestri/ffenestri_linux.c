
#ifndef __FFENESTRI_LINUX_H__
#define __FFENESTRI_LINUX_H__

#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
#include <time.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <stdarg.h>

// References to assets
extern const unsigned char *assets[];
extern const unsigned char runtime;
extern const char *icon[];

// Constants
#define PRIMARY_MOUSE_BUTTON 1
#define MIDDLE_MOUSE_BUTTON 2
#define SECONDARY_MOUSE_BUTTON 3

// MAIN DEBUG FLAG
int debug;

// Credit: https://stackoverflow.com/a/8465083
char *concat(const char *s1, const char *s2)
{
    const size_t len1 = strlen(s1);
    const size_t len2 = strlen(s2);
    char *result = malloc(len1 + len2 + 1);
    memcpy(result, s1, len1);
    memcpy(result + len1, s2, len2 + 1);
    return result;
}

// Debug works like sprintf but mutes if the global debug flag is true
// Credit: https://stackoverflow.com/a/20639708
void Debug(char *message, ...)
{
    if (debug)
    {
        char *temp = concat("TRACE | Ffenestri (C) | ", message);
        message = concat(temp, "\n");
        free(temp);
        va_list args;
        va_start(args, message);
        vprintf(message, args);
        free(message);
        va_end(args);
    }
}

extern void messageFromWindowCallback(const char *);
typedef void (*ffenestriCallback)(const char *);

struct Application
{

    // Gtk Data
    GtkApplication *application;
    GtkWindow *mainWindow;
    GtkWidget *webView;
    int signalInvoke;
    int signalWindowDrag;
    int signalButtonPressed;
    int signalButtonReleased;
    int signalLoadChanged;

    // Saves the events for the drag mouse button
    GdkEventButton *dragButtonEvent;

    // The number of the default drag button
    int dragButton;

    // Window Data
    const char *title;
    char *id;
    int width;
    int height;
    int resizable;
    int devtools;
    int startHidden;
    int fullscreen;
    int minWidth;
    int minHeight;
    int maxWidth;
    int maxHeight;
    int frame;

    // User Data
    char *HTML;

    // Callback
    ffenestriCallback sendMessageToBackend;

    // Bindings
    const char *bindings;

    // Lock - used for sync operations (Should we be using g_mutex?)
    int lock;
};

void *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden)
{
    // Setup main application struct
    struct Application *result = malloc(sizeof(struct Application));
    result->title = title;
    result->width = width;
    result->height = height;
    result->resizable = resizable;
    result->devtools = devtools;
    result->fullscreen = fullscreen;
    result->minWidth = 0;
    result->minHeight = 0;
    result->maxWidth = 0;
    result->maxHeight = 0;
    result->frame = 1;
    result->startHidden = startHidden;

    // Default drag button is PRIMARY
    result->dragButton = PRIMARY_MOUSE_BUTTON;

    result->sendMessageToBackend = (ffenestriCallback)messageFromWindowCallback;

    // Create a unique ID based on the current unix timestamp
    char temp[11];
    sprintf(&temp[0], "%d", (int)time(NULL));
    result->id = concat("wails.app", &temp[0]);

    // Create the main GTK application
    GApplicationFlags flags = G_APPLICATION_FLAGS_NONE;
    result->application = gtk_application_new(result->id, flags);

    // Return the application struct
    return (void *)result;
}

void DestroyApplication(struct Application *app)
{
    Debug("Destroying Application");

    g_application_quit(G_APPLICATION(app->application));

    // Release the GTK ID string
    if (app->id != NULL)
    {
        free(app->id);
        app->id = NULL;
    }
    else
    {
        Debug("Almost a double free for app->id");
    }

    // Free the bindings
    if (app->bindings != NULL)
    {
        free((void *)app->bindings);
        app->bindings = NULL;
    }
    else
    {
        Debug("Almost a double free for app->bindings");
    }

    // Disconnect signal handlers
    WebKitUserContentManager *manager = webkit_web_view_get_user_content_manager((WebKitWebView *)app->webView);
    g_signal_handler_disconnect(manager, app->signalInvoke);
    if( app->frame == 0) {
        g_signal_handler_disconnect(manager, app->signalWindowDrag);
        g_signal_handler_disconnect(app->webView, app->signalButtonPressed);
        g_signal_handler_disconnect(app->webView, app->signalButtonReleased);
    }
    g_signal_handler_disconnect(app->webView, app->signalLoadChanged);

    // Release the main GTK Application
    if (app->application != NULL)
    {
        g_object_unref(app->application);
        app->application = NULL;
    }
    else
    {
        Debug("Almost a double free for app->application");
    }
    Debug("Finished Destroying Application");
}

// Quit will stop the gtk application and free up all the memory
// used by the application
void Quit(struct Application *app)
{
    Debug("Quit Called");
    gtk_window_close((GtkWindow *)app->mainWindow);
    g_application_quit((GApplication *)app->application);
    DestroyApplication(app);
}

// SetTitle sets the main window title to the given string
void SetTitle(struct Application *app, const char *title)
{
    gtk_window_set_title(app->mainWindow, title);
}

// Fullscreen sets the main window to be fullscreen
void Fullscreen(struct Application *app)
{
    gtk_window_fullscreen(app->mainWindow);
}

// UnFullscreen resets the main window after a fullscreen
void UnFullscreen(struct Application *app)
{
    gtk_window_unfullscreen(app->mainWindow);
}

void setMinMaxSize(struct Application *app)
{
    GdkGeometry size;
    size.min_width = size.min_height = size.max_width = size.max_height = 0;
    int flags = 0;
    if (app->maxHeight > 0 && app->maxWidth > 0)
    {
        size.max_height = app->maxHeight;
        size.max_width = app->maxWidth;
        flags |= GDK_HINT_MAX_SIZE;
    }
    if (app->minHeight > 0 && app->minWidth > 0)
    {
        size.min_height = app->minHeight;
        size.min_width = app->minWidth;
        flags |= GDK_HINT_MIN_SIZE;
    }
    gtk_window_set_geometry_hints(app->mainWindow, NULL, &size, flags);
}

char *fileDialogInternal(struct Application *app, GtkFileChooserAction chooserAction, char **args) {
    GtkFileChooserNative *native;
    GtkFileChooserAction action = chooserAction;
    gint res;
    char *filename;

    char *title = args[0];
    char *filter = args[1];

    native = gtk_file_chooser_native_new(title,
                                         app->mainWindow,
                                         action,
                                         "_Open",
                                         "_Cancel");

    GtkFileChooser *chooser = GTK_FILE_CHOOSER(native);

    // If we have filters, process them
    if (filter[0] != '\0') {
        GtkFileFilter *file_filter = gtk_file_filter_new();
        gchar **filters  = g_strsplit(filter, ",", -1);
        gint i;
        for(i = 0; filters && filters[i]; i++) {
            gtk_file_filter_add_pattern(file_filter, filters[i]);
            // Debug("Adding filter pattern: %s\n", filters[i]);
        }
        gtk_file_filter_set_name(file_filter, filter);
        gtk_file_chooser_add_filter(chooser, file_filter);
        g_strfreev(filters);
    }

    res = gtk_native_dialog_run(GTK_NATIVE_DIALOG(native));
    if (res == GTK_RESPONSE_ACCEPT)
    {
        filename = gtk_file_chooser_get_filename(chooser);
    }

    g_object_unref(native);

    return filename;
}

// openFileDialogInternal opens a dialog to select a file
// NOTE: The result is a string that will need to be freed!
char *openFileDialogInternal(struct Application *app, char **args)
{
    return fileDialogInternal(app, GTK_FILE_CHOOSER_ACTION_OPEN, args);
}

// saveFileDialogInternal opens a dialog to select a file
// NOTE: The result is a string that will need to be freed!
char *saveFileDialogInternal(struct Application *app, char **args)
{
    return fileDialogInternal(app, GTK_FILE_CHOOSER_ACTION_SAVE, args);
}


// openDirectoryDialogInternal opens a dialog to select a directory
// NOTE: The result is a string that will need to be freed!
char *openDirectoryDialogInternal(struct Application *app, char **args)
{
    return fileDialogInternal(app, GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER, args);
}

void SetMinWindowSize(struct Application *app, int minWidth, int minHeight)
{
    app->minWidth = minWidth;
    app->minHeight = minHeight;
}

void SetMaxWindowSize(struct Application *app, int maxWidth, int maxHeight)
{
    app->maxWidth = maxWidth;
    app->maxHeight = maxHeight;
}

// SetColour sets the colour of the webview to the given colour string
int SetColour(struct Application *app, const char *colourString)
{
    GdkRGBA rgba;
    gboolean result = gdk_rgba_parse(&rgba, colourString);
    if (result == FALSE)
    {
        return 0;
    }
    // Debug("Setting webview colour to: %s", colourString);
    webkit_web_view_get_background_color((WebKitWebView *)(app->webView), &rgba);
    return 1;
}

// DisableFrame disables the window frame
void DisableFrame(struct Application *app)
{
    app->frame = 0;
}

void syncCallback(GObject *source_object,
                  GAsyncResult *res,
                  void *data)
{
    struct Application *app = (struct Application *)data;
    app->lock = 0;
}

void syncEval(struct Application *app, const gchar *script)
{

    WebKitWebView *webView = (WebKitWebView *)(app->webView);

    // wait for lock to free
    while (app->lock == 1)
    {
        g_main_context_iteration(0, true);
    }
    // Set lock
    app->lock = 1;

    webkit_web_view_run_javascript(
        webView,
        script,
        NULL, syncCallback, (void*)app);

    while (app->lock == 1)
    {
        g_main_context_iteration(0, true);
    }
}

void asyncEval(WebKitWebView *webView, const gchar *script)
{
    webkit_web_view_run_javascript(
        webView,
        script,
        NULL, NULL, NULL);
}

typedef void (*dispatchMethod)(struct Application *app, void *);

struct dispatchData
{
    struct Application *app;
    dispatchMethod method;
    void *args;
};

gboolean executeMethod(gpointer data)
{
    struct dispatchData *d = (struct dispatchData *)data;
    (d->method)(d->app, d->args);
    g_free(d);
    return FALSE;
}

void ExecJS(struct Application *app, char *js)
{
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)syncEval;
    data->args = js;
    data->app = app;

    gdk_threads_add_idle(executeMethod, data);
}

typedef char *(*dialogMethod)(struct Application *app, void *);

struct dialogCall
{
    struct Application *app;
    dialogMethod method;
    void *args;
    void *filter;
    char *result;
    int done;
};

gboolean executeMethodWithReturn(gpointer data)
{
    struct dialogCall *d = (struct dialogCall *)data;

    d->result = (d->method)(d->app, d->args);
    d->done = 1;
    return FALSE;
}

char *OpenFileDialog(struct Application *app, char *title, char *filter)
{
    struct dialogCall *data = (struct dialogCall *)g_new(struct dialogCall, 1);
    data->result = NULL;
    data->done = 0;
    data->method = (dialogMethod)openFileDialogInternal;
    const char* dialogArgs[]={ title, filter };
    data->args = dialogArgs;
    data->app = app;

    gdk_threads_add_idle(executeMethodWithReturn, data);

    while (data->done == 0)
    {
        usleep(100000);
    }
    g_free(data);
    return data->result;
}

char *SaveFileDialog(struct Application *app, char *title, char *filter)
{
    struct dialogCall *data = (struct dialogCall *)g_new(struct dialogCall, 1);
    data->result = NULL;
    data->done = 0;
    data->method = (dialogMethod)saveFileDialogInternal;
    const char* dialogArgs[]={ title, filter };
    data->args = dialogArgs;
    data->app = app;

    gdk_threads_add_idle(executeMethodWithReturn, data);

    while (data->done == 0)
    {
        usleep(100000);
    }
    Debug("Dialog done");
    Debug("Result = %s\n", data->result);

    g_free(data);
    // Fingers crossed this wasn't freed by g_free above
    return data->result;
}

char *OpenDirectoryDialog(struct Application *app, char *title, char *filter)
{
    struct dialogCall *data = (struct dialogCall *)g_new(struct dialogCall, 1);
    data->result = NULL;
    data->done = 0;
    data->method = (dialogMethod)openDirectoryDialogInternal;
    const char* dialogArgs[]={ title, filter };
    data->args = dialogArgs;
    data->app = app;

    gdk_threads_add_idle(executeMethodWithReturn, data);

    while (data->done == 0)
    {
        usleep(100000);
    }
    Debug("Directory Dialog done");
    Debug("Result = %s\n", data->result);
    g_free(data);
    // Fingers crossed this wasn't freed by g_free above
    return data->result;
}

// Sets the icon to the XPM stored in icon
void setIcon(struct Application *app)
{
    GdkPixbuf *appIcon = gdk_pixbuf_new_from_xpm_data((const char **)icon);
    gtk_window_set_icon(app->mainWindow, appIcon);
}

static void load_finished_cb(WebKitWebView *webView,
                             WebKitLoadEvent load_event,
                             struct Application *app)
{
    switch (load_event)
    {
    // case WEBKIT_LOAD_STARTED:
    //     /* New load, we have now a provisional URI */
    //     // printf("Start downloading %s\n", webkit_web_view_get_uri(web_view));
    //     /* Here we could start a spinner or update the
    //      * location bar with the provisional URI */
    //     break;
    // case WEBKIT_LOAD_REDIRECTED:
    //     // printf("Redirected to: %s\n", webkit_web_view_get_uri(web_view));
    //     break;
    // case WEBKIT_LOAD_COMMITTED:
    //     /* The load is being performed. Current URI is
    //      * the final one and it won't change unless a new
    //      * load is requested or a navigation within the
    //      * same page is performed */
    //     // printf("Loading: %s\n", webkit_web_view_get_uri(web_view));
    //     break;
    case WEBKIT_LOAD_FINISHED:
        /* Load finished, we can now stop the spinner */
        // printf("Finished loading: %s\n", webkit_web_view_get_uri(web_view));

        // Bindings
        Debug("Binding Methods");
        syncEval(app, app->bindings);

        // Runtime
        Debug("Setting up Wails runtime");
        syncEval(app, &runtime);

        // Loop over assets
        int index = 1;
        while (1)
        {
            // Get next asset pointer
            const char *asset = assets[index];

            // If we have no more assets, break
            if (asset == 0x00)
            {
                break;
            }

            // sync eval the asset
            syncEval(app, asset);
            index++;
        };

        // Set the icon
        setIcon(app);

        // Setup fullscreen
        if (app->fullscreen)
        {
            Debug("Going fullscreen");
            Fullscreen(app);
        }

        // Setup resize
        gtk_window_resize(GTK_WINDOW(app->mainWindow), app->width, app->height);

        if (app->resizable)
        {
            gtk_window_set_default_size(GTK_WINDOW(app->mainWindow), app->width, app->height);
        }
        else
        {
            gtk_widget_set_size_request(GTK_WIDGET(app->mainWindow), app->width, app->height);
            gtk_window_resize(GTK_WINDOW(app->mainWindow), app->width, app->height);
            // Fix the min/max to the window size for good measure
            app->minHeight = app->maxHeight = app->height;
            app->minWidth = app->maxWidth = app->width;
        }
        gtk_window_set_resizable(GTK_WINDOW(app->mainWindow), app->resizable ? TRUE : FALSE);
        setMinMaxSize(app);

        // Centre by default
        gtk_window_set_position(app->mainWindow, GTK_WIN_POS_CENTER);

        // Show window and focus
        if( app->startHidden == 0) {
            gtk_widget_show_all(GTK_WIDGET(app->mainWindow));
            gtk_widget_grab_focus(app->webView);
        }
        break;
    }
}

static gboolean disable_context_menu_cb(
    WebKitWebView *web_view,
    WebKitContextMenu *context_menu,
    GdkEvent *event,
    WebKitHitTestResult *hit_test_result,
    gpointer user_data)
{
    return TRUE;
}

static void printEvent(const char *message, GdkEventButton *event)
{
    Debug("%s: [button:%d] [x:%f] [y:%f] [time:%d]",
          message,
          event->button,
          event->x_root,
          event->y_root,
          event->time);
}


static void dragWindow(WebKitUserContentManager *contentManager,
                       WebKitJavascriptResult *result,
                       struct Application *app)
{
    // If we get this message erroneously, ignore
    if (app->dragButtonEvent == NULL)
    {
        return;
    }

    // Ignore non-toplevel widgets
    GtkWidget *window = gtk_widget_get_toplevel(GTK_WIDGET(app->webView));
    if (!GTK_IS_WINDOW(window))
    {
        return;
    }

    // Initiate the drag
    printEvent("Starting drag with event", app->dragButtonEvent);

    gtk_window_begin_move_drag(app->mainWindow,
                               app->dragButton,
                               app->dragButtonEvent->x_root,
                               app->dragButtonEvent->y_root,
                               app->dragButtonEvent->time);
    // Clear the event
    app->dragButtonEvent = NULL;

    return;
}

gboolean buttonPress(GtkWidget *widget, GdkEventButton *event, struct Application *app)
{
    if (event->type == GDK_BUTTON_PRESS && event->button == app->dragButton)
    {
        app->dragButtonEvent = event;
    }
    return FALSE;
}

gboolean buttonRelease(GtkWidget *widget, GdkEventButton *event, struct Application *app)
{
    if (event->type == GDK_BUTTON_RELEASE && event->button == app->dragButton)
    {
        app->dragButtonEvent = NULL;
    }
    return FALSE;
}

static void sendMessageToBackend(WebKitUserContentManager *contentManager,
                                 WebKitJavascriptResult *result,
                                 struct Application *app)
{
#if WEBKIT_MAJOR_VERSION >= 2 && WEBKIT_MINOR_VERSION >= 22
    JSCValue *value = webkit_javascript_result_get_js_value(result);
    char *message = jsc_value_to_string(value);
#else
    JSGlobalContextRef context = webkit_javascript_result_get_global_context(result);
    JSValueRef value = webkit_javascript_result_get_value(result);
    JSStringRef js = JSValueToStringCopy(context, value, NULL);
    size_t messageSize = JSStringGetMaximumUTF8CStringSize(js);
    char *message = g_new(char, messageSize);
    JSStringGetUTF8CString(js, message, messageSize);
    JSStringRelease(js);
#endif
    app->sendMessageToBackend(message);
    g_free(message);
}

void SetDebug(struct Application *app, int flag)
{
    debug = flag;
}

// getCurrentMonitorGeometry gets the geometry of the monitor
// that the window is mostly on.
GdkRectangle getCurrentMonitorGeometry(GtkWindow *window) {
    // Get the monitor that the window is currently on
    GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
    GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
    GdkMonitor *monitor = gdk_display_get_monitor_at_window (display, gdk_window);

    // Get the geometry of the monitor
    GdkRectangle result;
    gdk_monitor_get_geometry (monitor,&result);

    return result;
}

/*******************
 * Window Position *
 *******************/

// Position holds an x/y corrdinate
struct Position {
    int x;
    int y;
};

// Internal call for setting the position of the window.
void setPositionInternal(struct Application *app, struct Position *pos) {

    // Get the monitor geometry
    GdkRectangle m = getCurrentMonitorGeometry(app->mainWindow);

    // Move the window relative to the monitor
    gtk_window_move(app->mainWindow, m.x + pos->x, m.y + pos->y);
    
    // Free memory
    free(pos);
}

// SetPosition sets the position of the window to the given x/y 
// coordinates. The x/y values are relative to the monitor 
// the window is mostly on.
void SetPosition(struct Application *app, int x, int y) {
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)setPositionInternal;
    struct Position *pos = malloc(sizeof(struct Position));
    pos->x = x;
    pos->y = y;
    data->args = pos;
    data->app = app;

    gdk_threads_add_idle(executeMethod, data);  
}

/***************
 * Window Size *
 ***************/

// Size holds a width/height
struct Size {
    int width;
    int height;
};

// Internal call for setting the size of the window.
void setSizeInternal(struct Application *app, struct Size *size) {
    gtk_window_resize(app->mainWindow, size->width, size->height);
    free(size);
}

// SetSize sets the size of the window to the given width/height
void SetSize(struct Application *app, int width, int height) {
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)setSizeInternal;
    struct Size *size = malloc(sizeof(struct Size));
    size->width = width;
    size->height = height;
    data->args = size;
    data->app = app;

    gdk_threads_add_idle(executeMethod, data);  
}


// centerInternal will center the main window on the monitor it is mostly in
void centerInternal(struct Application *app)
{
    // Get the geometry of the monitor
    GdkRectangle m = getCurrentMonitorGeometry(app->mainWindow);

    // Get the window width/height
    int windowWidth, windowHeight;
    gtk_window_get_size(app->mainWindow, &windowWidth, &windowHeight);

    // Place the window at the center of the monitor
    gtk_window_move(app->mainWindow, ((m.width - windowWidth) / 2) + m.x, ((m.height - windowHeight) / 2) + m.y);
}

// Center the window
void Center(struct Application *app) {

    // Setup a call to centerInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)centerInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}

// hideInternal hides the main window
void hideInternal(struct Application *app) {
    gtk_widget_hide (GTK_WIDGET(app->mainWindow));
}

// Hide places the hideInternal method onto the main thread for execution
void Hide(struct Application *app) {

    // Setup a call to hideInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)hideInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}

// showInternal shows the main window 
void showInternal(struct Application *app) {
    gtk_widget_show_all(GTK_WIDGET(app->mainWindow));
    gtk_widget_grab_focus(app->webView);
}

// Show places the showInternal method onto the main thread for execution
void Show(struct Application *app) {
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)showInternal;
    data->app = app;

    gdk_threads_add_idle(executeMethod, data);  
}


// maximiseInternal maximises the main window
void maximiseInternal(struct Application *app) {
    gtk_window_maximize(GTK_WIDGET(app->mainWindow));
}

// Maximise places the maximiseInternal method onto the main thread for execution
void Maximise(struct Application *app) {

    // Setup a call to maximiseInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)maximiseInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}

// unmaximiseInternal unmaximises the main window
void unmaximiseInternal(struct Application *app) {
    gtk_window_unmaximize(GTK_WIDGET(app->mainWindow));
}

// Unmaximise places the unmaximiseInternal method onto the main thread for execution
void Unmaximise(struct Application *app) {

    // Setup a call to unmaximiseInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)unmaximiseInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}


// minimiseInternal minimises the main window
void minimiseInternal(struct Application *app) {
    gtk_window_iconify(app->mainWindow);
}

// Minimise places the minimiseInternal method onto the main thread for execution
void Minimise(struct Application *app) {

    // Setup a call to minimiseInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)minimiseInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}

// unminimiseInternal unminimises the main window
void unminimiseInternal(struct Application *app) {
    gtk_window_present(app->mainWindow);
}

// Unminimise places the unminimiseInternal method onto the main thread for execution
void Unminimise(struct Application *app) {

    // Setup a call to unminimiseInternal on the main thread
    struct dispatchData *data = (struct dispatchData *)g_new(struct dispatchData, 1);
    data->method = (dispatchMethod)unminimiseInternal;
    data->app = app;

    // Add call to main thread
    gdk_threads_add_idle(executeMethod, data);  
}


void SetBindings(struct Application *app, const char *bindings)
{
    const char *temp = concat("window.wailsbindings = \"", bindings);
    const char *jscall = concat(temp, "\";");
    free((void *)temp);
    app->bindings = jscall;
}

// This is called when the close button on the window is pressed
gboolean close_button_pressed(GtkWidget *widget,
                              GdkEvent *event,
                              struct Application *app)
{
    app->sendMessageToBackend("WC"); // Window Close
    return TRUE;
}

static void setupWindow(struct Application *app)
{

    // Create the window
    GtkWidget *mainWindow = gtk_application_window_new(app->application);
    // Save reference
    app->mainWindow = GTK_WINDOW(mainWindow);

    // Setup frame
    gtk_window_set_decorated((GtkWindow *)mainWindow, app->frame);

    // Setup title
    gtk_window_set_title(GTK_WINDOW(mainWindow), app->title);

    // Setup script handler
    WebKitUserContentManager *contentManager = webkit_user_content_manager_new();

    // Setup the invoke handler
    webkit_user_content_manager_register_script_message_handler(contentManager, "external");
    app->signalInvoke = g_signal_connect(contentManager, "script-message-received::external", G_CALLBACK(sendMessageToBackend), app);

    // Setup the window drag handler if this is a frameless app
    if ( app->frame == 0 ) {
        webkit_user_content_manager_register_script_message_handler(contentManager, "windowDrag");
        app->signalWindowDrag = g_signal_connect(contentManager, "script-message-received::windowDrag", G_CALLBACK(dragWindow), app);
        // Setup the mouse handlers
        app->signalButtonPressed = g_signal_connect(app->webView, "button-press-event", G_CALLBACK(buttonPress), app);
        app->signalButtonReleased = g_signal_connect(app->webView, "button-release-event", G_CALLBACK(buttonRelease), app);
    }
    GtkWidget *webView = webkit_web_view_new_with_user_content_manager(contentManager);

    // Save reference
    app->webView = webView;

    // Add the webview to the window
    gtk_container_add(GTK_CONTAINER(mainWindow), webView);


    // Load default HTML
    app->signalLoadChanged = g_signal_connect(G_OBJECT(webView), "load-changed", G_CALLBACK(load_finished_cb), app);

    // Load the user's HTML
    // assets[0] is the HTML because the asset array is bundled like that by convention
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webView), assets[0]);

    // Check if we want to enable the dev tools
    if (app->devtools)
    {
        WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webView));
        // webkit_settings_set_enable_write_console_messages_to_stdout(settings, true);
        webkit_settings_set_enable_developer_extras(settings, true);
    }
    else
    {
        g_signal_connect(G_OBJECT(webView), "context-menu", G_CALLBACK(disable_context_menu_cb), app);
    }

    // Listen for close button signal
    g_signal_connect(GTK_WIDGET(mainWindow), "delete-event", G_CALLBACK(close_button_pressed), app);
}

static void activate(GtkApplication* _, struct Application *app)
{
    setupWindow(app);
}

void Run(struct Application *app, int argc, char **argv)
{
    g_signal_connect(app->application, "activate", G_CALLBACK(activate), app);
    g_application_run(G_APPLICATION(app->application), argc, argv);
}

#endif
