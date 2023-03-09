package templates

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/debug"
	"io/fs"
	"os"

	"github.com/wailsapp/wails/v3/internal/flags"

	"github.com/leaanthony/gosod"

	"github.com/samber/lo"
)

//go:embed lit
var lit embed.FS

//go:embed lit-ts
var litTS embed.FS

//go:embed vue
var vue embed.FS

//go:embed vue-ts
var vueTS embed.FS

//go:embed react
var react embed.FS

//go:embed react-ts
var reactTS embed.FS

//go:embed react-swc
var reactSWC embed.FS

//go:embed react-swc-ts
var reactSWCTS embed.FS

//go:embed svelte
var svelte embed.FS

//go:embed svelte-ts
var svelteTS embed.FS

//go:embed preact
var preact embed.FS

//go:embed preact-ts
var preactTS embed.FS

//go:embed vanilla
var vanilla embed.FS

//go:embed vanilla-ts
var vanillaTS embed.FS

type TemplateData struct {
	Name        string
	Description string
	FS          embed.FS
}

var defaultTemplates = []TemplateData{
	{
		Name:        "lit",
		Description: "Template using Lit Web Components: https://lit.dev",
		FS:          lit,
	},
	{
		Name:        "lit-ts",
		Description: "Template using Lit Web Components (TypeScript) : https://lit.dev",
		FS:          litTS,
	},
	{
		Name:        "vue",
		Description: "Template using Vue: https://vuejs.org",
		FS:          vue,
	},
	{
		Name:        "vue-ts",
		Description: "Template using Vue (TypeScript): https://vuejs.org",
		FS:          vueTS,
	},
	{
		Name:        "react",
		Description: "Template using React: https://reactjs.org",
		FS:          react,
	},
	{
		Name:        "react-ts",
		Description: "Template using React (TypeScript): https://reactjs.org",
		FS:          reactTS,
	},
	{
		Name:        "react-swc",
		Description: "Template using React with SWC: https://reactjs.org & https://swc.rs",
		FS:          reactSWC,
	},
	{
		Name:        "react-swc-ts",
		Description: "Template using React with SWC (TypeScript): https://reactjs.org & https://swc.rs",
		FS:          reactSWCTS,
	},
	{
		Name:        "svelte",
		Description: "Template using Svelte: https://svelte.dev",
		FS:          svelte,
	},
	{
		Name:        "svelte-ts",
		Description: "Template using Svelte (TypeScript): https://svelte.dev",
		FS:          svelteTS,
	},
	{
		Name:        "preact",
		Description: "Template using Preact: https://preactjs.com",
		FS:          preact,
	},
	{
		Name:        "preact-ts",
		Description: "Template using Preact (TypeScript): https://preactjs.com",
		FS:          preactTS,
	},
	{
		Name:        "vanilla",
		Description: "Template using Vanilla JS",
		FS:          vanilla,
	},
	{
		Name:        "vanilla-ts",
		Description: "Template using Vanilla JS (TypeScript)",
		FS:          vanillaTS,
	},
}

func ValidTemplateName(name string) bool {
	return lo.ContainsBy(defaultTemplates, func(template TemplateData) bool {
		return template.Name == name
	})
}

func GetDefaultTemplates() []TemplateData {
	return defaultTemplates
}

type TemplateOptions struct {
	*flags.Init
	LocalModulePath string
}

func Install(options *flags.Init) error {

	templateData := TemplateOptions{
		options,
		debug.LocalModulePath,
	}
	template, found := lo.Find(defaultTemplates, func(template TemplateData) bool {
		return template.Name == options.TemplateName
	})
	if !found {
		return fmt.Errorf("template '%s' not found", options.TemplateName)
	}

	if options.ProjectDir == "." || options.ProjectDir == "" {
		templateData.ProjectDir = lo.Must(os.Getwd())
	}
	templateData.ProjectDir = fmt.Sprintf("%s/%s", options.ProjectDir, options.ProjectName)
	fmt.Printf("Installing template '%s' into '%s'\n", options.TemplateName, options.ProjectDir)
	tfs, err := fs.Sub(template.FS, options.TemplateName)
	if err != nil {
		return err
	}

	return gosod.New(tfs).Extract(options.ProjectDir, templateData)
}
