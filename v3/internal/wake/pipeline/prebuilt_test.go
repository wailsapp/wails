package pipeline

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestPrebuiltPlansDoNotRequireBuildTools(t *testing.T) {
	for _, target := range []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "arm64"}, {OS: "darwin", Arch: "universal"}} {
		t.Run(target.OS+"/"+target.Arch, func(t *testing.T) {
			config := testConfig(t)
			config.Binaries = []manifest.Binary{{Name: "helper", Package: "./helper", Kind: "cli"}}
			config.Hooks = allPlannerHooks()
			format := map[string]string{"linux": "deb", "windows": "msix", "darwin": "dmg"}[target.OS]
			request := Request{Verb: "package", TargetOS: target.OS, TargetArch: target.Arch, Formats: []string{format}, SkipBuild: true}
			plan, err := PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", nil, nil))
			if err != nil {
				t.Fatal(err)
			}
			imports, hooks := 0, 0
			for _, node := range plan.Nodes {
				switch node.Kind {
				case ImportBinary:
					imports++
					spec := node.Spec.(PrebuiltSpec)
					if spec.Target != target || !strings.HasPrefix(spec.Source, "bin/") || len(node.Dependencies) != 0 {
						t.Fatalf("unexpected import: %#v", node)
					}
				case CompileApplication, MergeUniversalBinaries, InstallFrontendDependencies, BuildFrontend, GenerateBindings:
					t.Fatalf("prebuilt plan contains build node %s", node.Kind)
				case RunHook:
					phase := node.Spec.(HookSpec).Phase
					if phase == manifest.BeforeBuild || phase == manifest.AfterBuild {
						t.Fatalf("prebuilt plan contains build hook %s", phase)
					}
					hooks++
				case PublishArtifact:
					if node.Artifact.Kind == ArtifactBinary {
						t.Fatal("prebuilt input would be republished")
					}
				}
			}
			if imports != 2 || hooks != 2 {
				t.Fatalf("imports=%d hooks=%d; want 2 imports and package hooks", imports, hooks)
			}
			if len(config.Hooks) != 6 || config.Frontend.Disabled {
				t.Fatal("planning mutated the caller's configuration")
			}
		})
	}
}

func TestPrebuiltMultiTargetAndRunPaths(t *testing.T) {
	config := testConfig(t)
	config.Project.Kind = "cli"
	config.Build.OutputDirectory = "release"
	config.Binaries = []manifest.Binary{{Name: "helper", Package: ".", Kind: "cli"}}
	plan, err := PlanBuild(config, Request{Verb: "build", SkipBuild: true, Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "arm64"}}})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range plan.Artifacts {
		node := plan.Nodes[key]
		spec, ok := node.Spec.(PrebuiltSpec)
		if !ok {
			t.Fatalf("expected imported runnable, got %T", node.Spec)
		}
		want := filepath.ToSlash(filepath.Join("release", spec.Target.OS+"-"+spec.Target.Arch, filepath.Base(spec.Output)))
		if spec.Source != want {
			t.Fatalf("source=%s want %s", spec.Source, want)
		}
	}
	for _, target := range []string{"linux/amd64", "windows/arm64"} {
		parts := strings.Split(target, "/")
		key := NodeKey("binary:" + target + ":helper:" + parts[1])
		if plan.Nodes[key].Kind != ImportBinary {
			t.Fatalf("missing imported companion for %s", target)
		}
	}
}

func TestPrebuiltRejectsMobileAndDevelopment(t *testing.T) {
	for _, request := range []Request{
		{TargetOS: "android", TargetArch: "arm64", SkipBuild: true},
		{TargetOS: "ios", TargetArch: "arm64", SkipBuild: true},
		{TargetOS: "linux", TargetArch: "amd64", SkipBuild: true, Development: true},
	} {
		if _, err := PlanBuild(testConfig(t), request); err == nil || !strings.Contains(err.Error(), "--skip-build") {
			t.Fatalf("request=%#v error=%v", request, err)
		}
	}
}

func TestPrebuiltPreservesSignedPublicationAndHooks(t *testing.T) {
	config := testConfig(t)
	config.Hooks = allPlannerHooks()
	config.Package.Darwin.Formats = nil
	plan, err := PlanBuild(config, Request{Verb: "sign", TargetOS: "darwin", TargetArch: "arm64", SkipBuild: true})
	if err != nil {
		t.Fatal(err)
	}
	afterSign := NodeKey("hook:after_sign:darwin-arm64")
	if _, ok := plan.Nodes[afterSign]; !ok {
		t.Fatal("after_sign hook lost from prebuilt signing plan")
	}
	for _, key := range plan.Artifacts {
		node := plan.Nodes[key]
		if node.Kind != PublishArtifact || !node.Artifact.Signed || !strings.HasSuffix(node.Output, ".signed.app") {
			t.Fatalf("signed publication lost: %#v", node)
		}
	}
}
