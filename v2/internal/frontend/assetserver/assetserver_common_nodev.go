//go:build !dev

package assetserver

import (
	"context"
	"io/fs"
)

func isAssetDirOverride(ctx context.Context, serverType string, assets fs.FS) (fs.FS, error) {
	return nil, nil
}
