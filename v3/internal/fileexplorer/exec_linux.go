//go:build linux

package fileexplorer

import "strings"

var desktopStringEscapes = strings.NewReplacer(
	`\s`, " ", `\n`, "\n", `\t`, "\t", `\r`, "\r", `\\`, `\`,
)

// splitDesktopExec undoes string escapes and double quoting before field expansion.
// https://specifications.freedesktop.org/desktop-entry/latest/exec-variables.html
func splitDesktopExec(command string) []string {
	command = desktopStringEscapes.Replace(command)
	var args []string
	var word strings.Builder
	quoted, active := false, false
	for i := 0; i < len(command); i++ {
		c := command[i]
		switch {
		case c == '"':
			quoted = !quoted
			active = true
		case quoted && c == '\\' && i+1 < len(command):
			i++
			if strings.IndexByte("\"`$\\", command[i]) < 0 {
				return nil
			}
			word.WriteByte(command[i])
			active = true
		case !quoted && strings.IndexByte(" \t\r\n", c) >= 0:
			if active {
				args = append(args, word.String())
				word.Reset()
				active = false
			}
		default:
			word.WriteByte(c)
			active = true
		}
	}
	if quoted {
		return nil
	}
	if active {
		args = append(args, word.String())
	}
	if len(args) == 0 || args[0] == "" {
		return nil
	}
	return args
}
