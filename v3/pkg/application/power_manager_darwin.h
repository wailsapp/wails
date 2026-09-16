//go:build darwin && !ios && !server

#ifndef power_manager_darwin_h
#define power_manager_darwin_h

#include <stdbool.h>

typedef struct {
    bool lowPowerMode;
    // 0 unknown, 1 nominal, 2 fair, 3 serious, 4 critical (ThermalState in Go).
    int thermalState;
    bool onBattery;
    // 0..1, or -1 when unknown or no battery.
    double batteryLevel;
    bool charging;
} WailsPowerState;

// wailsPowerBeginActivity starts an NSProcessInfo activity that disables
// idle system sleep (and idle display sleep when display is set). The
// returned token is retained and must be passed to wailsPowerEndActivity.
void* wailsPowerBeginActivity(const char* reason, bool display);
void wailsPowerEndActivity(void* token);

WailsPowerState wailsPowerState(void);

// wailsPowerObserverStart registers for NSProcessInfo power and thermal
// notifications; they are forwarded to Go through powerNotification.
void wailsPowerObserverStart(void);

#endif /* power_manager_darwin_h */
