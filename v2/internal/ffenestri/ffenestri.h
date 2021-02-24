#ifndef __FFENESTRI_H__
#define __FFENESTRI_H__

#include <stdio.h>
struct Application;

extern struct Application *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel, int hideWindowOnClose);
extern void SetMinWindowSize(struct Application*, int minWidth, int minHeight);
extern void SetMaxWindowSize(struct Application*, int maxWidth, int maxHeight);
extern void Run(struct Application*, int argc, char **argv);
extern void DestroyApplication(struct Application*);
extern void SetDebug(struct Application*, int flag);
extern void SetBindings(struct Application*, const char *bindings);
extern void ExecJS(struct Application*, const char *script);
extern void Hide(struct Application*);
extern void Show(struct Application*);
extern void Center(struct Application*);
extern void Maximise(struct Application*);
extern void Unmaximise(struct Application*);
extern void ToggleMaximise(struct Application*);
extern void Minimise(struct Application*);
extern void Unminimise(struct Application*);
extern void ToggleMinimise(struct Application*);
extern void SetColour(struct Application*, int red, int green, int blue, int alpha);
extern void SetSize(struct Application*, int width, int height);
extern void SetPosition(struct Application*, int x, int y);
extern void Quit(struct Application*);
extern void SetTitle(struct Application*, const char *title);
extern void Fullscreen(struct Application*);
extern void UnFullscreen(struct Application*);
extern void ToggleFullscreen(struct Application*);
extern void DisableFrame(struct Application*);
extern void OpenDialog(struct Application*, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolvesAliases, int treatPackagesAsDirectories);
extern void SaveDialog(struct Application*, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories);
extern void MessageDialog(struct Application*, char *callbackID, char *type, char *title, char *message, char *icon, char *button1, char *button2, char *button3, char *button4, char *defaultButton, char *cancelButton);
extern void DarkModeEnabled(struct Application*, char *callbackID);
extern void SetApplicationMenu(struct Application*, const char *);
extern void AddTrayMenu(struct Application*, const char *menuTrayJSON);
extern void SetTrayMenu(struct Application*, const char *menuTrayJSON);
extern void DeleteTrayMenuByID(struct Application*, const char *id);
extern void UpdateTrayMenuLabel(struct Application*, const char* JSON);
extern void AddContextMenu(struct Application*, char *contextMenuJSON);
extern void UpdateContextMenu(struct Application*, char *contextMenuJSON);

#endif
