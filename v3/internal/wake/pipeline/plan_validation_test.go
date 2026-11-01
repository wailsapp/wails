package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanValidationRejectsEveryStructuralGraphInvariant(t *testing.T) {
	valid := func() Plan {
		return Plan{Roots: []NodeKey{"root"}, Nodes: map[NodeKey]Node{
			"dependency": {Key: "dependency", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheArtifact},
			"root":       {Key: "root", Kind: CompileApplication, Spec: CompileSpec{}, Dependencies: []NodeKey{"dependency"}, Cache: CacheNever},
		}}
	}
	require.NoError(t, valid().Validate("."))

	tests := map[string]func(*Plan){
		"no roots":       func(plan *Plan) { plan.Roots = nil },
		"missing root":   func(plan *Plan) { plan.Roots = []NodeKey{"missing"} },
		"duplicate root": func(plan *Plan) { plan.Roots = []NodeKey{"root", "root"} },
		"invalid key":    func(plan *Plan) { node := plan.Nodes["root"]; node.Key = "other"; plan.Nodes["root"] = node },
		"invalid spec":   func(plan *Plan) { node := plan.Nodes["root"]; node.Spec = FrontendSpec{}; plan.Nodes["root"] = node },
		"unknown kind":   func(plan *Plan) { node := plan.Nodes["root"]; node.Kind = "Unknown"; plan.Nodes["root"] = node },
		"invalid cache":  func(plan *Plan) { node := plan.Nodes["root"]; node.Cache = "sometimes"; plan.Nodes["root"] = node },
		"escaping output": func(plan *Plan) {
			node := plan.Nodes["root"]
			node.Output = t.TempDir()
			plan.Nodes["root"] = node
		},
		"duplicate output owner": func(plan *Plan) {
			dependency := plan.Nodes["dependency"]
			dependency.Output = "bin/shared"
			plan.Nodes["dependency"] = dependency
			root := plan.Nodes["root"]
			root.Output = "bin/shared"
			plan.Nodes["root"] = root
		},
		"missing artifact": func(plan *Plan) {
			plan.Artifacts = []NodeKey{"missing"}
		},
		"duplicate artifact": func(plan *Plan) {
			node := plan.Nodes["root"]
			node.Output = "bin/root"
			node.Artifact = ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: "linux", Arch: "amd64"}}
			plan.Nodes["root"] = node
			plan.Artifacts = []NodeKey{"root", "root"}
		},
		"untyped artifact": func(plan *Plan) {
			plan.Artifacts = []NodeKey{"root"}
		},
		"missing dependency": func(plan *Plan) {
			node := plan.Nodes["root"]
			node.Dependencies = []NodeKey{"missing"}
			plan.Nodes["root"] = node
		},
		"duplicate dependency": func(plan *Plan) {
			node := plan.Nodes["root"]
			node.Dependencies = []NodeKey{"dependency", "dependency"}
			plan.Nodes["root"] = node
		},
		"cycle": func(plan *Plan) {
			node := plan.Nodes["dependency"]
			node.Dependencies = []NodeKey{"root"}
			plan.Nodes["dependency"] = node
		},
		"unreachable": func(plan *Plan) {
			plan.Nodes["orphan"] = Node{Key: "orphan", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			plan := valid()
			mutate(&plan)
			assert.Error(t, plan.Validate("."))
		})
	}
}

func TestPlanValidationRejectsEmptyPlans(t *testing.T) {
	assert.ErrorContains(t, (Plan{}).Validate("."), "no nodes")
}
