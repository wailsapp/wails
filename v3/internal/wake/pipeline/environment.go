package pipeline

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

const CrossContainerImage = "ghcr.io/wailsapp/wails-cross:latest"

func preferredCrossContainerImages() []string {
	return []string{CrossContainerImage, "wails-cross"}
}

// HostCapabilities is a value snapshot of the execution facts that may affect
// Plan feasibility. Its collections are private and copied at construction so
// callers cannot mutate a capability snapshot after planning starts.
type HostCapabilities struct {
	os                   string
	arch                 string
	tools                []string
	credentials          []string
	dockerImages         []string
	containerRuntime     string
	containerImages      []string
	containerImagesKnown bool
	appleSDKs            []string
	androidSDK           bool
	androidNDK           bool
}

type HostFacts struct {
	DockerImages     []string
	ContainerRuntime string
	ContainerImages  []string
	AppleSDKs        []string
	AndroidSDK       bool
	AndroidNDK       bool
}

func NewHostCapabilities(hostOS, hostArch string, tools, credentials []string) HostCapabilities {
	return NewHostCapabilitiesWithFacts(hostOS, hostArch, tools, credentials, HostFacts{})
}

func NewHostCapabilitiesWithFacts(hostOS, hostArch string, tools, credentials []string, facts HostFacts) HostCapabilities {
	sortedTools := uniqueSorted(tools)
	containerRuntime := facts.ContainerRuntime
	containerImages := facts.ContainerImages
	if containerRuntime == "" && len(facts.DockerImages) != 0 {
		containerRuntime = "docker"
		containerImages = facts.DockerImages
	}
	if containerRuntime == "" {
		for _, candidate := range []string{"docker", "podman"} {
			if containsSorted(sortedTools, candidate) {
				containerRuntime = candidate
				break
			}
		}
	}
	return HostCapabilities{
		os: hostOS, arch: hostArch, tools: sortedTools, credentials: uniqueSorted(credentials),
		dockerImages: uniqueSorted(facts.DockerImages), appleSDKs: uniqueSorted(facts.AppleSDKs),
		containerRuntime: containerRuntime, containerImages: uniqueSorted(containerImages),
		containerImagesKnown: true,
		androidSDK:           facts.AndroidSDK, androidNDK: facts.AndroidNDK,
	}
}

func CurrentHostCapabilities(credentials ...string) HostCapabilities {
	return currentHostCapabilitiesWithOperations(credentials, hostProbeOperations{
		hostOS: runtime.GOOS, hostArch: runtime.GOARCH, lookPath: exec.LookPath, lookupEnv: os.LookupEnv, getenv: os.Getenv,
		stat: os.Stat, glob: filepath.Glob, run: func(name string, arguments ...string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return exec.CommandContext(ctx, name, arguments...).Run()
		},
	})
}

type hostProbeOperations struct {
	hostOS, hostArch string
	lookPath         func(string) (string, error)
	lookupEnv        func(string) (string, bool)
	getenv           func(string) string
	stat             func(string) (fs.FileInfo, error)
	glob             func(string) ([]string, error)
	run              func(string, ...string) error
}

func currentHostCapabilitiesWithOperations(credentials []string, operations hostProbeOperations) HostCapabilities {
	tools, presentCredentials := probeHostToolsAndCredentials(credentials, operations)
	facts := currentHostFacts(tools, operations, true)
	return NewHostCapabilitiesWithFacts(operations.hostOS, operations.hostArch, tools, presentCredentials, facts)
}

func probeHostToolsAndCredentials(credentials []string, operations hostProbeOperations) ([]string, []string) {
	knownTools := []string{"go", "garble", "zig", "docker", "podman", "cc", "gcc", "clang", "npm", "pnpm", "yarn", "bun", "makensis", "MakeAppx.exe", "signtool.exe", "hdiutil", "xcrun", "codesign", "ditto", "zip", "java", "jarsigner", "apksigner", "dpkg-sig", "rpmsign"}
	tools := make([]string, 0, len(knownTools))
	for _, tool := range knownTools {
		if _, err := operations.lookPath(tool); err == nil {
			tools = append(tools, tool)
		}
	}
	presentCredentials := make([]string, 0, len(credentials))
	for _, credential := range credentials {
		if credential != "" {
			if _, ok := operations.lookupEnv(credential); ok {
				presentCredentials = append(presentCredentials, credential)
			}
		}
	}
	return tools, presentCredentials
}

func (h HostCapabilities) hasTool(name string) bool       { return containsSorted(h.tools, name) }
func (h HostCapabilities) hasCredential(name string) bool { return containsSorted(h.credentials, name) }
func (h HostCapabilities) hasDockerImage(name string) bool {
	return containsSorted(h.containerImages, name) || containsSorted(h.dockerImages, name)
}
func (h HostCapabilities) hasAppleSDK(name string) bool { return containsSorted(h.appleSDKs, name) }

func currentHostFacts(tools []string, operations hostProbeOperations, probeContainerImages bool) HostFacts {
	sortedTools := uniqueSorted(tools)
	hasTool := func(name string) bool { return containsSorted(sortedTools, name) }
	facts := HostFacts{}
	sdk := firstNonemptyEnvironment(operations.getenv, "ANDROID_HOME", "ANDROID_SDK_ROOT")
	facts.AndroidSDK = directoryExists(sdk, operations.stat)
	ndk := operations.getenv("ANDROID_NDK_HOME")
	if !directoryExists(ndk, operations.stat) && facts.AndroidSDK {
		matches, _ := operations.glob(filepath.Join(sdk, "ndk", "*"))
		for _, match := range matches {
			if directoryExists(match, operations.stat) {
				ndk = match
				break
			}
		}
	}
	facts.AndroidNDK = directoryExists(ndk, operations.stat)
	for _, containerRuntime := range []string{"docker", "podman"} {
		if !hasTool(containerRuntime) {
			continue
		}
		if facts.ContainerRuntime == "" {
			facts.ContainerRuntime = containerRuntime
		}
		if probeContainerImages {
			for _, image := range preferredCrossContainerImages() {
				if operations.run(containerRuntime, "image", "inspect", image) == nil {
					facts.ContainerRuntime = containerRuntime
					facts.ContainerImages = []string{image}
					if containerRuntime == "docker" {
						facts.DockerImages = []string{image}
					}
					break
				}
			}
		}
		if !probeContainerImages || len(facts.ContainerImages) != 0 {
			break
		}
	}
	if hasTool("xcrun") {
		for _, sdkName := range []string{"macosx", "iphoneos", "iphonesimulator"} {
			if operations.run("xcrun", "--sdk", sdkName, "--show-sdk-path") == nil {
				facts.AppleSDKs = append(facts.AppleSDKs, sdkName)
			}
		}
	}
	return facts
}

// PlanBuildForCurrentHost keeps native planning free of container subprocesses.
// It probes image availability only after the structural host plan selects the
// container-backed toolchain, then re-plans with the complete immutable facts.
func PlanBuildForCurrentHost(config manifest.Config, request Request, credentials ...string) (Plan, error) {
	operations := hostProbeOperations{
		hostOS: runtime.GOOS, hostArch: runtime.GOARCH, lookPath: exec.LookPath, lookupEnv: os.LookupEnv, getenv: os.Getenv,
		stat: os.Stat, glob: filepath.Glob, run: func(name string, arguments ...string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return exec.CommandContext(ctx, name, arguments...).Run()
		},
	}
	return planBuildForCurrentHostWithOperations(config, request, credentials, operations)
}

func planBuildForCurrentHostWithOperations(config manifest.Config, request Request, credentials []string, operations hostProbeOperations) (Plan, error) {
	tools, presentCredentials := probeHostToolsAndCredentials(credentials, operations)
	facts := currentHostFacts(tools, operations, false)
	host := NewHostCapabilitiesWithFacts(operations.hostOS, operations.hostArch, tools, presentCredentials, facts)
	host.containerImagesKnown = false
	plan, err := PlanBuildForHost(config, request, host)
	if err != nil || !planUsesContainerToolchain(plan) {
		return plan, err
	}
	facts = currentHostFacts(tools, operations, true)
	host = NewHostCapabilitiesWithFacts(operations.hostOS, operations.hostArch, tools, presentCredentials, facts)
	return PlanBuildForHost(config, request, host)
}

func planUsesContainerToolchain(plan Plan) bool {
	for _, node := range plan.Nodes {
		if node.Kind == CompileApplication && node.Spec.(CompileSpec).Toolchain == "docker" {
			return true
		}
	}
	return false
}

func firstNonemptyEnvironment(getenv func(string) string, names ...string) string {
	for _, name := range names {
		if value := getenv(name); value != "" {
			return value
		}
	}
	return ""
}

func directoryExists(path string, stat func(string) (fs.FileInfo, error)) bool {
	info, err := stat(path)
	return err == nil && info.IsDir()
}

func PlanBuildForHost(config manifest.Config, request Request, host HostCapabilities) (Plan, error) {
	if config.Selected.Name == "" && len(request.Targets) == 0 {
		if request.TargetOS == "" {
			request.TargetOS = host.os
		}
		if request.TargetArch == "" {
			request.TargetArch = host.arch
		}
	}
	plan, err := PlanBuild(config, request)
	if err != nil {
		return Plan{}, err
	}
	if hostBit(host.os) == 0 || host.arch == "" {
		return Plan{}, fmt.Errorf("unsupported build host %s/%s", host.os, host.arch)
	}
	keys := make([]string, 0, len(plan.Nodes))
	for key := range plan.Nodes {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	for _, rawKey := range keys {
		key := NodeKey(rawKey)
		node := plan.Nodes[key]
		switch node.Kind {
		case InstallFrontendDependencies:
			spec := node.Spec.(InstallSpec)
			if err := requireCommandTool(host, "frontend dependency installation", firstCommand(spec.Arguments, spec.Manager)); err != nil {
				return Plan{}, err
			}
		case GenerateBindings:
			if err := requireHostTools(host, "bindings generation", "go"); err != nil {
				return Plan{}, err
			}
		case BuildFrontend:
			spec := node.Spec.(FrontendSpec)
			if err := requireCommandTool(host, "frontend build", firstCommand(spec.Arguments, spec.Manager)); err != nil {
				return Plan{}, err
			}
		case CompileApplication:
			spec := node.Spec.(CompileSpec)
			resolved, resolveErr := resolveToolchain(spec, host)
			if resolveErr != nil {
				return Plan{}, resolveErr
			}
			spec.Toolchain = resolved
			if resolved == "docker" {
				spec.ContainerRuntime = host.containerRuntime
				spec.ContainerImage = host.crossContainerImage()
			}
			if err := validateCompileEnvironment(spec, host); err != nil {
				return Plan{}, err
			}
			node.Spec = spec
			plan.Nodes[key] = node
		case AssembleApplication:
			spec := node.Spec.(PackageSpec)
			if spec.TargetOS == "ios" {
				if host.os != "darwin" {
					return Plan{}, fmt.Errorf("iOS application assembly for %s/%s requires a darwin host", spec.TargetOS, spec.TargetArch)
				}
				if err := requireHostTools(host, "iOS application assembly for "+spec.TargetOS+"/"+spec.TargetArch, "xcrun", "codesign"); err != nil {
					return Plan{}, err
				}
			}
		case PackageArtifact:
			spec := node.Spec.(PackageSpec)
			capability, _ := lookupFormat(spec.Format)
			operation := spec.Format + " packaging for " + spec.TargetOS + "/" + spec.TargetArch
			if !capability.SupportsHost(host.os) {
				return Plan{}, fmt.Errorf("%s requires a %s host", operation, requiredHostName(capability.Hosts))
			}
			for index := range int(capability.ToolCount) {
				if err := requireHostTools(host, operation, capability.RequiredTool(index)); err != nil {
					return Plan{}, err
				}
			}
			if spec.TargetOS == "android" && !host.androidSDK {
				return Plan{}, fmt.Errorf("%s requires an Android SDK", operation)
			}
		case SignArtifact:
			if err := validateSigningHost(node.Spec.(SignSpec), host); err != nil {
				return Plan{}, err
			}
		}
	}
	return plan, nil
}

func firstCommand(arguments []string, fallback string) string {
	if len(arguments) != 0 {
		return arguments[0]
	}
	return fallback
}

func requireCommandTool(host HostCapabilities, operation, command string) error {
	if command == "" || strings.ContainsAny(command, `/\\`) {
		return nil
	}
	return requireHostTools(host, operation, command)
}

func validateCompileEnvironment(spec CompileSpec, host HostCapabilities) error {
	operation := "compiling " + spec.TargetOS + "/" + spec.TargetArch
	if spec.Toolchain == "docker" {
		if host.containerImagesKnown && host.crossContainerImage() == "" {
			return fmt.Errorf("%s with Docker or Podman requires container image %q", operation, CrossContainerImage)
		}
		return nil
	}
	tool := "go"
	if spec.Obfuscated {
		tool = "garble"
	}
	if err := requireHostTools(host, operation, "go", tool); err != nil {
		return err
	}
	if spec.TargetOS == "linux" && spec.Toolchain == "native" && !host.hasTool("cc") && !host.hasTool("gcc") && !host.hasTool("clang") {
		return fmt.Errorf("%s with the native toolchain requires a C compiler (cc, gcc, or clang)", operation)
	}
	if spec.TargetOS == "android" {
		if !host.androidSDK {
			return fmt.Errorf("%s requires an Android SDK", operation)
		}
		if !host.androidNDK {
			return fmt.Errorf("%s requires an Android NDK", operation)
		}
	}
	if spec.TargetOS == "ios" {
		sdk := "iphonesimulator"
		if spec.Destination == "device" {
			sdk = "iphoneos"
		}
		if !host.hasAppleSDK(sdk) {
			return fmt.Errorf("%s requires Apple SDK %q", operation, sdk)
		}
	}
	return nil
}

func (h HostCapabilities) crossContainerImage() string {
	for _, image := range preferredCrossContainerImages() {
		if h.hasDockerImage(image) {
			return image
		}
	}
	return ""
}

func resolveToolchain(spec CompileSpec, host HostCapabilities) (string, error) {
	target := spec.TargetOS + "/" + spec.TargetArch
	requested := spec.Toolchain
	if requested == "" {
		requested = "auto"
	}
	if requested == "auto" {
		if nativeToolchainSupports(spec, host) {
			return "native", nil
		}
		if spec.TargetOS == "windows" && host.hasTool("zig") {
			return "zig", nil
		}
		if dockerToolchainSupports(spec, host) && host.containerRuntime != "" {
			return "docker", nil
		}
		if spec.TargetOS == "linux" {
			return "", fmt.Errorf("target %s on %s/%s requires the wails-cross image with Docker or Podman; plain Zig does not provide the target Linux desktop sysroot", target, host.os, host.arch)
		}
		return "", fmt.Errorf("target %s on %s/%s requires zig, Docker, or Podman; no compatible toolchain is available", target, host.os, host.arch)
	}
	if requested == "native" && !nativeToolchainSupports(spec, host) {
		return "", fmt.Errorf("toolchain %q cannot build %s on %s/%s", requested, target, host.os, host.arch)
	}
	if requested == "zig" && !host.hasTool("zig") {
		return "", fmt.Errorf("toolchain %q requires tool %q", requested, "zig")
	}
	if requested == "docker" && !dockerToolchainSupports(spec, host) {
		return "", fmt.Errorf("toolchain %q cannot build %s on %s/%s", requested, target, host.os, host.arch)
	}
	if requested == "docker" && host.containerRuntime == "" {
		return "", fmt.Errorf("toolchain %q requires Docker or Podman", requested)
	}
	return requested, nil
}

func dockerToolchainSupports(spec CompileSpec, host HostCapabilities) bool {
	capability, ok := lookupTarget(spec.TargetOS, spec.TargetArch)
	return ok && host.os != "windows" && capability.SupportsToolchain("docker")
}

func nativeToolchainSupports(spec CompileSpec, host HostCapabilities) bool {
	switch spec.TargetOS {
	case "android":
		return hostBit(host.os) != 0
	case "ios":
		return host.os == "darwin"
	default:
		return host.os == spec.TargetOS && (host.arch == spec.TargetArch || spec.TargetOS == "darwin")
	}
}

func validateSigningHost(spec SignSpec, host HostCapabilities) error {
	operation := "signing " + spec.Format + " for " + spec.TargetOS
	if !spec.Config.Enabled {
		return fmt.Errorf("signing is not enabled for %s", spec.TargetOS)
	}
	switch spec.TargetOS {
	case "darwin":
		if host.os != "darwin" {
			return fmt.Errorf("%s requires a darwin host", operation)
		}
		if spec.Config.Identity == "" {
			return fmt.Errorf("darwin signing requires signing.darwin.identity")
		}
		if err := requireHostTools(host, operation, "codesign"); err != nil {
			return err
		}
		if spec.Config.Notarize {
			if spec.Config.NotarizationCredential == "" {
				return fmt.Errorf("notarization requires signing.darwin.notarization credential")
			}
			if err := requireHostTools(host, "darwin notarization", "xcrun"); err != nil {
				return err
			}
			if spec.Format == "app" {
				if err := requireHostTools(host, "darwin app notarization", "ditto"); err != nil {
					return err
				}
			}
		}
	case "ios":
		if host.os != "darwin" {
			return fmt.Errorf("%s requires a darwin host", operation)
		}
		if spec.Config.Identity == "" {
			return fmt.Errorf("iOS signing requires signing.ios.identity")
		}
		return requireHostTools(host, operation, "codesign")
	case "windows":
		if host.os != "windows" {
			return fmt.Errorf("%s requires a windows host", operation)
		}
		if spec.Config.Certificate == "" && spec.Config.Thumbprint == "" {
			return fmt.Errorf("windows signing requires certificate or thumbprint")
		}
		return requireHostTools(host, operation, "signtool.exe")
	case "android":
		if spec.Config.Certificate == "" || spec.Config.KeyAlias == "" || spec.Config.Credential == "" {
			return fmt.Errorf("android signing requires certificate, key_alias, and credential")
		}
		if !host.hasCredential(spec.Config.Credential) {
			return fmt.Errorf("android signing credential %q is unavailable", spec.Config.Credential)
		}
		tool := "jarsigner"
		if spec.Format == "apk" {
			tool = "apksigner"
		}
		return requireHostTools(host, operation, tool)
	case "linux":
		if spec.Config.Certificate == "" {
			return fmt.Errorf("linux signing requires signing.linux.certificate as the PGP key identifier")
		}
		switch spec.Format {
		case "deb":
			return requireHostTools(host, operation, "dpkg-sig")
		case "rpm":
			return requireHostTools(host, operation, "rpmsign")
		default:
			return fmt.Errorf("signing format %q is not supported for linux", spec.Format)
		}
	}
	return nil
}

func requireHostTools(host HostCapabilities, operation string, tools ...string) error {
	for _, tool := range tools {
		if !host.hasTool(tool) {
			return fmt.Errorf("%s requires tool %q", operation, tool)
		}
	}
	return nil
}

func requiredHostName(mask buildinfo.HostMask) string { return buildinfo.RequiredHostName(mask) }

func uniqueSorted(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	if len(result) == 0 {
		return result
	}
	write := 1
	for read := 1; read < len(result); read++ {
		if result[read] != result[write-1] {
			result[write] = result[read]
			write++
		}
	}
	return result[:write]
}

func containsSorted(values []string, want string) bool {
	index := sort.SearchStrings(values, want)
	return index < len(values) && values[index] == want
}
