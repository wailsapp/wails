//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1245
#define EventApplicationDidEnterBackground 1246
#define EventApplicationDidFinishLaunching 1247
#define EventApplicationDidReceiveMemoryWarning 1248
#define EventApplicationWillEnterForeground 1249
#define EventApplicationWillResignActive 1250
#define EventApplicationWillTerminate 1251
#define EventWindowDidLoad 1252
#define EventWindowWillAppear 1253
#define EventWindowDidAppear 1254
#define EventWindowWillDisappear 1255
#define EventWindowDidDisappear 1256
#define EventWindowSafeAreaInsetsChanged 1257
#define EventWindowOrientationChanged 1258
#define EventWindowTouchBegan 1259
#define EventWindowTouchMoved 1260
#define EventWindowTouchEnded 1261
#define EventWindowTouchCancelled 1262
#define EventWebViewDidStartNavigation 1263
#define EventWebViewDidFinishNavigation 1264
#define EventWebViewDidFailNavigation 1265
#define EventWebViewDecidePolicyForNavigationAction 1266

#define MAX_EVENTS 1267


#endif