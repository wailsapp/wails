// updater-ui-driver-byo: BYO-window variant. Instead of letting the updater
// auto-create its builtin window, we create our own *application.WebviewWindow
// with a totally custom HTML template and pass it via Config.Window. The
// updater then drives that window through its full lifecycle (events
// in, user-action events back) without ever touching the framework's
// default template.
//
// This is the "Bring Your Own UI" path the PR description promises but the
// matrix had open.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

// A deliberately minimal custom UI that looks NOTHING like the builtin
// template — pink background, a single big heading, no buttons. The point
// is to prove the event channel is what's animating the page, not the
// builtin CSS classes.
const customHTML = `<!doctype html>
<html><head><meta charset="utf-8"><title>BYO Demo</title>
<style>
  html,body { margin:0; height:100%; }
  body {
    background: linear-gradient(135deg,#ff6f9c,#ffd56b);
    color:#1a1a1a;
    font: 600 22px/1.3 -apple-system, "Segoe UI Variable", "Segoe UI", sans-serif;
    display:flex; align-items:center; justify-content:center;
    text-align:center;
  }
  #card {
    background: rgba(255,255,255,.88);
    padding: 28px 36px;
    border-radius: 20px;
    box-shadow: 0 6px 30px rgba(0,0,0,.18);
    max-width: 80%;
  }
  #state { font-size: 14px; color:#666; text-transform: uppercase; letter-spacing: .15em; margin-bottom: 8px; }
  #version { font-size: 12px; color:#444; margin-top: 8px; font-weight: 400; }
  #notes { font-size: 11px; color:#555; margin-top: 12px; font-weight: 400; max-height: 100px; overflow:auto; text-align:left; white-space: pre-wrap; }
</style>
</head>
<body>
<div id="card">
  <div id="state">starting</div>
  <div id="title">BYO Updater</div>
  <div id="version">—</div>
  <div id="notes"></div>
</div>
<script>
(function () {
  var w = window._wails = window._wails || {};
  var L = Object.create(null);
  w.dispatchWailsEvent = function (ev) {
    if (!ev || !ev.name) return;
    var cbs = L[ev.name]; if (!cbs) return;
    for (var i = 0; i < cbs.length; i++) try { cbs[i](ev); } catch (_) {}
  };
  function on(n, cb) { (L[n] = L[n] || []).push(cb); }
  function emit(n) { if (typeof w.invoke === "function") w.invoke("wails:event:emit:" + n); }
  on("wails:updater:CheckStarted",    function () { setState("Checking…", "blue"); });
  on("wails:updater:UpdateAvailable", function (e) {
    var rel = e && (e.data != null ? e.data : e);
    setState("Update Available", "green");
    document.getElementById("title").textContent = "BYO Heard You";
    document.getElementById("version").textContent = "v" + (rel && rel.version || "?");
    if (rel && rel.notes) document.getElementById("notes").textContent = rel.notes;
  });
  on("wails:updater:DownloadStarted",  function () { setState("Downloading…", "blue"); });
  on("wails:updater:DownloadProgress", function (e) {
    var p = e && (e.data != null ? e.data : e);
    if (p && p.total) {
      var pct = Math.round((p.written / p.total) * 100);
      document.getElementById("state").textContent = "Downloading " + pct + "%";
    }
  });
  on("wails:updater:Verifying",  function () { setState("Verifying…", "blue"); });
  on("wails:updater:Installing", function () { setState("Installing…", "blue"); });
  on("wails:updater:UpdateReady", function () { setState("Ready to restart", "green"); });
  on("wails:updater:Error",      function (e) {
    var info = e && (e.data != null ? e.data : e);
    setState("Error", "red");
    document.getElementById("notes").textContent = (info && info.message) || "Unknown error";
  });
  function setState(text, _color) {
    document.getElementById("state").textContent = text;
  }
  // Wait for invoke, then signal runtime ready so queued events flush.
  (function announce() {
    if (window._wails && typeof window._wails.invoke === "function") {
      window._wails.invoke("wails:runtime:ready");
      emit("wails:updater:window:Ready");
    } else { setTimeout(announce, 30); }
  })();
})();
</script>
</body></html>`

func main() {
	app := application.New(application.Options{
		Name:        "Updater BYO Driver",
		Description: "Runs the updater against a user-supplied window with a custom HTML template.",
	})

	// User-owned window with custom HTML.
	myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "My Updater",
		Width:  520,
		Height: 460,
		HTML:   customHTML,
		// Required for the custom HTML's buttons to reach the host via the
		// `wails:event:emit:` postMessage path — see WebviewWindowOptions
		// GoDoc for the threat model.
		AllowSimpleEventEmit: true,
	})

	gh, err := github.New(github.Config{
		Repository:    envOr("GH_REPOSITORY", "wailsapp/updater-demo"),
		ChecksumAsset: "SHA256SUMS",
	})
	if err != nil {
		log.Fatalf("github.New: %v", err)
	}

	if err := app.Updater.Init(updater.Config{
		CurrentVersion: envOr("APP_VERSION", "1.0.0"),
		Providers:      []updater.Provider{gh},
		Window:         updater.BYOWindow(myWin.AsUpdaterWindow()),
	}); err != nil {
		log.Fatalf("Init: %v", err)
	}

	go func() {
		time.Sleep(800 * time.Millisecond)
		if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
			log.Printf("CheckAndInstall: %v", err)
		}
	}()

	if budget := os.Getenv("EXIT_AFTER_SECONDS"); budget != "" {
		go func() {
			d, _ := time.ParseDuration(budget + "s")
			if d == 0 {
				d = 8 * time.Second
			}
			time.Sleep(d)
			os.Exit(0)
		}()
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
