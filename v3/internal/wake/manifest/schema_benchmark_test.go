package manifest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var (
	benchmarkFileSink    *hcl.File
	benchmarkRawSink     hclDocument
	benchmarkConfigSink  Config
	benchmarkOriginsSink map[string]Origin
	benchmarkBytesSink   []byte
)

func BenchmarkManifestDecode(b *testing.B) {
	root := b.TempDir()
	for _, fixture := range manifestBenchmarkFixtures() {
		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := decodeHCL(root, "wails.hcl", fixture.source, ""); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkManifestEncode(b *testing.B) {
	root := b.TempDir()
	for _, fixture := range manifestBenchmarkFixtures() {
		loaded, err := decodeHCL(root, "wails.hcl", fixture.source, "")
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			var encoded []byte
			for b.Loop() {
				encoded, err = EncodeConfig(loaded.Config)
				if err != nil {
					b.Fatal(err)
				}
			}
			benchmarkBytesSink = encoded
			b.ReportMetric(float64(len(encoded)), "output_bytes")
		})
	}
}

func BenchmarkManifestStages(b *testing.B) {
	root := b.TempDir()
	for _, fixture := range manifestBenchmarkFixtures() {
		b.Run(fixture.name, func(b *testing.B) {
			file, diagnostics := hclsyntax.ParseConfig(fixture.source, "wails.hcl", hcl.InitialPos)
			if diagnostics.HasErrors() {
				b.Fatal(diagnostics.Error())
			}
			body := file.Body.(*hclsyntax.Body)
			raw, err := decodeManifestSchema(body)
			if err != nil {
				b.Fatal(err)
			}
			document, err := documentFromHCL(raw)
			if err != nil {
				b.Fatal(err)
			}
			config := configFromDocument(root, "", document)

			b.Run("parse", func(b *testing.B) {
				b.ReportAllocs()
				var parsed *hcl.File
				for b.Loop() {
					var diagnostics hcl.Diagnostics
					parsed, diagnostics = hclsyntax.ParseConfig(fixture.source, "wails.hcl", hcl.InitialPos)
					if diagnostics.HasErrors() {
						b.Fatal(diagnostics.Error())
					}
				}
				benchmarkFileSink = parsed
			})

			b.Run("structural-decode", func(b *testing.B) {
				b.ReportAllocs()
				var decoded hclDocument
				for b.Loop() {
					decoded, err = decodeManifestSchema(body)
					if err != nil {
						b.Fatal(err)
					}
				}
				benchmarkRawSink = decoded
			})

			b.Run("semantic-resolution", func(b *testing.B) {
				b.ReportAllocs()
				var resolved Config
				for b.Loop() {
					document, err = documentFromHCL(raw)
					if err != nil {
						b.Fatal(err)
					}
					resolved = configFromDocument(root, "", document)
				}
				benchmarkConfigSink = resolved
			})

			b.Run("origin-tracking", func(b *testing.B) {
				b.ReportAllocs()
				var origins map[string]Origin
				for b.Loop() {
					origins = manifestOrigins(body)
				}
				benchmarkOriginsSink = origins
			})

			b.Run("semantic-validation", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					if err := validateConfig(config); err != nil {
						b.Fatal(err)
					}
				}
			})

			b.Run("encode", func(b *testing.B) {
				b.ReportAllocs()
				var encoded []byte
				for b.Loop() {
					encoded, err = EncodeConfig(config)
					if err != nil {
						b.Fatal(err)
					}
				}
				benchmarkBytesSink = encoded
			})
		})
	}
}

type manifestBenchmarkFixture struct {
	name   string
	source []byte
}

func manifestBenchmarkFixtures() []manifestBenchmarkFixture {
	minimal := []byte(`version = 3

project {
  name = "bench"
  product_name = "Bench"
  identifier = "com.example.bench"
  version = "1.0.0"
}
`)
	representative := []byte(`version = 3

project {
  name = "bench"
  product_name = "Bench"
  identifier = "com.example.bench"
  version = "1.0.0"
  company = "Example"
  binary_name = "bench"
}

frontend {
  directory = "frontend"
  install = ["pnpm", "install", "--frozen-lockfile"]
  build = ["pnpm", "run", "build"]
  dev = ["pnpm", "run", "dev"]
  output = "frontend/dist"
}

build {
  output = "dist"
  tags = ["sqlite_fts5"]
  trim_path = true
  strip = true
}

windows {
  publisher = "CN=Example"
  signing {
    credential = "windows-release"
    timestamp_server = "https://timestamp.example.com"
  }
}

target "windows/amd64" {
  tags = ["enterprise"]
}

profile "release" {
  target "windows/amd64" {
    formats = ["nsis"]
    sign = true
  }
}
`)
	var large strings.Builder
	large.Write(minimal)
	for index := range 1000 {
		fmt.Fprintf(&large, "\nfile_association \"kind-%04d\" {\n  extensions = [\"kind-%04d\"]\n  name = \"Kind %04d\"\n  platforms = [\"windows\", \"darwin\", \"linux\"]\n}\n", index, index, index)
	}
	return []manifestBenchmarkFixture{
		{name: "minimal", source: minimal},
		{name: "representative", source: representative},
		{name: "large", source: []byte(large.String())},
	}
}
