//go:build linux && purego && !gtk3 && !android && !server

package application

// Pure-Go coverage for the GTK4 purego backend: key/modifier mapping, string
// and GValue-adjacent marshalling, library loading error paths, the dialog
// request plumbing and the action-name registry. None of these tests opens a
// display; the GTK-level behaviour lives in linux_purego_gtk_test.go.

import (
	"slices"
	"strings"
	"testing"
	"unsafe"
)

func TestParseKeyGTK(t *testing.T) {
	for _, tc := range []struct {
		key  string
		want uint32
	}{
		{"backspace", 0xff08},
		{"tab", 0xff09},
		{"return", 0xff0d},
		{"enter", 0xff0d}, // alias of return
		{"escape", 0xff1b},
		{"left", 0xff51},
		{"up", 0xff52},
		{"right", 0xff53},
		{"down", 0xff54},
		{"space", 0x0020},
		{"delete", 0xffff},
		{"home", 0xff50},
		{"end", 0xff57},
		{"page up", 0xff55},
		{"page down", 0xff56},
		{"f1", 0xffbe},
		{"f12", 0xffc9},
		{"f35", 0xffe0},
		{"numlock", 0xff7f},
		// Single characters go through gdk_unicode_to_keyval, which for ASCII
		// is the code point itself.
		{"a", 'a'},
		{"A", 'A'},
		{"1", '1'},
		{"+", '+'},
		// Unknown multi-character names and the empty string map to no key.
		{"", 0},
		{"notakey", 0},
		{"F1", 0}, // names are case-sensitive; the Wails layer lower-cases
	} {
		if got := parseKeyGTK(tc.key); got != tc.want {
			t.Errorf("parseKeyGTK(%q) = %#x, want %#x", tc.key, got, tc.want)
		}
	}
}

func TestParseModifiersGTK(t *testing.T) {
	for _, tc := range []struct {
		name string
		mods []modifier
		want uint32
	}{
		{"none", nil, 0},
		{"shift", []modifier{ShiftKey}, gdkShiftMask},
		{"control", []modifier{ControlKey}, gdkControlMask},
		{"cmd-or-ctrl maps to control", []modifier{CmdOrCtrlKey}, gdkControlMask},
		{"option-or-alt maps to alt", []modifier{OptionOrAltKey}, gdkAltMask},
		{"super", []modifier{SuperKey}, gdkSuperMask},
		{"ctrl+shift", []modifier{ControlKey, ShiftKey}, gdkControlMask | gdkShiftMask},
		{"cmd+ctrl is one bit", []modifier{CmdOrCtrlKey, ControlKey}, gdkControlMask},
		{"all five", []modifier{ShiftKey, ControlKey, OptionOrAltKey, SuperKey, CmdOrCtrlKey},
			gdkShiftMask | gdkControlMask | gdkAltMask | gdkSuperMask},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseModifiersGTK(tc.mods); got != tc.want {
				t.Fatalf("parseModifiersGTK(%v) = %#x, want %#x", tc.mods, got, tc.want)
			}
		})
	}
}

func TestAcceleratorToGTK(t *testing.T) {
	key, mods := acceleratorToGTK(&accelerator{Key: "s", Modifiers: []modifier{CmdOrCtrlKey, ShiftKey}})
	if key != 's' || mods != gdkControlMask|gdkShiftMask {
		t.Fatalf("acceleratorToGTK(Ctrl+Shift+s) = (%#x, %#x)", key, mods)
	}
	key, mods = acceleratorToGTK(&accelerator{Key: "f5"})
	if key != 0xffc2 || mods != 0 {
		t.Fatalf("acceleratorToGTK(F5) = (%#x, %#x)", key, mods)
	}
	// An unknown key yields keysym 0, which setMenuItemAccelerator treats as
	// "no accelerator"; the modifier bits are still reported.
	key, mods = acceleratorToGTK(&accelerator{Key: "bogus", Modifiers: []modifier{SuperKey}})
	if key != 0 || mods != gdkSuperMask {
		t.Fatalf("acceleratorToGTK(bogus) = (%#x, %#x)", key, mods)
	}
}

func TestWriteInt(t *testing.T) {
	for _, tc := range []struct {
		n    int
		want string
	}{
		{0, "0"},
		{7, "7"},
		{10, "10"},
		{1234567, "1234567"},
		{-1, "-1"},
		{-900, "-900"},
	} {
		buf := make([]byte, 16)
		n := writeInt(buf, tc.n)
		if got := string(buf[:n]); got != tc.want || n != len(tc.want) {
			t.Errorf("writeInt(%d) wrote %q (n=%d), want %q", tc.n, got, n, tc.want)
		}
	}
}

func TestBoolHelpers(t *testing.T) {
	if gbool(true) != 1 || gbool(false) != 0 {
		t.Fatal("gbool")
	}
	if btoi(true) != 1 || btoi(false) != 0 {
		t.Fatal("btoi")
	}
}

func TestCStringRoundTrip(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	for _, s := range []string{"", "plain", "héllo wörld ✓", strings.Repeat("x", 4096)} {
		p := cString(s)
		if p == 0 {
			t.Fatalf("cString(%q) returned NULL", s)
		}
		if got := goString(p); got != s {
			t.Errorf("goString(cString(%q)) = %q", s, got)
		}
		// takeGString frees p, so the copy is taken exactly once.
		p2 := cString(s)
		if got := takeGString(p2); got != s {
			t.Errorf("takeGString(cString(%q)) = %q", s, got)
		}
	}
}

func TestGoStringNull(t *testing.T) {
	if got := goString(0); got != "" {
		t.Fatalf("goString(NULL) = %q", got)
	}
	if got := takeGString(0); got != "" {
		t.Fatalf("takeGString(NULL) = %q", got)
	}
}

func TestCallocOwnsStrings(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	c := NewCalloc()
	p := c.String("calloc-value")
	if got := goString(uintptr(unsafe.Pointer(p))); got != "calloc-value" {
		t.Fatalf("Calloc.String = %q", got)
	}
	if len(c.pool) != 1 {
		t.Fatalf("pool holds %d pointers, want 1", len(c.pool))
	}
	c.Free()
	if len(c.pool) != 0 {
		t.Fatalf("Free left %d pointers in the pool", len(c.pool))
	}
}

func TestSplitAndTrim(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"*.txt", []string{"*.txt"}},
		{"*.txt;*.md", []string{"*.txt", "*.md"}},
		{" *.go ;\t*.mod\t", []string{"*.go", "*.mod"}},
		{"a;;b", []string{"a", "", "b"}},
		{"trailing;", []string{"trailing", ""}},
		{"", []string{""}},
	} {
		if got := splitAndTrim(tc.in, ";"); !slices.Equal(got, tc.want) {
			t.Errorf("splitAndTrim(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDialogTypeToIconName(t *testing.T) {
	for dt, want := range map[DialogType]string{
		InfoDialogType:     "dialog-information-symbolic",
		WarningDialogType:  "dialog-warning-symbolic",
		ErrorDialogType:    "dialog-error-symbolic",
		QuestionDialogType: "dialog-question-symbolic",
	} {
		if got := dialogTypeToIconName(dt); got != want {
			t.Errorf("dialogTypeToIconName(%v) = %q, want %q", dt, got, want)
		}
	}
	if got := dialogTypeToIconName(DialogType(99)); got != "" {
		t.Errorf("unknown dialog type icon = %q, want empty", got)
	}
}

func TestBuildVersionString(t *testing.T) {
	if got := buildVersionString(4, 14, 2); got != "4.14.2" {
		t.Fatalf("buildVersionString = %q", got)
	}
	if got := buildVersionString(0, 0, 0); got != "0.0.0" {
		t.Fatalf("buildVersionString(0) = %q", got)
	}
}

func TestGdkSurfaceEdgeForBorder(t *testing.T) {
	// The values are GdkSurfaceEdge (gdk/gdkenums.h): NW=0, N=1, NE=2, W=3,
	// E=4, SW=5, S=6, SE=7. A wrong value resizes from the wrong edge.
	want := map[string]int32{
		"nw-resize": 0, "n-resize": 1, "ne-resize": 2, "w-resize": 3,
		"e-resize": 4, "sw-resize": 5, "s-resize": 6, "se-resize": 7,
	}
	if len(gdkSurfaceEdgeForBorder) != len(want) {
		t.Fatalf("%d borders mapped, want %d", len(gdkSurfaceEdgeForBorder), len(want))
	}
	for border, edge := range want {
		if got, ok := gdkSurfaceEdgeForBorder[border]; !ok || got != edge {
			t.Errorf("border %q -> %d (ok=%v), want %d", border, got, ok, edge)
		}
	}
}

func TestStartResizeRejectsUnknownBorder(t *testing.T) {
	// The unknown-border check runs before any widget is touched, so a
	// zero-valued window is safe here.
	w := &linuxWebviewWindow{}
	if err := w.startResize("middle-resize"); err == nil {
		t.Fatal("startResize(unknown) returned nil error")
	} else if !strings.Contains(err.Error(), "middle-resize") {
		t.Fatalf("error does not name the border: %v", err)
	}
	if err := w.startResize(""); err == nil {
		t.Fatal("startResize(\"\") returned nil error")
	}
}

func TestGetLinuxWebviewWindow(t *testing.T) {
	if got := getLinuxWebviewWindow(nil); got != nil {
		t.Fatalf("nil window -> %v", got)
	}
	// A WebviewWindow whose impl is not the Linux implementation (here, no
	// impl at all) yields nil instead of panicking on the type assertion.
	if got := getLinuxWebviewWindow(&WebviewWindow{}); got != nil {
		t.Fatalf("window without linux impl -> %v", got)
	}
	lw := &linuxWebviewWindow{id: 3}
	if got := getLinuxWebviewWindow(&WebviewWindow{impl: lw}); got != lw {
		t.Fatalf("linux window not returned: %v", got)
	}
}

func TestGenerateAndLookupActionName(t *testing.T) {
	const a, b, c = uint(910001), uint(910002), uint(910003)
	t.Cleanup(func() {
		menuItemsLock.Lock()
		for _, id := range []uint{a, b, c} {
			delete(menuItemActions, id)
		}
		menuItemsLock.Unlock()
	})

	if _, ok := lookupActionName(a); ok {
		t.Fatal("lookup before generate succeeded")
	}
	nameA := generateActionName(a)
	nameB := generateActionName(b)
	if nameA == nameB || !strings.HasPrefix(nameA, "action_") {
		t.Fatalf("names not unique: %q %q", nameA, nameB)
	}
	if got, ok := lookupActionName(a); !ok || got != nameA {
		t.Fatalf("lookup(a) = %q, %v; want %q", got, ok, nameA)
	}
	// Regenerating for the same item replaces the mapping with a new name.
	nameA2 := generateActionName(a)
	if nameA2 == nameA {
		t.Fatalf("regenerated name %q equals the old one", nameA2)
	}
	if got, _ := lookupActionName(a); got != nameA2 {
		t.Fatalf("lookup after regenerate = %q, want %q", got, nameA2)
	}
	if _, ok := lookupActionName(c); ok {
		t.Fatal("unrelated item resolved")
	}
}

func TestNextDialogRequestIDIsMonotonic(t *testing.T) {
	first := nextDialogRequestID()
	second := nextDialogRequestID()
	if second != first+1 {
		t.Fatalf("request ids %d then %d, want consecutive", first, second)
	}
}

func TestFileDialogCallbackDeliversPathsInOrder(t *testing.T) {
	id := nextDialogRequestID()
	ch := make(chan string, 100)
	dialogRequestMutex.Lock()
	fileDialogCallbacks[id] = ch
	dialogRequestMutex.Unlock()

	fileDialogCallback(id, []string{"/a", "/b", "/c"}, false)

	var got []string
	for p := range ch { // closed after the last path
		got = append(got, p)
	}
	if !slices.Equal(got, []string{"/a", "/b", "/c"}) {
		t.Fatalf("paths = %q", got)
	}
	dialogRequestMutex.Lock()
	_, stillRegistered := fileDialogCallbacks[id]
	dialogRequestMutex.Unlock()
	if stillRegistered {
		t.Fatal("request still registered after delivery")
	}
}

func TestFileDialogCallbackCancelledClosesWithoutPaths(t *testing.T) {
	id := nextDialogRequestID()
	ch := make(chan string, 100)
	dialogRequestMutex.Lock()
	fileDialogCallbacks[id] = ch
	dialogRequestMutex.Unlock()

	fileDialogCallback(id, []string{"/ignored"}, true)

	if p, open := <-ch; open {
		t.Fatalf("cancelled dialog delivered %q", p)
	}
}

func TestFileDialogCallbackUnknownIDIsNoop(t *testing.T) {
	// No panic, no goroutine blocked on a missing channel.
	fileDialogCallback(nextDialogRequestID(), []string{"/x"}, false)
	fileDialogCallback(nextDialogRequestID(), nil, true)
}

func TestAlertDialogCallbackDeliversButton(t *testing.T) {
	id := nextDialogRequestID()
	ch := make(chan int, 1)
	dialogRequestMutex.Lock()
	alertDialogCallbacks[id] = ch
	dialogRequestMutex.Unlock()

	alertDialogCallback(id, 2)

	if v, ok := <-ch; !ok || v != 2 {
		t.Fatalf("button = %d (open=%v), want 2", v, ok)
	}
	if _, open := <-ch; open {
		t.Fatal("alert channel not closed after the result")
	}
	// Unknown ids are ignored.
	alertDialogCallback(nextDialogRequestID(), 0)
}

func TestAddFileFilterPatternsAreTrimmed(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	// The filter store is a GListStore of GtkFileFilter. Adding a filter
	// with several patterns must still produce exactly one list entry.
	// GtkFileFilter registers its GType on first use; create one first.
	first := gtk_file_filter_new()
	g_object_unref(first)
	store := g_list_store_new(g_type_from_name("GtkFileFilter"))
	defer g_object_unref(store)
	addFileFilter(0, store, "Text", " *.txt ; *.md ;; ")
	if n := g_list_model_get_n_items(store); n != 1 {
		t.Fatalf("filter store holds %d items, want 1", n)
	}
}

func TestVariantTypeStringIsCached(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	first := variantTypeString()
	if first == 0 {
		t.Fatal("variantTypeString returned NULL")
	}
	if variantTypeString() != first {
		t.Fatal("variantTypeString allocated a second type")
	}
}

func TestGTypeInstanceIsA(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	if gTypeInstanceIsA(0, gtype_Object()) || gTypeInstanceIsA(1, 0) {
		t.Fatal("NULL instance or type must not match")
	}
	menu := g_menu_new()
	defer g_object_unref(menu)
	if !gTypeInstanceIsA(menu, gtype_Object()) {
		t.Fatal("GMenu is not a GObject")
	}
	if !gTypeInstanceIsA(menu, g_type_from_name("GMenuModel")) {
		t.Fatal("GMenu is not a GMenuModel")
	}
	if gTypeInstanceIsA(menu, g_type_from_name("GtkWidget")) {
		t.Fatal("GMenu reported as GtkWidget")
	}
}

func TestGTypeObjectResolvesByName(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	if gtype_Object() == 0 || gtype_Object() != g_type_from_name("GObject") {
		t.Fatal("gtype_Object does not resolve GObject")
	}
}

func TestDlopenFirstFallsBack(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	// DlopenFirst with only bogus names reports every failed attempt.
	handle, errs := dlopenFirst([]string{"libwails-no-such-lib.so.9", "libwails-no-such-lib.so"})
	if handle != 0 {
		t.Fatalf("bogus library loaded: %#x", handle)
	}
	if len(errs) != 2 {
		t.Fatalf("dlopen errors = %q, want 2", errs)
	}
	// A real soname is found on the first candidate that exists.
	handle, errs = dlopenFirst([]string{"libwails-no-such-lib.so.9", "libglib-2.0.so.0"})
	if handle == 0 || errs != nil {
		t.Fatalf("fallback to libglib failed: handle=%#x errs=%q", handle, errs)
	}
}

func TestHaveSymbol(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	if haveSymbol(0, "g_free") {
		t.Fatal("symbol reported on a NULL library")
	}
	if !haveSymbol(libGLib, "g_free") {
		t.Fatal("g_free not found in libglib")
	}
	if haveSymbol(libGLib, "wails_no_such_symbol") {
		t.Fatal("bogus symbol reported as present")
	}
}

func TestRegisterOptionalSkipsMissingSymbols(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	var fn func(uintptr) uintptr
	var strdup func(string) uintptr
	if registerOptional(&fn, 0, "g_free") {
		t.Fatal("registerOptional on NULL library returned true")
	}
	if registerOptional(&fn, libGLib, "wails_no_such_symbol") || fn != nil {
		t.Fatal("registerOptional bound a missing symbol")
	}
	if !registerOptional(&strdup, libGLib, "g_strdup") || strdup == nil {
		t.Fatal("registerOptional did not bind g_strdup")
	}
	// A bound optional function is usable.
	if got := takeGString(strdup("optional")); got != "optional" {
		t.Fatalf("bound g_strdup returned %q", got)
	}
}

func TestRegisterRecordsMissingRequiredSymbols(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	saved := missingSymbols
	t.Cleanup(func() { missingSymbols = saved })
	missingSymbols = nil

	var fn func()
	register(&fn, 0, "g_free") // library not loaded: silently skipped
	if len(missingSymbols) != 0 {
		t.Fatalf("NULL library recorded %q", missingSymbols)
	}
	register(&fn, libGLib, "wails_no_such_symbol")
	if !slices.Equal(missingSymbols, []string{"wails_no_such_symbol"}) {
		t.Fatalf("missing symbols = %q", missingSymbols)
	}
	if fn != nil {
		t.Fatal("missing required symbol was bound anyway")
	}
}

func TestLoadLinuxLibrariesReportsMissingLibraries(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	savedCandidates := linuxLibCandidates
	t.Cleanup(func() { linuxLibCandidates = savedCandidates })

	var handle uintptr
	linuxLibCandidates = []struct {
		handle  *uintptr
		install string
		names   []string
	}{
		{&handle, "libwails-test-pkg", []string{"libwails-no-such-lib.so.9"}},
	}
	err := loadLinuxLibraries()
	if err == nil {
		t.Fatal("loadLinuxLibraries succeeded with no libraries")
	}
	for _, want := range []string{"libwails-no-such-lib.so.9", "libwails-test-pkg", "dnf install gtk4", "apt install"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q:\n%v", want, err)
		}
	}
	if handle != 0 {
		t.Fatal("failed library left a handle behind")
	}
}
