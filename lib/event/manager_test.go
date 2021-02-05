package event

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/wailsapp/wails/lib/renderer"
)

func TestLifecycle(t *testing.T) {

	// test fn
	emptyTest := func(data ...interface{}) {
		fmt.Println("emptyTest ran")
	}

	// spin up an instance
	// HELP - WHAT IS THE RIGHT WAY TO DO THIS?
	eventManager := NewManager()
	eventManager.Start(renderer.NewWebView()) // is WebView ok to use?

	// load it with a test event
	eventManager.On("test", emptyTest) // OK

	// begin tests
	// eventManager.Emit("test") // NOPE - blows chunks

	// proper Shutdown
	eventManager.Shutdown()
}

func TestInstantiation(t *testing.T) {
	// NewManager should expose certain public functions
	eventManager := NewManager()

	if reflect.ValueOf(eventManager.Emit).Kind() != reflect.Func {
		t.Errorf("Emit func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.On).Kind() != reflect.Func {
		t.Errorf("On func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.Once).Kind() != reflect.Func {
		t.Errorf("Once func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.OnMultiple).Kind() != reflect.Func {
		t.Errorf("OnMultiple func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.PushEvent).Kind() != reflect.Func {
		t.Errorf("PushEvent func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.Shutdown).Kind() != reflect.Func {
		t.Errorf("Shutdown func is missing from NewManager")
	}

	if reflect.ValueOf(eventManager.Start).Kind() != reflect.Func {
		t.Errorf("Start func is missing from NewManager")
	}
}
