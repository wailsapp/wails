package cmds

import (
	"os"
	"os/exec"
)

type GoBuildOptions struct {
	Output    string
	Tags      []string
	Ldflags   string
	Gcflags   string
	Trimpath  bool
	Buildvcs  bool
	Package   string
	Race      bool
	Mod       string
	ExtraArgs []string
}

type GoBuildCmd struct {
	Opts GoBuildOptions
	Dir  string
	Env  []string
}

func GoBuild(opts GoBuildOptions) *GoBuildCmd {
	return &GoBuildCmd{Opts: opts}
}

func (g *GoBuildCmd) Run() error {
	args := []string{"build"}
	if g.Opts.Output != "" {
		args = append(args, "-o", g.Opts.Output)
	}
	if len(g.Opts.Tags) > 0 {
		args = append(args, "-tags", joinTags(g.Opts.Tags))
	}
	if g.Opts.Ldflags != "" {
		args = append(args, "-ldflags", g.Opts.Ldflags)
	}
	if g.Opts.Gcflags != "" {
		args = append(args, "-gcflags", g.Opts.Gcflags)
	}
	if g.Opts.Trimpath {
		args = append(args, "-trimpath")
	}
	if !g.Opts.Buildvcs {
		args = append(args, "-buildvcs=false")
	}
	if g.Opts.Race {
		args = append(args, "-race")
	}
	if g.Opts.Mod != "" {
		args = append(args, "-mod", g.Opts.Mod)
	}
	args = append(args, g.Opts.ExtraArgs...)
	if g.Opts.Package != "" {
		args = append(args, g.Opts.Package)
	}

	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoRunOptions struct {
	Tags     []string
	Ldflags  string
	Trimpath bool
	Args     []string
}

type GoRunCmd struct {
	Opts GoRunOptions
	Dir  string
	Env  []string
}

func GoRun(opts GoRunOptions) *GoRunCmd {
	return &GoRunCmd{Opts: opts}
}

func (g *GoRunCmd) Run() error {
	args := []string{"run"}
	if len(g.Opts.Tags) > 0 {
		args = append(args, "-tags", joinTags(g.Opts.Tags))
	}
	if g.Opts.Ldflags != "" {
		args = append(args, "-ldflags", g.Opts.Ldflags)
	}
	if g.Opts.Trimpath {
		args = append(args, "-trimpath")
	}
	args = append(args, g.Opts.Args...)

	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoTestOptions struct {
	Package string
	Verbose bool
	Cover   bool
	Race    bool
	Tags    []string
	Count   int
	Run     string
	Extra   []string
}

type GoTestCmd struct {
	Opts GoTestOptions
	Dir  string
	Env  []string
}

func GoTest(opts GoTestOptions) *GoTestCmd {
	return &GoTestCmd{Opts: opts}
}

func (g *GoTestCmd) Run() error {
	args := []string{"test"}
	if g.Opts.Verbose {
		args = append(args, "-v")
	}
	if g.Opts.Cover {
		args = append(args, "-cover")
	}
	if g.Opts.Race {
		args = append(args, "-race")
	}
	if len(g.Opts.Tags) > 0 {
		args = append(args, "-tags", joinTags(g.Opts.Tags))
	}
	if g.Opts.Count > 0 {
		args = append(args, "-count", string(rune('0'+g.Opts.Count)))
	}
	if g.Opts.Run != "" {
		args = append(args, "-run", g.Opts.Run)
	}
	args = append(args, g.Opts.Extra...)
	if g.Opts.Package != "" {
		args = append(args, g.Opts.Package)
	}

	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoModTidyCmd struct {
	Dir string
	Env []string
}

func GoModTidy() *GoModTidyCmd {
	return &GoModTidyCmd{}
}

func (g *GoModTidyCmd) Run() error {
	c := exec.Command("go", "mod", "tidy")
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoInstallOptions struct {
	Package string
}

type GoInstallCmd struct {
	Opts GoInstallOptions
	Dir  string
	Env  []string
}

func GoInstall(opts GoInstallOptions) *GoInstallCmd {
	return &GoInstallCmd{Opts: opts}
}

func (g *GoInstallCmd) Run() error {
	args := []string{"install"}
	if g.Opts.Package != "" {
		args = append(args, g.Opts.Package)
	}
	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoModOptions struct {
	Subcommand string
	Args       []string
}

type GoModCmd struct {
	Opts GoModOptions
	Dir  string
	Env  []string
}

func GoMod(opts GoModOptions) *GoModCmd {
	return &GoModCmd{Opts: opts}
}

func (g *GoModCmd) Run() error {
	args := append([]string{"mod", g.Opts.Subcommand}, g.Opts.Args...)
	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoVetCmd struct {
	Package string
	Dir     string
	Env     []string
}

func GoVet() *GoVetCmd {
	return &GoVetCmd{}
}

func (g *GoVetCmd) Run() error {
	args := []string{"vet"}
	if g.Package != "" {
		args = append(args, g.Package)
	}
	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type GoFmtCmd struct {
	Write bool
	Paths []string
	Dir   string
	Env   []string
}

func GoFmt() *GoFmtCmd {
	return &GoFmtCmd{}
}

func (g *GoFmtCmd) Run() error {
	args := []string{"fmt"}
	if g.Write {
		args = append(args, "-w")
	}
	args = append(args, g.Paths...)
	c := exec.Command("go", args...)
	c.Dir = g.Dir
	c.Env = g.Env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func joinTags(tags []string) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += " "
		}
		result += t
	}
	return result
}
