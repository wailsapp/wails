package parser

import (
	"errors"
	"go/ast"
	"go/doc"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/flags"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	pkgs  []*Package
	Stats Stats
}

type Package struct {
	*packages.Package
	services         []*Service
	anonymousStructs map[string]string
	doc              *doc.Package
}

func WrapPackages(pkgs []*packages.Package, services []*Service) []*Package {
	pkgMap := make(map[*types.Package]*Package)

	for _, pkg := range pkgs {
		pkgMap[pkg.Types] = &Package{
			Package:          pkg,
			services:         []*Service{},
			anonymousStructs: make(map[string]string),
			doc:              NewDoc(pkg),
		}
	}

	for _, service := range services {
		pkgMap[service.Pkg()].addService(service)
	}
	return lo.Values(pkgMap)
}

type Service struct {
	*types.TypeName
}

func (s *Service) Methods() (methods []*BoundMethod) {
	if named, ok := s.Type().(*types.Named); ok {
		for i := 0; i < named.NumMethods(); i++ {
			methods = append(methods, &BoundMethod{
				Func: named.Method(i),
				//TODO assign ID
				ID:      0,
				Service: s,
			})
		}
	}
	return
}

type Stats struct {
	NumPackages int
	NumStructs  int
	NumMethods  int
	NumEnums    int
	NumModels   int
	StartTime   time.Time
	EndTime     time.Time
}

func ParseProject(patterns []string, options *flags.GenerateBindingsOptions) (*Project, error) {
	buildFlags, err := options.BuildFlags()
	if err != nil {
		return nil, err
	}

	pkgs, err := LoadPackages(buildFlags, true,
		append(patterns, WailsAppPkgPath)...,
	)
	if err != nil {
		return nil, err
	}

	services, err := Services(pkgs)
	if err != nil {
		return nil, err
	}

	return &Project{
		pkgs: WrapPackages(pkgs, services),
	}, nil
}

func Services(pkgs []*packages.Package) (services []*Service, err error) {
	var app *packages.Package
	otherPkgs := append(make([]*packages.Package, 0, len(pkgs)), pkgs...)
	if index := slices.IndexFunc(pkgs, func(pkg *packages.Package) bool { return pkg.PkgPath == WailsAppPkgPath }); index >= 0 {
		app = pkgs[index]
		otherPkgs = slices.Delete(otherPkgs, index, index+1)
	}

	if app == nil {
		err = errors.New("LoadPackages() did not load the application package")
		return
	}

	found, err := FindServices(app, otherPkgs)
	if err != nil {
		return
	}

	for _, service := range found {
		services = append(services, &Service{service})
	}
	return
}

func (p *Package) anonymousStructID(s *types.Struct) string {
	key := s.String()

	if _, ok := p.anonymousStructs[key]; !ok {
		p.anonymousStructs[key] = "$anon" + strconv.Itoa(len(p.anonymousStructs)+1)
	}
	return p.anonymousStructs[key]
}

func RelativeBindingsDir(base *types.Package, target *types.Package) string {
	if base == target {
		return "."
	}

	basePath := base.Path()
	if base.Name() == "main" {
		basePath = filepath.Join(basePath, "main")
	}

	targetPath := target.Path()
	if target.Name() == "main" {
		targetPath = filepath.Join(targetPath, "main")
	}

	relativePath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		panic(err)
	}

	return filepath.ToSlash(relativePath)
}

// Credit: https://stackoverflow.com/a/70999797/3140799
func (p *Package) constantsOf(t *types.Named) (values map[string]*types.Const) {
	values = make(map[string]*types.Const)

	for _, file := range p.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, specs := range genDecl.Specs {
				valueSpec, ok := specs.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					c := p.TypesInfo.ObjectOf(name).(*types.Const)
					if strings.HasSuffix(c.Type().String(), t.Obj().Name()) {
						values[name.Name] = c
					}
				}
			}
		}
	}
	return
}
