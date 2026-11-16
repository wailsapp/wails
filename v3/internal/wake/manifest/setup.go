package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// SetupState owns only the primary application backend and one named production
// profile. Package templates, hooks, companions and other profiles stay intact.
type SetupState struct {
	Kind             string
	PureGo           bool
	FrontendDisabled bool
	Profile          Profile
	Signing          map[string]SigningPlatform
}

// PreviewSetup validates the exact document that WriteSetup would publish.
// A non-nil initial state creates a new manifest; existing files are never replaced.
func PreviewSetup(start string, initial *InitialState, state SetupState) (Config, error) {
	_, config, _, _, err := prepareSetup(start, initial, state)
	return config, err
}

func WriteSetup(start string, initial *InitialState, state SetupState) error {
	data, _, path, mode, err := prepareSetup(start, initial, state)
	if err != nil {
		return err
	}
	if initial != nil {
		return exclusiveWrite(path, data, mode)
	}
	return atomicWrite(path, data, mode)
}

func prepareSetup(start string, initial *InitialState, state SetupState) ([]byte, Config, string, os.FileMode, error) {
	root, err := filepath.Abs(start)
	if err != nil {
		return nil, Config{}, "", 0, err
	}
	path, mode := filepath.Join(root, Filename), os.FileMode(0644)
	var raw []byte
	if initial != nil {
		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
			if err == nil {
				err = fmt.Errorf("%s already exists; setup will not replace it", path)
			}
			return nil, Config{}, "", 0, err
		}
		raw = EncodeInitial(*initial)
	} else {
		loaded, err := Load(start, "")
		if err != nil {
			return nil, Config{}, "", 0, err
		}
		root, path, raw = loaded.Config.Root, loaded.Path, loaded.Raw
		info, err := os.Stat(path)
		if err != nil {
			return nil, Config{}, "", 0, err
		}
		mode = info.Mode().Perm()
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, Config{}, "", 0, err
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, Config{}, "", 0, errors.New("setup requires a YAML mapping")
	}
	if err := applySetup(doc.Content[0], state); err != nil {
		return nil, Config{}, "", 0, err
	}
	data, err := marshalYAML(&doc)
	if err != nil {
		return nil, Config{}, "", 0, err
	}
	loaded, err := decodeYAML(root, path, data, state.Profile.Name)
	if err != nil {
		return nil, Config{}, "", 0, err
	}
	return data, loaded.Config, path, mode, nil
}

func applySetup(root *yaml.Node, state SetupState) error {
	if state.Kind != "gui" && state.Kind != "cli" && state.Kind != "server" {
		return errors.New("application type must be gui, cli or server")
	}
	if state.Profile.Name != "setup" {
		return errors.New("setup may edit only the setup profile")
	}
	seen := make(map[string]bool, len(state.Profile.Targets))
	for _, target := range state.Profile.Targets {
		if seen[target.Target] {
			return fmt.Errorf("duplicate setup target %q", target.Target)
		}
		seen[target.Target] = true
	}
	yamlUpdate(yamlEnsureMap(root, "project"), "kind", yamlString(state.Kind))
	yamlUpdate(yamlEnsureMap(root, "frontend"), "disabled", yamlBool(state.FrontendDisabled))
	build := yamlEnsureMap(root, "build")
	mode := ""
	if state.PureGo {
		mode = "0"
	}
	resetZero := false
	for _, target := range state.Profile.Targets {
		if state.Kind == "gui" && (strings.HasPrefix(target.Target, "linux/") || strings.HasPrefix(target.Target, "darwin/")) {
			resetZero = true
		}
	}
	setSetupBackend(build, state.PureGo, mode, resetZero)
	targets := yamlEnsureMap(root, "targets")
	for _, selected := range state.Profile.Targets {
		target := yamlEnsureMap(targets, selected.Target)
		// Target environments replace build.environment. Preserve inherited keys
		// when introducing the first target-specific backend override.
		if mappingValue(target, "environment") == nil {
			var inherited map[string]string
			if err := mappingValue(build, "environment").Decode(&inherited); err != nil {
				return err
			}
			yamlSetOptionalStringMap(target, "environment", inherited)
		}
		setSetupBackend(target, state.PureGo, mode, resetZero)
	}
	profiles := yamlEnsureMap(root, "profiles")
	profile := yamlEnsureMap(profiles, state.Profile.Name)
	selected := yamlMap()
	for _, target := range state.Profile.Targets {
		node := yamlMap()
		yamlSet(node, "formats", yamlStrings(target.Formats))
		yamlSet(node, "sign", yamlBool(target.Sign))
		if target.Notarize {
			yamlSet(node, "notarize", yamlBool(true))
		}
		yamlSetOptionalString(node, "destination", target.Destination)
		yamlSet(selected, target.Target, node)
	}
	yamlUpdate(profile, "targets", selected)
	platforms := make([]string, 0, len(state.Signing))
	for platform := range state.Signing {
		platforms = append(platforms, platform)
	}
	sort.Strings(platforms)
	for _, platform := range platforms {
		signing := state.Signing[platform]
		encoded := encodeYAMLPlatform(platform, Platform{}, signing)
		if encoded == nil {
			continue
		}
		node := mappingValue(encoded, "signing")
		if node == nil {
			continue
		}
		existing := yamlEnsureMap(yamlEnsureMap(root, platform), "signing")
		for i := 0; i+1 < len(node.Content); i += 2 {
			yamlUpdate(existing, node.Content[i].Value, node.Content[i+1])
		}
	}
	return nil
}

func setSetupBackend(node *yaml.Node, pureGo bool, mode string, resetZero bool) {
	var tags []string
	if previous := mappingValue(node, "tags"); previous != nil {
		_ = previous.Decode(&tags)
	}
	resetZero = resetZero || slices.Contains(tags, "purego")
	tags = slices.DeleteFunc(tags, func(tag string) bool { return tag == "purego" })
	if pureGo {
		tags = append(tags, "purego")
	}
	yamlUpdate(node, "tags", yamlStrings(tags))
	environment := yamlEnsureMap(node, "environment")
	if mode != "" {
		yamlUpdate(environment, "CGO_ENABLED", yamlString(mode))
	} else if previous := mappingValue(environment, "CGO_ENABLED"); previous != nil && previous.Value == "0" && resetZero {
		yamlRemove(environment, "CGO_ENABLED")
	}
}
