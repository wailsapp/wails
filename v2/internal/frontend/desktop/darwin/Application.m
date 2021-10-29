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

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int debug) {
    
    WailsContext *result = [WailsContext new];

    result.debug = debug;
    
    [result CreateWindow:width :height :frameless :resizable :fullscreen :fullSizeContent :hideTitleBar :titlebarAppearsTransparent :hideTitle :useToolbar :hideToolbarSeparator :webviewIsTransparent :hideWindowOnClose :appearance :windowIsTranslucent];
    [result SetTitle:title];
    [result Center];
    
    result.alwaysOnTop = alwaysOnTop;
    result.hideOnClose = hideWindowOnClose;
        
    return result;
}

void ProcessURLResponse(void *inctx, const char *url, const char *contentType, const char* data, int datalength) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *nsurl = [[NSString alloc] initWithUTF8String:url];
    NSString *nsContentType = [[NSString alloc] initWithUTF8String:contentType];
    NSData *nsdata = [NSData dataWithBytes:data length:datalength];
    
    [ctx processURLResponse:nsurl :nsContentType :nsdata];
}

void ExecJS(void* inctx, const char *script) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx ExecJS:script];
    );
}

void SetTitle(void* inctx, const char *title) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetTitle:title];
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
       [ctx SetSize:x :y];
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

const char* GetSize(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSRect frame = [ctx.mainWindow frame];
    NSString *result = [NSString stringWithFormat:@"%d,%d", (int)frame.size.width, (int)frame.size.height];
    return [result UTF8String];
}

const char* GetPos(void *inctx) {
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

void MessageDialog(void *inctx, const char* dialogType, const char* title, const char* message, const char* button1, const char* button2, const char* button3, const char* button4, const char* defaultButton, const char* cancelButton) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
                   [ctx MessageDialog:dialogType :title :message :button1 :button2 :button3 :button4 :defaultButton :cancelButton];
    )
}

void OpenFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int allowDirectories, int allowFiles, int canCreateDirectories, int treatPackagesAsDirectories, int resolveAliases, int showHiddenFiles, int allowMultipleSelection, const char* filters) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
                   [ctx OpenFileDialog:title :defaultFilename :defaultDirectory :allowDirectories :allowFiles :canCreateDirectories :treatPackagesAsDirectories :resolveAliases :showHiddenFiles :allowMultipleSelection :filters];
    )
}

void SaveFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int canCreateDirectories, int treatPackagesAsDirectories, int showHiddenFiles, const char* filters) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
                   [ctx SaveFileDialog:title :defaultFilename :defaultDirectory :canCreateDirectories :treatPackagesAsDirectories :showHiddenFiles :filters];
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

void SetAbout(void *inctx, const char* title, const char* description, void* imagedata, int datalen) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [ctx SetAbout :title :description :imagedata :datalen];
}

void* AppendMenuItem(void* inctx, void* inMenu, const char* label, const char* shortcutKey, int modifiers, int disabled, int checked, int menuItemID) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    WailsMenu *menu = (__bridge WailsMenu*) inMenu;
    return [menu AppendMenuItem:ctx :label :shortcutKey :modifiers :disabled :checked :menuItemID];
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

    [ctx loadRequest:@"wails://wails/"];
    [app setMainMenu:ctx.applicationMenu];
    [app run];
    [ctx release];
}
