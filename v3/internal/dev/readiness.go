package dev

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/devruntime"
)

// Readiness is unique to one candidate launch; an old application cannot mark
// its replacement ready. Desktop listeners bind loopback; physical iOS uses
// the explicitly selected LAN interface.
type Readiness struct {
	server         *http.Server
	address, token string
	ready          chan struct{}
	once           sync.Once
}

func NewReadiness() (*Readiness, error) { return NewReadinessOn("127.0.0.1") }
func NewReadinessOn(host string) (*Readiness, error) {
	l, e := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if e != nil {
		return nil, e
	}
	var secret [32]byte
	if _, e = rand.Read(secret[:]); e != nil {
		l.Close()
		return nil, e
	}
	r := &Readiness{address: "http://" + l.Addr().String() + "/ready", token: hex.EncodeToString(secret[:]), ready: make(chan struct{})}
	r.server = &http.Server{ReadHeaderTimeout: 2 * time.Second, IdleTimeout: 2 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
		if q.Method != http.MethodPost || q.URL.Path != "/ready" || subtle.ConstantTimeCompare([]byte(q.Header.Get("Authorization")), []byte("Bearer "+r.token)) != 1 {
			http.Error(w, "unauthorised", http.StatusUnauthorized)
			return
		}
		r.once.Do(func() { close(r.ready) })
		w.WriteHeader(http.StatusNoContent)
	})}
	go r.server.Serve(l)
	return r, nil
}
func (r *Readiness) Environment() []string {
	return []string{devruntime.AddressEnv + "=" + r.address, devruntime.TokenEnv + "=" + r.token}
}
func (r *Readiness) Close() {
	if r != nil {
		r.server.Close()
	}
}
func (r *Readiness) Wait(ctx context.Context, p Process, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	defer r.Close()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.Done():
		return processFailure(p, "backend", "readiness", "process exited before Wails runtime became ready")
	case <-timer.C:
		return &Diagnostic{Component: "backend", Phase: "readiness", ExitCode: -1, Err: fmt.Errorf("timed out waiting for Wails runtime readiness after %s", timeout)}
	case <-r.ready:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		select {
		case <-p.Done():
			return processFailure(p, "backend", "readiness", "backend exited during readiness")
		default:
			return nil
		}
	}
}

func (r *Readiness) Port() string {
	u := strings.TrimPrefix(r.address, "http://")
	h, _, _ := strings.Cut(u, "/")
	_, p, _ := net.SplitHostPort(h)
	return p
}
