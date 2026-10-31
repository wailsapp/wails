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

// Keep local replacement/workspace modules conservatively tracked as before.
// Within the project module, each executable tracks its imported Go packages,
// including inactive platform files, plus the existing broad native-file input.
// Unknown or incomplete source graphs retain the original conservative inputs.
func refineCompileInputs(plan *Plan, root string) {
	for key, node := range plan.Nodes {
		if spec, ok := node.Spec.(CompileSpec); ok {
			node.Inputs = executableInputs(root, spec.Package, node.Inputs)
			plan.Nodes[key] = node
		}
	}
}

func executableInputs(root, entry string, original []InputSpec) []InputSpec {
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
		return original
	}
	moduleRoot, modulePath := projectModule(root)
	if modulePath == "" {
		return original
	}
	directory := filepath.Join(root, entry)
	graph := packageInputsGraph{root: root, moduleRoot: moduleRoot, modulePath: modulePath, source: source, seen: map[string]bool{}}
	if !graph.collect(directory) {
		return original
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
	return result
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
	for _, entry := range entries {
		if entry.IsDir() {
			input.ExcludeDirs = appendUnique(input.ExcludeDirs, entry.Name())
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
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
