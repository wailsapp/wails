//go:build windows

// Native regression probe: build and run in an interactive Windows session.
// It requires no debugging port or external CDP client.
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed page.html
var page string

type result struct {
	Name, Outcome string
	Milliseconds  int64
}
type probe struct {
	mu          sync.Mutex
	inflight    sync.WaitGroup
	results     map[string]result
	closing     *application.WebviewWindow
	externalURL string
}

func (p *probe) record(name, outcome string, elapsed int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.results[name] = result{name, outcome, elapsed}
}
func privateMarker(r *http.Request) bool {
	return r.Header.Get("X-Wails-Request-Id") != "" || strings.Contains(r.URL.String(), "__wails_keepalive")
}
func (p *probe) slow(w http.ResponseWriter, r *http.Request) {
	p.inflight.Add(1)
	defer p.inflight.Done()
	name := r.URL.Query().Get("name")
	if r.Header.Get("X-Case") != "" {
		name = r.Header.Get("X-Case")
	}
	start := time.Now()
	outcome := "complete"
	if privateMarker(r) {
		outcome = "marker leaked"
	}
	if name == "post" || name == "keepalive-post" {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "upload body" || r.Header.Get("X-Probe") != "preserved" {
			outcome = "body/header changed"
		}
	}
	if name == "close" {
		go func() { time.Sleep(700 * time.Millisecond); p.closing.Close() }()
	}
	duration := 4 * time.Second
	if name == "app-shutdown" {
		duration = 30 * time.Second
	}
	select {
	case <-r.Context().Done():
		if outcome == "complete" {
			outcome = "cancelled"
		}
	case <-time.After(duration):
	}
	p.record(name, outcome, time.Since(start).Milliseconds())
	fmt.Fprint(w, "done")
}
func (p *probe) external(w http.ResponseWriter, r *http.Request) {
	outcome := "complete"
	if privateMarker(r) {
		outcome = "marker leaked"
	}
	if r.Method == "OPTIONS" {
		outcome = "unexpected preflight"
	}
	p.record(r.URL.Query().Get("name"), outcome, 0)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, "external response")
}
func (p *probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/slow":
		p.slow(w, r)
	case "/external":
		p.external(w, r)
	case "/redirect":
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "redirect"
		}
		http.Redirect(w, r, "/slow?name="+name, http.StatusTemporaryRedirect)
	case "/external-redirect":
		http.Redirect(w, r, p.externalURL+"/external?name="+r.URL.Query().Get("name"), http.StatusTemporaryRedirect)
	case "/client":
		p.record(r.URL.Query().Get("name"), r.URL.Query().Get("outcome"), 0)
	case "/worker.js":
		w.Header().Set("Content-Type", "text/javascript")
		fmt.Fprint(w, `const name=new URL(location).searchParams.get('name')||'worker';const c=new AbortController();fetch('/slow?name='+name,{signal:c.signal,keepalive:name.includes('keepalive')}).catch(()=>{});if(!name.endsWith('terminated'))setTimeout(()=>c.abort(),700);`)
	case "/done":
		fmt.Fprint(w, "navigated")
	default:
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, page)
	}
}
func (p *probe) validate() error {
	expected := map[string][]string{
		"cancelled": {"fetch", "xhr", "post", "redirect", "same-abort", "frame", "navigation", "close", "worker", "worker-keepalive", "keepalive-abort", "keepalive-redirect-abort", "worker-terminated", "keepalive-post", "app-shutdown"},
		"complete":  {"same-control", "normal", "keepalive", "http-keepalive", "keepalive-redirect", "external", "external-client", "external-keepalive", "external-keepalive-client", "keepalive-url", "worker-keepalive-terminated"},
	}
	for i := range 20 {
		for _, prefix := range []string{"stress-", "stress-keepalive-"} {
			expected["cancelled"] = append(expected["cancelled"], fmt.Sprintf("%s%d", prefix, i))
		}
	}
	for outcome, names := range expected {
		for _, name := range names {
			if p.results[name].Outcome != outcome {
				return fmt.Errorf("%s: got %q, want %s", name, p.results[name].Outcome, outcome)
			}
		}
	}
	return nil
}
func (p *probe) windows(app *application.App) {
	app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "HTTP keepalive", URL: p.externalURL + "/?scenario=http-keepalive", Hidden: true})
	for _, scenario := range []string{"main", "navigation", "keepalive", "close"} {
		window := app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Cancellation probe: " + scenario, URL: "/?scenario=" + scenario, Width: 400, Height: 200, Hidden: true})
		if scenario == "close" {
			p.closing = window
		}
	}
}
func run() error {
	logs, err := os.Create("native.log")
	if err != nil {
		return err
	}
	defer logs.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer listener.Close()
	p := &probe{results: make(map[string]result), externalURL: "http://" + listener.Addr().String()}
	go http.Serve(listener, p)
	app := application.New(application.Options{Name: "RequestCancellationRegression", Logger: slog.New(slog.NewTextHandler(logs, nil)), Windows: application.WindowsOptions{DisableQuitOnLastWindowClosed: true}, Assets: application.AssetOptions{Handler: p}})
	p.windows(app)
	go func() { time.Sleep(15 * time.Second); app.Quit() }()
	if err := app.Run(); err != nil {
		return err
	}
	p.inflight.Wait()
	p.mu.Lock()
	defer p.mu.Unlock()
	data, err := json.MarshalIndent(p.results, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("request-cancellation-results.json", data, 0600); err != nil {
		return err
	}
	return p.validate()
}
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
