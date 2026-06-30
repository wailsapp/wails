package cmds

import (
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Executor interface {
	Run() error
}

// Output is embedded in every command type to give the executor one uniform way
// to redirect a command's streams. The executor captures into a buffer by
// default and streams live when verbose; a nil writer discards (the os/exec
// default of /dev/null).
type Output struct {
	Stdout io.Writer
	Stderr io.Writer
}

// OutputSetter is implemented (via the embedded Output) by every command type,
// so the executor can set streams without a type switch.
type OutputSetter interface {
	SetOutput(stdout, stderr io.Writer)
}

func (o *Output) SetOutput(stdout, stderr io.Writer) {
	o.Stdout = stdout
	o.Stderr = stderr
}

func (o Output) apply(c *exec.Cmd) {
	c.Stdout = o.Stdout
	c.Stderr = o.Stderr
}

type ShellCmd struct {
	Output
	Cmd string
	Dir string
	Env []string
}

func (s *ShellCmd) Run() error {
	c := exec.Command("sh", "-c", s.Cmd)
	c.Dir = s.Dir
	c.Env = s.Env
	s.apply(c)
	return c.Run()
}

type RouteOptions struct {
	Dir string
	Env []string
}

func Route(cmdStr string, opts RouteOptions) Executor {
	routes := []struct {
		pattern *regexp.Regexp
		build   func([]string, RouteOptions) Executor
	}{
		{regexp.MustCompile(`^go\s+build\b`), buildGoBuild},
		{regexp.MustCompile(`^go\s+run\b`), buildGoRun},
		{regexp.MustCompile(`^go\s+test\b`), buildGoTest},
		{regexp.MustCompile(`^go\s+mod\s+tidy\b`), buildGoModTidy},
		{regexp.MustCompile(`^go\s+mod\s+(\S+)`), buildGoMod},
		{regexp.MustCompile(`^go\s+install\b`), buildGoInstall},
		{regexp.MustCompile(`^go\s+vet\b`), buildGoVet},
		{regexp.MustCompile(`^go\s+fmt\b`), buildGoFmt},
		{regexp.MustCompile(`^npm\s+install\b`), buildNpmInstall},
		{regexp.MustCompile(`^npm\s+run\s+(\S+)`), buildNpmRun},
		{regexp.MustCompile(`^bun\s+install\b`), buildBunInstall},
		{regexp.MustCompile(`^pnpm\s+install\b`), buildPnpmInstall},
		{regexp.MustCompile(`^yarn\s+install\b`), buildYarnInstall},
	}

	for _, r := range routes {
		if r.pattern.MatchString(cmdStr) {
			// Tokenise with shell-aware quoting so args like
			// `-ldflags="-w -s"` stay as a single token instead of being
			// split on the embedded space (the bug strings.Fields would
			// create). If quoting is unbalanced or otherwise too exotic to
			// parse, fall back to running through `sh -c` rather than risk
			// a wrong tokenisation.
			args := shellSplit(cmdStr)
			if args == nil {
				return &ShellCmd{Cmd: cmdStr, Dir: opts.Dir, Env: opts.Env}
			}
			ex := r.build(args, opts)
			if ex != nil {
				return ex
			}
		}
	}

	return &ShellCmd{Cmd: cmdStr, Dir: opts.Dir, Env: opts.Env}
}

// shellSplit tokenises a shell-style command line into argv-style tokens,
// respecting single- and double-quoted regions and `\` escapes. The quotes
// themselves are stripped. Adjacent quoted/unquoted fragments concatenate
// into one token (e.g. `-ldflags="-w -s"` → `-ldflags=-w -s` as ONE arg).
//
// Returns nil for unbalanced quoting; callers should fall back to `sh -c`
// in that case rather than risk mis-parsing the user's command.
func shellSplit(s string) []string {
	var (
		out    []string
		cur    strings.Builder
		inS    bool // inside single quote
		inD    bool // inside double quote
		active bool // current token has at least one (possibly empty-quoted) char
	)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\\' && !inS && i+1 < len(s):
			cur.WriteByte(s[i+1])
			i++
			active = true
		case c == '\'' && !inD:
			inS = !inS
			active = true
		case c == '"' && !inS:
			inD = !inD
			active = true
		case (c == ' ' || c == '\t') && !inS && !inD:
			if active {
				out = append(out, cur.String())
				cur.Reset()
				active = false
			}
		default:
			cur.WriteByte(c)
			active = true
		}
	}
	if inS || inD {
		return nil
	}
	if active {
		out = append(out, cur.String())
	}
	return out
}

func buildGoBuild(args []string, opts RouteOptions) Executor {
	parsed := parseGoBuildArgs(args)
	cmd := GoBuild(parsed)
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func parseGoBuildArgs(args []string) GoBuildOptions {
	opts := GoBuildOptions{}
	i := 2
	for i < len(args) {
		switch args[i] {
		case "-o":
			if i+1 < len(args) {
				opts.Output = args[i+1]
				i += 2
				continue
			}
		case "-tags":
			if i+1 < len(args) {
				opts.Tags = strings.Split(args[i+1], " ")
				i += 2
				continue
			}
		case "-ldflags":
			if i+1 < len(args) {
				opts.Ldflags = args[i+1]
				i += 2
				continue
			}
		case "-gcflags":
			if i+1 < len(args) {
				opts.Gcflags = args[i+1]
				i += 2
				continue
			}
		case "-trimpath":
			opts.Trimpath = true
		case "-race":
			opts.Race = true
		case "-buildvcs=true":
			opts.Buildvcs = true
		case "-buildvcs=false":
			opts.Buildvcs = false
		case "-mod":
			if i+1 < len(args) {
				opts.Mod = args[i+1]
				i += 2
				continue
			}
		default:
			if !strings.HasPrefix(args[i], "-") {
				opts.Package = args[i]
			} else {
				opts.ExtraArgs = append(opts.ExtraArgs, args[i])
			}
		}
		i++
	}
	return opts
}

func buildGoRun(args []string, opts RouteOptions) Executor {
	parsed := parseGoRunArgs(args)
	cmd := GoRun(parsed)
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func parseGoRunArgs(args []string) GoRunOptions {
	opts := GoRunOptions{}
	i := 2
	for i < len(args) {
		switch args[i] {
		case "-tags":
			if i+1 < len(args) {
				opts.Tags = strings.Split(args[i+1], " ")
				i += 2
				continue
			}
		case "-ldflags":
			if i+1 < len(args) {
				opts.Ldflags = args[i+1]
				i += 2
				continue
			}
		case "-trimpath":
			opts.Trimpath = true
		default:
			opts.Args = append(opts.Args, args[i])
		}
		i++
	}
	return opts
}

func buildGoTest(args []string, opts RouteOptions) Executor {
	parsed := parseGoTestArgs(args)
	cmd := GoTest(parsed)
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func parseGoTestArgs(args []string) GoTestOptions {
	opts := GoTestOptions{}
	i := 2
	for i < len(args) {
		switch args[i] {
		case "-v":
			opts.Verbose = true
		case "-cover":
			opts.Cover = true
		case "-race":
			opts.Race = true
		case "-tags":
			if i+1 < len(args) {
				opts.Tags = strings.Split(args[i+1], " ")
				i += 2
				continue
			}
		case "-count":
			if i+1 < len(args) {
				// Parse the full integer, not just the first character.
				// The previous implementation read args[i+1][0]-'0', so
				// `-count 10` became 1, `-count 100` became 1, etc.
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					opts.Count = n
				}
				i += 2
				continue
			}
		case "-run":
			if i+1 < len(args) {
				opts.Run = args[i+1]
				i += 2
				continue
			}
		default:
			if !strings.HasPrefix(args[i], "-") {
				opts.Package = args[i]
			} else {
				opts.Extra = append(opts.Extra, args[i])
			}
		}
		i++
	}
	return opts
}

func buildGoModTidy(args []string, opts RouteOptions) Executor {
	cmd := GoModTidy()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildGoMod(args []string, opts RouteOptions) Executor {
	subcommand := args[2]
	modArgs := args[3:]
	cmd := GoMod(GoModOptions{Subcommand: subcommand, Args: modArgs})
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildGoInstall(args []string, opts RouteOptions) Executor {
	pkg := ""
	if len(args) > 2 {
		pkg = args[2]
	}
	cmd := GoInstall(GoInstallOptions{Package: pkg})
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildGoVet(args []string, opts RouteOptions) Executor {
	cmd := GoVet()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildGoFmt(args []string, opts RouteOptions) Executor {
	cmd := GoFmt()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildNpmInstall(args []string, opts RouteOptions) Executor {
	cmd := NpmInstall()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildNpmRun(args []string, opts RouteOptions) Executor {
	if len(args) >= 3 {
		cmd := NpmRun(args[2])
		cmd.Dir = opts.Dir
		cmd.Env = opts.Env
		return cmd
	}
	return nil
}

func buildBunInstall(args []string, opts RouteOptions) Executor {
	cmd := BunInstall()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildPnpmInstall(args []string, opts RouteOptions) Executor {
	cmd := PnpmInstall()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}

func buildYarnInstall(args []string, opts RouteOptions) Executor {
	cmd := YarnInstall()
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	return cmd
}
