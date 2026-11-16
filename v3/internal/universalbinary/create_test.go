package universalbinary

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func thin(cpu macho.Cpu, sub uint32, order binary.ByteOrder) []byte {
	var data bytes.Buffer
	_ = binary.Write(&data, order, macho.FileHeader{Magic: macho.Magic64, Cpu: cpu, SubCpu: sub, Type: macho.TypeObj})
	_ = binary.Write(&data, order, uint32(0))
	return append(data.Bytes(), bytes.Repeat([]byte("payload"), 101)...)
}

func input(t testing.TB, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCreatePreservesSlices(t *testing.T) {
	amd := thin(macho.CpuAmd64, 3, binary.LittleEndian)
	arm := thin(macho.CpuArm64, 0, binary.LittleEndian)
	out := filepath.Join(t.TempDir(), "universal")
	if err := Create(out, []string{input(t, "arm", arm), input(t, "amd", amd)}); err != nil {
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
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range [][]byte{amd, arm} {
		a := f.Arches[i]
		if a.Offset%(1<<a.Align) != 0 || !bytes.Equal(data[a.Offset:a.Offset+a.Size], want) {
			t.Errorf("slice %d changed or misaligned", i)
		}
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0100 == 0 {
		t.Error("output lost executable mode")
	}
}

func archive(data []byte, extended bool) []byte {
	name := "object.o/"
	if extended {
		name = "#1/16"
		data = append([]byte("long-object.o\x00\x00\x00"), data...)
	}
	header := fmt.Sprintf("%-16s%-12d%-6d%-6d%-8s%-10d`\n", name, 0, 0, 0, "644", len(data))
	result := append([]byte("!<arch>\n"+header), data...)
	if len(data)%2 != 0 {
		result = append(result, '\n')
	}
	return result
}

func TestCreateArchivesAndUniversalInput(t *testing.T) {
	a := archive(thin(macho.CpuAmd64, 3, binary.LittleEndian), false)
	b := archive(thin(macho.CpuArm64, 0, binary.LittleEndian), true)
	first := filepath.Join(t.TempDir(), "first")
	if err := Create(first, []string{input(t, "a", a), input(t, "b", b)}); err != nil {
		t.Fatal(err)
	}
	second := filepath.Join(t.TempDir(), "second")
	if err := Create(second, []string{first}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(second)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range [][]byte{a, b} {
		row := data[8+20*i:]
		offset, size := binary.BigEndian.Uint32(row[8:]), binary.BigEndian.Uint32(row[12:])
		if !bytes.Equal(data[offset:offset+size], want) {
			t.Errorf("archive %d changed", i)
		}
	}
}

func TestCreateRejectsInvalidInputsWithoutReplacingOutput(t *testing.T) {
	valid := input(t, "amd", thin(macho.CpuAmd64, 3, binary.LittleEndian))
	for name, data := range map[string][]byte{"empty": nil, "text": []byte("not a binary"), "truncated": {0xcf, 0xfa, 0xed, 0xfe}, "empty_archive": []byte("!<arch>\n"), "bad_archive": []byte("!<arch>\ninvalid"), "bad_fat": {0xca, 0xfe, 0xba, 0xbe, 0, 0, 0, 2}} {
		t.Run(name, func(t *testing.T) {
			out := input(t, "existing", []byte("keep"))
			err := Create(out, []string{valid, input(t, name, data)})
			got, readErr := os.ReadFile(out)
			if err == nil || readErr != nil || string(got) != "keep" {
				t.Fatalf("err=%v, output=%q, read=%v", err, got, readErr)
			}
			entries, err := os.ReadDir(filepath.Dir(out))
			if err != nil || len(entries) != 1 {
				t.Fatalf("staging left behind: %v, %v", entries, err)
			}
		})
	}
	if err := Create(filepath.Join(t.TempDir(), "duplicate"), []string{valid, valid}); err == nil {
		t.Error("accepted duplicate architecture")
	}
	if err := Create(valid, []string{valid}); err == nil {
		t.Error("accepted output aliasing input")
	}
	if err := Create(filepath.Join(t.TempDir(), "none"), nil); err == nil {
		t.Error("accepted no input")
	}
}

func TestCreateByteOrdersAnd32Bit(t *testing.T) {
	var header bytes.Buffer
	_ = binary.Write(&header, binary.BigEndian, macho.FileHeader{Magic: macho.Magic32, Cpu: macho.CpuPpc, Type: macho.TypeObj})
	out := filepath.Join(t.TempDir(), "universal")
	if err := Create(out, []string{input(t, "ppc", header.Bytes()), input(t, "amd", thin(macho.CpuAmd64, 3, binary.LittleEndian))}); err != nil {
		t.Fatal(err)
	}
	f, err := macho.OpenFat(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if len(f.Arches) != 2 || f.Arches[0].Cpu != macho.CpuPpc {
		t.Fatalf("unexpected arches: %v", f.Arches)
	}
}

func BenchmarkCreate(b *testing.B) {
	for _, size := range []int{1 << 20, 16 << 20} {
		b.Run(fmt.Sprintf("%dMiB_per_arch", size>>20), func(b *testing.B) {
			a := append(thin(macho.CpuAmd64, 3, binary.LittleEndian), make([]byte, size)...)
			c := append(thin(macho.CpuArm64, 0, binary.LittleEndian), make([]byte, size)...)
			paths := []string{input(b, "amd", a), input(b, "arm", c)}
			out := filepath.Join(b.TempDir(), "universal")
			b.ReportAllocs()
			b.SetBytes(int64(len(a) + len(c)))
			b.ResetTimer()
			for range b.N {
				if err := Create(out, paths); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestCreateReadsFat64(t *testing.T) {
	payload := thin(macho.CpuArm64, 0, binary.LittleEndian)
	data := make([]byte, 40)
	binary.BigEndian.PutUint32(data, 0xcafebabf)
	binary.BigEndian.PutUint32(data[4:], 1)
	binary.BigEndian.PutUint32(data[8:], uint32(macho.CpuArm64))
	binary.BigEndian.PutUint64(data[16:], 40)
	binary.BigEndian.PutUint64(data[24:], uint64(len(payload)))
	binary.BigEndian.PutUint32(data[32:], 2)
	out := filepath.Join(t.TempDir(), "universal")
	if err := Create(out, []string{input(t, "fat64", append(data, payload...))}); err != nil {
		t.Fatal(err)
	}
	f, err := macho.OpenFat(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if len(f.Arches) != 1 || f.Arches[0].Cpu != macho.CpuArm64 {
		t.Fatal("lost fat64 input architecture")
	}
	binary.BigEndian.PutUint64(data[16:], ^uint64(0))
	if err := Create(out, []string{input(t, "bad-fat64", append(data, payload...))}); err == nil {
		t.Error("accepted fat64 overflow")
	}
}

func TestPublishFailureRemovesStage(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, "directory")
	if err := os.Mkdir(out, 0700); err != nil {
		t.Fatal(err)
	}
	source := input(t, "amd", thin(macho.CpuAmd64, 3, binary.LittleEndian))
	if err := Create(out, []string{source}); err == nil {
		t.Fatal("replaced directory with binary")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 || !entries[0].IsDir() {
		t.Fatalf("publication failure left staging files: %v, %v", entries, err)
	}
}

func TestLayoutRejectsOverflow(t *testing.T) {
	for _, img := range []image{{size: 1 << 32, align: 2}, {size: 4, align: 31}, {size: 0, align: 2}} {
		if _, _, err := layout([]image{img}); err == nil {
			t.Errorf("accepted invalid image: %+v", img)
		}
	}
}

func sectionImage(wide bool, alignment uint32) []byte {
	var data bytes.Buffer
	if wide {
		_ = binary.Write(&data, binary.LittleEndian, macho.FileHeader{Magic: macho.Magic64, Cpu: macho.CpuArm64, Type: macho.TypeObj, Ncmd: 1, Cmdsz: 152})
		_ = binary.Write(&data, binary.LittleEndian, uint32(0))
		_ = binary.Write(&data, binary.LittleEndian, macho.Segment64{Cmd: macho.LoadCmdSegment64, Len: 152, Nsect: 1})
		_ = binary.Write(&data, binary.LittleEndian, macho.Section64{Align: alignment})
	} else {
		_ = binary.Write(&data, binary.LittleEndian, macho.FileHeader{Magic: macho.Magic32, Cpu: macho.Cpu386, Type: macho.TypeObj, Ncmd: 1, Cmdsz: 124})
		_ = binary.Write(&data, binary.LittleEndian, macho.Segment32{Cmd: macho.LoadCmdSegment, Len: 124, Nsect: 1})
		_ = binary.Write(&data, binary.LittleEndian, macho.Section32{Align: alignment})
	}
	return data.Bytes()
}

func TestCreateHonorsSectionAlignment(t *testing.T) {
	for _, wide := range []bool{false, true} {
		data := sectionImage(wide, 18)
		out := filepath.Join(t.TempDir(), "universal")
		if err := Create(out, []string{input(t, "section", data)}); err != nil {
			t.Fatal(err)
		}
		f, err := macho.OpenFat(out)
		if err != nil {
			t.Fatal(err)
		}
		if f.Arches[0].Align != 18 || f.Arches[0].Offset%(1<<18) != 0 {
			t.Error("section alignment lost")
		}
		_ = f.Close()
	}
}

func TestRejectMalformedLoadCommands(t *testing.T) {
	for _, kind := range []string{"count", "length", "sections", "truncated"} {
		t.Run(kind, func(t *testing.T) {
			data := sectionImage(true, 4)
			switch kind {
			case "count":
				binary.LittleEndian.PutUint32(data[16:], ^uint32(0))
			case "length":
				binary.LittleEndian.PutUint32(data[36:], 0)
			case "sections":
				binary.LittleEndian.PutUint32(data[96:], 100)
			case "truncated":
				data = data[:64]
			}
			if _, _, _, err := machoArchitecture(bytes.NewReader(data), int64(len(data))); err == nil {
				t.Error("accepted malformed load commands")
			}
		})
	}
}

func FuzzMetadata(f *testing.F) {
	f.Add(thin(macho.CpuAmd64, 3, binary.LittleEndian))
	f.Add(sectionImage(true, 14))
	f.Add(archive(thin(macho.CpuArm64, 0, binary.LittleEndian), true))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64<<10 {
			t.Skip()
		}
		_, _, _, _ = machoArchitecture(bytes.NewReader(data), int64(len(data)))
		if bytes.HasPrefix(data, []byte("!<arch>\n")) {
			_, _, _ = archiveArchitecture(bytes.NewReader(data), int64(len(data)))
		}
	})
}
