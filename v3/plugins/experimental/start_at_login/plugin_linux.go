//go:build linux

package start_at_login

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type startTpl struct {
	Name    string
	Cmd     string
	Enabled bool
}

const tpl = `
[Desktop Entry]
Name={{.Name}}
Comment=Autostart service for {{.Name}}
Type=Application
Exec={{.Cmd}}
{{if .Enabled}}
Hidden=true
X-GNOME-Autostart-enabled=true
{{end}}
`

func (p *Plugin) init() error {
	return nil
}

func (p *Plugin) autoStartFileExists() bool {
	if _, err := os.Stat(p.autostartFile()); err == nil {
		return true
	}
	return false
}

func (p *Plugin) StartAtLogin(enabled bool) error {
	autostart := p.autostartFile()
	if !enabled && p.autoStartFileExists() {
		return os.Remove(autostart)
	} else if enabled && !p.autoStartFileExists() {
		p.createAutoStartFile(autostart, enabled)
	}

	return nil
}

func (p *Plugin) IsStartAtLogin() (bool, error) {
	result := p.autoStartFileExists()
	return result, nil
}

func (p *Plugin) autostartFile() string {
	homedir, _ := os.UserHomeDir()
	exe, _ := os.Executable()
	name := filepath.Base(exe)
	autostartFile := fmt.Sprintf("%s-autostart.desktop", name)
	return strings.Join([]string{homedir, ".config", "autostart", autostartFile}, "/")
}

func (p *Plugin) createAutoStartFile(filename string, enabled bool) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	defer file.Close()

	tmpl, err := template.New("autostart").Parse(tpl)
	if err != nil {
		panic(err)
	}
	exe, _ := os.Executable()
	input := startTpl{
		Name:    filepath.Base(exe),
		Cmd:     exe,
		Enabled: enabled,
	}
	err = tmpl.Execute(file, input)
	if err != nil {
		panic(err)
	}
}
