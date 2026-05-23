package cmds

import (
	"os"
	"os/exec"
	"path/filepath"
)

func DetectPackageManager(dir string) string {
	lockfiles := map[string]string{
		"package-lock.json": "npm",
		"bun.lock":          "bun",
		"bun.lockb":         "bun",
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
	}
	for name, pm := range lockfiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return pm
		}
	}
	return "npm"
}

type NpmInstallCmd struct {
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
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

type NpmRunCmd struct {
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
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

type BunInstallCmd struct {
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
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

type PnpmInstallCmd struct {
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
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

type YarnInstallCmd struct {
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
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}
