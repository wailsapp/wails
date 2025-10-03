//go:build linux
// +build linux

package linux

/*
#cgo webkit_6 CFLAGS: -DWEBKIT_6
#cgo !webkit_6 pkg-config: gtk+-3.0
#cgo webkit_6 pkg-config: gtk4

#include "gtk/gtk.h"

#ifdef WEBKIT_6
const GdkSurfaceEdge N_RESIZE = GDK_SURFACE_EDGE_NORTH;
const GdkSurfaceEdge NE_RESIZE = GDK_SURFACE_EDGE_NORTH_EAST;
const GdkSurfaceEdge E_RESIZE = GDK_SURFACE_EDGE_EAST;
const GdkSurfaceEdge SE_RESIZE = GDK_SURFACE_EDGE_SOUTH_EAST;
const GdkSurfaceEdge S_RESIZE = GDK_SURFACE_EDGE_SOUTH;
const GdkSurfaceEdge SW_RESIZE = GDK_SURFACE_EDGE_SOUTH_WEST;
const GdkSurfaceEdge W_RESIZE = GDK_SURFACE_EDGE_WEST;
const GdkSurfaceEdge NW_RESIZE = GDK_SURFACE_EDGE_NORTH_WEST;
#else
const GdkWindowEdge N_RESIZE = GDK_WINDOW_EDGE_NORTH;
const GdkWindowEdge NE_RESIZE = GDK_WINDOW_EDGE_NORTH_EAST;
const GdkWindowEdge E_RESIZE = GDK_WINDOW_EDGE_EAST;
const GdkWindowEdge SE_RESIZE = GDK_WINDOW_EDGE_SOUTH_EAST;
const GdkWindowEdge S_RESIZE = GDK_WINDOW_EDGE_SOUTH;
const GdkWindowEdge SW_RESIZE = GDK_WINDOW_EDGE_SOUTH_WEST;
const GdkWindowEdge W_RESIZE = GDK_WINDOW_EDGE_WEST;
const GdkWindowEdge NW_RESIZE = GDK_WINDOW_EDGE_NORTH_WEST;
#endif

void runMainLoop() {
	#ifdef WEBKIT_6
	GMainLoop *mainLoop = g_main_loop_new(NULL, true);
	setMainLoop(mainLoop);
	g_main_loop_run(mainLoop);
	#else
	gtk_main();
	#endif
}

gboolean initGtk() {
	#ifdef WEBKIT_6
	return gtk_init_check();
	#else
	return gtk_init_check(NULL, NULL);
	#endif
}

*/
import "C"

var edgeMap = map[string]uintptr{
	"n-resize":  C.N_RESIZE,
	"ne-resize": C.NE_RESIZE,
	"e-resize":  C.E_RESIZE,
	"se-resize": C.SE_RESIZE,
	"s-resize":  C.S_RESIZE,
	"sw-resize": C.SW_RESIZE,
	"w-resize":  C.W_RESIZE,
	"nw-resize": C.NW_RESIZE,
}

func runMainLoop() {
	C.runMainLoop()
}

func initGtk() C.gboolean {
	return C.initGtk()
}
