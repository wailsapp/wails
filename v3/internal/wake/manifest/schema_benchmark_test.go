package manifest

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	benchmarkFileSink    *yaml.Node
	benchmarkRawSink     manifestDocument
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
				if _, err := decodeYAML(root, "wails.yaml", fixture.source, ""); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkManifestEncode(b *testing.B) {
	root := b.TempDir()
	for _, fixture := range manifestBenchmarkFixtures() {
		loaded, err := decodeYAML(root, "wails.yaml", fixture.source, "")
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
			raw, origins, err := decodeManifestSchema(fixture.source, "wails.yaml")
			if err != nil {
				b.Fatal(err)
			}
			document, err := documentFromManifest(raw)
			if err != nil {
				b.Fatal(err)
			}
			config := configFromDocument(root, "", document)

			b.Run("parse", func(b *testing.B) {
				b.ReportAllocs()
				var parsed yaml.Node
				for b.Loop() {
					if err := yaml.Unmarshal(fixture.source, &parsed); err != nil {
						b.Fatal(err)
					}
				}
				benchmarkFileSink = &parsed
			})

			b.Run("structural-decode", func(b *testing.B) {
				b.ReportAllocs()
				var decoded manifestDocument
				for b.Loop() {
					decoded, _, err = decodeManifestSchema(fixture.source, "wails.yaml")
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
					document, err = documentFromManifest(raw)
					if err != nil {
						b.Fatal(err)
					}
					resolved = configFromDocument(root, "", document)
				}
				benchmarkConfigSink = resolved
			})

			b.Run("origin-tracking", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					_, origins, err = decodeManifestSchema(fixture.source, "wails.yaml")
					if err != nil {
						b.Fatal(err)
					}
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

func BenchmarkSchemaReferenceMarkdown(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if len(SchemaReferenceMarkdown()) == 0 {
			b.Fatal("empty schema reference")
		}
	}
}

type manifestBenchmarkFixture struct {
	name   string
	source []byte
}

func manifestBenchmarkFixtures() []manifestBenchmarkFixture {
	minimal := []byte(`version: 3
project:
  name: bench
  product_name: Bench
  identifier: com.example.bench
  version: 1.0.0`)
	representative := []byte(`version: 3
project:
  name: bench
  product_name: Bench
  identifier: com.example.bench
  version: 1.0.0
  company: Example
  binary_name: bench
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: dist
  tags: [sqlite_fts5]
  trim_path: true
  strip: true
windows:
  publisher: CN=Example
  signing:
    credential: windows-release
    timestamp_server: https://timestamp.example.com
targets:
  windows/amd64:
    tags: [enterprise]
profiles:
  release:
    targets:
      windows/amd64:
        formats: [nsis]
        sign: true`)
	var large strings.Builder
	large.Write(minimal)
	large.WriteString("\nfile_associations:\n")
	for index := range 1000 {
		fmt.Fprintf(&large, "  kind-%04d:\n    extensions: [kind-%04d]\n    name: Kind %04d\n    platforms: [windows, darwin, linux]\n", index, index, index)
	}
	return []manifestBenchmarkFixture{
		{name: "minimal", source: minimal},
		{name: "representative", source: representative},
		{name: "large", source: []byte(large.String())},
	}
}
