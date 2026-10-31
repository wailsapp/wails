//go:build !windows

package cmds

import (
	"bytes"
	"testing"
)

func TestProgrammaticGoTestOptions(t *testing.T) {
	base := installFakeGo(t, `printf '%s\n' "$@"`)
	cmd := GoTest(GoTestOptions{
		Package: "./pkg", Verbose: true, Cover: true, Race: true,
		Tags: []string{"gtk4", "production"}, Count: 10, Run: "^TestFoo$",
		Extra: []string{"-timeout", "5s"},
	})
	cmd.Dir = base
	var output bytes.Buffer
	cmd.SetOutput(&output, &output)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	want := "test\n-v\n-cover\n-race\n-tags\ngtk4 production\n-count\n10\n-run\n^TestFoo$\n-timeout\n5s\n./pkg\n"
	if got := output.String(); got != want {
		t.Errorf("programmatic arguments = %q, want %q", got, want)
	}
}
