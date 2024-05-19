package services

import (
	"embed"
	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v3/internal/debug"
	"github.com/wailsapp/wails/v3/internal/flags"
	"io/fs"
)

//go:embed create/*
var CreateTemplate embed.FS

func Install(options *flags.ServiceCreateOptions) error {
	type Data struct {
		*flags.ServiceCreateOptions
		LocalModulePath string
	}
	var data = &Data{
		ServiceCreateOptions: options,
		LocalModulePath:      debug.LocalModulePath,
	}
	tmpl, err := fs.Sub(CreateTemplate, "create")
	if err != nil {
		return err
	}
	return gosod.New(tmpl).Extract(options.OutputDir, data)
}
