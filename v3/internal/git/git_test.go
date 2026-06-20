package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// makeRepo creates a temp git repository with one commit and a v1.0.0 tag.
func makeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmds := [][]string{
		{"-C", dir, "init", "--quiet"},
		{"-C", dir, "-c", "user.email=t@t.com", "-c", "user.name=T", "commit", "--allow-empty", "--quiet", "-m", "init"},
		{"-C", dir, "tag", "v1.0.0"},
	}
	for _, args := range cmds {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("setup git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

func TestInit_Success(t *testing.T) {
	if err := Init(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}

func TestInit_Error(t *testing.T) {
	// git -C on a non-existent path fails
	if err := Init("/nonexistent_wails_test_path_xyz"); err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestRemoteAdd_Success(t *testing.T) {
	dir := t.TempDir()
	if err := Init(dir); err != nil {
		t.Fatal(err)
	}
	if err := RemoteAdd(dir, "origin", "https://example.com/repo.git"); err != nil {
		t.Fatal(err)
	}
}

func TestAddAll_Success(t *testing.T) {
	dir := t.TempDir()
	if err := Init(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AddAll(dir); err != nil {
		t.Fatal(err)
	}
}

func TestHeadHash_Success(t *testing.T) {
	src := makeRepo(t)
	hash, err := HeadHash(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 8 {
		t.Errorf("expected 8-char hash, got %q (len %d)", hash, len(hash))
	}
}

func TestHeadHash_Error(t *testing.T) {
	// not a git repository
	if _, err := HeadHash(t.TempDir()); err == nil {
		t.Fatal("expected error for non-repo dir")
	}
}

func TestClone_WithoutTag(t *testing.T) {
	src := makeRepo(t)
	dst := filepath.Join(t.TempDir(), "clone")
	if err := Clone(src, dst, ""); err != nil {
		t.Fatal(err)
	}
}

func TestClone_WithTag(t *testing.T) {
	src := makeRepo(t)
	dst := filepath.Join(t.TempDir(), "clone")
	if err := Clone(src, dst, "v1.0.0"); err != nil {
		t.Fatal(err)
	}
}

func TestRun_NotInstalled(t *testing.T) {
	t.Setenv("PATH", "/nonexistent_path_that_does_not_exist")
	err := Init(t.TempDir())
	if !errors.Is(err, ErrNotInstalled) {
		t.Fatalf("expected ErrNotInstalled, got %v", err)
	}
}

func TestOutput_NotInstalled(t *testing.T) {
	t.Setenv("PATH", "/nonexistent_path_that_does_not_exist")
	_, err := HeadHash(t.TempDir())
	if !errors.Is(err, ErrNotInstalled) {
		t.Fatalf("expected ErrNotInstalled, got %v", err)
	}
}
