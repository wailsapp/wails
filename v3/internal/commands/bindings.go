package commands

import (
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser"
)

func GenerateBindings(options *flags.GenerateBindingsOptions) error {
	return parser.GenerateBindingsAndModels(options)
}
