// Package ipcauth holds the per-launch capability that hardens access to the
// development IPC server against the browser vector.
//
// The `wails dev` IPC WebSocket forwards straight to the bound-method
// dispatcher. The token is minted once per launch, handed only to pages the dev
// server itself serves (as an HttpOnly, SameSite=Strict cookie), and required
// on every WebSocket upgrade. Together with the dev server's Host allowlist and
// same-origin check it raises the bar for a malicious page in the developer's
// browser.
//
// It is deliberately not a defence against other processes running as the same
// user: such a process can read the cookie from a served response, or read the
// developer's memory and traffic outright, so it is outside this package's
// threat model.
package ipcauth

import (
	"crypto/rand"
	"crypto/subtle"
	"sync"
)

// CookieName is the cookie the dev server sets on the pages it serves and
// requires back on the IPC WebSocket upgrade.
const CookieName = "wails_ipc_capability"

var (
	once  sync.Once
	token string
)

// Token returns the process-wide dev-IPC capability, minting it on first use.
// rand.Text is a CSPRNG string and cannot fail, so obtaining the capability can
// never be the thing that breaks the dev server.
func Token() string {
	once.Do(func() { token = rand.Text() })
	return token
}

// Valid reports whether presented equals the capability. subtle.ConstantTimeCompare
// returns early when the lengths differ, so the compare is constant-time only
// across equal-length inputs; that is enough here, since the token has a fixed
// length and a length mismatch already means the wrong value.
func Valid(presented string) bool {
	return subtle.ConstantTimeCompare([]byte(presented), []byte(Token())) == 1
}
