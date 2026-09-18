package application

import (
	"errors"
	"testing"
)

func TestRestorationStateRoundTrip(t *testing.T) {
	state := RestorationState{Data: map[string]string{
		"counter":  "42",
		"document": "/tmp/notes.md",
		"unicode":  "héllo \"quoted\" \\ back",
		"empty":    "",
	}}
	encoded, err := encodeRestorationState(state)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	decoded, err := decodeRestorationState(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(decoded.Data) != len(state.Data) {
		t.Fatalf("decoded %d keys, want %d", len(decoded.Data), len(state.Data))
	}
	for key, want := range state.Data {
		if got := decoded.Get(key); got != want {
			t.Fatalf("decoded[%q] = %q, want %q", key, got, want)
		}
	}
	if decoded.Get("missing") != "" {
		t.Fatal("Get on a missing key must return an empty string")
	}
}

func TestRestorationStateEncodesEmptyAndNil(t *testing.T) {
	for _, state := range []RestorationState{{}, {Data: map[string]string{}}} {
		encoded, err := encodeRestorationState(state)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		if encoded != "{}" {
			t.Fatalf("encoded empty state = %q, want {}", encoded)
		}
	}
	decoded, err := decodeRestorationState("")
	if err != nil || decoded.Data == nil || len(decoded.Data) != 0 {
		t.Fatalf("decode of an empty string = %+v, %v", decoded, err)
	}
	decoded, err = decodeRestorationState("null")
	if err != nil || decoded.Data == nil {
		t.Fatalf("decode of null = %+v, %v", decoded, err)
	}
	if (RestorationState{}).Get("x") != "" {
		t.Fatal("Get on a nil map must return an empty string")
	}
}

func TestRestorationStateRejectsMalformedData(t *testing.T) {
	for _, input := range []string{"{", "[1,2]", `{"a":1}`, "not json"} {
		decoded, err := decodeRestorationState(input)
		if err == nil {
			t.Fatalf("decode(%q) succeeded", input)
		}
		if decoded.Data == nil {
			t.Fatal("a failed decode must still return a usable state")
		}
	}
}

func TestRestorationSettersQueueBeforeCreation(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "pending", Mac: MacWindow{RestorationID: "from-options"}}}
	defer forgetMacRestoration(window)
	if got := window.RestorationID(); got != "from-options" {
		t.Fatalf("RestorationID from options = %q", got)
	}
	window.SetRestorationID("explicit")
	if got := window.RestorationID(); got != "explicit" {
		t.Fatalf("RestorationID after SetRestorationID = %q", got)
	}
	data := map[string]string{"counter": "3"}
	window.SetRestorationData(data)
	data["counter"] = "mutated"
	if got := window.RestorationData(); got["counter"] != "3" {
		t.Fatalf("RestorationData did not copy the map: %v", got)
	}
	window.SetRestorationData(nil)
	if window.RestorationData() != nil {
		t.Fatal("RestorationData after clearing must be nil")
	}
	window.SetRestorationID("")
	if window.RestorationID() != "" {
		t.Fatal("an empty SetRestorationID must override the option")
	}
}

func TestRestorationNilAndUncreatedWindowsAreSafe(t *testing.T) {
	var nilWindow *WebviewWindow
	if got, _ := nilWindow.SetRestorationID("x").(*WebviewWindow); got != nil {
		t.Fatal("SetRestorationID on a nil window must return the nil receiver")
	}
	if got, _ := nilWindow.SetRestorationData(map[string]string{"a": "b"}).(*WebviewWindow); got != nil {
		t.Fatal("SetRestorationData on a nil window must return the nil receiver")
	}
	if nilWindow.RestorationID() != "" || nilWindow.RestorationData() != nil {
		t.Fatal("nil window reported restoration state")
	}
	if _, err := nilWindow.InteractionState(); err == nil {
		t.Fatal("InteractionState on a nil window succeeded")
	}
	if err := nilWindow.RestoreInteractionState([]byte{1}); err == nil {
		t.Fatal("RestoreInteractionState on a nil window succeeded")
	}
	uncreated := &WebviewWindow{options: WebviewWindowOptions{Name: "uncreated"}}
	defer forgetMacRestoration(uncreated)
	_, err := uncreated.InteractionState()
	if macWindowRestorationSupported() {
		if !errors.Is(err, ErrMacWindowNotCreated) {
			t.Fatalf("InteractionState before creation = %v", err)
		}
		if err := uncreated.RestoreInteractionState([]byte{1}); !errors.Is(err, ErrMacWindowNotCreated) {
			t.Fatalf("RestoreInteractionState before creation = %v", err)
		}
	} else {
		if !errors.Is(err, ErrMacOnly) {
			t.Fatalf("InteractionState off macOS = %v", err)
		}
		if err := uncreated.RestoreInteractionState([]byte{1}); !errors.Is(err, ErrMacOnly) {
			t.Fatalf("RestoreInteractionState off macOS = %v", err)
		}
	}
}

func TestRestoreHandlerRegistrationAndPanicRecovery(t *testing.T) {
	defer setMacWindowRestoreHandler(nil)
	if getMacWindowRestoreHandler() != nil {
		t.Fatal("a handler was registered before the test")
	}
	manager := &WindowManager{}
	var seenID string
	var seenState RestorationState
	manager.OnRestore(func(id string, state RestorationState) Window {
		seenID, seenState = id, state
		return nil
	})
	handler := getMacWindowRestoreHandler()
	if handler == nil {
		t.Fatal("OnRestore did not register the handler")
	}
	if window := runMacWindowRestoreHandler(handler, "main", RestorationState{Data: map[string]string{"k": "v"}}); window != nil {
		t.Fatal("handler returning nil must yield nil")
	}
	if seenID != "main" || seenState.Get("k") != "v" {
		t.Fatalf("handler received %q %v", seenID, seenState)
	}
	panicking := func(string, RestorationState) Window { panic("boom") }
	if window := runMacWindowRestoreHandler(panicking, "main", RestorationState{}); window != nil {
		t.Fatal("a panicking handler must yield nil")
	}
	// Without a handler the request completes with an error and must not
	// block or panic.
	setMacWindowRestoreHandler(nil)
	handleMacWindowRestore(macWindowRestoreRequest{requestID: 1, id: "main", encoded: "{}"})
}
