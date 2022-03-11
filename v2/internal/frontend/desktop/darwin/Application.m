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
#import "WailsMenu.h"
#import "WailsMenuItem.h"

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int debug, int windowStartState, int startsHidden, int minWidth, int minHeight, int maxWidth, int maxHeight) {
    
    [NSApplication sharedApplication];

    WailsContext *result = [WailsContext new];

    result.debug = debug;
    
    if ( windowStartState == WindowStartsFullscreen ) {
        fullscreen = 1;
    }

    [result CreateWindow:width :height :frameless :resizable :fullscreen :fullSizeContent :hideTitleBar :titlebarAppearsTransparent :hideTitle :useToolbar :hideToolbarSeparator :webviewIsTransparent :hideWindowOnClose :safeInit(appearance) :windowIsTranslucent :minWidth :minHeight :maxWidth :maxHeight];
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
    
    result.alwaysOnTop = alwaysOnTop;
    result.hideOnClose = hideWindowOnClose;
        
    return result;
}

void ProcessURLResponse(void *inctx, const char *url, int statusCode, const char *contentType, void* data, int datalength) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *nsurl = safeInit(url);
    NSString *nsContentType = safeInit(contentType);
    NSData *nsdata = [NSData dataWithBytes:data length:datalength];
    
    [ctx processURLResponse:nsurl :statusCode :nsContentType :nsdata];

    [nsdata release];
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


void SetRGBA(void *inctx, int r, int g, int b, int a) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetRGBA:r :g :b :a];
    );
}

void SetSize(void* inctx, int width, int height) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetSize:width :height];
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

void UnMaximise(void* inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx UnMaximise];
    );
}

void Quit(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [NSApp stop:ctx];
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



void Run(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSApplication *app = [NSApplication sharedApplication];
    AppDelegate* delegate = [AppDelegate new];
    [app setDelegate:(id)delegate];
    ctx.appdelegate = delegate;
    delegate.mainWindow = ctx.mainWindow;
    delegate.alwaysOnTop = ctx.alwaysOnTop;
    delegate.startHidden = ctx.startHidden;
    delegate.startFullscreen = ctx.startFullscreen;

    [ctx loadRequest:@"wails://wails/"];
    [app setMainMenu:ctx.applicationMenu];
    [app run];
    [ctx release];
}
