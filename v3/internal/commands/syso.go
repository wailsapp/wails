package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

type SysoOptions struct {
	Example  bool   `description:"Generate example manifest & info files"`
	Manifest string `description:"The manifest file"`
	Info     string `description:"The info.json file"`
	Icon     string `description:"The icon file"`
	Out      string `description:"The output filename for the syso file"`
	Arch     string `description:"The target architecture"`
}

func (i *SysoOptions) Default() *SysoOptions {
	return &SysoOptions{
		Arch: runtime.GOARCH,
	}
}

func GenerateSyso(options *SysoOptions) error {

	// Generate example files?
	if options.Example {
		return generateExampleSyso()
	}

	if options.Manifest == "" {
		return fmt.Errorf("manifest is required")
	}
	if options.Icon == "" {
		return fmt.Errorf("icon is required")
	}

	rs := winres.ResourceSet{}

	// Process Icon
	iconFile, err := os.Open(options.Icon)
	if err != nil {
		return err
	}
	defer iconFile.Close()
	ico, err := winres.LoadICO(iconFile)
	if err != nil {
		return fmt.Errorf("couldn't load icon '%s': %v", options.Icon, err)
	}
	err = rs.SetIcon(winres.RT_ICON, ico)
	if err != nil {
		return err
	}

	// Process Manifest
	manifestData, err := os.ReadFile(options.Manifest)
	if err != nil {
		return err
	}

	xmlData, err := winres.AppManifestFromXML(manifestData)
	if err != nil {
		return err
	}
	rs.SetManifest(xmlData)

	if options.Info != "" {
		infoData, err := os.ReadFile(options.Info)
		if err != nil {
			return err
		}
		if len(infoData) != 0 {
			var v version.Info
			if err := v.UnmarshalJSON(infoData); err != nil {
				return err
			}
			rs.SetVersionInfo(v)
		}
	}

	targetFile := options.Out
	if targetFile == "" {
		targetFile = "rsrc_windows_" + options.Arch + ".syso"
	}
	fout, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	archs := map[string]winres.Arch{
		"amd64": winres.ArchAMD64,
		"arm64": winres.ArchARM64,
		"386":   winres.ArchI386,
	}
	targetArch, supported := archs[options.Arch]
	if !supported {
		return fmt.Errorf("arch '%s' not supported", options.Arch)
	}

	err = rs.WriteObject(fout, targetArch)
	if err != nil {
		return err
	}
	return nil
}

func generateExampleSyso() error {
	// Generate example info.json
	err := os.WriteFile("info.json", Info, 0644)
	if err != nil {
		return err
	}
	// Generate example manifest
	err = os.WriteFile("wails.exe.manifest", Manifest, 0644)
	if err != nil {
		return err
	}
	return nil
}
