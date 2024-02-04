//go:build linux

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventSystemThemeChanged 1024
#define EventWindowLoadChanged 1025
#define EventWindowDeleteEvent 1026
#define EventApplicationStartup 1027

#define MAX_EVENTS 1028


#endif