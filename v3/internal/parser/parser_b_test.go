package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func BenchmarkParser(b *testing.B) {
	benchmarks := []struct {
		pkg string
	}{
		{
			pkg: "struct_literal_single",
		},
		{
			pkg: "complex_json",
		},
		// {
		// 	pkg: "multiple_packages",
		// },
	}

	for _, bench := range benchmarks {

		options := &flags.GenerateBindingsOptions{
			ModelsFilename:   "models",
			TS:               true,
			ProjectDirectory: "github.com/wailsapp/wails/v3/internal/parser/testdata/" + bench.pkg,
		}

		b.Run(bench.pkg+"/load", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LoadPackages(nil, true,
					options.ProjectDirectory,
					WailsAppPkgPath,
				)
			}

		})

		pPkgs, err := LoadPackages(nil, true,
			options.ProjectDirectory, WailsAppPkgPath,
		)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/analyze", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseServices(pPkgs[1], pPkgs[0])
			}
		})

		b.Run(bench.pkg+"/parse", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseProject(options)
			}
		})

		project, err := ParseProject(options)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/bindings", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				project.GenerateBindings()
			}
		})

		b.Run(bench.pkg+"/models", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				project.GenerateModels()
			}
		})
	}
}

func BenchmarkLoad(b *testing.B) {

	benchmarks := []struct {
		pkg string
	}{
		{
			pkg: "struct_literal_single",
		},
		{
			pkg: "complex_json",
		},
	}

	for _, bench := range benchmarks {

		dir := "testdata/" + bench.pkg
		absDir, err := filepath.Abs(dir)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/packages.Load", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pkgs, err := LoadPackages(nil, true, absDir)
				if err != nil {
					b.Fatal(err)
				}
				fmt.Println(len(pkgs[0].Syntax))
			}
		})

		b.Run(bench.pkg+"/parser.ParseFile", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fset := token.NewFileSet()
				_, err := parser.ParseDir(fset, absDir, nil, parser.AllErrors|parser.ParseComments|parser.SkipObjectResolution)
				if err != nil {
					b.Fatal(err)
				}
			}

		})

	}

}
