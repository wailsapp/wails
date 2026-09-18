package application

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestValidateFourCharCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "core suite", code: "aevt"},
		{name: "open documents", code: "odoc"},
		{name: "get url", code: "GURL"},
		{name: "with spaces", code: "gs  "},
		{name: "empty", code: "", wantErr: true},
		{name: "too short", code: "abc", wantErr: true},
		{name: "too long", code: "abcde", wantErr: true},
		{name: "control character", code: "ab\x01c", wantErr: true},
		{name: "non ascii", code: "ab\xc3\xa9", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateFourCharCode(test.code)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateFourCharCode(%q) error = %v, wantErr %v", test.code, err, test.wantErr)
			}
		})
	}
}

func TestAppleEventsRegister(t *testing.T) {
	m := newAppleEventsManager(nil)
	noop := func(*Context, AppleEvent) (AppleEventReply, error) { return AppleEventReply{}, nil }

	if err := m.register("WAIL", "note", nil); err == nil {
		t.Fatal("expected an error for a nil handler")
	}
	if err := m.register("WAI", "note", noop); err == nil {
		t.Fatal("expected an error for a bad class code")
	}
	if err := m.register("WAIL", "notes", noop); err == nil {
		t.Fatal("expected an error for a bad ID code")
	}
	if err := m.register("WAIL", "note", noop); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := m.register("WAIL", "note", noop); err != nil {
		t.Fatalf("re-register: %v", err)
	}
	if err := m.register("aevt", "odoc", noop); err != nil {
		t.Fatalf("register: %v", err)
	}
	if len(m.order) != 2 {
		t.Fatalf("expected 2 registrations in order, got %d", len(m.order))
	}
	if m.handler("WAIL", "note") == nil || m.handler("aevt", "odoc") == nil {
		t.Fatal("handlers not stored")
	}
	if m.handler("WAIL", "find") != nil {
		t.Fatal("unexpected handler for an unregistered event")
	}
}

func TestScriptingDefinition(t *testing.T) {
	m := newAppleEventsManager(nil)
	noop := func(*Context, AppleEvent) (AppleEventReply, error) { return AppleEventReply{}, nil }

	empty := m.ScriptingDefinition()
	if !strings.Contains(empty, `<dictionary title="Application Terminology">`) || strings.Contains(empty, "<suite") {
		t.Fatalf("unexpected empty sdef:\n%s", empty)
	}

	for _, reg := range [][2]string{{"WAIL", "note"}, {"WAIL", "Find"}, {"aevt", "odoc"}, {"WAIL", "<&>\""}} {
		if err := m.register(reg[0], reg[1], noop); err != nil {
			t.Fatalf("register %v: %v", reg, err)
		}
	}
	sdef := m.ScriptingDefinition()

	var doc struct {
		Title  string `xml:"title,attr"`
		Suites []struct {
			Name     string `xml:"name,attr"`
			Code     string `xml:"code,attr"`
			Commands []struct {
				Name   string `xml:"name,attr"`
				Code   string `xml:"code,attr"`
				Direct struct {
					Type     string `xml:"type,attr"`
					Optional string `xml:"optional,attr"`
				} `xml:"direct-parameter"`
				Result struct {
					Type string `xml:"type,attr"`
				} `xml:"result"`
			} `xml:"command"`
		} `xml:"suite"`
	}
	if err := xml.Unmarshal([]byte(sdef), &doc); err != nil {
		t.Fatalf("sdef is not well-formed XML: %v\n%s", err, sdef)
	}
	if len(doc.Suites) != 2 {
		t.Fatalf("expected 2 suites, got %d\n%s", len(doc.Suites), sdef)
	}
	// Suites are sorted by class code: "WAIL" < "aevt".
	if doc.Suites[0].Code != "WAIL" || doc.Suites[1].Code != "aevt" {
		t.Fatalf("unexpected suite order: %s, %s", doc.Suites[0].Code, doc.Suites[1].Code)
	}
	wail := doc.Suites[0]
	if len(wail.Commands) != 3 {
		t.Fatalf("expected 3 WAIL commands, got %d", len(wail.Commands))
	}
	if wail.Commands[0].Name != "note" || wail.Commands[0].Code != "WAILnote" {
		t.Fatalf("unexpected first command: %+v", wail.Commands[0])
	}
	if wail.Commands[1].Name != "find" || wail.Commands[1].Code != "WAILFind" {
		t.Fatalf("unexpected second command: %+v", wail.Commands[1])
	}
	if !strings.HasPrefix(wail.Commands[2].Name, "command ") || wail.Commands[2].Code != `WAIL<&>"` {
		t.Fatalf("unexpected escaped command: %+v", wail.Commands[2])
	}
	if wail.Commands[0].Direct.Type != "text" || wail.Commands[0].Direct.Optional != "yes" || wail.Commands[0].Result.Type != "any" {
		t.Fatalf("unexpected parameter declarations: %+v", wail.Commands[0])
	}
	if doc.Suites[1].Commands[0].Code != "aevtodoc" {
		t.Fatalf("unexpected odoc command: %+v", doc.Suites[1].Commands[0])
	}
	if !strings.HasPrefix(sdef, `<?xml version="1.0" encoding="UTF-8"?>`) || !strings.Contains(sdef, "sdef.dtd") {
		t.Fatalf("missing XML prologue or DOCTYPE:\n%s", sdef)
	}
}

func TestSdefCommandName(t *testing.T) {
	tests := map[string]string{
		"note": "note",
		"Find": "find",
		"do X": "do x",
		"a__b": "a b",
		"1abc": "command 31616263",
		"    ": "command 20202020",
	}
	for in, want := range tests {
		if got := sdefCommandName(in); got != want {
			t.Errorf("sdefCommandName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAppleEventCodecRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{name: "nil", value: nil, want: nil},
		{name: "string", value: "hello", want: "hello"},
		{name: "int", value: 42, want: int64(42)},
		{name: "large int", value: int64(1) << 40, want: int64(1) << 40},
		{name: "uint", value: uint16(7), want: int64(7)},
		{name: "float", value: 1.5, want: 1.5},
		{name: "bool", value: true, want: true},
		{name: "file", value: AppleEventFile("/tmp/a.txt"), want: "/tmp/a.txt"},
		{name: "string list", value: []string{"a", "b"}, want: []string{"a", "b"}},
		{name: "file list", value: []AppleEventFile{"/a", "/b"}, want: []string{"/a", "/b"}},
		{name: "mixed list", value: []any{"a", 1, true, nil}, want: []any{"a", int64(1), true, nil}},
		{name: "nested list", value: []any{[]string{"x"}, "y"}, want: []any{[]string{"x"}, "y"}},
		{name: "raw", value: AppleEventRawData{Type: "obj ", Data: []byte{1, 2, 3}}, want: AppleEventRawData{Type: "obj ", Data: []byte{1, 2, 3}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wire, err := appleEventEncode(test.value)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			data, err := json.Marshal(wire)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var decodedWire appleEventWire
			if err := json.Unmarshal(data, &decodedWire); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			got, err := appleEventDecode(decodedWire)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("round trip = %#v, want %#v (wire %s)", got, test.want, data)
			}
		})
	}

	if _, err := appleEventEncode(struct{}{}); err == nil {
		t.Fatal("expected an error for an unsupported type")
	}
	if _, err := appleEventEncode(uint64(1) << 63); err == nil {
		t.Fatal("expected an error for an out of range uint64")
	}
	if _, err := appleEventEncode(AppleEventRawData{Type: "bad"}); err == nil {
		t.Fatal("expected an error for a bad raw type")
	}
	if _, err := appleEventDecode(appleEventWire{Type: "mystery"}); err == nil {
		t.Fatal("expected an error for an unknown wire type")
	}
}

func TestAppleEventDecodeEvent(t *testing.T) {
	event, err := appleEventDecodeEvent(`{"class":"aevt","id":"odoc","direct":{"t":"list","v":[{"t":"file","v":"/a.txt"},{"t":"file","v":"/b.txt"}]},"params":{"keyw":{"t":"int","v":3}}}`)
	if err != nil {
		t.Fatalf("decode event: %v", err)
	}
	if event.Class != "aevt" || event.ID != "odoc" {
		t.Fatalf("unexpected codes: %+v", event)
	}
	if !reflect.DeepEqual(event.DirectObject, []string{"/a.txt", "/b.txt"}) {
		t.Fatalf("unexpected direct object: %#v", event.DirectObject)
	}
	if !reflect.DeepEqual(event.Params, map[string]any{"keyw": int64(3)}) {
		t.Fatalf("unexpected params: %#v", event.Params)
	}
	if _, err := appleEventDecodeEvent(`not json`); err == nil {
		t.Fatal("expected an error for malformed JSON")
	}
}

func TestAppleEventDispatch(t *testing.T) {
	m := newAppleEventsManager(nil)
	var received AppleEvent
	err := m.register("WAIL", "note", func(ctx *Context, event AppleEvent) (AppleEventReply, error) {
		if ctx == nil {
			t.Error("context is nil")
		}
		received = event
		if event.DirectObject == "fail" {
			return AppleEventReply{}, errors.New("boom")
		}
		if event.DirectObject == "custom" {
			return AppleEventReply{ErrorNumber: 7, ErrorString: "seven"}, nil
		}
		if event.DirectObject == "panic" {
			panic("handler panic")
		}
		return AppleEventReply{Result: []string{"ok", event.DirectObject.(string)}}, nil
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	decode := func(replyJSON string) appleEventWireReply {
		var reply appleEventWireReply
		if err := json.Unmarshal([]byte(replyJSON), &reply); err != nil {
			t.Fatalf("bad reply %q: %v", replyJSON, err)
		}
		return reply
	}

	reply := decode(m.dispatch(`{"class":"WAIL","id":"note","direct":{"t":"string","v":"hello"},"params":{"extr":{"t":"bool","v":true}}}`))
	if reply.ErrorNumber != 0 || reply.Result == nil {
		t.Fatalf("unexpected reply: %+v", reply)
	}
	result, err := appleEventDecode(*reply.Result)
	if err != nil || !reflect.DeepEqual(result, []string{"ok", "hello"}) {
		t.Fatalf("unexpected result %#v (%v)", result, err)
	}
	if received.Params["extr"] != true {
		t.Fatalf("params not delivered: %#v", received.Params)
	}

	reply = decode(m.dispatch(`{"class":"WAIL","id":"note","direct":{"t":"string","v":"fail"}}`))
	if reply.ErrorNumber != appleEventFailedErrorNumber || reply.ErrorString != "boom" {
		t.Fatalf("unexpected error reply: %+v", reply)
	}

	reply = decode(m.dispatch(`{"class":"WAIL","id":"note","direct":{"t":"string","v":"custom"}}`))
	if reply.ErrorNumber != 7 || reply.ErrorString != "seven" {
		t.Fatalf("unexpected custom error reply: %+v", reply)
	}

	reply = decode(m.dispatch(`{"class":"WAIL","id":"note","direct":{"t":"string","v":"panic"}}`))
	if reply.ErrorNumber != appleEventFailedErrorNumber || !strings.Contains(reply.ErrorString, "handler panic") {
		t.Fatalf("unexpected panic reply: %+v", reply)
	}

	reply = decode(m.dispatch(`{"class":"WAIL","id":"none"}`))
	if reply.ErrorNumber != appleEventNotHandledErrorNumber {
		t.Fatalf("unexpected unhandled reply: %+v", reply)
	}

	reply = decode(m.dispatch(`garbage`))
	if reply.ErrorNumber != appleEventFailedErrorNumber {
		t.Fatalf("unexpected malformed reply: %+v", reply)
	}
}

func TestAppleEventDecodeReply(t *testing.T) {
	result, err := appleEventDecodeReply(`{"result":{"t":"int","v":5},"errn":0}`)
	if err != nil || result != int64(5) {
		t.Fatalf("unexpected reply %#v, %v", result, err)
	}
	_, err = appleEventDecodeReply(`{"errn":-1728,"errs":"no such object"}`)
	if err == nil || !strings.Contains(err.Error(), "no such object") || !strings.Contains(err.Error(), "-1728") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := appleEventDecodeReply(`{`); err == nil {
		t.Fatal("expected an error for malformed JSON")
	}
}
