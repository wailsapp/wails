package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/version"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type ConfigOptions struct {
	Profile string `name:"profile" description:"Manifest profile to resolve"`
	JSON    bool   `name:"json" description:"Print resolved configuration as JSON"`
}

func ConfigCheck(options *ConfigOptions) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(root, options.Profile)
	if err != nil {
		return err
	}
	fmt.Printf("%s is valid (%s)\n", manifest.Filename, loaded.Config.Project.Identifier)
	return nil
}

func ConfigShow(options *ConfigOptions) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(root, options.Profile)
	if err != nil {
		return err
	}
	if options.JSON {
		data, err := json.MarshalIndent(loaded.Config, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	data, err := manifest.EncodeConfig(loaded.Config)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

type EjectOptions struct {
	Backup bool `name:"backup" description:"Keep a timestamped backup of wails.toml"`
}

func Eject(options *EjectOptions, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: wails3 eject [profile] [--backup]")
	}
	profile := ""
	if len(args) == 1 {
		profile = strings.TrimSpace(args[0])
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := manifest.Eject(root, profile, version.String(), options.Backup); err != nil {
		if errors.Is(err, manifest.ErrEjectionSuggestionsUnavailable) {
			fmt.Println(err)
			return nil
		}
		return err
	}
	if profile == "" {
		fmt.Println("Ejected the complete default build configuration into wails.toml")
	} else {
		fmt.Printf("Ejected profile %q into wails.toml\n", profile)
	}
	return nil
}
