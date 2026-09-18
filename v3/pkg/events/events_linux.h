//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationStartup 1052
#define EventSystemDidWake 1053
#define EventSystemThemeChanged 1054
#define EventSystemWillSleep 1055
#define EventWindowDeleteEvent 1056
#define EventWindowDidMove 1057
#define EventWindowDidResize 1058
#define EventWindowFocusIn 1059
#define EventWindowFocusOut 1060
#define EventWindowLoadStarted 1061
#define EventWindowLoadRedirected 1062
#define EventWindowLoadCommitted 1063
#define EventWindowLoadFinished 1064

#define MAX_EVENTS 1065


#endif