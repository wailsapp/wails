package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestValidateRemoteURL_Valid(t *testing.T) {
	for _, u := range []string{
		"https://github.com/user/repo",
		"https://github.com/user/repo.git",
		"HTTPS://github.com/user/repo",
		"http://example.com/user/repo.git",
		"ssh://git@github.com/user/repo.git",
		"ssh://git@github.com:2222/user/repo.git",
		"git://example.com/user/repo.git",
		"git@github.com:user/repo",
		"git@github.com:user/repo.git",
		"github.com:user/repo",
	} {
		if err := ValidateRemoteURL(u); err != nil {
			t.Errorf("ValidateRemoteURL(%q) = %v, want nil", u, err)
		}
	}
}

func TestValidateRemoteURL_Invalid(t *testing.T) {
	for _, u := range []string{
		"",
		"--help",
		"-uhelp",
		"--upload-pack=touch /tmp/x",
		"ext::sh -c touch /tmp/x",
		"ext::sh",
		"fd::3",
		"fd::3,4",
		"EXT::sh",
		"file:///etc",
		"FILE:///etc",
		"ftp://example.com/repo",
		"github.com/user/repo",
		"/tmp/repo",
		"./repo",
		"C:/repo",
		`C:\repo`,
		"ssh://-oProxyCommand=touch%20x/repo",
		"-oProxyCommand=x:repo",
		"git@-oProxyCommand=x:repo",
		"https:///repo",
		"https://github.com/user/repo\n--help",
		"https://github.com/user/repo --help",
	} {
		err := ValidateRemoteURL(u)
		if !errors.Is(err, ErrInvalidRemoteURL) {
			t.Errorf("ValidateRemoteURL(%q) = %v, want ErrInvalidRemoteURL", u, err)
		}
	}
}

func TestCloneRemote_RejectsBeforeRunningGit(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "marker")
	for _, u := range []string{
		"--help",
		"-uhelp",
		"--upload-pack=touch " + marker,
		"ext::sh -c touch% " + marker,
		"fd::3",
		"file:///etc",
	} {
		dst := filepath.Join(t.TempDir(), "clone")
		err := CloneRemote(u, dst, "")
		if !errors.Is(err, ErrInvalidRemoteURL) {
			t.Errorf("CloneRemote(%q) = %v, want ErrInvalidRemoteURL", u, err)
		}
		if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
			t.Errorf("CloneRemote(%q) created %s", u, dst)
		}
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("marker file was created: %v", err)
	}
}

func TestCloneRemote_RejectsOptionTag(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "clone")
	if err := CloneRemote("https://github.com/user/repo", dst, "--upload-pack=x"); err == nil {
		t.Fatal("expected error for tag starting with '-'")
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Fatalf("clone directory was created: %v", err)
	}
}

func TestClone_RejectsOptions(t *testing.T) {
	src := makeRepo(t)
	dst := filepath.Join(t.TempDir(), "clone")
	cases := []struct{ url, dir, tag string }{
		{"--help", dst, ""},
		{"-uhelp", dst, ""},
		{src, "--help", ""},
		{src, dst, "--upload-pack=x"},
	}
	for _, c := range cases {
		if err := Clone(c.url, c.dir, c.tag); err == nil {
			t.Errorf("Clone(%q, %q, %q) = nil, want error", c.url, c.dir, c.tag)
		}
	}
}

func TestCloneRemote_RefusesLocalTransport(t *testing.T) {
	// Even if validation were bypassed, git itself must refuse the file
	// transport under the protocol restrictions used by CloneRemote.
	src := makeRepo(t)
	dst := filepath.Join(t.TempDir(), "clone")
	args, err := cloneArgs("file://"+src, dst, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := run(append(append([]string{}, remoteProtocolArgs...), args...)...); err == nil {
		t.Fatal("expected git to refuse the file transport")
	}
}

func TestRemoteAdd_RejectsOptions(t *testing.T) {
	dir := t.TempDir()
	if err := Init(dir); err != nil {
		t.Fatal(err)
	}
	if err := RemoteAdd(dir, "origin", "--mirror=fetch"); err == nil {
		t.Fatal("expected error for url starting with '-'")
	}
	if err := RemoteAdd(dir, "-f", "https://example.com/repo.git"); err == nil {
		t.Fatal("expected error for name starting with '-'")
	}
}

func TestValidationErrorsOmitCredentials(t *testing.T) {
	const secret = "s3cr3t-token"
	for _, raw := range []string{
		"https://user:" + secret + "@-host/repo",
		"https://user:" + secret + "@host:bad/repo",
		"ftp://user:" + secret + "@host/repo",
	} {
		err := ValidateRemoteURL(raw)
		if err == nil {
			t.Fatalf("ValidateRemoteURL(%q) = nil, want error", raw)
		}
		if strings.Contains(err.Error(), secret) {
			t.Errorf("ValidateRemoteURL(%q) error leaks credentials: %v", raw, err)
		}
	}
	if _, err := cloneArgs("-https://user:"+secret+"@host/repo", "dir", ""); err == nil || strings.Contains(err.Error(), secret) {
		t.Errorf("cloneArgs error = %v, want an error without credentials", err)
	}
	if err := RemoteAdd(t.TempDir(), "origin", "-https://user:"+secret+"@host/repo"); err == nil || strings.Contains(err.Error(), secret) {
		t.Errorf("RemoteAdd error = %v, want an error without credentials", err)
	}
}
