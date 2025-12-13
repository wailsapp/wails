//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1239
#define EventApplicationDidEnterBackground 1240
#define EventApplicationDidFinishLaunching 1241
#define EventApplicationDidReceiveMemoryWarning 1242
#define EventApplicationWillEnterForeground 1243
#define EventApplicationWillResignActive 1244
#define EventApplicationWillTerminate 1245
#define EventWindowDidLoad 1246
#define EventWindowWillAppear 1247
#define EventWindowDidAppear 1248
#define EventWindowWillDisappear 1249
#define EventWindowDidDisappear 1250
#define EventWindowSafeAreaInsetsChanged 1251
#define EventWindowOrientationChanged 1252
#define EventWindowTouchBegan 1253
#define EventWindowTouchMoved 1254
#define EventWindowTouchEnded 1255
#define EventWindowTouchCancelled 1256
#define EventWebViewDidStartNavigation 1257
#define EventWebViewDidFinishNavigation 1258
#define EventWebViewDidFailNavigation 1259
#define EventWebViewDecidePolicyForNavigationAction 1260

#define MAX_EVENTS 1261


#endif