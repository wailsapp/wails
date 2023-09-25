package main

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func main() {
	// Create a renderer to output data
	out, err := gomarkdoc.NewRenderer()
	if err != nil {
		// handle error
	}

	wd, err := os.Getwd()
	if err != nil {
		// handle error
	}

	packagePath := filepath.Join(wd, "../../v3/pkg/application")

	buildPkg, err := build.ImportDir(packagePath, build.ImportComment)
	if err != nil {
		// handle error
	}

	// Create a documentation package from the build representation of our
	// package.
	log := logger.New(logger.DebugLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		// handle error
		panic(err)
	}

	// Write the documentation out to console.
	data, err := out.Package(pkg)
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(cwd)

	err = os.WriteFile(filepath.Join("..", "docs", "API", "fullapi.md"), []byte(data), 0644)
	if err != nil {
		panic(err)
	}
}
