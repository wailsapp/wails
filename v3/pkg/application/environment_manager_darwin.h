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

#endif /* environment_manager_darwin_h */
