// Command puregogen generates purego function bindings from Clang's C AST.
// Only maintainers need Clang and development headers; application builds use
// the checked-in Go output and dynamically load the runtime libraries.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type function struct {
	Name         string `json:"name"`
	Library      string `json:"library"`
	GoName       string `json:"go_name,omitempty"`
	Optional     bool   `json:"optional,omitempty"`
	Strings      []int  `json:"strings,omitempty"`
	StringResult bool   `json:"string_result,omitempty"`
}

type manifest struct {
	Package   string     `json:"package"`
	Build     string     `json:"build"`
	Headers   []string   `json:"headers"`
	PkgConfig []string   `json:"pkg_config"`
	Functions []function `json:"functions"`
}

type ctype struct {
	QualType  string `json:"qualType"`
	Desugared string `json:"desugaredQualType"`
}

type node struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Type     ctype  `json:"type"`
	Variadic bool   `json:"variadic"`
	Inner    []node `json:"inner"`
}

type declarations struct {
	aliases   map[string]string
	functions map[string]node
	enums     map[string]bool
}

func main() {
	manifestPath := flag.String("manifest", "", "binding manifest")
	output := flag.String("out", "", "generated Go file")
	clang := flag.String("clang", "clang", "Clang executable")
	check := flag.Bool("check", false, "fail if the checked-in output differs")
	cgoManifest := flag.String("cgo-manifest", "", "mechanically port existing CGo source using compiler-derived ABI declarations")
	flag.Parse()
	var err error
	if *cgoManifest != "" {
		err = portCgo(*cgoManifest, *check)
	} else {
		err = run(*manifestPath, *output, *clang, *check)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "puregogen:", err)
		os.Exit(1)
	}
}

func run(manifestPath, output, clang string, check bool) error {
	if manifestPath == "" || output == "" {
		return fmt.Errorf("-manifest and -out are required")
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var config manifest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return err
	}
	args := []string{"-x", "c", "-fsyntax-only", "-Xclang", "-ast-dump=json", "-"}
	if len(config.PkgConfig) > 0 {
		flags, err := exec.Command("pkg-config", append([]string{"--cflags"}, config.PkgConfig...)...).Output()
		if err != nil {
			return fmt.Errorf("pkg-config: %w", err)
		}
		args = append(strings.Fields(string(flags)), args...)
	}
	var header strings.Builder
	for _, name := range config.Headers {
		fmt.Fprintf(&header, "#include <%s>\n", name)
	}
	// The generated ABI uses LP64. Fail at generation instead of silently
	// producing incorrect long/size_t declarations on a different target.
	header.WriteString("_Static_assert(sizeof(void*) == 8 && sizeof(long) == 8, \"puregogen requires LP64\");\n")
	cmd := exec.Command(clang, args...)
	cmd.Stdin = strings.NewReader(header.String())
	var diagnostics bytes.Buffer
	cmd.Stderr = &diagnostics
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	d, parseErr := readDeclarations(stdout)
	if parseErr != nil {
		_, _ = io.Copy(io.Discard, stdout)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("clang: %w\n%s", err, &diagnostics)
	}
	if parseErr != nil {
		return parseErr
	}
	generated, err := generate(config, d)
	if err != nil {
		return err
	}
	if check {
		old, err := os.ReadFile(output)
		if err != nil {
			return err
		}
		if !bytes.Equal(old, generated) {
			return fmt.Errorf("%s is stale; regenerate bindings", output)
		}
		return nil
	}
	return os.WriteFile(output, generated, 0644)
}

func readDeclarations(reader io.Reader) (declarations, error) {
	var root node
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return declarations{}, err
	}
	d := declarations{aliases: map[string]string{}, functions: map[string]node{}, enums: map[string]bool{}}
	var visit func(node)
	visit = func(n node) {
		switch n.Kind {
		case "TypedefDecl":
			d.aliases[n.Name] = n.Type.QualType
		case "EnumDecl":
			d.enums[n.Name] = true
		case "FunctionDecl":
			d.functions[n.Name] = n
		}
		for _, child := range n.Inner {
			visit(child)
		}
	}
	visit(root)
	return d, nil
}

func (d declarations) resolve(name string) string {
	seen := map[string]bool{}
	for !seen[name] {
		seen[name] = true
		next, ok := d.aliases[name]
		if !ok {
			break
		}
		name = next
	}
	return name
}

func (d declarations) goType(name string, asString bool) (string, error) {
	name = d.resolve(name)
	if asString {
		base := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(name, "*"), "const "))
		if !strings.HasSuffix(name, "*") || d.resolve(base) != "char" {
			return "", fmt.Errorf("string conversion requires char *, got %q", name)
		}
		return "string", nil
	}
	// Includes opaque pointers, pointer-to-pointer parameters, and callback
	// typedefs. Ownership and callback lifetime remain explicit in the caller.
	if strings.Contains(name, "*") {
		return "uintptr", nil
	}
	if strings.HasPrefix(name, "enum ") {
		return "int32", nil
	}
	if d.enums[name] {
		return "int32", nil
	}
	switch name {
	case "void":
		return "", nil
	case "_Bool", "bool":
		return "bool", nil
	case "char", "signed char":
		return "int8", nil
	case "unsigned char":
		return "uint8", nil
	case "short", "short int":
		return "int16", nil
	case "unsigned short", "unsigned short int":
		return "uint16", nil
	case "int":
		return "int32", nil
	case "unsigned int":
		return "uint32", nil
	case "long", "long int", "long long", "long long int":
		return "int64", nil
	case "unsigned long", "unsigned long int", "unsigned long long", "unsigned long long int":
		return "uint64", nil
	case "float":
		return "float32", nil
	case "double":
		return "float64", nil
	}
	return "", fmt.Errorf("unsupported C type %q (use an explicit wrapper for aggregates or variadics)", name)
}

func generate(config manifest, d declarations) ([]byte, error) {
	var out bytes.Buffer
	var wrappers bytes.Buffer
	fmt.Fprintln(&out, "// Code generated by internal/puregogen from C headers; DO NOT EDIT.")
	if config.Build != "" {
		fmt.Fprintf(&out, "//go:build %s\n", config.Build)
	}
	fmt.Fprintf(&out, "\npackage %s\n\nvar (\n", config.Package)
	functions := append([]function(nil), config.Functions...)
	sort.Slice(functions, func(i, j int) bool { return functions[i].Name < functions[j].Name })
	seen := map[string]bool{}
	for _, f := range functions {
		name := f.GoName
		if name == "" {
			name = f.Name
		}
		if seen[name] {
			return nil, fmt.Errorf("duplicate binding %s", name)
		}
		seen[name] = true
		n, ok := d.functions[f.Name]
		if !ok {
			return nil, fmt.Errorf("%s: declaration missing from headers", f.Name)
		}
		if n.Variadic {
			return nil, fmt.Errorf("%s: variadic functions require a fixed-signature wrapper", f.Name)
		}
		var params []node
		for _, child := range n.Inner {
			if child.Kind == "ParmVarDecl" {
				params = append(params, child)
			}
		}
		stringParams := map[int]bool{}
		for _, index := range f.Strings {
			if index < 0 || index >= len(params) {
				return nil, fmt.Errorf("%s: invalid string parameter %d", f.Name, index)
			}
			stringParams[index] = true
		}
		var args []string
		for i, param := range params {
			typ, err := d.goType(param.Type.QualType, stringParams[i])
			if err != nil {
				return nil, fmt.Errorf("%s parameter %d: %w", f.Name, i, err)
			}
			args = append(args, typ)
		}
		result, _, ok := strings.Cut(n.Type.QualType, "(")
		if !ok {
			return nil, fmt.Errorf("%s: unrecognised prototype %q", f.Name, n.Type.QualType)
		}
		ret, err := d.goType(strings.TrimSpace(result), f.StringResult)
		if err != nil {
			return nil, fmt.Errorf("%s return: %w", f.Name, err)
		}
		fmt.Fprintf(&out, "\t_native_%s func(%s) %s\n", name, strings.Join(args, ", "), ret)
		var paramsGo, values []string
		for i, typ := range args {
			param := fmt.Sprintf("p%d", i)
			paramsGo = append(paramsGo, param+" "+typ)
			values = append(values, param)
		}
		// RegisterFunc's function variable cannot carry compiler directives.
		// A wrapper makes uintptr(unsafe.Pointer(&goMemory)) escape and remain
		// alive across stack growth and callbacks during the foreign call.
		fmt.Fprintf(&wrappers, "//go:uintptrescapes\nfunc %s(%s) %s {\n", name, strings.Join(paramsGo, ", "), ret)
		if ret != "" {
			wrappers.WriteString("return ")
		}
		fmt.Fprintf(&wrappers, "_native_%s(%s)\n}\n\n", name, strings.Join(values, ", "))
	}
	fmt.Fprintln(&out, ")\n\nfunc registerGeneratedNativeFunctions() {")
	for _, f := range functions {
		register := "register"
		if f.Optional {
			register = "registerOptional"
		}
		name := f.GoName
		if name == "" {
			name = f.Name
		}
		fmt.Fprintf(&out, "\t%s(&_native_%s, %s, %q)\n", register, name, f.Library, f.Name)
	}
	fmt.Fprintln(&out, "}")
	out.Write(wrappers.Bytes())
	return format.Source(out.Bytes())
}
