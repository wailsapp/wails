//go:build darwin
// +build darwin

package darwin

import (
	"testing"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// TestShouldEnableRetinaDevicePixelRatio verifies that the Retina devicePixelRatio
// override is only enabled when the developer explicitly opts in via
// mac.Options.EnableRetinaDevicePixelRatio. The default-off behaviour is a safety
// invariant: the underlying _setOverrideDeviceScaleFactor: SPI is private and could
// cause Mac App Store rejection, so it must never be active unless requested.
func TestShouldEnableRetinaDevicePixelRatio(t *testing.T) {
	tests := []struct {
		name string
		opts *options.App
		want bool
	}{
		{
			name: "nil Mac options defaults to disabled",
			opts: &options.App{},
			want: false,
		},
		{
			name: "Mac options present but flag unset defaults to disabled",
			opts: &options.App{Mac: &mac.Options{}},
			want: false,
		},
		{
			name: "Mac options present with flag explicitly false stays disabled",
			opts: &options.App{Mac: &mac.Options{EnableRetinaDevicePixelRatio: false}},
			want: false,
		},
		{
			name: "Mac options present with flag enabled opts in",
			opts: &options.App{Mac: &mac.Options{EnableRetinaDevicePixelRatio: true}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldEnableRetinaDevicePixelRatio(tt.opts); got != tt.want {
				t.Errorf("shouldEnableRetinaDevicePixelRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}
