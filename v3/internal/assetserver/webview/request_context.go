package webview

import "context"

// Context returns the native request's lifetime context when the backend exposes
// one. Backends without a cancellation signal retain their existing behavior.
func Context(r Request) context.Context {
	if r, ok := r.(interface{ Context() context.Context }); ok {
		return r.Context()
	}
	return context.Background()
}
