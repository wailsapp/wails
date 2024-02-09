package commands

import (
	"github.com/wailsapp/wails/v3/internal/commands/build_assets/runtime"
	"os"
)

type GenerateRuntimeOptions struct {
}

func GenerateRuntime(options *GenerateRuntimeOptions) error {
	err := os.WriteFile("runtime.js", runtime.RuntimeJS, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile("runtime.debug.js", runtime.RuntimeDebugJS, 0644)
	if err != nil {
		return err
	}
	return nil
}
