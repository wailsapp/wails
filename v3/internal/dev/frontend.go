package dev

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WaitHTTP requires an actual HTTP response, rather than merely an open port.
// Client errors are accepted for custom root routes; server errors keep waiting.
func WaitHTTP(ctx context.Context, p Process, address string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	transport := &http.Transport{Proxy: nil}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: time.Second}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	var last error
	for {
		if ctx.Err() != nil {
			if last == nil {
				return fmt.Errorf("frontend HTTP readiness: %w", ctx.Err())
			}
			return fmt.Errorf("frontend HTTP readiness: %v: %w", last, ctx.Err())
		}
		select {
		case <-p.Done():
			return fmt.Errorf("frontend exited before HTTP readiness: %v", p.Err())
		default:
		}
		req, e := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
		if e != nil {
			return e
		}
		response, e := client.Do(req)
		if e == nil {
			io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
			response.Body.Close()
			if response.StatusCode < 500 {
				return nil
			}
			last = fmt.Errorf("HTTP %d", response.StatusCode)
		} else {
			last = e
		}
		select {
		case <-ctx.Done():
		case <-p.Done():
			return fmt.Errorf("frontend exited before HTTP readiness: %v", p.Err())
		case <-ticker.C:
		}
	}
}
