//go:build linux && webkit_6
// +build linux,webkit_6

#include <jsc/jsc.h>
#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#include <string.h>
#include <locale.h>
#include "window_webkit6.h"

// These are the x,y,time & button of the last mouse down event
// It's used for window dragging
static float xroot = 0.0f;
static float yroot = 0.0f;
static int dragTime = -1;
static uint mouseButton = 0;
static int wmIsWayland = -1;
static int decoratorWidth = -1;
static int decoratorHeight = -1;

// casts
void ExecuteOnMainThread(void *f, gpointer jscallback)
{
    g_idle_add((GSourceFunc)f, (gpointer)jscallback);
}

GtkWidget *GTKWIDGET(void *pointer)
{
    return GTK_WIDGET(pointer);
}

GtkWindow *GTKWINDOW(void *pointer)
{
    return GTK_WINDOW(pointer);
}

GtkBox *GTKBOX(void *pointer)
{
    return GTK_BOX(pointer);
}

extern void processMessage(char *);

static void sendMessageToBackend(WebKitUserContentManager *contentManager,
                                 JSCValue *value,
                                 void *data)
{
    char *message = jsc_value_to_string(value);

    processMessage(message);
    g_free(message);
}

static bool isNULLRectangle(GdkRectangle input)
{
    return input.x == -1 && input.y == -1 && input.width == -1 && input.height == -1;
}

static gboolean onWayland()
{
    switch (wmIsWayland)
    {
    case -1:
     char *gdkBackend = getenv("XDG_SESSION_TYPE");
        if(gdkBackend != NULL && strcmp(gdkBackend, "wayland") == 0) 
        {
            wmIsWayland = 1;
            return TRUE;
        }
        
        wmIsWayland = 0;
        return FALSE;
    case 1:
        return TRUE;
    default:
        return FALSE;
    }
}

static GdkMonitor *getCurrentMonitor(GtkWindow *window)
{
    // Get the monitor that the window is currently on
    GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));

    if(native == NULL) {
        return NULL;
    }

	GdkSurface *surface = gtk_native_get_surface(native);

	GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));

	GdkMonitor *currentMonitor = gdk_display_get_monitor_at_surface(display, surface);

    return currentMonitor;
}

static GdkRectangle getCurrentMonitorGeometry(GtkWindow *window)
{
    GdkMonitor *monitor = getCurrentMonitor(window);
    GdkRectangle result;
    if (monitor == NULL)
    {
        result.x = result.y = result.height = result.width = -1;
        return result;
    }

    // Get the geometry of the monitor
    gdk_monitor_get_geometry(monitor, &result);
    return result;
}

static int getCurrentMonitorScaleFactor(GtkWindow *window)
{
    GdkMonitor *monitor = getCurrentMonitor(window);

    return gdk_monitor_get_scale_factor(monitor);
}

// window

ulong SetupInvokeSignal(void *contentManager)
{
    return g_signal_connect((WebKitUserContentManager *)contentManager, "script-message-received::external", G_CALLBACK(sendMessageToBackend), NULL);
}

void SetWindowIcon(GtkWindow *window, const guchar *buf, gsize len)
{
    // GdkPixbufLoader *loader = gdk_pixbuf_loader_new();
    // if (!loader)
    // {
    //     return;
    // }
    // if (gdk_pixbuf_loader_write(loader, buf, len, NULL) && gdk_pixbuf_loader_close(loader, NULL))
    // {
    //     GdkPixbuf *pixbuf = gdk_pixbuf_loader_get_pixbuf(loader);
    //     if (pixbuf)
    //     {
    //         gtk_window_set_icon(window, pixbuf);
    //     }
    // }
    // g_object_unref(loader);
}

void SetWindowTransparency(GtkWidget *widget)
{
    //// TODO: gtk_widget_set_opacity might be able to be used here?

    // GdkScreen *screen = gtk_widget_get_screen(widget);
    // GdkVisual *visual = gdk_screen_get_rgba_visual(screen);

    // if (visual != NULL && gdk_screen_is_composited(screen))
    // {
    //     gtk_widget_set_app_paintable(widget, true);
    //     gtk_widget_set_visual(widget, visual);
    // }
}

static GtkCssProvider *windowCssProvider = NULL;

void SetBackgroundColour(void *data)
{
    // set webview's background color
    RGBAOptions *options = (RGBAOptions *)data;

    GdkRGBA colour = {options->r / 255.0, options->g / 255.0, options->b / 255.0, options->a / 255.0};
    if (options->windowIsTranslucent != NULL && options->windowIsTranslucent == TRUE)
    {
        colour.alpha = 0.0;
    }
    webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(options->webview), &colour);

    // set window's background color
    // Get the name of the current locale
    char *old_locale, *saved_locale;
    old_locale = setlocale(LC_ALL, NULL);

    // Copy the name so it won’t be clobbered by setlocale.
    saved_locale = strdup(old_locale);
    if (saved_locale == NULL)
        return;

    //Now change the locale to english for so printf always converts floats with a dot decimal separator
    setlocale(LC_ALL, "en_US.UTF-8");
    gchar *str = g_strdup_printf("#webview-box {background-color: rgba(%d, %d, %d, %1.1f);}", options->r, options->g, options->b, options->a / 255.0);

    //Restore the original locale.
    setlocale(LC_ALL, saved_locale);
    free(saved_locale);

    if (windowCssProvider == NULL)
    {
        windowCssProvider = gtk_css_provider_new();
        gtk_style_context_add_provider(
            gtk_widget_get_style_context(GTK_WIDGET(options->webviewBox)),
            GTK_STYLE_PROVIDER(windowCssProvider),
            GTK_STYLE_PROVIDER_PRIORITY_USER);
        g_object_unref(windowCssProvider);
    }

    // TODO: gtk_css_provider_load_from_data is deprecated since 4.12
    // but the user's system might not offer a compatible version.
    //
    // see: https://docs.gtk.org/gtk4/method.CssProvider.load_from_data.html
    gtk_css_provider_load_from_data(windowCssProvider, str, -1);

    g_free(str);
}

static gboolean setTitle(gpointer data)
{
    SetTitleArgs *args = (SetTitleArgs *)data;
    gtk_window_set_title(args->window, args->title);
    free((void *)args->title);
    free((void *)data);

    return G_SOURCE_REMOVE;
}

void SetTitle(GtkWindow *window, char *title)
{
    SetTitleArgs *args = malloc(sizeof(SetTitleArgs));
    args->window = window;
    args->title = title;
    ExecuteOnMainThread(setTitle, (gpointer)args);
}

//// TODO: gtk_window_move has been removed
// see: https://docs.gtk.org/gtk4/migrating-3to4.html#adapt-to-gtkwindow-api-changes
static gboolean setPosition(gpointer data)
{
    // SetPositionArgs *args = (SetPositionArgs *)data;
    // gtk_window_move((GtkWindow *)args->window, args->x, args->y);
    // free(args);

    return G_SOURCE_REMOVE;
}

//// TODO: gtk_window_move has been removed
// see: https://docs.gtk.org/gtk4/migrating-3to4.html#adapt-to-gtkwindow-api-changes
void SetPosition(void *window, int x, int y)
{
    // GdkRectangle monitorDimensions = getCurrentMonitorGeometry(window);
    // if (isNULLRectangle(monitorDimensions))
    // {
    //     return;
    // }
    // SetPositionArgs *args = malloc(sizeof(SetPositionArgs));
    // args->window = window;
    // args->x = monitorDimensions.x + x;
    // args->y = monitorDimensions.y + y;
    // ExecuteOnMainThread(setPosition, (gpointer)args);
}

//// TODO: gtk_window_set_geometry_hints has been removed
void SetMinMaxSize(GtkWindow *window, int min_width, int min_height, int max_width, int max_height)
{
    // GdkGeometry size;
    // size.min_width = size.min_height = size.max_width = size.max_height = 0;

    // GdkRectangle monitorSize = getCurrentMonitorGeometry(window);
    // if (isNULLRectangle(monitorSize))
    // {
    //     return;
    // }

    // int flags = GDK_HINT_MAX_SIZE | GDK_HINT_MIN_SIZE;

    // size.max_height = (max_height == 0 ? monitorSize.height : max_height);
    // size.max_width = (max_width == 0 ? monitorSize.width : max_width);
    // size.min_height = min_height;
    // size.min_width = min_width;

    // // On Wayland window manager get the decorators and calculate the differences from the windows' size.
    // if(onWayland()) 
    // {
    //     if(decoratorWidth == -1 && decoratorHeight == -1)
    //     {
    //         int windowWidth, windowHeight;
    //         gtk_window_get_size(window, &windowWidth, &windowHeight);

    //         GtkAllocation windowAllocation;
    //         gtk_widget_get_allocation(GTK_WIDGET(window), &windowAllocation);

    //         decoratorWidth = (windowAllocation.width-windowWidth);
    //         decoratorHeight = (windowAllocation.height-windowHeight);        
    //     }
    
    //     // Add the decorator difference to the window so fullscreen and maximise can fill the window.
    //     size.max_height = decoratorHeight+size.max_height;
    //     size.max_width = decoratorWidth+size.max_width;
    // }

    // gtk_window_set_geometry_hints(window, NULL, &size, flags);
}

// function to disable the context menu but propagate the event
static gboolean disableContextMenu(GtkWidget *widget, WebKitContextMenu *context_menu, GdkEvent *event, WebKitHitTestResult *hit_test_result, gpointer data)
{
    // return true to disable the context menu
    return TRUE;
}

void DisableContextMenu(void *webview)
{
    // Disable the context menu but propagate the event
    g_signal_connect(WEBKIT_WEB_VIEW(webview), "context-menu", G_CALLBACK(disableContextMenu), NULL);
}

static void buttonPress(GtkGestureClick* gesture, gint n_press, gdouble gesture_x, gdouble gesture_y, gpointer data)
{
    GdkEvent *event = gtk_event_controller_get_current_event(gesture);

    if (event == NULL)
    {
        xroot = yroot = 0.0f;
        dragTime = -1;
        return;
    }

    guint button = gtk_gesture_single_get_button(gesture);
    mouseButton = button;

    if (button == 3)
    {
        return;
    }

    if (gdk_event_get_event_type(event) == GDK_BUTTON_PRESS && button == 1)
    {
        double x, y;
        gboolean success = gdk_event_get_position(event, &x, &y);

        if(success) {
            xroot = x;
            yroot = y;
        }

        dragTime = gdk_event_get_time(event);
    }
}

static void buttonRelease(GtkGestureClick* gesture, gint n_press, gdouble gesture_x, gdouble gesture_y, gpointer data)
{
    GdkEvent *event = gtk_event_controller_get_current_event(gesture);

    if (event == NULL || 
        (gdk_event_get_event_type(event) == GDK_BUTTON_RELEASE && gtk_gesture_single_get_button(gesture) == 1))
    {
        xroot = yroot = 0.0f;
        dragTime = -1;
    }
}

void ConnectButtons(void *webview)
{
    GtkGesture *press = gtk_gesture_click_new();
    GtkGesture *release = gtk_gesture_click_new();

    gtk_widget_add_controller(GTK_WIDGET(webview), press);
    gtk_widget_add_controller(GTK_WIDGET(webview), release);

    g_signal_connect(press, "pressed", G_CALLBACK(buttonPress), NULL);
    g_signal_connect(release, "released", G_CALLBACK(buttonRelease), NULL);
}

int IsFullscreen(GtkWidget *widget)
{
    GtkWindow *gtkwindow = gtk_widget_get_root(widget);
    return gtk_window_is_fullscreen(gtkwindow);
}

int IsMaximised(GtkWidget *widget)
{
    GtkWindow *gtkwindow = gtk_widget_get_root(widget);
    return gtk_window_is_maximized(gtkwindow);
}

int IsMinimised(GtkWidget *widget)
{
    GtkWindow *gtkwindow = gtk_widget_get_root(widget);
    return !gtk_window_is_fullscreen(gtkwindow) && !gtk_window_is_maximized(gtkwindow);
}

//// TODO: gtk_window_move has been removed
// see: https://docs.gtk.org/gtk4/migrating-3to4.html#adapt-to-gtkwindow-api-changes
gboolean Center(gpointer data)
{
    // GtkWindow *window = (GtkWindow *)data;

    // // Get the geometry of the monitor
    // GdkRectangle m = getCurrentMonitorGeometry(window);
    // if (isNULLRectangle(m))
    // {
    //     return G_SOURCE_REMOVE;
    // }

    // // Get the window width/height
    // int windowWidth, windowHeight;
    // gtk_window_get_size(window, &windowWidth, &windowHeight);

    // int newX = ((m.width - windowWidth) / 2) + m.x;
    // int newY = ((m.height - windowHeight) / 2) + m.y;

    // // Place the window at the center of the monitor
    // gtk_window_move(window, newX, newY);

    return G_SOURCE_REMOVE;
}

gboolean Show(gpointer data)
{
    gtk_widget_show((GtkWidget *)data);

    return G_SOURCE_REMOVE;
}

gboolean Hide(gpointer data)
{
    gtk_widget_hide((GtkWidget *)data);

    return G_SOURCE_REMOVE;
}

gboolean Maximise(gpointer data)
{
    gtk_window_maximize((GtkWindow *)data);

    return G_SOURCE_REMOVE;
}

gboolean UnMaximise(gpointer data)
{
    gtk_window_unmaximize((GtkWindow *)data);

    return G_SOURCE_REMOVE;
}

gboolean Minimise(gpointer data)
{
    gtk_window_minimize((GtkWindow *)data);

    return G_SOURCE_REMOVE;
}

gboolean UnMinimise(gpointer data)
{
    gtk_window_present((GtkWindow *)data);

    return G_SOURCE_REMOVE;
}

gboolean Fullscreen(gpointer data)
{
    GtkWindow *window = (GtkWindow *)data;

    // Get the geometry of the monitor.
    GdkRectangle m = getCurrentMonitorGeometry(window);
    if (isNULLRectangle(m))
    {
        return G_SOURCE_REMOVE;
    }
    int scale = getCurrentMonitorScaleFactor(window);
    SetMinMaxSize(window, 0, 0, m.width * scale, m.height * scale);

    gtk_window_fullscreen(window);

    return G_SOURCE_REMOVE;
}

gboolean UnFullscreen(gpointer data)
{
    gtk_window_unfullscreen((GtkWindow *)data);

    return G_SOURCE_REMOVE;
}

static void webviewLoadChanged(WebKitWebView *web_view, WebKitLoadEvent load_event, gpointer data)
{
    if (load_event == WEBKIT_LOAD_FINISHED)
    {
        processMessage("DomReady");
    }
}

extern void processURLRequest(void *request);

void window_hide(GtkWindow* window, gpointer data) {
    gtk_widget_set_visible(GTK_WIDGET(window), false);
}

// This is called when the close button on the window is pressed
// gboolean close_button_pressed(GtkWidget *widget, GdkEvent *event, void *data)
gboolean close_button_pressed(GtkWindow* window, gpointer data)
{
    processMessage("Q");
    // since we handle the close in processMessage tell GTK to not invoke additional handlers - see:
    // https://docs.gtk.org/gtk3/signal.Widget.delete-event.html
    return TRUE;
}

char *droppedFiles = NULL;

// static void onDragDataReceived(GtkWidget *self, GdkDragContext *context, gint x, gint y, GtkSelectionData *selection_data, guint target_type, guint time, gpointer data)
// {
//     if(selection_data == NULL || (gtk_selection_data_get_length(selection_data) <= 0) || target_type != 2)
//     {
//         return;
//     }

//     if(droppedFiles != NULL) {
//         free(droppedFiles);
//         droppedFiles = NULL;
//     }

//     gchar **filenames = NULL;
//     filenames = g_uri_list_extract_uris((const gchar *)gtk_selection_data_get_data(selection_data));
//     if (filenames == NULL) // If unable to retrieve filenames:
//     {
//         g_strfreev(filenames);
//         return;
//     }

//     droppedFiles = calloc((size_t)gtk_selection_data_get_length(selection_data), 1);

//     int iter = 0;
//     while(filenames[iter] != NULL) // The last URI list element is NULL.
//     {
//         if(iter != 0)
//         {
//             strncat(droppedFiles, "\n", 1);
//         }
//         char *filename = g_filename_from_uri(filenames[iter], NULL, NULL);
//         if (filename == NULL)
//         {
//             break;
//         }
//         strncat(droppedFiles, filename, strlen(filename));

//         free(filename);
//         iter++;
//     }

//     g_strfreev(filenames);
// }

// static gboolean onDragDrop(GtkWidget* self, GdkDragContext* context, gint x, gint y, guint time, gpointer user_data)
// {
//     if(droppedFiles == NULL)
//     {
//         return FALSE;
//     }

//     size_t resLen = strlen(droppedFiles)+(sizeof(gint)*2)+6;
//     char *res = calloc(resLen, 1);

//     snprintf(res, resLen, "DD:%d:%d:%s", x, y, droppedFiles);

//     if(droppedFiles != NULL) {
//         free(droppedFiles);
//         droppedFiles = NULL;
//     }

//     processMessage(res);
//     return FALSE;
// }

static void onDelete(GtkWidget* self) {}

// WebView
GtkWidget *SetupWebview(void *contentManager, GtkWindow *window, int hideWindowOnClose, int gpuPolicy, int disableWebViewDragAndDrop, int enableDragAndDrop)
{
    GtkWidget *webview = GTK_WIDGET(g_object_new(WEBKIT_TYPE_WEB_VIEW, "user-content-manager", (WebKitUserContentManager *) contentManager, NULL));

    gtk_widget_set_vexpand(webview, true);

    WebKitWebContext *context = webkit_web_context_get_default();
    webkit_web_context_register_uri_scheme(context, "wails", (WebKitURISchemeRequestCallback)processURLRequest, NULL, NULL);
    g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webviewLoadChanged), NULL);

    // if(disableWebViewDragAndDrop)
    // {
    //     gtk_drag_dest_unset(webview);
    // }

    // if(enableDragAndDrop)
    // {
    //     g_signal_connect(G_OBJECT(webview), "drag-data-received", G_CALLBACK(onDragDataReceived), NULL);
    //     g_signal_connect(G_OBJECT(webview), "drag-drop", G_CALLBACK(onDragDrop), NULL);
    // }

    if (hideWindowOnClose)
    {
        g_signal_connect(window, "close-request", G_CALLBACK(window_hide), NULL);
    }
    else
    {
        g_signal_connect(window, "close-request", G_CALLBACK(close_button_pressed), NULL);
    }

    WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
    webkit_settings_set_user_agent_with_application_details(settings, "wails.io", "");

    switch (gpuPolicy)
    {
    // case 0:
    //     webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
    //     break;
    // case 1:
    //     webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER);
    //     break;
    default:
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
    }
    
    return webview;
}

void DevtoolsEnabled(void *webview, int enabled, bool showInspector)
{
    WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
    gboolean genabled = enabled == 1 ? true : false;
    webkit_settings_set_enable_developer_extras(settings, genabled);

    if (genabled && showInspector)
    {
        ShowInspector(webview);
    }
}

void LoadIndex(void *webview, char *url)
{
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), url);
}

static gboolean startDrag(gpointer data)
{
    DragOptions *options = (DragOptions *)data;

    // Ignore non-toplevel widgets
    GtkRoot *root = gtk_widget_get_root(GTK_WIDGET(options->webview)); 
    if (!GTK_IS_WINDOW(root))
    {
        free(data);
        return G_SOURCE_REMOVE;
    }

    gdk_toplevel_begin_move(options->mainwindow, NULL, mouseButton, xroot, yroot, dragTime);

    free(data);

    return G_SOURCE_REMOVE;
}

void StartDrag(void *webview, GtkWindow *mainwindow)
{
    DragOptions *data = malloc(sizeof(DragOptions));
    data->webview = webview;
    data->mainwindow = mainwindow;
    ExecuteOnMainThread(startDrag, (gpointer)data);
}

static gboolean startResize(gpointer data)
{
    ResizeOptions *options = (ResizeOptions *)data;

    // Ignore non-toplevel widgets
    GtkRoot *root = gtk_widget_get_root(GTK_WIDGET(options->webview)); 
    if (!GTK_IS_WINDOW(root))
    {
        free(data);
        return G_SOURCE_REMOVE;
    }

    gdk_toplevel_begin_resize(options->mainwindow, options->edge, NULL, mouseButton, xroot, yroot, dragTime);
    free(data);

    return G_SOURCE_REMOVE;
}

void StartResize(void *webview, GtkWindow *mainwindow, GdkSurfaceEdge edge)
{
    ResizeOptions *data = malloc(sizeof(ResizeOptions));
    data->webview = webview;
    data->mainwindow = mainwindow;
    data->edge = edge;
    ExecuteOnMainThread(startResize, (gpointer)data);
}

void ExecuteJS(void *data)
{
    struct JSCallback *js = data;
    webkit_web_view_evaluate_javascript(js->webview, js->script, -1, NULL, NULL, NULL, NULL, NULL);

    free(js->script);
}

void extern processMessageDialogResult(char *);

void messageResult(GtkDialog* dialog, gint response_id, gpointer user_data) {
    if(response_id == GTK_RESPONSE_YES) {
        processMessageDialogResult("Yes");
    } else if(response_id == GTK_RESPONSE_NO) {
        processMessageDialogResult("No");
    } else if(response_id == GTK_RESPONSE_OK) {
        processMessageDialogResult("OK");
    } else if(response_id == GTK_RESPONSE_CANCEL) {
        processMessageDialogResult("Cancel");
    } else {
        processMessageDialogResult("");
    }

    gtk_window_destroy(GTK_WINDOW(dialog));
}

void MessageDialog(void *data)
{
    GtkDialogFlags flags;
    GtkMessageType messageType;
    MessageDialogOptions *options = (MessageDialogOptions *)data;
    if (options->messageType == 0)
    {
        messageType = GTK_MESSAGE_INFO;
        flags = GTK_BUTTONS_OK;
    }
    else if (options->messageType == 1)
    {
        messageType = GTK_MESSAGE_ERROR;
        flags = GTK_BUTTONS_OK;
    }
    else if (options->messageType == 2)
    {
        messageType = GTK_MESSAGE_QUESTION;
        flags = GTK_BUTTONS_YES_NO;
    }
    else
    {
        messageType = GTK_MESSAGE_WARNING;
        flags = GTK_BUTTONS_OK;
    }

    // TODO: gtk_message_dialog_new is deprecated since 4.10
    // but the user's system might not offer a compatible version.
    //
    // see: https://docs.gtk.org/gtk4/ctor.MessageDialog.new.html
    GtkWidget *dialog;
    dialog = gtk_message_dialog_new(GTK_WINDOW(options->window),
                                    GTK_DIALOG_DESTROY_WITH_PARENT,
                                    messageType,
                                    flags,
                                    options->message, NULL);
    
    g_object_ref_sink(dialog);

    gtk_window_set_title(GTK_WINDOW(dialog), options->title);
    gtk_window_set_modal(GTK_WINDOW(dialog), true);

    g_signal_connect(dialog, "response", G_CALLBACK(messageResult), NULL);
    
    gtk_widget_show(dialog);

    free(options->title);
    free(options->message);
}

void extern processOpenFileResult(void *);

GtkFileFilter **AllocFileFilterArray(size_t ln)
{
    return (GtkFileFilter **)malloc(ln * sizeof(GtkFileFilter *));
}

void freeFileFilterArray(GtkFileFilter **filters)
{
    free(filters);
}

void openFileResult(GtkDialog *dialog, int response) {
    GtkFileChooser *fc = GTK_FILE_CHOOSER(dialog);

    // Max 1024 files to select
    char **result = calloc(1024, sizeof(char *));
    int resultIndex = 0;

    if(response == GTK_RESPONSE_ACCEPT) {
        GListModel *files = gtk_file_chooser_get_files(fc);

        GObject *item = g_list_model_get_object(files, resultIndex);

        while(item) {
            GFile *file = G_FILE(item);
            char *path = g_file_get_path(file);

            result[resultIndex] = path;
            resultIndex++;

            g_object_unref(file);

            if(resultIndex == 1024) {
                break;
            }

            item = g_list_model_get_object(files, resultIndex);
        }

        processOpenFileResult(result);

        for(int i = 0; i < resultIndex; i++) {
            g_free(result[i]);
        }

        g_object_unref(files);
    } else {
        processOpenFileResult(result);
    }
    free(result);

    gtk_window_destroy(GTK_WINDOW(dialog));
}

void Opendialog(void *data)
{
    struct OpenFileDialogOptions *options = data;
    char *label = "_Open";
    if (options->action == GTK_FILE_CHOOSER_ACTION_SAVE)
    {
        label = "_Save";
    }

    // TODO: gtk_file_chooser_dialog_new is deprecated since 4.10
    // but the user's system might not offer a compatible version.
    //
    // see: https://docs.gtk.org/gtk4/class.FileChooserDialog.html
    GtkWidget *dialog = gtk_file_chooser_dialog_new(options->title, options->window, options->action,
                                                       "_Cancel", GTK_RESPONSE_CANCEL,
                                                       label, GTK_RESPONSE_ACCEPT,
                                                       NULL);

    g_object_ref_sink(dialog);

    // TODO: GtkFileChooser is deprecated since 4.10
    // but the user's system might not offer a compatible version.
    //
    // see: https://docs.gtk.org/gtk4/iface.FileChooser.html
    GtkFileChooser *fc = GTK_FILE_CHOOSER(dialog);

    // filters
    if (options->filters != 0)
    {
        int index = 0;
        GtkFileFilter *thisFilter;
        while (options->filters[index] != NULL)
        {
            thisFilter = options->filters[index];
            gtk_file_chooser_add_filter(fc, thisFilter);
            index++;
        }
    }

    if (options->multipleFiles == 1)
    {
        gtk_file_chooser_set_select_multiple(fc, TRUE);
    }

    if (options->createDirectories == 1)
    {
        gtk_file_chooser_set_create_folders(fc, TRUE);
    }

    if (options->defaultDirectory != NULL)
    {
        // TODO: gtk_file_chooser_set_current_folder is deprecated since 4.10
        // but the user's system might not offer a compatible version.
        //
        // see: https://docs.gtk.org/gtk4/method.FileChooser.set_current_folder.html
        gtk_file_chooser_set_current_folder(fc, options->defaultDirectory, NULL);
        free(options->defaultDirectory);
    }

    if (options->action == GTK_FILE_CHOOSER_ACTION_SAVE)
    {
        if (options->defaultFilename != NULL)
        {
            gtk_file_chooser_set_current_name(fc, options->defaultFilename);
            free(options->defaultFilename);
        }
    }

    g_signal_connect(dialog, "response", G_CALLBACK(openFileResult), NULL);
    
    gtk_widget_show(dialog);

    // Release filters
    if (options->filters != NULL)
    {
        int index = 0;
        GtkFileFilter *thisFilter;
        while (options->filters[index] != 0)
        {
            thisFilter = options->filters[index];
            g_object_unref(thisFilter);
            index++;
        }
        freeFileFilterArray(options->filters);
    }

    free(options->title);
}

GtkFileFilter *newFileFilter()
{
    GtkFileFilter *result = gtk_file_filter_new();
    g_object_ref(result);
    return result;
}

void ShowInspector(void *webview) {
    WebKitWebInspector *inspector = webkit_web_view_get_inspector(WEBKIT_WEB_VIEW(webview));
    webkit_web_inspector_show(WEBKIT_WEB_INSPECTOR(inspector));
}

void sendShowInspectorMessage(GAction *action, GVariant *param) {
    processMessage("wails:showInspector");
}

// When the user presses Ctrl+Shift+F12, call ShowInspector
void InstallF12Hotkey(GtkApplication *app, GtkWindow *window)
{
    GSimpleAction *action = g_simple_action_new("show-inspector", NULL);
    g_signal_connect(action, "activate", G_CALLBACK(sendShowInspectorMessage), NULL);
    g_action_map_add_action(G_ACTION_MAP(window), G_ACTION(action));

    gtk_application_set_accels_for_action(
        app, 
        "win.show-inspector", 
        (const char *[]) { "<Control><Shift>F12", NULL });
}

extern void onActivate();

const int G_APPLICATION_DEFAULT_FLAGS = 0;

static void activate(GtkApplication *app, gpointer user_data) {
	onActivate();
}

GtkApplication* createApp(char *appId) {
	GtkApplication *app = gtk_application_new(appId, G_APPLICATION_DEFAULT_FLAGS);
	g_signal_connect(app, "activate", G_CALLBACK(activate), NULL);
	return app;
}

void runApp(GtkApplication *app) {
	g_application_run(G_APPLICATION(app), 0, NULL);
	g_object_unref(app);
}
