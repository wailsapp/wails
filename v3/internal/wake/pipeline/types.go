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
	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type NodeKey string
type NodeKind string

const (
	InstallFrontendDependencies NodeKind = "InstallFrontendDependencies"
	GenerateBindings            NodeKind = "GenerateBindings"
	BuildFrontend               NodeKind = "BuildFrontend"
	CompileApplication          NodeKind = "CompileApplication"
	MergeUniversalBinaries      NodeKind = "MergeUniversalBinaries"
	GeneratePlatformAssets      NodeKind = "GeneratePlatformAssets"
	AssembleApplication         NodeKind = "AssembleApplication"
	PackageArtifact             NodeKind = "PackageArtifact"
	SignArtifact                NodeKind = "SignArtifact"
	PublishArtifact             NodeKind = "PublishArtifact"
	CollectArtifacts            NodeKind = "CollectArtifacts"
	RunHook                     NodeKind = "RunHook"
)

type Scope string

const (
	ProjectScope    Scope = "project"
	TargetScope     Scope = "target"
	PackageScope    Scope = "package"
	InvocationScope Scope = "invocation"
)

type CachePolicy string

const (
	CacheArtifact CachePolicy = "artifact"
	CacheReceipt  CachePolicy = "receipt"
	CacheNever    CachePolicy = "never"
)

type ArtifactKind string

const (
	ArtifactBinary   ArtifactKind = "binary"
	ArtifactBundle   ArtifactKind = "bundle"
	ArtifactPackage  ArtifactKind = "package"
	ArtifactSymbols  ArtifactKind = "symbols"
	ArtifactMetadata ArtifactKind = "metadata"
)

type ArtifactIdentity struct {
	Kind      ArtifactKind
	Target    Target
	Format    string
	Signed    bool
	Notarized bool
}

func (a ArtifactIdentity) DisplayKind() string {
	if a.Format != "" {
		return a.Format
	}
	return string(a.Kind)
}

func (a ArtifactIdentity) Empty() bool { return a.Kind == "" }

type ResourceClaims struct {
	CPU       int
	MemoryMB  int64
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
	ExcludeSuffixes   []string
	UseGitIgnore      bool
	SemanticGo        bool
}

type Node struct {
	Key          NodeKey
	Kind         NodeKind
	Label        string
	Scope        Scope
	Dependencies []NodeKey
	Spec         NodeSpec
	Inputs       []InputSpec
	Output       string
	Marker       string
	Cache        CachePolicy
	Claims       ResourceClaims
	EstimateMS   int64
	Artifact     ArtifactIdentity
	Origins      []OriginReference
}

type OriginReference struct {
	Field  string
	Origin manifest.Origin
}

type Plan struct {
	Name      string
	Target    string
	Intent    BuildIntent
	Roots     []NodeKey
	Artifacts []NodeKey
	Nodes     map[NodeKey]Node
}

type BuildIntent struct {
	Command string
	Profile string
	Targets []TargetIntent
}

type TargetIntent struct {
	Target      Target
	Formats     []string
	Sign        bool
	Notarize    bool
	Destination string
}

type Result struct {
	Key              NodeKey
	Status           cache.LookupStatus
	Artifact, Output string
}

type Inspection struct {
	Operations map[NodeKey]OperationInspection
}

type OperationInspection struct {
	Decision string
	Status   cache.LookupStatus
	Inputs   []InputSnapshot
}

type InputSnapshot struct {
	Label  string
	Digest string
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
	Root          string
	Workers       int
	MemoryLimitMB int64
	Force         bool
	Reporter      report.Reporter
}

type Target = buildinfo.Target

type Request struct {
	Verb        string
	TargetOS    string
	TargetArch  string
	Targets     []Target
	Formats     []string
	Development bool
	ExtraTags   []string
	GarbleArgs  []string
	Obfuscated  bool
	resolved    bool
	sign        bool
	notarize    bool
	destination string
}

// NodeSpec is a closed set of typed planner data. It deliberately contains no
// execution method; handlers are selected by NodeKind.
type NodeSpec interface{ nodeSpec() }

type InstallSpec struct {
	Manager, Directory, Command string
	Arguments                   []string
	Environment                 map[string]string
}
type BindingsSpec struct {
	Config     manifest.Bindings
	Tags       []string
	Obfuscated bool
}
type FrontendSpec struct {
	Manager, Directory, Command, Output string
	Arguments                           []string
	Production                          bool
	Environment                         map[string]string
}
type CompileSpec struct {
	TargetOS, TargetArch, Output                 string
	Assets                                       string
	Destination                                  string
	MinimumVersion                               string
	Tags, LinkerFlags, CompilerFlags, GarbleArgs []string
	LocalRoots                                   []string
	Production, Obfuscated, TrimPath, Strip      bool
	VCSInfo                                      bool
	Toolchain                                    string
	ContainerRuntime                             string
	ContainerImage                               string
	Environment                                  map[string]string
}
type MergeSpec struct {
	Inputs []string
	Output string
}
type ComponentBinary struct {
	Arch string
	Path string
}
type AssetsSpec struct {
	TargetOS, TargetArch, Directory string
	MinimumVersion                  string
	Capabilities                    []string
	Project                         manifest.Project
	Associations                    []manifest.Association
	Protocols                       []manifest.Protocol
}
type PackageSpec struct {
	TargetOS, TargetArch, Format, Binary, Assets, Output string
	Binaries                                             []ComponentBinary
	Profile                                              string
	Destination                                          string
	MinimumVersion                                       string
	Config                                               manifest.PackageFormat
	Project                                              manifest.Project
	Capabilities                                         []string
	Associations                                         []manifest.Association
	Protocols                                            []manifest.Protocol
}
type SignSpec struct {
	TargetOS, TargetArch, Format, Input string
	Config                              manifest.SigningPlatform
}
type PublishSpec struct {
	Source      string
	Destination string
}
type ArtifactReference struct {
	Key      NodeKey
	Path     string
	Identity ArtifactIdentity
}
type CollectSpec struct {
	Artifacts []ArtifactReference
	Receipt   string
}
type HookSpec struct {
	Phase                               manifest.HookPhase
	Script, Directory, Profile, Command string
	TargetOS, TargetArch                string
	Scope                               Scope
	ScopeOutput                         string
	DeclaredOutputs                     []string
	ContextVersion                      int
}

func (InstallSpec) nodeSpec()  {}
func (BindingsSpec) nodeSpec() {}
func (FrontendSpec) nodeSpec() {}
func (CompileSpec) nodeSpec()  {}
func (MergeSpec) nodeSpec()    {}
func (AssetsSpec) nodeSpec()   {}
func (PackageSpec) nodeSpec()  {}
func (SignSpec) nodeSpec()     {}
func (PublishSpec) nodeSpec()  {}
func (CollectSpec) nodeSpec()  {}
func (HookSpec) nodeSpec()     {}

func (p Plan) Validate(root string) error {
	if len(p.Nodes) == 0 {
		return fmt.Errorf("wake: plan has no nodes")
	}
	if len(p.Roots) == 0 {
		return fmt.Errorf("wake: plan has no roots")
	}
	rootSet := make(map[NodeKey]bool, len(p.Roots))
	for _, key := range p.Roots {
		if _, ok := p.Nodes[key]; !ok {
			return fmt.Errorf("wake: plan root %s is missing", key)
		}
		if rootSet[key] {
			return fmt.Errorf("wake: duplicate plan root %s", key)
		}
		rootSet[key] = true
	}
	artifactSet := make(map[NodeKey]bool, len(p.Artifacts))
	for _, key := range p.Artifacts {
		node, ok := p.Nodes[key]
		if !ok {
			return fmt.Errorf("wake: plan artifact %s is missing", key)
		}
		if artifactSet[key] {
			return fmt.Errorf("wake: duplicate plan artifact %s", key)
		}
		if node.Output == "" || node.Artifact.Empty() {
			return fmt.Errorf("wake: plan artifact %s has no typed output identity", key)
		}
		artifactSet[key] = true
	}
	owners := map[string]NodeKey{}
	for key, node := range p.Nodes {
		if key == "" || node.Key != key {
			return fmt.Errorf("wake: invalid node key %q", key)
		}
		if !validNodeSpec(node) {
			return fmt.Errorf("wake: node %s kind %s carries invalid spec %T", key, node.Kind, node.Spec)
		}
		if node.Cache != CacheArtifact && node.Cache != CacheReceipt && node.Cache != CacheNever {
			return fmt.Errorf("wake: node %s has invalid cache policy %q", key, node.Cache)
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
		dependencies := make(map[NodeKey]bool, len(node.Dependencies))
		for _, dep := range node.Dependencies {
			if _, ok := p.Nodes[dep]; !ok {
				return fmt.Errorf("wake: node %s depends on missing %s", key, dep)
			}
			if dependencies[dep] {
				return fmt.Errorf("wake: node %s contains duplicate dependency %s", key, dep)
			}
			dependencies[dep] = true
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
	reachable := make(map[NodeKey]bool, len(p.Nodes))
	var markReachable func(NodeKey)
	markReachable = func(key NodeKey) {
		if reachable[key] {
			return
		}
		reachable[key] = true
		for _, dependency := range p.Nodes[key].Dependencies {
			markReachable(dependency)
		}
	}
	for _, key := range p.Roots {
		markReachable(key)
	}
	for _, key := range keys {
		if !reachable[NodeKey(key)] {
			return fmt.Errorf("wake: node %s is unreachable from plan roots", key)
		}
	}
	return nil
}

func validNodeSpec(node Node) bool {
	switch node.Kind {
	case InstallFrontendDependencies:
		_, ok := node.Spec.(InstallSpec)
		return ok
	case GenerateBindings:
		_, ok := node.Spec.(BindingsSpec)
		return ok
	case BuildFrontend:
		_, ok := node.Spec.(FrontendSpec)
		return ok
	case CompileApplication:
		_, ok := node.Spec.(CompileSpec)
		return ok
	case MergeUniversalBinaries:
		_, ok := node.Spec.(MergeSpec)
		return ok
	case GeneratePlatformAssets:
		_, ok := node.Spec.(AssetsSpec)
		return ok
	case AssembleApplication, PackageArtifact:
		_, ok := node.Spec.(PackageSpec)
		return ok
	case SignArtifact:
		_, ok := node.Spec.(SignSpec)
		return ok
	case PublishArtifact:
		_, ok := node.Spec.(PublishSpec)
		return ok
	case CollectArtifacts:
		_, ok := node.Spec.(CollectSpec)
		return ok
	case RunHook:
		_, ok := node.Spec.(HookSpec)
		return ok
	default:
		return false
	}
}
