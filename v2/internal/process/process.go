package process

import (
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// Process defines a process that can be executed
type Process struct {
	logger      *clilogger.CLILogger
	cmd         *exec.Cmd
	exitChannel chan bool
	Running     bool
}

// NewProcess creates a new process struct
func NewProcess(logger *clilogger.CLILogger, cmd string, args ...string) *Process {
	return &Process{
		logger:      logger,
		cmd:         exec.Command(cmd, args...),
		exitChannel: make(chan bool, 1),
	}
}

// Start the process
func (p *Process) Start() error {

	err := p.cmd.Start()
	if err != nil {
		return err
	}

	p.Running = true

	go func(cmd *exec.Cmd, running *bool, logger *clilogger.CLILogger, exitChannel chan bool) {
		logger.Println("Starting process (PID: %d)", cmd.Process.Pid)
		cmd.Wait()
		logger.Println("Exiting process (PID: %d)", cmd.Process.Pid)
		*running = false
		exitChannel <- true
	}(p.cmd, &p.Running, p.logger, p.exitChannel)

	return nil
}

// Kill the process
func (p *Process) Kill() error {
	if !p.Running {
		return nil
	}
	err := p.cmd.Process.Kill()

	// Wait for command to exit properly
	<-p.exitChannel

	return err
}

// PID returns the process PID
func (p *Process) PID() int {
	return p.cmd.Process.Pid
}
