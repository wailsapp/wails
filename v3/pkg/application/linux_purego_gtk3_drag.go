//go:build linux && purego && gtk3 && !android && !server

package application

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	gtk3DragTargets     func(uintptr) uintptr
	gtk3AtomName        func(uintptr) uintptr
	gtk3DragFinish      func(uintptr, int32, int32, uint32)
	gtk3DragStatus      func(uintptr, uint32, uint32)
	gtk3DragGetData     func(uintptr, uintptr, uintptr, uint32)
	gtk3SelectionLength func(uintptr) int32
	gtk3SelectionData   func(uintptr) uintptr
	gtk3ExtractURIs     func(uintptr) uintptr
	gtk3FilenameFromURI func(uintptr, uintptr, uintptr) uintptr
	gtk3FreeStrings     func(uintptr)
	gtk3DragEntered     bool
)

func loadGTK3DragFunctions(lib uintptr, missing *[]string) {
	for name, target := range map[string]any{
		"gdk_drag_context_list_targets": &gtk3DragTargets, "gdk_atom_name": &gtk3AtomName,
		"gtk_drag_finish": &gtk3DragFinish, "gdk_drag_status": &gtk3DragStatus,
		"gtk_drag_get_data":             &gtk3DragGetData,
		"gtk_selection_data_get_length": &gtk3SelectionLength, "gtk_selection_data_get_data": &gtk3SelectionData,
		"g_uri_list_extract_uris": &gtk3ExtractURIs, "g_filename_from_uri": &gtk3FilenameFromURI,
		"g_strfreev": &gtk3FreeStrings,
	} {
		bindGTK3Symbol(missing, lib, target, name)
	}
}

func gtk3IsFileDrag(context uintptr) bool {
	for list := (*_Ctype_GList)(unsafe.Pointer(gtk3DragTargets(context))); list != nil; list = list.next {
		name := gtk3AtomName(uintptr(list.data))
		match := _Cfunc_GoString((*_Ctype_char)(unsafe.Pointer(name))) == "text/uri-list"
		_Cfunc_free(unsafe.Pointer(name))
		if match {
			return true
		}
	}
	return false
}

var gtk3DragDrop = purego.NewCallback(func(widget, context uintptr, x, y int32, timestamp uint32, data uintptr) int32 {
	if !gtk3IsFileDrag(context) {
		return 0
	}
	gtk3DragEntered = false
	gtk3SetData(unsafe.Pointer(widget), "drop-x", uintptr(x))
	gtk3SetData(unsafe.Pointer(widget), "drop-y", uintptr(y))
	gtk3SetData(unsafe.Pointer(widget), "drop-data", data)
	gtk3DragGetData(widget, context, gtk3Atom("text/uri-list", 0), timestamp)
	return 1
})

var gtk3DragMotion = purego.NewCallback(func(widget, context uintptr, x, y int32, timestamp uint32, data uintptr) int32 {
	if !gtk3IsFileDrag(context) {
		return 0
	}
	gtk3DragStatus(context, 2, timestamp) // GDK_ACTION_COPY
	if !gtk3DragEntered {
		gtk3DragEntered = true
		onDragEnter(data)
	}
	onDragOver(_Ctype_gint(x), _Ctype_gint(y), data)
	return 1
})

var gtk3DragLeave = purego.NewCallback(func(widget, context uintptr, timestamp uint32, data uintptr) {
	if gtk3IsFileDrag(context) && gtk3DragEntered {
		gtk3DragEntered = false
		onDragLeave(data)
	}
})

var gtk3DragReceived = purego.NewCallback(func(widget, context uintptr, x, y int32, selection uintptr, target, timestamp uint32, data uintptr) {
	if target != 2 {
		return
	}
	if selection == 0 || gtk3SelectionLength(selection) <= 0 {
		gtk3DragFinish(context, 0, 0, timestamp)
		return
	}
	uris := gtk3ExtractURIs(gtk3SelectionData(selection))
	if uris == 0 {
		gtk3DragFinish(context, 0, 0, timestamp)
		return
	}
	defer gtk3FreeStrings(uris)
	var filenames []*_Ctype_char
	for index := uintptr(0); ; index++ {
		uri := *(*uintptr)(unsafe.Pointer(uris + index*unsafe.Sizeof(uintptr(0))))
		if uri == 0 {
			break
		}
		filename := gtk3FilenameFromURI(uri, 0, 0)
		if filename != 0 {
			filenames = append(filenames, (*_Ctype_char)(unsafe.Pointer(filename)))
		}
	}
	if len(filenames) == 0 {
		gtk3DragFinish(context, 0, 0, timestamp)
		return
	}
	defer func() {
		for _, filename := range filenames {
			_Cfunc_free(unsafe.Pointer(filename))
		}
	}()
	x = int32(gtk3GetData(unsafe.Pointer(widget), "drop-x"))
	y = int32(gtk3GetData(unsafe.Pointer(widget), "drop-y"))
	data = gtk3GetData(unsafe.Pointer(widget), "drop-data")
	filenames = append(filenames, nil)
	onUriList(&filenames[0], _Ctype_gint(x), _Ctype_gint(y), data)
	gtk3DragFinish(context, 1, 0, timestamp)
})

var gtk3DragDropBlocked = purego.NewCallback(func(widget, context uintptr, x, y int32, timestamp uint32, data uintptr) int32 {
	if !gtk3IsFileDrag(context) {
		return 0
	}
	gtk3DragFinish(context, 0, 0, timestamp)
	return 1
})
var gtk3DragMotionBlocked = purego.NewCallback(func(widget, context uintptr, x, y int32, timestamp uint32, data uintptr) int32 {
	if !gtk3IsFileDrag(context) {
		return 0
	}
	gtk3DragStatus(context, 0, timestamp)
	return 1
})

func gtk3ConnectSignal(widget *_Ctype_GtkWidget, name string, callback, data uintptr) {
	cname := _Cfunc_CString(name)
	defer _Cfunc_free(unsafe.Pointer(cname))
	gtk3Connect(unsafe.Pointer(widget), cname, unsafe.Pointer(callback), data, 0, 0)
}

func _Cfunc_enableDND(widget *_Ctype_GtkWidget, data _Ctype_uintptr_t) {
	gtk3ConnectSignal(widget, "drag-data-received", gtk3DragReceived, uintptr(data))
	gtk3ConnectSignal(widget, "drag-drop", gtk3DragDrop, uintptr(data))
	gtk3ConnectSignal(widget, "drag-motion", gtk3DragMotion, uintptr(data))
	gtk3ConnectSignal(widget, "drag-leave", gtk3DragLeave, uintptr(data))
}

func _Cfunc_disableDND(widget *_Ctype_GtkWidget, data _Ctype_uintptr_t) {
	gtk3ConnectSignal(widget, "drag-drop", gtk3DragDropBlocked, uintptr(data))
	gtk3ConnectSignal(widget, "drag-motion", gtk3DragMotionBlocked, uintptr(data))
}
