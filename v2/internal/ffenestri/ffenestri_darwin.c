
#ifdef FFENESTRI_DARWIN

#define OBJC_OLD_DISPATCH_PROTOTYPES 1
#include <objc/objc-runtime.h>
#include <CoreGraphics/CoreGraphics.h>
#include "json.h"

// Macros to make it slightly more sane
#define msg objc_msgSend

#define c(str) (id)objc_getClass(str)
#define s(str) sel_registerName(str)
#define u(str) sel_getUid(str)
#define str(input) msg(c("NSString"), s("stringWithUTF8String:"), input)
#define url(input) msg(c("NSURL"), s("fileURLWithPath:"), str(input))

#define GET_FRAME(receiver) ((CGRect(*)(id, SEL))objc_msgSend_stret)(receiver, s("frame"));

#define ON_MAIN_THREAD(str) dispatch( ^{ str; } );
#define MAIN_WINDOW_CALL(str) msg(app->mainWindow, s((str)));

#define NSBackingStoreBuffered 2

#define NSWindowStyleMaskBorderless 0
#define NSWindowStyleMaskTitled 1
#define NSWindowStyleMaskClosable 2
#define NSWindowStyleMaskMiniaturizable 4
#define NSWindowStyleMaskResizable 8
#define NSWindowStyleMaskFullscreen 1 << 14

#define NSWindowTitleHidden 1
#define NSWindowStyleMaskFullSizeContentView 1 << 15

// References to assets
extern const unsigned char *assets[];
extern const unsigned char runtime;
extern const char *icon[];

// MAIN DEBUG FLAG
int debug;

// Dispatch Method
typedef void (^dispatchMethod)(void);

// dispatch will execute the given `func` pointer
void dispatch(dispatchMethod func) {
    dispatch_async(dispatch_get_main_queue(), func);
}

// App Delegate
typedef struct AppDel {
	Class isa;
	id window;
} AppDelegate;

// Credit: https://stackoverflow.com/a/8465083
char* concat(const char *string1, const char *string2)
{
    const size_t len1 = strlen(string1);
    const size_t len2 = strlen(string2);
    char *result = malloc(len1 + len2 + 1);
    memcpy(result, string1, len1);
    memcpy(result + len1, string2, len2 + 1);
    return result;
}

// yes command simply returns YES!
BOOL yes(id self, SEL cmd)
{
    return YES;
}

// Debug works like sprintf but mutes if the global debug flag is true
// Credit: https://stackoverflow.com/a/20639708
void Debug(const char *message, ... ) {
    if ( debug ) {
        char *temp = concat("TRACE | Ffenestri (C) | ", message);
        message = concat(temp, "\n");
        free(temp);
        va_list args;
        va_start(args, message);
        vprintf(message, args);
        free((void*)message);
        va_end(args);
    }
}

extern void messageFromWindowCallback(const char *);
typedef void (*ffenestriCallback)(const char *);

struct Application {

    // Cocoa data
    id application;
    id mainWindow;
    id wkwebview;
    id manager;
    id config;
    id mouseEvent;

    // Window Data
    const char *title;
    int width;
    int height;
    int minWidth;
    int minHeight;
    int maxWidth;
    int maxHeight;
    int resizable;
    int devtools;
    int fullscreen;
    int red;
    int green;
    int blue;
    int alpha;

    // Features
    int frame;
    int maximised;
    int minimised;
    int titlebarAppearsTransparent;
    int hideTitle;
    int hideTitleBar;
    int fullSizeContent;
    int useToolBar;
    int hideToolbarSeparator;

    // User Data
    char *HTML;

    // Callback
    ffenestriCallback sendMessageToBackend;

    // Bindings
    const char *bindings;

    // Lock - used for sync operations (Should we be using g_mutex?)
    int lock;

};

void TitlebarAppearsTransparent(struct Application* app) {
    app->titlebarAppearsTransparent = 1;
    Debug("[x] setTitlebarAppearsTransparent %d", app->titlebarAppearsTransparent? YES:NO);
}

void HideTitle(struct Application *app) {
    app->hideTitle = 1;
}

void HideTitleBar(struct Application *app) {
    app->hideTitleBar = 1;
}

void HideToolbarSeparator(struct Application *app) {
    app->hideToolbarSeparator = 1;
}

void UseToolbar(struct Application *app) {
    app->useToolBar = 1;
}

void applyWindowColour(struct Application *app) {
    // Apply the colour only if the window has been created
    if( app->mainWindow != NULL ) {
        ON_MAIN_THREAD(
            id colour = msg(c("NSColor"), s("colorWithCalibratedRed:green:blue:alpha:"),
                                (float)app->red / 255.0, 
                                (float)app->green / 255.0, 
                                (float)app->blue / 255.0,
                                (float)app->alpha / 255.0);
            msg(app->mainWindow, s("setBackgroundColor:"), colour);
        )
    } 
} 

void SetColour(struct Application *app, int red, int green, int blue, int alpha) {
    app->red = red;
    app->green = green;
    app->blue = blue;
    app->alpha = alpha;

    applyWindowColour(app);
}

void FullSizeContent(struct Application *app) {
    app->fullSizeContent = 1;
}

void Hide(struct Application *app) { 
    ON_MAIN_THREAD( 
        msg(app->application, s("hide:")) 
    )
}

void Show(struct Application *app) { 
    ON_MAIN_THREAD(
        msg(app->mainWindow, s("makeKeyAndOrderFront:"), NULL);
        msg(app->application, s("activateIgnoringOtherApps:"), YES);
    )
}

// Sends messages to the backend
void messageHandler(id self, SEL cmd, id contentController, id message) {
    struct Application *app = (struct Application *)objc_getAssociatedObject(
                              self, "application");
    const char *name = (const char *)msg(msg(message, s("name")), s("UTF8String"));
    if( strcmp(name, "completed") == 0) {
        // Delete handler
        msg(app->manager, s("removeScriptMessageHandlerForName:"), str("completed"));
        // Show window after a short delay so rendering can catch up
        dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 10000000), dispatch_get_main_queue(), ^{
            Show(app); 
            msg(app->config, s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 0), str("suppressesIncrementalRendering"));
        });
    } else if( strcmp(name, "windowDrag") == 0 ) {
        Debug("DRAGGING!!!!!");
        msg(app->mainWindow, s("performWindowDragWithEvent:"), app->mouseEvent);
    } else {
        const char *m = (const char *)msg(msg(message, s("body")), s("UTF8String"));
        app->sendMessageToBackend(m);
    }
}

// closeWindow is called when the close button is pressed
void closeWindow(id self, SEL cmd, id sender) {
    struct Application *app = (struct Application *) objc_getAssociatedObject(self, "application");
    app->sendMessageToBackend("WC");
}

void* NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen) {
    // Setup main application struct
    struct Application *result = malloc(sizeof(struct Application));
    result->title = title;
    result->width = width;
    result->height = height;
    result->minWidth = 0;
    result->minHeight = 0;
    result->maxWidth = 0;
    result->maxHeight = 0;
    result->resizable = resizable;
    result->devtools = devtools;
    result->fullscreen = fullscreen;
    result->lock = 0;
    result->maximised = 0;
    result->minimised = 0;

    result->mainWindow = NULL;

    // Features
    result->frame = 1;
    result->hideTitle = 0;
    result->hideTitleBar = 0;
    result->fullSizeContent = 0;
    result->useToolBar = 0;
    result->hideToolbarSeparator = 0;

    result->titlebarAppearsTransparent = 0;

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

    msg(app->manager, s("removeScriptMessageHandlerForName:"), str("external"));
    msg(app->mainWindow, s("close"));
    msg(c("NSApp"), s("terminate:"), NULL);
    Debug("Finished Destroying Application");
}

// Quit will stop the gtk application and free up all the memory
// used by the application
void Quit(void *appPointer) {
    Debug("Quit Called");
    DestroyApplication(appPointer);
}

// SetTitle sets the main window title to the given string
void SetTitle(struct Application *app, const char *title) {
    Debug("SetTitle Called");
    ON_MAIN_THREAD(
        msg(app->mainWindow, s("setTitle:"), str(title));
    )
}

void ToggleFullscreen(struct Application *app) {
    ON_MAIN_THREAD(
        app->fullscreen = !app->fullscreen;
        MAIN_WINDOW_CALL("toggleFullScreen:")
    )
}

// Fullscreen sets the main window to be fullscreen
void Fullscreen(struct Application *app) {
    Debug("Fullscreen Called");
    if( app->fullscreen == 0) {
        ToggleFullscreen(app);
    }
}

// UnFullscreen resets the main window after a fullscreen
void UnFullscreen(struct Application *app) {
    Debug("UnFullscreen Called");
    if( app->fullscreen == 1) {
        ToggleFullscreen(app);
    }
}

void Center(struct Application *app) {
    Debug("Center Called");
    ON_MAIN_THREAD(
        MAIN_WINDOW_CALL("center")
    )
}

void ToggleMaximise(struct Application *app) {
    ON_MAIN_THREAD(
        app->maximised = !app->maximised;
        MAIN_WINDOW_CALL("zoom:")
    )
}

void Maximise(struct Application *app) { 
    if( app->maximised == 0) {
        ToggleMaximise(app);
    }
}

void Unmaximise(struct Application *app) { 
    if( app->maximised == 1) {
        ToggleMaximise(app);
    }
}

void ToggleMinimise(struct Application *app) {
    ON_MAIN_THREAD(
        MAIN_WINDOW_CALL(app->minimised ? "deminiaturize:" : "miniaturize:" );
        app->minimised = !app->minimised;
    )
}

void Minimise(struct Application *app) {
    if( app->minimised == 0) {
        ToggleMinimise(app);
    }
 }
void Unminimise(struct Application *app) {
    if( app->minimised == 1) {
        ToggleMinimise(app);
    }
 }

id getCurrentScreen(struct Application *app) {
    id screen = NULL;
    screen = msg(app->mainWindow, s("screen"));
    if( screen == NULL ) {
        screen = msg(c("NSScreen"), u("mainScreen"));
    }
    return screen;
}

void dumpFrame(const char *message, CGRect frame) {
    Debug(message);
    Debug("origin.x %f", frame.origin.x);
    Debug("origin.y %f", frame.origin.y);        
    Debug("size.width %f", frame.size.width);
    Debug("size.height %f", frame.size.height);
}

void SetSize(struct Application *app, int width, int height) { 
    ON_MAIN_THREAD(
        id screen = getCurrentScreen(app);

        // Get the rect for the window
        CGRect frame = GET_FRAME(app->mainWindow);

        // Credit: https://github.com/patr0nus/DeskGap/blob/73c0ac9f2c73f55b6e81f64f6673a7962b5719cd/lib/src/platform/mac/util/NSScreen%2BGeometry.m
        frame.origin.y = (frame.origin.y + frame.size.height) - (float)height;
        frame.size.width = (float)width;
        frame.size.height = (float)height;
        dumpFrame("after", frame);

        msg(app->mainWindow, s("setFrame:display:animate:"), frame, 1, 0);
    )
}

void SetPosition(struct Application *app, int x, int y) { 
    ON_MAIN_THREAD(
        id screen = getCurrentScreen(app);
        CGRect screenFrame = GET_FRAME(screen);
        CGRect windowFrame = GET_FRAME(app->mainWindow);

        windowFrame.origin.x = screenFrame.origin.x + (float)x;
        windowFrame.origin.y = (screenFrame.origin.y + screenFrame.size.height) - windowFrame.size.height - (float)y;
        msg(app->mainWindow, s("setFrame:display:animate:"), windowFrame, 1, 0);
    )
}

// OpenDialog opens a dialog to select files/directories
void OpenDialog(struct Application *app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolveAliases, int treatPackagesAsDirectories) {
    Debug("OpenDialog Called with callback id: %s", callbackID);

    // Create an open panel
    ON_MAIN_THREAD(

        // Create the dialog
        id dialog = msg(c("NSOpenPanel"), s("openPanel"));

        // Valid but appears to do nothing.... :/
        msg(dialog, s("setTitle:"), str(title));

        // Filters
        if( filters != NULL && strlen(filters) > 0) {
            id filterString = msg(str(filters), s("stringByReplacingOccurrencesOfString:withString:"), str("*."), str(""));
            filterString = msg(filterString, s("stringByReplacingOccurrencesOfString:withString:"), str(" "), str(""));
            id filterList = msg(filterString, s("componentsSeparatedByString:"), str(","));
            msg(dialog, s("setAllowedFileTypes:"), filterList);
        } else {
            msg(dialog, s("setAllowsOtherFileTypes:"), YES);
        }

        // Default Directory
        if( defaultDir != NULL && strlen(defaultDir) > 0 ) {
            msg(dialog, s("setDirectoryURL:"), url(defaultDir));
        }

        // Default Filename
        if( defaultFilename != NULL && strlen(defaultFilename) > 0 ) {
            msg(dialog, s("setNameFieldStringValue:"), str(defaultFilename));
        }

        // Setup Options
        msg(dialog, s("setCanChooseFiles:"), allowFiles);
        msg(dialog, s("setCanChooseDirectories:"), allowDirs);
        msg(dialog, s("setAllowsMultipleSelection:"), allowMultiple);
        msg(dialog, s("setShowsHiddenFiles:"), showHiddenFiles);
        msg(dialog, s("setCanCreateDirectories:"), canCreateDirectories);
        msg(dialog, s("setResolvesAliases:"), resolveAliases);
        msg(dialog, s("setTreatsFilePackagesAsDirectories:"), treatPackagesAsDirectories);

        // Setup callback handler
        msg(dialog, s("beginSheetModalForWindow:completionHandler:"), app->mainWindow, ^(id result) {
        
            // Create the response JSON object
            JsonNode *response = json_mkarray();

            // If the user selected some files
            if( result == (id)1 ) {
                // Grab the URLs returned
                id urls = msg(dialog, s("URLs"));

                // Iterate over all the selected files
                int noOfResults = (int)msg(urls, s("count"));
                for( int index = 0; index < noOfResults; index++ ) {

                    // Extract the filename
                    id url = msg(urls, s("objectAtIndex:"), index);
                    const char *filename = (const char *)msg(msg(url, s("path")), s("UTF8String"));

                    // Add the the response array
                    json_append_element(response, json_mkstring(filename));
                }
            }

            // Create JSON string and free json memory
            char *encoded = json_stringify(response, "");
            json_delete(response);

            // Construct callback message. Format "D<callbackID>|<json array of strings>"
            const char *callback = concat("DO", callbackID);
            const char *header = concat(callback, "|");
            const char *responseMessage = concat(header, encoded);

            // Send message to backend
            app->sendMessageToBackend(responseMessage); 

            // Free memory
            free((void*)header);
            free((void*)callback);
            free((void*)responseMessage);
        });

        msg( c("NSApp"), s("runModalForWindow:"), app->mainWindow);
    )
}

// SaveDialog opens a dialog to select files/directories
void SaveDialog(struct Application *app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories) {
    Debug("SaveDialog Called with callback id: %s", callbackID);

    // Create an open panel
    ON_MAIN_THREAD(

        // Create the dialog
        id dialog = msg(c("NSSavePanel"), s("savePanel"));

        // Valid but appears to do nothing.... :/
        msg(dialog, s("setTitle:"), str(title));

        // Filters
        if( filters != NULL && strlen(filters) > 0) {
            id filterString = msg(str(filters), s("stringByReplacingOccurrencesOfString:withString:"), str("*."), str(""));
            filterString = msg(filterString, s("stringByReplacingOccurrencesOfString:withString:"), str(" "), str(""));
            id filterList = msg(filterString, s("componentsSeparatedByString:"), str(","));
            msg(dialog, s("setAllowedFileTypes:"), filterList);
        } else {
            msg(dialog, s("setAllowsOtherFileTypes:"), YES);
        }

        // Default Directory
        if( defaultDir != NULL && strlen(defaultDir) > 0 ) {
            msg(dialog, s("setDirectoryURL:"), url(defaultDir));
        }

        // Default Filename
        if( defaultFilename != NULL && strlen(defaultFilename) > 0 ) {
            msg(dialog, s("setNameFieldStringValue:"), str(defaultFilename));
        }

        // Setup Options
        msg(dialog, s("setShowsHiddenFiles:"), showHiddenFiles);
        msg(dialog, s("setCanCreateDirectories:"), canCreateDirectories);
        msg(dialog, s("setTreatsFilePackagesAsDirectories:"), treatPackagesAsDirectories);

        // Setup callback handler
        msg(dialog, s("beginSheetModalForWindow:completionHandler:"), app->mainWindow, ^(id result) {
        
            // Default is blank
            const char *filename = "";

            // If the user selected some files
            if( result == (id)1 ) {
                // Grab the URL returned
                id url = msg(dialog, s("URL"));
                filename = (const char *)msg(msg(url, s("path")), s("UTF8String"));
            }

            // Construct callback message. Format "DS<callbackID>|<json array of strings>"
            const char *callback = concat("DS", callbackID);
            const char *header = concat(callback, "|");
            const char *responseMessage = concat(header, filename);

            // Send message to backend
            app->sendMessageToBackend(responseMessage); 

            // Free memory
            free((void*)header);
            free((void*)callback);
            free((void*)responseMessage);
        });

        msg( c("NSApp"), s("runModalForWindow:"), app->mainWindow);
    )
}

const char *invoke = "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}};";

// DisableFrame disables the window frame
void DisableFrame(struct Application *app)
{
   app->frame = 0;
}

void setMinMaxSize(struct Application *app)
{
    if (app->maxHeight > 0 && app->maxWidth > 0)
    {
        msg(app->mainWindow, s("setMaxSize:"), CGSizeMake(app->maxWidth, app->maxHeight));
    }
    if (app->minHeight > 0 && app->minWidth > 0)
    {
        msg(app->mainWindow, s("setMinSize:"), CGSizeMake(app->minWidth, app->minHeight));
    }
}

void SetMinWindowSize(struct Application *app, int minWidth, int minHeight)
{
    app->minWidth = minWidth;
    app->minHeight = minHeight;

    // Apply if the window is created
    if( app->mainWindow != NULL ) {
        ON_MAIN_THREAD(
            setMinMaxSize(app);
        )
    }
}

void SetMaxWindowSize(struct Application *app, int maxWidth, int maxHeight)
{
    app->maxWidth = maxWidth;
    app->maxHeight = maxHeight;
    
    // Apply if the window is created
    if( app->mainWindow != NULL ) {
        ON_MAIN_THREAD(
            setMinMaxSize(app);
        )
    }
}

void ExecJS(struct Application *app, const char *js) {
    ON_MAIN_THREAD(
        msg(app->wkwebview, 
            s("evaluateJavaScript:completionHandler:"),
            str(js), 
            NULL);
    )
}

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

void enableBoolConfig(id config, const char *setting) {
    msg(msg(config, s("preferences")), s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 1), str(setting));
}

void disableBoolConfig(id config, const char *setting) {
    msg(msg(config, s("preferences")), s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 0), str(setting));
}

void Run(void *applicationPointer, int argc, char **argv) {
    struct Application *app = (struct Application*) applicationPointer;

    int decorations = 0;

    if (app->frame == 1 ) { 
        if( app->hideTitleBar == 0) {
            decorations |= NSWindowStyleMaskTitled;
        }
        decorations |= NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable;
    }

    if (app->resizable) {
        decorations |= NSWindowStyleMaskResizable;
    }

    if (app->fullscreen) {
        decorations |= NSWindowStyleMaskFullscreen;
    }

    if( app->fullSizeContent || app->frame == 0) {
        decorations |= NSWindowStyleMaskFullSizeContentView;
    }


    id application = msg(c("NSApplication"), s("sharedApplication"));
    app->application = application;
    msg(application, s("setActivationPolicy:"), 0);

    // Define delegate
    Class delegateClass = objc_allocateClassPair((Class) c("NSResponder"), "AppDelegate", 0);
    class_addProtocol(delegateClass, objc_getProtocol("NSTouchBarProvider"));
    class_addMethod(delegateClass, s("applicationShouldTerminateAfterLastWindowClosed:"), (IMP) yes, "c@:@");

    // Script handler
    class_addMethod(delegateClass, s("userContentController:didReceiveScriptMessage:"), (IMP) messageHandler,
                    "v@:@@");
    objc_registerClassPair(delegateClass);

    // Create delegate
    id delegate = msg((id)delegateClass, s("new"));
    objc_setAssociatedObject(delegate, "application", (id)app, OBJC_ASSOCIATION_ASSIGN);
    msg(application, s("setDelegate:"), delegate);

    // Create main window
    id mainWindow = msg(c("NSWindow"),s("alloc"));
    mainWindow = msg(mainWindow, s("initWithContentRect:styleMask:backing:defer:"),
          CGRectMake(0, 0, app->width, app->height), decorations, NSBackingStoreBuffered, 0);
    msg(mainWindow, s("autorelease"));

    app->mainWindow = mainWindow;

    // Set the main window title
    SetTitle(app, app->title);

    // Center Window
    Center(app);

    // Set Style Mask
    msg(mainWindow, s("setStyleMask:"), decorations);

    // Set Colour
    applyWindowColour(app);
    
    // Setup webview
    id config = msg(c("WKWebViewConfiguration"), s("new"));
    app->config = config;
    id manager = msg(config, s("userContentController"));
    app->manager = manager;
    id wkwebview = msg(c("WKWebView"), s("alloc"));
    app->wkwebview = wkwebview;

    // Only show content when fully rendered
    msg(config, s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 1), str("suppressesIncrementalRendering"));
    
    // TODO: Fix "NSWindow warning: adding an unknown subview: <WKInspectorWKWebView: 0x465ed90>. Break on NSLog to debug." error
    if (app->devtools) {
      Debug("Enabling devtools");
      enableBoolConfig(config, "developerExtrasEnabled");
    }
    msg(wkwebview, s("initWithFrame:configuration:"), CGRectMake(0, 0, 0, 0), config);

    
    msg(manager, s("addScriptMessageHandler:name:"), delegate, str("external"));
    msg(manager, s("addScriptMessageHandler:name:"), delegate, str("completed"));
    msg(mainWindow, s("setContentView:"), wkwebview);

    if( app->frame == 0) {
        msg(mainWindow, s("setTitlebarAppearsTransparent:"), YES);
        msg(mainWindow, s("setTitleVisibility:"), NSWindowTitleHidden);

        // Setup drag message handler
        msg(manager, s("addScriptMessageHandler:name:"), delegate, str("windowDrag"));
        // Add mouse event hooks
        // msg(mainWindow, s("mouseDown:"), ^(id event) {
        //     Debug("Got mouse down event!!!!");
        //     app->mouseEvent = event;
        // });

    } else {
        Debug("setTitlebarAppearsTransparent %d", app->titlebarAppearsTransparent ? YES :NO);
        msg(mainWindow, s("setTitlebarAppearsTransparent:"), app->titlebarAppearsTransparent ? YES : NO);
        msg(app->mainWindow, s("setTitleVisibility:"), app->hideTitle);

        // Toolbar
        if( app->useToolBar ) {
            Debug("Setting Toolbar");
            id toolbar = msg(c("NSToolbar"),s("alloc"));
            msg(toolbar, s("initWithIdentifier:"), str("wails.toolbar"));
            msg(toolbar, s("autorelease"));

            // Separator
            if( app->hideToolbarSeparator ) {
                msg(toolbar, s("setShowsBaselineSeparator:"), NO);
            }

            msg(mainWindow, s("setToolbar:"), toolbar);
        }
    }

    // Fix up resizing
    if (app->resizable == 0) {
        app->minHeight = app->maxHeight = app->height;
        app->minWidth = app->maxWidth = app->width;
    }
    setMinMaxSize(app);

    // Load HTML
    id html = msg(c("NSURL"), s("URLWithString:"), str(assets[0]));
    msg(wkwebview, s("loadRequest:"), msg(c("NSURLRequest"), s("requestWithURL:"), html));
    
    Debug("Loading Internal Code");
    // We want to evaluate the internal code plus runtime before the assets
    const char *temp = concat(invoke, app->bindings);
    const char *internalCode = concat(temp, (const char*)&runtime);
    free((void*)temp);

      // Loop over assets and build up one giant Mother Of All Evals
    int index = 1;
    while(1) {
        // Get next asset pointer
        const unsigned char *asset = assets[index];

        // If we have no more assets, break
        if (asset == 0x00) {
            break;
        }

        temp = concat(internalCode, (const char *)asset);
        free((void*)internalCode);
        internalCode = temp;
        index++;
    };

    class_addMethod(delegateClass, s("closeWindow"), (IMP) closeWindow, "v@:@");

    // Include callback after evaluation
    temp = concat(internalCode, "webkit.messageHandlers.completed.postMessage(true);");
    free((void*)internalCode);
    internalCode = temp;

    // const char *viewportScriptString = "var meta = document.createElement('meta'); meta.setAttribute('name', 'viewport'); meta.setAttribute('content', 'width=device-width'); meta.setAttribute('initial-scale', '1.0'); meta.setAttribute('maximum-scale', '1.0'); meta.setAttribute('minimum-scale', '1.0'); meta.setAttribute('user-scalable', 'no'); document.getElementsByTagName('head')[0].appendChild(meta);";
    // ExecJS(app, viewportScriptString);
    

    // This evaluates the MOAE once the Dom has finished loading
    msg(manager, 
        s("addUserScript:"),
        msg(msg(c("WKUserScript"), s("alloc")),
                    s("initWithSource:injectionTime:forMainFrameOnly:"),
                    str(internalCode),
                    1, 
                    1));

    // Finally call run
    Debug("Run called");
    msg(application, s("activateIgnoringOtherApps:"), true);
    msg(application, s("run"));

    free((void*)internalCode);
}

#endif
