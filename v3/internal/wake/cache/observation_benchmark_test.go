package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkDirectoryObservation(b *testing.B) {
	root := b.TempDir()
	for directory := 0; directory < 100; directory++ {
		path := filepath.Join(root, fmt.Sprintf("pkg-%03d", directory))
		if err := os.MkdirAll(path, 0755); err != nil {
			b.Fatal(err)
		}
		for file := 0; file < 10; file++ {
			if err := os.WriteFile(filepath.Join(path, fmt.Sprintf("file-%02d.go", file)), []byte("package fixture\n"), 0644); err != nil {
				b.Fatal(err)
			}
		}
	}

	for _, workers := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("workers-%d", workers), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				files, err := observeDirectoryTreeWithWorkers(root, nil, nil, func(string) bool { return true }, workers)
				if err != nil || len(files) != 1000 {
					b.Fatalf("observed %d files: %v", len(files), err)
				}
			}
		})
	}
}
