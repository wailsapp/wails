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

extern void windowAddWebContentsView(void* nsWindow, void* view);
extern void windowRemoveWebContentsView(void* nsWindow, void* view);

#endif /* webcontentsview_darwin_h */
extern void webContentsViewGoBack(void* view);
extern const char* webContentsViewGetURL(void* view);
