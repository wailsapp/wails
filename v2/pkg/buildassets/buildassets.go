package buildassets

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/project"
)

//go:embed build
var assets embed.FS

const (
	rootFolder = "build"
)

// Install will install all default project assets
func Install(targetDir string) error {
	templateDir := gosod.New(assets)
	err := templateDir.Extract(targetDir, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetLocalPath returns the local path of the requested build asset file
func GetLocalPath(projectData *project.Project, file string) string {
	return filepath.Clean(filepath.Join(projectData.Path, rootFolder, filepath.FromSlash(file)))
}

// ReadFile reads the file from the project build folder.
// If the file does not exist it falls back to the embedded file and the file will be written
// to the disk for customisation.
func ReadFile(projectData *project.Project, file string) ([]byte, error) {
	fs := os.DirFS(filepath.ToSlash(projectData.Path)) // os.DirFs always operates on "/" as separatator
	file = path.Join(rootFolder, file)

	content, err := iofs.ReadFile(fs, file)
	if errors.Is(err, iofs.ErrNotExist) {
		// The file does not exist, let's read it from the assets FS and write it to disk
		content, err := iofs.ReadFile(assets, file)
		if err != nil {
			return nil, err
		}

		if err := writeFileSystemFile(projectData, file, content); err != nil {
			return nil, fmt.Errorf("Unable to create file in build folder: %s", err)
		}
		return content, nil
	}

	return content, err
}

// ReadFileWithProjectData reads the file from the project build folder and replaces ProjectInfo if necessary.
// If the file does not exist it falls back to the embedded file and the file will be written
// to the disk for customisation. The file written is the original unresolved one.
func ReadFileWithProjectData(projectData *project.Project, file string) ([]byte, error) {
	content, err := ReadFile(projectData, file)
	if err != nil {
		return nil, err
	}

	content, err = resolveProjectData(content, projectData)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve data in %s: %w", file, err)
	}
	return content, nil
}

// ReadOriginalFileWithProjectDataAndSave reads the file from the embedded assets and replaces
// ProjectInfo if necessary.
// It will also write the resolved final file back to the project build folder.
func ReadOriginalFileWithProjectDataAndSave(projectData *project.Project, file string) ([]byte, error) {
	file = path.Join(rootFolder, file)
	content, err := iofs.ReadFile(assets, file)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file %s: %w", file, err)
	}

	content, err = resolveProjectData(content, projectData)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve data in %s: %w", file, err)
	}

	if err := writeFileSystemFile(projectData, file, content); err != nil {
		return nil, fmt.Errorf("Unable to create file in build folder: %w", err)
	}
	return content, nil
}

type assetData struct {
	Name string
	Info project.Info
}

func resolveProjectData(content []byte, projectData *project.Project) ([]byte, error) {
	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return nil, err
	}

	data := &assetData{
		Name: projectData.Name,
		Info: projectData.Info,
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func writeFileSystemFile(projectData *project.Project, file string, content []byte) error {
	path := filepath.Clean(filepath.Join(projectData.Path, filepath.FromSlash(file)))

	if dir := filepath.Dir(path); !fs.DirExists(dir) {
		if err := fs.MkDirs(dir, 0755); err != nil {
			return fmt.Errorf("Unable to create directory: %w", err)
		}
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return err
	}
	return nil
}
