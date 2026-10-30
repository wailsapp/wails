//go:build !windows

package cmds

import (
	"bytes"
	"testing"
)

func TestInvalidTestCountReachesGoTool(t *testing.T) {
	base := installFakeGo(t, `printf '%s\n' "$@"`)
	command := Route("go test -count not-a-number ./...", RouteOptions{Dir: base})
	var output bytes.Buffer
	command.(OutputSetter).SetOutput(&output, &output)
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "test\n-count\nnot-a-number\n./...\n" {
		t.Errorf("forwarded Go arguments = %q, want invalid count retained", got)
	}
}
