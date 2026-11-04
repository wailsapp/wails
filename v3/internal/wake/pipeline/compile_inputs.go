package pipeline

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
)

// Retain the source roots so generated imports can be resolved after frontend
// work finishes, rather than freezing an incomplete graph during planning.
func refineCompileInputs(plan *Plan, root string) {
	for key, node := range plan.Nodes {
		spec, ok := node.Spec.(CompileSpec)
		if !ok {
			continue
		}
		allowFrontendReuse := spec.CGOEnabled == "0"
		node.OrderOnlyDependencies = nil
		node.Inputs = append([]InputSpec(nil), node.Inputs...)
		for i := range node.Inputs {
			if node.Inputs[i].Label == "go-sources" {
				node.Inputs[i].Root = root
				allowFrontendReuse = allowFrontendReuse && !node.Inputs[i].UseGitIgnore
			}
		}
		plan.Nodes[key] = node
		// Ignored generated Go files may still be consumed by the compiler.
		if !allowFrontendReuse {
			continue
		}
		if _, complete := collectExecutableInputs(root, spec.Package, node.Inputs); !complete {
			continue
		}
		for _, dependency := range node.Dependencies {
			if plan.Nodes[dependency].Kind == BuildFrontend {
				node.OrderOnlyDependencies = append(node.OrderOnlyDependencies, dependency)
			}
		}
		plan.Nodes[key] = node
	}
}

func resolveCompileInputs(node Node) Node {
	spec, ok := node.Spec.(CompileSpec)
	if !ok {
		return node
	}
	for _, input := range node.Inputs {
		if input.Label == "go-sources" && filepath.IsAbs(input.Root) {
			var complete bool
			node.Inputs, complete = collectExecutableInputs(input.Root, spec.Package, node.Inputs)
			if !complete {
				node.OrderOnlyDependencies = nil
			}
			break
		}
	}
	return node
}

func executableInputs(root, entry string, original []InputSpec) []InputSpec {
	inputs, _ := collectExecutableInputs(root, entry, original)
	return inputs
}

// Local replacement/workspace modules stay conservatively tracked. Within the
// project module, track imported packages (including inactive platform files),
// embeds and native sources; incomplete graphs keep the original broad inputs.
func collectExecutableInputs(root, entry string, original []InputSpec) ([]InputSpec, bool) {
	var source InputSpec
	found := false
	for _, input := range original {
		if input.Label == "go-sources" {
			source = input
			found = true
			break
		}
	}
	if !found {
		return original, false
	}
	moduleRoot, modulePath := projectModule(root)
	if modulePath == "" {
		return original, false
	}
	directory := filepath.Join(root, entry)
	graph := packageInputsGraph{root: root, moduleRoot: moduleRoot, modulePath: modulePath, source: source, seen: map[string]bool{}}
	if !graph.collect(directory) {
		return original, false
	}
	native := source
	native.Label, native.IncludeGoEmbed = "go-native-sources", false
	native.IncludeExtensions = nil
	for _, extension := range source.IncludeExtensions {
		if extension != ".go" {
			native.IncludeExtensions = append(native.IncludeExtensions, extension)
		}
	}
	result := []InputSpec{native}
	sort.Slice(graph.inputs, func(i, j int) bool { return graph.inputs[i].Root < graph.inputs[j].Root })
	result = append(result, graph.inputs...)
	for _, input := range original {
		if input.Label != "go-sources" {
			result = append(result, input)
		}
	}
	vendor := filepath.Join(moduleRoot, "vendor")
	if info, err := os.Stat(vendor); err == nil && info.IsDir() {
		input := source
		input.Label, input.Root = "go-vendor", vendor
		input.IncludeNames = append(append([]string(nil), input.IncludeNames...), "modules.txt")
		result = append(result, input)
	}
	return result, true
}

func projectModule(root string) (string, string) {
	for _, metadata := range goMetadataFiles(root) {
		if filepath.Base(metadata) != "go.mod" {
			continue
		}
		if !filepath.IsAbs(metadata) {
			metadata = filepath.Join(root, metadata)
		}
		data, err := os.ReadFile(metadata)
		if err != nil {
			return "", ""
		}
		return filepath.Dir(metadata), modfile.ModulePath(data)
	}
	return "", ""
}

type packageInputsGraph struct {
	root, moduleRoot, modulePath string
	source                       InputSpec
	seen                         map[string]bool
	inputs                       []InputSpec
}

func (graph *packageInputsGraph) collect(directory string) bool {
	directory = filepath.Clean(directory)
	if graph.seen[directory] {
		return true
	}
	graph.seen[directory] = true
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false
	}
	input := graph.source
	input.Label, input.Root = "go-package", directory
	input.IncludeNames = nil
	input.IncludeExtensions = append([]string(nil), graph.source.IncludeExtensions...)
	input.ExcludeDirs = append([]string(nil), input.ExcludeDirs...)
	var imports []string
	goFiles := 0
	for _, entry := range entries {
		if entry.IsDir() {
			input.ExcludeDirs = appendUnique(input.ExcludeDirs, entry.Name())
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		goFiles++
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, entry.Name()), nil, parser.ImportsOnly)
		if err != nil {
			return false
		}
		for _, spec := range file.Imports {
			name, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return false
			}
			imports = append(imports, name)
		}
	}
	if goFiles == 0 {
		return false
	}
	graph.inputs = append(graph.inputs, input)
	for _, name := range imports {
		relative, local := strings.CutPrefix(name, graph.modulePath+"/")
		if name == graph.modulePath {
			relative, local = ".", true
		}
		if local && !graph.collect(filepath.Join(graph.moduleRoot, filepath.FromSlash(relative))) {
			return false
		}
	}
	return true
}
