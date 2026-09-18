//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

// Result codes shared with sound_manager_darwin.go.
enum {
	WailsSoundOK = 0,
	WailsSoundNotFound = 1,
	WailsSoundPlayFailed = 2,
};

void soundBeep(void);
int soundPlayNamed(const char *name);
int soundPlayFile(const char *path);
int soundPlayData(const void *bytes, int length);
