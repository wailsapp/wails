package pipeline

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type buildOutcome struct {
	target      Target
	formats     []string
	sign        bool
	notarize    bool
	destination string
}

func resolveBuildOutcomes(config manifest.Config, request Request) ([]buildOutcome, error) {
	if config.Selected.Name != "" {
		if len(request.Targets) != 0 || request.TargetOS != "" || request.TargetArch != "" || len(request.Formats) != 0 {
			return nil, fmt.Errorf("profile %q is a complete build request and cannot be combined with --targets or --formats", config.Selected.Name)
		}
		return resolveProfileOutcomes(config.Selected)
	}
	return resolveAnonymousOutcomes(config, request)
}

func resolveProfileOutcomes(profile manifest.Profile) ([]buildOutcome, error) {
	if len(profile.Targets) == 0 {
		return nil, fmt.Errorf("profile %q requires at least one target", profile.Name)
	}
	result := make([]buildOutcome, 0, len(profile.Targets))
	seen := make(map[string]bool, len(profile.Targets))
	for _, selected := range profile.Targets {
		target, err := parseTargetName(selected.Target)
		if err != nil {
			return nil, fmt.Errorf("profile %q: %w", profile.Name, err)
		}
		name := target.OS + "/" + target.Arch
		if seen[name] {
			return nil, fmt.Errorf("profile %q contains duplicate target %s", profile.Name, name)
		}
		seen[name] = true
		capability, ok := lookupTarget(target.OS, target.Arch)
		if !ok {
			return nil, unsupportedTargetError(target.OS, target.Arch)
		}
		formats, err := resolveFormatsForTarget(capability, selected.Formats, false)
		if err != nil {
			return nil, fmt.Errorf("profile %q target %s: %w", profile.Name, name, err)
		}
		if capability.Runnable == runnableNone && len(formats) == 0 {
			return nil, fmt.Errorf("profile %q target %s must select the aab production format", profile.Name, name)
		}
		if target.OS != "ios" && selected.Destination != "" {
			return nil, fmt.Errorf("profile %q target %s: destination is only valid for ios/arm64", profile.Name, name)
		}
		if target.OS == "ios" {
			if selected.Destination != "simulator" && selected.Destination != "device" {
				return nil, fmt.Errorf("profile %q target %s requires destination = %q or %q", profile.Name, name, "simulator", "device")
			}
			if contains(formats, "ipa") && selected.Destination != "device" {
				return nil, fmt.Errorf("profile %q target %s IPA requires destination = %q", profile.Name, name, "device")
			}
		}
		if selected.Notarize && target.OS != "darwin" {
			return nil, fmt.Errorf("profile %q target %s cannot be notarized", profile.Name, name)
		}
		if selected.Notarize && !selected.Sign {
			return nil, fmt.Errorf("profile %q target %s must be signed before notarization", profile.Name, name)
		}
		result = append(result, buildOutcome{target: target, formats: formats, sign: selected.Sign, notarize: selected.Notarize, destination: selected.Destination})
	}
	return result, nil
}

func resolveAnonymousOutcomes(config manifest.Config, request Request) ([]buildOutcome, error) {
	targets := append([]Target(nil), request.Targets...)
	if len(targets) == 0 {
		targets = []Target{{OS: request.TargetOS, Arch: request.TargetArch}}
	}
	seen := make(map[string]bool, len(targets))
	for index := range targets {
		if targets[index].OS == "" {
			targets[index].OS = runtime.GOOS
		}
		if targets[index].Arch == "" {
			targets[index].Arch = runtime.GOARCH
		}
		name := targets[index].OS + "/" + targets[index].Arch
		if seen[name] {
			return nil, fmt.Errorf("duplicate target %s", name)
		}
		seen[name] = true
		if _, ok := lookupTarget(targets[index].OS, targets[index].Arch); !ok {
			return nil, unsupportedTargetError(targets[index].OS, targets[index].Arch)
		}
	}
	sort.Slice(targets, func(left, right int) bool {
		return targets[left].OS+"/"+targets[left].Arch < targets[right].OS+"/"+targets[right].Arch
	})
	switch request.Verb {
	case "", "build", "package", "sign":
	default:
		return nil, fmt.Errorf("unsupported pipeline verb %q", request.Verb)
	}

	result := make([]buildOutcome, len(targets))
	for index, target := range targets {
		result[index].target = target
	}
	if len(request.Formats) != 0 {
		formats, err := uniqueFormats(request.Formats)
		if err != nil {
			return nil, err
		}
		for _, format := range formats {
			if _, ok := lookupFormat(format); !ok {
				return nil, fmt.Errorf("unknown package format %q", format)
			}
			if format == "apk" && !request.Development {
				return nil, fmt.Errorf("production APK is no longer supported; select aab (APK remains available only to the development deployment flow)")
			}
		}
		matched := make(map[string]bool, len(formats))
		for index, target := range targets {
			capability, _ := lookupTarget(target.OS, target.Arch)
			for _, format := range formats {
				if capability.SupportsFormat(format, request.Development) {
					result[index].formats = append(result[index].formats, format)
					matched[format] = true
				}
			}
		}
		for _, format := range formats {
			if !matched[format] {
				return nil, fmt.Errorf("format %q is not supported for any selected target", format)
			}
		}
		for index, target := range targets {
			if len(result[index].formats) == 0 {
				return nil, fmt.Errorf("target %s/%s receives no compatible format from --formats", target.OS, target.Arch)
			}
		}
	} else {
		for index, target := range targets {
			capability, _ := lookupTarget(target.OS, target.Arch)
			switch request.Verb {
			case "package", "sign":
				configured := packagePlatform(config.Package, target.OS).Formats
				formats, err := resolveFormatsForTarget(capability, configured, request.Development)
				if err != nil {
					return nil, fmt.Errorf("target %s/%s: %w", target.OS, target.Arch, err)
				}
				result[index].formats = formats
			case "", "build":
				if capability.Runnable == runnableNone && !request.Development {
					result[index].formats = []string{"aab"}
				}
			}
		}
	}
	for index := range result {
		result[index].sign = request.Verb == "sign"
	}
	return result, nil
}

func resolveFormatsForTarget(target targetCapability, requested []string, development bool) ([]string, error) {
	formats, err := uniqueFormats(requested)
	if err != nil {
		return nil, err
	}
	for _, format := range formats {
		capability, ok := lookupFormat(format)
		if !ok {
			return nil, fmt.Errorf("unknown package format %q", format)
		}
		if format == "apk" && !development {
			return nil, fmt.Errorf("production APK is no longer supported; select aab (APK remains available only to the development deployment flow)")
		}
		if !target.SupportsFormat(capability.Name, development) {
			mode := "production"
			if development {
				mode = "development"
			}
			return nil, fmt.Errorf("package format %q is not supported for %s/%s in %s", format, target.Target.OS, target.Target.Arch, mode)
		}
	}
	return formats, nil
}

func uniqueFormats(values []string) ([]string, error) {
	result := append([]string(nil), values...)
	sort.Strings(result)
	for index := range result {
		if strings.TrimSpace(result[index]) == "" {
			return nil, fmt.Errorf("package format cannot be empty")
		}
		if index != 0 && result[index] == result[index-1] {
			return nil, fmt.Errorf("duplicate package format %s", result[index])
		}
	}
	return result, nil
}

func parseTargetName(value string) (Target, error) {
	platform, arch, found := strings.Cut(value, "/")
	if !found || platform == "" || arch == "" || strings.Contains(arch, "/") {
		return Target{}, fmt.Errorf("invalid target %q", value)
	}
	return Target{OS: platform, Arch: arch}, nil
}
