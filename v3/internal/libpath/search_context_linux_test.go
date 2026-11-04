//go:build linux

package libpath

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFindLibraryPathContextErrors(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, expire := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer expire()
	for _, tt := range []struct {
		name string
		ctx  context.Context
		want error
	}{
		{"cancelled", cancelled, context.Canceled},
		{"expired", expired, context.DeadlineExceeded},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path, err := findLibraryPathCtx(tt.ctx, "wails-test-cancelled-library-no-match")
			if path != "" || !errors.Is(err, tt.want) {
				t.Errorf("search = %q, %v; want empty path and %v", path, err, tt.want)
			}
		})
	}
}
