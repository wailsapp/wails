package main

import (
	"embed"
	"generate/s"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"os"
	"strings"
)

//go:embed assets/svelte/*
var svelte embed.FS

//go:embed assets/lit/*
var lit embed.FS

//go:embed assets/vue/*
var vue embed.FS

//go:embed assets/react/*
var react embed.FS

//go:embed assets/preact/*
var preact embed.FS

func checkError(err error) {
	if err != nil {
		println("\nERROR:", err.Error())
		os.Exit(1)
	}
}

type template struct {
	Name          string
	ShortName     string
	Description   string
	Assets        embed.FS
	FilesToDelete []string
	DirsToDelete  []string
}

var templates = []*template{
	{
		Name:          "Svelte + Vite",
		ShortName:     "Svelte",
		Description:   "Svelte + Vite development server",
		Assets:        svelte,
		FilesToDelete: []string{"frontend/index.html", "frontend/.gitignore", "frontend/src/assets/svelte.png"},
		DirsToDelete:  []string{"frontend/public", "frontend/src/lib"},
	},
	{
		Name:          "Lit + Vite",
		ShortName:     "Lit",
		Description:   "Lit + Vite development server",
		Assets:        lit,
		FilesToDelete: []string{"frontend/index.html", "frontend/vite.config.js"},
	},
	{
		Name:          "Vue + Vite",
		ShortName:     "Vue",
		Description:   "Vue + Vite development server",
		Assets:        vue,
		FilesToDelete: []string{"frontend/index.html", "frontend/.gitignore"},
		DirsToDelete:  []string{"frontend/src/assets", "frontend/src/components", "frontend/public"},
	},
	{
		Name:          "React + Vite",
		ShortName:     "React",
		Description:   "React + Vite development server",
		Assets:        react,
		FilesToDelete: []string{"frontend/src/index.css", "frontend/src/favicon.svg", "frontend/src/logo.svg", "frontend/.gitignore", "frontend/index.html"},
	},
	{
		Name:          "Preact + Vite",
		ShortName:     "Preact",
		Description:   "Preact + Vite development server",
		Assets:        preact,
		FilesToDelete: []string{"frontend/src/index.css", "frontend/src/favicon.svg", "frontend/src/logo.jsx", "frontend/.gitignore", "frontend/index.html"},
	},
}

func main() {

	for _, t := range templates {
		createTemplate(t)
	}
}

func createTemplate(template *template) {
	cwd := s.CWD()
	name := template.Name
	shortName := strings.ToLower(template.ShortName)
	assets, err := debme.FS(template.Assets, "assets/"+shortName)
	checkError(err)

	s.CD("..")
	s.ENDIR("testtemplates")
	s.CD("testtemplates")
	s.RMDIR(shortName)
	s.COPYDIR("../base", shortName)
	s.CD(shortName)
	s.ECHO("Generating vite template: " + shortName)
	s.EXEC("npm create vite@latest frontend --template " + shortName)

	// Clean up template
	for _, fileToDelete := range template.FilesToDelete {
		s.DELETE(fileToDelete)
	}
	for _, dirToDelete := range template.DirsToDelete {
		s.RMDIR(dirToDelete)
	}
	s.REPLACEALL("README.md", s.Sub{"$NAME": template.ShortName})
	s.REPLACEALL("template.json", s.Sub{"$NAME": name, "$SHORTNAME": shortName, "$DESCRIPTION": template.Description})

	// Add custom files
	g := gosod.New(assets)
	g.SetTemplateFilters([]string{})
	err = g.Extract(".", nil)
	checkError(err)

	// Do frontend
	s.CD("frontend")
	s.MKDIR("dist")
	s.TOUCH("dist/.gitignore")
	s.CD(cwd)
}
