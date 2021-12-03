//go:build dev
// +build dev

package assetserver

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/logger"
)

func isAssetDirOverride(ctx context.Context, serverType string, assets fs.FS) (fs.FS, error) {
	_logger, _ := ctx.Value("logger").(*logger.Logger)
	logWarning := func(msg string) {
		if _logger == nil {
			return
		}

		_logger.Warning("[%s] Reloading changed asset files won't work: %s", serverType, msg)
	}

	_, isEmbedFs := assets.(embed.FS)
	if assetdir := ctx.Value("assetdir"); assetdir == nil {
		// We are in dev mode, but have no assetdir, let's check the type of the assets FS.
		if isEmbedFs {
			logWarning("A 'embed.FS' has been provided, but no assetdir has been defined")
		} else {
			// That's fine, the user is using another fs.FS and also has not set an assetdir
		}
	} else {
		// We are in dev mode and an assetdir has been defined
		if isEmbedFs {
			// If the fs.FS is an embed.FS we assume it's serving the files from this directory, therefore replace the assets
			// fs.FS with this directory.
			// Are there any better ways to detect this? If we could somehow determine the compile time source path for the
			// embed.FS this would be great. Currently there seems no way to get to this information.
			absdir, err := filepath.Abs(assetdir.(string))
			if err != nil {
				return nil, err
			}

			if _logger != nil {
				_logger.Debug("[%s] Serving assets from disk: %s", serverType, absdir)
			}

			return os.DirFS(absdir), nil
		} else {
			logWarning("An assetdir has been defined, but no 'embed.FS' has been provided")
		}
	}

	return nil, nil
}
