package webview2runtime

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed MicrosoftEdgeWebview2Setup.exe
var setupexe []byte

// WriteInstallerToFile writes the installer file to the given file.
func WriteInstallerToFile(targetFile string) error {
	return os.WriteFile(targetFile, setupexe, 0o755)
}

// WriteInstaller writes the installer exe file to the given directory and returns the path to it.
func WriteInstaller(targetPath string) (string, error) {
	installer := filepath.Join(targetPath, `MicrosoftEdgeWebview2Setup.exe`)
	return installer, WriteInstallerToFile(installer)
}
