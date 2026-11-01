package ast

type Taskfile struct {
	Version  string
	Includes map[string]*Include
	Vars     map[string]*Var
	Tasks    map[string]*Task
	// Env is the Taskfile-level `env:` block: environment variables that are
	// added to every task's subprocess environment, overridden by task-local
	// `env:` on per-task collisions. Empty when the Taskfile declares no env.
	Env         map[string]string
	Shopt       []string
	Silent      bool
	Location    string
	Dotenv      []string
	Output      string
	Run         string
	Requires    *Requires
	Interval    string
	Unsupported map[string]any
}

type Requires struct {
	Preconditions []string
}

type Include struct {
	Taskfile    string
	Dir         string
	Optional    bool
	Internal    bool
	Aliases     []string
	Resolved    *Taskfile
	Vars        map[string]*Var
	Unsupported []string
}

type Var struct {
	Literal    bool
	Definition any
	Static     string
	Shell      string
	Ref        string
	Value      string
}

type Task struct {
	// Definition preserves fields the analyser does not interpret for exact stock comparison.
	Definition   any
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
