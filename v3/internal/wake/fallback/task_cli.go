package fallback

import (
	"fmt"
	"os"
	"os/exec"
)

type ErrUnsupported struct {
	Feature string
}

func (e *ErrUnsupported) Error() string {
	return fmt.Sprintf("unsupported feature: %s", e.Feature)
}

func TaskCLI(name, dir string, env []string) error {
	c := exec.Command("task", name)
	c.Dir = dir
	c.Env = append(os.Environ(), env...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func Available() bool {
	_, err := exec.LookPath("task")
	return err == nil
}
