package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

const bindingIDMetadataFile = "wails_obfuscated.gen.go"

// generateBindingIDMetadataAggregate writes the aggregated obfuscated bindings
// file from previously recorded registrations.
func (generator *Generator) generateBindingIDMetadataAggregate() {
	registrations := generator.drainObfuscatedRegistrations()
	if len(registrations) == 0 {
		return
	}

	outputDir, packageName, ok := generator.resolveObfuscatedOutput()
	if !ok {
		return
	}

	// If the destination directory matches a loaded package, services in that
	// same package can be referenced without an import (and must be, because
	// `package main` cannot import itself, and any package's self-import is
	// invalid Go).
	selfPkgPath := generator.dirToPkgPath[outputDir]

	source, err := bindingIDMetadataSource(packageName, selfPkgPath, registrations)
	if err != nil {
		generator.logger.Errorf("obfuscated bindings: %v", err)
		return
	}

	if generator.options.DryRun {
		return
	}

	generator.writeObfuscatedBindingsFile(outputDir, source)
}

func (generator *Generator) drainObfuscatedRegistrations() []bindingIDRegistration {
	generator.obfuscatedMu.Lock()
	registrations := generator.obfuscatedRegistrations
	generator.obfuscatedRegistrations = nil
	generator.obfuscatedMu.Unlock()

	slices.SortFunc(registrations, func(a, b bindingIDRegistration) int {
		if c := strings.Compare(a.PackagePath, b.PackagePath); c != 0 {
			return c
		}
		if c := strings.Compare(a.TypeName, b.TypeName); c != 0 {
			return c
		}
		return strings.Compare(a.MethodName, b.MethodName)
	})
	return registrations
}

// resolveObfuscatedOutput returns the destination dir and package clause for
// the aggregated file. Reports a fatal error and returns ok=false when the
// destination cannot be determined.
func (generator *Generator) resolveObfuscatedOutput() (dir, packageName string, ok bool) {
	dir = generator.options.ObfuscatedOutput
	if dir == "" {
		if generator.mainPackageDirErr != nil {
			generator.logger.Errorf(
				"obfuscated bindings: %v; pass -obfuscated-output to set the destination directory explicitly",
				generator.mainPackageDirErr,
			)
			return "", "", false
		}
		dir = generator.mainPackageDir
	} else {
		// User-routed destination: registration runs only if the chosen
		// package is reachable from main's import graph.
		generator.logger.Infof(
			"obfuscated bindings: writing metadata file to %s; ensure this package is imported (directly or transitively) by your main package so its init() runs",
			dir,
		)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		generator.logger.Errorf("obfuscated bindings: resolve output dir %s: %v", dir, err)
		return "", "", false
	}
	dir = absDir

	packageName, err = resolveTargetPackageName(dir)
	if err != nil {
		fallback := filepath.Base(dir)
		generator.logger.Warningf(
			"obfuscated bindings: %v; using directory name %q as the package clause",
			err, fallback,
		)
		packageName = fallback
	}
	return dir, packageName, true
}

func (generator *Generator) writeObfuscatedBindingsFile(dir string, source []byte) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		generator.logger.Errorf("obfuscated bindings: failed to create %s: %v", dir, err)
		return
	}
	path := filepath.Join(dir, bindingIDMetadataFile)
	if err := os.WriteFile(path, source, 0o644); err != nil {
		generator.logger.Errorf("obfuscated bindings: failed to write %s: %v", path, err)
	}
}

type bindingIDRegistration struct {
	PackagePath string
	PackageName string
	TypeName    string
	MethodName  string
	ID          string
}

type aliasImport struct {
	alias string
	path  string
}

// recordObfuscatedRegistrations appends a package's binding-ID registrations
// to the shared accumulator.
func (generator *Generator) recordObfuscatedRegistrations(index *collect.PackageIndex) {
	if index == nil || len(index.Services) == 0 {
		return
	}

	var local []bindingIDRegistration
	for _, service := range index.Services {
		for _, method := range service.Methods {
			local = append(local, bindingIDRegistration{
				PackagePath: index.Package.Path,
				PackageName: index.Package.Name,
				TypeName:    service.Name,
				MethodName:  method.Name,
				ID:          method.ID,
			})
		}
	}

	if len(local) == 0 {
		return
	}

	generator.obfuscatedMu.Lock()
	generator.obfuscatedRegistrations = append(generator.obfuscatedRegistrations, local...)
	generator.obfuscatedMu.Unlock()
}

// bindingIDMetadataSource renders the obfuscated bindings file. selfPkgPath,
// if non-empty, is the import path of the destination package: registrations
// from that path are emitted without an alias and without an import, since a
// package cannot import itself.
func bindingIDMetadataSource(packageName, selfPkgPath string, registrations []bindingIDRegistration) ([]byte, error) {
	if len(registrations) == 0 {
		return nil, errors.New("no registrations to emit")
	}
	if packageName == "" {
		return nil, errors.New("empty package name for obfuscated bindings file")
	}

	aliases := buildPackageAliases(registrations, selfPkgPath)

	var buf bytes.Buffer
	writeMetadataHeader(&buf, packageName)
	writeMetadataImports(&buf, aliases)
	writeMetadataInit(&buf, registrations, aliases, selfPkgPath)

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format.Source: %w", err)
	}
	return source, nil
}

func writeMetadataHeader(buf *bytes.Buffer, packageName string) {
	fmt.Fprintf(buf, "//go:build wails_obfuscated\n")
	fmt.Fprintf(buf, "// +build wails_obfuscated\n\n")
	fmt.Fprintf(buf, "// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL\n")
	fmt.Fprintf(buf, "// This file is automatically generated. DO NOT EDIT\n\n")
	fmt.Fprintf(buf, "package %s\n\n", packageName)
}

func writeMetadataImports(buf *bytes.Buffer, aliases map[string]string) {
	imports := make([]aliasImport, 0, len(aliases))
	for path, alias := range aliases {
		imports = append(imports, aliasImport{alias: alias, path: path})
	}
	slices.SortFunc(imports, func(a, b aliasImport) int { return strings.Compare(a.alias, b.alias) })

	fmt.Fprintf(buf, "import (\n")
	fmt.Fprintf(buf, "\t\"github.com/wailsapp/wails/v3/pkg/application\"\n\n")
	for _, imp := range imports {
		fmt.Fprintf(buf, "\t%s %q\n", imp.alias, imp.path)
	}
	fmt.Fprintf(buf, ")\n\n")
}

func writeMetadataInit(buf *bytes.Buffer, registrations []bindingIDRegistration, aliases map[string]string, selfPkgPath string) {
	fmt.Fprintf(buf, "func init() {\n")
	for _, reg := range registrations {
		if selfPkgPath != "" && reg.PackagePath == selfPkgPath {
			fmt.Fprintf(
				buf,
				"\tapplication.RegisterBindingMethodID((*%s).%s, %s)\n",
				reg.TypeName,
				reg.MethodName,
				reg.ID,
			)
			continue
		}
		fmt.Fprintf(
			buf,
			"\tapplication.RegisterBindingMethodID((*%s.%s).%s, %s)\n",
			aliases[reg.PackagePath],
			reg.TypeName,
			reg.MethodName,
			reg.ID,
		)
	}
	fmt.Fprintf(buf, "}\n")
}

// buildPackageAliases returns a deterministic alias per import path.
// Collisions on the seed name are resolved with a numeric suffix.
// selfPkgPath, if non-empty, is excluded: the destination package is referred
// to without an import.
func buildPackageAliases(registrations []bindingIDRegistration, selfPkgPath string) map[string]string {
	aliases := make(map[string]string)
	used := map[string]bool{
		"application": true,
	}

	seen := make(map[string]string)
	var paths []string
	for _, reg := range registrations {
		if reg.PackagePath == selfPkgPath {
			continue
		}
		if _, ok := seen[reg.PackagePath]; ok {
			continue
		}
		seen[reg.PackagePath] = reg.PackageName
		paths = append(paths, reg.PackagePath)
	}
	slices.Sort(paths)

	for _, path := range paths {
		base := sanitizeAlias(seen[path])
		if base == "" {
			base = "pkg"
		}
		alias := base
		for i := 2; used[alias]; i++ {
			alias = fmt.Sprintf("%s%d", base, i)
		}
		aliases[path] = alias
		used[alias] = true
	}

	return aliases
}

func sanitizeAlias(name string) string {
	var b strings.Builder
	first := true
	for _, r := range name {
		valid := r == '_' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(!first && r >= '0' && r <= '9')
		if valid {
			b.WriteRune(r)
			first = false
		}
	}
	return b.String()
}

// resolveTargetPackageName returns the package clause read from a .go file in
// dir, or an error when none can be determined.
func resolveTargetPackageName(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("target directory %s does not yet exist", dir)
		}
		return "", fmt.Errorf("read dir %s: %w", dir, err)
	}

	fset := token.NewFileSet()
	sawGoFile := false
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		// Match `go build`'s ignore rules.
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		// Don't read our own output as the truth source on a re-run.
		if name == bindingIDMetadataFile {
			continue
		}
		sawGoFile = true
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
		if err != nil {
			continue
		}
		if file.Name != nil && file.Name.Name != "" {
			return file.Name.Name, nil
		}
	}

	if !sawGoFile {
		return "", fmt.Errorf("no Go source files found in %s", dir)
	}
	return "", fmt.Errorf("no parseable Go source files in %s", dir)
}
