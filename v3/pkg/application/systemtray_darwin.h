//go:build darwin

#include <Cocoa/Cocoa.h>

@interface StatusItemController : NSObject
@property long id;
- (void)statusItemClicked:(id)sender;
@end

void* systemTrayNew(long id);
void systemTraySetLabel(void* nsStatusItem, char *label);
void systemTraySetANSILabel(void* nsStatusItem, void* attributedString);
void systemTraySetLabelColor(void* nsStatusItem, char *fg, char *bg);
void* createAttributedString(char *title, char *FG, char *BG);
void* appendAttributedString(void* original, char* label, char* fg, char* bg);
NSImage* imageFromBytes(const unsigned char *bytes, int length);
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate);
void systemTrayDestroy(void* nsStatusItem);
void showMenu(void* nsStatusItem, void *nsMenu);
void systemTrayGetBounds(void* nsStatusItem, NSRect *rect, void **screen);
NSRect NSScreen_frame(void* screen);
void windowSetScreen(void* window, void* screen, int yOffset);
int statusBarHeight();
void systemTrayPositionWindow(void* nsStatusItem, void* nsWindow, int offset);