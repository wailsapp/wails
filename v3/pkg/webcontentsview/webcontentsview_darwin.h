#ifndef webcontentsview_darwin_h
#define webcontentsview_darwin_h

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

typedef struct {
    bool devTools;
    bool javascript;
    bool webSecurity;
    bool images;
    bool plugins;
    double zoomFactor;
    int defaultFontSize;
    int defaultMonospaceFontSize;
    int minimumFontSize;
    const char* userAgent;

} WebContentsViewPreferences;

extern void* createWebContentsView(int x, int y, int w, int h, WebContentsViewPreferences prefs);
extern void webContentsViewSetBounds(void* view, int x, int y, int w, int h);
extern void webContentsViewSetURL(void* view, const char* url);
extern void webContentsViewExecJS(void* view, const char* js);
extern void webContentsViewGoBack(void* view);
extern const char* webContentsViewGetURL(void* view);

extern void windowAddWebContentsView(void* nsWindow, void* view);
extern void windowRemoveWebContentsView(void* nsWindow, void* view);

// Async snapshot
extern void webContentsViewTakeSnapshot(void* view, uintptr_t callbackID);
extern void webContentsViewConfigureAutomation(void* view, uintptr_t viewID, const char* pageScript, const char* automationScript);
extern void webContentsViewAutomationEvaluate(void* view, const char* expression, int world, bool awaitPromise, uintptr_t callbackID);
extern void webContentsViewAutomationInvoke(void* view, const char* method, const char* paramsJSON, uintptr_t callbackID);
extern void webContentsViewGetCookies(void* view, uintptr_t callbackID);
extern void webContentsViewSetCookie(void* view, const char* cookieJSON, uintptr_t callbackID);
extern void webContentsViewDeleteCookie(void* view, const char* name, const char* domain, const char* path, uintptr_t callbackID);
extern void webContentsViewClearCookies(void* view, uintptr_t callbackID);
extern void webContentsViewCreatePDF(void* view, uintptr_t callbackID);
extern bool webContentsViewAutomationSetInspectable(void* view, bool enabled);
extern bool webContentsViewAutomationSupportsInspection(void);
extern bool webContentsViewAutomationSupportsProxyCapture(void);

#endif /* webcontentsview_darwin_h */
