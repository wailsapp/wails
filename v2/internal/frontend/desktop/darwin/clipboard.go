//go:build darwin

package darwin

import (
	"os"
	"os/exec"
)

// ensureUTF8Env returns the current environment with LANG set to en_US.UTF-8
// if it is not already set. This is needed because packaged macOS apps do not
// inherit the terminal's LANG variable, causing pbpaste/pbcopy to default to
// an ASCII-compatible encoding that mangles non-ASCII text.
func ensureUTF8Env() []string {
	env := os.Environ()
	if _, ok := os.LookupEnv("LANG"); !ok {
		env = append(env, "LANG=en_US.UTF-8")
	}
	return env
}

func (f *Frontend) ClipboardGetText() (string, error) {
	pasteCmd := exec.Command("pbpaste")
	pasteCmd.Env = ensureUTF8Env()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (f *Frontend) ClipboardSetText(text string) error {
	copyCmd := exec.Command("pbcopy")
	copyCmd.Env = ensureUTF8Env()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
