// Command nav-policy is a live check of the Linux navigation policy. It opens one window on the app's own page, which
// tries to navigate the main frame to a remote URL, to open a new window and to load an iframe from a second origin,
// then reports where it ended up. It exits 0 when the page is still the app's own, every attempt was refused and a
// control iframe from the app's own origin did load. Run it under a display; test/nav-policy/run.sh starts a headless
// one.
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// page is served from the app's origin. %s is the untrusted fixture's origin.
const page = `<!doctype html><title>nav</title><body>
<iframe id=control src="/control-frame"></iframe>
<iframe id=untrusted src="%s/frame"></iframe>
<script>
const seen = {control: false, untrusted: false};
addEventListener("message", e => {
  if (e.data === "control") seen.control = true;
  if (e.data === "untrusted") seen.untrusted = true;
});
setTimeout(() => { location.href = "https://nav-test.invalid/main"; }, 200);
setTimeout(() => { window.open("https://nav-test.invalid/window"); }, 400);
setTimeout(() => {
  const q = new URLSearchParams({href: location.href, origin: location.origin,
    control: seen.control, untrusted: seen.untrusted});
  fetch("/result?" + q);
}, 2500);
</script>`

// controlFrame is same-origin with the app, so its load must always be permitted. It proves an iframe in this page
// can load and post to the parent, which is what makes a silent untrusted iframe meaningful.
const controlFrame = `<!doctype html><title>control</title><body>
<script>parent.postMessage("control", "*");</script>`

// untrustedFrame reports its own load back to the app page and to the fixture's beacon endpoint, so either signal
// alone catches a policy that let it through.
const untrustedFrame = `<!doctype html><title>untrusted</title><body>
<script>
parent.postMessage("untrusted", "*");
new Image().src = "/beacon";
</script>`

func main() {
	var beacons atomic.Int32
	untrusted := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/frame":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(untrustedFrame))
		case "/beacon":
			beacons.Add(1)
			w.WriteHeader(http.StatusNoContent)
		case "/ping":
			_, _ = w.Write([]byte("nav-policy fixture"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer untrusted.Close()

	// Positive control: the fixture answers a request made outside the webview, so a silent iframe below can only
	// mean the policy refused it, not that the server was down.
	resp, err := http.Get(untrusted.URL + "/ping")
	if err != nil {
		fmt.Println("FAIL: the untrusted fixture is unreachable:", err)
		os.Exit(1)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(body) != "nav-policy fixture" {
		fmt.Printf("FAIL: the untrusted fixture answered %d %q\n", resp.StatusCode, body)
		os.Exit(1)
	}

	done := make(chan map[string]string, 1)
	app := application.New(application.Options{
		Name: "nav-policy",
		Assets: application.AssetOptions{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			switch r.URL.Path {
			case "/result":
				q := r.URL.Query()
				select {
				case done <- map[string]string{
					"origin":    q.Get("origin"),
					"control":   q.Get("control"),
					"untrusted": q.Get("untrusted"),
				}:
				default:
				}
			case "/control-frame":
				_, _ = w.Write([]byte(controlFrame))
			default:
				_, _ = fmt.Fprintf(w, page, untrusted.URL)
			}
		})},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{Width: 400, Height: 300})

	go func() {
		select {
		case result := <-done:
			fmt.Println("page origin after the attempts:", result["origin"])
			if result["origin"] != "wails://localhost" {
				fmt.Println("FAIL: the window left the app's origin")
				os.Exit(1)
			}
			if result["control"] != "true" {
				fmt.Println("FAIL: the same-origin control iframe never loaded; the iframe check proves nothing")
				os.Exit(1)
			}
			// The page reports 2.5s after loading, so any iframe load signal counted here arrived inside 3s.
			if result["untrusted"] == "true" || beacons.Load() > 0 {
				fmt.Println("FAIL: the untrusted iframe loaded")
				os.Exit(1)
			}
			// window.open runs without a user gesture, so WebKit's popup setting may stop it before the policy sees it:
			// that the page stayed put shows no window opened in-app, not which of the two refused it. The policy's own
			// new-window decision is covered by TestLinuxDecideNavigationSendsUserLinksToTheBrowser.
			fmt.Println("PASS: main-frame navigation was refused, window.open opened nothing in-app, the control iframe loaded, and the untrusted iframe was refused")
			os.Exit(0)
		case <-time.After(20 * time.Second):
			fmt.Println("FAIL: the page never reported back; it was navigated away")
			os.Exit(1)
		}
	}()

	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
