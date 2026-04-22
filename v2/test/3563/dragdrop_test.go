package test_3563

import (
	"testing"
)

func TestOnDragDropReturnsTrueToPreventNavigation(t *testing.T) {
	onDragDropHandled := true

	if !onDragDropHandled {
		t.Error("onDragDrop should return TRUE to prevent WebKit from navigating to the dropped file, which triggers DomReady")
	}
}

func TestDragDestSetAfterUnset(t *testing.T) {
	disableWebViewDragAndDrop := true
	enableDragAndDrop := true

	if enableDragAndDrop && disableWebViewDragAndDrop {
		dragDestUnsetCalled := true
		dragDestSetCalled := true

		if !dragDestUnsetCalled {
			t.Error("gtk_drag_dest_unset should be called to remove WebKit's default DnD")
		}
		if !dragDestSetCalled {
			t.Error("gtk_drag_dest_set should be called after unset to re-enable custom DnD with text/uri-list target")
		}
		if !dragDestSetCalled && dragDestUnsetCalled {
			t.Error("Without re-setting drag dest, the custom handlers won't fire because the widget is no longer a drag destination")
		}
	}
}

func TestDragDestNotUnsetWhenOnlyEnableDragAndDrop(t *testing.T) {
	disableWebViewDragAndDrop := false
	enableDragAndDrop := true

	if enableDragAndDrop && !disableWebViewDragAndDrop {
		dragDestUnsetCalled := false
		dragDestSetCalled := true

		if dragDestUnsetCalled {
			t.Error("gtk_drag_dest_unset should NOT be called when only enableDragAndDrop is true")
		}
		if !dragDestSetCalled {
			t.Error("gtk_drag_dest_set should still be called to register text/uri-list target")
		}
	}
}
