package application

import (
	"context"
	"testing"
)

// Test service implementations
type testService struct {
	name string
}

func (s *testService) ServiceName() string {
	return s.name
}

type testServiceWithStartup struct {
	started bool
}

func (s *testServiceWithStartup) ServiceStartup(ctx context.Context, options ServiceOptions) error {
	s.started = true
	return nil
}

type testServiceWithShutdown struct {
	shutdown bool
}

func (s *testServiceWithShutdown) ServiceShutdown() error {
	s.shutdown = true
	return nil
}

type testServiceNoInterface struct {
	value int
}

func TestNewService(t *testing.T) {
	svc := &testService{name: "test"}
	service := NewService(svc)

	if service.instance == nil {
		t.Error("NewService should set instance")
	}
	if service.Instance() != svc {
		t.Error("Instance() should return the original service")
	}
}

func TestNewServiceWithOptions(t *testing.T) {
	svc := &testService{name: "original"}
	opts := ServiceOptions{
		Name:  "custom-name",
		Route: "/api",
	}
	service := NewServiceWithOptions(svc, opts)

	if service.Instance() != svc {
		t.Error("Instance() should return the original service")
	}
	if service.options.Name != "custom-name" {
		t.Errorf("options.Name = %q, want %q", service.options.Name, "custom-name")
	}
	if service.options.Route != "/api" {
		t.Errorf("options.Route = %q, want %q", service.options.Route, "/api")
	}
}

func TestGetServiceName_FromOptions(t *testing.T) {
	svc := &testService{name: "service-name"}
	opts := ServiceOptions{Name: "options-name"}
	service := NewServiceWithOptions(svc, opts)

	name := getServiceName(service)
	if name != "options-name" {
		t.Errorf("getServiceName() = %q, want %q (options takes precedence)", name, "options-name")
	}
}

func TestGetServiceName_FromInterface(t *testing.T) {
	svc := &testService{name: "interface-name"}
	service := NewService(svc)

	name := getServiceName(service)
	if name != "interface-name" {
		t.Errorf("getServiceName() = %q, want %q (from interface)", name, "interface-name")
	}
}

func TestGetServiceName_FromType(t *testing.T) {
	svc := &testServiceNoInterface{value: 42}
	service := NewService(svc)

	name := getServiceName(service)
	// Should contain the type name
	if name == "" {
		t.Error("getServiceName() should return type name for services without ServiceName interface")
	}
	// The name should contain "testServiceNoInterface"
	expected := "application.testServiceNoInterface"
	if name != expected {
		t.Errorf("getServiceName() = %q, want %q", name, expected)
	}
}

func TestService_Instance(t *testing.T) {
	svc := &testService{name: "test"}
	service := NewService(svc)

	instance := service.Instance()
	if instance == nil {
		t.Error("Instance() should not return nil")
	}

	// Type assertion to verify it's the correct type
	if _, ok := instance.(*testService); !ok {
		t.Error("Instance() should return the correct type")
	}
}

func TestDefaultServiceOptions(t *testing.T) {
	// Verify DefaultServiceOptions is zero-valued
	if DefaultServiceOptions.Name != "" {
		t.Errorf("DefaultServiceOptions.Name should be empty, got %q", DefaultServiceOptions.Name)
	}
	if DefaultServiceOptions.Route != "" {
		t.Errorf("DefaultServiceOptions.Route should be empty, got %q", DefaultServiceOptions.Route)
	}
	if DefaultServiceOptions.MarshalError != nil {
		t.Error("DefaultServiceOptions.MarshalError should be nil")
	}
}

func TestNewService_UsesDefaultOptions(t *testing.T) {
	svc := &testService{name: "test"}
	service := NewService(svc)

	// Service created with NewService should use DefaultServiceOptions
	if service.options.Name != DefaultServiceOptions.Name {
		t.Error("NewService should use DefaultServiceOptions")
	}
}

func TestServiceOptions_WithMarshalError(t *testing.T) {
	customMarshal := func(err error) []byte {
		return []byte(`{"error": "custom"}`)
	}

	svc := &testService{name: "test"}
	opts := ServiceOptions{
		MarshalError: customMarshal,
	}
	service := NewServiceWithOptions(svc, opts)

	if service.options.MarshalError == nil {
		t.Error("MarshalError should be set")
	}

	result := service.options.MarshalError(nil)
	expected := `{"error": "custom"}`
	if string(result) != expected {
		t.Errorf("MarshalError result = %q, want %q", string(result), expected)
	}
}

func TestServiceStartupInterface(t *testing.T) {
	svc := &testServiceWithStartup{}
	service := NewService(svc)

	// Verify the service implements ServiceStartup
	instance := service.Instance()
	if startup, ok := instance.(ServiceStartup); ok {
		err := startup.ServiceStartup(context.Background(), ServiceOptions{})
		if err != nil {
			t.Errorf("ServiceStartup returned error: %v", err)
		}
		if !svc.started {
			t.Error("ServiceStartup should have been called")
		}
	} else {
		t.Error("testServiceWithStartup should implement ServiceStartup")
	}
}

func TestServiceShutdownInterface(t *testing.T) {
	svc := &testServiceWithShutdown{}
	service := NewService(svc)

	// Verify the service implements ServiceShutdown
	instance := service.Instance()
	if shutdown, ok := instance.(ServiceShutdown); ok {
		err := shutdown.ServiceShutdown()
		if err != nil {
			t.Errorf("ServiceShutdown returned error: %v", err)
		}
		if !svc.shutdown {
			t.Error("ServiceShutdown should have been called")
		}
	} else {
		t.Error("testServiceWithShutdown should implement ServiceShutdown")
	}
}
