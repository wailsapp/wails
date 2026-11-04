// Native regression probe: build with desktop,production,webkit2_41 tags.
// Run on a Linux desktop. Trusted main-frame calls must succeed; cross-origin
// and sandboxed child-frame calls must not. The app closes after six seconds.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"testing/fstest"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	var childLoads atomic.Int32
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childLoads.Add(1)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<script>window.webkit.messageHandlers.external.postMessage('C'+JSON.stringify({name:'main.App.Witness',args:['child'],callbackID:'child'}));</script>`)
	})}
	defer server.Close()
	go server.Serve(listener)
	childURL, _ := json.Marshal("http://" + listener.Addr().String() + "/child")
	page := fmt.Sprintf(`<html><body><script>setTimeout(()=>{
 window.webkit.messageHandlers.external.postMessage('C'+JSON.stringify({name:'main.App.Witness',args:['trusted-main'],callbackID:'main'}));
 for(let i=0;i<2;i++){let f=document.createElement('iframe');if(i===1)f.setAttribute('sandbox','allow-scripts');f.src=%s;document.body.appendChild(f)}
 },1000)</script></body></html>`, childURL)
	assets := fstest.MapFS{"index.html": {Data: []byte(page)}}
	err = wails.Run(&options.App{
		Title: "Wails bridge security probe", Width: 500, Height: 300,
		AssetServer: &assetserver.Options{Assets: assets}, Bind: []interface{}{app},
		OnStartup: func(ctx context.Context) { go func() { time.Sleep(6 * time.Second); runtime.Quit(ctx) }() },
	})
	if err != nil {
		panic(err)
	}
	if app.mainCalls.Load() != 1 || app.childCalls.Load() != 0 || childLoads.Load() != 2 {
		panic(fmt.Sprintf("bridge admission failed: main=%d child=%d loaded=%d", app.mainCalls.Load(), app.childCalls.Load(), childLoads.Load()))
	}
	fmt.Println("PASS: trusted main frame admitted; both child frames loaded and blocked")
}
