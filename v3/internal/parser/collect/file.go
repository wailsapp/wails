package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"sync"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/parser/config"
	"golang.org/x/tools/go/packages"
)

// FileInfo records information about an [ast.File]
// and provides infrastructure for qualified name resolution
// within its scope.
type FileInfo struct {
	fset *token.FileSet
	file *ast.File

	scope  *types.Scope
	logger config.Logger
	once   sync.Once
}

// newFileInfo initialises an info struct for the given package and file.
// The logger will be used to log errors that may occur while type-checking
// the file's import declarations.
func newFileInfo(pkg *packages.Package, file *ast.File, logger config.Logger) *FileInfo {
	// We only care about imports hence we discard everything else.
	info := &FileInfo{
		fset: pkg.Fset,
		file: &ast.File{
			Package: file.Package,
			Name:    file.Name,
			Decls: lo.Filter(file.Decls, func(decl ast.Decl, _ int) bool {
				gen, _ := decl.(*ast.GenDecl)
				return gen != nil && gen.Tok == token.IMPORT
			}),
			FileStart: file.FileStart,
			FileEnd:   file.FileEnd,
			Imports:   file.Imports,
			GoVersion: file.GoVersion,
		},
		logger: logger,
	}

	// If file scope is available, cache it.
	if pkg.TypesInfo != nil {
		info.scope = pkg.TypesInfo.Scopes[file]
	}

	return info
}

// TypeOf computes the type of any expression
// within the scope of the given file and package.
//
// TypeOf returns nil for invalid expressions.
//
// The behaviour of TypeOf is unspecified if the file
// described by the receiver is not part of the given package.
//
// Resolve is safe for concurrent use.
func (info *FileInfo) TypeOf(pkg *types.Package, expr ast.Expr) types.Type {
	// Sanity check.
	if pkg.Name() != info.file.Name.Name {
		return nil
	}

	info.Scope(pkg)

	typesInfo := types.Info{
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	err := types.CheckExpr(info.fset, pkg, expr.Pos(), expr, &typesInfo)
	if err != nil {
		return nil
	}

	return typesInfo.TypeOf(expr)
}

// Scope returns the scope of the file described by info
// relative to the given package. The returned scope is never nil.
//
// In some instances, the returned scope
// may contain just imports and nothing else.
// Therefore, it is only useful in resolving import names.
//
// Scope is safe for concurrent calls.
func (info *FileInfo) Scope(pkg *types.Package) *types.Scope {
	info.once.Do(func() {
		logger := info.logger
		info.logger = nil

		if info.scope != nil {
			return
		}

		// Setup type checker.
		checker := types.NewChecker(
			&types.Config{
				// Report type-checking errors as warnings.
				Error:    func(err error) { logger.Warningf("%v", err) },
				Importer: newImporter(pkg),

				IgnoreFuncBodies:         true,
				DisableUnusedImportCheck: true,
			},
			info.fset,
			pkg,
			&types.Info{
				Scopes: make(map[ast.Node]*types.Scope),
			},
		)

		if err := checker.Files([]*ast.File{info.file}); err != nil {
			// Report error.
			logger.Warningf("%v", err)
		}

		info.scope = checker.Scopes[info.file]
		if info.scope == nil {
			// In case of errors, return an empty scope.
			info.scope = types.NewScope(pkg.Scope(), info.file.FileStart, info.file.FileEnd, "")
		}
	})

	return info.scope
}

type importer map[string]*types.Package

func newImporter(pkg *types.Package) importer {
	result := make(importer)
	for _, imp := range pkg.Imports() {
		result[imp.Path()] = imp
	}
	return result
}

func (i importer) Import(path string) (*types.Package, error) {
	pkg, ok := i[path]
	if !ok {
		// Create a temporary fake package; this will never be involved in lookups anyways.
		pkg = types.NewPackage(path, path)
		pkg.MarkComplete()
	}
	return pkg, nil
}
