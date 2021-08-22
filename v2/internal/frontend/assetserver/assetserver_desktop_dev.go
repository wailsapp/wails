//go:build dev

package assetserver

import (
	"os"
	"path/filepath"
)

func (a *DesktopAssetServer) ReadFile(filename string) ([]byte, error) {
	a.LogInfo("Loading file from disk: %s", filename)
	return os.ReadFile(filepath.Join(a.assetdir, filename))
}
