package commands

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

// BenchmarkManifestLinuxPackageParallelism measures complete format adapters,
// including staging and publication, with identical synthetic binary payloads.
// Half the payload is deterministic entropy and half is zero-filled so the
// compression workload includes both incompressible and compressible data.
func BenchmarkManifestLinuxPackageParallelism(b *testing.B) {
	for _, workers := range []int{1, 4} {
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			handler, base := linuxPackageFixture(b)
			base.Project.Description = "Linux package performance fixture"
			payload := make([]byte, 4<<20)
			random := rand.New(rand.NewPCG(1, 2))
			for offset := 0; offset < len(payload)/2; offset += 8 {
				binary.LittleEndian.PutUint64(payload[offset:], random.Uint64())
			}
			for _, name := range []string{"example", "helper"} {
				if err := os.WriteFile(filepath.Join(handler.root, "bin", name), payload, 0755); err != nil {
					b.Fatal(err)
				}
			}
			var specs []pipeline.PackageSpec
			for _, format := range []string{"deb", "rpm", "archlinux", "void"} {
				spec := base
				spec.Format = format
				spec.Output = filepath.Join("dist", "example."+format)
				specs = append(specs, spec)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(payload) * 2 * len(specs)))
			b.ResetTimer()
			for range b.N {
				jobs := make(chan pipeline.PackageSpec, len(specs))
				failures := make(chan error, len(specs))
				for _, spec := range specs {
					jobs <- spec
				}
				close(jobs)
				var group sync.WaitGroup
				for range workers {
					group.Go(func() {
						for spec := range jobs {
							var err error
							if spec.Format == "void" {
								_, err = handler.packageVoid(context.Background(), spec)
							} else {
								_, err = handler.packageLinux(spec)
							}
							if err != nil {
								failures <- fmt.Errorf("%s: %w", spec.Format, err)
							}
						}
					})
				}
				group.Wait()
				close(failures)
				for err := range failures {
					b.Fatal(err)
				}
				for _, spec := range specs {
					info, err := os.Stat(filepath.Join(handler.root, spec.Output))
					if err != nil || info.Size() == 0 {
						b.Fatalf("%s output missing or empty: %v", spec.Format, err)
					}
				}
			}
		})
	}
}
