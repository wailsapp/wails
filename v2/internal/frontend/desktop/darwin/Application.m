//go:build darwin
//
//  Application.m
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WailsContext.h"
#import "Application.h"
#import "AppDelegate.h"
#import "WindowDelegate.h"
#import "WailsMenu.h"
#import "WailsMenuItem.h"

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int devtoolsEnabled, int defaultContextMenuEnabled, int windowStartState, int startsHidden, int minWidth, int minHeight, int maxWidth, int maxHeight, bool fraudulentWebsiteWarningEnabled, struct Preferences preferences, int singleInstanceLockEnabled, const char* singleInstanceUniqueId) {

    [NSApplication sharedApplication];

    WailsContext *result = [WailsContext new];

    result.devtoolsEnabled = devtoolsEnabled;
    result.defaultContextMenuEnabled = defaultContextMenuEnabled;

    if ( windowStartState == WindowStartsFullscreen ) {
        fullscreen = 1;
    }

    [result CreateWindow:width :height :frameless :resizable :fullscreen :fullSizeContent :hideTitleBar :titlebarAppearsTransparent :hideTitle :useToolbar :hideToolbarSeparator :webviewIsTransparent :hideWindowOnClose :safeInit(appearance) :windowIsTranslucent :minWidth :minHeight :maxWidth :maxHeight :fraudulentWebsiteWarningEnabled :preferences];
    [result SetTitle:safeInit(title)];
    [result Center];

    switch( windowStartState ) {
        case WindowStartsMaximised:
            [result.mainWindow zoom:nil];
            break;
        case WindowStartsMinimised:
            //TODO: Can you start a mac app minimised?
            break;
    }

    if ( startsHidden == 1 ) {
        result.startHidden = true;
    }

    if ( fullscreen == 1 ) {
        result.startFullscreen = true;
    }

    if ( singleInstanceLockEnabled == 1 ) {
        result.singleInstanceLockEnabled = true;
        result.singleInstanceUniqueId = safeInit(singleInstanceUniqueId);
    }

    result.alwaysOnTop = alwaysOnTop;
    result.hideOnClose = hideWindowOnClose;

    return result;
}

void ExecJS(void* inctx, const char *script) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *nsscript = safeInit(script);
    ON_MAIN_THREAD(
       [ctx ExecJS:nsscript];
       [nsscript release];
    );
}

void SetTitle(void* inctx, const char *title) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *_title = safeInit(title);
    ON_MAIN_THREAD(
       [ctx SetTitle:_title];
    );
}


void SetBackgroundColour(void *inctx, int r, int g, int b, int a) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetBackgroundColour:r :g :b :a];
    );
}

void SetSize(void* inctx, int width, int height) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetSize:width :height];
    );
}

void SetAlwaysOnTop(void* inctx, int onTop) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetAlwaysOnTop:onTop];
    );
}

void SetMinSize(void* inctx, int width, int height) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetMinSize:width :height];
    );
}

void SetMaxSize(void* inctx, int width, int height) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetMaxSize:width :height];
    );
}

void SetPosition(void* inctx, int x, int y) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetPosition:x :y];
    );
}

void Center(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Center];
    );
}

void Fullscreen(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Fullscreen];
    );
}

void UnFullscreen(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx UnFullscreen];
    );
}

void Minimise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Minimise];
    );
}

void UnMinimise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx UnMinimise];
    );
}

void Maximise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Maximise];
    );
}

void ToggleMaximise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx ToggleMaximise];
    );
}

const char* GetSize(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSRect frame = [ctx.mainWindow frame];
    NSString *result = [NSString stringWithFormat:@"%d,%d", (int)frame.size.width, (int)frame.size.height];
    return [result UTF8String];
}

const char* GetPosition(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSScreen* screen = [ctx getCurrentScreen];
    NSRect windowFrame = [ctx.mainWindow frame];
    NSRect screenFrame = [screen visibleFrame];
    int x = windowFrame.origin.x - screenFrame.origin.x;
    int y = windowFrame.origin.y - screenFrame.origin.y;
    y = screenFrame.size.height - y - windowFrame.size.height;
    NSString *result = [NSString stringWithFormat:@"%d,%d",x,y];
    return [result UTF8String];
}

const bool IsFullScreen(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    return [ctx IsFullScreen];
}

const bool IsMinimised(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    return [ctx IsMinimised];
}

const bool IsMaximised(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    return [ctx IsMaximised];
}

void UnMaximise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx UnMaximise];
    );
}

void Quit(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [NSApp stop:ctx];
    [NSApp abortModal];
}

void Hide(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Hide];
    );
}

void Show(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx Show];
    );
}


void HideApplication(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx HideApplication];
    );
}

void ShowApplication(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx ShowApplication];
    );
}

NSString* safeInit(const char* input) {
    NSString *result = nil;
    if (input != nil) {
        result = [NSString stringWithUTF8String:input];
    }
    return result;
}

void MessageDialog(void *inctx, const char* dialogType, const char* title, const char* message, const char* button1, const char* button2, const char* button3, const char* button4, const char* defaultButton, const char* cancelButton, void* iconData, int iconDataLength) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;

    NSString *_dialogType = safeInit(dialogType);
    NSString *_title = safeInit(title);
    NSString *_message = safeInit(message);
    NSString *_button1 = safeInit(button1);
    NSString *_button2 = safeInit(button2);
    NSString *_button3 = safeInit(button3);
    NSString *_button4 = safeInit(button4);
    NSString *_defaultButton = safeInit(defaultButton);
    NSString *_cancelButton = safeInit(cancelButton);

    ON_MAIN_THREAD(
                   [ctx MessageDialog:_dialogType :_title :_message :_button1 :_button2 :_button3 :_button4 :_defaultButton :_cancelButton :iconData :iconDataLength];
    )
}

void OpenFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int allowDirectories, int allowFiles, int canCreateDirectories, int treatPackagesAsDirectories, int resolveAliases, int showHiddenFiles, int allowMultipleSelection, const char* filters) {

    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *_title = safeInit(title);
    NSString *_defaultFilename = safeInit(defaultFilename);
    NSString *_defaultDirectory = safeInit(defaultDirectory);
    NSString *_filters = safeInit(filters);

    ON_MAIN_THREAD(
                   [ctx OpenFileDialog:_title :_defaultFilename :_defaultDirectory :allowDirectories :allowFiles :canCreateDirectories :treatPackagesAsDirectories :resolveAliases :showHiddenFiles :allowMultipleSelection :_filters];
    )
}

void SaveFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int canCreateDirectories, int treatPackagesAsDirectories, int showHiddenFiles, const char* filters) {

    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *_title = safeInit(title);
    NSString *_defaultFilename = safeInit(defaultFilename);
    NSString *_defaultDirectory = safeInit(defaultDirectory);
    NSString *_filters = safeInit(filters);

    ON_MAIN_THREAD(
                   [ctx SaveFileDialog:_title :_defaultFilename :_defaultDirectory :canCreateDirectories :treatPackagesAsDirectories :showHiddenFiles :_filters];
    )
}

void AppendRole(void *inctx, void *inMenu, int role) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    WailsMenu *menu = (__bridge WailsMenu*) inMenu;
    [menu appendRole :ctx :role];
}

void* NewMenu(const char *name) {
    NSString *title = @"";
    if (name != nil) {
        title = [NSString stringWithUTF8String:name];
    }
    WailsMenu *result = [[WailsMenu new] initWithNSTitle:title];
    return result;
}

void AppendSubmenu(void* inparent, void* inchild) {
    WailsMenu *parent = (__bridge WailsMenu*) inparent;
    WailsMenu *child = (__bridge WailsMenu*) inchild;
    [parent appendSubmenu:child];
}

void SetAsApplicationMenu(void *inctx, void *inMenu) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    WailsMenu *menu = (__bridge WailsMenu*) inMenu;
    ctx.applicationMenu = menu;
}

void UpdateApplicationMenu(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
                   NSApplication *app = [NSApplication sharedApplication];
                   [app setMainMenu:ctx.applicationMenu];
    )
}

void SetAbout(void *inctx, const char* title, const char* description, void* imagedata, int datalen) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *_title = safeInit(title);
    NSString *_description = safeInit(description);

    [ctx SetAbout :_title :_description :imagedata :datalen];
}

void* AppendMenuItem(void* inctx, void* inMenu, const char* label, const char* shortcutKey, int modifiers, int disabled, int checked, int menuItemID) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    WailsMenu *menu = (__bridge WailsMenu*) inMenu;
    NSString *_label = safeInit(label);
    NSString *_shortcutKey = safeInit(shortcutKey);

    return [menu AppendMenuItem:ctx :_label :_shortcutKey :modifiers :disabled :checked :menuItemID];
}

void UpdateMenuItem(void* nsmenuitem, int checked) {
    ON_MAIN_THREAD(
        WailsMenuItem *menuItem = (__bridge WailsMenuItem*) nsmenuitem;
        [menuItem setState:(checked == 1?NSControlStateValueOn:NSControlStateValueOff)];
                   )
}


void AppendSeparator(void* inMenu) {
    WailsMenu *menu = (__bridge WailsMenu*) inMenu;
    [menu AppendSeparator];
}



void Run(void *inctx, const char* url) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSApplication *app = [NSApplication sharedApplication];
    AppDelegate* delegate = [AppDelegate new];
    [app setDelegate:(id)delegate];
    ctx.appdelegate = delegate;
    delegate.mainWindow = ctx.mainWindow;
    delegate.alwaysOnTop = ctx.alwaysOnTop;
    delegate.startHidden = ctx.startHidden;
    delegate.singleInstanceLockEnabled = ctx.singleInstanceLockEnabled;
    delegate.singleInstanceUniqueId = ctx.singleInstanceUniqueId;
    delegate.startFullscreen = ctx.startFullscreen;

    NSString *_url = safeInit(url);
    [ctx loadRequest:_url];
    [_url release];

    [app setMainMenu:ctx.applicationMenu];
}

void RunMainLoop(void) {
    NSApplication *app = [NSApplication sharedApplication];
    [app run];
}

void ReleaseContext(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [ctx release];
}

// Credit: https://stackoverflow.com/q/33319295
void WindowPrint(void *inctx) {

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	if (@available(macOS 11.0, *)) {
        ON_MAIN_THREAD(
            WailsContext *ctx = (__bridge WailsContext*) inctx;
            WKWebView* webView = ctx.webview;

            // I think this should be exposed as a config
            // It directly affects the printed output/PDF
            NSPrintInfo *pInfo = [NSPrintInfo sharedPrintInfo];
            pInfo.horizontalPagination = NSPrintingPaginationModeAutomatic;
            pInfo.verticalPagination = NSPrintingPaginationModeAutomatic;
            pInfo.verticallyCentered = YES;
            pInfo.horizontallyCentered = YES;
            pInfo.orientation = NSPaperOrientationLandscape;
            pInfo.leftMargin = 0;
            pInfo.rightMargin = 0;
            pInfo.topMargin = 0;
            pInfo.bottomMargin = 0;

            NSPrintOperation *po = [webView printOperationWithPrintInfo:pInfo];
            po.showsPrintPanel = YES;
            po.showsProgressPanel = YES;

            po.view.frame = webView.bounds;

            [po runOperationModalForWindow:ctx.mainWindow delegate:ctx.mainWindow.delegate didRunSelector:nil contextInfo:nil];
        )
	}
#endif
}
