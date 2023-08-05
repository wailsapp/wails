//go:build darwin

@interface StatusItemController : NSObject
@property long id;
- (void)statusItemClicked:(id)sender;
@end

void* systemTrayNew(long id);
void systemTraySetLabel(void* nsStatusItem, char *label);
void systemTraySetANSILabel(void* nsStatusItem, char *label, char *FG, char *BG);
NSMutableAttributedString* createAttributedString(char *title, char *FG, char *BG);
NSImage* imageFromBytes(const unsigned char *bytes, int length);
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate);
void systemTraySetMenu(void* nsStatusItem, void* nsMenu);
void systemTrayDestroy(void* nsStatusItem);
void showMenu(void* nsStatusItem, void *nsMenu);
void systemTrayGetBounds(void* nsStatusItem, NSRect *rect);
int statusBarHeight();