//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1237
#define EventApplicationDidEnterBackground 1238
#define EventApplicationDidFinishLaunching 1239
#define EventApplicationDidReceiveMemoryWarning 1240
#define EventApplicationWillEnterForeground 1241
#define EventApplicationWillResignActive 1242
#define EventApplicationWillTerminate 1243
#define EventWindowDidLoad 1244
#define EventWindowWillAppear 1245
#define EventWindowDidAppear 1246
#define EventWindowWillDisappear 1247
#define EventWindowDidDisappear 1248
#define EventWindowSafeAreaInsetsChanged 1249
#define EventWindowOrientationChanged 1250
#define EventWindowTouchBegan 1251
#define EventWindowTouchMoved 1252
#define EventWindowTouchEnded 1253
#define EventWindowTouchCancelled 1254
#define EventWebViewDidStartNavigation 1255
#define EventWebViewDidFinishNavigation 1256
#define EventWebViewDidFailNavigation 1257
#define EventWebViewDecidePolicyForNavigationAction 1258

#define MAX_EVENTS 1259


#endif