#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>
#ifdef G_APPLICATION_DEFAULT_FLAGS
#define APPLICATION_DEFAULT_FLAGS G_APPLICATION_DEFAULT_FLAGS
#else
#define APPLICATION_DEFAULT_FLAGS G_APPLICATION_FLAGS_NONE
#endif

typedef struct CallbackID
{
	unsigned int value;
} CallbackID;

extern void dispatchOnMainThreadCallback(unsigned int);

static gboolean dispatchCallback(gpointer data)
{
	struct CallbackID *args = data;
	unsigned int cid = args->value;
	dispatchOnMainThreadCallback(cid);
	free(args);

	return G_SOURCE_REMOVE;
};

static void dispatchOnMainThread(unsigned int id)
{
	CallbackID *args = malloc(sizeof(CallbackID));
	args->value = id;
	g_idle_add((GSourceFunc)dispatchCallback, (gpointer)args);
}

typedef struct WindowEvent
{
	uint id;
	uint event;
} WindowEvent;

// exported below
void activateLinux(gpointer data);
extern void emit(WindowEvent *data);
extern gboolean handleDeleteEvent(GtkWidget *, GdkEvent *, uintptr_t);
extern gboolean handleFocusEvent(GtkWidget *, GdkEvent *, uintptr_t);
extern void handleLoadChanged(WebKitWebView *, WebKitLoadEvent, uintptr_t);
void handleClick(void *);
extern gboolean onButtonEvent(GtkWidget *widget, GdkEventButton *event, uintptr_t user_data);
extern gboolean onMenuButtonEvent(GtkWidget *widget, GdkEventButton *event, uintptr_t user_data);
extern void onUriList(char **extracted, gpointer data);
extern gboolean onKeyPressEvent(GtkWidget *widget, GdkEventKey *event, uintptr_t user_data);
extern void onProcessRequest(WebKitURISchemeRequest *request, uintptr_t user_data);
extern void sendMessageToBackend(WebKitUserContentManager *contentManager, WebKitJavascriptResult *result, void *data);
// exported below (end)

static void signal_connect(void *widget, char *event, void *cb, void *data)
{
	// g_signal_connect is a macro and can't be called directly
	g_signal_connect(widget, event, cb, data);
}

static WebKitWebView *webkit_web_view(GtkWidget *webview)
{
	return WEBKIT_WEB_VIEW(webview);
}

static void *new_message_dialog(GtkWindow *parent, const gchar *msg, int dialogType, bool hasButtons)
{
	// gtk_message_dialog_new is variadic!  Can't call from cgo directly
	GtkWidget *dialog;
	int buttonMask;

	// buttons will be added after creation
	buttonMask = GTK_BUTTONS_OK;
	if (hasButtons)
	{
		buttonMask = GTK_BUTTONS_NONE;
	}

	dialog = gtk_message_dialog_new(
		parent,
		GTK_DIALOG_MODAL | GTK_DIALOG_DESTROY_WITH_PARENT,
		dialogType,
		buttonMask,
		"%s",
		msg);

	// g_signal_connect_swapped (dialog,
	//                           "response",
	//                           G_CALLBACK (callback),
	//                           dialog);
	return dialog;
};

extern void messageDialogCB(gint button);

static void *gtkFileChooserDialogNew(char *title, GtkWindow *window, GtkFileChooserAction action, char *cancelLabel, char *acceptLabel)
{
	// gtk_file_chooser_dialog_new is variadic!  Can't call from cgo directly
	return (GtkFileChooser *)gtk_file_chooser_dialog_new(
		title,
		window,
		action,
		cancelLabel,
		GTK_RESPONSE_CANCEL,
		acceptLabel,
		GTK_RESPONSE_ACCEPT,
		NULL);
}

typedef struct Screen
{
	const char *id;
	const char *name;
	int p_width;
	int p_height;
	int width;
	int height;
	int x;
	int y;
	int w_width;
	int w_height;
	int w_x;
	int w_y;
	float scale;
	double rotation;
	bool isPrimary;
} Screen;

// CREDIT: https://github.com/rainycape/magick
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

static void fix_signal(int signum)
{
	struct sigaction st;

	if (sigaction(signum, NULL, &st) < 0)
	{
		goto fix_signal_error;
	}
	st.sa_flags |= SA_ONSTACK;
	if (sigaction(signum, &st, NULL) < 0)
	{
		goto fix_signal_error;
	}
	return;
fix_signal_error:
	fprintf(stderr, "error fixing handler for signal %d, please "
					"report this issue to "
					"https://github.com/wailsapp/wails: %s\n",
			signum, strerror(errno));
}

static void install_signal_handlers()
{
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

static int GetNumScreens()
{
	return 0;
}

void on_data_received(GtkWidget *widget, GdkDragContext *context, gint x, gint y,
					  GtkSelectionData *selection_data, guint target_type, guint time,
					  gpointer data);
void enableDND(GtkWidget *widget, gpointer data);