package main

import (
	"embed"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/initialise/templates/generate/s"
	"os"
	"strings"
)

//go:embed assets/common/*
var common embed.FS

//go:embed assets/svelte/*
var svelte embed.FS

//go:embed assets/svelte-ts/*
var sveltets embed.FS

//go:embed assets/lit/*
var lit embed.FS

//go:embed assets/lit-ts/*
var litts embed.FS

//go:embed assets/vue/*
var vue embed.FS

//go:embed assets/vue-ts/*
var vuets embed.FS

//go:embed assets/react/*
var react embed.FS

//go:embed assets/react-ts/*
var reactts embed.FS

//go:embed assets/preact/*
var preact embed.FS

//go:embed assets/preact-ts/*
var preactts embed.FS

//go:embed assets/vanilla/*
var vanilla embed.FS

//go:embed assets/vanilla-ts/*
var vanillats embed.FS

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
		Name:          "Svelte + Vite (Typescript)",
		ShortName:     "Svelte-TS",
		Description:   "Svelte + TS + Vite development server",
		Assets:        sveltets,
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
		Name:          "Lit + Vite (Typescript)",
		ShortName:     "Lit-TS",
		Description:   "Lit + TS + Vite development server",
		Assets:        litts,
		FilesToDelete: []string{"frontend/index.html", "frontend/src/favicon.svg"},
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
		Name:          "Vue + Vite (Typescript)",
		ShortName:     "Vue-TS",
		Description:   "Vue + Vite development server",
		Assets:        vuets,
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
		Name:          "React + Vite (Typescript)",
		ShortName:     "React-TS",
		Description:   "React + Vite development server",
		Assets:        reactts,
		FilesToDelete: []string{"frontend/src/index.css", "frontend/src/favicon.svg", "frontend/src/logo.svg", "frontend/.gitignore", "frontend/index.html"},
	},
	{
		Name:          "Preact + Vite",
		ShortName:     "Preact",
		Description:   "Preact + Vite development server",
		Assets:        preact,
		FilesToDelete: []string{"frontend/src/index.css", "frontend/src/favicon.svg", "frontend/src/logo.jsx", "frontend/.gitignore", "frontend/index.html"},
	},
	{
		Name:          "Preact + Vite (Typescript)",
		ShortName:     "Preact-TS",
		Description:   "Preact + Vite development server",
		Assets:        preactts,
		FilesToDelete: []string{"frontend/src/index.css", "frontend/src/favicon.svg", "frontend/src/logo.tsx", "frontend/.gitignore", "frontend/index.html"},
	},
	{
		Name:          "Vanilla + Vite",
		ShortName:     "Vanilla",
		Description:   "Vanilla + Vite development server",
		Assets:        vanilla,
		FilesToDelete: []string{"frontend/.gitignore", "frontend/index.html", "frontend/favicon.svg", "frontend/main.js", "frontend/style.css"},
	},
	{
		Name:          "Vanilla + Vite (Typescript)",
		ShortName:     "Vanilla-TS",
		Description:   "Vanilla + Vite development server",
		Assets:        vanillats,
		FilesToDelete: []string{"frontend/.gitignore", "frontend/index.html", "frontend/favicon.svg", "frontend/src/main.ts", "frontend/src/style.css"},
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
	commonAssets, err := debme.FS(common, "assets/common")
	checkError(err)

	s.CD("..")
	s.ENDIR("templates")
	s.CD("templates")
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

	// Add common files
	g := gosod.New(commonAssets)
	g.SetTemplateFilters([]string{})
	err = g.Extract(".", nil)
	checkError(err)

	// Add custom files
	g = gosod.New(assets)
	g.SetTemplateFilters([]string{})
	err = g.Extract(".", nil)
	checkError(err)

	//s.ECHO("HERE")
	//s.EXEC("wails init -n " + shortName + "test -t ./" + shortName)
	//s.ECHO("HERE")
	//s.CD(shortName + "test")
	//s.ECHO("HERE")
	//s.REPLACEALL("go.mod", s.Sub{"// replace": "replace"})
	//s.ECHO("HERE")
	//s.EXEC("wails build -debug")
	//s.ECHO("HERE")
	s.CD(cwd)
}
