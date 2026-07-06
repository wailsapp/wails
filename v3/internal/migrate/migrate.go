// Package migrate contains the engine behind `wails3 migrate`: it parses a
// Wails v2 project (wails.json + the declarative options.App literal passed to
// wails.Run) and produces the pieces of an equivalent v3 project. The
// orchestration (scaffolding via the init machinery, build assets, etc.) lives
// in internal/commands/migrate.go; everything in this package is pure
// parse/transform logic so it can be unit tested in isolation.
package migrate

import (
	"go/ast"
	"go/token"
)

// V2RuntimeImport is the Wails v2 runtime package rewritten by the migrator.
// It is replaced with the project-local compatibility bridge (see
// CompatRuntimeImport and WriteCompatBridge); the generated package is also
// named "runtime", so call sites compile unchanged.
const V2RuntimeImport = "github.com/wailsapp/wails/v2/pkg/runtime"

// V2Project is everything the migrator learned about the source project.
type V2Project struct {
	Dir    string
	Config V2Config // parsed wails.json

	ModulePath  string // from go.mod
	GoModPath   string // absolute path to go.mod
	FrontendDir string // absolute path to the frontend directory

	Main *MainInfo // parsed main package (wails.Run call site)

	// BoundTypes are the resolved struct types listed in Bind, in order.
	BoundTypes []*BoundType

	// GoFiles are all .go files of the project (absolute paths), excluding
	// the file containing wails.Run (handled separately).
	GoFiles []string

	// UsesV2Runtime is true when any project file imports the v2 runtime
	// package (the migrated project then needs the compatibility bridge).
	UsesV2Runtime bool

	// Report accumulates human-readable notes about everything that needs
	// manual attention. It is written to MIGRATION.md in the output project.
	Report *Report
}

// MainInfo describes the file containing the wails.Run call.
type MainInfo struct {
	Path   string // absolute path
	Source []byte // original file bytes
	File   *ast.File
	Fset   *token.FileSet

	// RunStmt is the statement containing the wails.Run(...) call.
	RunStmt ast.Stmt
	// RunCall is the wails.Run call expression itself.
	RunCall *ast.CallExpr
	// AppLit is the &options.App{...} composite literal passed to wails.Run,
	// nil if the argument was not a literal.
	AppLit *ast.CompositeLit
	// ErrIdent is the identifier the wails.Run error was assigned to
	// ("err" in `err := wails.Run(...)`), empty for a bare call.
	ErrIdent string
	// AssignTok is the assignment token used (token.DEFINE or token.ASSIGN)
	// when ErrIdent is set.
	AssignTok token.Token

	// Imports maps local package names to import paths for this file.
	Imports map[string]string
}

// BoundType is a struct type bound to the frontend via options.App.Bind.
type BoundType struct {
	// Expr is the original Go expression used in the Bind slice (e.g. "app").
	Expr string
	// PkgName is the Go package name the type is declared in (e.g. "main").
	PkgName string
	// PkgPath is the full package path used in binding FQNs. For types in the
	// main package this is "main"; for other in-module packages it is
	// module/path/to/pkg.
	PkgPath string
	// Name is the struct type name (e.g. "App").
	Name string
	// Methods are the exported methods, in declaration order.
	Methods []*BoundMethod
}

// BoundMethod is one exported method of a bound type.
type BoundMethod struct {
	Name    string
	Params  []Param // excluding the receiver
	Results []Param
}

// Param is a parameter or result of a bound method.
type Param struct {
	Name string
	// GoType is the printed Go type (e.g. "string", "[]int", "*Person").
	GoType string
	// TSType is the best-effort TypeScript equivalent used in generated .d.ts
	// shims (e.g. "string", "number", "any").
	TSType string
}
