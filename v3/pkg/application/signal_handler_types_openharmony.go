//go:build openharmony && !server

package application

import "os"

// signalType is a placeholder type for OHOS where signals are not used.
type signalType = os.Signal
