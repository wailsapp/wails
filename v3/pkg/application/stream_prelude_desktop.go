//go:build !server

package application

// streamPrelude is empty in webview builds: there is no listener to upgrade on,
// so Stream() uses the held-poll transport and needs no factory. See the server
// build of this file for why the choice cannot be left to custom.js.
func streamPrelude() []byte { return nil }
