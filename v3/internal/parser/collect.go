package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/packages"
)

type (
	PackageInfo struct {
		Path   string
		Name   string
		Groups []*GroupInfo
		Types  map[string]*TypeInfo
		Consts map[string]*ConstInfo

		pkg  *packages.Package
		err  error
		once sync.Once
	}

	GroupInfo struct {
		Doc   string
		Group *GroupInfo
	}

	TypeInfo struct {
		Name    string
		Doc     string
		Group   *GroupInfo
		Alias   bool
		Def     ast.Expr
		Methods map[string]*MethodInfo
	}

	MethodInfo struct {
		Name string
		Doc  string
	}

	ConstInfo struct {
		Name  string
		Doc   string
		Group *GroupInfo
	}
)

// NewPackageInfo initialises an empty information structure
// for the given package path.
//
// If pkg is not nil, the Collect method will use it
// to get syntax information instead of loading the package
// from disk.
//
// The cost of this function must be as low as possible
// and it must not perform any significant work,
// as it might be called multiple times for the same package
// and its result might be discarded often.
func NewPackageInfo(path string, pkg *packages.Package) *PackageInfo {
	return &PackageInfo{
		Path: path,
		pkg:  pkg,
	}
}

// Collect loads information for the package described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// If the package has not been loaded yet,
// it will be loaded with the specified build flags.
//
// After Collect returns, the calling goroutine is free to access
// the receiver's fields indefinitely.
func (pi *PackageInfo) Collect(buildFlags []string) error {
	pi.once.Do(func() {
		pkg := pi.pkg
		pi.pkg = nil

		if pkg == nil {
			pkgs, err := LoadPackages(buildFlags, false, pi.Path)
			if err != nil {
				pi.err = err
				return
			} else if len(pkgs) < 1 {
				pi.err = fmt.Errorf("%s: package not found", pi.Path)
				return
			} else if len(pkgs) > 1 {
				pi.err = fmt.Errorf("%s: multiple packages loaded for the same path", pi.Path)
				return
			}

			pkg = pkgs[0]
		}

		pi.Types = make(map[string]*TypeInfo)
		pi.Consts = make(map[string]*ConstInfo)

		// Collect all methods here temporarily.
		methods := make(map[string]map[string]*MethodInfo)

		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					var group *GroupInfo
					empty := true

					doc := decl.Doc.Text()
					if len(doc) != 0 {
						group = &GroupInfo{
							Doc:   doc,
							Group: nil,
						}
					}

					switch decl.Tok {
					case token.TYPE:
						for _, spec := range decl.Specs {
							tspec, ok := spec.(*ast.TypeSpec)
							if !ok {
								// Ignore invalid or malformed specs.
								continue
							}

							empty = false
							pi.Types[tspec.Name.Name] = &TypeInfo{
								Name:  tspec.Name.Name,
								Doc:   tspec.Doc.Text(),
								Group: group,
								Alias: tspec.Assign.IsValid(),
								Def:   tspec.Type,
							}
						}

					case token.CONST:
						for _, spec := range decl.Specs {
							vspec, ok := spec.(*ast.ValueSpec)
							if !ok || len(vspec.Names) == 0 {
								// Ignore invalid or malformed specs.
								continue
							}

							sgroup := group
							sempty := true

							doc = vspec.Doc.Text()
							if len(doc) != 0 && len(vspec.Names) > 1 {
								sgroup = &GroupInfo{
									Doc:   doc,
									Group: group,
								}
								doc = ""
							}

							for _, name := range vspec.Names {
								if !name.IsExported() {
									// Ignore unexported constants.
									return
								}

								empty = false
								sempty = false
								pi.Consts[name.Name] = &ConstInfo{
									Name:  name.Name,
									Doc:   doc,
									Group: sgroup,
								}
							}

							if !sempty && sgroup != group {
								pi.Groups = append(pi.Groups, sgroup)
							}
						}
					}

					if !empty && group != nil {
						pi.Groups = append(pi.Groups, group)
					}

				case *ast.FuncDecl:
					if decl.Recv == nil || !decl.Name.IsExported() {
						// Ignore functions and unexported methods.
						continue
					}

					if len(decl.Recv.List) != 1 || len(decl.Recv.List[0].Names) > 1 {
						// Malformed receiver.
						continue
					}

					var recv string

					switch expr := decl.Recv.List[0].Type.(type) {
					case *ast.Ident:
						recv = expr.Name
					case *ast.StarExpr:
						ident, ok := expr.X.(*ast.Ident)
						if !ok {
							// Invalid receiver type
							continue
						}

						recv = ident.Name
					default:
						// Invalid receiver type.
						continue
					}

					mmap := methods[recv]
					if mmap == nil {
						mmap = make(map[string]*MethodInfo)
						methods[recv] = mmap
					}

					mmap[decl.Name.Name] = &MethodInfo{
						Name: decl.Name.Name,
						Doc:  decl.Doc.Text(),
					}
				}
			}
		}

		// Assign method maps to discovered types.
		for name, info := range pi.Types {
			info.Methods = methods[name]
		}
	})
	return pi.err
}
