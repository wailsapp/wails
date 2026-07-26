package templates

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestList(t *testing.T) {

	is2 := is.New(t)
	templateList, err := List()
	is2.NoErr(err)

	is2.Equal(len(templateList), 13)
}

func TestShortname(t *testing.T) {

	is2 := is.New(t)

	vanillaTemplate, err := getTemplateByShortname("vanilla")
	is2.NoErr(err)

	is2.Equal(vanillaTemplate.Name, "Vanilla + Vite")
}

func TestInstall(t *testing.T) {

	is2 := is.New(t)

	// Change to the directory of this file
	_, filename, _, _ := runtime.Caller(0)

	err := os.Chdir(filepath.Dir(filename))
	is2.NoErr(err)

	options := &Options{
		ProjectName:  "test",
		TemplateName: "vanilla",
		AuthorName:   "Lea Anthony",
		AuthorEmail:  "lea.anthony@gmail.com",
	}

	defer func() {
		_ = os.RemoveAll(options.ProjectName)
	}()
	_, _, err = Install(options)
	is2.NoErr(err)

}

func TestVueTSTemplateAppScriptBlock(t *testing.T) {
	paths := []string{
		"templates/vue-ts/frontend/src/App.vue",
		"generate/assets/vue-ts/frontend/src/App.vue",
	}
	const scriptBlock = `<script lang="ts" setup>
import HelloWorld from './components/HelloWorld.vue'
</script>`

	var canonical string
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}

		if !strings.Contains(string(contents), scriptBlock) {
			t.Errorf("%s does not contain a valid App script block", path)
		}
		if canonical == "" {
			canonical = string(contents)
		} else if string(contents) != canonical {
			t.Errorf("%s has drifted from %s", path, paths[0])
		}
	}
}

func TestInstallFailsInNonEmptyDirectory(t *testing.T) {
	is2 := is.New(t)

	// Create a temp directory with a file in it
	tempDir, err := os.MkdirTemp("", "wails-test-nonempty-*")
	is2.NoErr(err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create a file in the directory to make it non-empty
	err = os.WriteFile(filepath.Join(tempDir, "existing-file.txt"), []byte("test"), 0644)
	is2.NoErr(err)

	options := &Options{
		ProjectName:  "test",
		TemplateName: "vanilla",
		TargetDir:    tempDir,
	}

	_, _, err = Install(options)
	is2.True(err != nil) // Should fail
	is2.True(err.Error() == "cannot initialise project in non-empty directory: "+tempDir)
}

func TestInstallSucceedsInEmptyDirectory(t *testing.T) {
	is2 := is.New(t)

	// Create an empty temp directory
	tempDir, err := os.MkdirTemp("", "wails-test-empty-*")
	is2.NoErr(err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	options := &Options{
		ProjectName:  "test",
		TemplateName: "vanilla",
		TargetDir:    tempDir,
	}

	_, _, err = Install(options)
	is2.NoErr(err) // Should succeed in empty directory
}
