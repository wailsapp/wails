package application

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ErrAppleEventsNotSupported is returned by AppleEventsManager methods on
// platforms without an Apple Event Manager.
var ErrAppleEventsNotSupported = errors.New("apple events are not supported on this platform")

// appleEventFailedErrorNumber is errAEEventFailed, the generic "the handler
// failed" OSStatus reported to the sender when a handler returns an error
// without choosing its own ErrorNumber.
const appleEventFailedErrorNumber = -10000

// appleEventNotHandledErrorNumber is errAEEventNotHandled.
const appleEventNotHandledErrorNumber = -1708

// AppleEvent is a decoded incoming Apple Event.
//
// DirectObject is the event's direct parameter ("----"): a string, a
// []string of file paths (for example the files of an "odoc" event), an
// int64, a float64, a bool, a []any for mixed lists, an AppleEventRawData
// for anything the codec does not understand, or nil when the event has no
// direct parameter. Params holds every other parameter keyed by its
// four-character keyword, decoded the same way.
type AppleEvent struct {
	Class        string
	ID           string
	DirectObject any
	Params       map[string]any
}

// AppleEventReply is what a handler returns to the sender.
//
// Result becomes the reply's direct parameter and accepts the same Go values
// the decoder produces (string, integers, floats, bool, AppleEventFile,
// []string, []any, AppleEventRawData or nil). A non-zero ErrorNumber is
// reported to the sender as the "errn" parameter together with ErrorString
// ("errs"); AppleScript surfaces both as a script error.
type AppleEventReply struct {
	Result      any
	ErrorNumber int
	ErrorString string
}

// AppleEventHandler handles one Apple Event. Returning a non-nil error
// reports errAEEventFailed (-10000) with the error text to the sender unless
// the reply already carries an ErrorNumber.
type AppleEventHandler func(*Context, AppleEvent) (AppleEventReply, error)

// AppleEventFile is a file path that encodes as a file URL descriptor rather
// than plain text, for example as the direct object of Send or in a reply.
type AppleEventFile string

// AppleEventRawData carries a descriptor the codec does not decode: Type is
// the four-character descriptor type and Data its raw bytes. It round-trips
// unchanged when used in a reply or as the direct object of Send.
type AppleEventRawData struct {
	Type string
	Data []byte
}

// AppleEventsManager registers handlers for Apple Events so scripts and other
// applications can drive this one on macOS.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type AppleEventsManager struct {
	app *App

	lock     sync.Mutex
	handlers map[appleEventKey]AppleEventHandler
	order    []appleEventKey
}

type appleEventKey struct {
	class string
	id    string
}

func newAppleEventsManager(app *App) *AppleEventsManager {
	return &AppleEventsManager{
		app:      app,
		handlers: map[appleEventKey]AppleEventHandler{},
	}
}

// validateFourCharCode checks that code is a four-character Apple Event code
// such as "aevt", "odoc" or "GURL": exactly four printable ASCII bytes.
func validateFourCharCode(code string) error {
	if len(code) != 4 {
		return fmt.Errorf("apple events: code %q must be exactly four characters", code)
	}
	for i := 0; i < len(code); i++ {
		if code[i] < 0x20 || code[i] > 0x7e {
			return fmt.Errorf("apple events: code %q must contain printable ASCII only", code)
		}
	}
	return nil
}

// Handle registers handler for events of the given class and ID. The codes
// are four-character strings, for example Handle("WAIL", "note", ...) for a
// custom command or Handle("aevt", "odoc", ...) for Open Documents.
// Registering the same class and ID again replaces the previous handler.
//
// Handlers run on their own goroutine while the event is suspended, so they
// may take their time and call back into the application; the sender
// receives the reply when the handler returns.
//
// Wails already handles the Get URL event ("GURL"/"GURL") to deliver custom
// URL schemes as events.Common.ApplicationLaunchedWithUrl. A handler
// registered for "GURL"/"GURL" chains to that built-in behaviour: the URL is
// still delivered as an application event, then the handler runs and its
// reply is returned to the sender. A handler for "aevt"/"odoc" replaces
// AppKit's Open Documents delivery, so application:openFile: is no longer
// called for those files while the handler is registered.
//
// macOS: NSAppleEventManager. Other platforms return ErrAppleEventsNotSupported.
func (m *AppleEventsManager) Handle(eventClass, eventID string, handler AppleEventHandler) error {
	if err := m.register(eventClass, eventID, handler); err != nil {
		return err
	}
	return appleEventsRegisterNative(m, eventClass, eventID)
}

// register validates and stores a handler without touching the platform.
func (m *AppleEventsManager) register(eventClass, eventID string, handler AppleEventHandler) error {
	if err := validateFourCharCode(eventClass); err != nil {
		return err
	}
	if err := validateFourCharCode(eventID); err != nil {
		return err
	}
	if handler == nil {
		return errors.New("apple events: handler must not be nil")
	}
	key := appleEventKey{class: eventClass, id: eventID}
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.handlers[key]; !exists {
		m.order = append(m.order, key)
	}
	m.handlers[key] = handler
	return nil
}

func (m *AppleEventsManager) handler(eventClass, eventID string) AppleEventHandler {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.handlers[appleEventKey{class: eventClass, id: eventID}]
}

// Send sends an Apple Event to another application identified by bundle ID,
// for example Send("com.apple.finder", "misc", "actv", nil), and waits for
// the reply. directObject accepts the same values as AppleEventReply.Result.
// The reply's direct parameter is decoded and returned; an "errn" parameter
// in the reply becomes an error carrying the "errs" text.
//
// The target must be running. Sending events to other applications requires
// the NSAppleEventsUsageDescription Info.plist key in a bundled application
// and the user's consent (System Settings > Privacy & Security > Automation).
// Send blocks the calling goroutine, never the main thread.
//
// macOS 10.11+: NSAppleEventDescriptor sendEventWithOptions:timeout:error:.
// Other platforms return ErrAppleEventsNotSupported.
func (m *AppleEventsManager) Send(target, eventClass, eventID string, directObject any) (any, error) {
	if strings.TrimSpace(target) == "" {
		return nil, errors.New("apple events: target bundle identifier is empty")
	}
	if err := validateFourCharCode(eventClass); err != nil {
		return nil, err
	}
	if err := validateFourCharCode(eventID); err != nil {
		return nil, err
	}
	wire, err := appleEventEncode(directObject)
	if err != nil {
		return nil, err
	}
	directJSON, err := json.Marshal(wire)
	if err != nil {
		return nil, err
	}
	replyJSON, err := appleEventsSendNative(target, eventClass, eventID, string(directJSON))
	if err != nil {
		return nil, err
	}
	return appleEventDecodeReply(replyJSON)
}

// ScriptingDefinition returns a minimal scripting definition (sdef) XML
// document describing the registered handlers: one suite per event class
// whose code is the class, and one command per event ID whose code is the
// class and ID joined ("WAILnote"). Each command takes an optional text direct
// parameter and returns a value, which is enough for AppleScript to call
// custom commands by name instead of the raw «event WAILnote» syntax.
//
// To ship it, write the XML to a file such as Contents/Resources/App.sdef in
// the application bundle and add two keys to Info.plist:
//
//	<key>NSAppleScriptEnabled</key>
//	<true/>
//	<key>OSAScriptingDefinition</key>
//	<string>App.sdef</string>
//
// NSAppleScriptEnabled tells the system the application is scriptable and
// OSAScriptingDefinition names the sdef file (relative to Resources) that
// Script Editor opens from File > Open Dictionary. Handlers can be invoked
// from a script without an sdef by using the raw event syntax, for example
// tell application id "com.example.app" to «event WAILnote» "hello".
//
// The document is generated on every platform; only delivery needs macOS.
func (m *AppleEventsManager) ScriptingDefinition() string {
	appName := "Application"
	if m.app != nil && strings.TrimSpace(m.app.options.Name) != "" {
		appName = strings.TrimSpace(m.app.options.Name)
	}

	m.lock.Lock()
	keys := append([]appleEventKey(nil), m.order...)
	m.lock.Unlock()

	suites := map[string][]appleEventKey{}
	var classes []string
	for _, key := range keys {
		if _, seen := suites[key.class]; !seen {
			classes = append(classes, key.class)
		}
		suites[key.class] = append(suites[key.class], key)
	}
	sort.Strings(classes)

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<!DOCTYPE dictionary SYSTEM "file://localhost/System/Library/DTDs/sdef.dtd">` + "\n")
	fmt.Fprintf(&b, `<dictionary title="%s Terminology">`+"\n", sdefEscape(appName))
	for _, class := range classes {
		fmt.Fprintf(&b, `  <suite name="%s Suite" code="%s" description="Commands handled by %s.">`+"\n",
			sdefEscape(appName+" "+class), sdefEscape(class), sdefEscape(appName))
		for _, key := range suites[class] {
			fmt.Fprintf(&b, `    <command name="%s" code="%s%s" description="Sends the %s/%s Apple Event to %s.">`+"\n",
				sdefEscape(sdefCommandName(key.id)), sdefEscape(key.class), sdefEscape(key.id),
				sdefEscape(key.class), sdefEscape(key.id), sdefEscape(appName))
			b.WriteString(`      <direct-parameter type="text" optional="yes" description="The direct parameter."/>` + "\n")
			b.WriteString(`      <result type="any" description="The handler's result."/>` + "\n")
			b.WriteString("    </command>\n")
		}
		b.WriteString("  </suite>\n")
	}
	b.WriteString("</dictionary>\n")
	return b.String()
}

// sdefCommandName turns an event ID into a command name AppleScript can
// parse: lower-case letters, digits and single spaces, starting with a
// letter. IDs with nothing usable fall back to "command" plus the hex of the
// code.
func sdefCommandName(id string) string {
	var b strings.Builder
	lastSpace := true
	for _, r := range strings.ToLower(id) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastSpace = false
		case r == ' ' || r == '-' || r == '_':
			if !lastSpace {
				b.WriteByte(' ')
				lastSpace = true
			}
		}
	}
	name := strings.TrimSpace(b.String())
	if name == "" || (name[0] >= '0' && name[0] <= '9') {
		return "command " + fmt.Sprintf("%x", []byte(id))
	}
	return name
}

func sdefEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Wire codec.
//
// Descriptors cross the cgo boundary as JSON so the Go side never touches
// NSAppleEventDescriptor. Every value is {"t": kind, "v": payload}:
//
//	string  {"t":"string","v":"text"}
//	int     {"t":"int","v":123}            (64-bit)
//	float   {"t":"float","v":1.5}
//	bool    {"t":"bool","v":true}
//	file    {"t":"file","v":"/path"}       (file URL descriptors)
//	list    {"t":"list","v":[...]}
//	null    {"t":"null"}
//	raw     {"t":"raw","type":"xxxx","v":"base64"}
//
// An incoming event is {"class":"WAIL","id":"note","direct":value,
// "params":{"keyw":value}} and a reply is {"result":value,"errn":0,"errs":""}.

type appleEventWire struct {
	Type     string          `json:"t"`
	Value    json.RawMessage `json:"v,omitempty"`
	DescType string          `json:"type,omitempty"`
}

type appleEventWireEvent struct {
	Class  string                    `json:"class"`
	ID     string                    `json:"id"`
	Direct *appleEventWire           `json:"direct,omitempty"`
	Params map[string]appleEventWire `json:"params,omitempty"`
}

type appleEventWireReply struct {
	Result      *appleEventWire `json:"result,omitempty"`
	ErrorNumber int             `json:"errn"`
	ErrorString string          `json:"errs,omitempty"`
}

func appleEventWireValue(kind string, value any) (appleEventWire, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return appleEventWire{}, err
	}
	return appleEventWire{Type: kind, Value: raw}, nil
}

// appleEventEncode converts a Go value into its wire form.
func appleEventEncode(value any) (appleEventWire, error) {
	switch v := value.(type) {
	case nil:
		return appleEventWire{Type: "null"}, nil
	case string:
		return appleEventWireValue("string", v)
	case AppleEventFile:
		return appleEventWireValue("file", string(v))
	case bool:
		return appleEventWireValue("bool", v)
	case int:
		return appleEventWireValue("int", int64(v))
	case int8:
		return appleEventWireValue("int", int64(v))
	case int16:
		return appleEventWireValue("int", int64(v))
	case int32:
		return appleEventWireValue("int", int64(v))
	case int64:
		return appleEventWireValue("int", v)
	case uint:
		return appleEventWireValue("int", int64(v))
	case uint8:
		return appleEventWireValue("int", int64(v))
	case uint16:
		return appleEventWireValue("int", int64(v))
	case uint32:
		return appleEventWireValue("int", int64(v))
	case uint64:
		if v > 1<<63-1 {
			return appleEventWire{}, fmt.Errorf("apple events: %d does not fit a 64-bit signed integer", v)
		}
		return appleEventWireValue("int", int64(v))
	case float32:
		return appleEventWireValue("float", float64(v))
	case float64:
		return appleEventWireValue("float", v)
	case AppleEventRawData:
		if err := validateFourCharCode(v.Type); err != nil {
			return appleEventWire{}, err
		}
		wire, err := appleEventWireValue("raw", base64.StdEncoding.EncodeToString(v.Data))
		if err != nil {
			return appleEventWire{}, err
		}
		wire.DescType = v.Type
		return wire, nil
	case []string:
		items := make([]any, len(v))
		for i := range v {
			items[i] = v[i]
		}
		return appleEventEncodeList(items)
	case []AppleEventFile:
		items := make([]any, len(v))
		for i := range v {
			items[i] = v[i]
		}
		return appleEventEncodeList(items)
	case []any:
		return appleEventEncodeList(v)
	default:
		return appleEventWire{}, fmt.Errorf("apple events: cannot encode %T", value)
	}
}

func appleEventEncodeList(items []any) (appleEventWire, error) {
	wires := make([]appleEventWire, 0, len(items))
	for _, item := range items {
		wire, err := appleEventEncode(item)
		if err != nil {
			return appleEventWire{}, err
		}
		wires = append(wires, wire)
	}
	return appleEventWireValue("list", wires)
}

// appleEventDecode converts a wire value back into a Go value. Lists whose
// items are all strings (including file paths) decode as []string.
func appleEventDecode(wire appleEventWire) (any, error) {
	switch wire.Type {
	case "null", "":
		return nil, nil
	case "string", "file":
		var s string
		if err := json.Unmarshal(wire.Value, &s); err != nil {
			return nil, fmt.Errorf("apple events: bad %s value: %w", wire.Type, err)
		}
		return s, nil
	case "bool":
		var b bool
		if err := json.Unmarshal(wire.Value, &b); err != nil {
			return nil, fmt.Errorf("apple events: bad bool value: %w", err)
		}
		return b, nil
	case "int":
		var n json.Number
		if err := json.Unmarshal(wire.Value, &n); err != nil {
			return nil, fmt.Errorf("apple events: bad int value: %w", err)
		}
		i, err := strconv.ParseInt(n.String(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("apple events: bad int value: %w", err)
		}
		return i, nil
	case "float":
		var f float64
		if err := json.Unmarshal(wire.Value, &f); err != nil {
			return nil, fmt.Errorf("apple events: bad float value: %w", err)
		}
		return f, nil
	case "list":
		var wires []appleEventWire
		if err := json.Unmarshal(wire.Value, &wires); err != nil {
			return nil, fmt.Errorf("apple events: bad list value: %w", err)
		}
		items := make([]any, 0, len(wires))
		allStrings := true
		for _, item := range wires {
			value, err := appleEventDecode(item)
			if err != nil {
				return nil, err
			}
			if _, ok := value.(string); !ok {
				allStrings = false
			}
			items = append(items, value)
		}
		if allStrings {
			strs := make([]string, len(items))
			for i := range items {
				strs[i] = items[i].(string)
			}
			return strs, nil
		}
		return items, nil
	case "raw":
		var encoded string
		if err := json.Unmarshal(wire.Value, &encoded); err != nil {
			return nil, fmt.Errorf("apple events: bad raw value: %w", err)
		}
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("apple events: bad raw value: %w", err)
		}
		return AppleEventRawData{Type: wire.DescType, Data: data}, nil
	default:
		return nil, fmt.Errorf("apple events: unknown wire type %q", wire.Type)
	}
}

// appleEventDecodeEvent parses the JSON produced by the native handler.
func appleEventDecodeEvent(eventJSON string) (AppleEvent, error) {
	var wire appleEventWireEvent
	if err := json.Unmarshal([]byte(eventJSON), &wire); err != nil {
		return AppleEvent{}, fmt.Errorf("apple events: bad event: %w", err)
	}
	event := AppleEvent{Class: wire.Class, ID: wire.ID, Params: map[string]any{}}
	if wire.Direct != nil {
		direct, err := appleEventDecode(*wire.Direct)
		if err != nil {
			return AppleEvent{}, err
		}
		event.DirectObject = direct
	}
	for keyword, value := range wire.Params {
		decoded, err := appleEventDecode(value)
		if err != nil {
			return AppleEvent{}, err
		}
		event.Params[keyword] = decoded
	}
	return event, nil
}

// appleEventEncodeReply serialises a reply for the native side.
func appleEventEncodeReply(reply AppleEventReply, handlerErr error) string {
	wire := appleEventWireReply{ErrorNumber: reply.ErrorNumber, ErrorString: reply.ErrorString}
	if handlerErr != nil && wire.ErrorNumber == 0 {
		wire.ErrorNumber = appleEventFailedErrorNumber
		if wire.ErrorString == "" {
			wire.ErrorString = handlerErr.Error()
		}
	}
	if reply.Result != nil {
		result, err := appleEventEncode(reply.Result)
		if err != nil {
			if wire.ErrorNumber == 0 {
				wire.ErrorNumber = appleEventFailedErrorNumber
				wire.ErrorString = err.Error()
			}
		} else {
			wire.Result = &result
		}
	}
	data, err := json.Marshal(wire)
	if err != nil {
		return `{"errn":-10000,"errs":"apple events: could not encode reply"}`
	}
	return string(data)
}

// appleEventDecodeReply parses a reply produced by the native Send.
func appleEventDecodeReply(replyJSON string) (any, error) {
	var wire appleEventWireReply
	if err := json.Unmarshal([]byte(replyJSON), &wire); err != nil {
		return nil, fmt.Errorf("apple events: bad reply: %w", err)
	}
	var result any
	if wire.Result != nil {
		decoded, err := appleEventDecode(*wire.Result)
		if err != nil {
			return nil, err
		}
		result = decoded
	}
	if wire.ErrorNumber != 0 {
		message := wire.ErrorString
		if message == "" {
			message = "the target reported an error"
		}
		return result, fmt.Errorf("apple events: %s (error %d)", message, wire.ErrorNumber)
	}
	return result, nil
}

// dispatch decodes an incoming event, runs the matching handler and returns
// the encoded reply. Unknown events and panicking handlers produce an error
// reply rather than a missing one so the sender never waits on a reply that
// is not coming.
func (m *AppleEventsManager) dispatch(eventJSON string) (reply string) {
	event, err := appleEventDecodeEvent(eventJSON)
	if err != nil {
		return appleEventEncodeReply(AppleEventReply{}, err)
	}
	handler := m.handler(event.Class, event.ID)
	if handler == nil {
		return appleEventEncodeReply(AppleEventReply{ErrorNumber: appleEventNotHandledErrorNumber},
			fmt.Errorf("apple events: no handler for %s/%s", event.Class, event.ID))
	}
	defer func() {
		if r := recover(); r != nil {
			panicErr, ok := r.(error)
			if !ok {
				panicErr = fmt.Errorf("%v", r)
			}
			if m.app != nil {
				m.app.handleError(fmt.Errorf("apple events: handler for %s/%s panicked: %w", event.Class, event.ID, panicErr))
			}
			reply = appleEventEncodeReply(AppleEventReply{}, panicErr)
		}
	}()
	result, handlerErr := handler(newContext(), event)
	if handlerErr != nil && m.app != nil {
		m.app.handleError(fmt.Errorf("apple events: handler for %s/%s: %w", event.Class, event.ID, handlerErr))
	}
	return appleEventEncodeReply(result, handlerErr)
}
