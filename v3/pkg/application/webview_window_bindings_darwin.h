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

extern void registerListener(unsigned int event);

void *windowNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences);
void *panelNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences);
void *windowOrPanelNew(bool isWindow, unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences);
void printWindowStyle(void *window);
void setInvisibleTitleBarHeight(void *window, unsigned int height);
void windowSetTransparent(void *nsWindow);
void windowSetInvisibleTitleBar(void *nsWindow, unsigned int height);
void windowSetTitle(void *nsWindow, char *title);
void windowSetSize(void *nsWindow, int width, int height);
void windowSetAlwaysOnTop(void *nsWindow, bool alwaysOnTop);
void setNormalWindowLevel(void *nsWindow);
void setFloatingWindowLevel(void *nsWindow);
void setPopUpMenuWindowLevel(void *nsWindow);
void setMainMenuWindowLevel(void *nsWindow);
void setStatusWindowLevel(void *nsWindow);
void setModalPanelWindowLevel(void *nsWindow);
void setScreenSaverWindowLevel(void *nsWindow);
void setTornOffMenuWindowLevel(void *nsWindow);
void navigationLoadURL(void *nsWindow, char *url);
void windowSetResizable(void *nsWindow, bool resizable);
void windowSetMinSize(void *nsWindow, int width, int height);
void windowSetMaxSize(void *nsWindow, int width, int height);
void windowZoomReset(void *nsWindow);
void windowZoomSet(void *nsWindow, double zoom);
float windowZoomGet(void *nsWindow);
void windowZoomIn(void *nsWindow);
void windowZoomOut(void *nsWindow);
void windowSetRelativePosition(void *nsWindow, int x, int y);
void windowExecJS(void *nsWindow, const char *js);
void windowSetTranslucent(void *nsWindow);
void webviewSetTransparent(void *nsWindow);
void webviewSetBackgroundColour(void *nsWindow, int r, int g, int b, int alpha);
void windowSetBackgroundColour(void *nsWindow, int r, int g, int b, int alpha);
bool windowIsMaximised(void *nsWindow);
bool windowIsFullscreen(void *nsWindow);
bool windowIsMinimised(void *nsWindow);
bool windowIsFocused(void *nsWindow);
void windowFullscreen(void *nsWindow);
void windowUnFullscreen(void *nsWindow);
void windowRestore(void *nsWindow);
void setFullscreenButtonEnabled(void *nsWindow, bool enabled);
void windowSetTitleBarAppearsTransparent(void *nsWindow, bool transparent);
void windowSetFullSizeContent(void *nsWindow, bool fullSize);
void windowSetHideTitleBar(void *nsWindow, bool hideTitlebar);
void windowSetHideTitle(void *nsWindow, bool hideTitle);
void windowSetUseToolbar(void *nsWindow, bool useToolbar);
void windowSetToolbarStyle(void *nsWindow, int style);
void windowSetHideToolbarSeparator(void *nsWindow, bool hideSeparator);
void windowSetShowToolbarWhenFullscreen(void *window, bool setting);
void windowSetAppearanceTypeByName(void *nsWindow, const char *appearanceName);
void windowCenter(void *nsWindow);
void windowGetSize(void *nsWindow, int *width, int *height);
int windowGetWidth(void *nsWindow);
int windowGetHeight(void *nsWindow);
void windowGetRelativePosition(void *nsWindow, int *x, int *y);
void windowGetPosition(void *nsWindow, int *x, int *y);
void windowSetPosition(void *nsWindow, int x, int y);
void windowDestroy(void *nsWindow);
void windowSetShadow(void *nsWindow, bool hasShadow);
void windowClose(void *window);
void windowZoom(void *window);
void windowRenderHTML(void *window, const char *html);
void windowInjectCSS(void *window, const char *css);
void windowMinimise(void *window);
void windowMaximise(void *window);
bool isFullScreen(void *window);
bool isVisible(void *window);
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
void setWindowEnabled(void *window, bool enabled);
void windowSetEnabled(void *window, bool enabled);
void windowFocus(void *window);
bool isIgnoreMouseEvents(void *nsWindow);
void setIgnoreMouseEvents(void *nsWindow, bool ignore);

#endif // WEBVIEW_WINDOW_BINDINGS_DARWIN