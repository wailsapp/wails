package setupwizard

import (
	"errors"
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestWakeToolRequirementsFollowResolvedCompiler(t *testing.T) {
	for _, test := range []struct {
		name, host, arch, target, cgo string
		console                       bool
		environment                   map[string]string
		installed, want, absent       []string
	}{
		{name: "native CGo", host: "linux", arch: "amd64", target: "linux", cgo: "1", installed: []string{"go", "cc"}, want: []string{"go", "cc", "pkg-config"}, absent: []string{"docker", "npm"}},
		{name: "pure Go cross", host: "linux", arch: "amd64", target: "darwin", cgo: "0", want: []string{"go"}, absent: []string{"docker", "cc", "clang", "zig", "npm"}},
		{name: "CGo cross container", host: "linux", arch: "amd64", target: "darwin", cgo: "1", want: []string{"docker"}, absent: []string{"go", "npm", "cc"}},
		{name: "CGo cross Zig", host: "linux", arch: "amd64", target: "windows", cgo: "1", installed: []string{"zig"}, want: []string{"go", "zig"}, absent: []string{"docker", "npm"}},
		{name: "explicit cross compiler", host: "linux", arch: "amd64", target: "darwin", cgo: "1", environment: map[string]string{"CC": "target-clang"}, want: []string{"go", "target-clang"}, absent: []string{"docker", "npm", "xcrun"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			plan := pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{"compile": {Label: "Compile", Kind: pipeline.CompileApplication, Spec: pipeline.CompileSpec{TargetOS: test.target, TargetArch: "amd64", CGOEnabled: test.cgo, Environment: test.environment, Console: test.console}}}}
			lookup := func(name string) (string, error) {
				if slices.Contains(test.installed, name) {
					return "/tools/" + name, nil
				}
				return "", errors.New("not installed")
			}
			requirements, warnings := wakeToolRequirements(plan, t.TempDir(), test.host, test.arch, lookup)
			if len(warnings) != 0 {
				t.Fatal(warnings)
			}
			names := []string{}
			for _, r := range requirements {
				names = append(names, r.Name)
			}
			for _, want := range test.want {
				if !slices.Contains(names, want) {
					t.Fatalf("missing %s: %v", want, names)
				}
			}
			for _, absent := range test.absent {
				if slices.Contains(names, absent) {
					t.Fatalf("unexpected %s: %v", absent, names)
				}
			}
		})
	}
}
