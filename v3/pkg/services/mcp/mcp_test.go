package mcp

import (
	"testing"
	"time"
)

func TestNewWithConfigDefaults(t *testing.T) {
	service := New()
	if service.config.Host != "127.0.0.1" {
		t.Errorf("unexpected default host: %s", service.config.Host)
	}
	if service.config.Port != 9099 {
		t.Errorf("unexpected default port: %d", service.config.Port)
	}
	if service.config.EvalTimeout != 30*time.Second {
		t.Errorf("unexpected default eval timeout: %s", service.config.EvalTimeout)
	}
}

func TestNewWithConfigOverrides(t *testing.T) {
	service := NewWithConfig(Config{Host: "0.0.0.0", Port: -1, EvalTimeout: time.Second, HideCursor: true})
	if service.config.Host != "0.0.0.0" {
		t.Errorf("host override ignored: %s", service.config.Host)
	}
	if service.config.Port != -1 {
		t.Errorf("port override ignored: %d", service.config.Port)
	}
	if service.config.EvalTimeout != time.Second {
		t.Errorf("timeout override ignored: %s", service.config.EvalTimeout)
	}
	if !service.config.HideCursor {
		t.Error("HideCursor override ignored")
	}
}

func TestServiceName(t *testing.T) {
	if New().ServiceName() != "github.com/wailsapp/wails/v3/services/mcp" {
		t.Errorf("unexpected service name: %s", New().ServiceName())
	}
}
