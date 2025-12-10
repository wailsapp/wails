//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1235
#define EventApplicationDidEnterBackground 1236
#define EventApplicationDidFinishLaunching 1237
#define EventApplicationDidReceiveMemoryWarning 1238
#define EventApplicationWillEnterForeground 1239
#define EventApplicationWillResignActive 1240
#define EventApplicationWillTerminate 1241
#define EventWindowDidLoad 1242
#define EventWindowWillAppear 1243
#define EventWindowDidAppear 1244
#define EventWindowWillDisappear 1245
#define EventWindowDidDisappear 1246
#define EventWindowSafeAreaInsetsChanged 1247
#define EventWindowOrientationChanged 1248
#define EventWindowTouchBegan 1249
#define EventWindowTouchMoved 1250
#define EventWindowTouchEnded 1251
#define EventWindowTouchCancelled 1252
#define EventWebViewDidStartNavigation 1253
#define EventWebViewDidFinishNavigation 1254
#define EventWebViewDidFailNavigation 1255
#define EventWebViewDecidePolicyForNavigationAction 1256

#define MAX_EVENTS 1257


#endif