package assetserver

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
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
func (a *AssetServer) ServeWebViewRequest(req webview.Request) {
	a.dispatchInit.Do(func() {
		workers := a.dispatchWorkers
		if workers <= 0 {
			return
		}

		workerC := make(chan webview.Request, workers*2)
		for i := 0; i < workers; i++ {
			go func() {
				for req := range workerC {
					a.processWebViewRequest(req)
				}
			}()
		}

		dispatchC := make(chan webview.Request)
		go queueingDispatcher(50, dispatchC, workerC)

		a.dispatchReqC = dispatchC
	})

	if a.dispatchReqC == nil {
		go a.processWebViewRequest(req)
	} else {
		a.dispatchReqC <- req
	}
}

func (a *AssetServer) processWebViewRequest(r webview.Request) {
	uri, _ := r.URL()
	a.processWebViewRequestInternal(r)
	if err := r.Close(); err != nil {
		a.options.Logger.Error("Unable to call close for request for uri.", "uri", uri)
	}
}

// processHTTPRequest processes the HTTP Request by faking a golang HTTP Server.
// The request will be finished with a StatusNotImplemented code if no handler has written to the response.
func (a *AssetServer) processWebViewRequestInternal(r webview.Request) {
	uri := "unknown"
	var err error

	wrw := r.Response()
	defer func() {
		if err := wrw.Finish(); err != nil {
			a.options.Logger.Error("Error finishing request '%s': %s", uri, err)
		}
	}()

	var rw http.ResponseWriter = &contentTypeSniffer{rw: wrw} // Make sure we have a Content-Type sniffer
	defer rw.WriteHeader(http.StatusNotImplemented)           // This is a NOP when a handler has already written and set the status

	uri, err = r.URL()
	if err != nil {
		a.options.Logger.Error(fmt.Sprintf("Error processing request, unable to get URL: %s (HttpResponse=500)", err))
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	method, err := r.Method()
	if err != nil {
		a.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Method: %w", err))
		return
	}

	header, err := r.Header()
	if err != nil {
		a.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Header: %w", err))
		return
	}

	body, err := r.Body()
	if err != nil {
		a.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Body: %w", err))
		return
	}

	if body == nil {
		body = http.NoBody
	}
	defer body.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		a.webviewRequestErrorHandler(uri, rw, fmt.Errorf("HTTP-Request: %w", err))
		return
	}

	// For server requests, the URL is parsed from the URI supplied on the Request-Line as stored in RequestURI. For
	// most requests, fields other than Path and RawQuery will be empty. (See RFC 7230, Section 5.3)
	req.URL.Scheme = ""
	req.URL.Host = ""
	req.URL.Fragment = ""
	req.URL.RawFragment = ""

	if requestURL := req.URL; req.RequestURI == "" && requestURL != nil {
		req.RequestURI = requestURL.String()
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

	if expectedHost := a.ExpectedWebViewHost; expectedHost != "" && expectedHost != req.Host {
		a.webviewRequestErrorHandler(uri, rw, fmt.Errorf("expected host '%s' in request, but was '%s'", expectedHost, req.Host))
		return
	}

	a.ServeHTTP(rw, req)
}

func (a *AssetServer) webviewRequestErrorHandler(uri string, rw http.ResponseWriter, err error) {
	logInfo := uri
	if uri, err := url.ParseRequestURI(uri); err == nil {
		logInfo = strings.Replace(logInfo, fmt.Sprintf("%s://%s", uri.Scheme, uri.Host), "", 1)
	}

	a.options.Logger.Error("Error processing request (HttpResponse=500)", "details", logInfo, "error", err.Error())
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
