
#ifdef FFENESTRI_DARWIN

#define OBJC_OLD_DISPATCH_PROTOTYPES 1
#include <objc/objc-runtime.h>
#include <CoreGraphics/CoreGraphics.h>
#include "json.h"
#include "hashmap.h"

// Macros to make it slightly more sane
#define msg objc_msgSend

#define c(str) (id)objc_getClass(str)
#define s(str) sel_registerName(str)
#define u(str) sel_getUid(str)
#define str(input) msg(c("NSString"), s("stringWithUTF8String:"), input)
#define strunicode(input) msg(c("NSString"), s("stringWithFormat:"), str("%C"), (unsigned short)input)
#define cstr(input) (const char *)msg(input, s("UTF8String"))
#define url(input) msg(c("NSURL"), s("fileURLWithPath:"), str(input))

#define ALLOC(classname) msg(c(classname), s("alloc"))
#define GET_FRAME(receiver) ((CGRect(*)(id, SEL))objc_msgSend_stret)(receiver, s("frame"))
#define GET_BOUNDS(receiver) ((CGRect(*)(id, SEL))objc_msgSend_stret)(receiver, s("bounds"))

#define STREQ(a,b) strcmp(a, b) == 0
#define STRCOPY(a) concat(a, "")
#define MEMFREE(input) free((void*)input); input = NULL;

#define ON_MAIN_THREAD(str) dispatch( ^{ str; } )
#define MAIN_WINDOW_CALL(str) msg(app->mainWindow, s((str)))

#define NSBackingStoreBuffered 2

#define NSWindowStyleMaskBorderless 0
#define NSWindowStyleMaskTitled 1
#define NSWindowStyleMaskClosable 2
#define NSWindowStyleMaskMiniaturizable 4
#define NSWindowStyleMaskResizable 8
#define NSWindowStyleMaskFullscreen 1 << 14

#define NSVisualEffectMaterialWindowBackground 12
#define NSVisualEffectBlendingModeBehindWindow 0
#define NSVisualEffectStateFollowsWindowActiveState 0
#define NSVisualEffectStateActive 1
#define NSVisualEffectStateInactive 2

#define NSViewWidthSizable 2
#define NSViewHeightSizable 16

#define NSWindowBelow -1
#define NSWindowAbove 1

#define NSWindowTitleHidden 1
#define NSWindowStyleMaskFullSizeContentView 1 << 15

#define NSEventModifierFlagCommand 1 << 20
#define NSEventModifierFlagOption 1 << 19
#define NSEventModifierFlagControl 1 << 18
#define NSEventModifierFlagShift 1 << 17

#define NSControlStateValueMixed -1
#define NSControlStateValueOff 0
#define NSControlStateValueOn 1

// Unbelievably, if the user swaps their button preference
// then right buttons are reported as left buttons
#define NSEventMaskLeftMouseDown 1 << 1
#define NSEventMaskLeftMouseUp 1 << 2
#define NSEventMaskRightMouseDown 1 << 3
#define NSEventMaskRightMouseUp 1 << 4

#define NSEventTypeLeftMouseDown 1
#define NSEventTypeLeftMouseUp 2
#define NSEventTypeRightMouseDown 3
#define NSEventTypeRightMouseUp 4


// References to assets
extern const unsigned char *assets[];
extern const unsigned char runtime;
extern const char *icon[];

// Tray icon
extern const unsigned int trayIconLength;
extern const unsigned char *trayIcon[];

// MAIN DEBUG FLAG
int debug;

// MenuItem map for the application menu
struct hashmap_s menuItemMapForApplicationMenu;

// RadioGroup map for the application menu. Maps a menuitem id with its associated radio group items
struct hashmap_s radioGroupMapForApplicationMenu;

// MenuItem map for the tray menu
struct hashmap_s menuItemMapForTrayMenu;

// RadioGroup map for the tray menu. Maps a menuitem id with its associated radio group items
struct hashmap_s radioGroupMapForTrayMenu;

// contextMenuMap is a hashmap of context menus keyed on a string ID
struct hashmap_s contextMenuMap;

// MenuItem map for the context menus
struct hashmap_s menuItemMapForContextMenus;

// RadioGroup map for the context menus. Maps a menuitem id with its associated radio group items
struct hashmap_s radioGroupMapForContextMenus;

// Context menu data is given by the frontend when clicking a context menu.
// We send this to the backend when an item is selected;
const char *contextMenuData;

// Dispatch Method
typedef void (^dispatchMethod)(void);

// dispatch will execute the given `func` pointer
void dispatch(dispatchMethod func) {
	dispatch_async(dispatch_get_main_queue(), func);
}

// App Delegate
typedef struct AppDel {
	Class isa;
	id window;
} AppDelegate;

// Credit: https://stackoverflow.com/a/8465083
char* concat(const char *string1, const char *string2)
{
	const size_t len1 = strlen(string1);
	const size_t len2 = strlen(string2);
	char *result = malloc(len1 + len2 + 1);
	memcpy(result, string1, len1);
	memcpy(result + len1, string2, len2 + 1);
	return result;
}

// yes command simply returns YES!
BOOL yes(id self, SEL cmd)
{
	return YES;
}

// Prints a hashmap entry
int hashmap_log(void *const context, struct hashmap_element_s *const e) {
  printf("%s: %p ", (char*)e->key, e->data);
  return 0;
}

// Utility function to visualise a hashmap
void dumpHashmap(const char *name, struct hashmap_s *hashmap) {
  printf("%s = { ", name);
  if (0!=hashmap_iterate_pairs(hashmap, hashmap_log, NULL)) {
	fprintf(stderr, "Failed to dump hashmap entries\n");
  }
  printf("}\n");
}

extern void messageFromWindowCallback(const char *);
typedef void (*ffenestriCallback)(const char *);

void HideMouse() {
	msg(c("NSCursor"), s("hide"));
}

void ShowMouse() {
	msg(c("NSCursor"), s("unhide"));
}

struct Application {

	// Cocoa data
	id application;
	id delegate;
	id mainWindow;
	id wkwebview;
	id manager;
	id config;
	id mouseEvent;
	id mouseDownMonitor;
	id mouseUpMonitor;
	id vibrancyLayer;


	// Window Data
	const char *title;
	int width;
	int height;
	int minWidth;
	int minHeight;
	int maxWidth;
	int maxHeight;
	int resizable;
	int devtools;
	int fullscreen;
	int red;
	int green;
	int blue;
	int alpha;
	int webviewIsTranparent;
	const char *appearance;
	int decorations;
	int logLevel;

	// Features
	int frame;
	int startHidden;
	int maximised;
	int minimised;
	int titlebarAppearsTransparent;
	int hideTitle;
	int hideTitleBar;
	int fullSizeContent;
	int useToolBar;
	int hideToolbarSeparator;
	int windowBackgroundIsTranslucent;

	// Menu
	const char *menuAsJSON;
	id menubar;
	JsonNode *processedMenu;

	// Tray
	const char *trayMenuAsJSON;
	JsonNode *processedTrayMenu;
	id statusItem;

	// Context Menus
	const char *contextMenusAsJSON;
	JsonNode *processedContextMenus;

	// User Data
	char *HTML;

	// Callback
	ffenestriCallback sendMessageToBackend;

	// Bindings
	const char *bindings;

	// Lock - used for sync operations (Should we be using g_mutex?)
	int lock;

};

// Debug works like sprintf but mutes if the global debug flag is true
// Credit: https://stackoverflow.com/a/20639708

// 5k is more than enough for a log message
#define MAXMESSAGE 1024*10
char logbuffer[MAXMESSAGE];
void Debug(struct Application *app, const char *message, ... ) {
	if ( debug ) {
		const char *temp = concat("LTFfenestri (C) | ", message);
		va_list args;
		va_start(args, message);
		vsnprintf(logbuffer, MAXMESSAGE, temp, args);
		app->sendMessageToBackend(&logbuffer[0]);
		MEMFREE(temp);
		va_end(args);
	}
}

void Fatal(struct Application *app, const char *message, ... ) {
  const char *temp = concat("LFFfenestri (C) | ", message);
  va_list args;
  va_start(args, message);
  vsnprintf(logbuffer, MAXMESSAGE, temp, args);
  app->sendMessageToBackend(&logbuffer[0]);
  MEMFREE(temp);
  va_end(args);
}

void TitlebarAppearsTransparent(struct Application* app) {
	app->titlebarAppearsTransparent = 1;
}

void HideTitle(struct Application *app) {
	app->hideTitle = 1;
}

void HideTitleBar(struct Application *app) {
	app->hideTitleBar = 1;
}

void HideToolbarSeparator(struct Application *app) {
	app->hideToolbarSeparator = 1;
}

void UseToolbar(struct Application *app) {
	app->useToolBar = 1;
}

// WebviewIsTransparent will make the webview transparent
// revealing the Cocoa window underneath
void WebviewIsTransparent(struct Application *app) {
	app->webviewIsTranparent = 1;
}

// SetAppearance will set the window's Appearance to the
// given value
void SetAppearance(struct Application *app, const char *appearance) {
	app->appearance = appearance;
}


void applyWindowColour(struct Application *app) {
	// Apply the colour only if the window has been created
	if( app->mainWindow != NULL ) {
		ON_MAIN_THREAD(
			id colour = msg(c("NSColor"), s("colorWithCalibratedRed:green:blue:alpha:"),
								(float)app->red / 255.0,
								(float)app->green / 255.0,
								(float)app->blue / 255.0,
								(float)app->alpha / 255.0);
			msg(app->mainWindow, s("setBackgroundColor:"), colour);
		);
	}
}

void showContextMenu(struct Application *app, const char *contextMenuID) {

	// If no context menu ID was given
	if( contextMenuID == NULL ) {
		// Show default context menu if we have one
		return;
	}

	printf("contextMenuID = %s\n", contextMenuID);

	// Look for the context menu for this ID
	id contextMenu = (id)hashmap_get(&contextMenuMap, (char*)contextMenuID, strlen(contextMenuID));

	printf("CONTEXT MENU = %p\n", contextMenu);

	// Free menu id
	MEMFREE(contextMenuID);

	if( contextMenu == NULL ) {
		printf("\n\n\n\n\n\n\n");
		dumpHashmap("contextMenuMap", &contextMenuMap);
		return;
	}

	// Grab the content view and show the menu
	id contentView = msg(app->mainWindow, s("contentView"));
	printf("contentView = %p\n", contentView);

	// Get the triggering event
	id menuEvent = msg(app->mainWindow, s("currentEvent"));
	printf("menuEvent = %p\n", menuEvent);

	// Show popup
	msg(c("NSMenu"), s("popUpContextMenu:withEvent:forView:"), contextMenu, menuEvent, contentView);

}

void SetColour(struct Application *app, int red, int green, int blue, int alpha) {
	app->red = red;
	app->green = green;
	app->blue = blue;
	app->alpha = alpha;

	applyWindowColour(app);
}

void FullSizeContent(struct Application *app) {
	app->fullSizeContent = 1;
}

void Hide(struct Application *app) {
	ON_MAIN_THREAD(
		msg(app->application, s("hide:"))
	);
}

void Show(struct Application *app) {
	ON_MAIN_THREAD(
		msg(app->mainWindow, s("makeKeyAndOrderFront:"), NULL);
		msg(app->application, s("activateIgnoringOtherApps:"), YES);
	);
}

void SetWindowBackgroundIsTranslucent(struct Application *app) {
	app->windowBackgroundIsTranslucent = 1;
}

// Sends messages to the backend
void messageHandler(id self, SEL cmd, id contentController, id message) {
	struct Application *app = (struct Application *)objc_getAssociatedObject(
							  self, "application");
	const char *name = (const char *)msg(msg(message, s("name")), s("UTF8String"));
	if( strcmp(name, "completed") == 0) {
		// Delete handler
		msg(app->manager, s("removeScriptMessageHandlerForName:"), str("completed"));
		if (app->startHidden == 0) {
			Show(app);
		}
		msg(app->config, s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 0), str("suppressesIncrementalRendering"));
	} else if( strcmp(name, "windowDrag") == 0 ) {
		// Guard against null events
		if( app->mouseEvent != NULL ) {
			HideMouse();
			ON_MAIN_THREAD(
				msg(app->mainWindow, s("performWindowDragWithEvent:"), app->mouseEvent);
			);
		}
	} else if( strcmp(name, "contextMenu") == 0 ) {

		// Did we get a context menu selector?
		if( message == NULL) {
			return;
		}

		const char *contextMenuMessage = cstr(msg(message, s("body")));

		if( contextMenuMessage == NULL ) {
			printf("EMPTY CONTEXT MENU MESSAGE!!\n");
			return;
		}

		// Parse the message
		JsonNode *contextMenuMessageJSON = json_decode(contextMenuMessage);
		if( contextMenuMessageJSON == NULL ) {
			Debug(app, "Error decoding context menu message: %s", contextMenuMessage);
			return;
		}

		// Get menu ID
		JsonNode *contextMenuIDNode = json_find_member(contextMenuMessageJSON, "id");
		if( contextMenuIDNode == NULL ) {
			Debug(app, "Error decoding context menu ID: %s", contextMenuMessage);
			return;
		}
		if( contextMenuIDNode->tag != JSON_STRING ) {
			Debug(app, "Error decoding context menu ID (Not a string): %s", contextMenuMessage);
			return;
		}

		// Get menu Data
		JsonNode *contextMenuDataNode = json_find_member(contextMenuMessageJSON, "data");
		if( contextMenuDataNode == NULL ) {
			Debug(app, "Error decoding context menu data: %s", contextMenuMessage);
			return;
		}
		if( contextMenuDataNode->tag != JSON_STRING ) {
			Debug(app, "Error decoding context menu data (Not a string): %s", contextMenuMessage);
			return;
		}

		// Save a copy of the context menu data
		if ( contextMenuData != NULL ) {
			MEMFREE(contextMenuData);
		}
		contextMenuData = STRCOPY(contextMenuDataNode->string_);

		ON_MAIN_THREAD(
			showContextMenu(app, contextMenuIDNode->string_);
		);

	} else {
		// const char *m = (const char *)msg(msg(message, s("body")), s("UTF8String"));
		const char *m = cstr(msg(message, s("body")));
		app->sendMessageToBackend(m);
	}
}

// Creates a JSON message for the given menuItemID and data
const char* createContextMenuMessage(const char *menuItemID, const char *contextMenuData) {
	JsonNode *jsonObject = json_mkobject();
	json_append_member(jsonObject, "menuItemID", json_mkstring(menuItemID));
	json_append_member(jsonObject, "data", json_mkstring(contextMenuData));
	const char *result = json_encode(jsonObject);
	json_delete(jsonObject);
	return result;
}

// Callback for menu items
void menuItemPressedForApplicationMenu(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));
  // Notify the backend
  const char *message = concat("MC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}

// Callback for tray items
void menuItemPressedForTrayMenu(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));
  // Notify the backend
  const char *message = concat("TC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}


// Callback for context menu items
void menuItemPressedForContextMenus(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));
  // Notify the backend
  const char *contextMenuMessage = createContextMenuMessage(menuItemID, contextMenuData);
  const char *message = concat("XC", contextMenuMessage);
  messageFromWindowCallback(message);
  MEMFREE(message);
  MEMFREE(contextMenuMessage);
}

// Callback for menu items
void checkboxMenuItemPressedForApplicationMenu(id self, SEL cmd, id sender, struct hashmap_s *menuItemMap) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForApplicationMenu, (char*)menuItemID, strlen(menuItemID));

  // Get the current state
  bool state = msg(menuItem, s("state"));

  // Toggle the state
  msg(menuItem, s("setState:"), (state? NSControlStateValueOff : NSControlStateValueOn));

  // Notify the backend
  const char *message = concat("MC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}

// Callback for tray menu items
void checkboxMenuItemPressedForTrayMenu(id self, SEL cmd, id sender, struct hashmap_s *menuItemMap) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForTrayMenu, (char*)menuItemID, strlen(menuItemID));

  // Get the current state
  bool state = msg(menuItem, s("state"));

  // Toggle the state
  msg(menuItem, s("setState:"), (state? NSControlStateValueOff : NSControlStateValueOn));

  // Notify the backend
  const char *message = concat("TC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}

// Callback for context menu items
void checkboxMenuItemPressedForContextMenus(id self, SEL cmd, id sender, struct hashmap_s *menuItemMap) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForContextMenus, (char*)menuItemID, strlen(menuItemID));

  // Get the current state
  bool state = msg(menuItem, s("state"));

  // Toggle the state
  msg(menuItem, s("setState:"), (state? NSControlStateValueOff : NSControlStateValueOn));

  // Notify the backend
  const char *contextMenuMessage = createContextMenuMessage(menuItemID, contextMenuData);
  const char *message = concat("XC", contextMenuMessage);
  messageFromWindowCallback(message);
  MEMFREE(message);
  MEMFREE(contextMenuMessage);
}

// radioMenuItemPressedForApplicationMenu
void radioMenuItemPressedForApplicationMenu(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForApplicationMenu, (char*)menuItemID, strlen(menuItemID));

  // Check the menu items' current state
  bool selected = msg(menuItem, s("state"));

  // If it's already selected, exit early
  if (selected) {
	return;
  }

  // Get this item's radio group members and turn them off
  id *members = (id*)hashmap_get(&radioGroupMapForApplicationMenu, (char*)menuItemID, strlen(menuItemID));

  // Uncheck all members of the group
  id thisMember = members[0];
  int count = 0;
  while(thisMember != NULL) {
	msg(thisMember, s("setState:"), NSControlStateValueOff);
	count = count + 1;
	thisMember = members[count];
  }

  // check the selected menu item
  msg(menuItem, s("setState:"), NSControlStateValueOn);

  // Notify the backend
  const char *message = concat("MC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}


// radioMenuItemPressedForTrayMenu
void radioMenuItemPressedForTrayMenu(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForTrayMenu, (char*)menuItemID, strlen(menuItemID));

  // Check the menu items' current state
  bool selected = msg(menuItem, s("state"));

  // If it's already selected, exit early
  if (selected) {
	return;
  }

  // Get this item's radio group members and turn them off
  id *members = (id*)hashmap_get(&radioGroupMapForTrayMenu, (char*)menuItemID, strlen(menuItemID));

  // Uncheck all members of the group
  id thisMember = members[0];
  int count = 0;
  while(thisMember != NULL) {
	msg(thisMember, s("setState:"), NSControlStateValueOff);
	count = count + 1;
	thisMember = members[count];
  }

  // check the selected menu item
  msg(menuItem, s("setState:"), NSControlStateValueOn);

  // Notify the backend
  const char *message = concat("TC", menuItemID);
  messageFromWindowCallback(message);
  MEMFREE(message);
}

// radioMenuItemPressedForContextMenus
void radioMenuItemPressedForContextMenus(id self, SEL cmd, id sender) {
  const char *menuItemID = (const char *)msg(msg(sender, s("representedObject")), s("pointerValue"));

  // Get the menu item from the menu item map
  id menuItem = (id)hashmap_get(&menuItemMapForContextMenus, (char*)menuItemID, strlen(menuItemID));

  // Check the menu items' current state
  bool selected = msg(menuItem, s("state"));

  // If it's already selected, exit early
  if (selected) {
	return;
  }

  // Get this item's radio group members and turn them off
  id *members = (id*)hashmap_get(&radioGroupMapForContextMenus, (char*)menuItemID, strlen(menuItemID));

  // Uncheck all members of the group
  id thisMember = members[0];
  int count = 0;
  while(thisMember != NULL) {
	msg(thisMember, s("setState:"), NSControlStateValueOff);
	count = count + 1;
	thisMember = members[count];
  }

  // check the selected menu item
  msg(menuItem, s("setState:"), NSControlStateValueOn);

    // Notify the backend
    const char *contextMenuMessage = createContextMenuMessage(menuItemID, contextMenuData);
    const char *message = concat("XC", contextMenuMessage);
    messageFromWindowCallback(message);
    MEMFREE(message);
    MEMFREE(contextMenuMessage);
}

// closeWindow is called when the close button is pressed
void closeWindow(id self, SEL cmd, id sender) {
	struct Application *app = (struct Application *) objc_getAssociatedObject(self, "application");
	app->sendMessageToBackend("WC");
}

void willFinishLaunching(id self, SEL cmd, id sender) {
	struct Application *app = (struct Application *) objc_getAssociatedObject(self, "application");
	printf("\n\n\n\n\n\n\n\n\n\n\n\nI AM HERE!!!!!!!\n\n\n\n\n\n\n\n\n\n\n");
}

bool isDarkMode(struct Application *app) {
	id userDefaults = msg(c("NSUserDefaults"), s("standardUserDefaults"));
	const char *mode = cstr(msg(userDefaults,  s("stringForKey:"), str("AppleInterfaceStyle")));
	return ( mode != NULL && strcmp(mode, "Dark") == 0 );
}

void ExecJS(struct Application *app, const char *js) {
	ON_MAIN_THREAD(
		msg(app->wkwebview,
			s("evaluateJavaScript:completionHandler:"),
			str(js),
			NULL);
	);
}

void themeChanged(id self, SEL cmd, id sender) {
	struct Application *app = (struct Application *)objc_getAssociatedObject(
							  self, "application");
	bool currentThemeIsDark = isDarkMode(app);
	if ( currentThemeIsDark ) {
		ExecJS(app, "window.wails.Events.Emit( 'wails:system:themechange', true );");
	} else {
		ExecJS(app, "window.wails.Events.Emit( 'wails:system:themechange', false );");
	}
}

// void willFinishLaunching(id self) {
//     struct Application *app = (struct Application *) objc_getAssociatedObject(self, "application");
//     Debug(app, "willFinishLaunching called!");
// }

void allocateMenuHashMaps(struct Application *app) {
	// Allocate new menuItem map
	if( 0 != hashmap_create((const unsigned)16, &menuItemMapForApplicationMenu)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate menuItemMapForApplicationMenu!");
	   return;
	}

	// Allocate the Radio Group Cache
	if( 0 != hashmap_create((const unsigned)4, &radioGroupMapForApplicationMenu)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate radioGroupMapForApplicationMenu!");
	   return;
	}
}

void allocateTrayHashMaps(struct Application *app) {
	// Allocate new menuItem map
	if( 0 != hashmap_create((const unsigned)16, &menuItemMapForTrayMenu)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate menuItemMapForTrayMenu!");
	   return;
	}

	// Allocate the Radio Group Cache
	if( 0 != hashmap_create((const unsigned)4, &radioGroupMapForTrayMenu)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate radioGroupMapForTrayMenu!");
	   return;
	}
}

void allocateContextMenuHashMaps(struct Application *app) {

	// Allocate new context menu map
	if( 0 != hashmap_create((const unsigned)4, &contextMenuMap)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate contextMenuMap!");
	}

	// Allocate new menuItem map
	if( 0 != hashmap_create((const unsigned)16, &menuItemMapForContextMenus)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate menuItemMapForContextMenus!");
	}

	// Allocate the Radio Group Cache
	if( 0 != hashmap_create((const unsigned)4, &radioGroupMapForContextMenus)) {
	   // Couldn't allocate map
	   Fatal(app, "Not enough memory to allocate radioGroupMapForContextMenus!");
	   return;
	}
}

void* NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel) {
	// Setup main application struct
	struct Application *result = malloc(sizeof(struct Application));
	result->title = title;
	result->width = width;
	result->height = height;
	result->minWidth = 0;
	result->minHeight = 0;
	result->maxWidth = 0;
	result->maxHeight = 0;
	result->resizable = resizable;
	result->devtools = devtools;
	result->fullscreen = fullscreen;
	result->lock = 0;
	result->maximised = 0;
	result->minimised = 0;
	result->startHidden = startHidden;
	result->decorations = 0;
	result->logLevel = logLevel;

	result->mainWindow = NULL;
	result->mouseEvent = NULL;
	result->mouseDownMonitor = NULL;
	result->mouseUpMonitor = NULL;

	// Features
	result->frame = 1;
	result->hideTitle = 0;
	result->hideTitleBar = 0;
	result->fullSizeContent = 0;
	result->useToolBar = 0;
	result->hideToolbarSeparator = 0;
	result->appearance = NULL;
	result->windowBackgroundIsTranslucent = 0;

	// Window data
	result->delegate = NULL;

	// Menu
	result->menuAsJSON = NULL;
	result->processedMenu = NULL;

	// Tray
	result->trayMenuAsJSON = NULL;
	result->processedTrayMenu = NULL;
	result->statusItem = NULL;

	// Context Menus
	result->contextMenusAsJSON = NULL;
	contextMenuData = NULL;

	// Window Appearance
	result->vibrancyLayer = NULL;
	result->titlebarAppearsTransparent = 0;
	result->webviewIsTranparent = 0;

	result->sendMessageToBackend = (ffenestriCallback) messageFromWindowCallback;

	return (void*) result;
}

int freeHashmapItem(void *const context, struct hashmap_element_s *const e) {
  free(e->data);
  return -1;
}

void destroyMenu(struct Application *app) {

	// Free menu item hashmap
	hashmap_destroy(&menuItemMapForApplicationMenu);

	// Free radio group members
	if( hashmap_num_entries(&radioGroupMapForApplicationMenu) > 0 ) {
		if (0!=hashmap_iterate_pairs(&radioGroupMapForApplicationMenu, freeHashmapItem, NULL)) {
			Fatal(app, "failed to deallocate hashmap entries!");
		}
	}

	//Free radio groups hashmap
	hashmap_destroy(&radioGroupMapForApplicationMenu);

	// Release the menu json if we have it
	if ( app->menuAsJSON != NULL ) {
		MEMFREE(app->menuAsJSON);
	}

	// Release processed menu
	if( app->processedMenu != NULL) {
		json_delete(app->processedMenu);
		app->processedMenu = NULL;
	}
}

void destroyContextMenus(struct Application *app) {

	// If we don't have a context menu, return
	if( app->contextMenusAsJSON == NULL ) {
		return;
	}

	// Free menu item hashmap
	hashmap_destroy(&menuItemMapForContextMenus);

	// Free radio group members
	if( hashmap_num_entries(&radioGroupMapForContextMenus) > 0 ) {
		if (0!=hashmap_iterate_pairs(&radioGroupMapForContextMenus, freeHashmapItem, NULL)) {
			Fatal(app, "failed to deallocate hashmap entries!");
		}
	}

	//Free radio groups hashmap
	hashmap_destroy(&radioGroupMapForContextMenus);

	//Free context menu map
	hashmap_destroy(&contextMenuMap);

    // Destroy processed Context Menus
	if( app->processedContextMenus != NULL) {
		json_delete(app->processedContextMenus);
		app->processedContextMenus = NULL;
	}

	// Release the menu json
    MEMFREE(app->contextMenusAsJSON);

}


void destroyTray(struct Application *app) {

	// If we don't have a tray, exit!
	if( app->trayMenuAsJSON == NULL ) {
		return;
	}

	// Free menu item hashmap
	hashmap_destroy(&menuItemMapForTrayMenu);

	// Free radio group members
	if( hashmap_num_entries(&radioGroupMapForTrayMenu) > 0 ) {
		if (0!=hashmap_iterate_pairs(&radioGroupMapForTrayMenu, freeHashmapItem, NULL)) {
			Fatal(app, "failed to deallocate hashmap entries!");
		}
	}

	//Free radio groups hashmap
	hashmap_destroy(&radioGroupMapForTrayMenu);

	// Release the menu json
	MEMFREE(app->trayMenuAsJSON);


	// Release processed tray
	if( app->processedTrayMenu != NULL) {
		json_delete(app->processedTrayMenu);
		app->processedTrayMenu = NULL;
	}
}

void DestroyApplication(struct Application *app) {
	Debug(app, "Destroying Application");

	// Free the bindings
	if (app->bindings != NULL) {
		MEMFREE(app->bindings);
	} else {
		Debug(app, "Almost a double free for app->bindings");
	}

	// Remove mouse monitors
	if( app->mouseDownMonitor != NULL ) {
		msg( c("NSEvent"), s("removeMonitor:"), app->mouseDownMonitor);
	}
	if( app->mouseUpMonitor != NULL ) {
		msg( c("NSEvent"), s("removeMonitor:"), app->mouseUpMonitor);
	}

	// Destroy the menu
	destroyMenu(app);

	// Destroy the tray
	destroyTray(app);

	// Destroy the context menus
	destroyContextMenus(app);

	// Clear context menu data if we have it
	if( contextMenuData != NULL ) {
		MEMFREE(contextMenuData);
	}

	// Remove script handlers
	msg(app->manager, s("removeScriptMessageHandlerForName:"), str("contextMenu"));
	msg(app->manager, s("removeScriptMessageHandlerForName:"), str("windowDrag"));
	msg(app->manager, s("removeScriptMessageHandlerForName:"), str("external"));

	// Close main window
	msg(app->mainWindow, s("close"));

	// Terminate app
	msg(c("NSApp"), s("terminate:"), NULL);
	Debug(app, "Finished Destroying Application");
}

// Quit will stop the cocoa application and free up all the memory
// used by the application
void Quit(struct Application *app) {
	Debug(app, "Quit Called");
	DestroyApplication(app);
}

// SetTitle sets the main window title to the given string
void SetTitle(struct Application *app, const char *title) {
	Debug(app, "SetTitle Called");
	ON_MAIN_THREAD(
		msg(app->mainWindow, s("setTitle:"), str(title));
	);
}

void ToggleFullscreen(struct Application *app) {
	ON_MAIN_THREAD(
		app->fullscreen = !app->fullscreen;
		MAIN_WINDOW_CALL("toggleFullScreen:");
	);
}

bool isFullScreen(struct Application *app) {
	int mask = (int)msg(app->mainWindow, s("styleMask"));
	bool result = (mask & NSWindowStyleMaskFullscreen) == NSWindowStyleMaskFullscreen;
	return result;
}

// Fullscreen sets the main window to be fullscreen
void Fullscreen(struct Application *app) {
	Debug(app, "Fullscreen Called");
	if( ! isFullScreen(app) ) {
		ToggleFullscreen(app);
	}
}

// UnFullscreen resets the main window after a fullscreen
void UnFullscreen(struct Application *app) {
	Debug(app, "UnFullscreen Called");
	if( isFullScreen(app) ) {
		ToggleFullscreen(app);
	}
}

void Center(struct Application *app) {
	Debug(app, "Center Called");
	ON_MAIN_THREAD(
		MAIN_WINDOW_CALL("center");
	);
}

void ToggleMaximise(struct Application *app) {
	ON_MAIN_THREAD(
		app->maximised = !app->maximised;
		MAIN_WINDOW_CALL("zoom:");
	);
}

void Maximise(struct Application *app) {
	if( app->maximised == 0) {
		ToggleMaximise(app);
	}
}

void Unmaximise(struct Application *app) {
	if( app->maximised == 1) {
		ToggleMaximise(app);
	}
}

void Minimise(struct Application *app) {
	ON_MAIN_THREAD(
		MAIN_WINDOW_CALL("miniaturize:");
	);
 }
void Unminimise(struct Application *app) {
	ON_MAIN_THREAD(
		MAIN_WINDOW_CALL("deminiaturize:");
	);
}

id getCurrentScreen(struct Application *app) {
	id screen = NULL;
	screen = msg(app->mainWindow, s("screen"));
	if( screen == NULL ) {
		screen = msg(c("NSScreen"), u("mainScreen"));
	}
	return screen;
}

void dumpFrame(struct Application *app, const char *message, CGRect frame) {
	Debug(app, message);
	Debug(app, "origin.x %f", frame.origin.x);
	Debug(app, "origin.y %f", frame.origin.y);
	Debug(app, "size.width %f", frame.size.width);
	Debug(app, "size.height %f", frame.size.height);
}

void SetSize(struct Application *app, int width, int height) {
	ON_MAIN_THREAD(
		id screen = getCurrentScreen(app);

		// Get the rect for the window
		CGRect frame = GET_FRAME(app->mainWindow);

		// Credit: https://github.com/patr0nus/DeskGap/blob/73c0ac9f2c73f55b6e81f64f6673a7962b5719cd/lib/src/platform/mac/util/NSScreen%2BGeometry.m
		frame.origin.y = (frame.origin.y + frame.size.height) - (float)height;
		frame.size.width = (float)width;
		frame.size.height = (float)height;

		msg(app->mainWindow, s("setFrame:display:animate:"), frame, 1, 0);
	);
}

void SetPosition(struct Application *app, int x, int y) {
	ON_MAIN_THREAD(
		id screen = getCurrentScreen(app);
		CGRect screenFrame = GET_FRAME(screen);
		CGRect windowFrame = GET_FRAME(app->mainWindow);

		windowFrame.origin.x = screenFrame.origin.x + (float)x;
		windowFrame.origin.y = (screenFrame.origin.y + screenFrame.size.height) - windowFrame.size.height - (float)y;
		msg(app->mainWindow, s("setFrame:display:animate:"), windowFrame, 1, 0);
	);
}

// OpenDialog opens a dialog to select files/directories
void OpenDialog(struct Application *app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolvesAliases, int treatPackagesAsDirectories) {
	Debug(app, "OpenDialog Called with callback id: %s", callbackID);

	// Create an open panel
	ON_MAIN_THREAD(

		// Create the dialog
		id dialog = msg(c("NSOpenPanel"), s("openPanel"));

		// Valid but appears to do nothing.... :/
		msg(dialog, s("setTitle:"), str(title));

		// Filters
		if( filters != NULL && strlen(filters) > 0) {
			id filterString = msg(str(filters), s("stringByReplacingOccurrencesOfString:withString:"), str("*."), str(""));
			filterString = msg(filterString, s("stringByReplacingOccurrencesOfString:withString:"), str(" "), str(""));
			id filterList = msg(filterString, s("componentsSeparatedByString:"), str(","));
			msg(dialog, s("setAllowedFileTypes:"), filterList);
		} else {
			msg(dialog, s("setAllowsOtherFileTypes:"), YES);
		}

		// Default Directory
		if( defaultDir != NULL && strlen(defaultDir) > 0 ) {
			msg(dialog, s("setDirectoryURL:"), url(defaultDir));
		}

		// Default Filename
		if( defaultFilename != NULL && strlen(defaultFilename) > 0 ) {
			msg(dialog, s("setNameFieldStringValue:"), str(defaultFilename));
		}

		// Setup Options
		msg(dialog, s("setCanChooseFiles:"), allowFiles);
		msg(dialog, s("setCanChooseDirectories:"), allowDirs);
		msg(dialog, s("setAllowsMultipleSelection:"), allowMultiple);
		msg(dialog, s("setShowsHiddenFiles:"), showHiddenFiles);
		msg(dialog, s("setCanCreateDirectories:"), canCreateDirectories);
		msg(dialog, s("setResolvesAliases:"), resolvesAliases);
		msg(dialog, s("setTreatsFilePackagesAsDirectories:"), treatPackagesAsDirectories);

		// Setup callback handler
		msg(dialog, s("beginSheetModalForWindow:completionHandler:"), app->mainWindow, ^(id result) {

			// Create the response JSON object
			JsonNode *response = json_mkarray();

			// If the user selected some files
			if( result == (id)1 ) {
				// Grab the URLs returned
				id urls = msg(dialog, s("URLs"));

				// Iterate over all the selected files
				int noOfResults = (int)msg(urls, s("count"));
				for( int index = 0; index < noOfResults; index++ ) {

					// Extract the filename
					id url = msg(urls, s("objectAtIndex:"), index);
					const char *filename = (const char *)msg(msg(url, s("path")), s("UTF8String"));

					// Add the the response array
					json_append_element(response, json_mkstring(filename));
				}
			}

			// Create JSON string and free json memory
			char *encoded = json_stringify(response, "");
			json_delete(response);

			// Construct callback message. Format "D<callbackID>|<json array of strings>"
			const char *callback = concat("DO", callbackID);
			const char *header = concat(callback, "|");
			const char *responseMessage = concat(header, encoded);

			// Send message to backend
			app->sendMessageToBackend(responseMessage);

			// Free memory
			MEMFREE(header);
			MEMFREE(callback);
			MEMFREE(responseMessage);
		});

		msg( c("NSApp"), s("runModalForWindow:"), app->mainWindow);
	);
}

// SaveDialog opens a dialog to select files/directories
void SaveDialog(struct Application *app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories) {
	Debug(app, "SaveDialog Called with callback id: %s", callbackID);

	// Create an open panel
	ON_MAIN_THREAD(

		// Create the dialog
		id dialog = msg(c("NSSavePanel"), s("savePanel"));

		// Valid but appears to do nothing.... :/
		msg(dialog, s("setTitle:"), str(title));

		// Filters
		if( filters != NULL && strlen(filters) > 0) {
			id filterString = msg(str(filters), s("stringByReplacingOccurrencesOfString:withString:"), str("*."), str(""));
			filterString = msg(filterString, s("stringByReplacingOccurrencesOfString:withString:"), str(" "), str(""));
			id filterList = msg(filterString, s("componentsSeparatedByString:"), str(","));
			msg(dialog, s("setAllowedFileTypes:"), filterList);
		} else {
			msg(dialog, s("setAllowsOtherFileTypes:"), YES);
		}

		// Default Directory
		if( defaultDir != NULL && strlen(defaultDir) > 0 ) {
			msg(dialog, s("setDirectoryURL:"), url(defaultDir));
		}

		// Default Filename
		if( defaultFilename != NULL && strlen(defaultFilename) > 0 ) {
			msg(dialog, s("setNameFieldStringValue:"), str(defaultFilename));
		}

		// Setup Options
		msg(dialog, s("setShowsHiddenFiles:"), showHiddenFiles);
		msg(dialog, s("setCanCreateDirectories:"), canCreateDirectories);
		msg(dialog, s("setTreatsFilePackagesAsDirectories:"), treatPackagesAsDirectories);

		// Setup callback handler
		msg(dialog, s("beginSheetModalForWindow:completionHandler:"), app->mainWindow, ^(id result) {

			// Default is blank
			const char *filename = "";

			// If the user selected some files
			if( result == (id)1 ) {
				// Grab the URL returned
				id url = msg(dialog, s("URL"));
				filename = (const char *)msg(msg(url, s("path")), s("UTF8String"));
			}

			// Construct callback message. Format "DS<callbackID>|<json array of strings>"
			const char *callback = concat("DS", callbackID);
			const char *header = concat(callback, "|");
			const char *responseMessage = concat(header, filename);

			// Send message to backend
			app->sendMessageToBackend(responseMessage);

			// Free memory
			MEMFREE(header);
			MEMFREE(callback);
			MEMFREE(responseMessage);
		});

		msg( c("NSApp"), s("runModalForWindow:"), app->mainWindow);
	);
}

const char *invoke = "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}};";

// DisableFrame disables the window frame
void DisableFrame(struct Application *app)
{
   app->frame = 0;
}

void setMinMaxSize(struct Application *app)
{
	if (app->maxHeight > 0 && app->maxWidth > 0)
	{
		msg(app->mainWindow, s("setMaxSize:"), CGSizeMake(app->maxWidth, app->maxHeight));
	}
	if (app->minHeight > 0 && app->minWidth > 0)
	{
		msg(app->mainWindow, s("setMinSize:"), CGSizeMake(app->minWidth, app->minHeight));
	}
}

void SetMinWindowSize(struct Application *app, int minWidth, int minHeight)
{
	app->minWidth = minWidth;
	app->minHeight = minHeight;

	// Apply if the window is created
	if( app->mainWindow != NULL ) {
		ON_MAIN_THREAD(
			setMinMaxSize(app);
		);
	}
}

void SetMaxWindowSize(struct Application *app, int maxWidth, int maxHeight)
{
	app->maxWidth = maxWidth;
	app->maxHeight = maxHeight;

	// Apply if the window is created
	if( app->mainWindow != NULL ) {
		ON_MAIN_THREAD(
			setMinMaxSize(app);
		);
	}
}


void SetDebug(void *applicationPointer, int flag) {
	debug = flag;
}

// SetMenu sets the initial menu for the application
void SetMenu(struct Application *app, const char *menuAsJSON) {
	app->menuAsJSON = menuAsJSON;
}

// SetTray sets the initial tray menu for the application
void SetTray(struct Application *app, const char *trayMenuAsJSON) {
	app->trayMenuAsJSON = trayMenuAsJSON;
}

// SetContextMenus sets the context menu map for this application
void SetContextMenus(struct Application *app, const char *contextMenusAsJSON) {
	app->contextMenusAsJSON = contextMenusAsJSON;
}

void SetBindings(struct Application *app, const char *bindings) {
	const char* temp = concat("window.wailsbindings = \"", bindings);
	const char* jscall = concat(temp, "\";");
	MEMFREE(temp);
	app->bindings = jscall;
}

void makeWindowBackgroundTranslucent(struct Application *app) {
	id contentView = msg(app->mainWindow, s("contentView"));
	id effectView = msg(c("NSVisualEffectView"), s("alloc"));
	CGRect bounds = GET_BOUNDS(contentView);
	effectView = msg(effectView, s("initWithFrame:"), bounds);

	msg(effectView, s("setAutoresizingMask:"), NSViewWidthSizable | NSViewHeightSizable);
	msg(effectView, s("setBlendingMode:"), NSVisualEffectBlendingModeBehindWindow);
	msg(effectView, s("setState:"), NSVisualEffectStateActive);
	msg(contentView, s("addSubview:positioned:relativeTo:"), effectView, NSWindowBelow, NULL);

	app->vibrancyLayer = effectView;
	Debug(app, "effectView: %p", effectView);
}

void enableBoolConfig(id config, const char *setting) {
	msg(msg(config, s("preferences")), s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 1), str(setting));
}

void disableBoolConfig(id config, const char *setting) {
	msg(msg(config, s("preferences")), s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 0), str(setting));
}

void processDecorations(struct Application *app) {

	int decorations = 0;

	if (app->frame == 1 ) {
		if( app->hideTitleBar == 0) {
			decorations |= NSWindowStyleMaskTitled;
		}
		decorations |= NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable;
	}

	if (app->resizable) {
		decorations |= NSWindowStyleMaskResizable;
	}

	if (app->fullscreen) {
		decorations |= NSWindowStyleMaskFullscreen;
	}

	if( app->fullSizeContent || app->frame == 0) {
		decorations |= NSWindowStyleMaskFullSizeContentView;
	}

	app->decorations = decorations;
}

void createApplication(struct Application *app) {
	id application = msg(c("NSApplication"), s("sharedApplication"));
	app->application = application;
	msg(application, s("setActivationPolicy:"), 0);
}

void DarkModeEnabled(struct Application *app, const char *callbackID) {
	ON_MAIN_THREAD(
		const char *result = isDarkMode(app) ? "T" : "F";

		// Construct callback message. Format "SD<callbackID>|<json array of strings>"
		const char *callback = concat("SD", callbackID);
		const char *header = concat(callback, "|");
		const char *responseMessage = concat(header, result);
		// Send message to backend
		app->sendMessageToBackend(responseMessage);

		// Free memory
		MEMFREE(header);
		MEMFREE(callback);
		MEMFREE(responseMessage);
	);
}

void createDelegate(struct Application *app) {
		// Define delegate
	Class delegateClass = objc_allocateClassPair((Class) c("NSObject"), "AppDelegate", 0);
	  bool resultAddProtoc = class_addProtocol(delegateClass, objc_getProtocol("NSApplicationDelegate"));
	class_addMethod(delegateClass, s("applicationShouldTerminateAfterLastWindowClosed:"), (IMP) yes, "c@:@");
	class_addMethod(delegateClass, s("applicationWillTerminate:"), (IMP) closeWindow, "v@:@");
	class_addMethod(delegateClass, s("applicationWillFinishLaunching:"), (IMP) willFinishLaunching, "v@:@");

	// Menu Callbacks
	class_addMethod(delegateClass, s("menuCallbackForApplicationMenu:"), (IMP)menuItemPressedForApplicationMenu, "v@:@");
	class_addMethod(delegateClass, s("checkboxMenuCallbackForApplicationMenu:"), (IMP) checkboxMenuItemPressedForApplicationMenu, "v@:@");
	class_addMethod(delegateClass, s("radioMenuCallbackForApplicationMenu:"), (IMP) radioMenuItemPressedForApplicationMenu, "v@:@");
	class_addMethod(delegateClass, s("menuCallbackForTrayMenu:"), (IMP)menuItemPressedForTrayMenu, "v@:@");
	class_addMethod(delegateClass, s("checkboxMenuCallbackForTrayMenu:"), (IMP) checkboxMenuItemPressedForTrayMenu, "v@:@");
	class_addMethod(delegateClass, s("radioMenuCallbackForTrayMenu:"), (IMP) radioMenuItemPressedForTrayMenu, "v@:@");
	class_addMethod(delegateClass, s("menuCallbackForContextMenus:"), (IMP)menuItemPressedForContextMenus, "v@:@");
	class_addMethod(delegateClass, s("checkboxMenuCallbackForContextMenus:"), (IMP) checkboxMenuItemPressedForContextMenus, "v@:@");
	class_addMethod(delegateClass, s("radioMenuCallbackForContextMenus:"), (IMP) radioMenuItemPressedForContextMenus, "v@:@");

	// Script handler
	class_addMethod(delegateClass, s("userContentController:didReceiveScriptMessage:"), (IMP) messageHandler, "v@:@@");
	objc_registerClassPair(delegateClass);

	// Create delegate
	id delegate = msg((id)delegateClass, s("new"));
	objc_setAssociatedObject(delegate, "application", (id)app, OBJC_ASSOCIATION_ASSIGN);

	// Theme change listener
	class_addMethod(delegateClass, s("themeChanged:"), (IMP) themeChanged, "v@:@@");

	// Get defaultCenter
	id defaultCenter = msg(c("NSDistributedNotificationCenter"), s("defaultCenter"));
	msg(defaultCenter, s("addObserver:selector:name:object:"), delegate, s("themeChanged:"), str("AppleInterfaceThemeChangedNotification"), NULL);

	app->delegate = delegate;

	msg(app->application, s("setDelegate:"), delegate);
}

void createMainWindow(struct Application *app) {
	// Create main window
	id mainWindow = ALLOC("NSWindow");
	mainWindow = msg(mainWindow, s("initWithContentRect:styleMask:backing:defer:"),
		  CGRectMake(0, 0, app->width, app->height), app->decorations, NSBackingStoreBuffered, NO);
	msg(mainWindow, s("autorelease"));

	// Set Appearance
	if( app->appearance != NULL ) {
		msg(mainWindow, s("setAppearance:"),
			msg(c("NSAppearance"), s("appearanceNamed:"), str(app->appearance))
		);
	}

	// Set Title appearance
	msg(mainWindow, s("setTitlebarAppearsTransparent:"), app->titlebarAppearsTransparent ? YES : NO);
	msg(mainWindow, s("setTitleVisibility:"), app->hideTitle);

	app->mainWindow = mainWindow;
}

const char* getInitialState(struct Application *app) {
	const char *result = "";
	if( isDarkMode(app) ) {
		result = "window.wails.System.IsDarkMode.set(true);";
	} else {
		result = "window.wails.System.IsDarkMode.set(false);";
	}
	char buffer[999];
	snprintf(&buffer[0], sizeof(buffer), "window.wails.System.LogLevel.set(%d);", app->logLevel);
	result = concat(result, &buffer[0]);
	Debug(app, "initialstate = %s", result);
	return result;
}

id createMenuItem(id title, const char *action, const char *key) {
  id item = ALLOC("NSMenuItem");
  msg(item, s("initWithTitle:action:keyEquivalent:"), title, s(action), str(key));
  msg(item, s("autorelease"));
  return item;
}

id createMenuItemNoAutorelease( id title, const char *action, const char *key) {
  id item = ALLOC("NSMenuItem");
  msg(item, s("initWithTitle:action:keyEquivalent:"), title, s(action), str(key));
  return item;
}

id createMenu(id title) {
  id menu = ALLOC("NSMenu");
  msg(menu, s("initWithTitle:"), title);
  msg(menu, s("setAutoenablesItems:"), NO);
  msg(menu, s("autorelease"));
  return menu;
}

id addMenuItem(id menu, const char *title, const char *action, const char *key, bool disabled) {
	id item = createMenuItem(str(title), action, key);
	msg(item, s("setEnabled:"), !disabled);
	msg(menu, s("addItem:"), item);
	return item;
}

void addSeparator(id menu) {
  id item = msg(c("NSMenuItem"), s("separatorItem"));
  msg(menu, s("addItem:"), item);
}

void createDefaultAppMenu(id parentMenu) {
// App Menu
  id appName = msg(msg(c("NSProcessInfo"), s("processInfo")), s("processName"));
  id appMenuItem = createMenuItemNoAutorelease(appName, NULL, "");
  id appMenu = createMenu(appName);

  msg(appMenuItem, s("setSubmenu:"), appMenu);
  msg(parentMenu, s("addItem:"), appMenuItem);

  id title = msg(str("Hide "), s("stringByAppendingString:"), appName);
  id item = createMenuItem(title, "hide:", "h");
  msg(appMenu, s("addItem:"), item);

  id hideOthers = addMenuItem(appMenu, "Hide Others", "hideOtherApplications:", "h", FALSE);
  msg(hideOthers, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagCommand));

  addMenuItem(appMenu, "Show All", "unhideAllApplications:", "", FALSE);

  addSeparator(appMenu);

  title = msg(str("Quit "), s("stringByAppendingString:"), appName);
  item = createMenuItem(title, "terminate:", "q");
  msg(appMenu, s("addItem:"), item);
}


void createDefaultEditMenu(id parentMenu) {
  // Edit Menu
  id editMenuItem = createMenuItemNoAutorelease(str("Edit"), NULL, "");
  id editMenu = createMenu(str("Edit"));

  msg(editMenuItem, s("setSubmenu:"), editMenu);
  msg(parentMenu, s("addItem:"), editMenuItem);

  addMenuItem(editMenu, "Undo", "undo:", "z", FALSE);
  addMenuItem(editMenu, "Redo", "redo:", "y", FALSE);
  addSeparator(editMenu);
  addMenuItem(editMenu, "Cut", "cut:", "x", FALSE);
  addMenuItem(editMenu, "Copy", "copy:", "c", FALSE);
  addMenuItem(editMenu, "Paste", "paste:", "v", FALSE);
  addMenuItem(editMenu, "Select All", "selectAll:", "a", FALSE);
}

void parseMenuRole(struct Application *app, id parentMenu, JsonNode *item) {
  const char *roleName = item->string_;

  if ( STREQ(roleName, "appMenu") ) {
	createDefaultAppMenu(parentMenu);
	return;
  }
  if ( STREQ(roleName, "editMenu")) {
	createDefaultEditMenu(parentMenu);
	return;
  }
  if ( STREQ(roleName, "hide")) {
	addMenuItem(parentMenu, "Hide Window", "hide:", "h", FALSE);
	return;
  }
  if ( STREQ(roleName, "hideothers")) {
	id hideOthers = addMenuItem(parentMenu, "Hide Others", "hideOtherApplications:", "h", FALSE);
	msg(hideOthers, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagCommand));
	return;
  }
  if ( STREQ(roleName, "unhide")) {
	addMenuItem(parentMenu, "Show All", "unhideAllApplications:", "", FALSE);
	return;
  }
  if ( STREQ(roleName, "front")) {
	addMenuItem(parentMenu, "Bring All to Front", "arrangeInFront:", "", FALSE);
	return;
  }
  if ( STREQ(roleName, "undo")) {
	addMenuItem(parentMenu, "Undo", "undo:", "z", FALSE);
	return;
  }
  if ( STREQ(roleName, "redo")) {
	addMenuItem(parentMenu, "Redo", "redo:", "y", FALSE);
	return;
  }
  if ( STREQ(roleName, "cut")) {
	addMenuItem(parentMenu, "Cut", "cut:", "x", FALSE);
	return;
  }
  if ( STREQ(roleName, "copy")) {
	addMenuItem(parentMenu, "Copy", "copy:", "c", FALSE);
	return;
  }
  if ( STREQ(roleName, "paste")) {
	addMenuItem(parentMenu, "Paste", "paste:", "v", FALSE);
	return;
  }
  if ( STREQ(roleName, "delete")) {
	addMenuItem(parentMenu, "Delete", "delete:", "", FALSE);
	return;
  }
  if( STREQ(roleName, "pasteandmatchstyle")) {
	id pasteandmatchstyle = addMenuItem(parentMenu, "Paste and Match Style", "pasteandmatchstyle:", "v", FALSE);
	msg(pasteandmatchstyle, s("setKeyEquivalentModifierMask:"), (NSEventModifierFlagOption | NSEventModifierFlagShift | NSEventModifierFlagCommand));
  }
  if ( STREQ(roleName, "selectall")) {
	addMenuItem(parentMenu, "Select All", "selectAll:", "a", FALSE);
	return;
  }
  if ( STREQ(roleName, "minimize")) {
	addMenuItem(parentMenu, "Minimize", "miniaturize:", "m", FALSE);
	return;
  }
  if ( STREQ(roleName, "zoom")) {
	addMenuItem(parentMenu, "Zoom", "performZoom:", "", FALSE);
	return;
  }
  if ( STREQ(roleName, "quit")) {
	addMenuItem(parentMenu, "Quit (More work TBD)", "terminate:", "q", FALSE);
	return;
  }
  if ( STREQ(roleName, "togglefullscreen")) {
	addMenuItem(parentMenu, "Toggle Full Screen", "toggleFullScreen:", "f", FALSE);
	return;
  }

}

const char* getJSONString(JsonNode *item, const char* key) {
  // Get key
  JsonNode *node = json_find_member(item, key);
  const char *result = "";
  if ( node != NULL && node->tag == JSON_STRING) {
	result = node->string_;
  }
  return result;
}

bool getJSONBool(JsonNode *item, const char* key, bool *result) {
  JsonNode *node = json_find_member(item, key);
  if ( node != NULL && node->tag == JSON_BOOL) {
	*result = node->bool_;
	return true;
  }
  return false;
}

bool getJSONInt(JsonNode *item, const char* key, int *result) {
  JsonNode *node = json_find_member(item, key);
  if ( node != NULL && node->tag == JSON_NUMBER) {
	*result = (int) node->number_;
	return true;
  }
  return false;
}

// This converts a string array of modifiers into the
// equivalent MacOS Modifier Flags
unsigned long parseModifiers(const char **modifiers) {

  // Our result is a modifier flag list
  unsigned long result = 0;

  const char *thisModifier = modifiers[0];
  int count = 0;
  while( thisModifier != NULL ) {
	// Determine flags
	if( STREQ(thisModifier, "CmdOrCtrl") ) {
	  result |= NSEventModifierFlagCommand;
	}
	if( STREQ(thisModifier, "OptionOrAlt") ) {
	  result |= NSEventModifierFlagOption;
	}
	if( STREQ(thisModifier, "Shift") ) {
	  result |= NSEventModifierFlagShift;
	}
	if( STREQ(thisModifier, "Super") ) {
	  result |= NSEventModifierFlagCommand;
	}
	if( STREQ(thisModifier, "Control") ) {
	  result |= NSEventModifierFlagControl;
	}
	count++;
	thisModifier = modifiers[count];
  }
  return result;
}

id processAcceleratorKey(const char *key) {

	// Guard against no accelerator key
	if( key == NULL ) {
		return str("");
	}

  if( STREQ(key, "Backspace") ) {
	return strunicode(0x0008);
  }
  if( STREQ(key, "Tab") ) {
	return strunicode(0x0009);
  }
  if( STREQ(key, "Return") ) {
	return strunicode(0x000d);
  }
  if( STREQ(key, "Escape") ) {
	return strunicode(0x001b);
  }
  if( STREQ(key, "Left") ) {
	return strunicode(0x001c);
  }
  if( STREQ(key, "Right") ) {
	return strunicode(0x001d);
  }
  if( STREQ(key, "Up") ) {
	return strunicode(0x001e);
  }
  if( STREQ(key, "Down") ) {
	return strunicode(0x001f);
  }
  if( STREQ(key, "Space") ) {
	return strunicode(0x0020);
  }
  if( STREQ(key, "Delete") ) {
	return strunicode(0x007f);
  }
  if( STREQ(key, "Home") ) {
	return strunicode(0x2196);
  }
  if( STREQ(key, "End") ) {
	return strunicode(0x2198);
  }
  if( STREQ(key, "Page Up") ) {
	return strunicode(0x21de);
  }
  if( STREQ(key, "Page Down") ) {
	return strunicode(0x21df);
  }
  if( STREQ(key, "F1") ) {
	return strunicode(0xf704);
  }
  if( STREQ(key, "F2") ) {
	return strunicode(0xf705);
  }
  if( STREQ(key, "F3") ) {
	return strunicode(0xf706);
  }
  if( STREQ(key, "F4") ) {
	return strunicode(0xf707);
  }
  if( STREQ(key, "F5") ) {
	return strunicode(0xf708);
  }
  if( STREQ(key, "F6") ) {
	return strunicode(0xf709);
  }
  if( STREQ(key, "F7") ) {
	return strunicode(0xf70a);
  }
  if( STREQ(key, "F8") ) {
	return strunicode(0xf70b);
  }
  if( STREQ(key, "F9") ) {
	return strunicode(0xf70c);
  }
  if( STREQ(key, "F10") ) {
	return strunicode(0xf70d);
  }
  if( STREQ(key, "F11") ) {
	return strunicode(0xf70e);
  }
  if( STREQ(key, "F12") ) {
	return strunicode(0xf70f);
  }
  if( STREQ(key, "F13") ) {
	return strunicode(0xf710);
  }
  if( STREQ(key, "F14") ) {
	return strunicode(0xf711);
  }
  if( STREQ(key, "F15") ) {
	return strunicode(0xf712);
  }
  if( STREQ(key, "F16") ) {
	return strunicode(0xf713);
  }
  if( STREQ(key, "F17") ) {
	return strunicode(0xf714);
  }
  if( STREQ(key, "F18") ) {
	return strunicode(0xf715);
  }
  if( STREQ(key, "F19") ) {
	return strunicode(0xf716);
  }
  if( STREQ(key, "F20") ) {
	return strunicode(0xf717);
  }
  if( STREQ(key, "F21") ) {
	return strunicode(0xf718);
  }
  if( STREQ(key, "F22") ) {
	return strunicode(0xf719);
  }
  if( STREQ(key, "F23") ) {
	return strunicode(0xf71a);
  }
  if( STREQ(key, "F24") ) {
	return strunicode(0xf71b);
  }
  if( STREQ(key, "F25") ) {
	return strunicode(0xf71c);
  }
  if( STREQ(key, "F26") ) {
	return strunicode(0xf71d);
  }
  if( STREQ(key, "F27") ) {
	return strunicode(0xf71e);
  }
  if( STREQ(key, "F28") ) {
	return strunicode(0xf71f);
  }
  if( STREQ(key, "F29") ) {
	return strunicode(0xf720);
  }
  if( STREQ(key, "F30") ) {
	return strunicode(0xf721);
  }
  if( STREQ(key, "F31") ) {
	return strunicode(0xf722);
  }
  if( STREQ(key, "F32") ) {
	return strunicode(0xf723);
  }
  if( STREQ(key, "F33") ) {
	return strunicode(0xf724);
  }
  if( STREQ(key, "F34") ) {
	return strunicode(0xf725);
  }
  if( STREQ(key, "F35") ) {
	return strunicode(0xf726);
  }
//  if( STREQ(key, "Insert") ) {
//	return strunicode(0xf727);
//  }
//  if( STREQ(key, "PrintScreen") ) {
//	return strunicode(0xf72e);
//  }
//  if( STREQ(key, "ScrollLock") ) {
//	return strunicode(0xf72f);
//  }
  if( STREQ(key, "NumLock") ) {
	return strunicode(0xf739);
  }

  return str(key);
}


id parseTextMenuItem(struct Application *app, id parentMenu, const char *title, const char *menuid, bool disabled, const char *acceleratorkey, const char **modifiers, const char *menuCallback) {
	id item = ALLOC("NSMenuItem");
	id wrappedId = msg(c("NSValue"), s("valueWithPointer:"), menuid);
	msg(item, s("setRepresentedObject:"), wrappedId);

	id key = processAcceleratorKey(acceleratorkey);
	msg(item, s("initWithTitle:action:keyEquivalent:"), str(title),
			  s(menuCallback), key);

	msg(item, s("setEnabled:"), !disabled);
	msg(item, s("autorelease"));

	// Process modifiers
	if( modifiers != NULL ) {
		unsigned long modifierFlags = parseModifiers(modifiers);
		msg(item, s("setKeyEquivalentModifierMask:"), modifierFlags);
	}
	msg(parentMenu, s("addItem:"), item);

	return item;
}

id parseCheckboxMenuItem(struct Application *app, id parentmenu, const char
*title, const char *menuid, bool disabled, bool checked, const char *key,
struct hashmap_s *menuItemMap, const char *checkboxCallbackFunction) {
	id item = ALLOC("NSMenuItem");

	// Store the item in the menu item map
	hashmap_put(menuItemMap, (char*)menuid, strlen(menuid), item);

	id wrappedId = msg(c("NSValue"), s("valueWithPointer:"), menuid);
	msg(item, s("setRepresentedObject:"), wrappedId);
	msg(item, s("initWithTitle:action:keyEquivalent:"), str(title), s(checkboxCallbackFunction), str(key));
	msg(item, s("setEnabled:"), !disabled);
	msg(item, s("autorelease"));
	msg(item, s("setState:"), (checked ? NSControlStateValueOn : NSControlStateValueOff));
	msg(parentmenu, s("addItem:"), item);
	return item;
}

id parseRadioMenuItem(struct Application *app, id parentmenu, const char *title,
 const char *menuid, bool disabled, bool checked, const char *acceleratorkey,
 struct hashmap_s *menuItemMap, const char *radioCallbackFunction) {
	id item = ALLOC("NSMenuItem");

	// Store the item in the menu item map
	hashmap_put(menuItemMap, (char*)menuid, strlen(menuid), item);

	id wrappedId = msg(c("NSValue"), s("valueWithPointer:"), menuid);
	msg(item, s("setRepresentedObject:"), wrappedId);

	id key = processAcceleratorKey(acceleratorkey);

	msg(item, s("initWithTitle:action:keyEquivalent:"), str(title), s(radioCallbackFunction), key);

	msg(item, s("setEnabled:"), !disabled);
	msg(item, s("autorelease"));
	msg(item, s("setState:"), (checked ? NSControlStateValueOn : NSControlStateValueOff));

	msg(parentmenu, s("addItem:"), item);
	return item;

}

void parseMenuItem(struct Application *app, id parentMenu, JsonNode *item,
struct hashmap_s *menuItemMap, const char *checkboxCallbackFunction, const char
*radioCallbackFunction, const char *menuCallbackFunction) {

  // Check if this item is hidden and if so, exit early!
  bool hidden = false;
  getJSONBool(item, "Hidden", &hidden);
  if( hidden ) {
	return;
  }

  // Get the role
  JsonNode *role = json_find_member(item, "Role");
  if( role != NULL ) {
	parseMenuRole(app, parentMenu, role);
	return;
  }

  // Check if this is a submenu
  JsonNode *submenu = json_find_member(item, "SubMenu");
  if( submenu != NULL ) {
	// Get the label
	JsonNode *menuNameNode = json_find_member(item, "Label");
	const char *name = "";
	if ( menuNameNode != NULL) {
	  name = menuNameNode->string_;
	}

	id thisMenuItem = createMenuItemNoAutorelease(str(name), NULL, "");
	id thisMenu = createMenu(str(name));

	msg(thisMenuItem, s("setSubmenu:"), thisMenu);
	msg(parentMenu, s("addItem:"), thisMenuItem);

	// Loop over submenu items
	JsonNode *item;
	json_foreach(item, submenu) {
	  // Get item label
	  parseMenuItem(app, thisMenu, item, menuItemMap, checkboxCallbackFunction, radioCallbackFunction, menuCallbackFunction);
	}

	return;
  }

  // This is a user menu. Get the common data
  // Get the label
  const char *label = getJSONString(item, "Label");
  if ( label == NULL) {
	label = "(empty)";
  }

  const char *menuid = getJSONString(item, "ID");
  if ( menuid == NULL) {
	menuid = "";
  }

  bool disabled = false;
  getJSONBool(item, "Disabled", &disabled);

  // Get the Accelerator
  JsonNode *accelerator = json_find_member(item, "Accelerator");
  const char *acceleratorkey = NULL;
  const char **modifiers = NULL;

  // If we have an accelerator
  if( accelerator != NULL ) {
	// Get the key
	  acceleratorkey = getJSONString(accelerator, "Key");
	  // Check if there are modifiers
	  JsonNode *modifiersList = json_find_member(accelerator, "Modifiers");
	  if ( modifiersList != NULL ) {
		// Allocate an array of strings
		int noOfModifiers = json_array_length(modifiersList);
		modifiers = malloc(sizeof(const char *) * (noOfModifiers + 1));
		JsonNode *modifier;
		int count = 0;
		// Iterate the modifiers and save a reference to them in our new array
		json_foreach(modifier, modifiersList) {
		  // Get modifier name
		  modifiers[count] = modifier->string_;
		  count++;
		}
		// Null terminate the modifier list
		modifiers[count] = NULL;
	  }
  }


  // Get the Type
  JsonNode *type = json_find_member(item, "Type");
  if( type != NULL ) {

	if( STREQ(type->string_, "Text")) {
	  parseTextMenuItem(app, parentMenu, label, menuid, disabled, acceleratorkey, modifiers, menuCallbackFunction);
	}
	else if ( STREQ(type->string_, "Separator")) {
	  addSeparator(parentMenu);
	}
	else if ( STREQ(type->string_, "Checkbox")) {
	  // Get checked state
	  bool checked = false;
	  getJSONBool(item, "Checked", &checked);

	  parseCheckboxMenuItem(app, parentMenu, label, menuid, disabled, checked, "", menuItemMap, checkboxCallbackFunction);
	}
	else if ( STREQ(type->string_, "Radio")) {
	  // Get checked state
	  bool checked = false;
	  getJSONBool(item, "Checked", &checked);

	  parseRadioMenuItem(app, parentMenu, label, menuid, disabled, checked, "", menuItemMap, radioCallbackFunction);
	}

	if ( modifiers != NULL ) {
	  free(modifiers);
	}

	return;
  }
}

void parseMenu(struct Application *app, id parentMenu, JsonNode *menu, struct hashmap_s *menuItemMap, const char *checkboxCallbackFunction, const char *radioCallbackFunction, const char *menuCallbackFunction) {
  JsonNode *items = json_find_member(menu, "Items");
  if( items == NULL ) {
	// Parse error!
	Fatal(app, "Unable to find Items:", app->menuAsJSON);
	return;
  }

  // Iterate items
  JsonNode *item;
  json_foreach(item, items) {
	// Get item label
	parseMenuItem(app, parentMenu, item, menuItemMap, checkboxCallbackFunction, radioCallbackFunction, menuCallbackFunction);
  }
}

void dumpMemberList(const char *name, id *memberList) {
  void *member = memberList[0];
  int count = 0;
  printf("%s = %p -> [ ", name, memberList);
  while( member != NULL ) {
	printf("%p ", member);
	count = count + 1;
	member = memberList[count];
  }
  printf("]\n");
}

void processRadioGroup(JsonNode *radioGroup, struct hashmap_s *menuItemMap,
struct hashmap_s *radioGroupMap) {

  int groupLength;
  getJSONInt(radioGroup, "Length", &groupLength);
  JsonNode *members = json_find_member(radioGroup, "Members");
  JsonNode *member;

  // Allocate array
  size_t arrayLength = sizeof(id)*(groupLength+1);
  id memberList[arrayLength];

  // Build the radio group items
  int count=0;
  json_foreach(member, members) {
	// Get menu by id
	id menuItem = (id)hashmap_get(menuItemMap, (char*)member->string_, strlen(member->string_));
	// Save Member
	memberList[count] = menuItem;
	count = count + 1;
  }
  // Null terminate array
  memberList[groupLength] = 0;

  // dumpMemberList("memberList", memberList);

  // Store the members
  json_foreach(member, members) {
	// Copy the memberList
	char *newMemberList = (char *)malloc(arrayLength);
	memcpy(newMemberList, memberList, arrayLength);
	// dumpMemberList("newMemberList", newMemberList);
	// printf("Address of newMemberList = %p\n", newMemberList);

	// add group to each member of group
	hashmap_put(radioGroupMap, member->string_, strlen(member->string_), newMemberList);
  }

  // dumpHashmap("radioGroupMap", &radioGroupMap);

}

void parseMenuData(struct Application *app) {

	// Allocate the hashmaps we need
	allocateMenuHashMaps(app);

	// Create a new menu bar
	id menubar = createMenu(str(""));

	// Parse the processed menu json
	app->processedMenu = json_decode(app->menuAsJSON);

	if( app->processedMenu == NULL ) {
		// Parse error!
		Fatal(app, "Unable to parse Menu JSON: %s", app->menuAsJSON);
		return;
	}


	// Pull out the Menu
	JsonNode *menuData = json_find_member(app->processedMenu, "Menu");
	if( menuData == NULL ) {
		// Parse error!
		Fatal(app, "Unable to find Menu data: %s", app->processedMenu);
		return;
	}


	parseMenu(app, menubar, menuData, &menuItemMapForApplicationMenu,
	"checkboxMenuCallbackForApplicationMenu:", "radioMenuCallbackForApplicationMenu:", "menuCallbackForApplicationMenu:");

	// Create the radiogroup cache
	JsonNode *radioGroups = json_find_member(app->processedMenu, "RadioGroups");
	if( radioGroups == NULL ) {
		// Parse error!
		Fatal(app, "Unable to find RadioGroups data: %s", app->processedMenu);
		return;
	}

	// Iterate radio groups
	JsonNode *radioGroup;
	json_foreach(radioGroup, radioGroups) {
		// Get item label
		processRadioGroup(radioGroup, &menuItemMapForApplicationMenu, &radioGroupMapForApplicationMenu);
	}

	// Apply the menu bar
	msg(msg(c("NSApplication"), s("sharedApplication")), s("setMainMenu:"), menubar);

}


// UpdateMenu replaces the current menu with the given one
void UpdateMenu(struct Application *app, const char *menuAsJSON) {
	Debug(app, "Menu is now: %s", menuAsJSON);
	ON_MAIN_THREAD (

		// Remove the current Menu
		id menubar = msg(msg(c("NSApplication"), s("sharedApplication")), s("mainMenu"));
		Debug(app, "Got menubar: %p", menubar);
		msg(menubar, s("removeAllItems"));

		// Free up memory
		destroyMenu(app);

		// Set the menu JSON
		app->menuAsJSON = menuAsJSON;
		parseMenuData(app);
	);
}

void dumpContextMenus(struct Application *app) {
	dumpHashmap("menuItemMapForContextMenus", &menuItemMapForContextMenus);
	printf("&menuItemMapForContextMenus = %p\n", &menuItemMapForContextMenus);

	//Free radio groups hashmap
	dumpHashmap("radioGroupMapForContextMenus", &radioGroupMapForContextMenus);
	printf("&radioGroupMapForContextMenus = %p\n", &radioGroupMapForContextMenus);

	//Free context menu map
	dumpHashmap("contextMenuMap", &contextMenuMap);
	printf("&contextMenuMap = %p\n", &contextMenuMap);
}

void parseContextMenus(struct Application *app) {

	// Allocation the hashmaps we need
	allocateContextMenuHashMaps(app);

	// Parse the context menu json
	app->processedContextMenus = json_decode(app->contextMenusAsJSON);

	if( app->processedContextMenus == NULL ) {
		// Parse error!
		Fatal(app, "Unable to parse Context Menus JSON: %s", app->contextMenusAsJSON);
		return;
	}

	JsonNode *contextMenuItems = json_find_member(app->processedContextMenus, "Items");
	if( contextMenuItems == NULL ) {
		// Parse error!
		Fatal(app, "Unable to find Items:", app->processedContextMenus);
		return;
	}
	// Iterate context menus
	JsonNode *contextMenu;
	json_foreach(contextMenu, contextMenuItems) {
		// Create a new menu
		id menu = createMenu(str(""));
		printf("Context menu NSMenu pointer = %p\n", menu);

		// parse the menu
		parseMenu(app, menu, contextMenu, &menuItemMapForContextMenus,
			"checkboxMenuCallbackForContextMenus:", "radioMenuCallbackForContextMenus:", "menuCallbackForContextMenus:");

		// Store the item in the context menu map
		printf("Putting context menu %p with key '%s' in contextMenuMap %p\n", menu, contextMenu->key, &contextMenuMap);
		hashmap_put(&contextMenuMap, (char*)contextMenu->key, strlen(contextMenu->key), menu);
	}

	dumpContextMenus(app);
}

void parseTrayData(struct Application *app) {

	// Allocate the hashmaps we need
	allocateTrayHashMaps(app);

	// Create a new menu
	id traymenu = createMenu(str(""));

	id statusItem = app->statusItem;

	// Create a new menu bar if we need to
	if ( statusItem == NULL ) {
		id statusBar = msg( c("NSStatusBar"), s("systemStatusBar") );
		statusItem = msg(statusBar, s("statusItemWithLength:"), -1.0);
		app->statusItem = statusItem;
		msg(statusItem, s("retain"));
		id statusBarButton = msg(statusItem, s("button"));

		// If we have a tray icon
		if ( trayIconLength > 0 ) {
			id imageData = msg(c("NSData"), s("dataWithBytes:length:"), trayIcon, trayIconLength);
			id trayImage = ALLOC("NSImage");
			msg(trayImage, s("initWithData:"), imageData);
			msg(statusBarButton, s("setImage:"), trayImage);
		}
	}

	// Parse the processed menu json
	app->processedTrayMenu = json_decode(app->trayMenuAsJSON);

	if( app->processedTrayMenu == NULL ) {
		// Parse error!
		Fatal(app, "Unable to parse Tray JSON: %s", app->trayMenuAsJSON);
		return;
	}


	// Pull out the Menu
	JsonNode *trayMenuData = json_find_member(app->processedTrayMenu, "Menu");
	if( trayMenuData == NULL ) {
		// Parse error!
		Fatal(app, "Unable to find Menu data: %s", app->processedTrayMenu);
		return;
	}


	parseMenu(app, traymenu, trayMenuData, &menuItemMapForTrayMenu,
	"checkboxMenuCallbackForTrayMenu:", "radioMenuCallbackForTrayMenu:", "menuCallbackForTrayMenu:");

	// Create the radiogroup cache
	JsonNode *radioGroups = json_find_member(app->processedTrayMenu, "RadioGroups");
	if( radioGroups == NULL ) {
		// Parse error!
		Fatal(app, "Unable to find RadioGroups data: %s", app->processedTrayMenu);
		return;
	}

	// Iterate radio groups
	JsonNode *radioGroup;
	json_foreach(radioGroup, radioGroups) {
		// Get item label
		processRadioGroup(radioGroup, &menuItemMapForTrayMenu, &radioGroupMapForTrayMenu);
	}


//     msg(statusBarButton, s("setImage:"),
//        msg(c("NSImage"), s("imageNamed:"),
//          msg(c("NSString"), s("stringWithUTF8String:"), tray->icon)));


	msg(statusItem, s("setMenu:"), traymenu);
 }


// UpdateTray replaces the current tray menu with the given one
void UpdateTray(struct Application *app, const char *trayMenuAsJSON) {
	ON_MAIN_THREAD (
		// Free up memory
		destroyTray(app);

		// Set the menu JSON
		app->trayMenuAsJSON = trayMenuAsJSON;
		parseTrayData(app);
	);
}

void UpdateContextMenus(struct Application *app, const char *contextMenusAsJSON) {
	ON_MAIN_THREAD (

		dumpContextMenus(app);

		// Free up memory
		destroyContextMenus(app);

		// Set the context menu JSON
		app->contextMenusAsJSON = contextMenusAsJSON;
		parseContextMenus(app);
	);
}


void Run(struct Application *app, int argc, char **argv) {

	// Process window decorations
	processDecorations(app);

	// Create the application
	createApplication(app);

	// Define delegate
	createDelegate(app);

	// Create the main window
	createMainWindow(app);

	// Create Content View
	id contentView = msg( ALLOC("NSView"), s("init") );
	msg(app->mainWindow, s("setContentView:"), contentView);

	// Set the main window title
	SetTitle(app, app->title);

	// Center Window
	Center(app);

	// Set Colour
	applyWindowColour(app);

	// Process translucency
	if (app->windowBackgroundIsTranslucent) {
		makeWindowBackgroundTranslucent(app);
	}

	// Setup webview
	id config = msg(c("WKWebViewConfiguration"), s("new"));
	msg(config, s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 1), str("suppressesIncrementalRendering"));
	if (app->devtools) {
	  Debug(app, "Enabling devtools");
	  enableBoolConfig(config, "developerExtrasEnabled");
	}
	app->config = config;

	id manager = msg(config, s("userContentController"));
	msg(manager, s("addScriptMessageHandler:name:"), app->delegate, str("external"));
	msg(manager, s("addScriptMessageHandler:name:"), app->delegate, str("completed"));
	app->manager = manager;

	id wkwebview = msg(c("WKWebView"), s("alloc"));
	app->wkwebview = wkwebview;

	msg(wkwebview, s("initWithFrame:configuration:"), CGRectMake(0, 0, 0, 0), config);

	msg(contentView, s("addSubview:"), wkwebview);
	msg(wkwebview, s("setAutoresizingMask:"), NSViewWidthSizable | NSViewHeightSizable);
	CGRect contentViewBounds = GET_BOUNDS(contentView);
	msg(wkwebview, s("setFrame:"), contentViewBounds );

	// Disable damn smart quotes
	// Credit: https://stackoverflow.com/a/31640511
	id userDefaults = msg(c("NSUserDefaults"), s("standardUserDefaults"));
	msg(userDefaults, s("setBool:forKey:"), NO, str("NSAutomaticQuoteSubstitutionEnabled"));

	// Setup drag message handler
	msg(manager, s("addScriptMessageHandler:name:"), app->delegate, str("windowDrag"));
	// Add mouse event hooks
	app->mouseDownMonitor = msg(c("NSEvent"), u("addLocalMonitorForEventsMatchingMask:handler:"), NSEventMaskLeftMouseDown, ^(id incomingEvent) {
		// Make sure the mouse click was in the window, not the tray
		id window = msg(incomingEvent, s("window"));
		if (window == app->mainWindow) {
			app->mouseEvent = incomingEvent;
		}
		return incomingEvent;
	});
	app->mouseUpMonitor = msg(c("NSEvent"), u("addLocalMonitorForEventsMatchingMask:handler:"), NSEventMaskLeftMouseUp, ^(id incomingEvent) {
		app->mouseEvent = NULL;
		ShowMouse();
		return incomingEvent;
	});

	// Setup context menu message handler
	msg(manager, s("addScriptMessageHandler:name:"), app->delegate, str("contextMenu"));

	// Toolbar
	if( app->useToolBar ) {
		Debug(app, "Setting Toolbar");
		id toolbar = msg(c("NSToolbar"),s("alloc"));
		msg(toolbar, s("initWithIdentifier:"), str("wails.toolbar"));
		msg(toolbar, s("autorelease"));

		// Separator
		if( app->hideToolbarSeparator ) {
			msg(toolbar, s("setShowsBaselineSeparator:"), NO);
		}

		msg(app->mainWindow, s("setToolbar:"), toolbar);
	}

	// Fix up resizing
	if (app->resizable == 0) {
		app->minHeight = app->maxHeight = app->height;
		app->minWidth = app->maxWidth = app->width;
	}
	setMinMaxSize(app);

	// Load HTML
	id html = msg(c("NSURL"), s("URLWithString:"), str(assets[0]));
	msg(wkwebview, s("loadRequest:"), msg(c("NSURLRequest"), s("requestWithURL:"), html));

	Debug(app, "Loading Internal Code");
	// We want to evaluate the internal code plus runtime before the assets
	const char *temp = concat(invoke, app->bindings);
	const char *internalCode = concat(temp, (const char*)&runtime);
	MEMFREE(temp);

	// Add code that sets up the initial state, EG: State Stores.
	temp = concat(internalCode, getInitialState(app));
	MEMFREE(internalCode);
	internalCode = temp;

	  // Loop over assets and build up one giant Mother Of All Evals
	int index = 1;
	while(1) {
		// Get next asset pointer
		const unsigned char *asset = assets[index];

		// If we have no more assets, break
		if (asset == 0x00) {
			break;
		}

		temp = concat(internalCode, (const char *)asset);
		MEMFREE(internalCode);
		internalCode = temp;
		index++;
	};

	// Disable context menu if not in debug mode
	if( debug != 1 ) {
		temp = concat(internalCode, "wails._.DisableDefaultContextMenu();");
		MEMFREE(internalCode);
		internalCode = temp;
	}

	// class_addMethod(delegateClass, s("applicationWillFinishLaunching:"), (IMP) willFinishLaunching, "@@:@");
	// Include callback after evaluation
	temp = concat(internalCode, "webkit.messageHandlers.completed.postMessage(true);");
	MEMFREE(internalCode);
	internalCode = temp;

	// const char *viewportScriptString = "var meta = document.createElement('meta'); meta.setAttribute('name', 'viewport'); meta.setAttribute('content', 'width=device-width'); meta.setAttribute('initial-scale', '1.0'); meta.setAttribute('maximum-scale', '1.0'); meta.setAttribute('minimum-scale', '1.0'); meta.setAttribute('user-scalable', 'no'); document.getElementsByTagName('head')[0].appendChild(meta);";
	// ExecJS(app, viewportScriptString);


	// This evaluates the MOAE once the Dom has finished loading
	msg(manager,
		s("addUserScript:"),
		msg(msg(c("WKUserScript"), s("alloc")),
					s("initWithSource:injectionTime:forMainFrameOnly:"),
					str(internalCode),
					1,
					1));



	if( app->webviewIsTranparent == 1 ) {
		msg(wkwebview, s("setValue:forKey:"), msg(c("NSNumber"), s("numberWithBool:"), 0), str("drawsBackground"));
	}

	// If we have a menu, process it
	if( app->menuAsJSON != NULL ) {
	  parseMenuData(app);
	}

	// If we have a tray menu, process it
	if( app->trayMenuAsJSON != NULL ) {
	  parseTrayData(app);
	}

	// If we have context menus, process them
	if( app->contextMenusAsJSON != NULL ) {
		parseContextMenus(app);
	}

	// We set it to be invisible by default. It will become visible when everything has initialised
	msg(app->mainWindow, s("setIsVisible:"), NO);

	// Finally call run
	Debug(app, "Run called");
	msg(app->application, s("run"));

	MEMFREE(internalCode);
}

#endif
