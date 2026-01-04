//go:build linux && !gtk3

#ifndef LINUX_CGO_GTK4_H
#define LINUX_CGO_GTK4_H

#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <jsc/jsc.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#include <errno.h>
#include <signal.h>
#include <string.h>

// Application flags for GTK4
#define APPLICATION_DEFAULT_FLAGS G_APPLICATION_DEFAULT_FLAGS

// ============================================================================
// Type definitions
// ============================================================================

typedef struct CallbackID {
    unsigned int value;
} CallbackID;

typedef struct WindowEvent {
    uint id;
    uint event;
} WindowEvent;

typedef struct Screen {
    const char* id;
    const char* name;
    int p_width;
    int p_height;
    int x;
    int y;
    int w_width;
    int w_height;
    int w_x;
    int w_y;
    float scaleFactor;
    double rotation;
    bool isPrimary;
} Screen;

typedef struct MenuItemData {
    guint id;
    GSimpleAction *action;
} MenuItemData;

typedef struct FileDialogData {
    guint request_id;
    gboolean allow_multiple;
    gboolean is_save;
    gboolean is_folder;
} FileDialogData;

typedef struct AlertDialogData {
    guint request_id;
} AlertDialogData;

// ============================================================================
// External Go callback declarations (implemented in Go with //export)
// ============================================================================

extern void dispatchOnMainThreadCallback(unsigned int);
extern void emit(WindowEvent* data);
extern gboolean handleCloseRequest(GtkWindow*, uintptr_t);
extern void handleNotifyState(GObject*, GParamSpec*, uintptr_t);
extern gboolean handleFocusEnter(GtkEventController*, uintptr_t);
extern gboolean handleFocusLeave(GtkEventController*, uintptr_t);
extern void handleLoadChanged(WebKitWebView*, WebKitLoadEvent, uintptr_t);
extern void handleButtonPressed(GtkGestureClick*, gint, gdouble, gdouble, uintptr_t);
extern void handleButtonReleased(GtkGestureClick*, gint, gdouble, gdouble, uintptr_t);
extern gboolean handleKeyPressed(GtkEventControllerKey*, guint, guint, GdkModifierType, uintptr_t);
extern void onProcessRequest(WebKitURISchemeRequest *request, uintptr_t user_data);
extern void sendMessageToBackend(WebKitUserContentManager *contentManager, JSCValue *value, void *data);
extern void menuActionActivated(guint id);
extern void fileDialogCallback(guint request_id, char **files, int count, gboolean cancelled);
extern void alertDialogCallback(guint request_id, int button_index);
extern void onDropEnter(uintptr_t);
extern void onDropLeave(uintptr_t);
extern void onDropMotion(gint, gint, uintptr_t);
extern void onDropFiles(char**, gint, gint, uintptr_t);

// Forward declaration for activate callback
void activateLinux(gpointer data);

// ============================================================================
// Main thread dispatch
// ============================================================================

void dispatchOnMainThread(unsigned int id);

// ============================================================================
// Signal handling
// ============================================================================

void install_signal_handlers(void);

// ============================================================================
// Object data helpers
// ============================================================================

void save_window_id(void *object, uint value);
guint get_window_id(void *object);
void save_webview_to_content_manager(void *contentManager, void *webview);
WebKitWebView* get_webview_from_content_manager(void *contentManager);

// ============================================================================
// Signal connection (wrapper for macro)
// ============================================================================

void signal_connect(void *widget, char *event, void *cb, void* data);

// ============================================================================
// WebView helpers
// ============================================================================

WebKitWebView* webkit_web_view(GtkWidget *webview);
GtkWidget* create_webview_with_user_content_manager(WebKitUserContentManager *manager);

// ============================================================================
// Menu system (GMenu/GAction for GTK4)
// ============================================================================

void init_app_action_group(void);
void set_app_menu_model(GMenu *menu);
GMenuItem* create_menu_item(const char *label, const char *action_name, guint item_id);
GMenuItem* create_check_menu_item(const char *label, const char *action_name, guint item_id, gboolean initial_state);
GMenuItem* create_radio_menu_item(const char *label, const char *action_name, const char *target, const char *initial_value, guint item_id);
GtkWidget* create_menu_bar_from_model(GMenu *menu_model);
GtkWidget* create_header_bar_with_menu(GMenu *menu_model);
void attach_action_group_to_widget(GtkWidget *widget);
void set_action_accelerator(GtkApplication *app, const char *action_name, const char *accel);
char* build_accelerator_string(guint key, GdkModifierType mods);
void set_action_enabled(const char *action_name, gboolean enabled);
void set_action_state(const char *action_name, gboolean state);
gboolean get_action_state(const char *action_name);
void menu_remove_item(GMenu *menu, gint position);
void menu_insert_item(GMenu *menu, gint position, GMenuItem *item);

// ============================================================================
// Window event controllers (GTK4 style)
// ============================================================================

void setupWindowEventControllers(GtkWindow *window, GtkWidget *webview, uintptr_t winID);

// ============================================================================
// Window drag/resize (GdkToplevel for GTK4)
// ============================================================================

void beginWindowDrag(GtkWindow *window, int button, double x, double y, guint32 timestamp);
void beginWindowResize(GtkWindow *window, GdkSurfaceEdge edge, int button, double x, double y, guint32 timestamp);

// ============================================================================
// Drag and drop (GtkDropTarget for GTK4)
// ============================================================================

void enableDND(GtkWidget *widget, gpointer data);
void disableDND(GtkWidget *widget, gpointer data);

// ============================================================================
// File dialogs (GtkFileDialog for GTK4)
// ============================================================================

GtkFileDialog* create_file_dialog(const char *title);
void add_file_filter(GtkFileDialog *dialog, GListStore *filters, const char *name, const char *pattern);
void set_file_dialog_filters(GtkFileDialog *dialog, GListStore *filters);
void show_open_file_dialog(GtkWindow *parent, GtkFileDialog *dialog, guint request_id, gboolean allow_multiple, gboolean is_folder);
void show_save_file_dialog(GtkWindow *parent, GtkFileDialog *dialog, guint request_id);

// ============================================================================
// Alert dialogs (GtkAlertDialog for GTK4)
// ============================================================================

void show_alert_dialog(GtkWindow *parent, const char *message, const char *detail,
                       const char **buttons, int button_count, int default_button,
                       int cancel_button, guint request_id);

// ============================================================================
// Misc
// ============================================================================

int GetNumScreens(void);

#endif // LINUX_CGO_GTK4_H
