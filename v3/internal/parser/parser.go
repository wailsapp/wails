package parser

import (
	"errors"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	pkgs             []*packages.Package
	options          *flags.GenerateBindingsOptions
	BoundMethods     []*BoundMethod
	anonymousStructs map[string]string
	Stats            Stats
}

type Service struct {
	*types.TypeName
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

func (p *Project) Services() (services []*types.TypeName, err error) {
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

	services, err = FindServices(app, pkgs)
	return
}

func (p *Project) anonymousStructID(s *types.Struct) string {
	key := s.String()

	if _, ok := p.anonymousStructs[key]; !ok {
		p.anonymousStructs[key] = "$anon" + strconv.Itoa(len(p.anonymousStructs)+1)
	}
	return p.anonymousStructs[key]
}

func (p *Project) RelativeBindingsDir(base *types.Package, target *types.Package) string {
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
