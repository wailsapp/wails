//go:build darwin && !ios && !server

#ifndef environment_manager_darwin_h
#define environment_manager_darwin_h

#include <stdbool.h>

typedef struct {
    bool reduceMotion;
    bool reduceTransparency;
    bool increaseContrast;
    bool differentiateWithoutColor;
    bool invertColors;
    bool voiceOverEnabled;
    bool switchControlEnabled;
} WailsAccessibilitySettings;

WailsAccessibilitySettings wailsAccessibilitySettings(void);

// wailsKeyboardLayoutJSON returns {"id","name","languages":[...]} for the
// current keyboard input source. Caller frees. Call on the main thread.
char* wailsKeyboardLayoutJSON(void);

// wailsLocaleJSON returns {"identifier","language","region","preferred":[...]}
// for NSLocale currentLocale. Caller frees.
char* wailsLocaleJSON(void);

// wailsEnvironmentObserverStart registers for accessibility, keyboard layout
// and locale change notifications; they are forwarded to Go through
// environmentNotification.
void wailsEnvironmentObserverStart(void);

// wailsWorkspaceSetDefaultHandler makes this bundle the default application
// for contentType (a UTI) or scheme; pass NULL for the one not used. The
// outcome arrives through wailsCompletionCallback(id, message) with a NULL
// message on success. Call on the main thread.
void wailsWorkspaceSetDefaultHandler(unsigned long long id, const char *contentType, const char *scheme);

// wailsWorkspaceDefaultHandlerJSON returns {"name","bundleID","path"} for
// the application that opens contentType or scheme (pass NULL for the one
// not used), or "null" when nothing is registered. Caller frees.
char* wailsWorkspaceDefaultHandlerJSON(const char *contentType, const char *scheme);

#endif /* environment_manager_darwin_h */
