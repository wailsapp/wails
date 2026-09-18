//go:build wails_native && !wails_single_instance

package application

/*
#include <stdlib.h>
*/
import "C"

import "errors"

var alreadyRunningError = errors.New("application is already running")

// SecondInstanceData is retained so shared source can describe second-instance
// callbacks in both build modes.
type SecondInstanceData struct {
	Args           []string          `json:"args"`
	WorkingDir     string            `json:"workingDir"`
	AdditionalData map[string]string `json:"additionalData,omitempty"`
}

// SingleInstanceOptions is source-compatible with the full implementation.
// Build native applications that use it with both wails_native and
// wails_single_instance.
type SingleInstanceOptions struct {
	// UniqueID identifies the application instance.
	UniqueID string

	// OnSecondInstanceLaunch receives data from a subsequent launch.
	OnSecondInstanceLaunch func(data SecondInstanceData)

	// AdditionalData is sent from a subsequent launch to the first instance.
	AdditionalData map[string]string

	// ExitCode is returned by a subsequent instance after notifying the first.
	ExitCode int

	// EncryptionKey encrypts instance communication when it is non-zero.
	EncryptionKey [32]byte
}

type singleInstanceManager struct{}

func newSingleInstanceManager(_ *App, options *SingleInstanceOptions) (*singleInstanceManager, error) {
	if options == nil {
		return nil, nil
	}
	return nil, errors.New("single-instance support is unavailable in this wails_native build; add the wails_single_instance build tag")
}

func (*singleInstanceManager) notifyFirstInstance() error { return nil }
func (*singleInstanceManager) cleanup()                   {}

//export handleSecondInstanceData
func handleSecondInstanceData(*C.char) {}
