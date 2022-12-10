//go:build darwin

#ifndef _events_h
#define _events_h

extern void systemEventHandler(char*);

#define EventApplicationDidFinishLaunching "mac:ApplicationDidFinishLaunching"
#define EventApplicationWillTerminate "mac:ApplicationWillTerminate"
#define EventApplicationDidBecomeActive "mac:ApplicationDidBecomeActive"
#define EventApplicationWillUpdate "mac:ApplicationWillUpdate"
#define EventApplicationDidUpdate "mac:ApplicationDidUpdate"
#define EventApplicationWillFinishLaunching "mac:ApplicationWillFinishLaunching"
#define EventApplicationWillHide "mac:ApplicationWillHide"
#define EventApplicationWillUnhide "mac:ApplicationWillUnhide"
#define EventApplicationDidHide "mac:ApplicationDidHide"
#define EventApplicationDidUnhide "mac:ApplicationDidUnhide"


#endif