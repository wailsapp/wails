package buildmeta

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSourceAndStrictSemanticVersions(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "version.go")
	input := "package app\nconst core = `v1.2.3`\nvar AppVersion = (core + \"-alpha-12\") + \"+build.5\"\nfunc init(){ panic(\"must not execute\") }\n"
	if err := os.WriteFile(path, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	info, err := Source(path, "AppVersion")
	if err != nil || info.Version != "1.2.3-alpha-12+build.5" || info.Core != "1.2.3" || !info.Prerelease() {
		t.Fatalf("info=%+v err=%v", info, err)
	}
	for _, value := range []string{"1.2", "1", "01.2.3", "1.2.3-01", "1.2.3;echo", "1.2.3\n"} {
		if _, err := Parse(value); err == nil {
			t.Errorf("accepted %q", value)
		}
	}
	info, err = Parse("1.2.3+build-with-dash")
	if err != nil || info.Prerelease() {
		t.Fatalf("build metadata became prerelease: %+v %v", info, err)
	}
	for _, input := range []string{"var AppVersion = compute()", "var AppVersion = Missing", "var AppVersion = Other\nvar Other = AppVersion", "const AppVersion = 123"} {
		if err := os.WriteFile(path, []byte("package app\n"+input), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := Source(path, "AppVersion"); err == nil {
			t.Errorf("accepted non-static source %q", input)
		}
	}
}

func TestWindowsVersionBounds(t *testing.T) {
	info, _ := Parse("65535.2.3-beta.2")
	got, err := info.Windows()
	if err != nil || got != [4]uint16{65535, 2, 3, 0} {
		t.Fatalf("numeric=%v %v", got, err)
	}
	info, _ = Parse("65536.2.3")
	if _, err := info.Windows(); err == nil {
		t.Fatal("Windows version silently overflowed")
	}
}

func TestRevisionWorktreesPackedAndDetachedHEAD(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("Git fixture tool unavailable")
	}
	root := t.TempDir()
	run := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command(git, args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.test", "GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.test")
		data, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s %v", args, data, err)
		}
		return strings.TrimSpace(string(data))
	}
	run(root, "init", "-q", "-b", "main")
	run(root, "commit", "-qm", "initial", "--allow-empty")
	want := run(root, "rev-parse", "HEAD")
	work := filepath.Join(t.TempDir(), "work tree")
	run(root, "worktree", "add", "-qb", "work", work)
	run(root, "pack-refs", "--all", "--prune")
	for _, path := range []string{root, work} {
		got, err := Revision(context.Background(), path)
		if err != nil || got != want {
			t.Fatalf("revision in %s=%q want=%s err=%v", path, got, want, err)
		}
	}
	run(work, "checkout", "--detach", "-q", want)
	got, err := Revision(context.Background(), work)
	if err != nil || got != want {
		t.Fatalf("detached=%q %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Revision(ctx, root); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func BenchmarkSourceVersion(b *testing.B) {
	path := filepath.Join(b.TempDir(), "version.go")
	if err := os.WriteFile(path, []byte("package app\nconst Core = `1.2.3`\nvar AppVersion = Core + `-alpha.12`\n"), 0600); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Source(path, "AppVersion"); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRevisionRejectsUnsupportedRepositoryFormats(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("Git fixture tool unavailable")
	}
	for _, option := range []string{"--object-format=sha256", "--ref-format=reftable"} {
		t.Run(option, func(t *testing.T) {
			root := t.TempDir()
			cmd := exec.Command(git, "init", "-q", option, root)
			cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Skipf("Git fixture format unavailable: %s", output)
			}
			commit := exec.Command(git, "-C", root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "--allow-empty", "-qm", "initial")
			commit.Env = cmd.Env
			if output, err := commit.CombinedOutput(); err != nil {
				t.Fatalf("commit fixture: %s %v", output, err)
			}
			if value, err := Revision(context.Background(), root); err == nil {
				t.Fatalf("unsupported repository returned %q", value)
			}
		})
	}
}
