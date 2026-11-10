package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// The desktop scenarios cover what other programs do to the app: a second
// copy started by the user hands its arguments to the first, and (Linux) the
// panel shows the tray icon (tray_linux.go).

// secondInstanceEnv marks the copy the second-instance scenario starts.
const secondInstanceEnv = "E2E_SECOND_INSTANCE"

// secondInstanceExit is the copy's exit code once it has handed over.
const secondInstanceExit = 3

// singleInstance makes the probe a single-instance app. The key is fixed:
// both copies need the same one, and it makes them encrypt what they send.
func singleInstance() *application.SingleInstanceOptions {
	return &application.SingleInstanceOptions{
		UniqueID:       "com.wails.e2e-probe",
		ExitCode:       secondInstanceExit,
		EncryptionKey:  [32]byte{'w', 'a', 'i', 'l', 's', '-', 'e', '2', 'e'},
		AdditionalData: map[string]string{"copy": os.Getenv(secondInstanceEnv)},
		OnSecondInstanceLaunch: func(d application.SecondInstanceData) {
			if slices.Contains(d.Args, "--e2e-handover") {
				emit("second-instance " + d.AdditionalData["copy"])
			}
		},
	}
}

// secondInstance starts a second copy, which must hand its arguments to this
// one and exit instead of running.
func secondInstance(p *probe) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "--e2e-handover")
	cmd.Env = append(os.Environ(), secondInstanceEnv+"=second")
	out := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { out <- cmd.Wait() }()
	if err := await("second-instance second", 15*time.Second); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("the first copy heard nothing: %w", err)
	}
	select {
	case err := <-out:
		if code := cmd.ProcessState.ExitCode(); code != secondInstanceExit {
			return fmt.Errorf("the second copy exited with %d (%v), want %d", code, err, secondInstanceExit)
		}
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		return fmt.Errorf("the second copy kept running after handing over")
	}
	return nil
}
