package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/term"
)

type androidAttachOperations struct {
	deviceState  func(context.Context, string) (string, bool, error)
	processID    func(context.Context, string, string) (string, error)
	startLogs    func(context.Context, string, string) (<-chan error, error)
	pollInterval time.Duration
	startupWait  time.Duration
}

type androidLogSession struct {
	cancel context.CancelFunc
	done   <-chan error
}

func attachAndroidApplication(ctx context.Context, device androidDevice, packageID string) error {
	term.Header("Android Logs")
	pterm.Info.Printfln("Streaming %s on %s · press Ctrl+C to stop", packageID, androidDeviceLabel(device))
	return attachAndroidApplicationWithOperations(ctx, device, packageID, androidAttachOperations{
		deviceState:  androidDeviceConnectionState,
		processID:    androidApplicationProcessID,
		startLogs:    startAndroidLogcat,
		pollInterval: time.Second,
		startupWait:  15 * time.Second,
	})
}

func androidDeviceConnectionState(ctx context.Context, serial string) (string, bool, error) {
	devices, err := listAndroidDevices(ctx)
	if err != nil {
		return "", false, err
	}
	for _, device := range devices {
		if device.Serial == serial {
			return device.State, true, nil
		}
	}
	return "", false, nil
}

func androidApplicationProcessID(ctx context.Context, serial, packageID string) (string, error) {
	output, err := runADB(ctx, "-s", serial, "shell", "pidof", packageID)
	if err != nil {
		if strings.TrimSpace(output) == "" {
			return "", nil
		}
		return "", fmt.Errorf("adb pidof: %s: %w", strings.TrimSpace(output), err)
	}
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return "", nil
	}
	pid, parseErr := strconv.Atoi(fields[0])
	if parseErr != nil || pid <= 0 {
		return "", fmt.Errorf("adb returned invalid process ID %q for %s", fields[0], packageID)
	}
	return fields[0], nil
}

func startAndroidLogcat(ctx context.Context, serial, pid string) (<-chan error, error) {
	return startAndroidLogcatWithWriters(ctx, serial, pid, os.Stdout, os.Stderr)
}

func startAndroidLogcatWithWriters(ctx context.Context, serial, pid string, stdout, stderr io.Writer) (<-chan error, error) {
	command := exec.CommandContext(ctx, "adb", "-s", serial, "logcat", "--pid="+pid, "-v", "brief")
	command.Stdout = stdout
	command.Stderr = stderr
	configureManifestProcess(command)
	command.Cancel = func() error {
		if command.Process == nil {
			return os.ErrProcessDone
		}
		return killManifestProcess(command.Process)
	}
	command.WaitDelay = 5 * time.Second
	if err := command.Start(); err != nil {
		return nil, err
	}
	result := make(chan error, 1)
	go func() {
		result <- command.Wait()
	}()
	return result, nil
}

func attachAndroidApplicationWithOperations(ctx context.Context, device androidDevice, packageID string, operations androidAttachOperations) error {
	pid, err := waitForAndroidApplicationProcess(ctx, device, packageID, operations)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	logContext, cancelLogs := context.WithCancel(ctx)
	logsDone, err := operations.startLogs(logContext, device.Serial, pid)
	if err != nil {
		cancelLogs()
		return fmt.Errorf("start Android application logs on %s: %w", device.Serial, err)
	}
	logs := androidLogSession{cancel: cancelLogs, done: logsDone}
	defer func() { logs.cancel() }()
	pollInterval := operations.pollInterval
	if pollInterval <= 0 {
		pollInterval = time.Second
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logs.cancel()
			<-logs.done
			return nil
		case logErr := <-logs.done:
			if ctx.Err() != nil {
				return nil
			}
			if logErr != nil {
				if targetErr := androidTargetAvailabilityError(ctx, device, operations); targetErr != nil {
					return targetErr
				}
				return fmt.Errorf("Android logcat stopped on %s: %w", device.Serial, logErr)
			}
			if targetErr := androidTargetAvailabilityError(ctx, device, operations); targetErr != nil {
				return targetErr
			}
			currentPID, pidErr := operations.processID(ctx, device.Serial, packageID)
			if pidErr != nil {
				return fmt.Errorf("inspect %s after Android logcat stopped on %s: %w", packageID, device.Serial, pidErr)
			}
			if currentPID == "" {
				return nil
			}
			return fmt.Errorf("Android logcat stopped while %s is still running on %s", packageID, device.Serial)
		case <-ticker.C:
			state, found, stateErr := operations.deviceState(ctx, device.Serial)
			if stateErr != nil {
				logs.cancel()
				<-logs.done
				return fmt.Errorf("monitor Android target %s: %w", device.Serial, stateErr)
			}
			if !found {
				logs.cancel()
				<-logs.done
				return fmt.Errorf("Android target %s disconnected", device.Serial)
			}
			if state != "device" {
				logs.cancel()
				<-logs.done
				return fmt.Errorf("Android target %s became %s", device.Serial, state)
			}
			currentPID, pidErr := operations.processID(ctx, device.Serial, packageID)
			if pidErr != nil {
				logs.cancel()
				<-logs.done
				if targetErr := androidTargetAvailabilityError(ctx, device, operations); targetErr != nil {
					return targetErr
				}
				return fmt.Errorf("monitor Android application %s on %s: %w", packageID, device.Serial, pidErr)
			}
			if currentPID == "" {
				logs.cancel()
				<-logs.done
				return nil
			}
			if currentPID != pid {
				logs.cancel()
				<-logs.done
				pid = currentPID
				nextContext, cancelNextLogs := context.WithCancel(ctx)
				nextDone, pidErr := operations.startLogs(nextContext, device.Serial, pid)
				if pidErr != nil {
					cancelNextLogs()
					return fmt.Errorf("restart Android application logs on %s for process %s: %w", device.Serial, pid, pidErr)
				}
				logs = androidLogSession{cancel: cancelNextLogs, done: nextDone}
			}
		}
	}
}

func androidTargetAvailabilityError(ctx context.Context, device androidDevice, operations androidAttachOperations) error {
	state, found, err := operations.deviceState(ctx, device.Serial)
	if err != nil {
		return fmt.Errorf("monitor Android target %s: %w", device.Serial, err)
	}
	if !found {
		return fmt.Errorf("Android target %s disconnected", device.Serial)
	}
	if state != "device" {
		return fmt.Errorf("Android target %s became %s", device.Serial, state)
	}
	return nil
}

func waitForAndroidApplicationProcess(ctx context.Context, device androidDevice, packageID string, operations androidAttachOperations) (string, error) {
	pollInterval := operations.pollInterval
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	startupWait := operations.startupWait
	if startupWait <= 0 {
		startupWait = 15 * time.Second
	}
	deadline := time.NewTimer(startupWait)
	defer deadline.Stop()
	for {
		pid, err := operations.processID(ctx, device.Serial, packageID)
		if err != nil {
			if targetErr := androidTargetAvailabilityError(ctx, device, operations); targetErr != nil {
				return "", targetErr
			}
			return "", fmt.Errorf("find Android application %s on %s: %w", packageID, device.Serial, err)
		}
		if pid != "" {
			return pid, nil
		}
		state, found, err := operations.deviceState(ctx, device.Serial)
		if err != nil {
			return "", fmt.Errorf("inspect Android target %s while waiting for %s: %w", device.Serial, packageID, err)
		}
		if !found {
			return "", fmt.Errorf("Android target %s disconnected while waiting for %s", device.Serial, packageID)
		}
		if state != "device" {
			return "", fmt.Errorf("Android target %s became %s while waiting for %s", device.Serial, state, packageID)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-deadline.C:
			return "", fmt.Errorf("Android application %s did not start on %s within %s", packageID, device.Serial, startupWait)
		case <-time.After(pollInterval):
		}
	}
}
