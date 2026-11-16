package commands

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type ConfigCheckOptions struct {
	Profile string `name:"profile" description:"Validate only this named profile"`
}

func ConfigCheck(options *ConfigCheckOptions, arguments []string) error {
	profile, err := manifestProfile(options.Profile, arguments)
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(".", "")
	if err != nil {
		return err
	}
	if _, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"}); err != nil {
		return err
	}
	profiles := make([]string, 0, len(loaded.Config.Profiles))
	if profile != "" {
		profiles = append(profiles, profile)
	} else {
		for name := range loaded.Config.Profiles {
			profiles = append(profiles, name)
		}
		sort.Strings(profiles)
	}
	for _, name := range profiles {
		profileConfig, loadErr := manifest.Load(loaded.Config.Root, name)
		if loadErr != nil {
			return loadErr
		}
		if _, planErr := pipeline.PlanBuild(profileConfig.Config, pipeline.Request{Verb: "build"}); planErr != nil {
			return fmt.Errorf("profile %q: %w", name, planErr)
		}
	}
	if profile == "" {
		pterm.Success.Printfln("%s is valid (%d profiles checked)", manifest.Filename, len(profiles))
	} else {
		pterm.Success.Printfln("%s profile %q is valid", manifest.Filename, profile)
	}
	return nil
}
