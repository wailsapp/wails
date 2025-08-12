package process

import (
	"os"
	"os/exec"
)

// Process defines a process that can be executed
type Process struct {
	cmd         *exec.Cmd
	exitChannel chan bool
	Running     bool
}

// NewProcess creates a new process struct
func NewProcess(cmd string, args ...string) *Process {
	result := &Process{
		cmd:         exec.Command(cmd, args...),
		exitChannel: make(chan bool, 1),
	}
	result.cmd.Stdout = os.Stdout
	result.cmd.Stderr = os.Stderr
	return result
}

// Start the process
func (p *Process) Start(exitCodeChannel chan int) error {
	err := p.cmd.Start()
	if err != nil {
		return err
	}

	p.Running = true

	go func(cmd *exec.Cmd, running *bool, exitChannel chan bool, exitCodeChannel chan int) {
		err := cmd.Wait()
		if err == nil {
			exitCodeChannel <- 0
		}
		*running = false
		exitChannel <- true
	}(p.cmd, &p.Running, p.exitChannel, exitCodeChannel)

	return nil
}

// Kill the process
func (p *Process) Kill() error {
	if !p.Running {
		return nil
	}
	err := p.cmd.Process.Kill()
	if err != nil {
		return err
	}
	err = p.cmd.Process.Release()
	if err != nil {
		return err
	}

	// Wait for command to exit properly
	<-p.exitChannel

	return err
}

// PID returns the process PID
func (p *Process) PID() int {
	return p.cmd.Process.Pid
}

func (p *Process) SetDir(dir string) {
	p.cmd.Dir = dir
}
