package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/hash"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// The runtime scenario drives the app from its page, through the Wails
// JavaScript runtime, as most apps do: window methods, events, the
// clipboard, screens, system queries, bound calls, streams and dialogs. Each
// step runs in the page and reports back; Go answers the dialogs a step opens.

// RuntimeService is bound for the page's Call steps.
type RuntimeService struct{}

// Echo returns its argument.
func (*RuntimeService) Echo(s string) string { return s }

// Add returns a+b: arguments are converted from JSON.
func (*RuntimeService) Add(a, b int) int { return a + b }

// Fail returns an error, which the page sees as a rejected call.
func (*RuntimeService) Fail(msg string) error { return errors.New(msg) }

// callCancelled hears that Wait's context was cancelled. It is a channel of
// its own: the page's report of the step can arrive before or after it, and
// waiting for the report discards other signals.
var callCancelled = make(chan struct{}, 1)

// Wait returns when the page cancels the call.
func (*RuntimeService) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		select {
		case callCancelled <- struct{}{}:
		default:
		}
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return nil
	}
}

// services are bound when the app starts, and setups run once it exists.
var (
	services = []application.Service{application.NewService(&RuntimeService{})}
	setups   = []func(*probe){func(p *probe) {
		// Events: the page emits e2e-js, Go answers with e2e-echo.
		p.app.Event.On("e2e-js", func(e *application.CustomEvent) {
			p.app.Event.Emit("e2e-echo", e.Data)
		})
		// Streams: e2e-echo sends back each JSON value it receives.
		p.app.HandleStream("e2e-echo", func(c *application.StreamConn) {
			defer c.Close()
			for {
				var v any
				if err := c.ReceiveJSON(&v); err != nil {
					return
				}
				if err := c.SendJSON(v); err != nil {
					return
				}
			}
		})
	}}
)

// runtimePage runs the steps rt.run(name, arg) asks for and reports each
// result to /rt.
const runtimePage = `<!doctype html><html><head><script type="module" src="/wails/runtime.js"></script></head>
<body style="margin:0;font:16px system-ui,-apple-system,sans-serif;background:#0f172a;color:white">
<p style="margin:20px">E2E runtime</p>
<script type="module">
import {Application, Call, Clipboard, Dialogs, Events, JSONStream, Screens, System, Window} from '/wails/runtime.js';
const q = new URLSearchParams(location.search);
const w = q.get('w') || 'rt';
const states = q.get('states') === '1';
const os = q.get('os');
const report = (path, params) => fetch(path + '?' + new URLSearchParams({w, ...params}));
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const show = (v) => JSON.stringify(v);
function expect(what, ok, got) {
	if (!ok) throw new Error(what + ': got ' + show(got));
}
// until polls get() until ok(value) holds.
async function until(what, get, ok, ms = 5000) {
	const end = Date.now() + ms;
	let v = await get();
	while (!ok(v)) {
		if (Date.now() > end) throw new Error(what + ': got ' + show(v));
		await sleep(25);
		v = await get();
	}
}
function within(what, promise, ms = 5000) {
	return Promise.race([promise, sleep(ms).then(() => { throw new Error(what + ': timed out'); })]);
}
async function rejects(what, promise) {
	try {
		await promise;
	} catch (e) {
		return e;
	}
	throw new Error(what + ': resolved, want rejected');
}

const steps = {
	async 'window-geometry'() {
		expect('Name', (await Window.Name()) === 'rt', await Window.Name());
		const main = await Window.Get('main').Name();
		expect('Get("main").Name', main === 'main', main);
		await Window.SetTitle('E2E rt renamed');
		await Window.SetTitle('E2E rt');
		await Window.SetSize(520, 360);
		await until('SetSize(520, 360)', () => Window.Size(), (s) => s.width === 520 && s.height === 360);
		expect('Width', (await Window.Width()) === 520, await Window.Width());
		expect('Height', (await Window.Height()) === 360, await Window.Height());
		await Window.SetPosition(220, 230);
		await until('SetPosition(220, 230)', () => Window.Position(), (p) => p.x === 220 && p.y === 230);
		const rel = await Window.RelativePosition();
		await Window.SetRelativePosition(rel.x, rel.y);
		await Window.Center();
		await Window.SetMinSize(300, 200);
		await Window.SetMaxSize(800, 600);
		await Window.DisableSizeConstraints();
		await Window.EnableSizeConstraints();
		await Window.SetMinSize(0, 0);
		await Window.SetMaxSize(0, 0);
		await Window.SetResizable(false);
		expect('Resizable after SetResizable(false)', (await Window.Resizable()) === false, await Window.Resizable());
		await Window.SetResizable(true);
		expect('Resizable after SetResizable(true)', (await Window.Resizable()) === true, await Window.Resizable());
		await Window.SetAlwaysOnTop(true);
		await Window.SetAlwaysOnTop(false);
		await Window.SetBackgroundColour(15, 23, 42, 255);
		await Window.SetFullscreenButtonEnabled(false);
		await Window.SetFullscreenButtonEnabled(true);
		const screen = await Window.GetScreen();
		expect('GetScreen', screen && typeof screen.ID === 'string', screen);
	},
	async 'window-zoom'() {
		await Window.SetZoom(1.5);
		expect('GetZoom after SetZoom(1.5)', (await Window.GetZoom()) === 1.5, await Window.GetZoom());
		await Window.ZoomIn();
		await until('ZoomIn from 1.5', () => Window.GetZoom(), (z) => z > 1.5);
		await Window.ZoomOut();
		await Window.ZoomReset();
		await until('ZoomReset', () => Window.GetZoom(), (z) => z === 1);
	},
	async 'window-frame'() {
		await Window.SetFrameless(true);
		await Window.SetFrameless(false);
		await Window.ToggleFrameless();
		await Window.ToggleFrameless();
		await Window.Hide();
		await Window.Show();
		await Window.Focus();
		await until('Focus', () => Window.IsFocused(), (f) => f === true);
	},
	async 'window-states'() {
		// Checked where something grants the states (see windowStatesChange).
		// macOS animates full screen and ignores a change made during one.
		const pause = os === 'darwin' ? 1500 : 0;
		const step = async (what, act, get, ok) => {
			await act();
			if (states) await until(what, get, ok, 10000);
			if (pause) await sleep(pause);
		};
		const max = () => Window.IsMaximised();
		const full = () => Window.IsFullscreen();
		const min = () => Window.IsMinimised();
		await step('Maximise', () => Window.Maximise(), max, (v) => v);
		await step('UnMaximise', () => Window.UnMaximise(), max, (v) => !v);
		await step('ToggleMaximise', () => Window.ToggleMaximise(), max, (v) => v);
		await step('ToggleMaximise back', () => Window.ToggleMaximise(), max, (v) => !v);
		await step('Zoom', () => Window.Zoom(), async () => true, () => true);
		await step('Zoom back', () => Window.Zoom(), async () => true, () => true);
		await step('Fullscreen', () => Window.Fullscreen(), full, (v) => v);
		await step('UnFullscreen', () => Window.UnFullscreen(), full, (v) => !v);
		await step('ToggleFullscreen', () => Window.ToggleFullscreen(), full, (v) => v);
		await step('ToggleFullscreen back', () => Window.ToggleFullscreen(), full, (v) => !v);
		await step('Minimise', () => Window.Minimise(), min, (v) => v);
		await step('UnMinimise', () => Window.UnMinimise(), min, (v) => !v);
		await step('Restore', () => Window.Restore(), async () => (await max()) || (await full()) || (await min()), (v) => !v);
	},
	async events() {
		const echoed = new Promise((res) => Events.Once('e2e-echo', (e) => res(e.data)));
		await Events.Emit('e2e-js', {n: 1});
		const d = await within('the echo of e2e-js', echoed);
		expect('echoed data', d && d.n === 1, d);
		// OnMultiple stops after its count.
		let count = 0;
		Events.OnMultiple('e2e-echo', () => count++, 2);
		for (let i = 0; i < 3; i++) await Events.Emit('e2e-js', {n: i});
		await until('OnMultiple(2) after three echoes', async () => count, (c) => c === 2);
		await sleep(500);
		expect('OnMultiple(2) after three echoes', count === 2, count);
		// A cancelled listener and Off hear nothing.
		let heard = 0;
		const cancel = Events.On('e2e-echo', () => heard++);
		cancel();
		Events.On('e2e-echo', () => heard++);
		Events.Off('e2e-echo');
		await Events.Emit('e2e-js', {n: 9});
		await sleep(500);
		expect('listeners after cancel and Off', heard === 0, heard);
	},
	async clipboard() {
		await Clipboard.SetText('wails-rt-clip');
		await until('Clipboard.Text', () => Clipboard.Text(), (t) => t === 'wails-rt-clip');
	},
	async screens() {
		const all = await Screens.GetAll();
		expect('GetAll', Array.isArray(all) && all.length > 0, all);
		const primary = await Screens.GetPrimary();
		expect('GetPrimary', primary && primary.Size.Width > 0, primary);
		const current = await Screens.GetCurrent();
		expect('GetCurrent', current && typeof current.ID === 'string', current);
		const byID = await Screens.GetByID(primary.ID);
		expect('GetByID', byID && byID.ID === primary.ID, byID);
		const first = await Screens.GetByIndex(0);
		expect('GetByIndex(0)', first && typeof first.ID === 'string', first);
	},
	async system() {
		const dark = await System.IsDarkMode();
		expect('IsDarkMode', typeof dark === 'boolean', dark);
		const env = await System.Environment();
		expect('Environment().OS', env && env.OS === os, env);
		const caps = await System.Capabilities();
		expect('Capabilities', caps && typeof caps === 'object', caps);
		const flags = {linux: System.IsLinux(), darwin: System.IsMac(), windows: System.IsWindows()};
		expect('Is' + os, flags[os] === true, flags);
		expect('IsDesktop', System.IsDesktop() === true, System.IsDesktop());
	},
	async application() {
		await Application.Hide();
		await Application.Show();
		await Window.Show();
		await Window.Focus();
	},
	async call() {
		expect('ByName Echo', (await Call.ByName('main.RuntimeService.Echo', 'hi')) === 'hi', null);
		expect('ByID Echo', (await Call.ByID(Number(q.get('echo')), 'id')) === 'id', null);
		expect('ByName Add', (await Call.ByName('main.RuntimeService.Add', 2, 3)) === 5, null);
		const failed = await rejects('Fail', Call.ByName('main.RuntimeService.Fail', 'boom'));
		expect('Fail message', String(failed && failed.message).includes('boom'), String(failed));
		await rejects('an unknown method', Call.ByName('main.RuntimeService.Missing'));
		await rejects('wrong argument types', Call.ByName('main.RuntimeService.Add', 'a', 'b'));
		// Cancelling a running call cancels its context in Go.
		const waiting = Call.ByName('main.RuntimeService.Wait');
		await sleep(300);
		waiting.cancel();
		await rejects('a cancelled call', waiting);
	},
	async stream() {
		const s = JSONStream('e2e-echo');
		await within('stream open', new Promise((res, rej) => {
			s.addEventListener('open', res, {once: true});
			s.addEventListener('error', () => rej(new Error('stream error')), {once: true});
		}));
		const got = [];
		const two = new Promise((res) => {
			s.onmessage = (ev) => { got.push(ev.data); if (got.length === 2) res(); };
		});
		s.send({n: 1});
		s.send({n: 2, text: 'two'});
		await within('two echoed frames', two);
		expect('echoed frames, in order', got[0].n === 1 && got[1].n === 2 && got[1].text === 'two', got);
		const closed = new Promise((res) => s.addEventListener('close', res, {once: true}));
		s.close();
		await within('stream close', closed);
	},
	// The dialog steps return what the dialog returned; Go presses its keys.
	'dialog-info': () => Dialogs.Info({Title: 'E2E rt info', Message: 'Return presses the default button',
		Buttons: [{Label: 'Later'}, {Label: 'OK', IsDefault: true}]}),
	'dialog-question': () => Dialogs.Question({Title: 'E2E rt question', Message: 'Escape presses the cancel button',
		Detached: true, Buttons: [{Label: 'Yes', IsDefault: true}, {Label: 'No', IsCancel: true}]}),
	'dialog-warning': () => Dialogs.Warning({Title: 'E2E rt warning', Message: 'A warning', Detached: true}),
	'dialog-error': () => Dialogs.Error({Title: 'E2E rt error', Message: 'An error', Detached: true}),
	'dialog-open': (dir) => Dialogs.OpenFile({Title: 'E2E rt open', Message: 'Pick a file', ButtonText: 'Pick',
		Directory: dir, CanChooseFiles: true, AllowsMultipleSelection: true, ShowHiddenFiles: true,
		CanCreateDirectories: false, ResolvesAliases: true, Detached: true,
		Filters: [{DisplayName: 'Text', Pattern: '*.txt'}, {DisplayName: 'Everything', Pattern: '*'}]}),
	'dialog-save': (dir) => Dialogs.SaveFile({Title: 'E2E rt save', Message: 'Save it', ButtonText: 'Keep',
		Directory: dir, Filename: 'rt.txt', CanCreateDirectories: true, Detached: true,
		Filters: [{DisplayName: 'Text', Pattern: '*.txt'}]}),
	// These end the page: Go waits for it to load again, or to close.
	reload: () => { Window.Reload(); },
	'force-reload': () => { Window.ForceReload(); },
	close: () => { Window.Close(); },
};

window.rt = {
	async run(name, arg) {
		try {
			const v = await steps[name](arg);
			report('/rt', {step: name, r: 'ok', v: show(v ?? null)});
		} catch (e) {
			report('/rt', {step: name, r: 'fail', v: String((e && e.message) || e)});
		}
	},
};
report('/ready', {dpr: devicePixelRatio});
</script></body></html>`

// rtResult waits for the page's report of step, and returns its value.
func rtResult(step string, timeout time.Duration) (string, error) {
	prefix := "rt " + step + " "
	deadline := time.After(timeout)
	for {
		select {
		case got := <-signals:
			if !strings.HasPrefix(got, prefix) {
				continue
			}
			result := strings.TrimPrefix(got, prefix)
			if v, ok := strings.CutPrefix(result, "ok "); ok {
				return v, nil
			}
			return "", fmt.Errorf("%s: %s", step, strings.TrimPrefix(result, "fail "))
		case <-deadline:
			return "", fmt.Errorf("%s: no result within %s", step, timeout)
		}
	}
}

// runtimeAPI opens a window on the runtime page and runs its steps.
func runtimeAPI(p *probe) error {
	const rtTitle = "E2E rt"
	states := "0"
	if windowStatesChange() {
		states = "1"
	}
	url := fmt.Sprintf("/runtime?w=rt&states=%s&os=%s&echo=%d",
		states, runtime.GOOS, hash.Fnv("main.RuntimeService.Echo"))
	w := p.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "rt", Title: rtTitle, Width: 500, Height: 350, X: 250, Y: 150, URL: url,
	})
	closed := make(chan struct{}, 1)
	w.OnWindowEvent(events.Common.WindowClosing, func(*application.WindowEvent) {
		select {
		case closed <- struct{}{}:
		default:
		}
	})
	isClosed := false
	defer func() {
		if isClosed {
			return
		}
		w.Close()
		select {
		case <-closed:
		case <-time.After(5 * time.Second):
		}
	}()
	if err := await("ready rt", 15*time.Second); err != nil {
		return err
	}
	start := func(step string, arg any) {
		drain()
		js, _ := json.Marshal(arg)
		w.ExecJS(fmt.Sprintf("rt.run(%q, %s)", step, js))
	}
	run := func(step string) error {
		start(step, nil)
		_, err := rtResult(step, 30*time.Second)
		return err
	}

	for _, step := range []string{
		"window-geometry", "window-zoom", "window-frame", "window-states",
		"events", "clipboard", "screens", "system", "application",
	} {
		if err := run(step); err != nil {
			shot("runtime-" + step)
			return err
		}
	}
	select {
	case <-callCancelled:
	default:
	}
	if err := run("call"); err != nil {
		return err
	}
	select {
	case <-callCancelled:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("the cancelled call's context was not cancelled")
	}
	if err := run("stream"); err != nil {
		return err
	}

	// Message dialogs: Go presses the key a user would.
	for _, d := range []struct{ step, key, want string }{
		{"dialog-info", "return", `"OK"`},
		{"dialog-question", "escape", `"No"`},
		{"dialog-warning", "return", ""},
		{"dialog-error", "return", ""},
	} {
		start(d.step, nil)
		// Kept: as in dialog-default-button, the dialog cannot be detected,
		// so the probe waits for it to paint.
		time.Sleep(1500 * time.Millisecond)
		shot("runtime-" + d.step)
		if err := p.d.key(d.key); err != nil {
			return err
		}
		got, err := rtResult(d.step, 10*time.Second)
		if err != nil {
			return err
		}
		if d.want != "" && got != d.want {
			return fmt.Errorf("%s returned %s, want %s", d.step, got, d.want)
		}
	}
	if err := runtimeFileDialogs(p, start); err != nil {
		return err
	}

	// Reloads, then the page closes its own window.
	for _, step := range []string{"reload", "force-reload"} {
		drain()
		start(step, nil)
		if err := await("ready rt", 10*time.Second); err != nil {
			return fmt.Errorf("%s: %w", step, err)
		}
	}
	start("close", nil)
	select {
	case <-closed:
		isClosed = true
	case <-time.After(5 * time.Second):
		return fmt.Errorf("Window.Close from the page: the window did not close")
	}
	return nil
}

// runtimeFileDialogs opens the file dialogs from the page, and answers them
// as the open-file and save-file scenarios do.
func runtimeFileDialogs(p *probe, start func(string, any)) error {
	dir, err := os.MkdirTemp("", "wails-e2e-rt")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	want := filepath.Join(dir, "picked.txt")
	if err := os.WriteFile(want, []byte("e2e"), 0o644); err != nil {
		return err
	}
	_, cancels := p.d.(dialogCanceller)
	// A lane that skips open-file cannot type a path (Caps Lock changes the
	// letters): it cancels this dialog with Escape instead of choosing.
	escapes := !cancels && list("E2E_SKIP")["open-file"]

	// Open, several files allowed: the result is a list.
	start("dialog-open", dir)
	time.Sleep(2 * time.Second)
	shot("runtime-dialog-open")
	switch {
	case cancels:
		err = p.d.(dialogCanceller).cancelDialog()
	case escapes:
		if err = waitFocus(p.d, "E2E rt open"); err == nil {
			err = p.d.key("escape")
		}
	default:
		if c, ok := p.d.(fileChooser); ok {
			err = c.chooseFile("E2E rt open", want)
		} else {
			err = gtkChooseFile(p.d, "E2E rt open", want)
		}
	}
	if err != nil {
		return err
	}
	got, err := rtResult("dialog-open", 10*time.Second)
	if err != nil {
		return err
	}
	var paths []string
	if got != "null" {
		if err := json.Unmarshal([]byte(got), &paths); err != nil {
			return fmt.Errorf("dialog-open returned %s: %w", got, err)
		}
	}
	if cancels || escapes {
		if len(paths) != 0 {
			return fmt.Errorf("cancelled dialog-open returned %s", got)
		}
	} else if len(paths) != 1 || !samePath(paths[0], want) {
		return fmt.Errorf("dialog-open returned %s, want [%q]", got, want)
	}

	// Save: Return accepts the suggested name in the suggested folder.
	start("dialog-save", dir)
	time.Sleep(2 * time.Second)
	shot("runtime-dialog-save")
	if cancels {
		err = p.d.(dialogCanceller).cancelDialog()
	} else if err = waitFocus(p.d, "E2E rt save"); err == nil {
		err = p.d.key("return")
	}
	if err != nil {
		return err
	}
	got, err = rtResult("dialog-save", 10*time.Second)
	if err != nil {
		return err
	}
	var saved string
	if err := json.Unmarshal([]byte(got), &saved); err != nil && got != "null" {
		return fmt.Errorf("dialog-save returned %s: %w", got, err)
	}
	if cancels {
		if saved != "" {
			return fmt.Errorf("cancelled dialog-save returned %s", got)
		}
	} else if !samePath(saved, filepath.Join(dir, "rt.txt")) {
		return fmt.Errorf("dialog-save returned %s, want %q", got, filepath.Join(dir, "rt.txt"))
	}
	return nil
}

// samePath compares paths after resolving symbolic links in their folders
// (macOS's /var is /private/var).
func samePath(got, want string) bool {
	return mustEval(filepath.Dir(got)) == mustEval(filepath.Dir(want)) && filepath.Base(got) == filepath.Base(want)
}
