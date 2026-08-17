//go:build darwin && !ios && wails_native

package webview

import (
	"errors"
	"io"
	"net/http"
	"unsafe"
)

var errUnavailableInNativeBuild = errors.New("WebView requests are unavailable in a wails_native build")

// NewRequest exists only so the additive v3 application package retains its
// source shape. No native build can produce the WKURLSchemeTask pointer this
// function would normally wrap.
func NewRequest(unsafe.Pointer) Request { return nativeBuildRequest{} }

type nativeBuildRequest struct{}

func (nativeBuildRequest) URL() (string, error)         { return "", errUnavailableInNativeBuild }
func (nativeBuildRequest) Method() (string, error)      { return "", errUnavailableInNativeBuild }
func (nativeBuildRequest) Header() (http.Header, error) { return nil, errUnavailableInNativeBuild }
func (nativeBuildRequest) Body() (io.ReadCloser, error) { return nil, errUnavailableInNativeBuild }
func (nativeBuildRequest) Response() ResponseWriter     { return nativeBuildResponse{} }
func (nativeBuildRequest) Close() error                 { return nil }

type nativeBuildResponse struct{}

func (nativeBuildResponse) Header() http.Header       { return http.Header{} }
func (nativeBuildResponse) Write([]byte) (int, error) { return 0, errUnavailableInNativeBuild }
func (nativeBuildResponse) WriteHeader(int)           {}
func (nativeBuildResponse) Finish() error             { return errUnavailableInNativeBuild }
func (nativeBuildResponse) Code() int                 { return 0 }
