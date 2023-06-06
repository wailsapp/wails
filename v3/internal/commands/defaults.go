package commands

import (
	_ "embed"
	"os"
)

//go:embed defaults/info.json
var Info []byte

//go:embed defaults/wails.exe.manifest
var Manifest []byte

//go:embed defaults/appicon.png
var AppIcon []byte

//go:embed defaults/icon.ico
var IconIco []byte

//go:embed defaults/Info.plist
var InfoPlist []byte

//go:embed defaults/Info.dev.plist
var InfoDevPlist []byte

//go:embed defaults/icons.icns
var IconsIcns []byte

var AllAssets = map[string][]byte{
	"info.json":          Info,
	"wails.exe.manifest": Manifest,
	"appicon.png":        AppIcon,
	"icon.ico":           IconIco,
	"Info.plist":         InfoPlist,
	"Info.dev.plist":     InfoDevPlist,
	"icons.icns":         IconsIcns,
}

type DefaultsOptions struct {
	Dir string `description:"The directory to generate the files into"`
}

func Defaults(options *DefaultsOptions) error {
	dir := options.Dir
	if dir == "" {
		dir = "."
	}
	for filename, data := range AllAssets {
		// If file exists, skip it
		if _, err := os.Stat(dir + "/" + filename); err == nil {
			println("Skipping " + filename)
			continue
		}
		err := os.WriteFile(dir+"/"+filename, data, 0644)
		if err != nil {
			return err
		}
		println("Generated " + filename)
	}
	return nil
}
