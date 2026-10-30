package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkCp(b *testing.B) {
	previousArgs, previousFooter := os.Args, DisableFooter
	b.Cleanup(func() { os.Args, DisableFooter = previousArgs, previousFooter })
	for _, size := range []int{1 << 20, 16 << 20} {
		for _, destinations := range []int{1, 2} {
			b.Run(fmt.Sprintf("%dMiB/%d_destinations", size>>20, destinations), func(b *testing.B) {
				root := b.TempDir()
				source := filepath.Join(root, "source")
				content := bytes.Repeat([]byte("binary payload!\x00"), size/16)
				if err := os.WriteFile(source, content, 0600); err != nil {
					b.Fatal(err)
				}
				os.Args = []string{"wails3", "tool", "cp", source}
				for i := range destinations {
					os.Args = append(os.Args, filepath.Join(root, fmt.Sprintf("target-%d", i)))
				}
				b.ReportAllocs()
				b.SetBytes(int64(size * destinations))
				b.ResetTimer()
				for range b.N {
					if err := Cp(&CpOptions{}); err != nil {
						b.Fatal(err)
					}
				}
				b.StopTimer()
				for _, target := range os.Args[4:] {
					got, err := os.ReadFile(target)
					if err != nil || !bytes.Equal(got, content) {
						b.Fatalf("copied content differs: target=%s, error=%v", target, err)
					}
				}
			})
		}
	}
}
