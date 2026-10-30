package cmds

import "testing"

func TestGoRunPreservesApplicationFlags(t *testing.T) {
	_, output, err := runOfflineGo(t, "go run -buildvcs=false ./app -tags runtime-tag -ldflags -s -trimpath", map[string]string{
		"go.mod":      "module example.com/wake-run-check\n\ngo 1.25.0\n",
		"app/main.go": "package main\n\nimport (\"fmt\"; \"os\")\n\nfunc main() { for _, arg := range os.Args[1:] { fmt.Println(arg) } }\n",
	})
	if err != nil {
		t.Fatalf("offline run failed: %v; output: %s", err, output)
	}
	if want := "-tags\nruntime-tag\n-ldflags\n-s\n-trimpath\n"; output != want {
		t.Errorf("application arguments = %q, want %q", output, want)
	}
}
