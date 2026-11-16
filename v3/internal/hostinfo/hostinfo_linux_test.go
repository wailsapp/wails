package hostinfo

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

const cpuFixture = `processor : 0
model name : Example CPU
physical id : 0
core id : 0
cpu cores : 2

processor : 1
model name : Example CPU
physical id : 0
core id : 0
cpu cores : 2

processor : 2
model name : Example CPU
physical id : 0
core id : 1
cpu cores : 2

processor : 3
model name : Example CPU
physical id : 0
core id : 1
cpu cores : 2`

func TestReadProcessors(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		want        []Processor
	}{
		{"SMT", cpuFixture, []Processor{{"Example CPU", 2}}},
		{"two sockets", cpuFixture + "\n\nprocessor : 4\nphysical id : 1\nmodel name : Second CPU\ncore id : 0\n", []Processor{{"Example CPU", 2}, {"Second CPU", 1}}},
		{"ARM without topology", "processor : 0\nProcessor : ARMv7\n\nprocessor : 1\nProcessor : ARMv7\n", []Processor{{"ARMv7", 2}}},
		{"no model", "processor : 0\n", []Processor{{"Unknown", 1}}},
		{"whitespace", "processor\t: 0\nmodel name\t:  Example CPU  \n\n\n", []Processor{{"Example CPU", 1}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readProcessors(strings.NewReader(tt.input))
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v, %v; want %+v", got, err, tt.want)
			}
		})
	}
}

func TestReadProcessorsErrors(t *testing.T) {
	for _, input := range []string{"", "Hardware : example\n", strings.Repeat("x", 70*1024)} {
		if _, err := readProcessors(strings.NewReader(input)); err == nil {
			t.Fatal("expected malformed or empty CPU data to fail")
		}
	}
	if _, err := readProcessors(errorReader{}); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("reader error lost: %v", err)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestPhysicalMemory(t *testing.T) {
	for _, tt := range []struct {
		name  string
		files map[string]string
		want  uint64
	}{
		{"online blocks", map[string]string{
			"sys/devices/system/memory/block_size_bytes": "8000000\n",
			"sys/devices/system/memory/memory0/state":    "online\n",
			"sys/devices/system/memory/memory1/state":    "offline\n",
			"sys/devices/system/memory/memory2/state":    "online\n",
			"sys/devices/system/memory/memorybad/state":  "online\n",
			"sys/devices/system/memory/memory3/state":    "going-offline\n",
			"proc/meminfo": "MemTotal: 100000 kB\n",
		}, 256 * 1024 * 1024},
		{"proc fallback", map[string]string{"proc/meminfo": "MemFree: 99 kB\nMemTotal: 2048 kB\n"}, 2 * 1024 * 1024},
		{"bad block size", map[string]string{"sys/devices/system/memory/block_size_bytes": "invalid", "proc/meminfo": "MemTotal: 2048 kB\n"}, 2 * 1024 * 1024},
		{"all offline", map[string]string{"sys/devices/system/memory/block_size_bytes": "8000000", "sys/devices/system/memory/memory0/state": "offline", "proc/meminfo": "MemTotal: 2048 kB"}, 2 * 1024 * 1024},
		{"block overflow", map[string]string{"sys/devices/system/memory/block_size_bytes": "ffffffffffffffff", "sys/devices/system/memory/memory0/state": "online", "sys/devices/system/memory/memory1/state": "online", "proc/meminfo": "MemTotal: 2048 kB"}, 2 * 1024 * 1024},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := memoryFixture(tt.files)
			got, err := physicalMemory(root)
			if err != nil || got != tt.want {
				t.Fatalf("got %d, %v; want %d", got, err, tt.want)
			}
		})
	}
	for _, input := range []string{"", "MemTotal: -1 kB", "MemTotal: 18446744073709551615 kB", "MemTotal: 123 MB", "MemTotal: 0 kB"} {
		if _, err := physicalMemory(memoryFixture(map[string]string{"proc/meminfo": input})); err == nil {
			t.Fatalf("accepted invalid memory: %q", input)
		}
	}
	if _, err := physicalMemory(fstest.MapFS{}); err == nil {
		t.Fatal("missing memory information accepted")
	}
}

func memoryFixture(files map[string]string) fstest.MapFS {
	result := make(fstest.MapFS, len(files))
	for path, data := range files {
		result[path] = &fstest.MapFile{Data: []byte(data)}
	}
	return result
}

func TestReadGraphics(t *testing.T) {
	directory := t.TempDir()
	files := map[string]string{
		"0000:01:00.0/class": "0x030000\n", "0000:01:00.0/vendor": "0x10de\n", "0000:01:00.0/device": "0x1234\n",
		"0000:02:00.0/class": "0x030200\n", "0000:02:00.0/vendor": "0xabcd\n", "0000:02:00.0/device": "0x5678\n",
		"0000:03:00.0/class": "0x020000\n",
		"0000:04:00.0/class": "0x030000\n", "0000:04:00.0/vendor": "invalid",
		"0000:05:00.0/class": "0x030000\n", "0000:05:00.0/vendor": "0xabcd", "0000:05:00.0/device": "invalid",
		"0000:06:00.0/class": "invalid",
	}
	for path, data := range files {
		path = filepath.Join(directory, path)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("../../drivers/example", filepath.Join(directory, "0000:01:00.0/driver")); err != nil {
		t.Fatal(err)
	}
	calls := 0
	lookup := func() (nameLookup, error) {
		calls++
		return func(vendor, device string) (string, string) {
			if vendor == "10de" && device == "1234" {
				return "Example vendor", "Example GPU"
			}
			return "", ""
		}, nil
	}
	got, err := readGraphics(directory, lookup)
	want := []GraphicsCard{{"Example GPU", "Example vendor", "example"}, {"PCI abcd:5678", "", ""}}
	if err != nil || !reflect.DeepEqual(got, want) || calls != 1 {
		t.Fatalf("got %+v, %v, lookups %d; want %+v", got, err, calls, want)
	}
	lookupError := errors.New("database unavailable")
	if _, err := readGraphics(directory, func() (nameLookup, error) { return nil, lookupError }); !errors.Is(err, lookupError) {
		t.Fatalf("database error lost: %v", err)
	}
	if _, err := readGraphics(filepath.Join(directory, "missing"), lookup); err == nil {
		t.Fatal("missing PCI directory accepted")
	}
	calls = 0
	got, err = readGraphics(t.TempDir(), lookup)
	if err != nil || len(got) != 0 || calls != 0 {
		t.Fatalf("empty system should not load database: %+v, %v, %d", got, err, calls)
	}
}

func BenchmarkReadProcessors(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := readProcessors(strings.NewReader(cpuFixture)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPhysicalMemory(b *testing.B) {
	root := memoryFixture(map[string]string{
		"sys/devices/system/memory/block_size_bytes": "8000000\n",
		"sys/devices/system/memory/memory0/state":    "online\n",
		"sys/devices/system/memory/memory1/state":    "offline\n",
	})
	b.ReportAllocs()
	for b.Loop() {
		if _, err := physicalMemory(root); err != nil {
			b.Fatal(err)
		}
	}
}

func FuzzReadProcessors(f *testing.F) {
	f.Add(cpuFixture)
	f.Add("processor : 0\nProcessor : ARMv7\n")
	f.Fuzz(func(t *testing.T, input string) {
		processors, err := readProcessors(strings.NewReader(input))
		if err == nil {
			for _, p := range processors {
				if p.Cores < 1 || p.Model == "" {
					t.Fatalf("invalid processor: %+v", p)
				}
			}
		}
	})
}
