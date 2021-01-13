
#ifndef FFENESTRI_DARWIN_H
#define FFENESTRI_DARWIN_H


#define OBJC_OLD_DISPATCH_PROTOTYPES 1
#include <objc/objc-runtime.h>
#include <CoreGraphics/CoreGraphics.h>
#include "json.h"
#include "hashmap.h"
#include "stdlib.h"

// Macros to make it slightly more sane
#define msg objc_msgSend

#define c(str) (id)objc_getClass(str)
#define s(str) sel_registerName(str)
#define u(str) sel_getUid(str)
#define str(input) msg(c("NSString"), s("stringWithUTF8String:"), input)
#define strunicode(input) msg(c("NSString"), s("stringWithFormat:"), str("%C"), (unsigned short)input)
#define cstr(input) (const char *)msg(input, s("UTF8String"))
#define url(input) msg(c("NSURL"), s("fileURLWithPath:"), str(input))

#define ALLOC(classname) msg(c(classname), s("alloc"))
#define ALLOC_INIT(classname) msg(msg(c(classname), s("alloc")), s("init"))
#define GET_FRAME(receiver) ((CGRect(*)(id, SEL))objc_msgSend_stret)(receiver, s("frame"))
#define GET_BOUNDS(receiver) ((CGRect(*)(id, SEL))objc_msgSend_stret)(receiver, s("bounds"))
#define GET_BACKINGSCALEFACTOR(receiver) ((CGFloat(*)(id, SEL))msg)(receiver, s("backingScaleFactor"))

#define ON_MAIN_THREAD(str) dispatch( ^{ str; } )
#define MAIN_WINDOW_CALL(str) msg(app->mainWindow, s((str)))

#define NSBackingStoreBuffered 2

#define NSWindowStyleMaskBorderless 0
#define NSWindowStyleMaskTitled 1
#define NSWindowStyleMaskClosable 2
#define NSWindowStyleMaskMiniaturizable 4
#define NSWindowStyleMaskResizable 8
#define NSWindowStyleMaskFullscreen 1 << 14

#define NSVisualEffectMaterialWindowBackground 12
#define NSVisualEffectBlendingModeBehindWindow 0
#define NSVisualEffectStateFollowsWindowActiveState 0
#define NSVisualEffectStateActive 1
#define NSVisualEffectStateInactive 2

#define NSViewWidthSizable 2
#define NSViewHeightSizable 16

#define NSWindowBelow -1
#define NSWindowAbove 1

#define NSSquareStatusItemLength   -2.0
#define NSVariableStatusItemLength -1.0

#define NSWindowTitleHidden 1
#define NSWindowStyleMaskFullSizeContentView 1 << 15

#define NSEventModifierFlagCommand 1 << 20
#define NSEventModifierFlagOption 1 << 19
#define NSEventModifierFlagControl 1 << 18
#define NSEventModifierFlagShift 1 << 17

#define NSControlStateValueMixed -1
#define NSControlStateValueOff 0
#define NSControlStateValueOn 1

// Unbelievably, if the user swaps their button preference
// then right buttons are reported as left buttons
#define NSEventMaskLeftMouseDown 1 << 1
#define NSEventMaskLeftMouseUp 1 << 2
#define NSEventMaskRightMouseDown 1 << 3
#define NSEventMaskRightMouseUp 1 << 4

#define NSEventTypeLeftMouseDown 1
#define NSEventTypeLeftMouseUp 2
#define NSEventTypeRightMouseDown 3
#define NSEventTypeRightMouseUp 4

#define NSNoImage       0
#define NSImageOnly     1
#define NSImageLeft     2
#define NSImageRight    3
#define NSImageBelow    4
#define NSImageAbove    5
#define NSImageOverlaps 6

#define NSAlertStyleWarning 0
#define NSAlertStyleInformational 1
#define NSAlertStyleCritical 2

#define NSAlertFirstButtonReturn   1000
#define NSAlertSecondButtonReturn  1001
#define NSAlertThirdButtonReturn   1002

struct Application;
int releaseNSObject(void *const context, struct hashmap_element_s *const e);
void TitlebarAppearsTransparent(struct Application* app);
void HideTitle(struct Application* app);
void HideTitleBar(struct Application* app);
void FullSizeContent(struct Application* app);
void UseToolbar(struct Application* app);
void HideToolbarSeparator(struct Application* app);
void DisableFrame(struct Application* app);
void SetAppearance(struct Application* app, const char *);
void WebviewIsTransparent(struct Application* app);
void WindowBackgroundIsTranslucent(struct Application* app);
void SetTray(struct Application* app, const char *, const char *, const char *);
void SetContextMenus(struct Application* app, const char *);
void AddTrayMenu(struct Application* app, const char *);

#endif