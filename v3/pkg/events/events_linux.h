//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationStartup 1048
#define EventSystemThemeChanged 1049
#define EventWindowDeleteEvent 1050
#define EventWindowDidMove 1051
#define EventWindowDidResize 1052
#define EventWindowFocusIn 1053
#define EventWindowFocusOut 1054
#define EventWindowLoadChanged 1055

#define MAX_EVENTS 1056


#endif