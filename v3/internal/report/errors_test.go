package report

import (
	"errors"
	"fmt"
	"testing"
)

func TestReportedErrorsKeepIdentityAndRenderOnce(t *testing.T) {
	cause := errors.New("build failed")
	marked := MarkReported(cause)
	if !errors.Is(marked, cause) || !IsReported(fmt.Errorf("dev: %w", marked)) {
		t.Fatal("lost cause or rendered marker")
	}
	if MarkReported(marked) != marked {
		t.Fatal("wrapped reported error twice")
	}
	if MarkReported(nil) != nil || IsReported(cause) {
		t.Fatal("marked an unreported error")
	}
}
