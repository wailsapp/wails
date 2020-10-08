package build

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/html"
	"github.com/wailsapp/wails/v2/internal/logger"
)

// DesktopBuilder builds applications for the desktop
type DesktopBuilder struct {
	*BaseBuilder
}

func newDesktopBuilder() *DesktopBuilder {
	return &DesktopBuilder{
		BaseBuilder: NewBaseBuilder(),
	}
}

// BuildAssets builds the assets for the desktop application
func (d *DesktopBuilder) BuildAssets(options *Options) error {
	var err error

	// Get a list of assets from the HTML
	assets, err := d.BaseBuilder.ExtractAssets()
	if err != nil {
		return err
	}

	// Build base assets (HTML/JS/CSS/etc)
	err = d.BuildBaseAssets(assets, options.Logger)
	if err != nil {
		return err
	}

	// Build static assets
	err = d.buildStaticAssets(d.projectData)
	if err != nil {
		return err
	}

	return nil
}

// BuildBaseAssets builds the assets for the desktop application
func (d *DesktopBuilder) BuildBaseAssets(assets *html.AssetBundle, outputLogger *logger.Logger) error {
	var err error

	outputLogger.Write("    - Embedding Assets...")

	// Get target asset directory
	assetDir := fs.RelativePath("../../../internal/ffenestri")

	// Dump assets as C
	assetsFile, err := assets.WriteToCFile(assetDir)
	if err != nil {
		return err
	}
	d.addFileToDelete(assetsFile)

	// Process Icon
	err = d.processIcon(assetDir)
	if err != nil {
		return err
	}

	outputLogger.Writeln("done.")

	return nil
}

// processIcon will copy a default icon if one doesn't exist, then, if
// needed, will compile the icon
func (d *DesktopBuilder) processIcon(assetDir string) error {

	// Copy default icon if one doesn't exist
	iconFile := filepath.Join(d.projectData.Path, "icon.png")
	if !fs.FileExists(iconFile) {
		err := fs.CopyFile(defaultIconPath(), iconFile)
		if err != nil {
			return err
		}
	}

	// Compile Icon
	return d.compileIcon(assetDir, iconFile)
}

// BuildRuntime builds the Wails javascript runtime and then converts it into a C file
func (d *DesktopBuilder) BuildRuntime(options *Options) error {

	outputLogger := options.Logger

	sourceDir := fs.RelativePath("../../../internal/runtime/js")

	if err := d.NpmInstall(sourceDir); err != nil {
		return err
	}

	outputLogger.Write("    - Embedding Runtime...")
	envvars := []string{"WAILSPLATFORM=" + options.Platform}
	if err := d.NpmRunWithEnvironment(sourceDir, "build:desktop", false, envvars); err != nil {
		return err
	}

	wailsJS := fs.RelativePath("../../../internal/runtime/assets/desktop.js")
	runtimeData, err := ioutil.ReadFile(wailsJS)
	if err != nil {
		return err
	}
	outputLogger.Writeln("done.")

	// Convert to C structure
	runtimeC := `
// runtime.c (c) 2019-Present Lea Anthony.
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file was auto-generated. DO NOT MODIFY.
const unsigned char runtime[]={`
	for _, b := range runtimeData {
		runtimeC += fmt.Sprintf("0x%x, ", b)
	}
	runtimeC += "0x00};"

	// Save file
	outputFile := fs.RelativePath("../../../internal/ffenestri/runtime.c")

	if err := ioutil.WriteFile(outputFile, []byte(runtimeC), 0600); err != nil {
		return err
	}

	return nil
}
