package webcontentsview

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// serializedTestImpl intentionally has no locks. The public WebContentsView
// API must serialize access to it even when multiple callers
// issue operations concurrently.
type serializedTestImpl struct {
	url       string
	visible   bool
	destroyed bool
}

func (i *serializedTestImpl) setBounds(application.Rect) {}
func (i *serializedTestImpl) setURL(url string)          { i.url = url }
func (i *serializedTestImpl) setHTML(string)             { i.url = "about:blank" }
func (i *serializedTestImpl) execJS(string)              {}
func (i *serializedTestImpl) goBack()                    {}
func (i *serializedTestImpl) getURL() string             { return i.url }
func (i *serializedTestImpl) takeSnapshot() string       { return "" }
func (i *serializedTestImpl) setVisible(visible bool)    { i.visible = visible }
func (i *serializedTestImpl) attach(application.Window)  {}
func (i *serializedTestImpl) detach()                    {}
func (i *serializedTestImpl) destroy()                   { i.destroyed = true }
func (i *serializedTestImpl) nativeView() unsafe.Pointer { return nil }

func TestWebContentsViewSwitchesBetweenURLAndInlineHTML(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{URL: "https://example.com"})

	view.SetHTML("<h1>inline</h1>")
	if view.options.URL != "" {
		t.Fatalf("URL after SetHTML = %q, want empty", view.options.URL)
	}
	if view.options.HTML != "<h1>inline</h1>" {
		t.Fatalf("HTML after SetHTML = %q", view.options.HTML)
	}

	view.SetURL("https://wails.io")
	if view.options.URL != "https://wails.io" {
		t.Fatalf("URL after SetURL = %q", view.options.URL)
	}
	if view.options.HTML != "" {
		t.Fatalf("HTML after SetURL = %q, want empty", view.options.HTML)
	}
}

func TestWebContentsViewVisibilityDefaultsToShown(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{})
	if !view.visible {
		t.Fatal("new WebContentsView is not visible by default")
	}

	view.Hide()
	if view.visible {
		t.Fatal("Hide did not record the hidden state")
	}

	view.Show()
	if !view.visible {
		t.Fatal("Show did not record the visible state")
	}
}

func TestWebPreferencesOptionalBooleanHelpers(t *testing.T) {
	if !Enabled.IsSet() || !Enabled.Get() {
		t.Fatal("Enabled is not an explicit true preference")
	}
	if !Disabled.IsSet() || Disabled.Get() {
		t.Fatal("Disabled is not an explicit false preference")
	}
}

func TestWebContentsViewDestroyIsPermanent(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{})
	view.Destroy()
	if !view.destroyed {
		t.Fatal("Destroy did not record the destroyed state")
	}

	view.Show()
	if view.visible {
		t.Fatal("Show revived a destroyed WebContentsView")
	}
	if got := view.TakeSnapshot(); got != "" {
		t.Fatalf("snapshot after Destroy = %q, want empty", got)
	}
	if got := view.GetURL(); got != "" {
		t.Fatalf("URL after Destroy = %q, want empty", got)
	}
	view.Destroy()
}

func TestWebContentsViewSerializesConcurrentOperations(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{})
	impl := &serializedTestImpl{}
	view.impl = impl

	var group sync.WaitGroup
	for worker := 0; worker < 16; worker++ {
		group.Add(1)
		go func(worker int) {
			defer group.Done()
			for iteration := 0; iteration < 64; iteration++ {
				switch (worker + iteration) % 5 {
				case 0:
					view.SetURL("https://example.com")
				case 1:
					view.SetHTML("<p>test</p>")
				case 2:
					view.Show()
				case 3:
					view.Hide()
				default:
					view.ExecJS("void 0")
				}
				_ = view.GetURL()
			}
		}(worker)
	}
	group.Wait()

	view.Destroy()
	if !impl.destroyed {
		t.Fatal("Destroy did not reach the serialized native implementation")
	}
}
