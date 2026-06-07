package cmds

import (
	"os"
	"os/exec"
	"path/filepath"
)

func DetectPackageManager(dir string) string {
	// Ordered slice, not a map: when more than one lockfile is present
	// (e.g. a repo migrated from npm to pnpm but kept the old lockfile)
	// the choice has to be deterministic across runs. Map iteration order
	// is randomised in Go, so the previous implementation could return
	// "npm" on one run and "pnpm" on the next from the same tree.
	for _, lf := range []struct {
		name, pm string
	}{
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
		{"bun.lock", "bun"},
		{"bun.lockb", "bun"},
		{"package-lock.json", "npm"},
	} {
		if _, err := os.Stat(filepath.Join(dir, lf.name)); err == nil {
			return lf.pm
		}
	}
	return "npm"
}

type NpmInstallCmd struct {
	Output
	Dir string
	Env []string
}

func NpmInstall() *NpmInstallCmd {
	return &NpmInstallCmd{}
}

func (n *NpmInstallCmd) Run() error {
	c := exec.Command("npm", "install")
	c.Dir = n.Dir
	c.Env = n.Env
	n.apply(c)
	return c.Run()
}

type NpmRunCmd struct {
	Output
	Script string
	Dir    string
	Env    []string
}

func NpmRun(script string) *NpmRunCmd {
	return &NpmRunCmd{Script: script}
}

func (n *NpmRunCmd) Run() error {
	c := exec.Command("npm", "run", n.Script)
	c.Dir = n.Dir
	c.Env = n.Env
	n.apply(c)
	return c.Run()
}

type BunInstallCmd struct {
	Output
	Dir string
	Env []string
}

func BunInstall() *BunInstallCmd {
	return &BunInstallCmd{}
}

func (b *BunInstallCmd) Run() error {
	c := exec.Command("bun", "install")
	c.Dir = b.Dir
	c.Env = b.Env
	b.apply(c)
	return c.Run()
}

type PnpmInstallCmd struct {
	Output
	Dir string
	Env []string
}

func PnpmInstall() *PnpmInstallCmd {
	return &PnpmInstallCmd{}
}

func (p *PnpmInstallCmd) Run() error {
	c := exec.Command("pnpm", "install")
	c.Dir = p.Dir
	c.Env = p.Env
	p.apply(c)
	return c.Run()
}

type YarnInstallCmd struct {
	Output
	Dir string
	Env []string
}

func YarnInstall() *YarnInstallCmd {
	return &YarnInstallCmd{}
}

func (y *YarnInstallCmd) Run() error {
	c := exec.Command("yarn", "install")
	c.Dir = y.Dir
	c.Env = y.Env
	y.apply(c)
	return c.Run()
}
