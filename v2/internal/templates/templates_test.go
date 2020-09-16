package templates

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestList(t *testing.T) {

	is := is.New(t)
	templates, err := List()
	is.Equal(err, nil)

	println("Found these templates:")
	for _, template := range templates {
		fmt.Printf("%+v\n", template)
	}
}
