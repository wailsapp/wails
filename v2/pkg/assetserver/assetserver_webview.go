package assetserver

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
)

type assetServerWebView struct {
	// ExpectedWebViewHost is checked against the Request Host of every WebViewRequest, other hosts won't be processed.
	ExpectedWebViewHost string

	dispatchInit    sync.Once
	dispatchReqC    chan<- webview.Request
	dispatchWorkers int
}

// ServeWebViewRequest processes the HTTP Request asynchronously by faking a golang HTTP Server.
// The request will be finished with a StatusNotImplemented code if no handler has written to the response.
// The AssetServer takes ownership of the request and the caller mustn't close it or access it in any other way.
func (d *AssetServer) ServeWebViewRequest(req webview.Request) {
	d.dispatchInit.Do(func() {
		workers := d.dispatchWorkers
		if workers == 0 {
			workers = 10
		}

		workerC := make(chan webview.Request, workers*2)
		for i := 0; i < workers; i++ {
			go func() {
				for req := range workerC {
					uri, _ := req.URL()
					d.processWebViewRequest(req)
					if err := req.Close(); err != nil {
						d.logError("Unable to call close for request for uri '%s'", uri)
					}
				}
			}()
		}

		dispatchC := make(chan webview.Request)
		go queueingDispatcher(50, dispatchC, workerC)

		d.dispatchReqC = dispatchC
	})

	d.dispatchReqC <- req
}

// processHTTPRequest processes the HTTP Request by faking a golang HTTP Server.
// The request will be finished with a StatusNotImplemented code if no handler has written to the response.
func (d *AssetServer) processWebViewRequest(r webview.Request) {
	wrw := r.Response()
	defer wrw.Finish()

	var rw http.ResponseWriter = &contentTypeSniffer{rw: wrw} // Make sure we have a Content-Type sniffer
	defer rw.WriteHeader(http.StatusNotImplemented)           // This is a NOP when a handler has already written and set the status

	uri, err := r.URL()
	if err != nil {
		d.logError("Error processing request, unable to get URL: %s (HttpResponse=500)", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	method, err := r.Method()
	if err != nil {
		d.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Method: %w", err))
		return
	}

	header, err := r.Header()
	if err != nil {
		d.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Header: %w", err))
		return
	}

	body, err := r.Body()
	if err != nil {
		d.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Body: %w", err))
		return
	}

	if body == nil {
		body = http.NoBody
	}
	defer body.Close()

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		d.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Request: %w", err))
		return
	}
	req.Header = header

	if req.RemoteAddr == "" {
		// 192.0.2.0/24 is "TEST-NET" in RFC 5737
		req.RemoteAddr = "192.0.2.1:1234"
	}

	if req.RequestURI == "" && req.URL != nil {
		req.RequestURI = req.URL.String()
	}

	if req.ContentLength == 0 {
		req.ContentLength, _ = strconv.ParseInt(req.Header.Get(HeaderContentLength), 10, 64)
	} else {
		req.Header.Set(HeaderContentLength, fmt.Sprintf("%d", req.ContentLength))
	}

	if host := req.Header.Get(HeaderHost); host != "" {
		req.Host = host
	}

	if expectedHost := d.ExpectedWebViewHost; expectedHost != "" && expectedHost != req.Host {
		d.webviewRequestErrorHandler(uri, rw, fmt.Errorf("expected host '%s' in request, but was '%s'", expectedHost, req.Host))
		return
	}

	d.ServeHTTP(rw, req)
}

func (d *AssetServer) webviewRequestErrorHandler(uri string, rw http.ResponseWriter, err error) {
	logInfo := uri
	if uri, err := url.ParseRequestURI(uri); err == nil {
		logInfo = strings.Replace(logInfo, fmt.Sprintf("%s://%s", uri.Scheme, uri.Host), "", 1)
	}

	d.logError("Error processing request '%s': %s (HttpResponse=500)", logInfo, err)
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

func queueingDispatcher[T any](minQueueSize uint, inC <-chan T, outC chan<- T) {
	q := newRingqueue[T](minQueueSize)
	for {
		in, ok := <-inC
		if !ok {
			return
		}

		q.Add(in)
		for q.Len() != 0 {
			out, _ := q.Peek()
			select {
			case outC <- out:
				q.Remove()
			case in, ok := <-inC:
				if !ok {
					return
				}

				q.Add(in)
			}
		}
	}
}
