package application

import (
	"errors"
	"fmt"
	"html"
	"sort"
	"strings"
	"sync"
)

// ServiceDefinition describes one entry this application contributes to the
// macOS Services menu (the "Services" submenu of every application menu and
// the context menu of selected text or files).
type ServiceDefinition struct {
	// Name identifies the service and doubles as the Objective-C selector
	// AppKit invokes (`<Name>:userData:error:`). It must be a valid
	// identifier: letters, digits and underscores, not starting with a
	// digit. It is also the value of the NSMessage Info.plist key.
	Name string
	// MenuTitle is the user-visible title in the Services menu, for example
	// "Summarise with MyApp".
	MenuTitle string
	// SendTypes lists the pasteboard types the service accepts from the
	// requesting application, as UTIs ("public.utf8-plain-text",
	// "public.file-url") or legacy pasteboard type names. The service is
	// offered only when the selection provides at least one of them. Leave
	// it empty for a service that only produces data.
	SendTypes []string
	// ReturnTypes lists the pasteboard types the service can write back.
	// Leave it empty for a service that only consumes data; the requesting
	// application then leaves its selection untouched.
	ReturnTypes []string
	// KeyEquivalent is an optional single character used with Command and
	// Shift as the menu shortcut ("S" for Cmd+Shift+S).
	KeyEquivalent string
	// Handler runs when the service is invoked. It receives the pasteboard
	// contents and returns what should be written back. An empty
	// ServiceResponse leaves the pasteboard untouched. The handler runs
	// while AppKit waits on the main thread, so keep it fast; hand long
	// work to a goroutine and return.
	Handler func(*Context, ServiceRequest) (ServiceResponse, error)
}

// ServiceRequest is the pasteboard content handed to a service handler.
type ServiceRequest struct {
	// Text is the plain text on the pasteboard, if any.
	Text string `json:"text"`
	// Files holds the paths of file URLs on the pasteboard, if any.
	Files []string `json:"files"`
	// Data maps each of the service's SendTypes present on the pasteboard
	// to its raw bytes.
	Data map[string][]byte `json:"data"`
	// Types lists every type the pasteboard offered.
	Types []string `json:"types"`
}

// ServiceResponse is what a service handler writes back to the pasteboard.
// The zero value writes nothing.
type ServiceResponse struct {
	// Text is written as plain text.
	Text string `json:"text"`
	// Files are written as file URLs.
	Files []string `json:"files"`
	// Data is written under each type identifier verbatim.
	Data map[string][]byte `json:"data"`
}

// IsEmpty reports whether the response carries nothing to write back.
func (r ServiceResponse) IsEmpty() bool {
	return r.Text == "" && len(r.Files) == 0 && len(r.Data) == 0
}

// ErrServicesUnsupported is returned by Register on platforms without a
// Services menu.
var ErrServicesUnsupported = errors.New("services: the Services menu is not supported on this platform")

// servicesProviderImpl is the platform side of ServicesProviderManager.
type servicesProviderImpl interface {
	// register makes the named selector resolvable on the services provider
	// object and refreshes the dynamic services registry.
	register(name string) error
}

// ServicesProviderManager registers handlers for the macOS Services menu so
// other applications can send this one text or files.
//
// Registration alone is not enough for the service to show up: the
// application bundle's Info.plist must also declare it under NSServices.
// InfoPlistEntries and InfoPlistXML produce those entries, and the CLI's
// build/config.yml accepts a `services` list that is rendered into the
// generated Info.plist (see the mac-integration example).
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type ServicesProviderManager struct {
	app  *App
	impl servicesProviderImpl

	lock     sync.RWMutex
	services map[string]ServiceDefinition
}

func newServicesProviderManager(app *App) *ServicesProviderManager {
	return &ServicesProviderManager{
		app:      app,
		impl:     newServicesProviderImpl(app),
		services: map[string]ServiceDefinition{},
	}
}

var errServiceName = errors.New("services: Name must be a valid identifier (letters, digits and underscores, not starting with a digit)")

func validServiceName(name string) bool {
	if name == "" {
		return false
	}
	for index, r := range name {
		switch {
		case r == '_':
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
			if index == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (def ServiceDefinition) validate() error {
	if !validServiceName(def.Name) {
		return errServiceName
	}
	if strings.TrimSpace(def.MenuTitle) == "" {
		return errors.New("services: MenuTitle is required")
	}
	if def.Handler == nil {
		return errors.New("services: Handler is required")
	}
	if len(def.SendTypes) == 0 && len(def.ReturnTypes) == 0 {
		return errors.New("services: at least one of SendTypes or ReturnTypes is required")
	}
	for _, t := range append(append([]string{}, def.SendTypes...), def.ReturnTypes...) {
		if strings.TrimSpace(t) == "" {
			return errors.New("services: pasteboard types must not be empty")
		}
	}
	if len([]rune(def.KeyEquivalent)) > 1 {
		return errors.New("services: KeyEquivalent must be a single character")
	}
	return nil
}

// Register adds a service. Its Name must be unique within the application.
// Register may be called before or after App.Run.
//
// macOS: NSApplication.servicesProvider receives a `<Name>:userData:error:`
// message when the user picks the service; the pasteboard is read into a
// ServiceRequest, the handler runs, and a non-empty ServiceResponse is
// written back. NSUpdateDynamicServices is called after each registration.
// The service still needs its NSServices Info.plist entry (see
// InfoPlistEntries); unbundled `go run` binaries have no Info.plist and so
// never appear in the Services menu.
//
// Other platforms: ErrServicesUnsupported.
func (m *ServicesProviderManager) Register(def ServiceDefinition) error {
	if err := def.validate(); err != nil {
		return err
	}
	if m.impl == nil {
		return ErrServicesUnsupported
	}
	m.lock.Lock()
	if _, exists := m.services[def.Name]; exists {
		m.lock.Unlock()
		return fmt.Errorf("services: %q is already registered", def.Name)
	}
	m.services[def.Name] = def
	m.lock.Unlock()
	return m.impl.register(def.Name)
}

// Definitions returns the registered services sorted by Name.
func (m *ServicesProviderManager) Definitions() []ServiceDefinition {
	m.lock.RLock()
	defer m.lock.RUnlock()
	result := make([]ServiceDefinition, 0, len(m.services))
	for _, def := range m.services {
		result = append(result, def)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

// InfoPlistEntries returns the NSServices array entries for the registered
// services, in the shape a plist serialiser or the CLI packaging expects:
// one map per service with NSMenuItem, NSMessage, NSPortName, NSSendTypes,
// NSReturnTypes and, when set, NSKeyEquivalent. NSPortName is the
// application's Name, which must match CFBundleName.
func (m *ServicesProviderManager) InfoPlistEntries() []map[string]any {
	entries := []map[string]any{}
	for _, def := range m.Definitions() {
		entries = append(entries, serviceInfoPlistEntry(def, m.portName()))
	}
	return entries
}

func (m *ServicesProviderManager) portName() string {
	if m.app == nil {
		return ""
	}
	return m.app.options.Name
}

func serviceInfoPlistEntry(def ServiceDefinition, portName string) map[string]any {
	entry := map[string]any{
		"NSMenuItem": map[string]any{"default": def.MenuTitle},
		"NSMessage":  def.Name,
	}
	if portName != "" {
		entry["NSPortName"] = portName
	}
	if len(def.SendTypes) > 0 {
		entry["NSSendTypes"] = append([]string{}, def.SendTypes...)
	}
	if len(def.ReturnTypes) > 0 {
		entry["NSReturnTypes"] = append([]string{}, def.ReturnTypes...)
	}
	if def.KeyEquivalent != "" {
		entry["NSKeyEquivalent"] = map[string]any{"default": def.KeyEquivalent}
	}
	return entry
}

// InfoPlistXML returns the NSServices key and array as plist XML, ready to
// paste inside the top-level <dict> of Info.plist. It returns "" when no
// services are registered.
func (m *ServicesProviderManager) InfoPlistXML() string {
	defs := m.Definitions()
	if len(defs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("<key>NSServices</key>\n<array>\n")
	for _, def := range defs {
		b.WriteString(serviceInfoPlistXML(def, m.portName()))
	}
	b.WriteString("</array>\n")
	return b.String()
}

func serviceInfoPlistXML(def ServiceDefinition, portName string) string {
	esc := html.EscapeString
	var b strings.Builder
	b.WriteString("    <dict>\n")
	b.WriteString("        <key>NSMenuItem</key>\n        <dict>\n            <key>default</key>\n            <string>" + esc(def.MenuTitle) + "</string>\n        </dict>\n")
	b.WriteString("        <key>NSMessage</key>\n        <string>" + esc(def.Name) + "</string>\n")
	if portName != "" {
		b.WriteString("        <key>NSPortName</key>\n        <string>" + esc(portName) + "</string>\n")
	}
	writeTypes := func(key string, types []string) {
		if len(types) == 0 {
			return
		}
		b.WriteString("        <key>" + key + "</key>\n        <array>\n")
		for _, t := range types {
			b.WriteString("            <string>" + esc(t) + "</string>\n")
		}
		b.WriteString("        </array>\n")
	}
	writeTypes("NSSendTypes", def.SendTypes)
	writeTypes("NSReturnTypes", def.ReturnTypes)
	if def.KeyEquivalent != "" {
		b.WriteString("        <key>NSKeyEquivalent</key>\n        <dict>\n            <key>default</key>\n            <string>" + esc(def.KeyEquivalent) + "</string>\n        </dict>\n")
	}
	b.WriteString("    </dict>\n")
	return b.String()
}

// handle runs the named service's handler. It is called by the platform
// layer with the decoded pasteboard and returns the response to write back.
func (m *ServicesProviderManager) handle(name string, request ServiceRequest) (ServiceResponse, error) {
	m.lock.RLock()
	def, ok := m.services[name]
	m.lock.RUnlock()
	if !ok {
		return ServiceResponse{}, fmt.Errorf("services: %q is not registered", name)
	}
	if request.Files == nil {
		request.Files = []string{}
	}
	if request.Data == nil {
		request.Data = map[string][]byte{}
	}
	if request.Types == nil {
		request.Types = []string{}
	}
	return def.Handler(newContext(), request)
}
