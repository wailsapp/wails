package buildmeta

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

// Revision reuses the Git library already linked by Wails, including its
// worktree/common-directory support. No git executable is launched.
func Revision(ctx context.Context, root string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	directory, err := repositoryDirectory(root)
	if err != nil {
		return "", err
	}
	repository, err := git.PlainOpenWithOptions(directory, &git.PlainOpenOptions{EnableDotGitCommonDir: true})
	if err != nil {
		return "", fmt.Errorf("open application Git repository: %w", err)
	}
	config, err := repository.Config()
	if err != nil {
		return "", err
	}
	if format := config.Raw.Section("extensions").Option("objectformat"); format != "" && format != "sha1" {
		return "", fmt.Errorf("application Git object format %q is not supported by revision lookup", format)
	}
	if storage := config.Raw.Section("extensions").Option("refStorage"); storage != "" && storage != "files" {
		return "", fmt.Errorf("application Git reference storage %q is not supported by revision lookup", storage)
	}
	head, err := repository.Head()
	if err != nil {
		return "", fmt.Errorf("resolve application Git HEAD: %w", err)
	}
	if head.Hash().IsZero() {
		return "", fmt.Errorf("application Git HEAD has no valid revision")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

// Find the nearest marker first: an unsupported inner repository must not make
// go-git silently search upward and report an unrelated ancestor's revision.
func repositoryDirectory(root string) (string, error) {
	directory, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Lstat(filepath.Join(directory, ".git")); err == nil {
			return directory, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("no application Git repository above %s", root)
		}
		directory = parent
	}
}
