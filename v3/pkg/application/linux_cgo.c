#include "linux_cgo.h"

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
    gtk_drag_dest_set(widget, GTK_DEST_DEFAULT_MOTION | GTK_DEST_DEFAULT_HIGHLIGHT | GTK_DEST_DEFAULT_DROP, target, 1, GDK_ACTION_COPY);

    signal_connect(widget, "drag-data-received", on_data_received, data);
}