//go:build windows

// Command kitchensink is an INTERACTIVE tour of the WebView2 API through the
// generated v2 bindings. It renders a control panel of buttons; clicking a
// button sends a command to Go over the WebView2 message bridge, Go invokes the
// corresponding (recent-leaning) WebView2 feature, and the result is shown back
// in the page's log. Several buttons have obvious visible effects — zoom,
// dark/light recolor, the real print dialog, a DevTools window, a Task Manager
// window, rasterization scaling, and live navigation — so you can see the
// bindings working, not just read a pass/fail list.
//
// Includes printing two ways: PrintToPdf (writes a PDF file) and ShowPrintUI
// (opens the browser print dialog).
//
// Run on an interactive desktop (console/RDP), not over SSH:
//
//	go run .
//
// Then click the buttons. The console also logs each action.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	webview2 "github.com/wailsapp/wails/v3/internal/webview2/pkg/webview2"
	"github.com/wailsapp/wails/v3/internal/webview2/webviewloader"
)

var (
	mainHWND      uintptr
	captureStream *webview2.IStream
	pdfPath       string
	shotPath      string

	// live COM objects (kept for the app lifetime)
	keepEnv        *webview2.ICoreWebView2Environment
	keepController *webview2.ICoreWebView2Controller
	keepWebView    *webview2.ICoreWebView2

	// recent versioned interfaces, acquired once at startup
	iC3       *webview2.ICoreWebView2Controller3
	iSettings2 *webview2.ICoreWebView2Settings2
	iProfile  *webview2.ICoreWebView2Profile
	iProfile2 *webview2.ICoreWebView2Profile2
	iCookies  *webview2.ICoreWebView2CookieManager
	iw6       *webview2.ICoreWebView2_6
	iw7       *webview2.ICoreWebView2_7
	iw16      *webview2.ICoreWebView2_16
	iw21      *webview2.ICoreWebView2_21

	// mutable UI state
	zoom     = 1.0
	raster   = 1.0
	schemeIx = 2 // start DARK
)

var schemes = []struct {
	name string
	val  webview2.COREWEBVIEW2_PREFERRED_COLOR_SCHEME
}{
	{"AUTO", webview2.COREWEBVIEW2_PREFERRED_COLOR_SCHEME_AUTO},
	{"LIGHT", webview2.COREWEBVIEW2_PREFERRED_COLOR_SCHEME_LIGHT},
	{"DARK", webview2.COREWEBVIEW2_PREFERRED_COLOR_SCHEME_DARK},
}

// ---- the control panel -------------------------------------------------------

const page = `<!doctype html>
<html><head><meta charset="utf-8"><title>WebView2 Kitchen Sink</title>
<style>
  :root{--bg:#0b1220;--fg:#e2e8f0;--card:#111c33;--edge:#1e293b;--accent:#3b82f6}
  @media (prefers-color-scheme: light){
    :root{--bg:#f1f5f9;--fg:#0f172a;--card:#ffffff;--edge:#cbd5e1;--accent:#2563eb}
  }
  body{font-family:Segoe UI,system-ui,sans-serif;margin:0;background:var(--bg);color:var(--fg)}
  .wrap{max-width:900px;margin:0 auto;padding:28px 32px}
  h1{font-size:24px;margin:0 0 2px}
  .sub{color:var(--accent);margin:0 0 20px;font-size:14px}
  .grid{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px}
  button{background:var(--card);color:var(--fg);border:1px solid var(--edge);
    border-radius:8px;padding:9px 13px;font-size:13px;cursor:pointer}
  button:hover{border-color:var(--accent)}
  .row{display:flex;gap:8px;margin:8px 0}
  input{flex:1;background:var(--card);color:var(--fg);border:1px solid var(--edge);
    border-radius:8px;padding:9px 12px;font-size:13px}
  #log{margin-top:16px;background:var(--card);border:1px solid var(--edge);border-radius:10px;
    padding:14px 16px;height:230px;overflow:auto;font-family:Consolas,monospace;font-size:13px;
    white-space:pre-wrap;line-height:1.5}
  .tag{color:var(--accent)}
</style></head>
<body><div class="wrap">
  <h1>WebView2 v2 — Kitchen Sink (interactive)</h1>
  <p class="sub">Click a button → Go calls a WebView2 feature via the generated bindings → result appears below.</p>

  <div class="grid">
    <button onclick="send('printpdf')">🖨 Print to PDF</button>
    <button onclick="send('printui')">🖨 Print dialog (ShowPrintUI)</button>
    <button onclick="send('zoom','+')">🔍 Zoom in</button>
    <button onclick="send('zoom','-')">🔍 Zoom out</button>
    <button onclick="send('scheme')">🌗 Toggle color scheme</button>
    <button onclick="send('raster')">🖥 Rasterization scale</button>
    <button onclick="send('devtools')">🛠 Open DevTools</button>
    <button onclick="send('taskmgr')">📊 Open Task Manager</button>
    <button onclick="send('version')">ℹ Browser version (DevTools protocol)</button>
    <button onclick="send('cookie')">🍪 Set cookie</button>
    <button onclick="send('clear')">🧹 Clear browsing data</button>
    <button onclick="send('ua')">🕵 Set User-Agent</button>
    <button onclick="send('screenshot')">📷 Screenshot to PNG</button>
  </div>

  <p class="sub" style="margin:4px 0">Read back what we set:</p>
  <div class="grid">
    <button onclick="send('getua')">📥 Get User-Agent</button>
    <button onclick="send('getcookies')">📥 Get cookies</button>
    <button onclick="send('getstate')">📥 Read state (zoom / raster / scheme)</button>
  </div>

  <div class="row">
    <input id="expr" value="1 + 2 + navigator.hardwareConcurrency">
    <button onclick="send('eval', document.getElementById('expr').value)">▶ ExecuteScriptWithResult</button>
  </div>
  <div class="row">
    <input id="url" value="https://wails.io">
    <button onclick="send('navigate', document.getElementById('url').value)">🌐 Navigate</button>
  </div>

  <div id="log"></div>
</div>
<script>
  var logEl = document.getElementById('log');
  function appendLog(line){
    var t = new Date().toLocaleTimeString();
    logEl.innerHTML = '<span class="tag">'+t+'</span>  '+line+'\n' + logEl.innerHTML;
  }
  function send(cmd, arg){
    appendLog('&rarr; '+cmd+(arg?(' ('+arg+')'):''));
    window.chrome.webview.postMessage(JSON.stringify({cmd:cmd, arg:arg||''}));
  }
  // Go replies via PostWebMessageAsJson({line:"..."}).
  window.chrome.webview.addEventListener('message', function(e){
    if (e.data && e.data.line) appendLog(e.data.line);
  });
  appendLog('control panel ready — click a button');
</script>
</body></html>`

func main() {
	runtime.LockOSThread()
	const coinitApartmentThreaded = 0x2
	if err := windows.CoInitializeEx(0, coinitApartmentThreaded); err != nil {
		log.Printf("CoInitializeEx: %v (continuing)", err)
	}
	defer windows.CoUninitialize()

	pdfPath, _ = filepath.Abs("kitchen.pdf")   // PrintToPdf needs an absolute path
	shotPath, _ = filepath.Abs("kitchen.png")

	mainHWND = createWindow("WebView2 v2 — kitchen sink", 960, 800)
	showWindow(mainHWND)
	log.Printf("kitchen sink: window open — click buttons in the panel")
	if err := webviewloader.CreateCoreWebView2Environment(&envHandler{}); err != nil {
		log.Fatalf("CreateCoreWebView2Environment: %v", err)
	}
	messageLoop()
	log.Printf("done")
}

// reply logs to the console and pushes a line into the page log via
// PostWebMessageAsJson (itself one of the demonstrated features).
func reply(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	log.Printf("  -> %s", line)
	if keepWebView != nil {
		b, _ := json.Marshal(map[string]string{"line": line})
		_ = keepWebView.PostWebMessageAsJson(string(b))
	}
}

// ---- bridge dispatch ---------------------------------------------------------

func dispatch(cmd, arg string) {
	switch cmd {
	case "printpdf":
		if iw7 == nil {
			reply("PrintToPdf unavailable (no ICoreWebView2_7)")
			return
		}
		if err := iw7.PrintToPdf(pdfPath, nil, webview2.NewICoreWebView2PrintToPdfCompletedHandler(&printPdfHandler{})); err != nil {
			reply("PrintToPdf error: %v", err)
		}
	case "printui":
		if iw16 == nil {
			reply("ShowPrintUI unavailable (no ICoreWebView2_16)")
			return
		}
		reply("opening print dialog…")
		if err := iw16.ShowPrintUI(webview2.COREWEBVIEW2_PRINT_DIALOG_KIND_BROWSER); err != nil {
			reply("ShowPrintUI error: %v", err)
		}
	case "zoom":
		if arg == "-" {
			zoom /= 1.25
		} else {
			zoom *= 1.25
		}
		if err := keepController.PutZoomFactor(zoom); err != nil {
			reply("PutZoomFactor error: %v", err)
		} else {
			reply("zoom = %.0f%%", zoom*100)
		}
	case "scheme":
		if iProfile == nil {
			reply("color scheme unavailable (no profile)")
			return
		}
		schemeIx = (schemeIx + 1) % len(schemes)
		s := schemes[schemeIx]
		if err := iProfile.PutPreferredColorScheme(s.val); err != nil {
			reply("PutPreferredColorScheme error: %v", err)
		} else {
			reply("preferred color scheme = %s (watch the page recolor)", s.name)
		}
	case "raster":
		if iC3 == nil {
			reply("rasterization scale unavailable (no Controller3)")
			return
		}
		switch {
		case raster >= 2.0:
			raster = 1.0
		default:
			raster += 0.5
		}
		if err := iC3.PutRasterizationScale(raster); err != nil {
			reply("PutRasterizationScale error: %v", err)
		} else {
			reply("rasterization scale = %.1fx", raster)
		}
	case "devtools":
		if err := keepWebView.OpenDevToolsWindow(); err != nil {
			reply("OpenDevToolsWindow error: %v", err)
		} else {
			reply("DevTools window opened")
		}
	case "taskmgr":
		if iw6 == nil {
			reply("Task Manager unavailable (no ICoreWebView2_6)")
			return
		}
		if err := iw6.OpenTaskManagerWindow(); err != nil {
			reply("OpenTaskManagerWindow error: %v", err)
		} else {
			reply("Task Manager window opened")
		}
	case "version":
		if err := keepWebView.CallDevToolsProtocolMethod("Browser.getVersion", "{}", webview2.NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(&devToolsHandler{})); err != nil {
			reply("CallDevToolsProtocolMethod error: %v", err)
		}
	case "eval":
		if iw21 == nil {
			reply("ExecuteScriptWithResult unavailable (no ICoreWebView2_21)")
			return
		}
		if err := iw21.ExecuteScriptWithResult(arg, webview2.NewICoreWebView2ExecuteScriptWithResultCompletedHandler(&execResultHandler{})); err != nil {
			reply("ExecuteScriptWithResult error: %v", err)
		}
	case "cookie":
		if iCookies == nil {
			reply("cookie manager unavailable")
			return
		}
		c, err := iCookies.CreateCookie("wv2demo", "hello", "localhost", "/")
		if err != nil {
			reply("CreateCookie error: %v", err)
			return
		}
		if err := iCookies.AddOrUpdateCookie(c); err != nil {
			reply("AddOrUpdateCookie error: %v", err)
		} else {
			reply("cookie set: wv2demo=hello (domain localhost)")
		}
	case "clear":
		if iProfile2 == nil {
			reply("clear browsing data unavailable (no Profile2)")
			return
		}
		if err := iProfile2.ClearBrowsingDataAll(webview2.NewICoreWebView2ClearBrowsingDataCompletedHandler(&clearDataHandler{})); err != nil {
			reply("ClearBrowsingDataAll error: %v", err)
		}
	case "ua":
		if iSettings2 == nil {
			reply("user-agent unavailable (no Settings2)")
			return
		}
		ua := "WailsWV2KitchenSink/1.0 (generated bindings)"
		if err := iSettings2.PutUserAgent(ua); err != nil {
			reply("PutUserAgent error: %v", err)
		} else {
			got, _ := iSettings2.GetUserAgent()
			reply("User-Agent set to: %s", got)
		}
	case "getua":
		if iSettings2 == nil {
			reply("user-agent unavailable (no Settings2)")
			return
		}
		ua, err := iSettings2.GetUserAgent()
		if err != nil {
			reply("GetUserAgent error: %v", err)
		} else {
			reply("current User-Agent: %s", ua)
		}
	case "getcookies":
		if iCookies == nil {
			reply("cookie manager unavailable")
			return
		}
		if err := iCookies.GetCookies("http://localhost/", webview2.NewICoreWebView2GetCookiesCompletedHandler(&getCookiesHandler{})); err != nil {
			reply("GetCookies error: %v", err)
		}
	case "getstate":
		z, _ := keepController.GetZoomFactor()
		line := fmt.Sprintf("zoom=%.0f%%", z*100)
		if iC3 != nil {
			if r, err := iC3.GetRasterizationScale(); err == nil {
				line += fmt.Sprintf(", raster=%.1fx", r)
			}
		}
		if iProfile != nil {
			if s, err := iProfile.GetPreferredColorScheme(); err == nil {
				line += ", scheme=" + schemeName(s)
			}
		}
		reply("state: %s", line)
	case "screenshot":
		capturePreview(shotPath)
	case "navigate":
		reply("navigating to %s (this replaces the panel — relaunch to return)", arg)
		if err := keepWebView.Navigate(arg); err != nil {
			reply("Navigate error: %v", err)
		}
	default:
		reply("unknown command: %s", cmd)
	}
}

func capturePreview(path string) {
	stream, err := createFileStream(path)
	if err != nil {
		reply("screenshot stream error: %v", err)
		return
	}
	captureStream = stream
	if err := keepWebView.CapturePreview(webview2.COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT_PNG, stream, webview2.NewICoreWebView2CapturePreviewCompletedHandler(&captureHandler{})); err != nil {
		reply("CapturePreview error: %v", err)
	}
}

// ---- handlers ----------------------------------------------------------------

type unknownStub struct{}

func (unknownStub) QueryInterface(refiid, object uintptr) uintptr { return 0x80004002 }
func (unknownStub) AddRef() uint32                                { return 1 }
func (unknownStub) Release() uint32                               { return 1 }

type envHandler struct{}

func (h *envHandler) EnvironmentCompleted(errorCode webviewloader.HRESULT, env *webviewloader.ICoreWebView2Environment) webviewloader.HRESULT {
	if errorCode != 0 || env == nil {
		log.Fatalf("environment creation failed: %s", hexHR(uint32(errorCode)))
	}
	keepEnv = (*webview2.ICoreWebView2Environment)(unsafe.Pointer(env))
	keepEnv.AddRef()
	if err := keepEnv.CreateCoreWebView2Controller(webview2.HWND(mainHWND), webview2.NewICoreWebView2CreateCoreWebView2ControllerCompletedHandler(&controllerHandler{})); err != nil {
		log.Fatalf("CreateCoreWebView2Controller: %v", err)
	}
	return 0
}

type controllerHandler struct{ unknownStub }

func (h *controllerHandler) CreateCoreWebView2ControllerCompleted(errorCode uintptr, controller *webview2.ICoreWebView2Controller) uintptr {
	if errorCode != 0 || controller == nil {
		log.Fatalf("controller creation failed: %s", hexHR(uint32(errorCode)))
	}
	keepController = controller
	keepController.AddRef()
	iC3, _ = keepController.GetICoreWebView2Controller3()

	wv, err := keepController.GetCoreWebView2()
	if err != nil {
		log.Fatalf("GetCoreWebView2: %v", err)
	}
	keepWebView = wv
	keepWebView.AddRef()

	// enable scripting + message bridge, then grab the recent interfaces.
	if s, err := keepWebView.GetSettings(); err == nil {
		_ = s.PutIsScriptEnabled(true)
		_ = s.PutIsWebMessageEnabled(true)
		iSettings2, _ = s.GetICoreWebView2Settings2()
	}
	iw6, _ = keepWebView.GetICoreWebView2_6()
	iw7, _ = keepWebView.GetICoreWebView2_7()
	iw16, _ = keepWebView.GetICoreWebView2_16()
	iw21, _ = keepWebView.GetICoreWebView2_21()
	if w2, err := keepWebView.GetICoreWebView2_2(); err == nil {
		iCookies, _ = w2.GetCookieManager()
	}
	if w13, err := keepWebView.GetICoreWebView2_13(); err == nil {
		if p, err := w13.GetProfile(); err == nil {
			iProfile = p
			iProfile2, _ = p.GetICoreWebView2Profile2()
			_ = p.PutPreferredColorScheme(schemes[schemeIx].val)
		}
	}

	var rc webview2.RECT
	getClientRect(mainHWND, &rc)
	_ = keepController.PutBounds(rc)
	_ = keepController.PutIsVisible(true)

	_, _ = keepWebView.AddWebMessageReceived(webview2.NewICoreWebView2WebMessageReceivedEventHandler(&webMsgHandler{}))
	if err := keepWebView.NavigateToString(page); err != nil {
		log.Fatalf("NavigateToString: %v", err)
	}
	return 0
}

type webMsgHandler struct{ unknownStub }

func (h *webMsgHandler) WebMessageReceived(sender *webview2.ICoreWebView2, args *webview2.ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	raw, err := args.TryGetWebMessageAsString()
	if err != nil {
		return 0
	}
	var m struct {
		Cmd string `json:"cmd"`
		Arg string `json:"arg"`
	}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		reply("bad message: %v", err)
		return 0
	}
	dispatch(m.Cmd, m.Arg)
	return 0
}

type printPdfHandler struct{ unknownStub }

func (h *printPdfHandler) PrintToPdfCompleted(errorCode uintptr, result bool) uintptr {
	if errorCode != 0 || !result {
		reply("PrintToPdf failed (hr=%s, ok=%v)", hexHR(uint32(errorCode)), result)
	} else {
		reply("PrintToPdf wrote %s — open it to see the printed page", pdfPath)
	}
	return 0
}

type devToolsHandler struct{ unknownStub }

func (h *devToolsHandler) CallDevToolsProtocolMethodCompleted(errorCode uintptr, result string) uintptr {
	if errorCode != 0 {
		reply("Browser.getVersion failed: %s", hexHR(uint32(errorCode)))
	} else {
		reply("Browser.getVersion → %s", truncate(result, 160))
	}
	return 0
}

type execResultHandler struct{ unknownStub }

func (h *execResultHandler) ExecuteScriptWithResultCompleted(errorCode uintptr, result *webview2.ICoreWebView2ExecuteScriptResult) uintptr {
	if errorCode != 0 || result == nil {
		reply("ExecuteScriptWithResult failed: %s", hexHR(uint32(errorCode)))
		return 0
	}
	if s, ok, err := result.TryGetResultAsString(); err == nil && ok {
		reply("eval result (string): %s", s)
	} else if j, err := result.GetResultAsJson(); err == nil {
		reply("eval result (json): %s", j)
	} else {
		reply("eval completed (no readable result)")
	}
	return 0
}

type clearDataHandler struct{ unknownStub }

func (h *clearDataHandler) ClearBrowsingDataCompleted(errorCode uintptr) uintptr {
	reply("ClearBrowsingDataAll done (hr=%s)", hexHR(uint32(errorCode)))
	return 0
}

type getCookiesHandler struct{ unknownStub }

func (h *getCookiesHandler) GetCookiesCompleted(errorCode uintptr, list *webview2.ICoreWebView2CookieList) uintptr {
	if errorCode != 0 || list == nil {
		reply("GetCookies failed: %s", hexHR(uint32(errorCode)))
		return 0
	}
	n, _ := list.GetCount()
	if n == 0 {
		reply("no cookies for http://localhost/")
		return 0
	}
	out := ""
	for i := uint32(0); i < n; i++ {
		c, err := list.GetValueAtIndex(i)
		if err != nil {
			continue
		}
		name, _ := c.GetName()
		val, _ := c.GetValue()
		if out != "" {
			out += "; "
		}
		out += name + "=" + val
	}
	reply("cookies for http://localhost/ (%d): %s", n, out)
	return 0
}

type captureHandler struct{ unknownStub }

func (h *captureHandler) CapturePreviewCompleted(errorCode uintptr) uintptr {
	if captureStream != nil {
		_ = captureStream.Release()
	}
	if errorCode != 0 {
		reply("CapturePreview failed: %s", hexHR(uint32(errorCode)))
	} else {
		reply("screenshot written: %s", shotPath)
	}
	return 0
}

// ---- helpers -----------------------------------------------------------------

func schemeName(s webview2.COREWEBVIEW2_PREFERRED_COLOR_SCHEME) string {
	if int(s) >= 0 && int(s) < len(schemes) {
		return schemes[s].name
	}
	return fmt.Sprintf("%d", int(s))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func hexHR(hr uint32) string {
	const digits = "0123456789abcdef"
	b := []byte("0x00000000")
	for i := 0; i < 8; i++ {
		b[2+7-i] = digits[hr&0xf]
		hr >>= 4
	}
	return string(b)
}

// ---- Win32 plumbing ----------------------------------------------------------

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shlwapi  = windows.NewLazySystemDLL("shlwapi.dll")

	pRegisterClassExW   = user32.NewProc("RegisterClassExW")
	pCreateWindowExW    = user32.NewProc("CreateWindowExW")
	pDefWindowProcW     = user32.NewProc("DefWindowProcW")
	pShowWindow         = user32.NewProc("ShowWindow")
	pUpdateWindow       = user32.NewProc("UpdateWindow")
	pGetMessageW        = user32.NewProc("GetMessageW")
	pTranslateMessage   = user32.NewProc("TranslateMessage")
	pDispatchMessageW   = user32.NewProc("DispatchMessageW")
	pPostQuitMessage    = user32.NewProc("PostQuitMessage")
	pGetClientRect      = user32.NewProc("GetClientRect")
	pGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	pSHCreateStreamFile = shlwapi.NewProc("SHCreateStreamOnFileEx")
)

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type msgW struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
}

const (
	wmDestroy          = 0x0002
	wsOverlappedWindow = 0x00CF0000
	cwUseDefault       = ^uintptr(0x7FFFFFFF) // 0x80000000
	swShow             = 5
)

func wndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	if msg == wmDestroy {
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func createWindow(title string, width, height int) uintptr {
	hInstance, _, _ := pGetModuleHandleW.Call(0)
	className := windows.StringToUTF16Ptr("WebView2KitchenSink")

	wc := wndClassExW{
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInstance,
		lpszClassName: className,
		hbrBackground: 6,
	}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	if atom, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		log.Fatalf("RegisterClassExW: %v", err)
	}
	hwnd, _, err := pCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))),
		wsOverlappedWindow,
		cwUseDefault, cwUseDefault, uintptr(width), uintptr(height),
		0, 0, hInstance, 0,
	)
	if hwnd == 0 {
		log.Fatalf("CreateWindowExW: %v", err)
	}
	return hwnd
}

func showWindow(hwnd uintptr) {
	pShowWindow.Call(hwnd, swShow)
	pUpdateWindow.Call(hwnd)
}

func getClientRect(hwnd uintptr, rc *webview2.RECT) {
	pGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(rc)))
}

func messageLoop() {
	var msg msgW
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 {
			return
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func createFileStream(path string) (*webview2.IStream, error) {
	const (
		stgmReadWrite = 0x00000002
		stgmCreate    = 0x00001000
		fileAttrNorm  = 0x00000080
	)
	var stream uintptr
	hr, _, _ := pSHCreateStreamFile.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		stgmReadWrite|stgmCreate,
		fileAttrNorm,
		1,
		0,
		uintptr(unsafe.Pointer(&stream)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return (*webview2.IStream)(unsafe.Pointer(stream)), nil
}
