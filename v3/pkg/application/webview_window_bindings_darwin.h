#ifndef WEBVIEW_WINDOW_BINDINGS_DARWIN
#define WEBVIEW_WINDOW_BINDINGS_DARWIN

#include "application_darwin.h"
#include "webview_window_darwin.h"
#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import <WebKit/WebKit.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

struct WebviewPreferences {
	bool *TabFocusesLinks;
	bool *TextInteractionEnabled;
	bool *FullscreenEnabled;
};

typedef enum {
	WindowTypeWindow,
	WindowTypePanel
} WindowType;

extern void registerListener(unsigned int event);

void *createWindow(WindowType windowType, unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences);
void printWindowStyle(void *window);
void setInvisibleTitleBarHeight(void *window, unsigned int height);
void windowSetTransparent(void *window);
void windowSetInvisibleTitleBar(void *window, unsigned int height);
void windowSetTitle(void *window, char *title);
void windowSetSize(void *window, int width, int height);
void windowSetAlwaysOnTop(void *window, bool alwaysOnTop);
void setNormalWindowLevel(void *window);
void setFloatingWindowLevel(void *window);
void setPopUpMenuWindowLevel(void *window);
void setMainMenuWindowLevel(void *window);
void setStatusWindowLevel(void *window);
void setModalPanelWindowLevel(void *window);
void setScreenSaverWindowLevel(void *window);
void setTornOffMenuWindowLevel(void *window);
void navigationLoadURL(void *window, char *url);
void windowSetResizable(void *window, bool resizable);
void windowSetMinSize(void *window, int width, int height);
void windowSetMaxSize(void *window, int width, int height);
void windowZoomReset(void *window);
void windowZoomSet(void *window, double zoom);
float windowZoomGet(void *window);
void windowZoomIn(void *window);
void windowZoomOut(void *window);
void windowSetRelativePosition(void *window, int x, int y);
void windowExecJS(void *window, const char *js);
void windowSetTranslucent(void *window);
void webviewSetTransparent(void *window);
void webviewSetBackgroundColour(void *window, int r, int g, int b, int alpha);
void windowSetBackgroundColour(void *window, int r, int g, int b, int alpha);
bool windowIsMaximised(void *window);
bool windowIsFullScreen (void *window);
bool windowIsMinimised(void *window);
bool windowIsFocused(void *window);
void windowFullscreen(void *window);
void windowUnFullscreen(void *window);
void windowRestore(void *window);
void setFullscreenButtonEnabled(void *window, bool enabled);
void windowSetTitleBarAppearsTransparent(void *window, bool transparent);
void windowSetFullSizeContent(void *window, bool fullSize);
void windowSetHideTitleBar(void *window, bool hideTitlebar);
void windowSetHideTitle(void *window, bool hideTitle);
void windowSetUseToolbar(void *window, bool useToolbar);
void windowSetToolbarStyle(void *window, int style);
void windowSetHideToolbarSeparator(void *window, bool hideSeparator);
void windowSetShowToolbarWhenFullscreen(void *window, bool setting);
void windowSetAppearanceTypeByName(void *window, const char *appearanceName);
void windowCenter(void *window);
void windowGetSize(void *window, int *width, int *height);
int windowGetWidth(void *window);
int windowGetHeight(void *window);
void windowGetRelativePosition(void *window, int *x, int *y);
void windowGetPosition(void *window, int *x, int *y);
void windowSetPosition(void *window, int x, int y);
void windowDestroy(void *window);
void windowSetShadow(void *window, bool hasShadow);
void windowClose(void *window);
void windowZoom(void *window);
void windowRenderHTML(void *window, const char *html);
void windowInjectCSS(void *window, const char *css);
void windowMinimise(void *window);
void windowMaximise(void *window);
bool windowIsVisible(void *window);
void windowSetFullScreen(void *window, bool fullscreen);
void windowUnminimise(void *window);
void windowUnmaximise(void *window);
void windowDisableSizeConstraints(void *window);
void windowShow(void *window);
void windowHide(void *window);
void setButtonState(void *button, int state);
void setMinimiseButtonState(void *window, int state);
void setMaximiseButtonState(void *window, int state);
void setCloseButtonState(void *window, int state);
void windowShowMenu(void *window, void *menu, int x, int y);
void windowSetFrameless(void *window, bool frameless);
void startDrag(void *window);
void windowPrint(void *window);
void windowSetEnabled(void *window, bool enabled);
void windowFocus(void *window);
bool isIgnoreMouseEvents(void *window);
void setIgnoreMouseEvents(void *window, bool ignore);

#endif // WEBVIEW_WINDOW_BINDINGS_DARWIN