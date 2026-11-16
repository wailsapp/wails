package buildmeta

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureRevision = "1234567890abcdef1234567890abcdef12345678"

func revisionFixture(t testing.TB, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, data := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRevisionLooseRefsAndRelativeGitfile(t *testing.T) {
	root := revisionFixture(t, map[string]string{
		".git": "gitdir: metadata\n", "metadata/config": "[core]\nrepositoryformatversion = 0\n",
		"metadata/HEAD": "ref: refs/heads/main\n", "metadata/refs/heads/main": fixtureRevision + "\n",
		"metadata/packed-refs": strings.Repeat("a", 40) + " refs/heads/main\n",
		"subdir/file":          "",
	})
	got, err := Revision(context.Background(), filepath.Join(root, "subdir"))
	if err != nil || got != fixtureRevision {
		t.Fatalf("revision=%q err=%v", got, err)
	}
}

func TestRevisionRejectsBrokenMetadata(t *testing.T) {
	for name, head := range map[string]string{
		"empty": "", "zero": strings.Repeat("0", 40), "nonhex": strings.Repeat("g", 40), "short": "1234",
		"traversal": "ref: refs/../secret", "absolute": "ref: /etc/passwd", "cycle": "ref: refs/heads/main",
		"control": "ref: refs/heads/bad\tname", "unborn": "ref: refs/heads/missing",
	} {
		t.Run(name, func(t *testing.T) {
			root := revisionFixture(t, map[string]string{".git/config": "[core]\nrepositoryformatversion = 0\n", ".git/HEAD": head, ".git/refs/heads/main": "ref: refs/heads/main"})
			if got, err := Revision(context.Background(), root); err == nil {
				t.Fatalf("accepted %q", got)
			}
		})
	}
}

func TestRevisionNeverUsesOuterRepositoryForBrokenInnerMarker(t *testing.T) {
	root := revisionFixture(t, map[string]string{".git/config": "[core]\nrepositoryformatversion = 0\n", ".git/HEAD": fixtureRevision, "nested/.git": "broken marker"})
	if got, err := Revision(context.Background(), filepath.Join(root, "nested")); err == nil {
		t.Fatalf("used outer revision %q", got)
	}
}

func TestRevisionFormatValues(t *testing.T) {
	for _, value := range []string{"objectformat = sha256", "refStorage = reftable", "objectformat =", "refStorage", `objectformat = "sha1#bad"`, `refStorage = "files;bad"`, `objectformat = "sha1 "`, `objectformat = "sha1`, `objectformat = "sha1\q"`} {
		t.Run(value, func(t *testing.T) {
			root := revisionFixture(t, map[string]string{".git/config": "[extensions]\n" + value + "\n", ".git/HEAD": fixtureRevision})
			if got, err := Revision(context.Background(), root); err == nil {
				t.Fatalf("accepted unsupported format: %s", got)
			}
		})
	}
}

func BenchmarkRevision(b *testing.B) {
	root := revisionFixture(b, map[string]string{".git/config": "[core]\nrepositoryformatversion = 0\n", ".git/HEAD": "ref: refs/heads/main\n", ".git/refs/heads/main": fixtureRevision})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Revision(context.Background(), root); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRevisionAcceptsQuotedFormatAndTrailingComments(t *testing.T) {
	for _, config := range []string{
		"[extensions] # formats\nobjectformat = \"sha1\" # trailing comment\nrefstorage = \"files\" ; comment\n",
		"[EXTENSIONS]\nobjectformat = \"sh\"a\"1\"\nrefstorage = files\n",
	} {
		root := revisionFixture(t, map[string]string{".git/config": config, ".git/HEAD": fixtureRevision})
		got, err := Revision(context.Background(), root)
		if err != nil || got != fixtureRevision {
			t.Fatalf("revision=%q err=%v config=%q", got, err, config)
		}
	}
}
