package parser

import (
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
	}

	for _, bench := range benchmarks {

		options := &flags.GenerateBindingsOptions{
			ModelsFilename: "models",
			TS:             true,
		}

		pattern := "github.com/wailsapp/wails/v3/internal/parser/testdata/" + bench.pkg

		b.Run(bench.pkg+"/load", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LoadPackages(nil, true,
					pattern,
					WailsAppPkgPath,
				)
			}

		})

		pPkgs, err := LoadPackages(nil, true,
			pattern, WailsAppPkgPath,
		)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/analyze", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Services(pPkgs)
			}
		})

		project, err := ParseProject([]string{pattern}, options)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(bench.pkg+"/bindings", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				project.GenerateBindings(options)
			}
		})

		b.Run(bench.pkg+"/models", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				project.GenerateModels(options)
			}
		})
	}
}
