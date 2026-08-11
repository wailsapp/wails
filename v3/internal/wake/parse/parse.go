package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"gopkg.in/yaml.v3"
)

func Parse(path string) (*ast.Taskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wake: read %s: %w", path, err)
	}

	var raw yaml.Node
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("wake: parse %s: %w", path, err)
	}

	if len(raw.Content) == 0 {
		return nil, fmt.Errorf("wake: empty Taskfile: %s", path)
	}

	tf := &ast.Taskfile{Location: path}
	root := raw.Content[0]

	for i := 0; i < len(root.Content); i += 2 {
		key := root.Content[i].Value
		val := root.Content[i+1]
		switch key {
		case "version":
			tf.Version = val.Value
		case "includes":
			tf.Includes, err = parseIncludes(val)
			if err != nil {
				return nil, err
			}
		case "vars":
			tf.Vars, err = parseVars(val)
			if err != nil {
				return nil, err
			}
		case "env":
			tf.Env, err = parseEnv(val)
			if err != nil {
				return nil, err
			}
		case "tasks":
			tf.Tasks, err = parseTasks(val)
			if err != nil {
				return nil, err
			}
		case "shopt":
			tf.Shopt = parseStringList(val)
		case "silent":
			tf.Silent = val.Value == "true"
		case "dotenv":
			tf.Dotenv = parseStringList(val)
		case "output":
			tf.Output = val.Value
		case "run":
			tf.Run = val.Value
		case "interval":
			tf.Interval = val.Value
		case "requires":
			tf.Requires = parseRequires(val)
		}
	}

	if tf.Version != "3" {
		return nil, fmt.Errorf("wake: unsupported Taskfile version %q in %s (only v3 supported)", tf.Version, path)
	}

	// A Taskfile without a `tasks:` key leaves Tasks nil; downstream code
	// (merges, builtins, DAG) assumes a non-nil map, so normalise it here.
	if tf.Tasks == nil {
		tf.Tasks = make(map[string]*ast.Task)
	}

	return tf, nil
}

func parseIncludes(node *yaml.Node) (map[string]*ast.Include, error) {
	includes := make(map[string]*ast.Include)
	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		inc := &ast.Include{}
		val := node.Content[i+1]

		if val.Kind == yaml.ScalarNode {
			inc.Taskfile = val.Value
		} else {
			for j := 0; j < len(val.Content); j += 2 {
				k := val.Content[j].Value
				v := val.Content[j+1]
				switch k {
				case "taskfile":
					inc.Taskfile = v.Value
				case "dir":
					inc.Dir = v.Value
				case "optional":
					inc.Optional = v.Value == "true"
				case "internal":
					inc.Internal = v.Value == "true"
				case "aliases":
					inc.Aliases = parseStringList(v)
				}
			}
		}
		includes[name] = inc
	}
	return includes, nil
}

func parseVars(node *yaml.Node) (map[string]*ast.Var, error) {
	vars := make(map[string]*ast.Var)
	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		val := node.Content[i+1]
		vr := &ast.Var{}

		if val.Kind == yaml.ScalarNode {
			vr.Static = val.Value
		} else if val.Kind == yaml.MappingNode {
			for j := 0; j < len(val.Content); j += 2 {
				k := val.Content[j].Value
				v := val.Content[j+1]
				switch k {
				case "sh":
					vr.Shell = v.Value
				case "ref":
					vr.Ref = v.Value
				default:
					vr.Static = v.Value
				}
			}
		} else {
			vr.Static = val.Value
		}
		vars[name] = vr
	}
	return vars, nil
}

func parseTasks(node *yaml.Node) (map[string]*ast.Task, error) {
	tasks := make(map[string]*ast.Task)
	for i := 0; i < len(node.Content); i += 2 {
		name := node.Content[i].Value
		task, err := parseTask(name, node.Content[i+1])
		if err != nil {
			return nil, err
		}
		tasks[name] = task
	}
	return tasks, nil
}

func parseTask(name string, node *yaml.Node) (*ast.Task, error) {
	task := &ast.Task{Name: name}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		val := node.Content[i+1]
		switch key {
		case "summary":
			task.Summary = val.Value
		case "desc":
			task.Desc = val.Value
		case "aliases":
			task.Aliases = parseStringList(val)
		case "label":
			task.Label = val.Value
		case "dir":
			task.Dir = val.Value
		case "silent":
			task.Silent = val.Value == "true"
		case "internal":
			task.Internal = val.Value == "true"
		case "interactive":
			task.Interactive = val.Value == "true"
		case "prompt":
			task.Prompt = val.Value
		case "vars":
			vars, err := parseVars(val)
			if err != nil {
				return nil, err
			}
			task.Vars = vars
		case "deps":
			deps, err := parseDeps(val)
			if err != nil {
				return nil, err
			}
			task.Deps = deps
		case "cmds":
			cmds, err := parseCmds(val)
			if err != nil {
				return nil, err
			}
			task.Cmds = cmds
		case "platforms":
			task.Platforms = parseStringList(val)
		case "sources":
			task.Sources = parseStringList(val)
		case "generates":
			task.Generates = parseStringList(val)
		case "status":
			task.Status = parseStringList(val)
		case "preconditions":
			preconds, err := parsePreconditions(val)
			if err != nil {
				return nil, err
			}
			task.Precondition = preconds
		case "env":
			env, err := parseEnv(val)
			if err != nil {
				return nil, err
			}
			task.Env = env
		case "method":
			task.Method = val.Value
		case "run":
			task.Run = val.Value
		case "short":
			task.Short = val.Value
		case "defer":
			task.Defer = parseStringList(val)
		case "interval":
			task.Interval = val.Value
		}
	}
	return task, nil
}

func parseDeps(node *yaml.Node) ([]*ast.Dep, error) {
	var deps []*ast.Dep
	for _, item := range node.Content {
		dep := &ast.Dep{}
		if item.Kind == yaml.ScalarNode {
			dep.Task = item.Value
		} else {
			for i := 0; i < len(item.Content); i += 2 {
				k := item.Content[i].Value
				v := item.Content[i+1]
				switch k {
				case "task":
					dep.Task = v.Value
				case "silent":
					dep.Silent = v.Value == "true"
				case "vars":
					vars, err := parseVars(v)
					if err != nil {
						return nil, err
					}
					dep.Vars = vars
				}
			}
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func parseCmds(node *yaml.Node) ([]*ast.Cmd, error) {
	var cmds []*ast.Cmd
	for _, item := range node.Content {
		cmd := &ast.Cmd{}
		if item.Kind == yaml.ScalarNode {
			cmd.Cmd = item.Value
		} else if item.Kind == yaml.MappingNode {
			hasStruct := false
			for i := 0; i < len(item.Content); i += 2 {
				k := item.Content[i].Value
				v := item.Content[i+1]
				switch k {
				case "cmd":
					cmd.Cmd = v.Value
					hasStruct = true
				case "task":
					cmd.Task = v.Value
					hasStruct = true
				case "for":
					forLoop, err := parseFor(v)
					if err != nil {
						return nil, err
					}
					cmd.For = forLoop
					hasStruct = true
				case "silent":
					cmd.Silent = v.Value == "true"
				case "ignore_error":
					cmd.IgnoreError = v.Value == "true"
				case "platforms":
					cmd.Platforms = parseStringList(v)
				case "vars":
					vars, err := parseVars(v)
					if err != nil {
						return nil, err
					}
					cmd.Vars = vars
				}
			}
			if !hasStruct && cmd.Cmd == "" && cmd.Task == "" && cmd.For == nil {
				cmd.Cmd = item.Content[0].Value
			}
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func parseFor(node *yaml.Node) (*ast.ForLoop, error) {
	fl := &ast.ForLoop{}
	switch node.Kind {
	case yaml.ScalarNode:
		// `for: someVar` — the loop iterates the value of `someVar`,
		// space-split at execution time. The variable name lands in
		// ForLoop.Var; runForLoop knows how to expand it.
		fl.Var = node.Value
	case yaml.SequenceNode:
		// `for: [a, b, c]` — inline item list.
		for _, it := range node.Content {
			fl.Items = append(fl.Items, it.Value)
		}
	case yaml.MappingNode:
		// `for: { var: …, task: …, vars: … }` — the explicit mapping form.
		// Also accepts a nested `items:` list inside the mapping; previously
		// this was silently dropped.
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i].Value
			v := node.Content[i+1]
			switch k {
			case "var":
				fl.Var = v.Value
			case "items":
				if v.Kind == yaml.SequenceNode {
					for _, it := range v.Content {
						fl.Items = append(fl.Items, it.Value)
					}
				}
			case "task":
				fl.Task = v.Value
			case "vars":
				vars, err := parseVars(v)
				if err != nil {
					return nil, err
				}
				fl.Vars = vars
			}
		}
	}
	return fl, nil
}

func parsePreconditions(node *yaml.Node) ([]*ast.Precondition, error) {
	var preconds []*ast.Precondition
	for _, item := range node.Content {
		pc := &ast.Precondition{}
		for i := 0; i < len(item.Content); i += 2 {
			k := item.Content[i].Value
			v := item.Content[i+1]
			switch k {
			case "sh":
				pc.Sh = v.Value
			case "msg":
				pc.Msg = v.Value
			}
		}
		preconds = append(preconds, pc)
	}
	return preconds, nil
}

func parseRequires(node *yaml.Node) *ast.Requires {
	req := &ast.Requires{}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i].Value
		v := node.Content[i+1]
		switch k {
		case "preconditions":
			req.Preconditions = parseStringList(v)
		}
	}
	return req
}

func parseEnv(node *yaml.Node) (map[string]string, error) {
	env := make(map[string]string)
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i].Value
		v := node.Content[i+1]
		env[k] = v.Value
	}
	return env, nil
}

func parseStringList(node *yaml.Node) []string {
	var list []string
	for _, item := range node.Content {
		if item.Kind == yaml.ScalarNode {
			list = append(list, item.Value)
		} else if item.Kind == yaml.MappingNode && len(item.Content) >= 2 {
			key := item.Content[0].Value
			val := item.Content[1].Value
			if key == "exclude" {
				list = append(list, "exclude: "+val)
			} else {
				list = append(list, val)
			}
		}
	}
	return list
}

func ResolveIncludes(tf *ast.Taskfile) error {
	resolved := make(map[string]*ast.Taskfile)
	return resolveIncludes(tf, resolved)
}

func resolveIncludes(tf *ast.Taskfile, resolved map[string]*ast.Taskfile) error {
	baseDir := filepath.Dir(tf.Location)
	for name, inc := range tf.Includes {
		incPath := inc.Taskfile
		if !filepath.IsAbs(incPath) {
			incPath = filepath.Join(baseDir, incPath)
		}

		incPath = resolveTaskfilePath(incPath)

		if inc.Dir == "" {
			inc.Dir = filepath.Dir(incPath)
		} else if !filepath.IsAbs(inc.Dir) {
			inc.Dir = filepath.Join(baseDir, inc.Dir)
		}

		if _, err := os.Stat(incPath); os.IsNotExist(err) {
			if inc.Optional {
				continue
			}
			return fmt.Errorf("wake: include %q not found: %s", name, incPath)
		}

		cleanPath, _ := filepath.EvalSymlinks(incPath)

		var tfResolved *ast.Taskfile
		if existing, ok := resolved[cleanPath]; ok {
			tfResolved = existing
		} else {
			var err error
			tfResolved, err = Parse(incPath)
			if err != nil {
				return fmt.Errorf("wake: include %q: %w", name, err)
			}
			resolved[cleanPath] = tfResolved

			if err := resolveIncludes(tfResolved, resolved); err != nil {
				return err
			}
		}

		inc.Resolved = tfResolved

		for taskName, task := range tfResolved.Tasks {
			namespaced := name + ":" + taskName
			cloned := task.Clone()
			// Inject the included file's top-level vars into each cloned task so
			// they're available during execution (e.g. CROSS_IMAGE in preconditions).
			// Task-specific vars have higher priority and overwrite these defaults.
			if len(tfResolved.Vars) > 0 {
				merged := make(map[string]*ast.Var, len(tfResolved.Vars)+len(cloned.Vars))
				// Deep-copy the included vars (consistent with Task.Clone) so the
				// later in-place mutation by ResolveVars/ResolveAllVarShells can't
				// leak across the tasks that share this included file.
				for k, v := range tfResolved.Vars {
					vClone := *v
					merged[k] = &vClone
				}
				for k, v := range cloned.Vars {
					merged[k] = v
				}
				cloned.Vars = merged
			}
			tf.Tasks[namespaced] = cloned
		}
	}
	return nil
}

func resolveTaskfilePath(path string) string {
	if _, err := os.Stat(path); err == nil {
		info, _ := os.Stat(path)
		if info.IsDir() {
			candidates := []string{
				filepath.Join(path, "Taskfile.yml"),
				filepath.Join(path, "Taskfile.yaml"),
			}
			for _, c := range candidates {
				if _, err := os.Stat(c); err == nil {
					return c
				}
			}
			return filepath.Join(path, "Taskfile.yml")
		}
	}
	return path
}

func PopulateBuiltins(tf *ast.Taskfile) {
	if tf.Vars == nil {
		tf.Vars = make(map[string]*ast.Var)
	}
	wd, _ := os.Getwd()
	builtins := map[string]string{
		"OS":                 runtime.GOOS,
		"ARCH":               runtime.GOARCH,
		"OSFAMILY":           osFamily(runtime.GOOS),
		"NUMCPU":             fmt.Sprintf("%d", runtime.NumCPU()),
		"ROOT_DIR":           filepath.Dir(tf.Location),
		"TASKFILE":           tf.Location,
		"TASKFILE_DIR":       filepath.Dir(tf.Location),
		"exeExt":             exeExt(runtime.GOOS),
		"BUILD_TAGS":         os.Getenv("BUILD_TAGS"),
		"USER_WORKING_DIR":   wd,
	}
	for k, v := range builtins {
		if _, ok := tf.Vars[k]; !ok {
			tf.Vars[k] = &ast.Var{Static: v, Value: v}
		}
	}
}

func osFamily(goos string) string {
	switch goos {
	case "linux", "darwin", "freebsd", "openbsd", "netbsd":
		return "unix"
	case "windows":
		return "windows"
	default:
		return goos
	}
}

func exeExt(goos string) string {
	if goos == "windows" {
		return ".exe"
	}
	return ""
}



func ResolveVars(vars map[string]*ast.Var) error {
	resolved := make(map[string]bool)
	for name := range vars {
		if err := resolveVar(name, vars, resolved, make(map[string]bool)); err != nil {
			return err
		}
	}
	return nil
}

func resolveVar(name string, vars map[string]*ast.Var, done, visiting map[string]bool) error {
	if done[name] {
		return nil
	}
	if visiting[name] {
		return fmt.Errorf("wake: variable reference cycle detected involving %q", name)
	}
	visiting[name] = true

	vr := vars[name]
	if vr == nil {
		return nil
	}

	if vr.Shell != "" && vr.Value == "" {
		// Intentionally leave Value empty here — ResolveAllVarShells (called
		// separately from wake.go) detects unresolved shell vars by
		// Shell != "" && Value == "" and executes the command. Setting
		// Value = Shell here would short-circuit that detection and the
		// raw command string would propagate into templates instead of
		// the command's stdout.
	} else if vr.Ref != "" {
		refName := strings.TrimPrefix(vr.Ref, ".")
		if err := resolveVar(refName, vars, done, visiting); err != nil {
			return err
		}
		ref := vars[refName]
		if ref != nil {
			vr.Value = ref.Value
		}
	} else if vr.Value == "" && !strings.Contains(vr.Static, "{{") {
		vr.Value = vr.Static
	}

	delete(visiting, name)
	done[name] = true
	return nil
}

// ExpandVarTemplates does a fixed-point pass over vars, evaluating any var
// whose Static is a template against the *same* var map. Used for vars at a
// scope where every reference is in-scope — typically the top-level Taskfile
// vars after [ResolveVars] has settled the static/shell/ref ones. Task-level
// vars must NOT go through this: a task-local template often references
// root-level vars that aren't yet in scope, and committing a half-expanded
// result here would freeze the wrong value before mergeVars at execution
// time can finish the job.
//
// Iterates up to 10 times for chained defaults like
// `OUTPUT: '{{ .OUTPUT | default .DEFAULT_OUTPUT }}'` where one var's
// resolved Value feeds the next.
func ExpandVarTemplates(vars map[string]*ast.Var) {
	for iter := 0; iter < 10; iter++ {
		changed := false
		for _, vr := range vars {
			if vr.Value != "" || !strings.Contains(vr.Static, "{{") {
				continue
			}
			expanded := ExpandTemplates(vr.Static, vars)
			if strings.Contains(expanded, "{{") {
				continue // template still references something unresolved
			}
			vr.Value = expanded
			changed = true
		}
		if !changed {
			break
		}
	}
}
