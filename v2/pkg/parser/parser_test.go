package parser

import (
	"testing"

	"github.com/leaanthony/slicer"
	"github.com/matryer/is"
	"github.com/wailsapp/wails/v2/internal/fs"
)

func TestParser(t *testing.T) {

	is := is.New(t)

	// Local project dir
	projectDir := fs.RelativePath("./testproject")

	p := NewParser()

	// Check parsing worked
	err := p.ParseProject(projectDir)
	is.NoErr(err)

	// Expected structs
	expectedBoundStructs := slicer.String()
	expectedBoundStructs.Add("main.Basic", "mypackage.Manager")

	// We expect these to be the same length
	is.Equal(expectedBoundStructs.Length(), len(p.BoundStructs))

	// Check bound structs
	for _, boundStruct := range p.BoundStructs {

		// Check the names are correct
		fqn := boundStruct.FullyQualifiedName()
		is.True(expectedBoundStructs.Contains(fqn))

		// Check that the structs have comments
		is.True(len(boundStruct.Comments) > 0)

		// Check that the structs have methods
		is.True(len(boundStruct.Methods) > 0)

	}

}
