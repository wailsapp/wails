//go:build linux

#ifndef _events_linux_h
#define _events_linux_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventSystemThemeChanged 1167

#define MAX_EVENTS 2


#endif