package commands

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseADBDevicesPreservesUsefulStateAndMetadata(t *testing.T) {
	devices := parseADBDevices(`List of devices attached
emulator-5554 device product:sdk_gphone64_x86_64 model:sdk_gphone64_x86_64 device:emu64xa transport_id:1
R5CT11ABC unauthorized usb:1-2 transport_id:2
offline-one offline transport_id:3

`)
	require.Len(t, devices, 3)
	assert.Equal(t, androidDevice{Serial: "emulator-5554", State: "device", Model: "sdk_gphone64_x86_64", Product: "sdk_gphone64_x86_64", Emulator: true}, devices[0])
	assert.Equal(t, "unauthorized", devices[1].State)
	assert.False(t, devices[1].Emulator)
}

func TestAndroidRunRequiresDeterministicSelectionWithoutATerminal(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.listDevices = func(context.Context) ([]androidDevice, error) {
		return []androidDevice{{Serial: "one", State: "device"}, {Serial: "two", State: "device"}}, nil
	}
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{}, "", operations)
	assert.ErrorContains(t, err, "--device or --emulator")
}

func TestAndroidRunRejectsConflictingSelectorsAndPropagatesABIFailure(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "one", Emulator: "pixel"}, "", operations)
	assert.ErrorContains(t, err, "mutually exclusive")

	operations.deviceABI = func(context.Context, string) (string, error) { return "", errors.New("unsupported") }
	err = androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone"}, "", operations)
	assert.ErrorContains(t, err, "inspect Android device phone ABI")
}

func TestAndroidRunCanInstallExistingAPKWithoutLaunching(t *testing.T) {
	apk := filepath.Join(t.TempDir(), "app.apk")
	require.NoError(t, os.WriteFile(apk, []byte("apk"), 0o644))
	operations := fakeAndroidDeployOperations()
	operations.packageID = func(string) (string, error) { return "com.example.existing", nil }
	var calls [][]string
	operations.adb = func(_ context.Context, args ...string) (string, error) {
		calls = append(calls, append([]string(nil), args...))
		return "Success", nil
	}

	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone", APK: apk, NoLaunch: true}, "", operations)
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"-s", "phone", "install", "-r", apk}}, calls)
}

func TestAndroidRunReportsInstallAndLaunchFailures(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.adb = func(_ context.Context, args ...string) (string, error) {
		return "permission denied", errors.New("exit 1")
	}
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone"}, "", operations)
	assert.ErrorContains(t, err, "install APK")
	assert.ErrorContains(t, err, "permission denied")

	install := true
	operations.adb = func(_ context.Context, args ...string) (string, error) {
		if install {
			install = false
			return "Success", nil
		}
		return "activity missing", errors.New("exit 1")
	}
	err = androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone"}, "", operations)
	assert.ErrorContains(t, err, "launch com.example.app")
}

func TestAndroidDeviceSelectionCanStartAnInteractivelySelectedAVD(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.interactive = func() bool { return true }
	operations.listDevices = func(context.Context) ([]androidDevice, error) { return nil, nil }
	operations.listAVDs = func(context.Context) ([]string, error) { return []string{"Pixel_8", "Tablet"}, nil }
	selections := []string{"@start-emulator", "Pixel_8"}
	operations.selectValue = func(_, _ string, choices []androidSelection) (string, error) {
		require.NotEmpty(t, choices)
		value := selections[0]
		selections = selections[1:]
		return value, nil
	}
	operations.startAVD = func(_ context.Context, name string) (androidDevice, error) {
		assert.Equal(t, "Pixel_8", name)
		return androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}, nil
	}

	device, err := chooseAndroidDevice(context.Background(), AndroidRunOptions{}, operations)
	require.NoError(t, err)
	assert.Equal(t, "emulator-5554", device.Serial)
}

func TestAndroidDeviceSelectionReportsMissingAVDsAndDevices(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.interactive = func() bool { return true }
	operations.listDevices = func(context.Context) ([]androidDevice, error) { return nil, nil }
	operations.selectValue = func(string, string, []androidSelection) (string, error) { return "@start-emulator", nil }
	device, err := chooseAndroidDevice(context.Background(), AndroidRunOptions{}, operations)
	assert.ErrorContains(t, err, "no Android Virtual Devices")
	assert.Empty(t, device.Serial)

	operations.listDevices = func(context.Context) ([]androidDevice, error) {
		return []androidDevice{{Serial: "other", State: "device"}}, nil
	}
	_, err = chooseAndroidDevice(context.Background(), AndroidRunOptions{Device: "missing"}, operations)
	assert.ErrorContains(t, err, "not connected")
}

func TestParseAndroidAVDNameHandlesADBLineEndings(t *testing.T) {
	assert.Equal(t, "Pixel_8", parseAndroidAVDName("Pixel_8\r\nOK\r\n"))
	assert.Empty(t, parseAndroidAVDName("OK\n"))
}

func TestAndroidCommandAdaptersDiscoverDevicesABIsAndAVDs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell adapter fixture")
	}
	tools := t.TempDir()
	adb := `#!/bin/sh
case "$*" in
  "devices -l") printf 'List of devices attached\nphone device model:Pixel product:pixel\n' ;;
  "-s phone shell getprop ro.product.cpu.abilist") printf 'arm64-v8a,armeabi-v7a\n' ;;
  "-s x86 shell getprop ro.product.cpu.abilist") printf 'x86_64,arm64-v8a,x86\n' ;;
  "-s phone shell pidof com.example.app") printf '123\n' ;;
  "-s phone logcat --pid=123 -v brief") printf 'W/Wails: application log\n' ;;
  "-s emulator-5554 emu kill") printf 'OK: killing emulator\n' ;;
  *) printf 'bad adb invocation: %s\n' "$*" >&2; exit 3 ;;
esac
`
	emulator := "#!/bin/sh\n[ \"$1\" = \"-list-avds\" ] || exit 4\nprintf 'Tablet\\nPixel_8\\n'\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "adb"), []byte(adb), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tools, "emulator"), []byte(emulator), 0o755))
	t.Setenv("PATH", tools)

	devices, err := listAndroidDevices(context.Background())
	require.NoError(t, err)
	require.Len(t, devices, 1)
	assert.Equal(t, "Pixel", devices[0].Model)
	arch, err := androidDeviceABI(context.Background(), "phone")
	require.NoError(t, err)
	assert.Equal(t, "arm64", arch)
	arch, err = androidDeviceABI(context.Background(), "x86")
	require.NoError(t, err)
	assert.Equal(t, "amd64", arch)
	_, err = androidDeviceABI(context.Background(), "unsupported")
	require.Error(t, err)
	pid, err := androidApplicationProcessID(context.Background(), "phone", "com.example.app")
	require.NoError(t, err)
	assert.Equal(t, "123", pid)
	var logs bytes.Buffer
	logsDone, err := startAndroidLogcatWithWriters(context.Background(), "phone", "123", &logs, &logs)
	require.NoError(t, err)
	require.NoError(t, <-logsDone)
	assert.Contains(t, logs.String(), "application log")
	require.NoError(t, stopAndroidEmulator(context.Background(), "emulator-5554"))
	avds, err := listAndroidAVDs(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"Pixel_8", "Tablet"}, avds)
}

func TestStartAndroidAVDReusesAnAlreadyRunningNamedEmulator(t *testing.T) {
	device := androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}
	actual, err := startAndroidAVD(context.Background(), "Pixel_8",
		func(context.Context) ([]androidDevice, error) { return []androidDevice{device}, nil },
		func(_ context.Context, arguments ...string) (string, error) {
			switch arguments[len(arguments)-1] {
			case "name":
				return "Pixel_8\r\nOK\r\n", nil
			case "sys.boot_completed":
				return "1\n", nil
			case "android":
				return "package:/system/framework/framework-res.apk\n", nil
			default:
				return "", errors.New("unexpected adb call")
			}
		},
	)
	require.NoError(t, err)
	assert.Equal(t, device, actual)
}

func TestInspectAndroidAVDReadinessRequiresBootAndPackageManager(t *testing.T) {
	device := androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}
	bootCompleted := ""
	packageManagerReady := false
	adb := func(_ context.Context, arguments ...string) (string, error) {
		switch arguments[len(arguments)-1] {
		case "name":
			return "Pixel_8\nOK\n", nil
		case "sys.boot_completed":
			return bootCompleted, nil
		case "android":
			if packageManagerReady {
				return "package:/system/framework/framework-res.apk\n", nil
			}
			return "", errors.New("package manager unavailable")
		default:
			return "", errors.New("unexpected adb call")
		}
	}

	_, running, ready := inspectAndroidAVDReadiness(context.Background(), "Pixel_8", []androidDevice{device}, adb)
	assert.True(t, running)
	assert.False(t, ready)

	bootCompleted = "1\n"
	_, running, ready = inspectAndroidAVDReadiness(context.Background(), "Pixel_8", []androidDevice{device}, adb)
	assert.True(t, running)
	assert.False(t, ready)

	packageManagerReady = true
	actual, running, ready := inspectAndroidAVDReadiness(context.Background(), "Pixel_8", []androidDevice{device}, adb)
	assert.True(t, running)
	assert.True(t, ready)
	assert.Equal(t, device, actual)
}

func TestStartAndroidAVDReportsInitialDeviceDiscoveryFailure(t *testing.T) {
	_, err := startAndroidAVD(context.Background(), "Pixel_8",
		func(context.Context) ([]androidDevice, error) { return nil, errors.New("adb server unavailable") },
		func(context.Context, ...string) (string, error) { return "", nil },
	)
	assert.ErrorContains(t, err, "list Android devices before starting emulator")
	assert.ErrorContains(t, err, "adb server unavailable")
}

func TestAndroidRunBuildsInstallsAndLaunchesTheSelectedDeviceABI(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.listDevices = func(context.Context) ([]androidDevice, error) {
		return []androidDevice{{Serial: "phone", State: "device", Model: "Pixel"}}, nil
	}
	operations.deviceABI = func(context.Context, string) (string, error) { return "arm64", nil }
	var builtProfile, builtArch string
	operations.buildAPK = func(_ context.Context, profile, arch string) (string, string, error) {
		builtProfile, builtArch = profile, arch
		return "/tmp/app.apk", "com.example.app", nil
	}
	var calls [][]string
	operations.adb = func(_ context.Context, args ...string) (string, error) {
		calls = append(calls, append([]string(nil), args...))
		return "Success", nil
	}

	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone"}, "mobile", operations)
	require.NoError(t, err)
	assert.Equal(t, "mobile", builtProfile)
	assert.Equal(t, "arm64", builtArch)
	assert.Equal(t, [][]string{
		{"-s", "phone", "install", "-r", "/tmp/app.apk"},
		{"-s", "phone", "shell", "am", "start", "-n", "com.example.app/com.wails.app.MainActivity"},
	}, calls)
}

func TestAndroidRunCanAttachApplicationLogsAfterLaunch(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	attached := false
	operations.attach = func(_ context.Context, device androidDevice, packageID string) error {
		attached = true
		assert.Equal(t, "phone", device.Serial)
		assert.Equal(t, "com.example.app", packageID)
		return nil
	}

	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone", Logs: true}, "", operations)
	require.NoError(t, err)
	assert.True(t, attached)
}

func TestAndroidRunRejectsLogsWithoutLaunching(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone", Logs: true, NoLaunch: true}, "", operations)
	assert.ErrorContains(t, err, "--logs cannot be combined with --no-launch")
}

func TestAndroidRunRequiresAttachedSessionForEmulatorCleanup(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Emulator: "Pixel_8", StopEmulator: true}, "", operations)
	assert.ErrorContains(t, err, "--stop-emulator requires --logs")
}

func TestAndroidRunNeverAppliesEmulatorCleanupToAPhysicalDevice(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone", Logs: true, StopEmulator: true}, "", operations)
	assert.ErrorContains(t, err, "--stop-emulator requires an emulator target")
}

func TestAndroidRunStopsAnEmulatorOnlyWhenExplicitlyRequested(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.startAVD = func(context.Context, string) (androidDevice, error) {
		return androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}, nil
	}
	operations.attach = func(context.Context, androidDevice, string) error { return nil }
	stopped := ""
	operations.stopEmulator = func(_ context.Context, serial string) error {
		stopped = serial
		return nil
	}

	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Emulator: "Pixel_8", Logs: true, StopEmulator: true}, "", operations)
	require.NoError(t, err)
	assert.Equal(t, "emulator-5554", stopped)
}

func TestAndroidRunReportsBothAttachedSessionAndCleanupFailures(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.startAVD = func(context.Context, string) (androidDevice, error) {
		return androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}, nil
	}
	operations.attach = func(context.Context, androidDevice, string) error { return errors.New("log stream failed") }
	operations.stopEmulator = func(context.Context, string) error { return errors.New("emulator kill failed") }

	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Emulator: "Pixel_8", Logs: true, StopEmulator: true}, "", operations)
	require.ErrorContains(t, err, "log stream failed")
	require.ErrorContains(t, err, "emulator kill failed")
}

func TestAndroidAttachedLogsExitCleanlyOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID:   func(context.Context, string, string) (string, error) { return "123", nil },
		startLogs: func(logContext context.Context, serial, pid string) (<-chan error, error) {
			assert.Equal(t, "phone", serial)
			assert.Equal(t, "123", pid)
			result := make(chan error, 1)
			close(started)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	done := make(chan error, 1)
	go func() {
		done <- attachAndroidApplicationWithOperations(ctx, androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	}()
	select {
	case <-started:
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Android log attachment")
	}
	cancel()
	require.NoError(t, <-done)
}

func TestAndroidAttachedLogsCancelCleanlyWhileWaitingForTheApplication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	operations := androidAttachOperations{
		deviceState:  func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID:    func(context.Context, string, string) (string, error) { return "", nil },
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	err := attachAndroidApplicationWithOperations(ctx, androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	require.NoError(t, err)
}

func TestAndroidAttachedLogsWaitForTheLaunchedApplicationProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0
	started := make(chan struct{})
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID: func(context.Context, string, string) (string, error) {
			attempts++
			if attempts < 3 {
				return "", nil
			}
			return "456", nil
		},
		startLogs: func(logContext context.Context, _, pid string) (<-chan error, error) {
			assert.Equal(t, "456", pid)
			result := make(chan error, 1)
			close(started)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	done := make(chan error, 1)
	go func() {
		done <- attachAndroidApplicationWithOperations(ctx, androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	}()
	select {
	case <-started:
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Android log attachment")
	}
	cancel()
	require.NoError(t, <-done)
	assert.GreaterOrEqual(t, attempts, 3)
}

func TestAndroidAttachedLogsReportDeviceLoss(t *testing.T) {
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "", false, nil },
		processID:   func(context.Context, string, string) (string, error) { return "123", nil },
		startLogs: func(logContext context.Context, _, _ string) (<-chan error, error) {
			result := make(chan error, 1)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	done := make(chan error, 1)
	go func() {
		done <- attachAndroidApplicationWithOperations(context.Background(), androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	}()
	select {
	case err := <-done:
		require.ErrorContains(t, err, "phone disconnected")
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for device-loss detection")
	}
}

func TestAndroidAttachedLogsPreferDeviceLossOverConcurrentPIDFailure(t *testing.T) {
	stateChecks := 0
	processChecks := 0
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) {
			stateChecks++
			if stateChecks == 1 {
				return "device", true, nil
			}
			return "", false, nil
		},
		processID: func(context.Context, string, string) (string, error) {
			processChecks++
			if processChecks == 1 {
				return "123", nil
			}
			return "", errors.New("adb transport closed")
		},
		startLogs: func(logContext context.Context, _, _ string) (<-chan error, error) {
			result := make(chan error, 1)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}

	err := attachAndroidApplicationWithOperations(context.Background(), androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	require.ErrorContains(t, err, "Android target phone disconnected")
	assert.NotContains(t, err.Error(), "adb transport closed")
}

func TestAndroidAttachedLogsFinishWhenTheApplicationStops(t *testing.T) {
	processChecks := 0
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID: func(context.Context, string, string) (string, error) {
			processChecks++
			if processChecks == 1 {
				return "123", nil
			}
			return "", nil
		},
		startLogs: func(logContext context.Context, _, _ string) (<-chan error, error) {
			result := make(chan error, 1)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	done := make(chan error, 1)
	go func() {
		done <- attachAndroidApplicationWithOperations(context.Background(), androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for application-exit detection")
	}
}

func TestAndroidAttachedLogsFollowAnApplicationRestart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	processChecks := 0
	var loggedPIDs []string
	secondStarted := make(chan struct{})
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID: func(context.Context, string, string) (string, error) {
			processChecks++
			if processChecks == 1 {
				return "123", nil
			}
			return "456", nil
		},
		startLogs: func(logContext context.Context, _, pid string) (<-chan error, error) {
			loggedPIDs = append(loggedPIDs, pid)
			if pid == "456" {
				close(secondStarted)
			}
			result := make(chan error, 1)
			go func() {
				<-logContext.Done()
				result <- logContext.Err()
			}()
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	done := make(chan error, 1)
	go func() {
		done <- attachAndroidApplicationWithOperations(ctx, androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	}()
	select {
	case <-secondStarted:
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for restarted application logs")
	}
	cancel()
	require.NoError(t, <-done)
	assert.Equal(t, []string{"123", "456"}, loggedPIDs)
}

func TestAndroidAttachedLogsReportUnexpectedLogcatExit(t *testing.T) {
	operations := androidAttachOperations{
		deviceState: func(context.Context, string) (string, bool, error) { return "device", true, nil },
		processID:   func(context.Context, string, string) (string, error) { return "123", nil },
		startLogs: func(context.Context, string, string) (<-chan error, error) {
			result := make(chan error, 1)
			result <- nil
			return result, nil
		},
		pollInterval: time.Millisecond,
		startupWait:  time.Second,
	}
	err := attachAndroidApplicationWithOperations(context.Background(), androidDevice{Serial: "phone", State: "device"}, "com.example.app", operations)
	require.ErrorContains(t, err, "logcat stopped while com.example.app is still running")
}

func TestAndroidRunReportsUnavailableSelectedDevice(t *testing.T) {
	operations := fakeAndroidDeployOperations()
	operations.listDevices = func(context.Context) ([]androidDevice, error) {
		return []androidDevice{{Serial: "phone", State: "unauthorized"}}, nil
	}
	err := androidRunWithOperations(context.Background(), AndroidRunOptions{Device: "phone"}, "", operations)
	assert.ErrorContains(t, err, "unauthorized")

	operations.listDevices = func(context.Context) ([]androidDevice, error) { return nil, errors.New("adb failed") }
	err = androidRunWithOperations(context.Background(), AndroidRunOptions{}, "", operations)
	assert.ErrorContains(t, err, "adb failed")
}

func fakeAndroidDeployOperations() androidDeployOperations {
	return androidDeployOperations{
		interactive: func() bool { return false },
		listDevices: func(context.Context) ([]androidDevice, error) {
			return []androidDevice{{Serial: "phone", State: "device"}}, nil
		},
		deviceABI: func(context.Context, string) (string, error) { return "amd64", nil },
		buildAPK: func(context.Context, string, string) (string, string, error) {
			return "/tmp/app.apk", "com.example.app", nil
		},
		packageID: func(string) (string, error) { return "com.example.app", nil },
		adb:       func(context.Context, ...string) (string, error) { return "", nil },
		selectValue: func(string, string, []androidSelection) (string, error) {
			return "", errors.New("unexpected selection")
		},
		listAVDs: func(context.Context) ([]string, error) { return nil, nil },
		startAVD: func(context.Context, string) (androidDevice, error) {
			return androidDevice{}, errors.New("unexpected emulator start")
		},
	}
}
