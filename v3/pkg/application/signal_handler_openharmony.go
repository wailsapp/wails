//go:build openharmony && !server

package application

import "os"

// On OHOS we don't handle Unix signals directly; the app lifecycle is managed
// by the ArkTS ability (mirrors signal_handler_android.go).
func (a *App) setupSignalHandler(options Options) {}

func handleSignal(_ os.Signal) {}
