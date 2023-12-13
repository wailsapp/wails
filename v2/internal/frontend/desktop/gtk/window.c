#include <JavaScriptCore/JavaScript.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#include <locale.h>
#include "window.h"

// These are the x,y,time & button of the last mouse down event
// It's used for window dragging
static float xroot = 0.0f;
static float yroot = 0.0f;
static int dragTime = -1;
static uint mouseButton = 0;

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

GtkContainer *GTKCONTAINER(void *pointer)
{
    return GTK_CONTAINER(pointer);
}

GtkBox *GTKBOX(void *pointer)
{
    return GTK_BOX(pointer);
}

extern void processMessage(char *);

static void sendMessageToBackend(WebKitUserContentManager *contentManager,
                                 WebKitJavascriptResult *result,
                                 void *data)
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
    processMessage(message);
    g_free(message);
}

static bool isNULLRectangle(GdkRectangle input)
{
    return input.x == -1 && input.y == -1 && input.width == -1 && input.height == -1;
}

static GdkMonitor *getCurrentMonitor(GtkWindow *window)
{
    // Get the monitor that the window is currently on
    GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
    GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
    if (gdk_window == NULL)
    {
        return NULL;
    }
    GdkMonitor *monitor = gdk_display_get_monitor_at_window(display, gdk_window);

    return GDK_MONITOR(monitor);
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
    GdkPixbufLoader *loader = gdk_pixbuf_loader_new();
    if (!loader)
    {
        return;
    }
    if (gdk_pixbuf_loader_write(loader, buf, len, NULL) && gdk_pixbuf_loader_close(loader, NULL))
    {
        GdkPixbuf *pixbuf = gdk_pixbuf_loader_get_pixbuf(loader);
        if (pixbuf)
        {
            gtk_window_set_icon(window, pixbuf);
        }
    }
    g_object_unref(loader);
}

void SetWindowTransparency(GtkWidget *widget)
{
    GdkScreen *screen = gtk_widget_get_screen(widget);
    GdkVisual *visual = gdk_screen_get_rgba_visual(screen);

    if (visual != NULL && gdk_screen_is_composited(screen))
    {
        gtk_widget_set_app_paintable(widget, true);
        gtk_widget_set_visual(widget, visual);
    }
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

    gtk_css_provider_load_from_data(windowCssProvider, str, -1, NULL);
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

static gboolean setPosition(gpointer data)
{
    SetPositionArgs *args = (SetPositionArgs *)data;
    gtk_window_move((GtkWindow *)args->window, args->x, args->y);
    free(args);

    return G_SOURCE_REMOVE;
}

void SetPosition(void *window, int x, int y)
{
    GdkRectangle monitorDimensions = getCurrentMonitorGeometry(window);
    if (isNULLRectangle(monitorDimensions))
    {
        return;
    }
    SetPositionArgs *args = malloc(sizeof(SetPositionArgs));
    args->window = window;
    args->x = monitorDimensions.x + x;
    args->y = monitorDimensions.y + y;
    ExecuteOnMainThread(setPosition, (gpointer)args);
}

void SetMinMaxSize(GtkWindow *window, int min_width, int min_height, int max_width, int max_height)
{
    GdkGeometry size;
    size.min_width = size.min_height = size.max_width = size.max_height = 0;

    GdkRectangle monitorSize = getCurrentMonitorGeometry(window);
    if (isNULLRectangle(monitorSize))
    {
        return;
    }
    int flags = GDK_HINT_MAX_SIZE | GDK_HINT_MIN_SIZE;
    size.max_height = (max_height == 0 ? monitorSize.height : max_height);
    size.max_width = (max_width == 0 ? monitorSize.width : max_width);
    size.min_height = min_height;
    size.min_width = min_width;
    gtk_window_set_geometry_hints(window, NULL, &size, flags);
}

// function to disable the context menu but propogate the event
static gboolean disableContextMenu(GtkWidget *widget, WebKitContextMenu *context_menu, GdkEvent *event, WebKitHitTestResult *hit_test_result, gpointer data)
{
    // return true to disable the context menu
    return TRUE;
}

void DisableContextMenu(void *webview)
{
    // Disable the context menu but propogate the event
    g_signal_connect(WEBKIT_WEB_VIEW(webview), "context-menu", G_CALLBACK(disableContextMenu), NULL);
}

static gboolean buttonPress(GtkWidget *widget, GdkEventButton *event, void *dummy)
{
    if (event == NULL)
    {
        xroot = yroot = 0.0f;
        dragTime = -1;
        return FALSE;
    }
    mouseButton = event->button;
    if (event->button == 3)
    {
        return FALSE;
    }

    if (event->type == GDK_BUTTON_PRESS && event->button == 1)
    {
        xroot = event->x_root;
        yroot = event->y_root;
        dragTime = event->time;
    }

    return FALSE;
}

static gboolean buttonRelease(GtkWidget *widget, GdkEventButton *event, void *dummy)
{
    if (event == NULL || (event->type == GDK_BUTTON_RELEASE && event->button == 1))
    {
        xroot = yroot = 0.0f;
        dragTime = -1;
    }
    return FALSE;
}

void ConnectButtons(void *webview)
{
    g_signal_connect(WEBKIT_WEB_VIEW(webview), "button-press-event", G_CALLBACK(buttonPress), NULL);
    g_signal_connect(WEBKIT_WEB_VIEW(webview), "button-release-event", G_CALLBACK(buttonRelease), NULL);
}

int IsFullscreen(GtkWidget *widget)
{
    GdkWindow *gdkwindow = gtk_widget_get_window(widget);
    GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
    return state & GDK_WINDOW_STATE_FULLSCREEN;
}

int IsMaximised(GtkWidget *widget)
{
    GdkWindow *gdkwindow = gtk_widget_get_window(widget);
    GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
    return state & GDK_WINDOW_STATE_MAXIMIZED && !(state & GDK_WINDOW_STATE_FULLSCREEN);
}

int IsMinimised(GtkWidget *widget)
{
    GdkWindow *gdkwindow = gtk_widget_get_window(widget);
    GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
    return state & GDK_WINDOW_STATE_ICONIFIED;
}

gboolean Center(gpointer data)
{
    GtkWindow *window = (GtkWindow *)data;

    // Get the geometry of the monitor
    GdkRectangle m = getCurrentMonitorGeometry(window);
    if (isNULLRectangle(m))
    {
        return G_SOURCE_REMOVE;
    }

    // Get the window width/height
    int windowWidth, windowHeight;
    gtk_window_get_size(window, &windowWidth, &windowHeight);

    int newX = ((m.width - windowWidth) / 2) + m.x;
    int newY = ((m.height - windowHeight) / 2) + m.y;

    // Place the window at the center of the monitor
    gtk_window_move(window, newX, newY);

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
    gtk_window_iconify((GtkWindow *)data);

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

// This is called when the close button on the window is pressed
gboolean close_button_pressed(GtkWidget *widget, GdkEvent *event, void *data)
{
    processMessage("Q");
    // since we handle the close in processMessage tell GTK to not invoke additional handlers - see:
    // https://docs.gtk.org/gtk3/signal.Widget.delete-event.html
    return TRUE;
}

// WebView
GtkWidget *SetupWebview(void *contentManager, GtkWindow *window, int hideWindowOnClose, int gpuPolicy)
{
    GtkWidget *webview = webkit_web_view_new_with_user_content_manager((WebKitUserContentManager *)contentManager);
    // gtk_container_add(GTK_CONTAINER(window), webview);
    WebKitWebContext *context = webkit_web_context_get_default();
    webkit_web_context_register_uri_scheme(context, "wails", (WebKitURISchemeRequestCallback)processURLRequest, NULL, NULL);
    g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webviewLoadChanged), NULL);
    if (hideWindowOnClose)
    {
        g_signal_connect(GTK_WIDGET(window), "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), NULL);
    }
    else
    {
        g_signal_connect(GTK_WIDGET(window), "delete-event", G_CALLBACK(close_button_pressed), NULL);
    }

    WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
    webkit_settings_set_user_agent_with_application_details(settings, "wails.io", "");

    switch (gpuPolicy)
    {
    case 0:
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
        break;
    case 1:
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND);
        break;
    case 2:
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER);
        break;
    default:
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND);
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
    GtkWidget *window = gtk_widget_get_toplevel(GTK_WIDGET(options->webview));
    if (!GTK_IS_WINDOW(window))
    {
        free(data);
        return G_SOURCE_REMOVE;
    }

    gtk_window_begin_move_drag(options->mainwindow, mouseButton, xroot, yroot, dragTime);
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
    GtkWidget *window = gtk_widget_get_toplevel(GTK_WIDGET(options->webview));
    if (!GTK_IS_WINDOW(window))
    {
        free(data);
        return G_SOURCE_REMOVE;
    }

    gtk_window_begin_resize_drag(options->mainwindow, options->edge, mouseButton, xroot, yroot, dragTime);
    free(data);

    return G_SOURCE_REMOVE;
}

void StartResize(void *webview, GtkWindow *mainwindow, GdkWindowEdge edge)
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
    webkit_web_view_run_javascript(js->webview, js->script, NULL, NULL, NULL);
    free(js->script);
}

void extern processMessageDialogResult(char *);

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

    GtkWidget *dialog;
    dialog = gtk_message_dialog_new(GTK_WINDOW(options->window),
                                    GTK_DIALOG_DESTROY_WITH_PARENT,
                                    messageType,
                                    flags,
                                    options->message, NULL);
    gtk_window_set_title(GTK_WINDOW(dialog), options->title);
    GtkResponseType result = gtk_dialog_run(GTK_DIALOG(dialog));
    if (result == GTK_RESPONSE_YES)
    {
        processMessageDialogResult("Yes");
    }
    else if (result == GTK_RESPONSE_NO)
    {
        processMessageDialogResult("No");
    }
    else if (result == GTK_RESPONSE_OK)
    {
        processMessageDialogResult("OK");
    }
    else if (result == GTK_RESPONSE_CANCEL)
    {
        processMessageDialogResult("Cancel");
    }
    else
    {
        processMessageDialogResult("");
    }

    gtk_widget_destroy(dialog);
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

void Opendialog(void *data)
{
    struct OpenFileDialogOptions *options = data;
    char *label = "_Open";
    if (options->action == GTK_FILE_CHOOSER_ACTION_SAVE)
    {
        label = "_Save";
    }
    GtkWidget *dlgWidget = gtk_file_chooser_dialog_new(options->title, options->window, options->action,
                                                       "_Cancel", GTK_RESPONSE_CANCEL,
                                                       label, GTK_RESPONSE_ACCEPT,
                                                       NULL);

    GtkFileChooser *fc = GTK_FILE_CHOOSER(dlgWidget);
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

    gtk_file_chooser_set_local_only(fc, FALSE);

    if (options->multipleFiles == 1)
    {
        gtk_file_chooser_set_select_multiple(fc, TRUE);
    }
    gtk_file_chooser_set_do_overwrite_confirmation(fc, TRUE);
    if (options->createDirectories == 1)
    {
        gtk_file_chooser_set_create_folders(fc, TRUE);
    }
    if (options->showHiddenFiles == 1)
    {
        gtk_file_chooser_set_show_hidden(fc, TRUE);
    }

    if (options->defaultDirectory != NULL)
    {
        gtk_file_chooser_set_current_folder(fc, options->defaultDirectory);
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

    gint response = gtk_dialog_run(GTK_DIALOG(dlgWidget));

    // Max 1024 files to select
    char **result = calloc(1024, sizeof(char *));
    int resultIndex = 0;

    if (response == GTK_RESPONSE_ACCEPT)
    {
        GSList *filenames = gtk_file_chooser_get_filenames(fc);
        GSList *iter = filenames;
        while (iter)
        {
            result[resultIndex++] = (char *)iter->data;
            iter = g_slist_next(iter);
            if (resultIndex == 1024)
            {
                break;
            }
        }
        processOpenFileResult(result);
        iter = filenames;
        while (iter)
        {
            g_free(iter->data);
            iter = g_slist_next(iter);
        }
    }
    else
    {
        processOpenFileResult(result);
    }
    free(result);

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
    gtk_widget_destroy(dlgWidget);
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

void sendShowInspectorMessage() {
    processMessage("wails:showInspector");
}

void InstallF12Hotkey(void *window)
{
    // When the user presses Ctrl+Shift+F12, call ShowInspector
    GtkAccelGroup *accel_group = gtk_accel_group_new();
    gtk_window_add_accel_group(GTK_WINDOW(window), accel_group);
    GClosure *closure = g_cclosure_new(G_CALLBACK(sendShowInspectorMessage), window, NULL);
    gtk_accel_group_connect(accel_group, GDK_KEY_F12, GDK_CONTROL_MASK | GDK_SHIFT_MASK, GTK_ACCEL_VISIBLE, closure);
}