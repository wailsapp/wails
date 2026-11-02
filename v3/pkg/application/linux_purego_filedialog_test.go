//go:build linux && purego && !gtk3 && !android && !server

package application

import "testing"

// The purego backend has to load and open file dialogs on GTK 4.8 (the floor,
// Debian 12): GtkFileDialog only exists from GTK 4.10, so those bindings are
// optional and the GtkFileChooser path replaces them when GTK lacks the class.
// These tests pin both halves of that contract without a display.

func TestFileDialogsLoadOnGtk48(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}

	// The version that decides which dialog implementation runs is the
	// presence of GtkFileDialog.
	hasFileDialog := _native_gtk_file_dialog_new != nil
	if got := dialogUsesGtkFileDialog(); got != hasFileDialog {
		t.Fatalf("dialogUsesGtkFileDialog() = %v, gtk_file_dialog_new bound = %v", got, hasFileDialog)
	}

	// GtkFileDialog shipped complete in 4.10, so its functions resolve
	// together: either all are bound (4.10+) or none are (4.8). A required
	// one would fail the whole backend on GTK 4.8, before this test could run.
	for name, bound := range map[string]bool{
		"gtk_file_dialog_new":                            _native_gtk_file_dialog_new != nil,
		"gtk_file_dialog_open":                           _native_gtk_file_dialog_open != nil,
		"gtk_file_dialog_open_finish":                    _native_gtk_file_dialog_open_finish != nil,
		"gtk_file_dialog_open_multiple":                  _native_gtk_file_dialog_open_multiple != nil,
		"gtk_file_dialog_open_multiple_finish":           _native_gtk_file_dialog_open_multiple_finish != nil,
		"gtk_file_dialog_save":                           _native_gtk_file_dialog_save != nil,
		"gtk_file_dialog_save_finish":                    _native_gtk_file_dialog_save_finish != nil,
		"gtk_file_dialog_select_folder":                  _native_gtk_file_dialog_select_folder != nil,
		"gtk_file_dialog_select_folder_finish":           _native_gtk_file_dialog_select_folder_finish != nil,
		"gtk_file_dialog_select_multiple_folders":        _native_gtk_file_dialog_select_multiple_folders != nil,
		"gtk_file_dialog_select_multiple_folders_finish": _native_gtk_file_dialog_select_multiple_folders_finish != nil,
		"gtk_file_dialog_set_accept_label":               _native_gtk_file_dialog_set_accept_label != nil,
		"gtk_file_dialog_set_filters":                    _native_gtk_file_dialog_set_filters != nil,
		"gtk_file_dialog_set_initial_folder":             _native_gtk_file_dialog_set_initial_folder != nil,
		"gtk_file_dialog_set_initial_name":               _native_gtk_file_dialog_set_initial_name != nil,
		"gtk_file_dialog_set_title":                      _native_gtk_file_dialog_set_title != nil,
	} {
		if bound != hasFileDialog {
			t.Errorf("%s bound = %v, GtkFileDialog available = %v", name, bound, hasFileDialog)
		}
	}

	// The GtkFileChooser replacements exist in every GTK4, so they are
	// required bindings: their absence would have failed linuxLibsErr above.
	for name, bound := range map[string]bool{
		"gtk_file_chooser_widget_new":          _native_gtk_file_chooser_widget_new != nil,
		"gtk_file_chooser_set_select_multiple": _native_gtk_file_chooser_set_select_multiple != nil,
		"gtk_file_chooser_set_create_folders":  _native_gtk_file_chooser_set_create_folders != nil,
		"gtk_file_chooser_set_current_folder":  _native_gtk_file_chooser_set_current_folder != nil,
		"gtk_file_chooser_set_current_name":    _native_gtk_file_chooser_set_current_name != nil,
		"gtk_file_chooser_add_filter":          _native_gtk_file_chooser_add_filter != nil,
		"gtk_file_chooser_get_files":           _native_gtk_file_chooser_get_files != nil,
		"gtk_button_new_with_mnemonic":         _native_gtk_button_new_with_mnemonic != nil,
	} {
		if !bound {
			t.Errorf("%s is not bound", name)
		}
	}
}

func TestFileChooserAcceptLabel(t *testing.T) {
	for _, tc := range []struct {
		label  string
		action int
		want   string
	}{
		{"", gtkFileChooserActionOpen, "_Open"},
		{"", gtkFileChooserActionSave, "_Save"},
		{"", gtkFileChooserActionSelectFolder, "_Select"},
		{"_Pick", gtkFileChooserActionOpen, "_Pick"},
		{"Pick", gtkFileChooserActionSave, "Pick"},
	} {
		if got := fileChooserAcceptLabel(tc.label, tc.action); got != tc.want {
			t.Errorf("fileChooserAcceptLabel(%q, %d) = %q, want %q", tc.label, tc.action, got, tc.want)
		}
	}
}

// TestFileDialogFallbackSelected forces the GTK 4.8 path by unbinding
// gtk_file_dialog_new and checks that the chooser is selected, and that the
// GtkFileDialog helpers refuse to call through the missing symbols. It needs
// no display: nothing is shown, only the selection logic runs.
func TestFileDialogFallbackSelected(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}

	saved := _native_gtk_file_dialog_new
	t.Cleanup(func() { _native_gtk_file_dialog_new = saved })

	_native_gtk_file_dialog_new = nil
	if dialogUsesGtkFileDialog() {
		t.Fatal("dialogUsesGtkFileDialog() = true with gtk_file_dialog_new unbound; want the chooser fallback")
	}
	if got := gtkFileDialogNew(); got != 0 {
		t.Fatalf("gtkFileDialogNew() = %#x with gtk_file_dialog_new unbound, want 0", got)
	}
}
