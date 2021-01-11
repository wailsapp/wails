#ifndef __FFENESTRI_H__
#define __FFENESTRI_H__

#include <stdio.h>

extern void *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel);
extern void SetMinWindowSize(void *app, int minWidth, int minHeight);
extern void SetMaxWindowSize(void *app, int maxWidth, int maxHeight);
extern void Run(void *app, int argc, char **argv);
extern void DestroyApplication(void *app);
extern void SetDebug(void *app, int flag);
extern void SetBindings(void *app, const char *bindings);
extern void ExecJS(void *app, const char *script);
extern void Hide(void *app);
extern void Show(void *app);
extern void Center(void *app);
extern void Maximise(void *app);
extern void Unmaximise(void *app);
extern void ToggleMaximise(void *app);
extern void Minimise(void *app);
extern void Unminimise(void *app);
extern void ToggleMinimise(void *app);
extern void SetColour(void *app, int red, int green, int blue, int alpha);
extern void SetSize(void *app, int width, int height);
extern void SetPosition(void *app, int x, int y);
extern void Quit(void *app);
extern void SetTitle(void *app, const char *title);
extern void Fullscreen(void *app);
extern void UnFullscreen(void *app);
extern void ToggleFullscreen(void *app);
extern void DisableFrame(void *app);
extern void OpenDialog(void *appPointer, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolvesAliases, int treatPackagesAsDirectories);
extern void SaveDialog(void *appPointer, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories);
extern void MessageDialog(void *appPointer, char *callbackID, char *type, char *title, char *message, char *icon, char *button1, char *button2, char *button3, char *button4, char *defaultButton, char *cancelButton);
extern void DarkModeEnabled(void *appPointer, char *callbackID);
extern void SetApplicationMenu(void *, const char *);
extern void UpdateTray(void *app, char *menuAsJSON);
extern void UpdateContextMenus(void *app, char *contextMenusAsJSON);
extern void UpdateTrayLabel(void *app, const char *label);
extern void UpdateTrayIcon(void *app, const char *label);

#endif
