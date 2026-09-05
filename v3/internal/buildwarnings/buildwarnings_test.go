package buildwarnings

import (
	"os"
	"testing"
)

func TestAddAndFlush(t *testing.T) {
	f, err := os.CreateTemp("", "wails-bw-test-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Setenv(EnvVar, f.Name())

	Add("tool has-cc", "deprecated message")
	Add("other-tool", "another warning")

	entries := read(f.Name())
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].source != "tool has-cc" || entries[0].message != "deprecated message" {
		t.Errorf("unexpected entry[0]: %+v", entries[0])
	}
	if entries[1].source != "other-tool" || entries[1].message != "another warning" {
		t.Errorf("unexpected entry[1]: %+v", entries[1])
	}

	// FlushAndPrint removes the file.
	FlushAndPrint()
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Error("expected warnings file to be removed after flush")
	}
}

func TestAddNoopWhenEnvUnset(t *testing.T) {
	t.Setenv(EnvVar, "") // ensure unset
	// Should not panic or create any file.
	Add("src", "msg")
}

func TestFlushNoopWhenEmpty(t *testing.T) {
	f, err := os.CreateTemp("", "wails-bw-test-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Setenv(EnvVar, f.Name())

	// No warnings added; FlushAndPrint should remove the file silently.
	FlushAndPrint()
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Error("expected empty warnings file to be removed after flush")
	}
}
