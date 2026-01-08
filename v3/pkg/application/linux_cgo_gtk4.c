//go:build linux && !gtk3

#include "linux_cgo_gtk4.h"

#ifdef WAILS_GTK_DEBUG
#define DEBUG_LOG(fmt, ...) fprintf(stderr, "[GTK4] " fmt "\n", ##__VA_ARGS__)
#else
#define DEBUG_LOG(fmt, ...)
#endif

// ============================================================================
// Global state
// ============================================================================

static GMenu *app_menu_model = NULL;
static GSimpleActionGroup *app_action_group = NULL;

// ============================================================================
// Main thread dispatch
// ============================================================================

static gboolean dispatchCallback(gpointer data) {
    struct CallbackID *args = data;
    unsigned int cid = args->value;
    dispatchOnMainThreadCallback(cid);
    free(args);
    return G_SOURCE_REMOVE;
}

void dispatchOnMainThread(unsigned int id) {
    CallbackID *args = malloc(sizeof(CallbackID));
    args->value = id;
    g_idle_add((GSourceFunc)dispatchCallback, (gpointer)args);
}

// ============================================================================
// Signal handling
// ============================================================================

static void fix_signal(int signum) {
    struct sigaction st;

    if (sigaction(signum, NULL, &st) < 0) {
        goto fix_signal_error;
    }
    st.sa_flags |= SA_ONSTACK;
    if (sigaction(signum, &st, NULL) < 0) {
        goto fix_signal_error;
    }
    return;
fix_signal_error:
    fprintf(stderr, "error fixing handler for signal %d, please "
            "report this issue to "
            "https://github.com/wailsapp/wails: %s\n",
            signum, strerror(errno));
}

void install_signal_handlers(void) {
#if defined(SIGCHLD)
    fix_signal(SIGCHLD);
#endif
#if defined(SIGHUP)
    fix_signal(SIGHUP);
#endif
#if defined(SIGINT)
    fix_signal(SIGINT);
#endif
#if defined(SIGQUIT)
    fix_signal(SIGQUIT);
#endif
#if defined(SIGABRT)
    fix_signal(SIGABRT);
#endif
#if defined(SIGFPE)
    fix_signal(SIGFPE);
#endif
#if defined(SIGTERM)
    fix_signal(SIGTERM);
#endif
#if defined(SIGBUS)
    fix_signal(SIGBUS);
#endif
#if defined(SIGSEGV)
    fix_signal(SIGSEGV);
#endif
#if defined(SIGXCPU)
    fix_signal(SIGXCPU);
#endif
#if defined(SIGXFSZ)
    fix_signal(SIGXFSZ);
#endif
}

// ============================================================================
// Object data helpers
// ============================================================================

void save_window_id(void *object, uint value) {
    g_object_set_data((GObject *)object, "windowid", GUINT_TO_POINTER((guint)value));
}

guint get_window_id(void *object) {
    return GPOINTER_TO_UINT(g_object_get_data((GObject *)object, "windowid"));
}

void save_webview_to_content_manager(void *contentManager, void *webview) {
    g_object_set_data(G_OBJECT((WebKitUserContentManager *)contentManager), "webview", webview);
}

WebKitWebView* get_webview_from_content_manager(void *contentManager) {
    return WEBKIT_WEB_VIEW(g_object_get_data(G_OBJECT(contentManager), "webview"));
}

// ============================================================================
// Signal connection (wrapper for macro)
// ============================================================================

void signal_connect(void *widget, char *event, void *cb, uintptr_t data) {
    g_signal_connect(widget, event, cb, (gpointer)data);
}

// ============================================================================
// WebView helpers
// ============================================================================

WebKitWebView* webkit_web_view(GtkWidget *webview) {
    return WEBKIT_WEB_VIEW(webview);
}

// WebKitGTK 6.0: webkit_web_view_new_with_user_content_manager() was removed
// Use g_object_new() with the "user-content-manager" property instead
GtkWidget* create_webview_with_user_content_manager(WebKitUserContentManager *manager) {
    return GTK_WIDGET(g_object_new(WEBKIT_TYPE_WEB_VIEW,
                                   "user-content-manager", manager,
                                   NULL));
}

// ============================================================================
// Menu system (GMenu/GAction for GTK4)
// ============================================================================

static void on_action_activated(GSimpleAction *action, GVariant *parameter, gpointer user_data) {
    MenuItemData *data = (MenuItemData *)user_data;
    if (data != NULL) {
        menuActionActivated(data->id);
    }
}

void init_app_action_group(void) {
    if (app_action_group == NULL) {
        app_action_group = g_simple_action_group_new();
    }
}

void set_app_menu_model(GMenu *menu) {
    app_menu_model = menu;
}

GMenuItem* create_menu_item(const char *label, const char *action_name, guint item_id) {
    init_app_action_group();

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    GMenuItem *item = g_menu_item_new(label, full_action_name);

    GSimpleAction *action = g_simple_action_new(action_name, NULL);
    MenuItemData *data = g_new0(MenuItemData, 1);
    data->id = item_id;
    data->action = action;
    g_signal_connect(action, "activate", G_CALLBACK(on_action_activated), data);
    g_action_map_add_action(G_ACTION_MAP(app_action_group), G_ACTION(action));

    return item;
}

GMenuItem* create_check_menu_item(const char *label, const char *action_name, guint item_id, gboolean initial_state) {
    init_app_action_group();

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    GMenuItem *item = g_menu_item_new(label, full_action_name);

    GSimpleAction *action = g_simple_action_new_stateful(action_name, NULL, g_variant_new_boolean(initial_state));
    MenuItemData *data = g_new0(MenuItemData, 1);
    data->id = item_id;
    data->action = action;
    g_signal_connect(action, "activate", G_CALLBACK(on_action_activated), data);
    g_action_map_add_action(G_ACTION_MAP(app_action_group), G_ACTION(action));

    return item;
}

static void on_radio_action_activated(GSimpleAction *action, GVariant *parameter, gpointer user_data) {
    const gchar *target = g_variant_get_string(parameter, NULL);
    g_simple_action_set_state(action, g_variant_new_string(target));
    guint item_id = (guint)atoi(target);
    menuActionActivated(item_id);
}

GMenuItem* create_radio_menu_item(const char *label, const char *action_name, const char *target, const char *initial_value, guint item_id) {
    init_app_action_group();

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    GMenuItem *item = g_menu_item_new(label, NULL);
    g_menu_item_set_action_and_target(item, full_action_name, "s", target);

    GAction *existing = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (existing == NULL) {
        GSimpleAction *action = g_simple_action_new_stateful(
            action_name,
            G_VARIANT_TYPE_STRING,
            g_variant_new_string(initial_value)
        );
        MenuItemData *data = g_new0(MenuItemData, 1);
        data->id = item_id;
        data->action = action;
        g_signal_connect(action, "activate", G_CALLBACK(on_radio_action_activated), data);
        g_action_map_add_action(G_ACTION_MAP(app_action_group), G_ACTION(action));
    }

    return item;
}

GtkWidget* create_menu_bar_from_model(GMenu *menu_model) {
    return gtk_popover_menu_bar_new_from_model(G_MENU_MODEL(menu_model));
}

GtkWidget* create_header_bar_with_menu(GMenu *menu_model) {
    GtkWidget *header_bar = gtk_header_bar_new();
    
    GtkWidget *menu_button = gtk_menu_button_new();
    gtk_menu_button_set_icon_name(GTK_MENU_BUTTON(menu_button), "open-menu-symbolic");
    gtk_menu_button_set_menu_model(GTK_MENU_BUTTON(menu_button), G_MENU_MODEL(menu_model));
    gtk_widget_set_tooltip_text(menu_button, "Main Menu");
    gtk_accessible_update_property(GTK_ACCESSIBLE(menu_button),
        GTK_ACCESSIBLE_PROPERTY_LABEL, "Main Menu", -1);
    
    gtk_header_bar_pack_end(GTK_HEADER_BAR(header_bar), menu_button);
    
    return header_bar;
}

void attach_action_group_to_widget(GtkWidget *widget) {
    init_app_action_group();
    gtk_widget_insert_action_group(widget, "app", G_ACTION_GROUP(app_action_group));
}

void set_action_accelerator(GtkApplication *app, const char *action_name, const char *accel) {
    if (app == NULL || accel == NULL || strlen(accel) == 0) return;

    char full_action_name[256];
    snprintf(full_action_name, sizeof(full_action_name), "app.%s", action_name);

    const char *accels[] = { accel, NULL };
    gtk_application_set_accels_for_action(app, full_action_name, accels);
}

char* build_accelerator_string(guint key, GdkModifierType mods) {
    return gtk_accelerator_name(key, mods);
}

void set_action_enabled(const char *action_name, gboolean enabled) {
    if (app_action_group == NULL) return;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL && G_IS_SIMPLE_ACTION(action)) {
        g_simple_action_set_enabled(G_SIMPLE_ACTION(action), enabled);
    }
}

void set_action_state(const char *action_name, gboolean state) {
    if (app_action_group == NULL) return;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL && G_IS_SIMPLE_ACTION(action)) {
        g_simple_action_set_state(G_SIMPLE_ACTION(action), g_variant_new_boolean(state));
    }
}

gboolean get_action_state(const char *action_name) {
    if (app_action_group == NULL) return FALSE;
    GAction *action = g_action_map_lookup_action(G_ACTION_MAP(app_action_group), action_name);
    if (action != NULL) {
        GVariant *state = g_action_get_state(action);
        if (state != NULL) {
            gboolean result = g_variant_get_boolean(state);
            g_variant_unref(state);
            return result;
        }
    }
    return FALSE;
}

void menu_remove_item(GMenu *menu, gint position) {
    g_menu_remove(menu, position);
}

void menu_insert_item(GMenu *menu, gint position, GMenuItem *item) {
    g_menu_insert_item(menu, position, item);
}

// ============================================================================
// Window event controllers (GTK4 style)
// ============================================================================

void setupWindowEventControllers(GtkWindow *window, GtkWidget *webview, uintptr_t winID) {
    // Close request (replaces delete-event)
    g_signal_connect(window, "close-request", G_CALLBACK(handleCloseRequest), (gpointer)winID);

    // Window state changes (maximize, fullscreen, etc)
    g_signal_connect(window, "notify::maximized", G_CALLBACK(handleNotifyState), (gpointer)winID);
    g_signal_connect(window, "notify::fullscreened", G_CALLBACK(handleNotifyState), (gpointer)winID);

    // Focus controller for window
    GtkEventController *focus_controller = gtk_event_controller_focus_new();
    gtk_widget_add_controller(GTK_WIDGET(window), focus_controller);
    g_signal_connect(focus_controller, "enter", G_CALLBACK(handleFocusEnter), (gpointer)winID);
    g_signal_connect(focus_controller, "leave", G_CALLBACK(handleFocusLeave), (gpointer)winID);

    // Click gesture for webview (button press/release)
    GtkGesture *click_gesture = gtk_gesture_click_new();
    gtk_gesture_single_set_button(GTK_GESTURE_SINGLE(click_gesture), 0); // Listen to all buttons
    gtk_widget_add_controller(webview, GTK_EVENT_CONTROLLER(click_gesture));
    g_signal_connect(click_gesture, "pressed", G_CALLBACK(handleButtonPressed), (gpointer)winID);
    g_signal_connect(click_gesture, "released", G_CALLBACK(handleButtonReleased), (gpointer)winID);

    // Key controller for webview
    GtkEventController *key_controller = gtk_event_controller_key_new();
    gtk_widget_add_controller(webview, key_controller);
    g_signal_connect(key_controller, "key-pressed", G_CALLBACK(handleKeyPressed), (gpointer)winID);
}

// ============================================================================
// Window drag/resize (GdkToplevel for GTK4)
// ============================================================================

void beginWindowDrag(GtkWindow *window, int button, double x, double y, guint32 timestamp) {
    GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));
    if (native == NULL) return;

    GdkSurface *surface = gtk_native_get_surface(native);
    if (surface == NULL || !GDK_IS_TOPLEVEL(surface)) return;

    GdkToplevel *toplevel = GDK_TOPLEVEL(surface);
    GdkDevice *device = NULL;
    GdkDisplay *display = gdk_surface_get_display(surface);
    GdkSeat *seat = gdk_display_get_default_seat(display);
    if (seat) {
        device = gdk_seat_get_pointer(seat);
    }

    gdk_toplevel_begin_move(toplevel, device, button, x, y, timestamp);
}

void beginWindowResize(GtkWindow *window, GdkSurfaceEdge edge, int button, double x, double y, guint32 timestamp) {
    GtkNative *native = gtk_widget_get_native(GTK_WIDGET(window));
    if (native == NULL) return;

    GdkSurface *surface = gtk_native_get_surface(native);
    if (surface == NULL || !GDK_IS_TOPLEVEL(surface)) return;

    GdkToplevel *toplevel = GDK_TOPLEVEL(surface);
    GdkDevice *device = NULL;
    GdkDisplay *display = gdk_surface_get_display(surface);
    GdkSeat *seat = gdk_display_get_default_seat(display);
    if (seat) {
        device = gdk_seat_get_pointer(seat);
    }

    gdk_toplevel_begin_resize(toplevel, edge, device, button, x, y, timestamp);
}

// ============================================================================
// Drag and drop (GtkDropTarget for GTK4)
// ============================================================================

static GdkDragAction on_drop_enter(GtkDropTarget *target, gdouble x, gdouble y, gpointer data) {
    onDropEnter((uintptr_t)data);
    return GDK_ACTION_COPY;
}

static void on_drop_leave(GtkDropTarget *target, gpointer data) {
    onDropLeave((uintptr_t)data);
}

static GdkDragAction on_drop_motion(GtkDropTarget *target, gdouble x, gdouble y, gpointer data) {
    onDropMotion((gint)x, (gint)y, (uintptr_t)data);
    return GDK_ACTION_COPY;
}

static gboolean on_drop(GtkDropTarget *target, const GValue *value, gdouble x, gdouble y, gpointer data) {
    if (!G_VALUE_HOLDS(value, GDK_TYPE_FILE_LIST)) {
        return FALSE;
    }

    GSList *file_list = g_value_get_boxed(value);
    if (file_list == NULL) {
        return FALSE;
    }

    guint count = g_slist_length(file_list);
    if (count == 0) {
        return FALSE;
    }

    char **paths = g_new0(char*, count + 1);
    guint i = 0;
    for (GSList *l = file_list; l != NULL; l = l->next) {
        GFile *file = G_FILE(l->data);
        paths[i++] = g_file_get_path(file);
    }
    paths[count] = NULL;

    onDropFiles(paths, (gint)x, (gint)y, (uintptr_t)data);

    for (i = 0; i < count; i++) {
        g_free(paths[i]);
    }
    g_free(paths);

    return TRUE;
}

void enableDND(GtkWidget *widget, uintptr_t winID) {
    GtkDropTarget *target = gtk_drop_target_new(GDK_TYPE_FILE_LIST, GDK_ACTION_COPY);
    g_signal_connect(target, "enter", G_CALLBACK(on_drop_enter), (gpointer)winID);
    g_signal_connect(target, "leave", G_CALLBACK(on_drop_leave), (gpointer)winID);
    g_signal_connect(target, "motion", G_CALLBACK(on_drop_motion), (gpointer)winID);
    g_signal_connect(target, "drop", G_CALLBACK(on_drop), (gpointer)winID);
    gtk_widget_add_controller(widget, GTK_EVENT_CONTROLLER(target));
}

void disableDND(GtkWidget *widget, uintptr_t winID) {
    // In GTK4, we don't add a drop target to block drops
    // The default behavior is to not accept drops
}

// ============================================================================
// File dialogs (GtkFileDialog for GTK4)
// ============================================================================

GtkFileDialog* create_file_dialog(const char *title) {
    GtkFileDialog *dialog = gtk_file_dialog_new();
    gtk_file_dialog_set_title(dialog, title);
    return dialog;
}

void add_file_filter(GtkFileDialog *dialog, GListStore *filters, const char *name, const char *pattern) {
    GtkFileFilter *filter = gtk_file_filter_new();
    gtk_file_filter_set_name(filter, name);

    gchar **patterns = g_strsplit(pattern, ";", -1);
    for (int i = 0; patterns[i] != NULL; i++) {
        gchar *p = g_strstrip(patterns[i]);
        if (strlen(p) > 0) {
            gtk_file_filter_add_pattern(filter, p);
        }
    }
    g_strfreev(patterns);

    g_list_store_append(filters, filter);
    g_object_unref(filter);
}

void set_file_dialog_filters(GtkFileDialog *dialog, GListStore *filters) {
    gtk_file_dialog_set_filters(dialog, G_LIST_MODEL(filters));
}

// File dialog callbacks

static void on_file_dialog_open_finish(GObject *source, GAsyncResult *result, gpointer user_data) {
    FileDialogData *data = (FileDialogData *)user_data;
    GtkFileDialog *dialog = GTK_FILE_DIALOG(source);
    GError *error = NULL;

    GFile *file = gtk_file_dialog_open_finish(dialog, result, &error);

    if (error != NULL) {
        DEBUG_LOG("open_finish error: %s", error->message);
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
        g_error_free(error);
    } else if (file != NULL) {
        char *path = g_file_get_path(file);
        char *files[1] = { path };
        fileDialogCallback(data->request_id, files, 1, FALSE);
        g_free(path);
        g_object_unref(file);
    } else {
        // Cancelled - no error but no file
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
    }

    g_free(data);
}

static void on_file_dialog_open_multiple_finish(GObject *source, GAsyncResult *result, gpointer user_data) {
    FileDialogData *data = (FileDialogData *)user_data;
    GtkFileDialog *dialog = GTK_FILE_DIALOG(source);
    GError *error = NULL;

    GListModel *files = gtk_file_dialog_open_multiple_finish(dialog, result, &error);

    if (error != NULL) {
        DEBUG_LOG("open_multiple_finish error: %s", error->message);
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
        g_error_free(error);
    } else if (files != NULL) {
        guint n = g_list_model_get_n_items(files);
        char **paths = g_new0(char*, n + 1);

        for (guint i = 0; i < n; i++) {
            GFile *file = G_FILE(g_list_model_get_item(files, i));
            paths[i] = g_file_get_path(file);
            g_object_unref(file);
        }

        fileDialogCallback(data->request_id, paths, (int)n, FALSE);

        for (guint i = 0; i < n; i++) {
            g_free(paths[i]);
        }
        g_free(paths);
        g_object_unref(files);
    } else {
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
    }

    g_free(data);
}

static void on_file_dialog_select_folder_finish(GObject *source, GAsyncResult *result, gpointer user_data) {
    FileDialogData *data = (FileDialogData *)user_data;
    GtkFileDialog *dialog = GTK_FILE_DIALOG(source);
    GError *error = NULL;

    GFile *file = gtk_file_dialog_select_folder_finish(dialog, result, &error);

    if (error != NULL) {
        DEBUG_LOG("select_folder_finish error: %s", error->message);
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
        g_error_free(error);
    } else if (file != NULL) {
        char *path = g_file_get_path(file);
        char *files[1] = { path };
        fileDialogCallback(data->request_id, files, 1, FALSE);
        g_free(path);
        g_object_unref(file);
    } else {
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
    }

    g_free(data);
}

static void on_file_dialog_select_multiple_folders_finish(GObject *source, GAsyncResult *result, gpointer user_data) {
    FileDialogData *data = (FileDialogData *)user_data;
    GtkFileDialog *dialog = GTK_FILE_DIALOG(source);
    GError *error = NULL;

    GListModel *files = gtk_file_dialog_select_multiple_folders_finish(dialog, result, &error);

    if (error != NULL) {
        DEBUG_LOG("select_multiple_folders_finish error: %s", error->message);
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
        g_error_free(error);
    } else if (files != NULL) {
        guint n = g_list_model_get_n_items(files);
        char **paths = g_new0(char*, n + 1);

        for (guint i = 0; i < n; i++) {
            GFile *file = G_FILE(g_list_model_get_item(files, i));
            paths[i] = g_file_get_path(file);
            g_object_unref(file);
        }

        fileDialogCallback(data->request_id, paths, (int)n, FALSE);

        for (guint i = 0; i < n; i++) {
            g_free(paths[i]);
        }
        g_free(paths);
        g_object_unref(files);
    } else {
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
    }

    g_free(data);
}

static void on_file_dialog_save_finish(GObject *source, GAsyncResult *result, gpointer user_data) {
    FileDialogData *data = (FileDialogData *)user_data;
    GtkFileDialog *dialog = GTK_FILE_DIALOG(source);
    GError *error = NULL;

    GFile *file = gtk_file_dialog_save_finish(dialog, result, &error);

    if (error != NULL) {
        DEBUG_LOG("save_finish error: %s", error->message);
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
        g_error_free(error);
    } else if (file != NULL) {
        char *path = g_file_get_path(file);
        char *files[1] = { path };
        fileDialogCallback(data->request_id, files, 1, FALSE);
        g_free(path);
        g_object_unref(file);
    } else {
        fileDialogCallback(data->request_id, NULL, 0, TRUE);
    }

    g_free(data);
}

void show_open_file_dialog(GtkWindow *parent, GtkFileDialog *dialog, guint request_id, gboolean allow_multiple, gboolean is_folder) {
    FileDialogData *data = g_new0(FileDialogData, 1);
    data->request_id = request_id;
    data->allow_multiple = allow_multiple;
    data->is_folder = is_folder;

    if (is_folder && allow_multiple) {
        gtk_file_dialog_select_multiple_folders(dialog, parent, NULL, on_file_dialog_select_multiple_folders_finish, data);
    } else if (is_folder) {
        gtk_file_dialog_select_folder(dialog, parent, NULL, on_file_dialog_select_folder_finish, data);
    } else if (allow_multiple) {
        gtk_file_dialog_open_multiple(dialog, parent, NULL, on_file_dialog_open_multiple_finish, data);
    } else {
        gtk_file_dialog_open(dialog, parent, NULL, on_file_dialog_open_finish, data);
    }
}

void show_save_file_dialog(GtkWindow *parent, GtkFileDialog *dialog, guint request_id) {
    FileDialogData *data = g_new0(FileDialogData, 1);
    data->request_id = request_id;
    data->is_save = TRUE;

    gtk_file_dialog_save(dialog, parent, NULL, on_file_dialog_save_finish, data);
}

// ============================================================================
// Custom Message Dialogs (GtkWindow-based for proper styling)
// ============================================================================

typedef struct {
    GtkWindow *dialog;
    guint request_id;
    int button_count;
    int cancel_button;
    GtkWidget **buttons;
} MessageDialogData;

static void message_dialog_cleanup(MessageDialogData *data) {
    if (data->buttons != NULL) {
        g_free(data->buttons);
    }
    g_free(data);
}

static void on_message_dialog_button_clicked(GtkButton *button, gpointer user_data) {
    MessageDialogData *data = (MessageDialogData *)user_data;
    int index = GPOINTER_TO_INT(g_object_get_data(G_OBJECT(button), "button-index"));
    
    alertDialogCallback(data->request_id, index);
    gtk_window_destroy(data->dialog);
    message_dialog_cleanup(data);
}

static gboolean on_message_dialog_close(GtkWindow *window, gpointer user_data) {
    MessageDialogData *data = (MessageDialogData *)user_data;
    int result = (data->cancel_button >= 0) ? data->cancel_button : -1;
    alertDialogCallback(data->request_id, result);
    message_dialog_cleanup(data);
    return TRUE;
}

static gboolean on_message_dialog_key_pressed(GtkEventControllerKey *controller,
                                               guint keyval, guint keycode,
                                               GdkModifierType state, gpointer user_data) {
    MessageDialogData *data = (MessageDialogData *)user_data;
    
    if (keyval == GDK_KEY_Escape && data->cancel_button >= 0 && data->cancel_button < data->button_count) {
        gtk_widget_activate(data->buttons[data->cancel_button]);
        return TRUE;
    }
    return FALSE;
}

void show_message_dialog(GtkWindow *parent, const char *heading, const char *body,
                         const char *icon_name, const unsigned char *icon_data, int icon_data_len,
                         const char **buttons, int button_count,
                         int default_button, int cancel_button, int destructive_button,
                         guint request_id) {
    
    GtkWidget *dialog = gtk_window_new();
    gtk_window_set_modal(GTK_WINDOW(dialog), TRUE);
    gtk_window_set_resizable(GTK_WINDOW(dialog), FALSE);
    gtk_window_set_decorated(GTK_WINDOW(dialog), TRUE);
    gtk_widget_add_css_class(dialog, "message");
    gtk_widget_set_size_request(dialog, 300, -1);
    
    if (parent != NULL) {
        gtk_window_set_transient_for(GTK_WINDOW(dialog), parent);
    }
    
    MessageDialogData *data = g_new0(MessageDialogData, 1);
    data->dialog = GTK_WINDOW(dialog);
    data->request_id = request_id;
    data->button_count = button_count;
    data->cancel_button = cancel_button;
    data->buttons = (button_count > 0) ? g_new0(GtkWidget*, button_count) : NULL;
    
    GtkWidget *content = gtk_box_new(GTK_ORIENTATION_VERTICAL, 12);
    gtk_widget_set_margin_start(content, 24);
    gtk_widget_set_margin_end(content, 24);
    gtk_widget_set_margin_top(content, 24);
    gtk_widget_set_margin_bottom(content, 24);
    
    const int max_icon_size = 64;
    GtkWidget *icon_widget = NULL;
    if (icon_data != NULL && icon_data_len > 0) {
        GBytes *bytes = g_bytes_new(icon_data, icon_data_len);
        GdkTexture *texture = gdk_texture_new_from_bytes(bytes, NULL);
        g_bytes_unref(bytes);
        if (texture != NULL) {
            GtkWidget *image = gtk_image_new_from_paintable(GDK_PAINTABLE(texture));
            gtk_image_set_pixel_size(GTK_IMAGE(image), max_icon_size);
            icon_widget = image;
            g_object_unref(texture);
        }
    } else if (icon_name != NULL && strlen(icon_name) > 0) {
        icon_widget = gtk_image_new_from_icon_name(icon_name);
        gtk_image_set_pixel_size(GTK_IMAGE(icon_widget), max_icon_size);
    }
    
    if (icon_widget != NULL) {
        gtk_widget_set_halign(icon_widget, GTK_ALIGN_CENTER);
        gtk_widget_set_margin_bottom(icon_widget, 12);
        gtk_box_append(GTK_BOX(content), icon_widget);
    }
    
    if (heading != NULL && strlen(heading) > 0) {
        GtkWidget *heading_label = gtk_label_new(heading);
        gtk_widget_add_css_class(heading_label, "title-2");
        gtk_widget_set_halign(heading_label, GTK_ALIGN_CENTER);
        gtk_label_set_wrap(GTK_LABEL(heading_label), TRUE);
        gtk_label_set_max_width_chars(GTK_LABEL(heading_label), 50);
        gtk_box_append(GTK_BOX(content), heading_label);
    }
    
    if (body != NULL && strlen(body) > 0) {
        GtkWidget *body_label = gtk_label_new(body);
        gtk_widget_set_halign(body_label, GTK_ALIGN_CENTER);
        gtk_label_set_wrap(GTK_LABEL(body_label), TRUE);
        gtk_label_set_max_width_chars(GTK_LABEL(body_label), 50);
        gtk_widget_add_css_class(body_label, "dim-label");
        gtk_box_append(GTK_BOX(content), body_label);
    }
    
    if (button_count > 0) {
        GtkWidget *button_box = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 8);
        gtk_widget_set_halign(button_box, GTK_ALIGN_CENTER);
        gtk_widget_set_margin_top(button_box, 12);
        
        for (int i = 0; i < button_count; i++) {
            GtkWidget *btn = gtk_button_new_with_label(buttons[i]);
            g_object_set_data(G_OBJECT(btn), "button-index", GINT_TO_POINTER(i));
            g_signal_connect(btn, "clicked", G_CALLBACK(on_message_dialog_button_clicked), data);
            data->buttons[i] = btn;
            
            if (i == default_button) {
                gtk_widget_add_css_class(btn, "suggested-action");
                gtk_widget_add_css_class(btn, "default");
            }
            if (i == destructive_button) {
                gtk_widget_add_css_class(btn, "destructive-action");
            }
            
            gtk_box_append(GTK_BOX(button_box), btn);
        }
        
        gtk_box_append(GTK_BOX(content), button_box);
    }
    
    gtk_window_set_child(GTK_WINDOW(dialog), content);
    
    GtkEventController *key_controller = gtk_event_controller_key_new();
    g_signal_connect(key_controller, "key-pressed", G_CALLBACK(on_message_dialog_key_pressed), data);
    gtk_widget_add_controller(dialog, key_controller);
    
    g_signal_connect(dialog, "close-request", G_CALLBACK(on_message_dialog_close), data);
    
    gtk_window_present(GTK_WINDOW(dialog));
    
    if (default_button >= 0 && default_button < button_count) {
        gtk_window_set_default_widget(GTK_WINDOW(dialog), data->buttons[default_button]);
        gtk_widget_grab_focus(data->buttons[default_button]);
    }
}

// ============================================================================
// Misc
// ============================================================================

int GetNumScreens(void) {
    return 0;
}
