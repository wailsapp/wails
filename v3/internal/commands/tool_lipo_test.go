package commands

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"github.com/wailsapp/wails/v3/internal/flags"
	"os"
	"path/filepath"
	"testing"
)

func TestToolLipo(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	dir := t.TempDir()
	var paths []string
	for i, cpu := range []macho.Cpu{macho.CpuAmd64, macho.CpuArm64} {
		var data bytes.Buffer
		_ = binary.Write(&data, binary.LittleEndian, macho.FileHeader{Magic: macho.Magic64, Cpu: cpu, Type: macho.TypeObj})
		_ = binary.Write(&data, binary.LittleEndian, uint32(0))
		path := filepath.Join(dir, []string{"amd", "arm"}[i])
		if err := os.WriteFile(path, data.Bytes(), 0755); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	out := filepath.Join(dir, "universal")
	if err := ToolLipo(&flags.Lipo{Inputs: paths, Output: out}); err != nil {
		t.Fatal(err)
	}
	f, err := macho.OpenFat(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if len(f.Arches) != 2 {
		t.Fatalf("arches=%d", len(f.Arches))
	}
	for _, opt := range []*flags.Lipo{{}, {Inputs: paths}, {Inputs: paths[:1], Output: out}} {
		if err := ToolLipo(opt); err == nil {
			t.Error("accepted missing inputs or output")
		}
	}
}
