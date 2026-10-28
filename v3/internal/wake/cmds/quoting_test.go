package cmds

import (
	"reflect"
	"testing"
)

func TestShellSplitBackslashes(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []string
	}{
		{`go test -run "^Test\w+$"`, []string{"go", "test", "-run", `^Test\w+$`}},
		{`go build -o "C:\apps\app.exe"`, []string{"go", "build", "-o", `C:\apps\app.exe`}},
		{`go test -run "Test\"Name"`, []string{"go", "test", "-run", `Test"Name`}},
		{`go test -run "Test\\w"`, []string{"go", "test", "-run", `Test\w`}},
		{"go test -run \"Test\\\nName\"", []string{"go", "test", "-run", "TestName"}},
		{"go\\\n test", []string{"go", "test"}},
		{`go build -o app\ name`, []string{"go", "build", "-o", "app name"}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			if got := shellSplit(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokens = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRoutePreservesTestPattern(t *testing.T) {
	command, ok := Route(`go test -run "^Test\w+$"`, RouteOptions{}).(*GoTestCmd)
	if !ok {
		t.Fatal("expected native Go test command")
	}
	if command.Opts.Run != `^Test\w+$` {
		t.Errorf("test pattern = %q, want original regex", command.Opts.Run)
	}
}
