//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationStartup 1049
#define EventSystemThemeChanged 1050
#define EventWindowDeleteEvent 1051
#define EventWindowDidMove 1052
#define EventWindowDidResize 1053
#define EventWindowFocusIn 1054
#define EventWindowFocusOut 1055
#define EventWindowLoadChanged 1056

#define MAX_EVENTS 1057


#endif