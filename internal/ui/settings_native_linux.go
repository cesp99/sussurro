//go:build linux

package ui

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

static void show_window(void *win) {
    gtk_widget_show_all(GTK_WIDGET(win));
    gtk_window_present(GTK_WINDOW(win));
}
static void hide_window(void *win) {
    gtk_widget_hide(GTK_WIDGET(win));
}

// Intercept the WM "X" close button: hide instead of destroy.
// Returning TRUE suppresses the default action (gtk_widget_destroy),
// keeping the window alive so it can be shown again later.
static gboolean on_settings_delete(GtkWidget *win, GdkEvent *ev, gpointer data) {
    (void)ev; (void)data;
    gtk_widget_hide(win);
    return TRUE;
}
static void setup_settings_hide_on_close(void *win) {
    g_signal_connect(GTK_WIDGET(win), "delete-event",
                     G_CALLBACK(on_settings_delete), NULL);
}
*/
import "C"
import "unsafe"

func showWebviewWindow(win unsafe.Pointer) {
	C.show_window(win)
}

func hideWebviewWindow(win unsafe.Pointer) {
	C.hide_window(win)
}

// interceptSettingsClose ensures the WM close button hides the window
// rather than destroying it, so it can be reopened.
func interceptSettingsClose(win unsafe.Pointer) {
	C.setup_settings_hide_on_close(win)
}
