package report

import (
	"sync"
	"testing"
)

func TestActiveReporterConcurrentAccess(t *testing.T) {
	defer SetActive(nil)
	reporters := []Reporter{Nop{}, &Nop{}, nil}
	var workers sync.WaitGroup
	for i := 0; i < 8; i++ {
		workers.Go(func() {
			for j := 0; j < 1000; j++ {
				SetActive(reporters[(i+j)%len(reporters)])
				if current := Active(); current == nil {
					t.Error("Active returned a nil reporter")
				} else {
					_ = current.Level()
				}
			}
		})
	}
	workers.Wait()
}
