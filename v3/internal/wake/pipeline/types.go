// Package pipeline plans and executes Wails build intent as a static typed DAG.
// The public interface is intentionally small: Plan resolves user intent, and
// Executor runs the resulting immutable graph through one Handler adapter.
package pipeline

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type NodeKey string
type NodeKind string

const (
	InstallFrontendDependencies NodeKind = "InstallFrontendDependencies"
	GenerateBindings            NodeKind = "GenerateBindings"
	GenerateIcons               NodeKind = "GenerateIcons"
	BuildFrontend               NodeKind = "BuildFrontend"
	CompileApplication          NodeKind = "CompileApplication"
	GeneratePlatformAssets      NodeKind = "GeneratePlatformAssets"
	AssembleBundle              NodeKind = "AssembleBundle"
	PackageArtifact             NodeKind = "PackageArtifact"
	SignArtifact                NodeKind = "SignArtifact"
	RunHook                     NodeKind = "RunHook"
)

type Scope string

const (
	ProjectScope Scope = "project"
	TargetScope  Scope = "target"
	PackageScope Scope = "package"
)

type CachePolicy string

const (
	CacheArtifact CachePolicy = "artifact"
	CacheReceipt  CachePolicy = "receipt"
	CacheNever    CachePolicy = "never"
)

type ResourceClaims struct {
	CPU       int
	MemoryMB  int
	Exclusive string
}

type InputSpec struct {
	Label             string
	Root              string
	Files             []string
	IncludeAll        bool
	IncludeNames      []string
	IncludeExtensions []string
	ExcludeDirs       []string
	SemanticGo        bool
}

type Node struct {
	Key          NodeKey
	Kind         NodeKind
	Label        string
	Scope        Scope
	Dependencies []NodeKey
	Spec         any
	Inputs       []InputSpec
	Output       string
	Marker       string
	Cache        CachePolicy
	Claims       ResourceClaims
	EstimateMS   int64
	ArtifactKind string
}

type Plan struct {
	Name   string
	Target string
	Roots  []NodeKey
	Nodes  map[NodeKey]Node
}

type Result struct {
	Key              NodeKey
	Status           cache.LookupStatus
	Artifact, Output string
}

type RunResult struct {
	Output string
	Detail string
}

type Handler interface {
	Identity(context.Context, Node) (string, error)
	Run(context.Context, Node) (RunResult, error)
}

type ExecuteOptions struct {
	Root     string
	Workers  int
	Force    bool
	Reporter report.Reporter
}

type Request struct {
	Verb        string
	TargetOS    string
	TargetArch  string
	Formats     []string
	Development bool
	ExtraTags   []string
	Obfuscated  bool
}

type InstallSpec struct{ Manager, Directory, Command string }
type BindingsSpec struct {
	Config     manifest.Bindings
	Tags       []string
	Obfuscated bool
}
type FrontendSpec struct {
	Manager, Directory, Command, Output string
	Production                          bool
}
type CompileSpec struct {
	TargetOS, TargetArch, Output                 string
	Assets                                       string
	Variant                                      string
	MinimumVersion                               string
	Tags, LinkerFlags, CompilerFlags, GarbleArgs []string
	Production, Obfuscated, TrimPath, Strip      bool
}
type AssetsSpec struct {
	TargetOS, TargetArch, Directory string
	MinimumVersion                  string
	Project                         manifest.Project
	Associations                    []manifest.Association
	Protocols                       []manifest.Protocol
}
type BundleSpec struct{ TargetOS, TargetArch, Binary, Assets, Output string }
type PackageSpec struct {
	TargetOS, TargetArch, Format, Binary, Assets, Output string
	Variant                                              string
	MinimumVersion                                       string
	Config                                               manifest.PackageFormat
	Project                                              manifest.Project
}
type SignSpec struct {
	TargetOS, Format, Input string
	Config                  manifest.SigningPlatform
}
type HookSpec struct {
	Phase, TargetOS, TargetArch, Profile, Output string
	Hook                                         manifest.Hook
}

func (p Plan) Validate(root string) error {
	if len(p.Nodes) == 0 {
		return fmt.Errorf("wake: plan has no nodes")
	}
	owners := map[string]NodeKey{}
	for key, node := range p.Nodes {
		if key == "" || node.Key != key {
			return fmt.Errorf("wake: invalid node key %q", key)
		}
		if node.Output != "" {
			clean := filepath.Clean(node.Output)
			if filepath.IsAbs(clean) || clean == ".." || len(clean) >= 3 && clean[:3] == ".."+string(filepath.Separator) {
				return fmt.Errorf("wake: node %s output escapes project: %s", key, node.Output)
			}
			if owner, exists := owners[clean]; exists {
				return fmt.Errorf("wake: output %s owned by both %s and %s", clean, owner, key)
			}
			owners[clean] = key
		}
		for _, dep := range node.Dependencies {
			if _, ok := p.Nodes[dep]; !ok {
				return fmt.Errorf("wake: node %s depends on missing %s", key, dep)
			}
		}
	}
	state := map[NodeKey]uint8{}
	var visit func(NodeKey) error
	visit = func(key NodeKey) error {
		if state[key] == 1 {
			return fmt.Errorf("wake: dependency cycle at %s", key)
		}
		if state[key] == 2 {
			return nil
		}
		state[key] = 1
		for _, dep := range p.Nodes[key].Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[key] = 2
		return nil
	}
	keys := make([]string, 0, len(p.Nodes))
	for key := range p.Nodes {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := visit(NodeKey(key)); err != nil {
			return err
		}
	}
	return nil
}
