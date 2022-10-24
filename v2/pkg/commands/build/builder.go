package build

import (
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// Builder defines a builder that can build Wails applications
type Builder interface {
	SetProjectData(projectData *project.Project)
	BuildFrontend(*clilogger.CLILogger) error
	CompileProject(*Options) error
	OutputFilename(*Options) string
	CleanUp()
}
