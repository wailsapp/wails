package service

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/flags"

	"github.com/leaanthony/gosod"

	"github.com/samber/lo"
)

//go:embed template
var serviceTemplate embed.FS

type TemplateOptions struct {
	*flags.ServiceInit
}

func Install(options *flags.ServiceInit) error {

	if options.OutputDir == "." || options.OutputDir == "" {
		options.OutputDir = filepath.Join(lo.Must(os.Getwd()), options.Name)
	}
	fmt.Printf("Generating service '%s' into '%s'\n", options.Name, options.OutputDir)
	tfs, err := fs.Sub(serviceTemplate, "template")
	if err != nil {
		return err
	}

	return gosod.New(tfs).Extract(options.OutputDir, options)
}
