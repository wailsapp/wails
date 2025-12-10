//go:build !ios && !android

package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

// Empty method name maps for logging on non-mobile platforms
var iosMethodNames = map[int]string{}
var androidMethodNames = map[int]string{}

// processIOSMethod is a stub for non-mobile platforms
func (m *MessageProcessor) processIOSMethod(req *RuntimeRequest, window Window) (any, error) {
	return nil, errs.NewInvalidIOSCallErrorf("iOS methods not available on this platform")
}

// processAndroidMethod is a stub for non-mobile platforms
func (m *MessageProcessor) processAndroidMethod(req *RuntimeRequest, window Window) (any, error) {
	return nil, errs.NewInvalidAndroidCallErrorf("Android methods not available on this platform")
}
