//go:build darwin

@interface StatusItemController : NSObject
@property long id;
- (void)statusItemClicked:(id)sender;
@end

void* systemTrayNew(long id);
void systemTraySetLabel(void* nsStatusItem, char *label);
NSImage* imageFromBytes(const unsigned char *bytes, int length);
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate);
void systemTraySetMenu(void* nsStatusItem, void* nsMenu);
void systemTrayDestroy(void* nsStatusItem);
void showMenu(void* nsStatusItem);
void systemTrayGetBounds(void* nsStatusItem, NSRect *rect);