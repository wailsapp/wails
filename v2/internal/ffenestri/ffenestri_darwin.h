

#ifndef FFENESTRI_DARWIN_H
#define FFENESTRI_DARWIN_H

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
extern void HideTitleBar(void *);
extern void FullSizeContent(void *);
extern void UseToolbar(void *);
extern void HideToolbarSeparator(void *);
extern void DisableFrame(void *);
extern void SetAppearance(void *, const char *);
extern void WebviewIsTransparent(void *);
extern void WindowBackgroundIsTranslucent(void *);
extern void SetMenu(void *, const char *);
extern void SetTray(void *, const char *);
extern void SetContextMenus(void *, const char *);

#endif