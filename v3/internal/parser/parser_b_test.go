package parser

import (
	"testing"

	"github.com/pterm/pterm"
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
		{
			pkg: "multiple_packages",
		},
	}

	// suppress warnings
	pterm.Warning.Debugger = true

	for _, bench := range benchmarks {

		options := &flags.GenerateBindingsOptions{
			ModelsFilename:   "models",
			TS:               true,
			ProjectDirectory: "github.com/wailsapp/wails/v3/internal/parser/testdata/" + bench.pkg,
		}

		// b.Run(bench.pkg+"/LoadPackages", func(b *testing.B) {
		// 	buildFlags, err := options.BuildFlags()
		// 	if err != nil {
		// 		b.Fatal(err)
		// 	}

		// 	for i := 0; i < b.N; i++ {
		// 		_, err := LoadPackages(buildFlags, true, options.ProjectDirectory, WailsAppPkgPath, JsonPkgPath)
		// 		if err != nil {
		// 			b.Fatal(err)
		// 		}
		// 	}
		// })

		b.Run(bench.pkg+"/ParseProject", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := ParseProject(options)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		project, err := ParseProject(options)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/ParsePackages", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := ParsePackages(project)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		project.pkgs, err = ParsePackages(project)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/GenerateBindings", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := project.GenerateBindings()
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(bench.pkg+"/GenerateModels", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := project.GenerateModels()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
