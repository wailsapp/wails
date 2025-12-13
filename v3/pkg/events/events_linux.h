//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationStartup 1051
#define EventSystemThemeChanged 1052
#define EventWindowDeleteEvent 1053
#define EventWindowDidMove 1054
#define EventWindowDidResize 1055
#define EventWindowFocusIn 1056
#define EventWindowFocusOut 1057
#define EventWindowLoadStarted 1058
#define EventWindowLoadRedirected 1059
#define EventWindowLoadCommitted 1060
#define EventWindowLoadFinished 1061

#define MAX_EVENTS 1062


#endif