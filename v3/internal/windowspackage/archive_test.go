package windowspackage

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) (Options, map[string][]byte) {
	t.Helper()
	root := t.TempDir()
	options := Options{Directory: filepath.Join(root, "input"), Output: filepath.Join(root, "output", "app.msix")}
	files := map[string][]byte{"AppxManifest.xml": []byte(`<?xml version="1.0"?><Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"/>`)}
	for _, size := range []int{0, 1, 65535, 65536, 65537, 131073} {
		data := make([]byte, size)
		_, _ = rand.New(rand.NewSource(int64(size))).Read(data)
		files[fmt.Sprintf("Assets/%d.bin", size)] = data
	}
	for name, data := range files {
		path := filepath.Join(options.Directory, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	return options, files
}

func TestArchiveBlockMapAndCompressionBoundaries(t *testing.T) {
	for _, level := range []int{1, 6, 9} {
		t.Run(fmt.Sprint(level), func(t *testing.T) {
			options, source := fixture(t)
			options.Compression = level
			if err := Create(context.Background(), options); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(options.Output)
			if err != nil {
				t.Fatal(err)
			}
			archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				t.Fatal(err)
			}
			entries := make(map[string]*zip.File)
			for _, file := range archive.File {
				name, err := url.PathUnescape(file.Name)
				if err != nil {
					t.Fatal(err)
				}
				entries[name] = file
				stream, err := file.Open()
				if err != nil {
					t.Fatal(err)
				}
				body, err := io.ReadAll(stream)
				stream.Close()
				if err != nil {
					t.Fatal(err) // ZIP reader independently checks the CRC and sizes.
				}
				if want, ok := source[name]; ok && !bytes.Equal(body, want) {
					t.Fatalf("payload %s changed", file.Name)
				}
				if file.Flags&8 != 0 {
					t.Fatalf("unexpected ZIP data descriptor: %s", file.Name)
				}
			}
			stream, err := entries["AppxBlockMap.xml"].Open()
			if err != nil {
				t.Fatal(err)
			}
			var blockMap struct {
				Files []struct {
					Name   string `xml:"Name,attr"`
					Size   uint64 `xml:"Size,attr"`
					Header uint64 `xml:"LfhSize,attr"`
					Blocks []struct {
						Hash string `xml:"Hash,attr"`
						Size uint64 `xml:"Size,attr"`
					} `xml:"Block"`
				} `xml:"File"`
			}
			if err = xml.NewDecoder(stream).Decode(&blockMap); err != nil {
				t.Fatal(err)
			}
			stream.Close()
			if len(blockMap.Files) != len(source) {
				t.Fatal("block map does not cover every payload/manifest")
			}
			for _, file := range blockMap.Files {
				name := strings.ReplaceAll(file.Name, "\\", "/")
				want := source[name]
				entry := entries[name]
				offset, err := entry.DataOffset()
				if err != nil {
					t.Fatal(err)
				}
				lfh := offset - int64(file.Header)
				if file.Size != uint64(len(want)) || binary.LittleEndian.Uint32(data[lfh:lfh+4]) != 0x04034b50 {
					t.Fatalf("bad size/local header for %s", name)
				}
				for i, block := range file.Blocks {
					start := i * 65536
					end := min(start+65536, len(want))
					hash := sha256.Sum256(want[start:end])
					if block.Hash != base64.StdEncoding.EncodeToString(hash[:]) {
						t.Fatalf("bad block hash %s:%d", name, i)
					}
					if entry.Method == zip.Deflate {
						compressed := append([]byte(nil), data[offset:offset+int64(block.Size)]...)
						if i != len(file.Blocks)-1 {
							compressed = append(compressed, 1, 0, 0, 255, 255)
						}
						reader := flate.NewReader(bytes.NewReader(compressed))
						decoded, err := io.ReadAll(reader)
						reader.Close()
						if err != nil || !bytes.Equal(decoded, want[start:end]) {
							t.Fatalf("block is not independently decompressible %s:%d: %v", name, i, err)
						}
						offset += int64(block.Size)
					}
				}
			}
			if err := Create(context.Background(), options); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(options.Output)
			if err != nil || !bytes.Equal(data, after) {
				t.Fatal("archive is not deterministic")
			}
		})
	}
}

func TestArchiveFailurePreservesPublishedFile(t *testing.T) {
	options, _ := fixture(t)
	if err := os.MkdirAll(filepath.Dir(options.Output), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(options.Output, []byte("previous"), 0644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Create(ctx, options); err == nil {
		t.Fatal("cancelled archive succeeded")
	}
	if err := os.WriteFile(filepath.Join(options.Directory, "AppxBlockMap.xml"), []byte("reserved"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Create(context.Background(), options); err == nil {
		t.Fatal("reserved input accepted")
	}
	if err := os.Remove(filepath.Join(options.Directory, "AppxBlockMap.xml")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(options.Directory, "AppxManifest.xml"), []byte("<Package>broken"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Create(context.Background(), options); err == nil {
		t.Fatal("malformed manifest accepted")
	}
	data, err := os.ReadFile(options.Output)
	if err != nil || string(data) != "previous" {
		t.Fatal("failure replaced previous package")
	}
	stages, err := filepath.Glob(filepath.Join(filepath.Dir(options.Output), ".wake-msix-*"))
	if err != nil || len(stages) != 0 {
		t.Fatalf("staging remains: %v %v", stages, err)
	}
}

func TestArchiveRejectsWindowsNameAliases(t *testing.T) {
	options, _ := fixture(t)
	if err := os.WriteFile(filepath.Join(options.Directory, "appxmanifest.xml"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Create(context.Background(), options); err == nil {
		t.Fatal("case alias accepted")
	}
}

func TestArchiveWithMicrosoftSDK(t *testing.T) {
	tool := os.Getenv("WAILS_TEST_MAKEMSIX")
	input := os.Getenv("WAILS_TEST_MSIX_INPUT")
	if tool == "" {
		t.Skip("set WAILS_TEST_MAKEMSIX for SDK interoperability")
	}
	if input == "" {
		input = sdkFixture(t)
	}
	output := filepath.Join(t.TempDir(), "native.msix")
	if err := Create(context.Background(), Options{Directory: input, Output: output}); err != nil {
		t.Fatal(err)
	}
	unpacked := filepath.Join(t.TempDir(), "unpacked")
	cmd := exec.Command(tool, "unpack", "/p", output, "/d", unpacked)
	if result, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Microsoft SDK rejected package: %v\n%s", err, result)
	}
	if err := filepath.WalkDir(input, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		relative, err := filepath.Rel(input, path)
		if err != nil {
			return err
		}
		original, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		extracted, err := os.ReadFile(filepath.Join(unpacked, relative))
		if err != nil {
			return err
		}
		if !bytes.Equal(original, extracted) {
			return fmt.Errorf("SDK extracted different bytes for %s", relative)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func sdkFixture(t *testing.T) string {
	t.Helper()
	options, _ := fixture(t)
	manifest := `<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10" xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10" xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities" IgnorableNamespaces="uap rescap">
<Identity Name="Wake.Native.Fixture" Publisher="CN=Wake fixture" Version="1.2.3.0" ProcessorArchitecture="x64"/>
<Properties><DisplayName>Wake fixture</DisplayName><PublisherDisplayName>Wake fixture</PublisherDisplayName><Logo>Assets/logo.png</Logo></Properties>
<Resources><Resource Language="en-us"/></Resources>
<Dependencies><TargetDeviceFamily Name="Windows.Desktop" MinVersion="10.0.17763.0" MaxVersionTested="10.0.26100.0"/></Dependencies>
<Applications><Application Id="App" Executable="app.exe" EntryPoint="Windows.FullTrustApplication"><uap:VisualElements DisplayName="Wake fixture" Description="Native SDK interoperability fixture" Square150x150Logo="Assets/logo.png" Square44x44Logo="Assets/logo.png" BackgroundColor="transparent"/></Application></Applications>
<Capabilities><rescap:Capability Name="runFullTrust"/></Capabilities></Package>`
	if err := os.WriteFile(filepath.Join(options.Directory, "AppxManifest.xml"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	var logo bytes.Buffer
	if err := png.Encode(&logo, image.NewRGBA(image.Rect(0, 0, 150, 150))); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(options.Directory, "Assets", "logo.png"), logo.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(options.Directory, "app.exe"), data, 0600); err != nil {
		t.Fatal(err)
	}
	return options.Directory
}

func BenchmarkPayloadCompression(b *testing.B) {
	path, err := os.Executable()
	if err != nil {
		b.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatal(err)
	}
	for _, level := range []int{1, 4, 6, 9} {
		b.Run(fmt.Sprint(level), func(b *testing.B) {
			buffer := make([]byte, 65536)
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for b.Loop() {
				header, _, err := encodePayload(context.Background(), io.Discard, bytes.NewReader(data), "app.exe", level, true, buffer)
				if err != nil {
					b.Fatal(err)
				}
				b.ReportMetric(float64(header.CompressedSize64), "compressed-bytes")
			}
		})
	}
}
