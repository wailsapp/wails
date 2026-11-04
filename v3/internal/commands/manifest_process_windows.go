//go:build windows

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func configureManifestProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func signalManifestProcess(process *os.Process, _ os.Signal) error {
	return exec.Command("taskkill", "/T", "/PID", strconv.Itoa(process.Pid)).Run()
}

func killManifestProcess(process *os.Process) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(process.Pid)).Run()
}

// Start suspended so even an immediately exiting wrapper cannot spawn children
// before it belongs to the job. The non-inherited job handle owns all descendants
// independently of the parent PID, and closing it terminates the entire job.
func startManifestOwnedProcess(cmd *exec.Cmd) (func(), error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("create development process job: %w", err)
	}
	closeJob := func() { _ = windows.CloseHandle(job) }
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		closeJob()
		return nil, fmt.Errorf("configure development process job: %w", err)
	}
	configureManifestProcess(cmd)
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_SUSPENDED
	if err := cmd.Start(); err != nil {
		closeJob()
		return nil, err
	}
	fail := func(err error) (func(), error) {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		closeJob()
		return nil, err
	}
	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fail(fmt.Errorf("open development process: %w", err))
	}
	err = windows.AssignProcessToJobObject(job, process)
	_ = windows.CloseHandle(process)
	if err != nil {
		return fail(fmt.Errorf("assign development process job: %w", err))
	}
	if err := resumeManifestProcess(uint32(cmd.Process.Pid)); err != nil {
		return fail(err)
	}
	return closeJob, nil
}

// os/exec closes the primary thread handle during Start. Enumerate the newly
// created, still-suspended process's thread to resume it using the documented API.
func resumeManifestProcess(pid uint32) error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("find development process thread: %w", err)
	}
	defer windows.CloseHandle(snapshot)
	entry := windows.ThreadEntry32{Size: uint32(unsafe.Sizeof(windows.ThreadEntry32{}))}
	for err = windows.Thread32First(snapshot, &entry); err == nil; err = windows.Thread32Next(snapshot, &entry) {
		if entry.OwnerProcessID != pid {
			continue
		}
		thread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, entry.ThreadID)
		if err != nil {
			return fmt.Errorf("open development process thread: %w", err)
		}
		_, err = windows.ResumeThread(thread)
		_ = windows.CloseHandle(thread)
		if err != nil {
			return fmt.Errorf("resume development process: %w", err)
		}
		return nil
	}
	return fmt.Errorf("find suspended development process %d thread: %w", pid, err)
}
