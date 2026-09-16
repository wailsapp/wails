//go:build darwin && !ios

#ifndef _DIALOGS_MAC_EXTRAS_H_
#define _DIALOGS_MAC_EXTRAS_H_

#import <Cocoa/Cocoa.h>
#include <stdbool.h>

// Go callbacks (exported from dialogs_mac_extras_darwin.go).
extern void dialogHelpCallback(int callbackID);
extern void saveFileDialogFormatCallback(unsigned int dialogID, int index);
extern void dialogPromptCallback(int callbackID, char *value, bool ok);
extern void dialogColorCallback(int callbackID, int red, int green, int blue, int alpha, bool closed, bool changed);
extern void dialogFontCallback(int callbackID, char *family, char *face, char *postScriptName, double size, bool closed, bool changed);

// Extra configuration for NSSavePanel. Every string is optional (NULL) and,
// when present, is freed by the callee. The joined strings use ';' for
// content types and tags, and '\x1f' (unit separator) for format columns so
// labels may contain any printable character.
typedef struct {
	char *contentTypes;
	char *nameFieldLabel;
	char *tags;
	char *formatLabels;
	char *formatExtensions;
	char *formatIdentifiers;
	int selectedFormat;
} WailsSavePanelExtras;

// NSAlert
void wailsAlertApplyExtras(void *alert, bool showsSuppression, char *suppressionLabel, bool showsHelp, int helpCallbackID, void *accessoryView);
bool wailsAlertSuppressionState(void *alert);

// NSOpenPanel / NSSavePanel. contentTypes is ';' separated and freed here.
// Returns the identifiers that were applied (autoreleased).
NSArray<NSString *> *wailsPanelApplyContentTypes(void *panel, char *contentTypes);
void wailsSavePanelApplyExtras(void *panel, unsigned int dialogID, WailsSavePanelExtras *extras);

// Text prompt alert. All strings are freed here.
void wailsShowPromptDialog(int callbackID, char *title, char *message, char *placeholder, char *defaultValue, bool secure, char *okLabel, char *cancelLabel, void *window);

// Shared colour and font panels. Both return false (and report nothing)
// when a session is already active.
bool wailsShowColorPanel(int callbackID, int red, int green, int blue, int alpha, bool showsAlpha, char *title, bool live);
bool wailsShowFontPanel(int callbackID, char *family, double size, bool live);

#endif
