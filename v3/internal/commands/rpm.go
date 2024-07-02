package commands

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/s"
	"os"
	"path/filepath"
)

type GenerateRPMOptions struct {
	Binary      string `description:"The binary to package including path"`
	Icon        string `description:"Path to the icon"`
	DesktopFile string `description:"Path to the desktop file"`
	OutputDir   string `description:"Path to the output directory" default:"."`
	BuildDir    string `description:"Path to the build directory"`
}

func GenerateRPM(options *GenerateRPMOptions) error {

	defer func() {
		pterm.DefaultSpinner.Stop()
	}()

	if options.Binary == "" {
		return fmt.Errorf("binary not provided")
	}
	if options.Icon == "" {
		return fmt.Errorf("icon path not provided")
	}
	if options.DesktopFile == "" {
		return fmt.Errorf("desktop file path not provided")
	}
	if options.BuildDir == "" {
		// Create temp directory
		var err error
		options.BuildDir, err = os.MkdirTemp("", "wails-rpm-*")
		if err != nil {
			return err
		}
	}
	var err error
	options.OutputDir, err = filepath.Abs(options.OutputDir)
	if err != nil {
		return err
	}

	pterm.Println(pterm.LightYellow("RPM Generator v1.0.0"))

	return generateRPM(options)
}

func generateRPM(options *GenerateRPMOptions) error {
	numberOfSteps := 5
	p, _ := pterm.DefaultProgressbar.WithTotal(numberOfSteps).WithTitle("Generating RPM").Start()

	// Get the last path of the binary and normalise the name
	name := normaliseName(filepath.Base(options.Binary))

	log(p, "Preparing RPMBUILD Directory: "+options.BuildDir)

	sources := filepath.Join(options.BuildDir, "SOURCES")
	s.MKDIR(sources)
	s.COPY(options.Binary, sources)
	s.COPY(options.DesktopFile, sources)
	s.COPY(filepath.Join(options.BuildDir, name+".png"), sources)
	s.CHMOD(filepath.Join(sources, filepath.Base(options.Binary)), 0755)

	// Build RPM
	if !s.EXISTS("/usr/bin/rpmbuild") {
		return errors.New("You need to install \"rpm-build\" tool to build RPM.")
	}
	log(p, fmt.Sprintf("Running rpmbuild -bb --define \"_rpmdir %s\" --define \"_sourcedir %s\" %s.spec", options.BuildDir, sources, name))
	_, err := s.EXEC(fmt.Sprintf("rpmbuild -bb --define \"_rpmdir %s\" --define \"_sourcedir %s\" %s.spec", options.BuildDir, sources, name))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	log(p, "RPM created.")
	return nil
}
