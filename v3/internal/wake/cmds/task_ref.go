package cmds

import (
	"fmt"
)

type TaskRefCmd struct {
	TaskName string
	Dir      string
	Env      []string
}

func (t *TaskRefCmd) Run() error {
	return fmt.Errorf("wake: task reference %q should be handled by executor, not routed", t.TaskName)
}
