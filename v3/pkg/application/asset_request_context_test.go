package application

import (
	"context"
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
