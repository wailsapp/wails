//go:build windows
// +build windows

package webview

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wailsapp/go-webview2/pkg/edge"
)

// NewRequest creates as new WebViewRequest for chromium. This Method must be called from the Main-Thread!
func NewRequest(env *edge.ICoreWebView2Environment, args *edge.ICoreWebView2WebResourceRequestedEventArgs, invokeSync func(fn func())) (Request, error) {
	req, err := args.GetRequest()
	if err != nil {
		return nil, fmt.Errorf("GetRequest failed: %s", err)
	}
	defer req.Release()

	r := &request{
		invokeSync: invokeSync,
	}

	code := http.StatusInternalServerError
	r.response, err = env.CreateWebResourceResponse(nil, code, http.StatusText(code), "")
	if err != nil {
		return nil, fmt.Errorf("CreateWebResourceResponse failed: %s", err)
	}

	if err := args.PutResponse(r.response); err != nil {
		r.finishResponse()
		return nil, fmt.Errorf("PutResponse failed: %s", err)
	}

	r.deferral, err = args.GetDeferral()
	if err != nil {
		r.finishResponse()
		return nil, fmt.Errorf("GetDeferral failed: %s", err)
	}

	r.url, r.urlErr = req.GetUri()
	r.method, r.methodErr = req.GetMethod()
	r.header, r.headerErr = getHeaders(req)

	if content, err := req.GetContent(); err != nil {
		r.bodyErr = err
	} else if content != nil {
		// It is safe to access Content from another Thread: https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/threading-model#thread-safety
		r.body = &iStreamReleaseCloser{stream: content}
	}

	return r, nil
}

var _ Request = &request{}

type request struct {
	response *edge.ICoreWebView2WebResourceResponse
	deferral *edge.ICoreWebView2Deferral

	url    string
	urlErr error

	method    string
	methodErr error

	header    http.Header
	headerErr error

	body    io.ReadCloser
	bodyErr error
	rw      *responseWriter

	invokeSync func(fn func())
}

func (r *request) URL() (string, error) {
	return r.url, r.urlErr
}

func (r *request) Method() (string, error) {
	return r.method, r.methodErr
}

func (r *request) Header() (http.Header, error) {
	return r.header, r.headerErr
}

func (r *request) Body() (io.ReadCloser, error) {
	return r.body, r.bodyErr
}

func (r *request) Response() ResponseWriter {
	if r.rw != nil {
		return r.rw
	}

	r.rw = &responseWriter{req: r}
	return r.rw
}

func (r *request) Close() error {
	var errs []error
	if r.body != nil {
		if err := r.body.Close(); err != nil {
			errs = append(errs, err)
		}
		r.body = nil
	}

	if err := r.Response().Finish(); err != nil {
		errs = append(errs, err)
	}

	return combineErrs(errs)
}

// finishResponse must be called on the main-thread
func (r *request) finishResponse() error {
	var errs []error
	if r.response != nil {
		if err := r.response.Release(); err != nil {
			errs = append(errs, err)
		}
		r.response = nil
	}
	if r.deferral != nil {
		if err := r.deferral.Complete(); err != nil {
			errs = append(errs, err)
		}

		if err := r.deferral.Release(); err != nil {
			errs = append(errs, err)
		}
		r.deferral = nil
	}
	return combineErrs(errs)
}

type iStreamReleaseCloser struct {
	stream *edge.IStream
	closed bool
}

func (i *iStreamReleaseCloser) Read(p []byte) (int, error) {
	if i.closed {
		return 0, io.ErrClosedPipe
	}
	return i.stream.Read(p)
}

func (i *iStreamReleaseCloser) Close() error {
	if i.closed {
		return nil
	}
	i.closed = true
	return i.stream.Release()
}

func getHeaders(req *edge.ICoreWebView2WebResourceRequest) (http.Header, error) {
	header := http.Header{}
	headers, err := req.GetHeaders()
	if err != nil {
		return nil, fmt.Errorf("GetHeaders Error: %s", err)
	}
	defer headers.Release()

	headersIt, err := headers.GetIterator()
	if err != nil {
		return nil, fmt.Errorf("GetIterator Error: %s", err)
	}
	defer headersIt.Release()

	for {
		has, err := headersIt.HasCurrentHeader()
		if err != nil {
			return nil, fmt.Errorf("HasCurrentHeader Error: %s", err)
		}
		if !has {
			break
		}

		name, value, err := headersIt.GetCurrentHeader()
		if err != nil {
			return nil, fmt.Errorf("GetCurrentHeader Error: %s", err)
		}

		header.Set(name, value)
		if _, err := headersIt.MoveNext(); err != nil {
			return nil, fmt.Errorf("MoveNext Error: %s", err)
		}
	}

	// WebView2 has problems when a request returns a 304 status code and the WebView2 is going to hang for other
	// requests including IPC calls.
	// So prevent 304 status codes by removing the headers that are used in combinationwith caching.
	header.Del("If-Modified-Since")
	header.Del("If-None-Match")
	return header, nil
}

func combineErrs(errs []error) error {
	// TODO use Go1.20 errors.Join
	if len(errs) == 0 {
		return nil
	}

	errStrings := make([]string, len(errs))
	for i, err := range errs {
		errStrings[i] = err.Error()
	}

	return fmt.Errorf(strings.Join(errStrings, "\n"))
}
