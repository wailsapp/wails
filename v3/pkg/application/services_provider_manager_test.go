package application

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type fakeServicesProvider struct {
	registered []string
}

func (f *fakeServicesProvider) register(name string) error {
	f.registered = append(f.registered, name)
	return nil
}

func newTestServicesProviderManager(impl servicesProviderImpl) *ServicesProviderManager {
	return &ServicesProviderManager{impl: impl, services: map[string]ServiceDefinition{}}
}

func echoService(*Context, ServiceRequest) (ServiceResponse, error) {
	return ServiceResponse{}, nil
}

func TestServiceDefinitionValidation(t *testing.T) {
	valid := ServiceDefinition{
		Name:        "summarise",
		MenuTitle:   "Summarise",
		SendTypes:   []string{"public.utf8-plain-text"},
		ReturnTypes: []string{"public.utf8-plain-text"},
		Handler:     echoService,
	}
	if err := valid.validate(); err != nil {
		t.Fatalf("valid definition rejected: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(*ServiceDefinition)
		want   string
	}{
		{"empty name", func(d *ServiceDefinition) { d.Name = "" }, "Name"},
		{"name with colon", func(d *ServiceDefinition) { d.Name = "do:it" }, "Name"},
		{"name with space", func(d *ServiceDefinition) { d.Name = "do it" }, "Name"},
		{"name starting with digit", func(d *ServiceDefinition) { d.Name = "1st" }, "Name"},
		{"no title", func(d *ServiceDefinition) { d.MenuTitle = "  " }, "MenuTitle"},
		{"no handler", func(d *ServiceDefinition) { d.Handler = nil }, "Handler"},
		{"no types", func(d *ServiceDefinition) { d.SendTypes = nil; d.ReturnTypes = nil }, "SendTypes or ReturnTypes"},
		{"blank type", func(d *ServiceDefinition) { d.SendTypes = []string{""} }, "types must not be empty"},
		{"long key equivalent", func(d *ServiceDefinition) { d.KeyEquivalent = "AB" }, "KeyEquivalent"},
	}
	for _, tc := range cases {
		def := valid
		tc.mutate(&def)
		err := def.validate()
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: got %v, want error mentioning %q", tc.name, err, tc.want)
		}
	}
}

func TestServicesProviderRegister(t *testing.T) {
	impl := &fakeServicesProvider{}
	m := newTestServicesProviderManager(impl)
	def := ServiceDefinition{Name: "summarise", MenuTitle: "Summarise", SendTypes: []string{"public.utf8-plain-text"}, Handler: echoService}
	if err := m.Register(def); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := m.Register(def); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("duplicate Register: got %v", err)
	}
	if err := m.Register(ServiceDefinition{}); err == nil {
		t.Fatal("invalid Register accepted")
	}
	if !reflect.DeepEqual(impl.registered, []string{"summarise"}) {
		t.Fatalf("impl registered %v", impl.registered)
	}
	if got := m.Definitions(); len(got) != 1 || got[0].Name != "summarise" {
		t.Fatalf("Definitions = %+v", got)
	}

	unsupported := newTestServicesProviderManager(nil)
	if err := unsupported.Register(def); !errors.Is(err, ErrServicesUnsupported) {
		t.Fatalf("unsupported platform: got %v", err)
	}
}

func TestServicesProviderHandle(t *testing.T) {
	m := newTestServicesProviderManager(&fakeServicesProvider{})
	var seen ServiceRequest
	_ = m.Register(ServiceDefinition{
		Name:        "upper",
		MenuTitle:   "Upper",
		SendTypes:   []string{"public.utf8-plain-text"},
		ReturnTypes: []string{"public.utf8-plain-text"},
		Handler: func(ctx *Context, req ServiceRequest) (ServiceResponse, error) {
			if ctx == nil {
				t.Error("nil context")
			}
			seen = req
			return ServiceResponse{Text: strings.ToUpper(req.Text)}, nil
		},
	})
	resp, err := m.handle("upper", ServiceRequest{Text: "hello"})
	if err != nil || resp.Text != "HELLO" {
		t.Fatalf("handle = %+v, %v", resp, err)
	}
	if seen.Files == nil || seen.Data == nil || seen.Types == nil {
		t.Fatalf("handle should normalise nil slices and maps: %+v", seen)
	}
	if _, err := m.handle("missing", ServiceRequest{}); err == nil {
		t.Fatal("unknown service handled")
	}
	if !(ServiceResponse{}).IsEmpty() || (ServiceResponse{Text: "x"}).IsEmpty() {
		t.Fatal("IsEmpty is wrong")
	}
}

func TestServicesInfoPlistEntries(t *testing.T) {
	m := newTestServicesProviderManager(&fakeServicesProvider{})
	m.app = &App{options: Options{Name: "Wails Example"}}
	if got := m.InfoPlistEntries(); len(got) != 0 {
		t.Fatalf("entries before registration = %v", got)
	}
	if got := m.InfoPlistXML(); got != "" {
		t.Fatalf("XML before registration = %q", got)
	}
	_ = m.Register(ServiceDefinition{
		Name:          "summarise",
		MenuTitle:     "Summarise with <Wails>",
		SendTypes:     []string{"public.utf8-plain-text", "public.file-url"},
		ReturnTypes:   []string{"public.utf8-plain-text"},
		KeyEquivalent: "S",
		Handler:       echoService,
	})
	_ = m.Register(ServiceDefinition{
		Name:      "archive",
		MenuTitle: "Archive",
		SendTypes: []string{"public.file-url"},
		Handler:   echoService,
	})

	entries := m.InfoPlistEntries()
	want := []map[string]any{
		{
			"NSMenuItem":  map[string]any{"default": "Archive"},
			"NSMessage":   "archive",
			"NSPortName":  "Wails Example",
			"NSSendTypes": []string{"public.file-url"},
		},
		{
			"NSMenuItem":      map[string]any{"default": "Summarise with <Wails>"},
			"NSMessage":       "summarise",
			"NSPortName":      "Wails Example",
			"NSSendTypes":     []string{"public.utf8-plain-text", "public.file-url"},
			"NSReturnTypes":   []string{"public.utf8-plain-text"},
			"NSKeyEquivalent": map[string]any{"default": "S"},
		},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("InfoPlistEntries =\n%#v\nwant\n%#v", entries, want)
	}

	xml := m.InfoPlistXML()
	for _, fragment := range []string{
		"<key>NSServices</key>",
		"<key>NSMessage</key>\n        <string>summarise</string>",
		"<string>Summarise with &lt;Wails&gt;</string>",
		"<key>NSPortName</key>\n        <string>Wails Example</string>",
		"<key>NSSendTypes</key>",
		"<string>public.file-url</string>",
		"<key>NSKeyEquivalent</key>",
	} {
		if !strings.Contains(xml, fragment) {
			t.Errorf("InfoPlistXML missing %q:\n%s", fragment, xml)
		}
	}
	if strings.Count(xml, "<dict>") != 2+3 { // two services, plus NSMenuItem twice and NSKeyEquivalent once
		t.Errorf("unexpected dict count in\n%s", xml)
	}
	if strings.Contains(xml[:strings.Index(xml, "archive")], "summarise") {
		t.Errorf("services are not sorted by name:\n%s", xml)
	}
}
