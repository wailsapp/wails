package pipeline

import (
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type frontendStages struct {
	config       manifest.Config
	request      Request
	tags         []string
	obfuscated   bool
	localInputs  []InputSpec
	dependencies []NodeKey
	add          func(Node) NodeKey
}

func (stages frontendStages) plan() NodeKey {
	var install NodeKey
	if !stages.config.Frontend.Disabled {
		install = stages.add(Node{Key: "frontend:install", Kind: InstallFrontendDependencies, Label: "Install frontend dependencies", Scope: ProjectScope, Dependencies: stages.dependencies,
			Spec:   InstallSpec{Manager: stages.config.Frontend.PackageManager, Directory: stages.config.Frontend.Directory, Command: stages.config.Frontend.InstallCommand, Arguments: append([]string(nil), stages.config.Frontend.Install...), Environment: cloneStringMap(stages.config.Frontend.Environment)},
			Inputs: []InputSpec{{Label: "frontend-install", Files: frontendInstallFiles(stages.config)}}, Marker: filepath.ToSlash(filepath.Join(stages.config.Frontend.Directory, "node_modules")), Cache: CacheReceipt, Claims: ResourceClaims{CPU: 1, MemoryMB: 512}, EstimateMS: 1800})
	}
	var frontend NodeKey
	if !stages.config.Frontend.Disabled {
		bindingsOut := filepath.ToSlash(filepath.Join(stages.config.Frontend.Directory, stages.config.Frontend.Bindings.OutputDirectory))
		bindingInputs := []InputSpec{
			{Label: "go-binding-api", Root: ".", SemanticGo: true, IncludeNames: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, IncludeExtensions: goSourceExtensions(), ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", stages.config.Frontend.Directory, "node_modules"}, ExcludeSuffixes: []string{"_test.go"}, UseGitIgnore: stages.request.Development && stages.config.Dev.UseGitIgnore},
			{Label: "go-module", Files: goMetadataFiles(stages.config.Root)},
		}
		for _, input := range stages.localInputs {
			bindingInputs = append(bindingInputs, InputSpec{Label: "go-binding-local-api", Root: input.Root, SemanticGo: true, IncludeNames: append([]string(nil), input.IncludeNames...), IncludeExtensions: append([]string(nil), input.IncludeExtensions...), ExcludeDirs: append([]string(nil), input.ExcludeDirs...), ExcludeSuffixes: append([]string(nil), input.ExcludeSuffixes...), UseGitIgnore: input.UseGitIgnore})
		}
		var bindings NodeKey
		if !stages.request.SkipBindings {
			bindings = stages.add(Node{Key: "frontend:bindings", Kind: GenerateBindings, Label: "Generate bindings", Scope: ProjectScope, Dependencies: stages.dependencies,
				Spec:   BindingsSpec{Package: stages.config.Project.Package, Config: stages.config.Frontend.Bindings, Tags: primaryBuildTags(stages.config, stages.tags), Obfuscated: stages.obfuscated},
				Inputs: bindingInputs,
				Output: bindingsOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "legacy-command-adapter"}, EstimateMS: 1000})
		}
		frontendDeps := []NodeKey{install}
		if bindings != "" {
			frontendDeps = append(frontendDeps, bindings)
		}
		frontendOut := frontendOutputPath(stages.config.Frontend.Directory, stages.config.Frontend.OutputDirectory)
		frontendExclude := filepath.ToSlash(filepath.Clean(stages.config.Frontend.OutputDirectory))
		if frontendExclude == stages.config.Frontend.Directory || strings.HasPrefix(frontendExclude, stages.config.Frontend.Directory+"/") {
			frontendExclude = strings.TrimPrefix(strings.TrimPrefix(frontendExclude, stages.config.Frontend.Directory), "/")
		}
		bindingsExclude := stages.config.Frontend.Bindings.OutputDirectory
		if stages.request.SkipBindings {
			bindingsExclude = ""
		}
		frontendInputs := []InputSpec{{Label: "frontend-source", Root: stages.config.Frontend.Directory, IncludeAll: true, ExcludeDirs: []string{".git", ".wails", "node_modules", bindingsExclude, frontendExclude}}}
		if stages.request.Development {
			// The persistent frontend Dev process owns source changes through HMR.
			// The finite frontend Node only bootstraps an embeddable output for Go;
			// keeping source out of its Action Key prevents a later backend edit from
			// rebuilding the frontend merely because HMR already handled a UI edit.
			frontendInputs = []InputSpec{{Label: "frontend-dev-bootstrap", Files: frontendInstallFiles(stages.config)}}
		}
		frontend = stages.add(Node{Key: "frontend:build", Kind: BuildFrontend, Label: "Build frontend", Scope: ProjectScope, Dependencies: frontendDeps,
			Spec:   FrontendSpec{Manager: stages.config.Frontend.PackageManager, Directory: stages.config.Frontend.Directory, Command: choose(stages.request.Development, stages.config.Frontend.BuildCommand+":dev", stages.config.Frontend.BuildCommand), Arguments: append([]string(nil), stages.config.Frontend.Build...), Output: stages.config.Frontend.OutputDirectory, Production: !stages.request.Development, Environment: cloneStringMap(stages.config.Frontend.Environment)},
			Inputs: frontendInputs,
			Output: frontendOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 2, MemoryMB: 1024}, EstimateMS: 900})
	}
	return frontend
}
