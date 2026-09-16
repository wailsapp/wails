package pipeline

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func FuzzBuildOutcomeSelection(f *testing.F) {
	seeds := [][]byte{
		{},
		{0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8},
		{255, 255, 255, 255, 255},
	}
	for _, seed := range seeds {
		f.Add(seed)
	}
	targets := supportedTargetNames()
	formats := []string{"", "nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux", "ipa", "aab", "apk", "unknown"}
	verbs := []string{"", "build", "package", "sign", "unknown"}
	destinations := []string{"", "simulator", "device", "unknown"}
	f.Fuzz(func(t *testing.T, input []byte) {
		if len(input) > 128 {
			t.Skip()
		}
		pick := func(values []string, index int) string {
			if len(input) <= index {
				return values[0]
			}
			return values[int(input[index])%len(values)]
		}
		count := 1
		if len(input) != 0 {
			count = int(input[0]%5) + 1
		}
		request := Request{Verb: pick(verbs, 1), Development: len(input) > 2 && input[2]&1 != 0}
		for index := range count {
			name := pick(targets, 3+index)
			target, err := parseTargetName(name)
			require.NoError(t, err)
			request.Targets = append(request.Targets, target)
		}
		for index := 0; index < count && index < 4; index++ {
			request.Formats = append(request.Formats, pick(formats, 9+index))
		}
		config := manifest.Config{}
		if len(input) > 14 && input[14]&1 != 0 {
			config.Selected.Name = "fuzz"
			request.Targets = nil
			request.Formats = nil
			for index := range count {
				config.Selected.Targets = append(config.Selected.Targets, manifest.ProfileTarget{
					Target:      pick(targets, 15+index),
					Formats:     []string{pick(formats, 21+index)},
					Destination: pick(destinations, 27+index),
					Sign:        len(input) > 33+index && input[33+index]&1 != 0,
					Notarize:    len(input) > 39+index && input[39+index]&1 != 0,
				})
			}
		}
		outcomes, err := resolveBuildOutcomes(config, request)
		if err != nil {
			return
		}
		seen := map[string]bool{}
		development := request.Development && config.Selected.Name == ""
		for _, outcome := range outcomes {
			name := outcome.target.OS + "/" + outcome.target.Arch
			if seen[name] {
				t.Fatalf("selection accepted duplicate target %s", name)
			}
			seen[name] = true
			capability, ok := lookupTarget(outcome.target.OS, outcome.target.Arch)
			if !ok {
				t.Fatalf("selection accepted unsupported target %s", name)
			}
			for _, format := range outcome.formats {
				if !capability.SupportsFormat(format, development) {
					t.Fatalf("selection accepted incompatible format %s for %s", format, name)
				}
			}
		}
	})
}

func FuzzPlanValidation(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 1, 1, 1, 1, 1})
	f.Add([]byte{8, 7, 6, 5, 4, 3, 2, 1})
	f.Fuzz(func(t *testing.T, input []byte) {
		if len(input) > 128 {
			t.Skip()
		}
		count := 1
		if len(input) != 0 {
			count = int(input[0]%32) + 1
		}
		plan := Plan{Nodes: make(map[NodeKey]Node, count)}
		for index := range count {
			key := NodeKey(fmt.Sprintf("node-%02d", index))
			dependencies := []NodeKey(nil)
			if index > 0 {
				dependencies = append(dependencies, NodeKey(fmt.Sprintf("node-%02d", index-1)))
				if index > 1 && len(input) > index {
					extra := int(input[index]) % (index - 1)
					dependencies = append(dependencies, NodeKey(fmt.Sprintf("node-%02d", extra)))
				}
			}
			plan.Nodes[key] = Node{Key: key, Kind: CompileApplication, Spec: CompileSpec{}, Dependencies: dependencies, Cache: CacheArtifact}
		}
		plan.Roots = []NodeKey{NodeKey(fmt.Sprintf("node-%02d", count-1))}
		if err := plan.Validate("."); err != nil {
			t.Fatalf("planner rejected generated acyclic graph: %v", err)
		}
		if count > 1 {
			root := plan.Nodes[plan.Roots[0]]
			root.Dependencies = append(root.Dependencies, "missing")
			plan.Nodes[root.Key] = root
			if err := plan.Validate("."); err == nil {
				t.Fatal("planner accepted a missing dependency")
			}
		}
	})
}
