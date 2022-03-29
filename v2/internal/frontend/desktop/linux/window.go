//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <JavaScriptCore/JavaScript.h>
#include "gtk/gtk.h"
#include "webkit2/webkit2.h"
#include <stdio.h>
#include <limits.h>
#include "stdint.h"


void ExecuteOnMainThread(void* f, gpointer jscallback) {
    g_idle_add((GSourceFunc)f, (gpointer)jscallback);
}

static GtkWidget* GTKWIDGET(void *pointer) {
	return GTK_WIDGET(pointer);
}

static GtkWindow* GTKWINDOW(void *pointer) {
	return GTK_WINDOW(pointer);
}

static GtkContainer* GTKCONTAINER(void *pointer) {
	return GTK_CONTAINER(pointer);
}

static GtkBox* GTKBOX(void *pointer) {
	return GTK_BOX(pointer);
}

static void SetMinMaxSize(GtkWindow* window, int min_width, int min_height, int max_width, int max_height) {
    GdkGeometry size;
    size.min_width = size.min_height = size.max_width = size.max_height = 0;
    int flags = GDK_HINT_MAX_SIZE | GDK_HINT_MIN_SIZE;
	size.max_height = (max_height == 0 ? INT_MAX : max_height);
	size.max_width = (max_width == 0 ? INT_MAX : max_width);
	size.min_height = min_height;
	size.min_width = min_width;
    gtk_window_set_geometry_hints(window, NULL, &size, flags);
}

GdkMonitor* getCurrentMonitor(GtkWindow *window) {
	// Get the monitor that the window is currently on
	GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
	GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
	GdkMonitor *monitor = gdk_display_get_monitor_at_window(display, gdk_window);

	return GDK_MONITOR(monitor);
}

GdkRectangle getCurrentMonitorGeometry(GtkWindow *window) {
	GdkMonitor *monitor = getCurrentMonitor(window);

	// Get the geometry of the monitor
	GdkRectangle result;
	gdk_monitor_get_geometry (monitor,&result);
	return result;
}

int getCurrentMonitorScaleFactor(GtkWindow *window) {
	GdkMonitor *monitor = getCurrentMonitor(window);

	return gdk_monitor_get_scale_factor(monitor);
}

gboolean Center(gpointer data) {
	GtkWindow *window = (GtkWindow*)data;

    // Get the geometry of the monitor
    GdkRectangle m = getCurrentMonitorGeometry(window);

    // Get the window width/height
    int windowWidth, windowHeight;
    gtk_window_get_size(window, &windowWidth, &windowHeight);

	int newX = ((m.width - windowWidth) / 2) + m.x;
	int newY = ((m.height - windowHeight) / 2) + m.y;

    // Place the window at the center of the monitor
    gtk_window_move(window, newX, newY);

    return G_SOURCE_REMOVE;
}

int IsFullscreen(GtkWidget *widget) {
	GdkWindow *gdkwindow = gtk_widget_get_window(widget);
	GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
	return state & GDK_WINDOW_STATE_FULLSCREEN;
}

int IsMaximised(GtkWidget *widget) {
	GdkWindow *gdkwindow = gtk_widget_get_window(widget);
	GdkWindowState state = gdk_window_get_state(GDK_WINDOW(gdkwindow));
	return state & GDK_WINDOW_STATE_MAXIMIZED;
}


extern void processMessage(char*);

static void sendMessageToBackend(WebKitUserContentManager *contentManager,
                                 WebKitJavascriptResult *result,
                                 void* data)
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

static void webviewLoadChanged(WebKitWebView *web_view, WebKitLoadEvent load_event, gpointer data)
{
    if (load_event == WEBKIT_LOAD_FINISHED) {
        processMessage("DomReady");
    }
}

ulong setupInvokeSignal(void* contentManager) {
	return g_signal_connect((WebKitUserContentManager*)contentManager, "script-message-received::external", G_CALLBACK(sendMessageToBackend), NULL);
}

// These are the x,y & time of the last mouse down event
// It's used for window dragging
float xroot = 0.0f;
float yroot = 0.0f;
int dragTime = -1;
bool contextMenuDisabled = false;

gboolean buttonPress(GtkWidget *widget, GdkEventButton *event, void* dummy)
{
	if( event == NULL ) {
		xroot = yroot = 0.0f;
		dragTime = -1;
		return FALSE;
	}

	if( event->button == 3 && contextMenuDisabled ) {
		return TRUE;
	}

	if (event->type == GDK_BUTTON_PRESS && event->button == 1)
    {
        xroot = event->x_root;
        yroot = event->y_root;
        dragTime = event->time;
    }

    return FALSE;
}

gboolean buttonRelease(GtkWidget *widget, GdkEventButton *event, void* dummy)
{
    if (event == NULL || (event->type == GDK_BUTTON_RELEASE && event->button == 1))
    {
		xroot = yroot = 0.0f;
		dragTime = -1;
    }
    return FALSE;
}

void connectButtons(void* webview) {
	g_signal_connect(WEBKIT_WEB_VIEW(webview), "button-press-event", G_CALLBACK(buttonPress), NULL);
	g_signal_connect(WEBKIT_WEB_VIEW(webview), "button-release-event", G_CALLBACK(buttonRelease), NULL);
}

extern void processURLRequest(void *request);

// This is called when the close button on the window is pressed
gboolean close_button_pressed(GtkWidget *widget, GdkEvent *event, void* data)
{
   	processMessage("Q");
    return FALSE;
}

GtkWidget* setupWebview(void* contentManager, GtkWindow* window, int hideWindowOnClose) {
	GtkWidget* webview = webkit_web_view_new_with_user_content_manager((WebKitUserContentManager*)contentManager);
	//gtk_container_add(GTK_CONTAINER(window), webview);
	WebKitWebContext *context = webkit_web_context_get_default();
	webkit_web_context_register_uri_scheme(context, "wails", (WebKitURISchemeRequestCallback)processURLRequest, NULL, NULL);
	g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webviewLoadChanged), NULL);
	if (hideWindowOnClose) {
		g_signal_connect(GTK_WIDGET(window), "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), NULL);
	} else {
		g_signal_connect(GTK_WIDGET(window), "delete-event", G_CALLBACK(close_button_pressed), NULL);
	}
	return webview;
}

void devtoolsEnabled(void* webview, int enabled) {
	WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
	gboolean genabled = enabled == 1 ? true : false;
	webkit_settings_set_enable_developer_extras(settings, genabled);
}

void loadIndex(void* webview) {
	webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), "wails:///");
}

typedef struct DragOptions {
	void *webview;
	GtkWindow* mainwindow;
} DragOptions;

static gboolean startDrag(gpointer data) {
	DragOptions* options = (DragOptions*)data;

	// Ignore non-toplevel widgets
	GtkWidget *window = gtk_widget_get_toplevel(GTK_WIDGET(options->webview));
	if (!GTK_IS_WINDOW(window)) {
		free(data);
		return G_SOURCE_REMOVE;
	}

	gtk_window_begin_move_drag(options->mainwindow, 1, xroot, yroot, dragTime);
	free(data);

	return G_SOURCE_REMOVE;
}

static void StartDrag(void *webview, GtkWindow* mainwindow) {
	DragOptions* data = malloc(sizeof(DragOptions));
	data->webview = webview;
	data->mainwindow = mainwindow;
	ExecuteOnMainThread(startDrag, (gpointer)data);
}

typedef struct JSCallback {
    void* webview;
    char* script;
} JSCallback;

gboolean executeJS(gpointer data) {
    struct JSCallback *js = data;
    webkit_web_view_run_javascript(js->webview, js->script, NULL, NULL, NULL);
    free(js->script);
    return G_SOURCE_REMOVE;
}

void extern processMessageDialogResult(char*);

typedef struct MessageDialogOptions {
	void* window;
	char* title;
	char* message;
	int messageType;
} MessageDialogOptions;

gboolean messageDialog(gpointer data) {

	GtkDialogFlags flags;
	GtkMessageType messageType;
	MessageDialogOptions *options = (MessageDialogOptions*) data;
	if( options->messageType == 0 ) {
		messageType = GTK_MESSAGE_INFO;
		flags = GTK_BUTTONS_OK;
	} else if( options->messageType == 1 ) {
		messageType = GTK_MESSAGE_ERROR;
		flags = GTK_BUTTONS_OK;
	} else if( options->messageType == 2 ) {
		messageType = GTK_MESSAGE_QUESTION;
		flags = GTK_BUTTONS_YES_NO;
	} else {
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
	if ( result == GTK_RESPONSE_YES ) {
		processMessageDialogResult("Yes");
	} else if ( result == GTK_RESPONSE_NO ) {
		processMessageDialogResult("No");
	} else if ( result == GTK_RESPONSE_OK ) {
		processMessageDialogResult("OK");
	} else if ( result == GTK_RESPONSE_CANCEL ) {
		processMessageDialogResult("Cancel");
	} else {
		processMessageDialogResult("");
	}

	gtk_widget_destroy(dialog);
	free(options->title);
	free(options->message);

	return G_SOURCE_REMOVE;
}

void extern processOpenFileResult(void*);

typedef struct OpenFileDialogOptions {
    void* webview;
    char* title;
	char* defaultFilename;
	char* defaultDirectory;
	int createDirectories;
	int multipleFiles;
	int showHiddenFiles;
 	GtkFileChooserAction action;
	GtkFileFilter** filters;
} OpenFileDialogOptions;

GtkFileFilter** allocFileFilterArray(size_t ln) {
	return (GtkFileFilter**) malloc(ln * sizeof(GtkFileFilter*));
}

void freeFileFilterArray(GtkFileFilter** filters) {
	free(filters);
}

gboolean opendialog(gpointer data) {
    struct OpenFileDialogOptions *options = data;
	char *label = "_Open";
	if (options->action == GTK_FILE_CHOOSER_ACTION_SAVE) {
		label = "_Save";
	}
    GtkWidget *dlgWidget = gtk_file_chooser_dialog_new(options->title, options->webview, options->action,
          "_Cancel", GTK_RESPONSE_CANCEL,
          label, GTK_RESPONSE_ACCEPT,
			NULL);

	GtkFileChooser *fc = GTK_FILE_CHOOSER(dlgWidget);
	// filters
	if (options->filters != 0) {
		int index = 0;
		GtkFileFilter* thisFilter;
		while(options->filters[index] != NULL) {
			thisFilter = options->filters[index];
			gtk_file_chooser_add_filter(fc, thisFilter);
			index++;
		}
	}

	gtk_file_chooser_set_local_only(fc, FALSE);

	if (options->multipleFiles == 1) {
		gtk_file_chooser_set_select_multiple(fc, TRUE);
	}
	gtk_file_chooser_set_do_overwrite_confirmation(fc, TRUE);
	if (options->createDirectories == 1) {
		gtk_file_chooser_set_create_folders(fc, TRUE);
	}
	if (options->showHiddenFiles == 1) {
		gtk_file_chooser_set_show_hidden(fc, TRUE);
	}

	if (options->defaultDirectory != NULL) {
		gtk_file_chooser_set_current_folder (fc, options->defaultDirectory);
		free(options->defaultDirectory);
	}

	if (options->action == GTK_FILE_CHOOSER_ACTION_SAVE) {
		if (options->defaultFilename != NULL) {
			gtk_file_chooser_set_current_name(fc, options->defaultFilename);
			free(options->defaultFilename);
		}
	}

	gint response = gtk_dialog_run(GTK_DIALOG(dlgWidget));

	// Max 1024 files to select
	char** result = calloc(1024, sizeof(char*));
	int resultIndex = 0;

    if (response == GTK_RESPONSE_ACCEPT) {
        GSList* filenames = gtk_file_chooser_get_filenames(fc);
		GSList *iter = filenames;
		while(iter) {
		  	result[resultIndex++] = (char *)iter->data;
		  	iter = g_slist_next(iter);
          	if (resultIndex == 1024) {
				break;
			}
		}
		processOpenFileResult(result);
		iter = filenames;
		while(iter) {
		  g_free(iter->data);
		  iter = g_slist_next(iter);
		}
    } else {
		processOpenFileResult(result);
	}
	free(result);

	// Release filters
	if (options->filters != NULL) {
		int index = 0;
		GtkFileFilter* thisFilter;
		while(options->filters[index] != 0) {
			thisFilter = options->filters[index];
			g_object_unref(thisFilter);
			index++;
		}
		freeFileFilterArray(options->filters);
	}
    gtk_widget_destroy(dlgWidget);
    free(options->title);
    return G_SOURCE_REMOVE;
}

GtkFileFilter* newFileFilter() {
	GtkFileFilter* result = gtk_file_filter_new();
	g_object_ref(result);
	return result;
}

typedef struct RGBAOptions {
	uint8_t r;
	uint8_t g;
	uint8_t b;
	uint8_t a;
	void *webview;
} RGBAOptions;

gboolean setRGBA(gpointer* data) {
	RGBAOptions* options = (RGBAOptions*)data;
	GdkRGBA colour = {options->r / 255.0, options->g / 255.0, options->b / 255.0, options->a / 255.0};
	webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(options->webview), &colour);

	return G_SOURCE_REMOVE;
}

typedef struct SetTitleArgs {
	GtkWindow* window;
	char* title;
} SetTitleArgs;

gboolean setTitle(gpointer data) {
	SetTitleArgs* args = (SetTitleArgs*)data;
	gtk_window_set_title(args->window, args->title);
	free((void*)args->title);
	free((void*)data);

	return G_SOURCE_REMOVE;
}

void SetTitle(GtkWindow* window, char* title) {
	SetTitleArgs* args = malloc(sizeof(SetTitleArgs));
	args->window = window;
	args->title = title;
	ExecuteOnMainThread(setTitle, (gpointer)args);
}

typedef struct SetPositionArgs {
	int x;
	int y;
	void* window;
} SetPositionArgs;

gboolean setPosition(gpointer data) {
	SetPositionArgs* args = (SetPositionArgs*)data;
	gtk_window_move((GtkWindow*)args->window, args->x, args->y);
	free(args);

	return G_SOURCE_REMOVE;
}

void SetPosition(void* window, int x, int y) {
	GdkRectangle monitorDimensions = getCurrentMonitorGeometry(window);
	SetPositionArgs* args = malloc(sizeof(SetPositionArgs));
	args->window = window;
	args->x = monitorDimensions.x + x;
	args->y = monitorDimensions.y + y;
	ExecuteOnMainThread(setPosition, (gpointer)args);
}

gboolean Show(gpointer data) {
	gtk_widget_show((GtkWidget*)data);

	return G_SOURCE_REMOVE;
}

gboolean Hide(gpointer data) {
	gtk_widget_hide((GtkWidget*)data);

	return G_SOURCE_REMOVE;
}

gboolean Maximise(gpointer data) {
	gtk_window_maximize((GtkWindow*)data);

	return G_SOURCE_REMOVE;
}

gboolean UnMaximise(gpointer data) {
	gtk_window_unmaximize((GtkWindow*)data);

	return G_SOURCE_REMOVE;
}

gboolean Minimise(gpointer data) {
	gtk_window_iconify((GtkWindow*)data);

	return G_SOURCE_REMOVE;
}

gboolean UnMinimise(gpointer data) {
	gtk_window_present((GtkWindow*)data);

	return G_SOURCE_REMOVE;
}

gboolean Fullscreen(gpointer data) {
	GtkWindow* window = (GtkWindow*)data;

	// Get the geometry of the monitor.
	GdkRectangle m = getCurrentMonitorGeometry(window);
	int scale = getCurrentMonitorScaleFactor(window);
	SetMinMaxSize(window, 0, 0, m.width * scale, m.height * scale);

	gtk_window_fullscreen(window);

	return G_SOURCE_REMOVE;
}

gboolean UnFullscreen(gpointer data) {
	gtk_window_unfullscreen((GtkWindow*)data);

	return G_SOURCE_REMOVE;
}

bool disableContextMenu(GtkWindow* window) {
	return TRUE;
}

void DisableContextMenu(void* webview) {
	contextMenuDisabled = TRUE;
	g_signal_connect(WEBKIT_WEB_VIEW(webview), "context-menu", G_CALLBACK(disableContextMenu), NULL);
}


*/
import "C"
import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"strings"
	"unsafe"
)

func gtkBool(input bool) C.gboolean {
	if input {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

type Window struct {
	appoptions                               *options.App
	debug                                    bool
	gtkWindow                                unsafe.Pointer
	contentManager                           unsafe.Pointer
	webview                                  unsafe.Pointer
	applicationMenu                          *menu.Menu
	menubar                                  *C.GtkWidget
	vbox                                     *C.GtkWidget
	accels                                   *C.GtkAccelGroup
	minWidth, minHeight, maxWidth, maxHeight int
}

func bool2Cint(value bool) C.int {
	if value {
		return C.int(1)
	}
	return C.int(0)
}

func NewWindow(appoptions *options.App, debug bool) *Window {

	result := &Window{
		appoptions: appoptions,
		debug:      debug,
		minHeight:  appoptions.MinHeight,
		minWidth:   appoptions.MinWidth,
		maxHeight:  appoptions.MaxHeight,
		maxWidth:   appoptions.MaxWidth,
	}

	gtkWindow := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.g_object_ref_sink(C.gpointer(gtkWindow))
	result.gtkWindow = unsafe.Pointer(gtkWindow)

	result.vbox = C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
	C.gtk_container_add(result.asGTKContainer(), result.vbox)

	result.contentManager = unsafe.Pointer(C.webkit_user_content_manager_new())
	external := C.CString("external")
	defer C.free(unsafe.Pointer(external))
	C.webkit_user_content_manager_register_script_message_handler(result.cWebKitUserContentManager(), external)
	C.setupInvokeSignal(result.contentManager)

	webview := C.setupWebview(result.contentManager, result.asGTKWindow(), bool2Cint(appoptions.HideWindowOnClose))
	result.webview = unsafe.Pointer(webview)
	buttonPressedName := C.CString("button-press-event")
	defer C.free(unsafe.Pointer(buttonPressedName))
	C.connectButtons(unsafe.Pointer(webview))

	if debug {
		C.devtoolsEnabled(unsafe.Pointer(webview), C.int(1))
	} else {
		C.DisableContextMenu(unsafe.Pointer(webview))
	}

	// Set background colour
	RGBA := appoptions.RGBA
	result.SetRGBA(RGBA.R, RGBA.G, RGBA.B, RGBA.A)

	// Setup window
	result.SetKeepAbove(appoptions.AlwaysOnTop)
	result.SetResizable(!appoptions.DisableResize)
	result.SetSize(appoptions.Width, appoptions.Height)
	result.SetDecorated(!appoptions.Frameless)
	result.SetTitle(appoptions.Title)
	result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
	result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)

	// Menu
	result.SetApplicationMenu(appoptions.Menu)

	return result
}

func (w *Window) asGTKWidget() *C.GtkWidget {
	return C.GTKWIDGET(w.gtkWindow)
}

func (w *Window) asGTKWindow() *C.GtkWindow {
	return C.GTKWINDOW(w.gtkWindow)
}

func (w *Window) asGTKContainer() *C.GtkContainer {
	return C.GTKCONTAINER(w.gtkWindow)
}

func (w *Window) cWebKitUserContentManager() *C.WebKitUserContentManager {
	return (*C.WebKitUserContentManager)(w.contentManager)
}

func (w *Window) Fullscreen() {
	C.ExecuteOnMainThread(C.Fullscreen, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnFullscreen() {
	if !w.IsFullScreen() {
		return
	}
	C.ExecuteOnMainThread(C.UnFullscreen, C.gpointer(w.asGTKWindow()))
	w.SetMinSize(w.minWidth, w.minHeight)
	w.SetMaxSize(w.maxWidth, w.maxHeight)
}

func (w *Window) Destroy() {
	C.gtk_widget_destroy(w.asGTKWidget())
	C.g_object_unref(C.gpointer(w.gtkWindow))
}

func (w *Window) Close() {
	C.gtk_window_close(w.asGTKWindow())
}

func (w *Window) Center() {
	C.ExecuteOnMainThread(C.Center, C.gpointer(w.asGTKWindow()))
}

func (w *Window) SetPosition(x int, y int) {
	C.SetPosition(unsafe.Pointer(w.asGTKWindow()), C.int(x), C.int(y))
}

func (w *Window) Size() (int, int) {
	var width, height C.int
	C.gtk_window_get_size(w.asGTKWindow(), &width, &height)
	return int(width), int(height)
}

func (w *Window) GetPosition() (int, int) {
	var width, height C.int
	C.gtk_window_get_position(w.asGTKWindow(), &width, &height)
	return int(width), int(height)
}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {
	w.maxHeight = maxHeight
	w.maxWidth = maxWidth
	C.SetMinMaxSize(w.asGTKWindow(), C.int(w.minWidth), C.int(w.minHeight), C.int(w.maxWidth), C.int(w.maxHeight))
}

func (w *Window) SetMinSize(minWidth int, minHeight int) {
	w.minHeight = minHeight
	w.minWidth = minWidth
	C.SetMinMaxSize(w.asGTKWindow(), C.int(w.minWidth), C.int(w.minHeight), C.int(w.maxWidth), C.int(w.maxHeight))
}

func (w *Window) Show() {
	C.ExecuteOnMainThread(C.Show, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Hide() {
	C.ExecuteOnMainThread(C.Hide, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Maximise() {
	C.ExecuteOnMainThread(C.Maximise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnMaximise() {
	C.ExecuteOnMainThread(C.UnMaximise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) Minimise() {
	C.ExecuteOnMainThread(C.Minimise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) UnMinimise() {
	C.ExecuteOnMainThread(C.UnMinimise, C.gpointer(w.asGTKWindow()))
}

func (w *Window) IsFullScreen() bool {
	result := C.IsFullscreen(w.asGTKWidget())
	if result != 0 {
		return true
	}
	return false
}

func (w *Window) IsMaximised() bool {
	result := C.IsMaximised(w.asGTKWidget())
	return result > 0
}

func (w *Window) SetRGBA(r uint8, g uint8, b uint8, a uint8) {
	data := C.RGBAOptions{
		r:       C.uchar(r),
		g:       C.uchar(g),
		b:       C.uchar(b),
		a:       C.uchar(a),
		webview: w.webview,
	}
	C.ExecuteOnMainThread(C.setRGBA, C.gpointer(&data))

}

func (w *Window) Run() {
	if w.menubar != nil {
		C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), w.menubar, 0, 0, 0)
	}
	C.gtk_box_pack_start(C.GTKBOX(unsafe.Pointer(w.vbox)), C.GTKWIDGET(w.webview), 1, 1, 0)
	C.loadIndex(w.webview)
	C.gtk_widget_show_all(w.asGTKWidget())
	w.Center()
	switch w.appoptions.WindowStartState {
	case options.Fullscreen:
		w.Fullscreen()
	case options.Minimised:
		w.Minimise()
	case options.Maximised:
		w.Maximise()
	}

	C.gtk_main()
	w.Destroy()
}

func (w *Window) SetKeepAbove(top bool) {
	C.gtk_window_set_keep_above(w.asGTKWindow(), gtkBool(top))
}

func (w *Window) SetResizable(resizable bool) {
	C.gtk_window_set_resizable(w.asGTKWindow(), gtkBool(resizable))
}

func (w *Window) SetSize(width int, height int) {
	C.gtk_window_resize(w.asGTKWindow(), C.gint(width), C.gint(height))
}

func (w *Window) SetDecorated(frameless bool) {
	C.gtk_window_set_decorated(w.asGTKWindow(), gtkBool(frameless))
}

func (w *Window) SetTitle(title string) {
	C.SetTitle(w.asGTKWindow(), C.CString(title))
}

func (w *Window) ExecJS(js string) {
	jscallback := C.JSCallback{
		webview: w.webview,
		script:  C.CString(js),
	}
	C.ExecuteOnMainThread(C.executeJS, C.gpointer(&jscallback))
}

func (w *Window) StartDrag() {
	C.StartDrag(w.webview, w.asGTKWindow())
}

func (w *Window) Quit() {
	C.gtk_main_quit()
}

func (w *Window) OpenFileDialog(dialogOptions frontend.OpenDialogOptions, multipleFiles int, action C.GtkFileChooserAction) {

	data := C.OpenFileDialogOptions{
		webview:       w.webview,
		title:         C.CString(dialogOptions.Title),
		multipleFiles: C.int(multipleFiles),
		action:        action,
	}

	if len(dialogOptions.Filters) > 0 {
		// Create filter array
		mem := NewCalloc()
		arraySize := len(dialogOptions.Filters) + 1
		data.filters = C.allocFileFilterArray((C.ulong)(arraySize))
		filters := (*[1 << 30]*C.struct__GtkFileFilter)(unsafe.Pointer(data.filters))
		for index, filter := range dialogOptions.Filters {
			thisFilter := C.gtk_file_filter_new()
			C.g_object_ref(C.gpointer(thisFilter))
			if filter.DisplayName != "" {
				cName := mem.String(filter.DisplayName)
				C.gtk_file_filter_set_name(thisFilter, cName)
			}
			if filter.Pattern != "" {
				for _, thisPattern := range strings.Split(filter.Pattern, ";") {
					cThisPattern := mem.String(thisPattern)
					C.gtk_file_filter_add_pattern(thisFilter, cThisPattern)
				}
			}
			// Add filter to array
			filters[index] = thisFilter
		}
		mem.Free()
		filters[arraySize-1] = nil
	}

	if dialogOptions.CanCreateDirectories {
		data.createDirectories = C.int(1)
	}

	if dialogOptions.ShowHiddenFiles {
		data.showHiddenFiles = C.int(1)
	}

	if dialogOptions.DefaultFilename != "" {
		data.defaultFilename = C.CString(dialogOptions.DefaultFilename)
	}

	if dialogOptions.DefaultDirectory != "" {
		data.defaultDirectory = C.CString(dialogOptions.DefaultDirectory)
	}

	C.ExecuteOnMainThread(C.opendialog, C.gpointer(&data))
}

func (w *Window) MessageDialog(dialogOptions frontend.MessageDialogOptions) {

	data := C.MessageDialogOptions{
		window:  w.gtkWindow,
		title:   C.CString(dialogOptions.Title),
		message: C.CString(dialogOptions.Message),
	}
	switch dialogOptions.Type {
	case frontend.InfoDialog:
		data.messageType = C.int(0)
	case frontend.ErrorDialog:
		data.messageType = C.int(1)
	case frontend.QuestionDialog:
		data.messageType = C.int(2)
	case frontend.WarningDialog:
		data.messageType = C.int(3)
	}
	C.ExecuteOnMainThread(C.messageDialog, C.gpointer(&data))
}

func (w *Window) ToggleMaximise() {
	if w.IsMaximised() {
		w.UnMaximise()
	} else {
		w.Maximise()
	}
}
