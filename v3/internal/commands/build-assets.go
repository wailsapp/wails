package commands

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/leaanthony/gosod"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed build_assets
var buildAssets embed.FS

type BuildAssetsOptions struct {
	Dir                string `description:"The directory to generate the files into" default:"build"`
	Name               string `description:"The name of the project"`
	ProductName        string `description:"The name of the product" default:"My Product"`
	ProductDescription string `description:"The description of the product" default:"My Product Description"`
	ProductVersion     string `description:"The version of the product" default:"0.1.0"`
	ProductCompany     string `description:"The company of the product" default:"My Company"`
	ProductCopyright   string `description:"The copyright notice"`
	ProductComments    string `description:"Comments to add to the generated files" default:"This is a comment"`
	ProductIdentifier  string `description:"The product identifier, e.g com.mycompany.myproduct"`
	Silent             bool   `description:"Suppress output to console"`
}

func GenerateBuildAssets(options *BuildAssetsOptions) error {

	var err error
	options.Dir, err = filepath.Abs(options.Dir)
	if err != nil {
		return err
	}

	// If directory doesn't exist, create it
	if _, err := os.Stat(options.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(options.Dir, 0755)
		if err != nil {
			return err
		}
	}

	if options.ProductComments == "" {
		options.ProductComments = fmt.Sprintf("(c) %d %s", time.Now().Year(), options.ProductCompany)
	}

	if options.ProductIdentifier == "" {
		options.ProductIdentifier = "com.wails." + normaliseName(options.Name)
	}

	tfs, err := fs.Sub(buildAssets, "build_assets")
	if err != nil {
		return err
	}

	if !options.Silent {
		println("Generating build assets in " + options.Dir)
	}
	return gosod.New(tfs).Extract(options.Dir, options)

}

func normaliseName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
