package commands

import (
	"github.com/wailsapp/wails/v3/internal/templates"
)

func GenerateTemplate(options *templates.BaseTemplate) error {
	return templates.GenerateTemplate(options)
}
