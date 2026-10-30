// Package releasenotes formats release Markdown for terminal output.
package releasenotes

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

// Render preserves the Markdown if the terminal renderer cannot be created.
func Render(result string, noColour bool) string {
	var renderer *glamour.TermRenderer
	var err error

	if noColour {
		renderer, err = glamour.NewTermRenderer(glamour.WithStyles(styles.NoTTYStyleConfig))
	} else {
		renderer, err = glamour.NewTermRenderer(glamour.WithAutoStyle())
	}
	if err != nil {
		return result
	}
	result, err = renderer.Render(result)
	if err != nil {
		return err.Error()
	}
	return result
}
