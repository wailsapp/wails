//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1243
#define EventApplicationDidEnterBackground 1244
#define EventApplicationDidFinishLaunching 1245
#define EventApplicationDidReceiveMemoryWarning 1246
#define EventApplicationWillEnterForeground 1247
#define EventApplicationWillResignActive 1248
#define EventApplicationWillTerminate 1249
#define EventWindowDidLoad 1250
#define EventWindowWillAppear 1251
#define EventWindowDidAppear 1252
#define EventWindowWillDisappear 1253
#define EventWindowDidDisappear 1254
#define EventWindowSafeAreaInsetsChanged 1255
#define EventWindowOrientationChanged 1256
#define EventWindowTouchBegan 1257
#define EventWindowTouchMoved 1258
#define EventWindowTouchEnded 1259
#define EventWindowTouchCancelled 1260
#define EventWebViewDidStartNavigation 1261
#define EventWebViewDidFinishNavigation 1262
#define EventWebViewDidFailNavigation 1263
#define EventWebViewDecidePolicyForNavigationAction 1264

#define MAX_EVENTS 1265


#endif