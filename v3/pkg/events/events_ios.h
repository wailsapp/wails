//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1241
#define EventApplicationDidEnterBackground 1242
#define EventApplicationDidFinishLaunching 1243
#define EventApplicationDidReceiveMemoryWarning 1244
#define EventApplicationWillEnterForeground 1245
#define EventApplicationWillResignActive 1246
#define EventApplicationWillTerminate 1247
#define EventWindowDidLoad 1248
#define EventWindowWillAppear 1249
#define EventWindowDidAppear 1250
#define EventWindowWillDisappear 1251
#define EventWindowDidDisappear 1252
#define EventWindowSafeAreaInsetsChanged 1253
#define EventWindowOrientationChanged 1254
#define EventWindowTouchBegan 1255
#define EventWindowTouchMoved 1256
#define EventWindowTouchEnded 1257
#define EventWindowTouchCancelled 1258
#define EventWebViewDidStartNavigation 1259
#define EventWebViewDidFinishNavigation 1260
#define EventWebViewDidFailNavigation 1261
#define EventWebViewDecidePolicyForNavigationAction 1262

#define MAX_EVENTS 1263


#endif