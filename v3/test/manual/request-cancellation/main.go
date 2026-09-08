//go:build linux && cgo && !android

// Native regression probe for #5963. Each scenario exits nonzero on failure.
package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	mode := "abort"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	aborted, cancelled, completed, bad := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	var a, c, d, b sync.Once
	var closeWindow func()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slow":
			id := r.Header.Get("X-Probe")
			fmt.Println("[request-cancellation] handler started", id, mode)
			if id == "B" {
				select {
				case <-r.Context().Done():
					fmt.Println("[request-cancellation] FAIL non-aborted request cancelled")
					b.Do(func() { close(bad) })
				case <-time.After(1200 * time.Millisecond):
					w.Write([]byte("ok"))
					fmt.Println("[request-cancellation] non-aborted request completed")
				}
				return
			}
			if mode == "stream" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(200)
				w.Write([]byte("started\n"))
				w.(http.Flusher).Flush()
				fmt.Println("[request-cancellation] response stream started")
			}
			select {
			case <-r.Context().Done():
				fmt.Println("[request-cancellation] aborted context cancelled")
				c.Do(func() { close(cancelled) })
			case <-time.After(3 * time.Second):
				fmt.Println("[request-cancellation] FAIL aborted handler reached timeout")
			}
		case "/next":
			a.Do(func() { close(aborted) })
			w.Write([]byte("navigated"))
		case "/close":
			w.Write([]byte("closing"))
			go func() { time.Sleep(100 * time.Millisecond); closeWindow(); a.Do(func() { close(aborted) }) }()
		case "/aborted":
			a.Do(func() { close(aborted) })
			w.Write([]byte("ok"))
		case "/completed":
			d.Do(func() { close(completed) })
			w.Write([]byte("ok"))
		case "/":
			js := ""
			if mode != "normal" {
				js += `const c=new AbortController();fetch('/slow',{headers:{'X-Probe':'A'},signal:c.signal}).then(r=>r.text()).catch(e=>{if(e.name==='AbortError')fetch('/aborted')});setTimeout(()=>c.abort(),500);`
			}
			if mode == "normal" || mode == "parallel" {
				js += `fetch('/slow',{headers:{'X-Probe':'B'}}).then(r=>r.text()).then(t=>{if(t==='ok')fetch('/completed')});`
			}
			if mode == "navigate" {
				js = `fetch('/slow',{headers:{'X-Probe':'A'}}).catch(()=>{});setTimeout(()=>location.replace('/next'),500);`
			}
			if mode == "close" {
				js = `fetch('/slow',{headers:{'X-Probe':'A'}}).catch(()=>{});setTimeout(()=>fetch('/close'),500);`
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body>5963 %s<script>%s</script></body></html>", mode, js)
		default:
			http.NotFound(w, r)
		}
	})
	app := application.New(application.Options{Name: "Wails5963Matrix", Linux: application.LinuxOptions{DisableQuitOnLastWindowClosed: true}, Assets: application.AssetOptions{Handler: handler}})
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "5963 Linux probe", Width: 400, Height: 150, URL: "/"})
	closeWindow = window.Close
	result := make(chan int, 1)
	go func() {
		code := 2
		defer func() { result <- code; app.Quit() }()
		if mode != "normal" {
			select {
			case <-aborted:
				fmt.Println("[request-cancellation] abort/teardown trigger observed")
			case <-time.After(10 * time.Second):
				fmt.Println("[request-cancellation] ERROR browser did not report abort")
				return
			}
			select {
			case <-cancelled:
			case <-time.After(time.Second):
				fmt.Println("[request-cancellation] FAIL context still active after browser abort")
				code = 1
				return
			}
		}
		if mode == "normal" || mode == "parallel" {
			select {
			case <-bad:
				code = 1
				return
			case <-completed:
			case <-time.After(10 * time.Second):
				fmt.Println("[request-cancellation] ERROR normal request did not complete")
				return
			}
		}
		fmt.Println("[request-cancellation] PASS", mode)
		code = 0
	}()
	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	os.Exit(<-result)
}
