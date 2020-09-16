package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestRelativePath(t *testing.T) {

	is := is.New(t)

	cwd, err := os.Getwd()
	is.Equal(err, nil)

	// Check current directory
	actual := RelativePath(".")
	is.Equal(actual, cwd)

	// Check 2 parameters
	actual = RelativePath("..", "fs")
	is.Equal(actual, cwd)

	// Check 3 parameters including filename
	actual = RelativePath("..", "fs", "fs.go")
	expected := filepath.Join(cwd, "fs.go")
	is.Equal(actual, expected)

}
