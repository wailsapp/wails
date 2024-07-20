//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventSystemThemeChanged 1024
#define EventWindowLoadChanged 1025
#define EventWindowDeleteEvent 1026
#define EventWindowDidMove 1027
#define EventWindowDidResize 1028
#define EventWindowFocusIn 1029
#define EventWindowFocusOut 1030
#define EventWindowDragDrop 1031
#define EventWindowDragBegin 1032
#define EventWindowDragEnd 1033
#define EventWindowDragLeave 1034
#define EventApplicationStartup 1035

#define MAX_EVENTS 1036


#endif