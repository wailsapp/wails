// Command nav-policy is a live check of the Linux navigation policy. It opens one window on the app's own page, which
// tries to navigate the main frame to a remote URL, to open a new window and to load a remote iframe, then reports
// where it ended up. It exits 0 when the page is still the app's own and every attempt was refused. Run it under a
// display; test/nav-policy/run.sh starts a headless one.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const page = `<!doctype html><title>nav</title><body>
<iframe id=f src="https://nav-test.invalid/frame"></iframe>
<script>
setTimeout(() => { location.href = "https://nav-test.invalid/main"; }, 200);
setTimeout(() => { window.open("https://nav-test.invalid/window"); }, 400);
setTimeout(() => {
  const q = new URLSearchParams({href: location.href, origin: location.origin});
  fetch("/result?" + q);
}, 2500);
</script>`

func main() {
	done := make(chan string, 1)
	app := application.New(application.Options{
		Name: "nav-policy",
		Assets: application.AssetOptions{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/result" {
				select {
				case done <- r.URL.Query().Get("origin"):
				default:
				}
				return
			}
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(page))
		})},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{Width: 400, Height: 300})

	go func() {
		select {
		case origin := <-done:
			fmt.Println("page origin after the attempts:", origin)
			if origin != "wails://localhost" {
				fmt.Println("FAIL: the window left the app's origin")
				os.Exit(1)
			}
			fmt.Println("PASS: main-frame navigation, window.open and the remote iframe were refused")
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
