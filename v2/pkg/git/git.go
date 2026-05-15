package git

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

func gitcommand() string {
	gitcommand := "git"
	if runtime.GOOS == "windows" {
		gitcommand = "git.exe"
	}

	return gitcommand
}

// IsInstalled returns true if git is installed for the given platform
func IsInstalled() bool {
	return shell.CommandExists(gitcommand())
}

// Email tries to retrieve the
func Email() (string, error) {
	stdout, _, err := shell.RunCommand(".", gitcommand(), "config", "user.email")
	return stdout, err
}

// Name tries to retrieve the
func Name() (string, error) {
	errMsg := "failed to retrieve git user name: %w"
	stdout, _, err := shell.RunCommand(".", gitcommand(), "config", "user.name")
	if err != nil {
		return "", fmt.Errorf(errMsg, err)
	}
	name := strings.TrimSpace(stdout)
	return EscapeName(name)
}

func EscapeName(str string) (string, error) {
	b, err := json.Marshal(str)
	if err != nil {
		return "", err
	}
	// Remove the surrounding quotes
	escaped := string(b[1 : len(b)-1])

	// Check if username is JSON compliant
	var js json.RawMessage
	jsonVal := fmt.Sprintf(`{"name": "%s"}`, escaped)
	err = json.Unmarshal([]byte(jsonVal), &js)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve git user name: %w", err)
	}
	return escaped, nil
}

func InitRepo(projectDir string) error {
	_, _, err := shell.RunCommand(projectDir, gitcommand(), "init")
	return err
}
