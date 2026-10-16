package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const indexHTML = `<!doctype html><html><head><meta charset="utf-8"><title>nested worker repro</title>
<script type="module" src="/wails/runtime.js"></script><style>html,body{background:#fff;color:#111}pre{white-space:pre-wrap;overflow-wrap:anywhere}</style></head>
<body style="font:14px monospace"><h3>Nested worker repro</h3>
<div id="phase">phase: loading</div>
<pre id="out">running...</pre>
<script>
window.__ready = true;
window.__results = null;
window.__control = null;
window.__phase = function(name, text) {
  document.getElementById('phase').textContent = 'phase: ' + name;
  document.getElementById('out').textContent = text;
};
window.__run = function() {
  window.__results = null;
  window.__phase('trigger', 'running nested worker tests...');
  var w = new Worker('/parent.js', { type: 'module' });
  w.onmessage = function(e) { window.__results = JSON.stringify(e.data); w.terminate(); };
  w.onerror = function(e) { window.__results = JSON.stringify({parentError: String(e.message)}); };
  setTimeout(function() { if (!window.__results) { window.__results = JSON.stringify({parentTimeout: true}); } }, 15000);
};
window.__runControl = function() {
  window.__phase('control', 'running top-level worker control...');
  window.__control = null;
  var w = new Worker('/child.js', { type: 'module' });
  w.onmessage = function() { window.__control = 'ok'; w.terminate(); };
  w.onerror = function() { window.__control = 'FAILED'; w.terminate(); };
  setTimeout(function() { if (!window.__control) { window.__control = 'timeout'; } }, 3000);
};
// Result channel: the page reports via same-origin relative fetch to the asset Handler.
window.__send = function(kind, value) {
  return fetch('/report?kind=' + encodeURIComponent(kind), { method: 'POST', body: String(value) });
};
window.__pollSend = function(kind, getter) {
  (function poll() {
    var v = getter();
    if (v) { window.__send(kind, v); } else { setTimeout(poll, 100); }
  })();
};
</script></body></html>`

const parentJS = `
const mk = (u, opts) => new Promise((res) => {
  let w;
  try { w = new Worker(u, opts); } catch (err) { res('THROW:' + String(err && err.name)); return; }
  w.onmessage = () => { w.terminate(); res('ok'); };
  w.onerror = (e) => { w.terminate(); res('FAILED(message=' + JSON.stringify(e.message) + ',isErrorEvent=' + (e instanceof ErrorEvent) + ')'); };
  setTimeout(() => { w.terminate(); res('timeout'); }, 3000);
});
(async () => {
  const abs = new URL('/child.js', self.location.href).href;
  const out = {};
  const r = await fetch(abs);
  out.fetchStatus = r.status;
  out.nestedUrlModule = await mk(abs, { type: 'module' });
  out.nestedUrlClassic = await mk(abs, {});
  out.nestedBlobImport = await mk(URL.createObjectURL(new Blob(['import ' + JSON.stringify(abs) + ';'], { type: 'text/javascript' })), { type: 'module' });
  out.nestedBlobInline = await mk(URL.createObjectURL(new Blob(["postMessage('hi')"], { type: 'text/javascript' })), {});
  postMessage(out);
})();
`

const childJS = `postMessage('child alive');`

var exitCode atomic.Int32

func main() {
	var finished int32

	controlCh := make(chan string, 4)
	resultsCh := make(chan string, 4)
	readyCh := make(chan struct{}, 4)

	mux := http.NewServeMux()
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	})
	mux.HandleFunc("/parent.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write([]byte(parentJS))
	})
	mux.HandleFunc("/child.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write([]byte(childJS))
	})
	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		// A single Read may return fewer bytes than the body holds, which would
		// truncate the report or drop a later field. Read the whole (bounded)
		// body to EOF and surface a read error instead of reporting a partial
		// result.
		raw, err := io.ReadAll(io.LimitReader(r.Body, 8192))
		if err != nil {
			http.Error(w, "reading report body", http.StatusBadRequest)
			return
		}
		body := string(raw)
		switch r.URL.Query().Get("kind") {
		case "ready":
			select {
			case readyCh <- struct{}{}:
			default:
			}
		case "control":
			select {
			case controlCh <- body:
			default:
			}
		case "results":
			select {
			case resultsCh <- body:
			default:
			}
		}
		w.WriteHeader(204)
	})

	app := application.New(application.Options{
		Name: "nested-worker-repro",
		Assets: application.AssetOptions{
			Handler: mux,
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "nested worker repro",
		Width:  700,
		Height: 300,
		URL:    "/index.html",
	})

	// Single outer watchdog.
	go func() {
		time.Sleep(45 * time.Second)
		if atomic.LoadInt32(&finished) == 0 {
			fmt.Println("REPRO_SETUP_ERROR: watchdog timeout")
			os.Exit(2)
		}
	}()

	setup := func(msg string) {
		fmt.Printf("REPRO_SETUP_ERROR: %s\n", msg)
		atomic.StoreInt32(&finished, 1)
		exitCode.Store(2)
		app.Quit()
	}

	go func() {
		// DOM readiness: page posts to /report?kind=ready once loaded.
		ready := false
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) && !ready {
			win.ExecJS(`if (window.__ready && window.__send) { window.__send('ready', '1'); }`)
			select {
			case <-readyCh:
				ready = true
			case <-time.After(500 * time.Millisecond):
			}
		}
		if !ready {
			setup("DOM not ready")
			return
		}

		// Positive control: top-level module worker from the document.
		win.ExecJS(`window.__runControl(); window.__pollSend('control', function(){ return window.__control; });`)
		var control string
		select {
		case control = <-controlCh:
		case <-time.After(10 * time.Second):
			setup("control timeout")
			return
		}
		win.ExecJS(fmt.Sprintf(`window.__phase('initial', 'top-level worker control: %s');`, control))
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("REPRO_OBSERVATION: phase=initial topLevelWorkerOk=%v\n", control == "ok")
		if control != "ok" {
			setup("positive control failed")
			return
		}
		fmt.Println("REPRO_READY: top-level worker control passed")

		// Trigger: nested workers created inside a module parent worker.
		win.ExecJS(`window.__run(); window.__pollSend('results', function(){ return window.__results; });`)
		var res string
		select {
		case res = <-resultsCh:
		case <-time.After(22 * time.Second):
			setup("results timeout")
			return
		}

		parentFailed := strings.Contains(res, "parentError") || strings.Contains(res, "parentTimeout")
		if parentFailed {
			setup("parent worker did not report")
			return
		}
		fetchOk := strings.Contains(res, `"fetchStatus":200`)
		urlModuleOk := strings.Contains(res, `"nestedUrlModule":"ok"`)
		urlClassicOk := strings.Contains(res, `"nestedUrlClassic":"ok"`)
		blobImportOk := strings.Contains(res, `"nestedBlobImport":"ok"`)
		blobInlineOk := strings.Contains(res, `"nestedBlobInline":"ok"`)
		// Reported signature: error event with undefined message, not an ErrorEvent.
		silentSig := strings.Contains(res, `isErrorEvent=false`)

		// Issue #6128 expects every recorded check to succeed, so any failing
		// check (the nested URL workers under test, or the fetch/blob controls)
		// is a failure and must drive the exit status.
		bug := !fetchOk || !urlModuleOk || !urlClassicOk || !blobImportOk || !blobInlineOk
		shown := fmt.Sprintf("fetchOk=%v urlModuleOk=%v urlClassicOk=%v blobImportOk=%v blobInlineOk=%v silentSignature=%v bugObserved=%v",
			fetchOk, urlModuleOk, urlClassicOk, blobImportOk, blobInlineOk, silentSig, bug)
		win.ExecJS(fmt.Sprintf(`window.__phase('trigger-done', %q);`, shown))
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("REPRO_OBSERVATION: phase=trigger fetchOk=%v nestedUrlModuleOk=%v nestedUrlClassicOk=%v nestedBlobImportOk=%v nestedBlobInlineOk=%v silentSignature=%v\n",
			fetchOk, urlModuleOk, urlClassicOk, blobImportOk, blobInlineOk, silentSig)

		// Reset check: top-level worker control still works.
		win.ExecJS(`window.__runControl(); window.__pollSend('control', function(){ return window.__control; });`)
		var control2 string
		select {
		case control2 = <-controlCh:
		case <-time.After(10 * time.Second):
			setup("reset control timeout")
			return
		}
		win.ExecJS(fmt.Sprintf(`window.__phase('reset', 'top-level worker control after trigger: %s; bugObserved=%v');`, control2, bug))
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("REPRO_OBSERVATION: phase=reset topLevelWorkerOk=%v bugObserved=%v\n", control2 == "ok", bug)

		if control2 != "ok" {
			setup("reset control failed")
			return
		}
		if bug {
			fmt.Printf("REPRO_ASSERTION: expected fetchOk=true nestedUrlModuleOk=true nestedUrlClassicOk=true nestedBlobImportOk=true nestedBlobInlineOk=true, actual fetchOk=%v nestedUrlModuleOk=%v nestedUrlClassicOk=%v nestedBlobImportOk=%v nestedBlobInlineOk=%v silentSignature=%v\n",
				fetchOk, urlModuleOk, urlClassicOk, blobImportOk, blobInlineOk, silentSig)
			fmt.Println("REPRO_DONE: fail")
			exitCode.Store(1)
		} else {
			fmt.Println("REPRO_DONE: pass")
			exitCode.Store(0)
		}
		atomic.StoreInt32(&finished, 1)
		app.Quit()
	}()

	if err := app.Run(); err != nil {
		fmt.Printf("REPRO_SETUP_ERROR: %v\n", err)
		os.Exit(2)
	}
	os.Exit(int(exitCode.Load()))
}
