package commands

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type RuntimeOptions struct {
	Directory string `name:"d" description:"Directory to generate runtime file in" default:"."`
}

func GenerateRuntime(options *RuntimeOptions) error {
	DisableFooter = true
	_, thisFile, _, _ := runtime.Caller(0)
	localDir := filepath.Dir(thisFile)
	bundledAssetsDir := filepath.Join(localDir, "..", "assetserver", "bundledassets")
	runtimeJS := filepath.Join(bundledAssetsDir, "runtime.js")
	err := CopyFile(runtimeJS, filepath.Join(options.Directory, "runtime.js"))
	if err != nil {
		return err
	}
	runtimeDebugJS := filepath.Join(bundledAssetsDir, "runtime.debug.js")
	err = CopyFile(runtimeDebugJS, filepath.Join(options.Directory, "runtime-debug.js"))
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(source string, target string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
