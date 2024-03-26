#include "linux_cgo.h"

gboolean on_drop(GtkWidget *widget, GdkDragContext *context, gint x, gint y, guint time, gpointer user_data)
{
    GdkAtom target_type;
    GList *targets = gdk_drag_context_list_targets(context);

    if (targets == NULL || g_list_length(targets) <= 0)
    {
        g_print("no targets\n");
        return FALSE;
    }
    /* Choose the best target type */
    target_type = GDK_POINTER_TO_ATOM(g_list_nth_data(targets, 0));

    /* Request the data from the source. */
    gtk_drag_get_data(
        widget,      /* will receive 'drag-data-received' signal */
        context,     /* represents the current state of the DnD */
        target_type, /* the target type we want */
        time         /* time stamp */
    );
}

void on_data_received(GtkWidget *widget, GdkDragContext *context, gint x, gint y,
                      GtkSelectionData *selection_data, guint target_type, guint time,
                      gpointer data)
{
    gint length = gtk_selection_data_get_length(selection_data);

    if (length < 0)
    {
        g_print("DnD failed!\n");
        gtk_drag_finish(context, FALSE, FALSE, time);
    }

    gchar *uri_data = (gchar *)gtk_selection_data_get_data(selection_data);
    gchar **uri_list = g_uri_list_extract_uris(uri_data);
    
    onUriList(uri_list, data);

    g_strfreev(uri_list);
    gtk_drag_finish(context, TRUE, TRUE, time);
}

// drag and drop tutorial: https://wiki.gnome.org/Newcomers/OldDragNDropTutorial
void enableDND(GtkWidget *widget, gpointer data)
{
    GtkTargetEntry *target = gtk_target_entry_new("text/uri-list", 0, 0);
    gtk_drag_dest_set(widget, GTK_DEST_DEFAULT_MOTION | GTK_DEST_DEFAULT_HIGHLIGHT, target, 1, GDK_ACTION_COPY);

    signal_connect(widget, "drag-drop", on_drop, data);
    signal_connect(widget, "drag-data-received", on_data_received, data);
}