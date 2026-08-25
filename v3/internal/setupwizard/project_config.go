package setupwizard

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// ProjectMetadataState is the setup wizard's editable project metadata.
type ProjectMetadataState struct {
	ProjectName       string `json:"projectName,omitempty"`
	CompanyName       string `json:"companyName"`
	ProductName       string `json:"productName"`
	ProductIdentifier string `json:"productIdentifier"`
	Description       string `json:"description"`
	Copyright         string `json:"copyright"`
	Comments          string `json:"comments,omitempty"`
	Version           string `json:"version"`
}

// ProjectConfigState is the complete project state edited by the setup wizard.
type ProjectConfigState struct {
	Info ProjectMetadataState `json:"info"`
}

var invalidWizardProjectName = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func (config ProjectConfigState) project(existing manifest.Project) (manifest.Project, error) {
	project := existing
	if config.Info.ProjectName != "" {
		project.Name = invalidWizardProjectName.ReplaceAllString(config.Info.ProjectName, "_")
	} else if project.Name == "" {
		project.Name = invalidWizardProjectName.ReplaceAllString(config.Info.ProductName, "_")
	}
	project.ProductName = config.Info.ProductName
	project.CompanyName = config.Info.CompanyName
	project.Identifier = config.Info.ProductIdentifier
	project.Description = config.Info.Description
	project.Version = config.Info.Version
	project.Copyright = config.Info.Copyright
	project.Comments = config.Info.Comments
	missing := make([]string, 0, 4)
	if project.Name == "" {
		missing = append(missing, "projectName")
	}
	if project.ProductName == "" {
		missing = append(missing, "productName")
	}
	if project.Identifier == "" {
		missing = append(missing, "productIdentifier")
	}
	if project.Version == "" {
		missing = append(missing, "version")
	}
	if len(missing) > 0 {
		return manifest.Project{}, fmt.Errorf("project state requires %s", strings.Join(missing, ", "))
	}
	return project, nil
}

func projectConfigStateFromProject(project manifest.Project) ProjectConfigState {
	return ProjectConfigState{Info: ProjectMetadataState{
		ProjectName:       project.Name,
		CompanyName:       project.CompanyName,
		ProductName:       project.ProductName,
		ProductIdentifier: project.Identifier,
		Description:       project.Description,
		Copyright:         project.Copyright,
		Comments:          project.Comments,
		Version:           project.Version,
	}}
}
