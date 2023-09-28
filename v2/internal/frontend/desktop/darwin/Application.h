//
//  Application.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef Application_h
#define Application_h

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WailsContext.h"

#define WindowStartsNormal 0
#define WindowStartsMaximised 1
#define WindowStartsMinimised 2
#define WindowStartsFullscreen 3

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int devtoolsEnabled, int defaultContextMenuEnabled, int windowStartState, int startsHidden, int minWidth, int minHeight, int maxWidth, int maxHeight, bool fraudulentWebsiteWarningEnabled, struct Preferences preferences, int singleInstanceEnabled, const char* singleInstanceUniqueId);
void Run(void*, const char* url);

void SetTitle(void* ctx, const char *title);
void Center(void* ctx);
void SetSize(void* ctx, int width, int height);
void SetAlwaysOnTop(void* ctx, int onTop);
void SetMinSize(void* ctx, int width, int height);
void SetMaxSize(void* ctx, int width, int height);
void SetPosition(void* ctx, int x, int y);
void Fullscreen(void* ctx);
void UnFullscreen(void* ctx);
void Minimise(void* ctx);
void UnMinimise(void* ctx);
void ToggleMaximise(void* ctx);
void Maximise(void* ctx);
void UnMaximise(void* ctx);
void Hide(void* ctx);
void Show(void* ctx);
void HideApplication(void* ctx);
void ShowApplication(void* ctx);
void SetBackgroundColour(void* ctx, int r, int g, int b, int a);
void ExecJS(void* ctx, const char*);
void Quit(void*);
void WindowPrint(void* ctx);

const char* GetSize(void *ctx);
const char* GetPosition(void *ctx);
const bool IsFullScreen(void *ctx);
const bool IsMinimised(void *ctx);
const bool IsMaximised(void *ctx);

/* Dialogs */

void MessageDialog(void *inctx, const char* dialogType, const char* title, const char* message, const char* button1, const char* button2, const char* button3, const char* button4, const char* defaultButton, const char* cancelButton, void* iconData, int iconDataLength);
void OpenFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int allowDirectories, int allowFiles, int canCreateDirectories, int treatPackagesAsDirectories, int resolveAliases, int showHiddenFiles, int allowMultipleSelection, const char* filters);
void SaveFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int canCreateDirectories, int treatPackagesAsDirectories, int showHiddenFiles, const char* filters);

/* Application Menu */
void* NewMenu(const char* name);
void AppendSubmenu(void* parent, void* child);
void AppendRole(void *inctx, void *inMenu, int role);
void SetAsApplicationMenu(void *inctx, void *inMenu);
void UpdateApplicationMenu(void *inctx);

void SetAbout(void *inctx, const char* title, const char* description, void* imagedata, int datalen);
void* AppendMenuItem(void* inctx, void* nsmenu, const char* label, const char* shortcutKey, int modifiers, int disabled, int checked, int menuItemID);
void AppendSeparator(void* inMenu);
void UpdateMenuItem(void* nsmenuitem, int checked);
void RunMainLoop(void);
void ReleaseContext(void *inctx);

NSString* safeInit(const char* input);

#endif /* Application_h */
