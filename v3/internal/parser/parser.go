package parser

import (
	"errors"
	"go/ast"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	pkgs             []*packages.Package
	options          *flags.GenerateBindingsOptions
	services         []*Service
	anonymousStructs map[string]string
	Stats            Stats
}

type Service struct {
	*types.TypeName
	Parent *Project
}

func (s *Service) Methods() (methods []*BoundMethod) {
	if named, ok := s.Type().(*types.Named); ok {
		for i := 0; i < named.NumMethods(); i++ {
			methods = append(methods, &BoundMethod{
				Func: named.Method(i),
				//TODO assign ID
				ID:     0,
				Parent: s,
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

	return &Project{
		pkgs:             pkgs,
		options:          options,
		anonymousStructs: make(map[string]string),
	}, nil
}

func (p *Project) Services() (services []*Service, err error) {
	if p.services != nil {
		return p.services, nil
	}

	var app *packages.Package
	pkgs := append(make([]*packages.Package, 0, len(p.pkgs)), p.pkgs...)
	if index := slices.IndexFunc(p.pkgs, func(pkg *packages.Package) bool { return pkg.PkgPath == WailsAppPkgPath }); index >= 0 {
		app = p.pkgs[index]
		p.pkgs = slices.Delete(p.pkgs, index, index+1)
	}

	if app == nil {
		err = errors.New("LoadPackages() did not load the application package")
		return
	}

	found, err := FindServices(app, pkgs)
	if err != nil {
		return
	}

	for _, service := range found {
		services = append(services, &Service{service, p})
	}
	p.services = services
	return
}

func (p *Project) anonymousStructID(s *types.Struct) string {
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
func (p *Project) constantsOf(t *types.Named) (values map[string]*types.Const) {
	values = make(map[string]*types.Const)

	pkgIndex := slices.IndexFunc(p.pkgs, func(pkg *packages.Package) bool { return pkg.Types == t.Obj().Pkg() })
	if pkgIndex < 0 {
		return
	}
	pkg := p.pkgs[pkgIndex]

	for _, file := range pkg.Syntax {
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
					c := pkg.TypesInfo.ObjectOf(name).(*types.Const)
					if strings.HasSuffix(c.Type().String(), t.Obj().Name()) {
						values[name.Name] = c
					}
				}
			}
		}
	}
	return
}
