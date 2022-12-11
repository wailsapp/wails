//go:build darwin

#ifndef _events_h
#define _events_h

extern void systemEventHandler(char*);

#define EventApplicationWillFinishLaunching "mac:ApplicationWillFinishLaunching"
#define EventApplicationDidFinishLaunching "mac:ApplicationDidFinishLaunching"
#define EventApplicationWillBecomeActive "mac:ApplicationWillBecomeActive"
#define EventApplicationDidBecomeActive "mac:ApplicationDidBecomeActive"
#define EventApplicationWillUpdate "mac:ApplicationWillUpdate"
#define EventApplicationDidUpdate "mac:ApplicationDidUpdate"
#define EventApplicationWillHide "mac:ApplicationWillHide"
#define EventApplicationDidHide "mac:ApplicationDidHide"
#define EventApplicationWillUnhide "mac:ApplicationWillUnhide"
#define EventApplicationDidUnhide "mac:ApplicationDidUnhide"
#define EventApplicationWillResignActive "mac:ApplicationWillResignActive"
#define EventApplicationDidResignActive "mac:ApplicationDidResignActive"
#define EventApplicationWillTerminate "mac:ApplicationWillTerminate"
#define EventApplicationDidChangeOcclusionState "mac:ApplicationDidChangeOcclusionState"
#define EventApplicationDidChangeScreenParameters "mac:ApplicationDidChangeScreenParameters"
#define EventApplicationDidChangeBackingProperties "mac:ApplicationDidChangeBackingProperties"
#define EventApplicationDidChangeIcon "mac:ApplicationDidChangeIcon"
#define EventApplicationDidChangeStatusBarOrientation "mac:ApplicationDidChangeStatusBarOrientation"
#define EventApplicationDidChangeStatusBarFrame "mac:ApplicationDidChangeStatusBarFrame"
#define EventApplicationDidChangeEffectiveAppearance "mac:ApplicationDidChangeEffectiveAppearance"


#endif