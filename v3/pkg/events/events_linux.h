//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationStartup 1050
#define EventSystemDidWake 1051
#define EventSystemThemeChanged 1052
#define EventSystemWillSleep 1053
#define EventWindowDeleteEvent 1054
#define EventWindowDidMove 1055
#define EventWindowDidResize 1056
#define EventWindowFocusIn 1057
#define EventWindowFocusOut 1058
#define EventWindowLoadStarted 1059
#define EventWindowLoadRedirected 1060
#define EventWindowLoadCommitted 1061
#define EventWindowLoadFinished 1062

#define MAX_EVENTS 1063


#endif