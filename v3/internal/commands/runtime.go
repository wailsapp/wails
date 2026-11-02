package commands

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
)

type RuntimeOptions struct {
	Directory string `name:"d" description:"Directory to generate runtime file in" default:"."`
}

func GenerateRuntime(options *RuntimeOptions) error {
	DisableFooter = true
	runtimeJS, runtimeDebugJS := bundledassets.RuntimeFiles()
	if err := os.WriteFile(filepath.Join(options.Directory, "runtime.js"), runtimeJS, 0666); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(options.Directory, "runtime-debug.js"), runtimeDebugJS, 0666)
}
