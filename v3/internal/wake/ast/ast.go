package ast

type Taskfile struct {
	Version  string
	Includes map[string]*Include
	Vars     map[string]*Var
	Tasks    map[string]*Task
	// Env is the Taskfile-level `env:` block: environment variables that are
	// added to every task's subprocess environment, overridden by task-local
	// `env:` on per-task collisions. Empty when the Taskfile declares no env.
	Env      map[string]string
	Shopt    []string
	Silent   bool
	Location string
	Dotenv   []string
	Output   string
	Run      string
	Requires *Requires
	Interval string
}

type Requires struct {
	Preconditions []string
}

type Include struct {
	Taskfile string
	Dir      string
	Optional bool
	Internal bool
	Aliases  []string
	Resolved *Taskfile
}

type Var struct {
	Static string
	Shell  string
	Ref    string
	Value  string
}

type Task struct {
	Name         string
	Summary      string
	Desc         string
	Aliases      []string
	Label        string
	Dir          string
	Silent       bool
	Internal     bool
	Interactive  bool
	Prompt       string
	Vars         map[string]*Var
	Deps         []*Dep
	Cmds         []*Cmd
	Platforms    []string
	Sources      []string
	Generates    []string
	Status       []string
	Precondition []*Precondition
	Env          map[string]string
	Method       string
	Run          string
	Short        string
	Defer        []string
	Interval     string
}

type Dep struct {
	Task   string
	Vars   map[string]*Var
	Silent bool
}

type Cmd struct {
	Cmd         string
	Task        string
	For         *ForLoop
	Silent      bool
	IgnoreError bool
	Platforms   []string
	Vars        map[string]*Var
}

type ForLoop struct {
	Var   string
	Items []string
	Task  string
	Vars  map[string]*Var
}

type Precondition struct {
	Sh  string
	Msg string
}

func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}
	clone := *t

	if t.Aliases != nil {
		clone.Aliases = make([]string, len(t.Aliases))
		copy(clone.Aliases, t.Aliases)
	}
	if t.Vars != nil {
		clone.Vars = make(map[string]*Var, len(t.Vars))
		for k, v := range t.Vars {
			vClone := *v
			clone.Vars[k] = &vClone
		}
	}
	if t.Deps != nil {
		clone.Deps = make([]*Dep, len(t.Deps))
		for i, d := range t.Deps {
			dClone := *d
			if d.Vars != nil {
				dClone.Vars = make(map[string]*Var, len(d.Vars))
				for k, v := range d.Vars {
					vClone := *v
					dClone.Vars[k] = &vClone
				}
			}
			clone.Deps[i] = &dClone
		}
	}
	if t.Cmds != nil {
		clone.Cmds = make([]*Cmd, len(t.Cmds))
		for i, c := range t.Cmds {
			cClone := *c
			if c.Vars != nil {
				cClone.Vars = make(map[string]*Var, len(c.Vars))
				for k, v := range c.Vars {
					vClone := *v
					cClone.Vars[k] = &vClone
				}
			}
			if c.For != nil {
				forClone := *c.For
				if c.For.Items != nil {
					forClone.Items = make([]string, len(c.For.Items))
					copy(forClone.Items, c.For.Items)
				}
				if c.For.Vars != nil {
					forClone.Vars = make(map[string]*Var, len(c.For.Vars))
					for k, v := range c.For.Vars {
						vClone := *v
						forClone.Vars[k] = &vClone
					}
				}
				cClone.For = &forClone
			}
			clone.Cmds[i] = &cClone
		}
	}
	if t.Platforms != nil {
		clone.Platforms = make([]string, len(t.Platforms))
		copy(clone.Platforms, t.Platforms)
	}
	if t.Sources != nil {
		clone.Sources = make([]string, len(t.Sources))
		copy(clone.Sources, t.Sources)
	}
	if t.Generates != nil {
		clone.Generates = make([]string, len(t.Generates))
		copy(clone.Generates, t.Generates)
	}
	if t.Status != nil {
		clone.Status = make([]string, len(t.Status))
		copy(clone.Status, t.Status)
	}
	if t.Precondition != nil {
		clone.Precondition = make([]*Precondition, len(t.Precondition))
		for i, p := range t.Precondition {
			pClone := *p
			clone.Precondition[i] = &pClone
		}
	}
	if t.Env != nil {
		clone.Env = make(map[string]string, len(t.Env))
		for k, v := range t.Env {
			clone.Env[k] = v
		}
	}
	if t.Defer != nil {
		clone.Defer = make([]string, len(t.Defer))
		copy(clone.Defer, t.Defer)
	}

	return &clone
}
