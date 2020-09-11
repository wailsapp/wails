
#ifdef FFENESTRI_DARWIN

#define OBJC_OLD_DISPATCH_PROTOTYPES 1
#include <objc/objc-runtime.h>
#include <CoreGraphics/CoreGraphics.h>

// Macros to make it slightly more sane
#define msg objc_msgSend

#define c(str) (id)objc_getClass(str)
#define s(str) sel_registerName(str)
#define str(input) msg(c("NSString"), s("stringWithUTF8String:"), input)

#define NSBackingStoreBuffered 2

#define NSBorderlessWindowMask 0
#define NSTitledWindowMask 1
#define NSClosableWindowMask 2
#define NSMiniaturizableWindowMask 4
#define NSResizableWindowMask 8

// References to assets
extern const unsigned char *assets[];
extern const unsigned char runtime;
extern const char *icon[];

// MAIN DEBUG FLAG
int debug;

// App Delegate
typedef struct AppDel {
	Class isa;
	id window;
} AppDelegate;

// Credit: https://stackoverflow.com/a/8465083
char* concat(const char *s1, const char *s2)
{
    const size_t len1 = strlen(s1);
    const size_t len2 = strlen(s2);
    char *result = malloc(len1 + len2 + 1);
    memcpy(result, s1, len1);
    memcpy(result + len1, s2, len2 + 1);
    return result;
}

// yes command simply returns YES!
BOOL yes(id self, SEL cmd)
{
    return YES;
}


// Debug works like sprintf but mutes if the global debug flag is true
// Credit: https://stackoverflow.com/a/20639708
void Debug(char *message, ... ) {
    if ( debug ) {
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


void messageHandler(id self, SEL cmd, id contentController, id msg) {
    // TODO: Pass message to callback in application object
    Debug("TODO: Process message");
}

extern void messageFromWindowCallback(const char *);
typedef void (*ffenestriCallback)(const char *);

struct Application {

    // Cocoa data
    id mainWindow;
    id wkwebview;

    // Window Data
    const char *title;
    int width;
    int height;
    int resizable;
    int devtools;
    int fullscreen;

    // Features
    int borderless;

    // User Data
    char *HTML;

    // Callback
    ffenestriCallback sendMessageToBackend;

    // Bindings
    const char *bindings;

    // Lock - used for sync operations (Should we be using g_mutex?)
    int lock;

};



void* NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen) {
    // Setup main application struct
    struct Application *result = malloc(sizeof(struct Application));
    result->title = title;
    result->width = width;
    result->height = height;
    result->resizable = resizable;
    result->devtools = devtools;
    result->fullscreen = fullscreen;
    result->lock = 0;

    // Features
    result->borderless = 0;

    result->sendMessageToBackend = (ffenestriCallback) messageFromWindowCallback;

    return (void*) result;
}

void DestroyApplication(void *appPointer) {
    Debug("Destroying Application");
    struct Application *app = (struct Application*) appPointer;

    // Free the bindings
    if (app->bindings != NULL) {
        free((void*)app->bindings);
        app->bindings = NULL;
    } else {
        Debug("Almost a double free for app->bindings");
    }

    msg(app->mainWindow, s("close"));
    msg(c("NSApp"), s("terminate:"), NULL);
    Debug("Finished Destroying Application");
}

// Quit will stop the gtk application and free up all the memory
// used by the application
void Quit(void *appPointer) {
    Debug("Quit Called");
    // struct Application *app = (struct Application*) appPointer;
    // gtk_window_close((GtkWindow*)app->mainWindow);
    // g_application_quit((GApplication*)app->application);
    DestroyApplication(appPointer);
}

// SetTitle sets the main window title to the given string
void SetTitle(void *appPointer, const char *title) {
    Debug("SetTitle Called");
    struct Application *app = (struct Application*) appPointer;
    msg(app->mainWindow, s("setTitle:"), msg(c("NSString"), s("stringWithUTF8String:"), title));
}

// Fullscreen sets the main window to be fullscreen
void Fullscreen(void *appPointer) {
    Debug("Fullscreen Called");
    // struct Application *app = (struct Application*) appPointer;
    // gtk_window_fullscreen (app->mainWindow);
}

// UnFullscreen resets the main window after a fullscreen
void UnFullscreen(void *appPointer) {
    Debug("UnFullscreen Called");
    // struct Application *app = (struct Application*) appPointer;
    // gtk_window_unfullscreen (app->mainWindow);
}

void Center(void *appPointer) {
    Debug("Center Called");
    // struct Application *app = (struct Application*) appPointer;
    // gtk_window_set_position(app->mainWindow, GTK_WIN_POS_CENTER);
}

void SetMaximumSize(void *appPointer, int width, int height) {
    Debug("SetMaximumSize Called");
    // struct Application *app = (struct Application*) appPointer;
    // GdkGeometry size;
    // size.max_height = (height == 0 ? INT_MAX: height);
    // size.max_width = (width == 0 ? INT_MAX: width);
    // gtk_window_set_geometry_hints(app->mainWindow, NULL, &size, GDK_HINT_MAX_SIZE);
}

void SetMinimumSize(void *appPointer, int width, int height) {
    Debug("SetMinimumSize Called");
    // struct Application *app = (struct Application*) appPointer;
    // GdkGeometry size;
    // size.max_height = height;
    // size.max_width = width;
    // gtk_window_set_geometry_hints(app->mainWindow, NULL, &size, GDK_HINT_MIN_SIZE);
}

// OpenFileDialog opens a dialog to select a file
// NOTE: The result is a string that will need to be freed!
char* OpenFileDialog(void *appPointer, char *title) {
    Debug("OpenFileDialog Called");

    // struct Application *app = (struct Application*) appPointer;
    // GtkFileChooserNative *native;
    // GtkFileChooserAction action = GTK_FILE_CHOOSER_ACTION_OPEN;
    // gint res;
    char *filename = "BogusFilename";

    // native = gtk_file_chooser_native_new (title,
    //                                     app->mainWindow,
    //                                     action,
    //                                     "_Open",
    //                                     "_Cancel");

    // res = gtk_native_dialog_run (GTK_NATIVE_DIALOG (native));
    // if (res == GTK_RESPONSE_ACCEPT)
    // {
    //     GtkFileChooser *chooser = GTK_FILE_CHOOSER (native);
    //     filename = gtk_file_chooser_get_filename (chooser);
    // }

    // g_object_unref (native);

    return filename;
}

// SaveFileDialog opens a dialog to select a file
// NOTE: The result is a string that will need to be freed!
char* SaveFileDialog(void *appPointer, char *title) {
    Debug("SaveFileDialog Called");
    char *filename = "BogusSaveFilename";
/*    struct Application *app = (struct Application*) appPointer;
    GtkFileChooserNative *native;
    GtkFileChooserAction action = GTK_FILE_CHOOSER_ACTION_SAVE;
    gint res;

    native = gtk_file_chooser_native_new (title,
                                        app->mainWindow,
                                        action,
                                        "_Save",
                                        "_Cancel");

    res = gtk_native_dialog_run (GTK_NATIVE_DIALOG (native));
    if (res == GTK_RESPONSE_ACCEPT)
    {
        GtkFileChooser *chooser = GTK_FILE_CHOOSER (native);
        filename = gtk_file_chooser_get_filename (chooser);
    }

    g_object_unref (native);
*/
    return filename;
}

// OpenDirectoryDialog opens a dialog to select a directory
// NOTE: The result is a string that will need to be freed!
char* OpenDirectoryDialog(void *appPointer, char *title) {
    Debug("OpenDirectoryDialog Called");
    char *foldername = "BogusDirectory";
/*
    struct Application *app = (struct Application*) appPointer;
    GtkFileChooserNative *native;
    GtkFileChooserAction action = GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER;
    gint res;

    native = gtk_file_chooser_native_new (title,
                                        app->mainWindow,
                                        action,
                                        "_Open",
                                        "_Cancel");

    res = gtk_native_dialog_run (GTK_NATIVE_DIALOG (native));
    if (res == GTK_RESPONSE_ACCEPT)
    {
        GtkFileChooser *chooser = GTK_FILE_CHOOSER (native);
        foldername = gtk_file_chooser_get_filename (chooser);
    }

    g_object_unref (native);
*/
    return foldername;
}

  // SetColour sets the colour of the webview to the given colour string
int SetColour(void *appPointer, const char *colourString) {
    Debug("SetColour Called");

    // struct Application *app = (struct Application*) appPointer;
    // GdkRGBA rgba;
    // gboolean result = gdk_rgba_parse(&rgba, colourString);
    // if (result == FALSE) {
    //     return 0;
    // }
    // Debug("Setting webview colour to: %s", colourString);
    // webkit_web_view_get_background_color((WebKitWebView*)(app->webView), &rgba);
    return 1;
}

// WindowBorderless will ensure that the main window gets created
// without a border
void WindowBorderless(void *appPointer) {
    Debug("WindowBoderless Called");
    
    struct Application *app = (struct Application*) appPointer;
    app->borderless = 1;
}
const char *invoke = "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}};";

// void syncCallback(id self, NSError cmd ) {
//     Debug("In syncCallback");
//     struct Application *app = (struct Application *)objc_getAssociatedObject(self, "application");
//     Debug("app: %p", app);
//     app->lock = 0;
// }
/*
void syncEval( struct Application *app, const char *script) {

    id wkwebview = app->wkwebview;
    Debug("[%p] wkwebview", wkwebview);
    Debug("[%p] Running sync", script);

    // wait for lock to free
    while (app->lock == 1) {
        Debug("Waiting for current lock to free");
        usleep(100000);
    }
    // Set lock
    app->lock = 1;

    Debug("Executing JS");

    msg(wkwebview, s("evaluateJavaScript:completionHandler:"),
    msg(c("NSString"), s("stringWithUTF8String:"), script), ^(id result, int *error) {
        Debug("Here");
        app->lock = 0;
    });

    while (app->lock == 1) {
        Debug("[%p] Waiting for callback", script);
        usleep(100000);
	}
    // Debug("[%p] Finished\n", script);
}

void asyncEval(  WebKitWebView  *webView, const gchar *script ) {
        webkit_web_view_run_javascript(
            webView,
            script,
            NULL, NULL, NULL);
}

typedef void (*dispatchMethod)(void *app, void *);

struct dispatchData {
    struct Application *app;
    dispatchMethod method;
    void *args;
};


gboolean executeMethod(gpointer data) {
    struct dispatchData *d = (struct dispatchData *)data;
    struct Application *app = (struct Application *)(d->app);
    // Debug("Webview %p\n", app->webView);
    // Debug("Args %s\n", d->args);
    // Debug("Method %p\n", (d->method));
    (d->method)(app, d->args);
    // Debug("Method Execute Complete. Freeing memory");
    g_free(d);
    return FALSE;
}
*/

// DisableFrame disables the window frame
void DisableFrame(void *appPointer)
{
   // TBD
}

void SetMinWindowSize(void *appPointer, int minWidth, int minHeight)
{
    // TBD
}

void SetMaxWindowSize(void *appPointer, int maxWidth, int maxHeight)
{
    // TBD
}

void ExecJS(void *appPointer, char *js) {
    Debug("ExecJS Called: %s", js);
    struct Application *app = (struct Application*) appPointer;
    Debug("App: %p", app);
    Debug("wkwebview: %p", app->wkwebview);
    msg(app->wkwebview, s("evaluateJavaScript:completionHandler:"),
        msg(c("NSString"), s("stringWithUTF8String:"), js), NULL);
}

// typedef char* (*dialogMethod)(void *app, void *);

// struct dialogCall {
//     struct Application *app;
//     dialogMethod method;
//     void *args;
//     char *result;
//     int done;
// };


// gboolean executeMethodWithReturn(gpointer data) {
//     struct dialogCall *d = (struct dialogCall *)data;
//     struct Application *app = (struct Application *)(d->app);
//     Debug("Webview %p\n", app->webView);
//     Debug("Args %s\n", d->args);
//     Debug("Method %p\n", (d->method));
//     d->result = (d->method)(app, d->args);
//     d->done = 1;
//     // Debug("Method Execute Complete. Freeing memory");
//     return FALSE;
// }


char* OpenFileDialogOnMainThread(void *app, char *title) {
    Debug("OpenFileDialogOnMainThread Called");

    //   struct dialogCall *data =
    //   (struct dialogCall *)g_new(struct dialogCall, 1);
    //   data->result = NULL;
    //   data->done = 0;
    //   data->method = (dialogMethod)OpenFileDialog;
    //   data->args = title;
    //   data->app = app;

    // gdk_threads_add_idle(executeMethodWithReturn, data);

    // while( data->done == 0 ) {
    //     // Debug("Waiting for dialog");
    //     usleep(100000);
    // }
    // Debug("Dialog done");
    // Debug("Result = %s\n", data->result);
    // g_free(data);
    // // Fingers crossed this wasn't freed by g_free above
    // return data->result;
    return "OpenFileDialogOnMainThread result";
}

char* SaveFileDialogOnMainThread(void *app, char *title) {
    Debug("SaveFileDialogOnMainThread Called");
    return "SaveFileDialogOnMainThread result";
}

char* OpenDirectoryDialogOnMainThread(void *app, char *title) {
    Debug("OpenDirectoryDialogOnMainThread Called");
    return "OpenDirectoryDialogOnMainThread result";
}

// // Sets the icon to the XPM stored in icon
// void setIcon( struct Application *app) {
//     GdkPixbuf *appIcon = gdk_pixbuf_new_from_xpm_data ((const char**)icon);
//     gtk_window_set_icon (app->mainWindow, appIcon);
// }


// static void load_finished_cb (  WebKitWebView  *webView,
//                                 WebKitLoadEvent load_event,
//                                 gpointer        userData)
// {
//     struct Application* app;

//     switch (load_event) {
//     // case WEBKIT_LOAD_STARTED:
//     //     /* New load, we have now a provisional URI */
//     //     // printf("Start downloading %s\n", webkit_web_view_get_uri(web_view));
//     //     /* Here we could start a spinner or update the
//     //      * location bar with the provisional URI */
//     //     break;
//     // case WEBKIT_LOAD_REDIRECTED:
//     //     // printf("Redirected to: %s\n", webkit_web_view_get_uri(web_view));
//     //     break;
//     // case WEBKIT_LOAD_COMMITTED:
//     //     /* The load is being performed. Current URI is
//     //      * the final one and it won't change unless a new
//     //      * load is requested or a navigation within the
//     //      * same page is performed */
//     //     // printf("Loading: %s\n", webkit_web_view_get_uri(web_view));
//     //     break;
//     case WEBKIT_LOAD_FINISHED:
//         /* Load finished, we can now stop the spinner */
//         // printf("Finished loading: %s\n", webkit_web_view_get_uri(web_view));
//         app = (struct Application*) userData;

//         Debug("Initialising Wails Window");
//         syncEval(app, invoke);

//         // Bindings
//         Debug("Binding Methods");
//         syncEval(app, app->bindings);

//         // Runtime
//         Debug("Setting up Wails runtime");
//         syncEval(app, &runtime);

//         // Loop over assets
//         int index = 0;
//         while(1) {
//             // Get next asset pointer
//             const char *asset = assets[index];

//             // If we have no more assets, break
//             if (asset == 0x00) {
//                 break;
//             }

//             // sync eval the asset
//             syncEval(app, asset);
//             index++;
//         };

//         // Set the icon
//         setIcon(app);

//         // Setup fullscreen
//         if( app->fullscreen ) {
//             Debug("Going fullscreen");
//             Fullscreen(app);
//         }

//         // Setup resize
//         gtk_window_resize(GTK_WINDOW (app->mainWindow), app->width, app->height);

//         if( app->resizable ) {
//             gtk_window_set_default_size(GTK_WINDOW (app->mainWindow), app->width, app->height);
//         } else {
//             gtk_widget_set_size_request(GTK_WIDGET (app->mainWindow), app->width, app->height);
//             SetMaximumSize(app, app->width, app->height);
//             SetMinimumSize(app, app->width, app->height);
//             gtk_window_resize(GTK_WINDOW (app->mainWindow), app->width, app->height);
//         }
//         gtk_window_set_resizable(GTK_WINDOW(app->mainWindow), app->resizable ? TRUE : FALSE );

//         // Centre by default
//         Center(app);

//         // Show window and focus
//         gtk_widget_show_all(GTK_WIDGET(app->mainWindow));
//         gtk_widget_grab_focus(app->webView);
//         break;
//     }
// }

// static
// gboolean
// disable_context_menu_cb(WebKitWebView       *web_view,
//                         WebKitContextMenu   *context_menu,
//                         GdkEvent            *event,
//                         WebKitHitTestResult *hit_test_result,
//                         gpointer             user_data) {
//   return TRUE;
// }

// static void sendMessageToBackend(   WebKitUserContentManager *contentManager,
//                                          WebKitJavascriptResult *result,
//                                          gpointer arg) {
//   struct Application *app = (struct Application *)arg;
// #if WEBKIT_MAJOR_VERSION >= 2 && WEBKIT_MINOR_VERSION >= 22
//   JSCValue *value = webkit_javascript_result_get_js_value(result);
//   char *message = jsc_value_to_string(value);
// #else
//   JSGlobalContextRef context = webkit_javascript_result_get_global_context(result);
//   JSValueRef value = webkit_javascript_result_get_value(result);
//   JSStringRef js = JSValueToStringCopy(context, value, NULL);
//   size_t messageSize = JSStringGetMaximumUTF8CStringSize(js);
//   char *message = g_new(char, messageSize);
//   JSStringGetUTF8CString(js, message, messageSize);
//   JSStringRelease(js);
// #endif
//   app->sendMessageToBackend(message);
//   g_free(message);
// }

void SetDebug(void *applicationPointer, int flag) {
    struct Application *app = (struct Application*) applicationPointer;
    debug = flag;
}

void SetBindings(void* applicationPointer, const char *bindings) {
    struct Application *app = (struct Application*) applicationPointer;

    const char* temp = concat("window.wailsbindings = \"", bindings);
    const char* jscall = concat(temp, "\";");
    free((void*)temp);
    app->bindings = jscall;
}

// This is called when the close button on the window is pressed
// void close_button_pressed () {
//     struct Application *app = (struct Application*) user_data;
//     app->sendMessageToBackend("WC"); // Window Close
//     return TRUE;
// }

// static void setupWindow(void *applicationPointer) {

//     struct Application *app = (struct Application*) applicationPointer;

//     // Create the window
//     GtkWidget *mainWindow = gtk_application_window_new (app->application);
//     // Save reference
//     app->mainWindow = GTK_WINDOW(mainWindow);

//     // Setup borderless
//     if (app->borderless) {
//         gtk_window_set_decorated((GtkWindow*)mainWindow, FALSE);
//     }

//     // Setup title
//     printf("Setting title to: %s\n", app->title);
//     gtk_window_set_title(GTK_WINDOW(mainWindow), app->title);

//     // Setup script handler
//     WebKitUserContentManager *contentManager = webkit_user_content_manager_new();
//     webkit_user_content_manager_register_script_message_handler(contentManager, "external");
//     g_signal_connect(contentManager, "script-message-received::external", G_CALLBACK(sendMessageToBackend), app);
//     GtkWidget *webView = webkit_web_view_new_with_user_content_manager(contentManager);
//     // Save reference
//     app->webView = webView;

//     // Add the webview to the window
//     gtk_container_add(GTK_CONTAINER(mainWindow), webView);

//     // Load default HTML
//     g_signal_connect(G_OBJECT(webView), "load-changed", G_CALLBACK(load_finished_cb), app);

//     // Load the user's HTML
//     webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webView), &userhtml);

//     // Check if we want to enable the dev tools
//     if( app->devtools ) {
//         WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webView));
//         // webkit_settings_set_enable_write_console_messages_to_stdout(settings, true);
//         webkit_settings_set_enable_developer_extras(settings, true);
//     } else {
//         g_signal_connect(G_OBJECT(webView), "context-menu", G_CALLBACK(disable_context_menu_cb), app);
//     }

//     // Listen for close button signal
//     g_signal_connect (GTK_WIDGET(mainWindow), "delete-event", G_CALLBACK (close_button_pressed), app);

// }


// static void
// activate (GtkApplication* app,
//           gpointer        user_data)
// {
//     struct Application *mainApp = (struct Application*) user_data;

//     // Main Window
//     setupWindow(mainApp);

// }

// void dispatch(void *func) {
//     dispatch_fn_t *f = malloc()
// dispatch_async_f(dispatch_get_main_queue(), new dispatch_fn_t(func),
//                     (dispatch_function_t)((void *arg) {
//                     (*arg)();
//                     }));
// }


void Run(void *applicationPointer, int argc, char **argv) {
    struct Application *app = (struct Application*) applicationPointer;

    int decorations = NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask;
    if (app->resizable) {
        decorations |= NSResizableWindowMask;
    }

    id application = msg(c("NSApplication"), s("sharedApplication"));
    msg(application, s("setActivationPolicy:"), 0);

    // Define delegate
    Class delegateClass = objc_allocateClassPair((Class) c("NSResponder"), "AppDelegate", 0);
    class_addProtocol(delegateClass, objc_getProtocol("NSTouchBarProvider"));
    class_addMethod(delegateClass, s("applicationShouldTerminateAfterLastWindowClosed:"), (IMP) yes, "c@:@");
    // TODO: add userContentController:didReceiveScriptMessage
    class_addMethod(delegateClass, s("userContentController:didReceiveScriptMessage:"), (IMP) messageHandler,
                    "v@:@@");
    objc_registerClassPair(delegateClass);

    // Create delegate
    id delegate = msg((id)delegateClass, s("new"));
    objc_setAssociatedObject(delegate, "application", (id)app, OBJC_ASSOCIATION_ASSIGN);
    msg(application, sel_registerName("setDelegate:"), delegate);

    // Create main window
    id mainWindow = msg(c("NSWindow"),s("alloc"));
    mainWindow = msg(mainWindow, s("initWithContentRect:styleMask:backing:defer:"),
          CGRectMake(0, 0, app->width, app->height), decorations, NSBackingStoreBuffered, 0);
    app->mainWindow = mainWindow;

    // Set the main window title
    SetTitle(app, app->title);

    // Setup webview
    id config = msg(c("WKWebViewConfiguration"), s("new"));
    id manager = msg(config, s("userContentController"));
    id wkwebview = msg(c("WKWebView"), s("alloc"));
    app->wkwebview = wkwebview;

    // TODO: Fix "NSWindow warning: adding an unknown subview: <WKInspectorWKWebView: 0x465ed90>. Break on NSLog to debug." error
    if (app->devtools) {
      Debug("Enabling devtools");
      msg(msg(config, s("preferences")), s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 1), str("developerExtrasEnabled"));
    }
    // TODO: Understand why this shouldn't be CGRectMake(0, 0, app->width, app->height)
    msg(wkwebview, s("initWithFrame:configuration:"), CGRectMake(0, 0, 0, 0), config);
    
    msg(manager, s("addScriptMessageHandler:name:"), delegate, str("external"));
    msg(mainWindow, s("setContentView:"), wkwebview);
    msg(mainWindow, s("makeKeyAndOrderFront:"), NULL);

    // Load HTML
    id html = msg(c("NSURL"), s("URLWithString:"), msg(c("NSString"), s("stringWithUTF8String:"), assets[0]));
    Debug("HTML: %p", html);
    msg(wkwebview, s("loadRequest:"), msg(c("NSURLRequest"), s("requestWithURL:"), html));

    // Load assets
    
    Debug("Loading Internal Code");
    // We want to evaluate the internal code plus runtime before the assets
    const char *temp = concat(invoke, app->bindings);
    const char *internalCode = concat(temp, (const char*)&runtime);
    // Debug("Internal code: %s", internalCode);
    free((void*)temp);

    // Loop over assets
    int index = 1;
    while(1) {
        // Get next asset pointer
        const unsigned char *asset = assets[index];

        // If we have no more assets, break
        if (asset == 0x00) {
            break;
        }

        temp = concat(internalCode, asset);
        free((void*)internalCode);
        internalCode = temp;
        index++;
    };

    Debug("MOAE: %s", internalCode);

    // Evaluate internal code
    // ---------------------------------------------------------------
    // TODO: WORK OUT WHY EXECUTING THE INTERNAL CODE BLOWS THE APP UP
    // ---------------------------------------------------------------
    msg(wkwebview, 
        s("evaluateJavaScript:completionHandler:"),
        msg(c("NSString"), s("stringWithUTF8String:"), (char*)internalCode), 
        ^() {

            // Debug("Result: %p", result);
            // Debug("Error: %p", error);
            // Free internalCode memory
            // free((void*)internalCode);
/*
            if( error != NULL) {
                Debug("Error message: %s", error);
                return;
            }
*/
        Debug("DO I CARE?");
 

        }
    );

    // Finally call run
    msg(application, s("run"));

    Debug("Run called");
}



#endif
