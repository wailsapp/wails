package git

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	name := template.JSEscapeString(strings.TrimSpace(stdout))

	// Check if username is JSON compliant
	var js json.RawMessage
	jsonVal := fmt.Sprintf(`{"name": "%s"}`, name)
	if json.Unmarshal([]byte(jsonVal), &js) != nil {
		return "", fmt.Errorf(errMsg, err)
	}

	return name, err
}

func InitRepo(projectDir string) error {
	_, _, err := shell.RunCommand(projectDir, gitcommand(), "init")
	return err
}
