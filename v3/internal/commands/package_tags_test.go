package commands

import (
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestPackageTagsReachManifestAndLegacyBuilds(t *testing.T) {
	t.Setenv(puregoEnvVar, "1")
	t.Setenv(mcpEnvVar, "")
	for _, active := range []bool{true, false} {
		t.Run(map[bool]string{true: "manifest", false: "legacy"}[active], func(t *testing.T) {
			called := false
			operations := manifestCommandOperations{
				active: func() (bool, error) { return active, nil },
				run: func(options manifestRunOptions) error {
					called = true
					if !slices.Equal(options.Tags, []string{"gtk3", "purego"}) {
						t.Fatalf("tags = %v", options.Tags)
					}
					return nil
				},
				task: func(verb string, args []string) error {
					called = true
					if verb != "package" || !slices.Contains(args, "EXTRA_TAGS=gtk3,purego") {
						t.Fatalf("legacy request = %q %v", verb, args)
					}
					return nil
				},
			}
			if err := packageWithOperations(&flags.Package{Tags: "gtk3,purego"}, nil, operations); err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("package request was not dispatched")
			}
		})
	}
}
