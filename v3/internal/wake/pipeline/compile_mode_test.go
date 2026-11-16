package pipeline

import (
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestPureGoDesktopHostTargetMatrix(t *testing.T) {
	targets := []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}, {OS: "darwin", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}, {OS: "windows", Arch: "amd64"}, {OS: "windows", Arch: "arm64"}}
	for _, host := range targets {
		for _, target := range targets {
			t.Run(host.OS+"-"+host.Arch+"/"+target.OS+"-"+target.Arch, func(t *testing.T) {
				config := testConfig(t)
				config.Frontend.Disabled = true
				capabilities := NewHostCapabilities(host.OS, host.Arch, []string{"go"}, nil)
				plan, err := PlanBuildForHost(config, Request{Verb: "build", CGOEnabled: "0", SkipBindings: true, TargetOS: target.OS, TargetArch: target.Arch}, capabilities)
				if err != nil {
					t.Fatal(err)
				}
				spec := plan.Nodes[NodeKey("target:"+target.OS+"/"+target.Arch+":compile")].Spec.(CompileSpec)
				if spec.Toolchain != "native" || spec.CGOEnabled != "0" {
					t.Fatalf("unexpected compiler: %+v", spec)
				}
				if target.OS != "windows" && !slices.Contains(spec.Tags, "purego") {
					t.Fatalf("missing backend tag: %v", spec.Tags)
				}
			})
		}
	}
}

func TestCompileModeUsesExecutableEnvironment(t *testing.T) {
	config := testConfig(t)
	config.Build.Environment = map[string]string{"CGO_ENABLED": "0"}
	config.Binaries = []manifest.Binary{{Name: "helper", Package: "./cmd/helper", Kind: "cli", Environment: map[string]string{"CGO_ENABLED": "1", "CC": "custom-cross-cc"}}}
	for _, development := range []bool{false, true} {
		plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "arm64", Development: development, CGOEnabled: "1"})
		if err != nil {
			t.Fatal(err)
		}
		for _, node := range plan.Nodes {
			if node.Kind != CompileApplication {
				continue
			}
			spec := node.Spec.(CompileSpec)
			want := "0"
			if spec.Package == "./cmd/helper" {
				want = "1"
			}
			if spec.CGOEnabled != want {
				t.Fatalf("development=%t package=%s: CGO_ENABLED=%s, want %s", development, spec.Package, spec.CGOEnabled, want)
			}
		}
	}
}

func TestExplicitWindowsCGoNeedsCrossCompiler(t *testing.T) {
	config := testConfig(t)
	request := Request{Verb: "build", TargetOS: "windows", TargetArch: "arm64", CGOEnabled: "1"}
	plan, err := PlanBuildForHost(config, request, testHost("linux", "amd64", "zig"))
	if err != nil {
		t.Fatal(err)
	}
	if spec := plan.Nodes["target:windows/arm64:compile"].Spec.(CompileSpec); spec.Toolchain != "zig" || spec.CGOEnabled != "1" {
		t.Fatalf("unexpected CGo toolchain: %+v", spec)
	}
	if _, err := PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", []string{"go", "npm", "cc"}, nil)); err == nil {
		t.Fatal("CGo cross build accepted without a cross compiler")
	}
}

func TestCompileModeRejectsInvalidValues(t *testing.T) {
	for _, mode := range []string{"true", "2", " 0"} {
		if _, err := PlanBuild(testConfig(t), Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", CGOEnabled: mode}); err == nil {
			t.Fatalf("accepted CGO_ENABLED=%q", mode)
		}
	}
}

func TestPureGoMultiTargetPlanSharesCompilerFreeBindings(t *testing.T) {
	plan, err := PlanBuild(testConfig(t), Request{Verb: "build", CGOEnabled: "0", Targets: []Target{{OS: "linux", Arch: "arm64"}, {OS: "windows", Arch: "amd64"}, {OS: "darwin", Arch: "universal"}}})
	if err != nil {
		t.Fatal(err)
	}
	spec := plan.Nodes["frontend:bindings"].Spec.(BindingsSpec)
	if spec.CGOEnabled != "0" || !slices.Contains(spec.Tags, "purego") {
		t.Fatalf("bindings require a host C compiler: %+v", spec)
	}
}

func TestExplicitCrossCompilerIsUsedWithoutContainers(t *testing.T) {
	config := testConfig(t)
	request := Request{Verb: "build", TargetOS: "darwin", TargetArch: "arm64", CGOEnabled: "1", CC: "custom-apple-clang"}
	plan, err := PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", []string{"go", "npm"}, nil))
	if err != nil {
		t.Fatal(err)
	}
	spec := plan.Nodes["target:darwin/arm64:compile"].Spec.(CompileSpec)
	if spec.Toolchain != "native" || spec.Environment["CC"] != request.CC {
		t.Fatalf("explicit cross compiler was ignored: %+v", spec)
	}
}

func TestWindowsCGoUsesZigWhenDefaultCompilerIsUnavailable(t *testing.T) {
	request := Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64", CGOEnabled: "1"}
	host := NewHostCapabilities("windows", "amd64", []string{"go", "npm", "zig"}, nil)
	plan, err := PlanBuildForHost(testConfig(t), request, host)
	if err != nil {
		t.Fatal(err)
	}
	if spec := plan.Nodes["target:windows/amd64:compile"].Spec.(CompileSpec); spec.Toolchain != "zig" {
		t.Fatalf("missing GCC should select available Zig: %+v", spec)
	}
}
