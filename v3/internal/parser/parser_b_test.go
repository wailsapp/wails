package parser

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os/exec"
	"path/filepath"
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

func BenchmarkLoadPackage(b *testing.B) {

	benchmarks := []struct {
		pkg string
	}{
		{
			pkg: "fmt",
		},
	}

	for _, bench := range benchmarks {

		b.Run(bench.pkg+"/LoadPackage", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := LoadPackage(nil, true, bench.pkg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(bench.pkg+"/ParseFile", func(b *testing.B) {
			type Package struct {
				Dir     string
				GoFiles []string
			}

			for i := 0; i < b.N; i++ {
				cmd := exec.Command("go", "list", "-json=Dir,GoFiles", bench.pkg)
				buf, err := cmd.Output()
				if err != nil {
					b.Fatal(err)
				}

				var p Package
				if err := json.Unmarshal(buf, &p); err != nil {
					b.Fatal(err)
				}

				fset := token.NewFileSet()
				for _, filename := range p.GoFiles {
					_, err := parser.ParseFile(fset, filepath.Join(p.Dir, filename), nil, parser.AllErrors|parser.ParseComments|parser.SkipObjectResolution)
					if err != nil {
						b.Fatal(err)
					}
				}

			}

		})

	}

}
