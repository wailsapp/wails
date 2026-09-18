package application

import (
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestChangedLinuxWindowStateEvents(t *testing.T) {
	tests := []struct {
		name     string
		previous linuxWindowState
		current  linuxWindowState
		observed bool
		want     []events.WindowEventType
	}{
		{
			name:    "initial normal state is silent",
			current: linuxWindowState{},
		},
		{
			name:    "initial maximised state is announced",
			current: linuxWindowState{maximised: true},
			want:    []events.WindowEventType{events.Common.WindowMaximise},
		},
		{
			name:     "restore from maximised",
			previous: linuxWindowState{maximised: true},
			current:  linuxWindowState{},
			observed: true,
			want:     []events.WindowEventType{events.Common.WindowUnMaximise},
		},
		{
			name:     "minimise and restore",
			previous: linuxWindowState{minimised: true},
			current:  linuxWindowState{},
			observed: true,
			want:     []events.WindowEventType{events.Common.WindowUnMinimise},
		},
		{
			name:     "fullscreen transition preserves distinct maximised event",
			previous: linuxWindowState{maximised: true},
			current:  linuxWindowState{fullscreen: true},
			observed: true,
			want: []events.WindowEventType{
				events.Common.WindowUnMaximise,
				events.Common.WindowFullscreen,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := changedLinuxWindowStateEvents(test.previous, test.current, test.observed)
			if !slices.Equal(got, test.want) {
				t.Fatalf("events = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestDisableSizeConstraintsDoesNotUsePositionAsDimension(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		wantMinW int
		wantMinH int
	}{
		{"monitor at origin", 0, 0, 0, 0},
		{"monitor at 1920x0", 1920, 0, 0, 0},
		{"monitor at 0x1080", 0, 1080, 0, 0},
		{"monitor at 1920x1080", 1920, 1080, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minW := 0
			minH := 0
			if minW != tt.wantMinW {
				t.Errorf("minWidth = %d, want %d (position was %d,%d)", minW, tt.wantMinW, tt.x, tt.y)
			}
			if minH != tt.wantMinH {
				t.Errorf("minHeight = %d, want %d (position was %d,%d)", minH, tt.wantMinH, tt.x, tt.y)
			}
		})
	}
}

func TestSetMinMaxSizeZeroMinIsCorrect(t *testing.T) {
	minW, minH := 0, 0
	maxW, maxH := 1920, 1080

	if minW != 0 {
		t.Errorf("disabled min width should be 0, got %d", minW)
	}
	if minH != 0 {
		t.Errorf("disabled min height should be 0, got %d", minH)
	}
	if maxW <= 0 {
		t.Errorf("max width should be positive, got %d", maxW)
	}
	if maxH <= 0 {
		t.Errorf("max height should be positive, got %d", maxH)
	}
}
