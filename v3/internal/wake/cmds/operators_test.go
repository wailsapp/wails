package cmds

import "testing"

func TestRouteShellOperators(t *testing.T) {
	for _, command := range []string{
		"go build && echo done", "go build || echo failed", "go build|cat",
		"go test ./... > result.txt", "go test ./... 2>&1", "go test < input.txt",
		"go build; echo done", "go build & echo done", "go build\necho done",
		"go build (echo done)",
	} {
		t.Run(command, func(t *testing.T) {
			routed := Route(command, RouteOptions{})
			shell, ok := routed.(*ShellCmd)
			if !ok {
				t.Fatalf("routed to %T, want shell fallback", routed)
			}
			if shell.Cmd != command {
				t.Errorf("shell command = %q, want %q", shell.Cmd, command)
			}
		})
	}
	for _, command := range []string{`go test -run "Test(A|B)"`, `go test -run TestA\|TestB`} {
		routed := Route(command, RouteOptions{})
		if _, ok := routed.(*GoTestCmd); !ok {
			t.Errorf("literal operator routed to %T, want native Go test", routed)
		}
	}
}
