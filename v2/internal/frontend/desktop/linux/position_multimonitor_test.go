//go:build linux

package linux

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// TestWindowPositionRoundTripOnSecondMonitor is a regression test for
// wailsapp/wails#4003: WindowSetPosition is relative to the current monitor, so
// WindowGetPosition must be too, otherwise SetPosition(GetPosition()) moves the
// window by the monitor's origin.
//
// It needs an X11 display with a second monitor to the right of the primary
// one, and the window must be placed on that monitor. It is skipped unless
// WAILS_TEST_MULTIMONITOR=1 and needs xdotool, e.g. with Xvfb plus
//
//	xrandr --fb 3840x1080 --setmonitor second 1920/508x1080/285+1920+0 none
//
// Run it in a dedicated process: GTK must be initialised on the test's thread.
func TestWindowPositionRoundTripOnSecondMonitor(t *testing.T) {
	if os.Getenv("WAILS_TEST_MULTIMONITOR") != "1" {
		t.Skip("set WAILS_TEST_MULTIMONITOR=1 with an X11 display that has a second monitor")
	}
	runtime.LockOSThread()

	log := logger.New(nil)
	f := NewFrontend(context.Background(), &options.App{Width: 300, Height: 200}, log, binding.NewBindings(log, nil, nil, false, nil), nil)

	xdotool, err := exec.LookPath("xdotool")
	if err != nil {
		t.Skip("xdotool is required to place the window on the second monitor")
	}

	result := make(chan [5]int, 1)
	go func() {
		time.Sleep(2 * time.Second) // let the window map
		// Move the window onto the second monitor (origin x=1920) with an absolute move.
		out, _ := exec.Command(xdotool, "search", "--pid", strconv.Itoa(os.Getpid())).Output()
		ids := strings.Fields(string(out))
		if len(ids) > 0 {
			_ = exec.Command(xdotool, "windowmove", ids[len(ids)-1], "2000", "100").Run()
		}
		time.Sleep(time.Second)
		f.WindowSetPosition(40, 30)
		time.Sleep(time.Second)
		x, y := f.WindowGetPosition()
		f.WindowSetPosition(x+10, y)
		time.Sleep(time.Second)
		x2, y2 := f.WindowGetPosition()
		abs, _ := f.mainWindow.absolutePosition()
		result <- [5]int{x, y, x2, y2, abs}
		invokeOnMainThread(f.mainWindow.Close)
	}()
	// Run must be called on the locked GTK thread, before the main loop.
	_ = f.Run(context.Background())
	f.RunMainLoop()

	got := <-result
	if got[4] < 1920 {
		t.Fatalf("window is not on the second monitor (absolute x = %d); the test setup is wrong", got[4])
	}
	if got[0] != 40 || got[1] != 30 {
		t.Errorf("GetPosition after SetPosition(40,30) = (%d,%d), want (40,30)", got[0], got[1])
	}
	if got[2] != 50 || got[3] != 30 {
		t.Errorf("GetPosition after SetPosition(x+10,y) = (%d,%d), want (50,30)", got[2], got[3])
	}
}
