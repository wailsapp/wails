package webcontentsview

import (
	"testing"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestLayoutResolveBottomHalf(t *testing.T) {
	layout := Layout{
		Width:          LayoutValue{Percent: 100},
		Height:         LayoutValue{Percent: 50},
		VerticalAnchor: LayoutAnchorEnd,
	}

	result := layout.Resolve(application.Rect{Width: 800, Height: 600})

	expected := application.Rect{X: 0, Y: 300, Width: 800, Height: 300}
	if result != expected {
		t.Fatalf("expected %#v, got %#v", expected, result)
	}
}

func TestLayoutResolveFromBottomRightWithOffset(t *testing.T) {
	layout := Layout{
		X:                LayoutValue{Pixels: 15},
		Y:                LayoutValue{Pixels: 10},
		Width:            LayoutValue{Pixels: 100},
		Height:           LayoutValue{Pixels: 50},
		HorizontalAnchor: LayoutAnchorEnd,
		VerticalAnchor:   LayoutAnchorEnd,
	}

	result := layout.Resolve(application.Rect{X: 10, Y: 20, Width: 400, Height: 300})

	expected := application.Rect{X: 295, Y: 260, Width: 100, Height: 50}
	if result != expected {
		t.Fatalf("expected %#v, got %#v", expected, result)
	}
}

func TestSetLayoutTracksExplicitLayoutUntilSetBounds(t *testing.T) {
	view := &WebContentsView{
		impl: &layoutTestWebContentsViewImpl{},
	}

	layout := Layout{Height: LayoutValue{Percent: 50}, VerticalAnchor: LayoutAnchorEnd}
	view.SetLayout(layout)

	if view.options.Layout == nil {
		t.Fatal("expected layout to be stored")
	}

	view.SetBounds(application.Rect{X: 10, Y: 20, Width: 300, Height: 200})

	if view.options.Layout != nil {
		t.Fatal("expected SetBounds to clear active layout")
	}
}

type layoutTestWebContentsViewImpl struct{}

func (m *layoutTestWebContentsViewImpl) setBounds(bounds application.Rect) {}
func (m *layoutTestWebContentsViewImpl) setURL(url string)                 {}
func (m *layoutTestWebContentsViewImpl) execJS(js string)                  {}
func (m *layoutTestWebContentsViewImpl) goBack()                           {}
func (m *layoutTestWebContentsViewImpl) getURL() string                    { return "" }
func (m *layoutTestWebContentsViewImpl) takeSnapshot() string              { return "" }
func (m *layoutTestWebContentsViewImpl) attach(window application.Window)  {}
func (m *layoutTestWebContentsViewImpl) detach()                           {}
func (m *layoutTestWebContentsViewImpl) layoutReferenceRect() application.Rect {
	return application.Rect{}
}
func (m *layoutTestWebContentsViewImpl) nativeView() unsafe.Pointer { return nil }
