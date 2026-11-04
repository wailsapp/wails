package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type stagedPublication struct {
	source, destination, staged, backup string
	previous, installed                 bool
}

func associatedPlanArtifacts(node pipeline.Node) []manifestPlanArtifact {
	var result []manifestPlanArtifact
	if spec, ok := node.Spec.(pipeline.PublishSpec); ok {
		for _, file := range spec.Associated {
			result = append(result, manifestPlanArtifact{Target: file.Identity.Target.OS + "/" + file.Identity.Target.Arch, Format: file.Identity.Format, Kind: file.Identity.DisplayKind(), Path: file.Destination, Signed: file.Identity.Signed})
		}
	}
	return result
}

// Associated files share one rollback transaction. Filesystems do not provide
// atomic multi-file renames, so readers should use the final artifact receipt.
func (h *manifestHandler) publishAssociated(spec pipeline.PublishSpec, rename func(string, string) error) (pipeline.RunResult, error) {
	files := append([]pipeline.PublishedFile{{Source: spec.Source, Destination: spec.Destination}}, spec.Associated...)
	paths := make([]stagedPublication, len(files))
	destinations := make(map[string]bool, len(files))
	for i, file := range files {
		source := filepath.ToSlash(filepath.Clean(file.Source))
		destination := filepath.ToSlash(filepath.Clean(file.Destination))
		if !strings.HasPrefix(source, ".wails/") || destination == ".wails" || strings.HasPrefix(destination, ".wails/") {
			return pipeline.RunResult{}, fmt.Errorf("invalid associated publication paths")
		}
		var err error
		paths[i].source, err = existingArtifactPathInsideProject(h.root, file.Source)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		paths[i].destination, err = pathInsideProject(h.root, file.Destination)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		if destinations[paths[i].destination] || i != 0 && filepath.Dir(paths[i].destination) != filepath.Dir(paths[0].destination) {
			return pipeline.RunResult{}, fmt.Errorf("associated publications need distinct sibling destinations")
		}
		destinations[paths[i].destination] = true
	}
	parent := filepath.Dir(paths[0].destination)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return pipeline.RunResult{}, err
	}
	transaction, err := os.MkdirTemp(parent, ".wails-publish-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(transaction)
		}
	}()
	for i := range paths {
		paths[i].staged = filepath.Join(transaction, fmt.Sprintf("file-%d", i))
		paths[i].backup = filepath.Join(transaction, fmt.Sprintf("previous-%d", i))
		if err := copyManifestPath(paths[i].source, paths[i].staged); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	if err := installPublications(paths, rename); err != nil {
		for _, file := range paths {
			if file.previous {
				cleanup = false
				return pipeline.RunResult{}, errors.Join(err, fmt.Errorf("publication backups retained at %s", transaction))
			}
		}
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Output: spec.Destination}, nil
}

func installPublications(files []stagedPublication, rename func(string, string) error) (err error) {
	defer func() {
		if err == nil {
			return
		}
		for i := len(files) - 1; i >= 0; i-- {
			file := files[i]
			if file.installed {
				err = errors.Join(err, os.Remove(file.destination))
			}
			if file.previous {
				restoreErr := rename(file.backup, file.destination)
				err = errors.Join(err, restoreErr)
				if restoreErr == nil {
					files[i].previous = false
				}
			}
		}
	}()
	for i := range files {
		if info, statErr := os.Lstat(files[i].destination); statErr == nil {
			if !info.Mode().IsRegular() {
				return fmt.Errorf("publication destination must be a regular file")
			}
			if err := rename(files[i].destination, files[i].backup); err != nil {
				return err
			}
			files[i].previous = true
		} else if !os.IsNotExist(statErr) {
			return statErr
		}
		if err := rename(files[i].staged, files[i].destination); err != nil {
			return err
		}
		files[i].installed = true
	}
	return nil
}
