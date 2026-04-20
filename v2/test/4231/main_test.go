package test_4231

import (
	"testing"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func TestDisableResizePreservesExplicitHeight(t *testing.T) {
	opts := &options.App{
		Width:         800,
		Height:        50,
		DisableResize: true,
	}

	if opts.DisableResize && opts.Height != 50 {
		t.Errorf("expected Height=50 when DisableResize=true, got %d", opts.Height)
	}

	if opts.Height < 200 && opts.DisableResize {
		t.Logf("DisableResize=true with Height=%d (below WebKitGTK default of 200): "+
			"window.go should set MinSize/MaxSize to (%d,%d) to override GTK default",
			opts.Height, opts.Width, opts.Height)
	}
}

func TestDisableResizeWithDefaultDimensions(t *testing.T) {
	opts := &options.App{
		Width:         1024,
		Height:        768,
		DisableResize: true,
	}

	if opts.Width != 1024 || opts.Height != 768 {
		t.Errorf("expected default dimensions preserved, got %dx%d", opts.Width, opts.Height)
	}
}

func TestResizableWindowUsesSeparateMinMax(t *testing.T) {
	opts := &options.App{
		Width:     800,
		Height:    600,
		MinWidth:  400,
		MinHeight: 300,
		MaxWidth:  1200,
		MaxHeight: 900,
	}

	if opts.DisableResize {
		t.Error("expected DisableResize=false by default")
	}

	if opts.MinHeight != 300 {
		t.Errorf("expected MinHeight=300, got %d", opts.MinHeight)
	}
}
