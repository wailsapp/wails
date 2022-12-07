//go:build darwin

#ifndef _events_h
#define _events_h

extern void systemEventHandler(char*);

#define EventApplicationDidFinishLaunching "mac:ApplicationDidFinishLaunching"
#define EventApplicationWillTerminate "mac:ApplicationWillTerminate"


#endif