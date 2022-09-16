package templates

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestList(t *testing.T) {

	is2 := is.New(t)
	templates, err := List()
	is2.NoErr(err)

	println("Found these templates:")
	for _, template := range templates {
		fmt.Printf("%+v\n", template)
	}
}

func TestShortname(t *testing.T) {

	is2 := is.New(t)

	template, err := getTemplateByShortname("vanilla")
	is2.NoErr(err)

	println("Found this template:")
	fmt.Printf("%+v\n", template)
}

func TestInstall(t *testing.T) {

	is2 := is.New(t)

	options := &Options{
		ProjectName:  "test",
		TemplateName: "vanilla",
		AuthorName:   "Lea Anthony",
		AuthorEmail:  "lea.anthony@gmail.com",
	}

	_, _, err := Install(options)
	is2.NoErr(err)
}
