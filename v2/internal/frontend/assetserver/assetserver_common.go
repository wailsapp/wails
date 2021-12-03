package assetserver

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/leaanthony/slicer"
)

func prepareAssetsForServing(ctx context.Context, serverType string, assets fs.FS) (fs.FS, bool, error) {
	// Let's check if we need an assetdir override
	if assetdir, err := isAssetDirOverride(ctx, serverType, assets); err != nil {
		return nil, false, err
	} else if assetdir != nil {
		return assetdir, true, nil
	}

	// Otherwise let's search for the index.html
	if _, err := assets.Open("."); err != nil {
		return nil, false, err
	}

	subDir, err := pathToIndexHTML(assets)
	if err != nil {
		return nil, false, err
	}

	assets, err = fs.Sub(assets, path.Clean(subDir))
	if err != nil {
		return nil, false, err
	}
	return assets, false, nil
}

func pathToIndexHTML(assets fs.FS) (string, error) {
	stat, _ := fs.Stat(assets, "index.html")
	if stat != nil {
		return ".", nil
	}
	var indexFiles slicer.StringSlicer
	err := fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "index.html") {
			indexFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if indexFiles.Length() > 1 {
		return "", fmt.Errorf("multiple 'index.html' files found in assets")
	}

	path, _ := filepath.Split(indexFiles.AsSlice()[0])
	return path, nil
}
