//go:build windows

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// systray-stress is a test harness that exercises SystemTray.SetMenu and
// SystemTray.OpenMenu under workloads that reproduce the Windows systray
// crash family described in the investigation handoff. It emits line-based
// log events to stderr so a supervisor can record iterations,
// GetGuiResources deltas, and exit cause.
//
// Modes:
//   - churn        rebuild a fresh Menu every iteration and SetMenu it.
//   - show         repeatedly OpenMenu + ESC-dismiss.
//   - churn+show   both of the above concurrently.
//   - mutate       keep the same Menu across iterations and call SetBitmap
//                  on items that have already been built. The first build
//                  assigns impls to the items; subsequent SetMenu calls
//                  destroy the prior Win32Menu and rebuild fresh impls.
//                  This is the "runtime SetBitmap" path the -bitmaps
//                  churn workload does NOT exercise (churn allocates a
//                  fresh Menu each iteration, so items never have an
//                  impl at SetBitmap time and the call is a no-op on
//                  the native side). The leak that previously orphaned
//                  an HBITMAP per SetBitmap call was fixed by having
//                  freeBitmaps walk menuMapping and release impl.bitmap
//                  before DestroyMenu — this workload regression-guards
//                  that cleanup path.

const (
	GR_GDIOBJECTS  = 0
	GR_USEROBJECTS = 1

	VK_ESCAPE    = 0x1B
	INPUT_KEY    = 1
	KEYEVENTF_KU = 0x0002 // KEYEVENTF_KEYUP
)

var (
	moduser32             = syscall.NewLazyDLL("user32.dll")
	modkernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetGuiResources   = moduser32.NewProc("GetGuiResources")
	procGetCurrentProcess = modkernel32.NewProc("GetCurrentProcess")
	procSendInput         = moduser32.NewProc("SendInput")
)

// keyboardInput maps to KEYBDINPUT packed into INPUT.
type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
	_padding    [8]byte // pad INPUT union to largest member (MOUSEINPUT)
}

// input corresponds to the Win32 INPUT struct with a keyboard union.
type input struct {
	inputType uint32
	_         uint32 // alignment to 8 bytes on amd64
	ki        keyboardInput
}

func getGuiResources(flag uint32) uint32 {
	hProc, _, _ := procGetCurrentProcess.Call()
	n, _, _ := procGetGuiResources.Call(hProc, uintptr(flag))
	return uint32(n)
}

func sendEscapeKey() {
	in := [2]input{
		{inputType: INPUT_KEY, ki: keyboardInput{wVk: VK_ESCAPE}},
		{inputType: INPUT_KEY, ki: keyboardInput{wVk: VK_ESCAPE, dwFlags: KEYEVENTF_KU}},
	}
	procSendInput.Call(uintptr(len(in)), uintptr(unsafe.Pointer(&in[0])), unsafe.Sizeof(in[0]))
}

//go:embed logo.png
var logo []byte

// buildMenu creates a non-trivial menu shape (many top-level items + a
// submenu) so every rebuild allocates multiple HMENU handles and a meaningful
// number of entries. The labels vary by iteration so the crash cannot be
// blamed on any single fixed string. When bitmaps is true, a handful of items
// receive a bitmap icon — exercising the SetMenuIcons path that leaks
// HBITMAPs (GR_GDIOBJECTS) on every rebuild.
func buildMenu(app *application.App, iter uint64, bitmaps bool) *application.Menu {
	m := app.NewMenu()
	header := m.Add(fmt.Sprintf("Iter #%d", iter)).SetEnabled(false)
	if bitmaps {
		header.SetBitmap(logo)
	}
	m.AddSeparator()
	for i := 0; i < 15; i++ {
		item := m.Add(fmt.Sprintf("Item %d-%d", iter, i))
		if bitmaps && i%3 == 0 {
			item.SetBitmap(logo)
		}
	}
	m.AddSeparator()
	m.AddCheckbox("Checkbox", iter%2 == 0)
	m.AddRadio("Radio A", true)
	m.AddRadio("Radio B", false)
	m.AddSeparator()
	sub := m.AddSubmenu("Submenu")
	for i := 0; i < 8; i++ {
		subItem := sub.Add(fmt.Sprintf("Sub %d-%d", iter, i))
		if bitmaps && i%2 == 0 {
			subItem.SetBitmap(logo)
		}
	}
	m.AddSeparator()
	m.Add("Quit").OnClick(func(ctx *application.Context) {
		os.Exit(0)
	})
	return m
}

// buildMutableMenu constructs a Menu whose items are returned separately
// so the caller can keep references to them across SetMenu rebuilds. The
// mutate workload uses these references to call SetBitmap after the
// initial build — that is the code path where windowsMenuItem.setBitmap
// runs with a non-nil impl and allocates an HBITMAP that is only tracked
// on impl.bitmap.
func buildMutableMenu(app *application.App, targetCount int) (*application.Menu, []*application.MenuItem) {
	m := app.NewMenu()
	m.Add("Mutable menu").SetEnabled(false)
	m.AddSeparator()
	targets := make([]*application.MenuItem, 0, targetCount)
	for i := 0; i < targetCount; i++ {
		targets = append(targets, m.Add(fmt.Sprintf("Mutable %d", i)))
	}
	m.AddSeparator()
	m.Add("Quit").OnClick(func(ctx *application.Context) {
		os.Exit(0)
	})
	return m, targets
}

type config struct {
	mode        string
	iters       int
	handleCap   uint
	churnGap    time.Duration
	showGap     time.Duration
	dismissGap  time.Duration
	logEvery    int
	runDuration time.Duration
	bitmaps     bool
}

func parseFlags() config {
	cfg := config{}
	flag.StringVar(&cfg.mode, "mode", "churn", "workload: churn | show | churn+show | mutate")
	flag.IntVar(&cfg.iters, "iters", 50000, "max SetMenu iterations before exiting cleanly (0 = no cap)")
	flag.UintVar(&cfg.handleCap, "handle-cap", 5000, "exit if GR_USEROBJECTS delta exceeds this")
	flag.DurationVar(&cfg.churnGap, "churn-gap", 2*time.Millisecond, "sleep between SetMenu calls")
	flag.DurationVar(&cfg.showGap, "show-gap", 80*time.Millisecond, "sleep between OpenMenu calls")
	flag.DurationVar(&cfg.dismissGap, "dismiss-gap", 30*time.Millisecond, "delay between popup open and ESC dismiss")
	flag.IntVar(&cfg.logEvery, "log-every", 500, "emit a progress line every N iterations")
	flag.DurationVar(&cfg.runDuration, "duration", 0, "wall-clock cap on the run (0 = unbounded)")
	flag.BoolVar(&cfg.bitmaps, "bitmaps", false, "attach a bitmap icon to a subset of menu items to exercise SetMenuIcons")
	flag.Parse()
	return cfg
}

func logEvent(event string, fields map[string]any) {
	// Plain event=... key=value log line (not JSON, despite the older name).
	// Keeps parsing trivial for the supervisor.
	out := fmt.Sprintf("event=%s", event)
	for k, v := range fields {
		out += fmt.Sprintf(" %s=%v", k, v)
	}
	fmt.Fprintln(os.Stderr, out)
}

func main() {
	runtime.LockOSThread()
	cfg := parseFlags()

	logEvent("start", map[string]any{
		"mode":        cfg.mode,
		"iters":       cfg.iters,
		"handle_cap":  cfg.handleCap,
		"pid":         os.Getpid(),
		"go_version":  runtime.Version(),
		"num_cpu":     runtime.NumCPU(),
		"churn_gap":   cfg.churnGap,
		"show_gap":    cfg.showGap,
		"dismiss_gap": cfg.dismissGap,
	})

	app := application.New(application.Options{
		Name:        "Systray Stress",
		Description: "Windows systray SetMenu stress test",
		Assets:      application.AlphaAssets,
	})

	tray := app.SystemTray.New()
	if len(logo) > 0 {
		tray.SetIcon(logo)
	}

	// Mutate mode owns its own Menu + item references, built once and reused.
	var mutateMenu *application.Menu
	var mutateTargets []*application.MenuItem
	if cfg.mode == "mutate" {
		mutateMenu, mutateTargets = buildMutableMenu(app, 5)
		tray.SetMenu(mutateMenu)
	} else {
		initial := buildMenu(app, 0, cfg.bitmaps)
		tray.SetMenu(initial)
	}

	baseHandles := getGuiResources(GR_USEROBJECTS)
	baseGDI := getGuiResources(GR_GDIOBJECTS)

	var iter uint64
	var exitCode int32 = -1 // set by the first terminating branch
	start := time.Now()
	exit := func(code int, reason string, extra map[string]any) {
		if !atomic.CompareAndSwapInt32(&exitCode, -1, int32(code)) {
			return
		}
		endUser := getGuiResources(GR_USEROBJECTS)
		endGDI := getGuiResources(GR_GDIOBJECTS)
		fields := map[string]any{
			"reason":        reason,
			"iter":          atomic.LoadUint64(&iter),
			"handles_start": baseHandles,
			"handles_end":   endUser,
			"handles_delta": int64(endUser) - int64(baseHandles),
			"gdi_start":     baseGDI,
			"gdi_end":       endGDI,
			"gdi_delta":     int64(endGDI) - int64(baseGDI),
			"runtime_ms":    time.Since(start).Milliseconds(),
		}
		for k, v := range extra {
			fields[k] = v
		}
		logEvent("exit", fields)
		os.Exit(code)
	}

	churn := func() {
		for {
			n := atomic.AddUint64(&iter, 1)
			tray.SetMenu(buildMenu(app, n, cfg.bitmaps))
			if cfg.logEvery > 0 && n%uint64(cfg.logEvery) == 0 {
				h := getGuiResources(GR_USEROBJECTS)
				g := getGuiResources(GR_GDIOBJECTS)
				logEvent("progress", map[string]any{
					"iter":          n,
					"handles":       h,
					"handles_delta": int64(h) - int64(baseHandles),
					"gdi":           g,
					"gdi_delta":     int64(g) - int64(baseGDI),
				})
				userOver := cfg.handleCap > 0 && h > baseHandles+uint32(cfg.handleCap)
				gdiOver := cfg.handleCap > 0 && g > baseGDI+uint32(cfg.handleCap)
				if userOver || gdiOver {
					kind := "user"
					if gdiOver {
						kind = "gdi"
					}
					exit(2, "handle_cap_exceeded", map[string]any{"handles": h, "gdi": g, "kind": kind})
				}
			}
			if cfg.iters > 0 && n >= uint64(cfg.iters) {
				exit(0, "iter_target_reached", nil)
			}
			if cfg.churnGap > 0 {
				time.Sleep(cfg.churnGap)
			}
		}
	}

	mutate := func() {
		for {
			n := atomic.AddUint64(&iter, 1)
			// Mutate bitmap on each target BEFORE rebuilding. Because the
			// items already have impls from the prior build, this dispatches
			// into windowsMenuItem.setBitmap, which allocates an HBITMAP
			// tracked only on impl.bitmap.
			for _, t := range mutateTargets {
				t.SetBitmap(logo)
			}
			// Rebuild the tray menu. updateMenu destroys the prior Win32Menu;
			// freeBitmaps walks menuMapping and releases each impl.bitmap
			// before DestroyMenu. This workload regression-guards that path:
			// if freeBitmaps ever stops walking menuMapping, GR_GDIOBJECTS
			// will climb one HBITMAP per iteration.
			tray.SetMenu(mutateMenu)

			if cfg.logEvery > 0 && n%uint64(cfg.logEvery) == 0 {
				h := getGuiResources(GR_USEROBJECTS)
				g := getGuiResources(GR_GDIOBJECTS)
				logEvent("progress", map[string]any{
					"iter":          n,
					"handles":       h,
					"handles_delta": int64(h) - int64(baseHandles),
					"gdi":           g,
					"gdi_delta":     int64(g) - int64(baseGDI),
				})
				userOver := cfg.handleCap > 0 && h > baseHandles+uint32(cfg.handleCap)
				gdiOver := cfg.handleCap > 0 && g > baseGDI+uint32(cfg.handleCap)
				if userOver || gdiOver {
					kind := "user"
					if gdiOver {
						kind = "gdi"
					}
					exit(2, "handle_cap_exceeded", map[string]any{"handles": h, "gdi": g, "kind": kind})
				}
			}
			if cfg.iters > 0 && n >= uint64(cfg.iters) {
				exit(0, "iter_target_reached", nil)
			}
			if cfg.churnGap > 0 {
				time.Sleep(cfg.churnGap)
			}
		}
	}

	show := func() {
		for atomic.LoadInt32(&exitCode) == -1 {
			// Fire a dismiss goroutine that presses ESC shortly after the popup appears.
			go func() {
				time.Sleep(cfg.dismissGap)
				sendEscapeKey()
			}()
			tray.OpenMenu() // blocks until popup dismissed
			if cfg.showGap > 0 {
				time.Sleep(cfg.showGap)
			}
		}
	}

	// Launch workloads once the app has fully started.
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		if cfg.runDuration > 0 {
			go func() {
				time.Sleep(cfg.runDuration)
				exit(0, "duration_reached", nil)
			}()
		}
		switch cfg.mode {
		case "churn":
			go churn()
		case "show":
			go show()
		case "churn+show":
			go churn()
			go show()
		case "mutate":
			go mutate()
		default:
			logEvent("fatal", map[string]any{"reason": "unknown_mode", "mode": cfg.mode})
			os.Exit(64)
		}
	})

	if err := app.Run(); err != nil {
		logEvent("fatal", map[string]any{"reason": "app_run_error", "err": err.Error()})
		log.Fatal(err)
	}
}
