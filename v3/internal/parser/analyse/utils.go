package analyse

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// Set represents a set of elements of type T by mapping them to booleans.
type Set[T comparable] map[T]bool

// Add adds an element to a set.
// It returns true if the element was added to the set,
// false if it was already present.
func (set Set[T]) Add(element T) (added bool) {
	added = !set[element]
	set[element] = true
	return
}

// Union adds all elements from the argument to the receiver
// and returns the number of newly added elements.
func (set Set[T]) Union(other Set[T]) (added int) {
	for el, _ := range other {
		if set.Add(el) {
			added++
		}
	}
	return
}

// Elements returns a slice containing all elements in the given set.
func (set Set[T]) Elements() []T {
	return lo.Keys(set)
}

// IsStdImportPath returns true if path has the shape of a Go standard import path,
// i.e. the first segment of the path contains no dots.
func IsStdImportPath(path string) bool {
	// TODO: verify that this std detection rule is always sound.
	top, _, _ := strings.Cut(path, "/")
	return !strings.ContainsRune(top, '.')
}

// SortAstFiles sorts the slice of files held by pkg.Syntax by starting position.
func SortAstFiles(pkg *packages.Package) {
	slices.SortFunc(pkg.Syntax, func(f1 *ast.File, f2 *ast.File) int {
		return cmp.Compare(f1.FileStart, f2.FileStart)
	})
}

// FindAstPath returns the node that encloses the source interval [start, end),
// and all its ancestors up to the [ast.File], in the syntax tree for pkg.
// If no source file can be found in pkg.Syntax for the specified interval,
// FindAstPath returns an empty slice.
//
// The slice of files held by pkg.Syntax must be sorted by starting position.
//
// See also [astutil.PathEnclosingInterval], [SortAstFiles].
func FindAstPath(pkg *packages.Package, start token.Pos, end token.Pos) []ast.Node {
	// Perform a binary search to find the file enclosing the node
	fileIndex, exact := slices.BinarySearchFunc(pkg.Syntax, start, func(f *ast.File, p token.Pos) int {
		return cmp.Compare(f.FileStart, p)
	})

	// If exact is true, pkg.Syntax[fileIndex] is the file we are looking for;
	// otherwise, it is the first file whose start position is _after_ ident.Pos()
	if !exact {
		fileIndex--
	}

	// When exact is false, the search could theoretically fail (this is bad and should never happen)
	if fileIndex < 0 || start < pkg.Syntax[fileIndex].FileStart || pkg.Syntax[fileIndex].FileEnd < end {
		return nil
	}

	path, _ := astutil.PathEnclosingInterval(pkg.Syntax[fileIndex], start, end)
	return path
}

// Reparen is the opposite of [ast.Unparen]: it travels up the given AST path
// until the immediate context `path[1]` is not a parenthesized expression.
func Reparen(path []ast.Node) []ast.Node {
	for ; len(path) > 1; path = path[1:] {
		if _, ok := path[1].(*ast.ParenExpr); !ok {
			break
		}
	}

	return path
}

// LongestCommonPrefix computes the length
// of the longest common prefix of two slices of integers.
func LongestCommonPrefix[S ~[]E, E comparable](p1 S, p2 S) (length int) {
	for length = 0; length < len(p1) && length < len(p2); length++ {
		if p1[length] != p2[length] {
			break
		}
	}

	return
}

// IsValidType returns true if and only if the given type is non-nil and valid.
func IsValidType(T types.Type) bool {
	return T != nil && !types.Identical(T, types.Typ[types.Invalid])
}

// IsAssignableTo returns true if and only if V, T are both valid types
// and values of type V are assignable to variables of type T.
//
// This is a custom version of [types.AssignableTo] that avoids
// undefined behaviour in case V, T are either nil or invalid.
func IsAssignableTo(V types.Type, T types.Type) bool {
	return IsValidType(V) && IsValidType(T) && types.AssignableTo(V, T)
}

// IsConvertibleTo returns true if and only if V, T are both valid types
// and values of type V are convertible to values of type T.
//
// This is a custom version of [types.ConvertibleTo] that avoids
// undefined behaviour in case V, T are either nil or invalid.
func IsConvertibleTo(V types.Type, T types.Type) bool {
	return IsValidType(V) && IsValidType(T) && types.ConvertibleTo(V, T)
}

// varPkgi computes the package index with which the given variable
// should be scheduled.
func varPkgi(analyser *Analyser, pkgi int, variable *types.Var) int {
	if variable.Parent() == variable.Pkg().Scope() && variable.Exported() {
		// Global exported variables must be analysed in every package.
		return -1
	}

	// Local or unexported variable: test current package first,
	// then search package slice.
	if variable.Pkg() != analyser.pkgs[pkgi].Types {
		for pkgi = 0; pkgi < len(analyser.pkgs); pkgi++ {
			if variable.Pkg() == analyser.pkgs[pkgi].Types {
				break
			}
		}
	}

	return pkgi
}
