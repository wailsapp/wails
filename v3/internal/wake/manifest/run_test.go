package manifest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunDefaultsPlatformOverridesAndEjection(t *testing.T) {
	root := t.TempDir()
	source := `version = 3
project {
 name = "example"
 product_name = "Example"
 identifier = "com.example.test"
 version = "1.0.0"
}
frontend { disabled = true }
run {
 tags = ["shared", "duplicate"]
 args = ["default"]
 environment = { SHARED = "yes", OVERRIDE = "shared" }
}
target "linux/amd64" {
 environment = { ARCH = "yes", OVERRIDE = "arch" }
 tags = ["arch"]
 run {
  tags = ["arch-run"]
  args = []
  environment = { OVERRIDE = "arch" }
 }
}
target "linux" {
 environment = { COMMON = "yes", OVERRIDE = "platform" }
 tags = ["platform"]
 run {
  tags = ["duplicate", "platform-run"]
  args = ["platform"]
  environment = { PLATFORM = "yes", OVERRIDE = "platform" }
 }
}
`
	loaded, err := decodeHCL(root, filepath.Join(root, Filename), []byte(source), "")
	require.NoError(t, err)
	require.True(t, loaded.Config.Frontend.Disabled)
	check := func(c Config) {
		run, err := c.RunForTarget("linux", "amd64")
		require.NoError(t, err)
		require.Empty(t, run.Args)
		require.True(t, run.ArgsSet)
		require.Equal(t, []string{"shared", "duplicate", "platform-run", "arch-run"}, run.Tags)
		require.Equal(t, map[string]string{"SHARED": "yes", "PLATFORM": "yes", "OVERRIDE": "arch"}, run.Environment)
		require.Equal(t, []string{"platform", "arch"}, c.Targets.Linux.AMD64.Tags)
		require.Equal(t, map[string]string{"COMMON": "yes", "ARCH": "yes", "OVERRIDE": "arch"}, c.Targets.Linux.AMD64.Environment)
		other, err := c.RunForTarget("darwin", "arm64")
		require.NoError(t, err)
		require.Equal(t, []string{"default"}, other.Args)
		require.Equal(t, []string{"shared", "duplicate"}, other.Tags)
		run.Environment["SHARED"] = "changed"
		require.Equal(t, "yes", c.Run.Environment["SHARED"])
	}
	check(loaded.Config)
	ejected, err := EncodeEjectedHCL(loaded.Config, "test")
	require.NoError(t, err)
	restored, err := decodeHCL(root, filepath.Join(root, EjectedFilename), ejected, "")
	require.NoError(t, err)
	check(restored.Config)
}

func TestRunManifestDiscoveryUsesExampleLocalRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/test\n"), 0600))
	example := filepath.Join(root, "examples", "test")
	require.NoError(t, os.MkdirAll(filepath.Join(example, "child"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(example, Filename), []byte("invalid but present"), 0600))
	found, _, err := Discover(filepath.Join(example, "child"))
	require.NoError(t, err)
	require.Equal(t, example, found)
	_, err = Load(example, "")
	require.Error(t, err)
}

// Discover entry points independently of manifests so a forgotten example
// cannot silently disappear from the launch smoke suite.
func TestExampleRunManifestsCoverRunnableDirectories(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "examples"))
	require.NoError(t, err)
	entries := map[string]bool{}
	manifests := map[string]bool{}
	require.NoError(t, filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "node_modules", ".wails", ".git", "bin", "frontend":
				return filepath.SkipDir
			case "build":
				if filepath.Dir(path) != root {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if entry.Name() == Filename {
			manifests[filepath.Dir(path)] = true
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		if file.Name.Name != "main" {
			return nil
		}
		for _, decl := range file.Decls {
			if function, ok := decl.(*ast.FuncDecl); ok && function.Name.Name == "main" && function.Recv == nil {
				entries[filepath.Dir(path)] = true
			}
		}
		return nil
	}))
	require.NotEmpty(t, entries)
	require.Equal(t, entries, manifests, "each runnable example needs exactly one local manifest")
	identifiers := map[string]string{}
	for directory := range entries {
		loaded, err := Load(directory, "")
		require.NoError(t, err, directory)
		identifier := loaded.Config.Project.Identifier
		previous, duplicate := identifiers[identifier]
		require.False(t, duplicate, "%s and %s share application identifier %s", previous, directory, identifier)
		identifiers[identifier] = directory
	}
	t.Logf("Validated %d independently discovered example entry directories", len(entries))
}
