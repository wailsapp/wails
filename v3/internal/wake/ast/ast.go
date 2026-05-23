package ast

type Taskfile struct {
	Version  string
	Includes map[string]*Include
	Vars     map[string]*Var
	Tasks    map[string]*Task
	Shopt    []string
	Silent   bool
	Location string
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
