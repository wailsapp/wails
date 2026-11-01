package setupwizard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/packagesign"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (w *Wizard) handleWakeRequirements(rw http.ResponseWriter, r *http.Request) {
	var selected WakeSetupConfig
	if !w.decodeWakeSetup(rw, r, &selected) {
		return
	}
	response := wakeRequirementsResponse{Requirements: []DependencyStatus{}, Errors: []string{}, Warnings: []string{}, Command: wakeSetupCommand(selected)}
	previewRoot, _, _, config, err := w.prepareInitWake(selected)
	if previewRoot != "" {
		defer os.RemoveAll(previewRoot)
	}
	if err == nil {
		config, err = signingconfig.ApplyDefaults(config)
	}
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
	} else {
		// Planning checks intent and inputs, but never executes hooks. Commit
		// injection is irrelevant to tool availability and must not invoke git here.
		config.Project.CommitSymbol = ""
		config.Binaries = slices.Clone(config.Binaries)
		for i := range config.Binaries {
			config.Binaries[i].CommitSymbol = ""
		}
		request := pipeline.Request{Verb: "build", SkipBuild: selected.Prebuilt, CGOEnabled: os.Getenv("CGO_ENABLED"), CC: os.Getenv("CC")}
		plan, planErr := pipeline.PlanBuild(config, request)
		if planErr != nil {
			response.Errors = append(response.Errors, planErr.Error())
		} else {
			response.Requirements, response.Warnings = wakeToolRequirements(plan, config.Root, runtime.GOOS, runtime.GOARCH, exec.LookPath)
			if _, hostErr := pipeline.PlanBuildForCurrentHost(config, request, wakeCredentialNames(config.Signing)...); hostErr != nil {
				response.Warnings = append(response.Warnings, hostErr.Error())
				response.Requirements = append(response.Requirements, DependencyStatus{Name: "Build configuration on this host", Installed: false, Required: true, Status: "needs_config", Message: hostErr.Error()})
			}
			response.Requirements = append(response.Requirements, wakeLinuxSDKRequirements(r.Context(), plan)...)
			if selected.Prebuilt {
				response.Warnings = append(response.Warnings, "Prebuilt mode only packages existing binaries from build.output; binary architecture and format are verified when the command runs.")
			}
			response.Warnings = append(response.Warnings, "The build backend is project-wide. Other profiles that inherit build settings may also use the selected backend.")
			if selected.Backend == "purego" {
				response.Warnings = append(response.Warnings, "Pure Go is experimental. GUI applications still require the target operating system's graphical libraries at runtime.")
			}
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(response)
}

func wakeCredentialNames(signing manifest.Signing) []string {
	var names []string
	for _, s := range []manifest.SigningPlatform{signing.Windows, signing.Darwin, signing.Linux, signing.IOS, signing.Android} {
		if s.Credential != "" {
			names = append(names, s.Credential)
		}
	}
	return names
}

func wakeToolRequirements(plan pipeline.Plan, root, hostOS, hostArch string, lookPath func(string) (string, error)) ([]DependencyStatus, []string) {
	requirements := map[string]DependencyStatus{}
	warnings := []string{}
	add := func(name, reason string) {
		if name == "" {
			return
		}
		item, exists := requirements[name]
		if !exists {
			_, err := lookPath(name)
			item = DependencyStatus{Name: name, Installed: err == nil, Required: true, Status: "not_installed"}
			if item.Installed {
				item.Status = "installed"
			}
		}
		if !strings.Contains(item.Message, reason) {
			if item.Message != "" {
				item.Message += "; "
			}
			item.Message += reason
		}
		requirements[name] = item
	}
	tools := []string{}
	for _, name := range []string{"go", "cc", "gcc", "clang", "zig", "docker", "podman"} {
		if _, err := lookPath(name); err == nil {
			tools = append(tools, name)
		}
	}
	host := pipeline.NewHostCapabilities(hostOS, hostArch, tools, nil)
	keys := make([]string, 0, len(plan.Nodes))
	for key := range plan.Nodes {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	for _, key := range keys {
		node := plan.Nodes[pipeline.NodeKey(key)]
		switch spec := node.Spec.(type) {
		case pipeline.InstallSpec:
			command := spec.Manager
			if len(spec.Arguments) > 0 {
				command = spec.Arguments[0]
			}
			add(command, node.Label)
		case pipeline.FrontendSpec:
			command := spec.Manager
			if len(spec.Arguments) > 0 {
				command = spec.Arguments[0]
			}
			add(command, node.Label)
		case pipeline.BindingsSpec:
			add("go", node.Label)
		case pipeline.CompileSpec:
			toolchain, err := pipeline.CompileToolchain(spec, host)
			if err != nil {
				suggested := append([]string{}, tools...)
				if spec.Toolchain == "zig" {
					suggested = append(suggested, "zig")
				} else {
					suggested = append(suggested, "docker")
				}
				toolchain, err = pipeline.CompileToolchain(spec, pipeline.NewHostCapabilities(hostOS, hostArch, suggested, nil))
			}
			if err != nil {
				warnings = append(warnings, err.Error())
				requirements["Toolchain for "+spec.TargetOS+"/"+spec.TargetArch] = DependencyStatus{Name: "Toolchain for " + spec.TargetOS + "/" + spec.TargetArch, Installed: false, Required: true, Status: "needs_config", Message: err.Error()}
				continue
			}
			reason := node.Label + " (" + spec.TargetOS + "/" + spec.TargetArch + ")"
			if spec.UPXLevel > 0 {
				add("upx", reason)
			}
			if toolchain == "docker" {
				container := "docker"
				if slices.Contains(tools, "podman") && !slices.Contains(tools, "docker") {
					container = "podman"
				}
				add(container, reason+" using the wails-cross image")
				continue
			}
			add("go", reason)
			if spec.Obfuscated {
				add("garble", reason)
			}
			if toolchain == "zig" {
				add("zig", reason)
			}
			if spec.CGOEnabled == "1" && toolchain == "native" {
				if command := strings.Fields(spec.Environment["CC"]); len(command) > 0 {
					add(command[0], reason)
				} else if hostOS == "darwin" && (spec.TargetOS == "darwin" || spec.TargetOS == "ios") {
					add("xcrun", reason)
				}
				if spec.TargetOS == "linux" || spec.TargetOS == "windows" {
					compiler := "cc"
					if spec.TargetOS == "windows" {
						compiler = "gcc"
					}
					if command := strings.Fields(spec.Environment["CC"]); len(command) > 0 {
						compiler = command[0]
					} else {
						for _, candidate := range []string{compiler, "gcc", "clang"} {
							if _, err := lookPath(candidate); err == nil {
								compiler = candidate
								break
							}
						}
					}
					add(compiler, reason)
					if spec.TargetOS == "linux" && !spec.Console {
						add("pkg-config", reason+" and target GTK/WebKit development files")
					}
				}
				if spec.TargetOS == "android" {
					for _, dependency := range wakeAndroidSDKRequirements() {
						requirements[dependency.Name] = dependency
					}
				}
			}
		case pipeline.PackageSpec:
			if node.Kind == pipeline.AssembleApplication && spec.TargetOS == "ios" {
				add("xcrun", node.Label)
				add("codesign", node.Label)
			}
			if node.Kind == pipeline.PackageArtifact {
				if capability, ok := buildinfo.LookupFormat(spec.Format); ok {
					for i := 0; i < int(capability.ToolCount); i++ {
						add(capability.RequiredTool(i), node.Label)
					}
				}
			}
		case pipeline.SignSpec:
			for _, tool := range wakeSigningTools(spec) {
				add(tool, node.Label)
			}
		case pipeline.PrebuiltSpec:
			path := filepath.Join(root, filepath.FromSlash(spec.Source))
			info, err := os.Stat(path)
			installed := err == nil && info.Mode().IsRegular()
			status := "not_installed"
			if installed {
				status = "installed"
			}
			requirements["Prebuilt "+spec.Source] = DependencyStatus{Name: "Prebuilt " + spec.Source, Installed: installed, Status: status, Required: true, Message: node.Label + "; verified during packaging"}
		}
	}
	result := make([]DependencyStatus, 0, len(requirements))
	for _, item := range requirements {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, warnings
}

func wakeAndroidSDKRequirements() []DependencyStatus {
	sdk := os.Getenv("ANDROID_HOME")
	if sdk == "" {
		sdk = os.Getenv("ANDROID_SDK_ROOT")
	}
	ndk := os.Getenv("ANDROID_NDK_HOME")
	if !dirExists(ndk) && sdk != "" {
		entries, _ := os.ReadDir(filepath.Join(sdk, "ndk"))
		for _, entry := range entries {
			if entry.IsDir() {
				ndk = filepath.Join(sdk, "ndk", entry.Name())
				break
			}
		}
	}
	result := []DependencyStatus{}
	for _, item := range []struct{ name, path string }{{"Android SDK", sdk}, {"Android NDK", ndk}} {
		installed := dirExists(item.path)
		status := "not_installed"
		if installed {
			status = "installed"
		}
		result = append(result, DependencyStatus{Name: item.name, Installed: installed, Required: true, Status: status, Message: "Required for Android source builds; configure ANDROID_HOME and ANDROID_NDK_HOME", HelpURL: "https://developer.android.com/studio"})
	}
	return result
}

func wakeSigningTools(spec pipeline.SignSpec) []string {
	switch spec.TargetOS {
	case "darwin":
		tools := []string{}
		if spec.Config.Certificate == "" && spec.Config.Identity != "-" {
			tools = append(tools, "codesign")
		}
		if spec.Config.Notarize && spec.Config.NotarizationKey == "" {
			tools = append(tools, "xcrun")
		}
		if spec.Format == "dmg" && !spec.BundleSigned {
			tools = append(tools, "hdiutil")
		}
		return tools
	case "windows":
		if spec.Config.Certificate == "" {
			return []string{"signtool.exe"}
		}
	case "ios":
		return []string{"codesign"}
	case "android":
		if spec.Format == "apk" {
			return []string{"apksigner"}
		}
		return []string{"jarsigner"}
	case "linux":
		if _, file := packagesign.KeyFile("", spec.Config.Certificate); file {
			return nil
		}
		if spec.Format == "deb" {
			return []string{"dpkg-sig"}
		}
		if spec.Format == "rpm" {
			return []string{"rpmsign"}
		}
	}
	return nil
}

func wakeLinuxSDKRequirements(ctx context.Context, plan pipeline.Plan) []DependencyStatus {
	if runtime.GOOS != "linux" {
		return nil
	}
	result := []DependencyStatus{}
	seen := map[string]bool{}
	for _, node := range plan.Nodes {
		spec, ok := node.Spec.(pipeline.CompileSpec)
		if !ok || spec.Console || spec.TargetOS != "linux" || spec.TargetArch != runtime.GOARCH || spec.CGOEnabled != "1" || spec.Toolchain == "docker" || spec.Environment["CC"] != "" {
			continue
		}
		libraries := []string{"gtk4", "webkitgtk-6.0"}
		if slices.Contains(spec.Tags, "gtk3") {
			libraries = []string{"gtk+-3.0", "webkit2gtk-4.1"}
		}
		for _, library := range libraries {
			if seen[library] {
				continue
			}
			seen[library] = true
			checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			err := exec.CommandContext(checkCtx, "pkg-config", "--exists", library).Run()
			cancel()
			status := "not_installed"
			if err == nil {
				status = "installed"
			}
			result = append(result, DependencyStatus{Name: library + " development files", Installed: err == nil, Required: true, Status: status, Message: fmt.Sprintf("Native CGo GUI compilation for %s/%s", spec.TargetOS, spec.TargetArch)})
		}
	}
	return result
}
