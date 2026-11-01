// Native regression probe: build with desktop,production,webkit2_41 tags.
// Run on a Linux desktop. The fixed bridge prints only WITNESS:trusted-main.
// Vulnerable builds additionally print WITNESS:cross-origin-child.
// Uses localhost port 18992 and closes itself after six seconds.
package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"sync/atomic"
	"testing/fstest"
	"time"
)

type App struct {
	mainCalls  atomic.Int32
	childCalls atomic.Int32
}

func (a *App) Witness(value string) {
	fmt.Println("WITNESS:" + value)
	if value == "trusted-main" {
		a.mainCalls.Add(1)
	} else {
		a.childCalls.Add(1)
	}
}
func main() {
	app := &App{}
	http.HandleFunc("/child", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<script>window.webkit.messageHandlers.external.postMessage('C'+JSON.stringify({name:'main.App.Witness',args:['cross-origin-child'],callbackID:'child'}));</script>`)
	})
	go http.ListenAndServe("127.0.0.1:18992", nil)
	assets := fstest.MapFS{"index.html": {Data: []byte(`<html><body><script>setTimeout(()=>{window.webkit.messageHandlers.external.postMessage('C'+JSON.stringify({name:'main.App.Witness',args:['trusted-main'],callbackID:'main'}));let f=document.createElement('iframe');f.src='http://127.0.0.1:18992/child';document.body.appendChild(f)},1000)</script></body></html>`)}}
	err := wails.Run(&options.App{Title: "Wails bridge security probe", Width: 500, Height: 300, AssetServer: &assetserver.Options{Assets: assets}, Bind: []interface{}{app}, OnStartup: func(ctx context.Context) { go func() { time.Sleep(6 * time.Second); runtime.Quit(ctx) }() }})
	if err != nil {
		panic(err)
	}
	if app.mainCalls.Load() != 1 || app.childCalls.Load() != 0 {
		panic(fmt.Sprintf("bridge admission failed: main=%d child=%d", app.mainCalls.Load(), app.childCalls.Load()))
	}
}
