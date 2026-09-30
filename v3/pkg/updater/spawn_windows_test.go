//go:build windows

package updater

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Each job test runs in a subprocess because job membership cannot be undone.
func TestStartHelperWindowsJobs(t *testing.T) {
	for _, tt := range []struct {
		name  string
		flags uint32
	}{
		{"benign job forbids breakaway", 0},
		{"kill on close forbids breakaway", windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE},
		{"kill on close allows breakaway", windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE | windows.JOB_OBJECT_LIMIT_BREAKAWAY_OK},
		{"kill on close allows silent breakaway", windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE | windows.JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK},
	} {
		t.Run(tt.name, func(t *testing.T) {
			self, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(self, "-test.run=^TestStartHelperWindowsJobProcess$")
			cmd.Env = append(os.Environ(), "WAILS_TEST_JOB_FLAGS="+strconv.FormatUint(uint64(tt.flags), 10))
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("job subprocess: %v\n%s", err, output)
			}
		})
	}
}

func TestStartHelperWindowsJobProcess(t *testing.T) {
	flagText := os.Getenv("WAILS_TEST_JOB_FLAGS")
	if flagText == "" {
		return
	}
	flags, err := strconv.ParseUint(flagText, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	t.Cleanup(func() {
		// Clearing kill-on-close lets the subprocess exit normally.
		info.BasicLimitInformation.LimitFlags = 0
		_, _ = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)))
		_ = windows.CloseHandle(job)
	})
	info.BasicLimitInformation.LimitFlags = uint32(flags)
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		t.Fatal(err)
	}
	if err := windows.AssignProcessToJobObject(job, windows.CurrentProcess()); err != nil {
		t.Fatal(err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := newDetachedCommand(self)
	cmd.Args = []string{self, "-test.run=^$"}
	err = startHelper(cmd)
	if flags == windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE {
		if err == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			t.Fatal("helper started in a kill-on-close job that forbids breakaway")
		}
		if !errors.Is(wrapHelperSpawnError(err), ErrJobBreakawayDenied) {
			t.Fatalf("unexpected spawn error: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
}
