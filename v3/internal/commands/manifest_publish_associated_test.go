package commands

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestAssociatedPublicationRollsBackBothFilesAfterSecondRenameFails(t *testing.T) {
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "second-install-failure"}[failure], func(t *testing.T) {
			root := t.TempDir()
			for name, data := range map[string]string{".wails/new/package": "new archive", ".wails/new/package.sig": "new signature", "dist/package": "old archive", "dist/package.sig": "old signature"} {
				path := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(data), 0640); err != nil {
					t.Fatal(err)
				}
			}
			spec := pipeline.PublishSpec{Source: ".wails/new/package", Destination: "dist/package", Associated: []pipeline.PublishedFile{{Source: ".wails/new/package.sig", Destination: "dist/package.sig"}}}
			failed := false
			_, err := (&manifestHandler{root: root}).publishAssociated(spec, func(from, to string) error {
				if failure && !failed && filepath.Base(from) == "file-1" {
					failed = true
					return errors.New("injected second publication failure")
				}
				return os.Rename(from, to)
			})
			if (err != nil) != failure || failed != failure {
				t.Fatalf("publication failure=%v, injected=%v, error=%v", failure, failed, err)
			}
			prefix := "new "
			if failure {
				prefix = "old "
			}
			for file, contents := range map[string]string{"package": prefix + "archive", "package.sig": prefix + "signature"} {
				path := filepath.Join(root, "dist", file)
				data, err := os.ReadFile(path)
				if err != nil || string(data) != contents {
					t.Fatalf("%s: %q, %v", file, data, err)
				}
				info, err := os.Stat(path)
				if err != nil || (runtime.GOOS != "windows" && info.Mode().Perm() != 0640) {
					t.Fatalf("%s permissions: %v, %v", file, info, err)
				}
			}
			leftovers, err := filepath.Glob(filepath.Join(root, "dist", ".wails-publish-*"))
			if err != nil || len(leftovers) != 0 {
				t.Fatalf("publication transaction was not cleaned up: %v, %v", leftovers, err)
			}
		})
	}
}
