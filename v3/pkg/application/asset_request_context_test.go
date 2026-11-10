package application

import (
	"context"
	"net/http"
	"testing"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type contextualAssetRequest struct {
	webview.Request
	ctx context.Context
}

func (r contextualAssetRequest) Context() context.Context { return r.ctx }

func TestAssetRequestPreservesNativeContext(t *testing.T) {
	ctx, abort := context.WithCancel(context.Background())
	defer abort()
	r := &webViewAssetRequest{Request: contextualAssetRequest{ctx: ctx}}
	actual, release := webview.RequestContext(r)
	defer release()
	abort()
	if actual.Err() != context.Canceled {
		t.Fatal("application wrapper discarded native cancellation")
	}
}

func TestAssetRequestWithoutNativeContext(t *testing.T) {
	r := &webViewAssetRequest{}
	ctx, release := webview.RequestContext(r)
	if ctx.Err() != nil {
		t.Fatal("fallback context starts cancelled")
	}
	release()
	if ctx.Err() != context.Canceled {
		t.Fatal("fallback context cannot be released")
	}
}

type headerAssetRequest struct {
	webview.Request
	header http.Header
}

func (r headerAssetRequest) Header() (http.Header, error) { return r.header, nil }

func TestAssetRequestWindowHeaders(t *testing.T) {
	cases := []struct {
		name, sent, want string
	}{
		// A runtime call from the page names its target window, as
		// Window.Get("main") does: that name is kept.
		{"page names a target", "main", "main"},
		// Otherwise the request carries its own window's name.
		{"no name sent", "", "caller"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sent := http.Header{}
			if c.sent != "" {
				sent.Set(webViewRequestHeaderWindowName, c.sent)
			}
			// The ID always comes from the platform: a page cannot pick it.
			sent.Set(webViewRequestHeaderWindowId, "99")
			r := &webViewAssetRequest{Request: headerAssetRequest{header: sent}, windowId: 7, windowName: "caller"}
			h, err := r.Header()
			if err != nil {
				t.Fatal(err)
			}
			if got := h.Get(webViewRequestHeaderWindowName); got != c.want {
				t.Errorf("window name %q, want %q", got, c.want)
			}
			if got := h.Get(webViewRequestHeaderWindowId); got != "7" {
				t.Errorf("window ID %q, want 7", got)
			}
		})
	}
}
