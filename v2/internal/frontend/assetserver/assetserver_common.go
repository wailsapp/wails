package assetserver

import (
	iofs "io/fs"
	"path"

	"github.com/wailsapp/wails/v2/internal/fs"
)

func prepareAssetsForServing(assets iofs.FS) (iofs.FS, error) {
	if _, err := assets.Open("."); err != nil {
		return nil, err
	}

	subDir, err := fs.FindPathToFile(assets, "index.html")
	if err != nil {
		return nil, err
	}

	assets, err = iofs.Sub(assets, path.Clean(subDir))
	if err != nil {
		return nil, err
	}
	return assets, nil
}
