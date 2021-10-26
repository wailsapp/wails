package webserver

import "github.com/wailsapp/wails/v2/internal/assetdb"

var (
	// WebAssets is our single asset db instance.
	// It will be constructed by a dynamically generated method in this directory.
	WebAssets *assetdb.AssetDB = assetdb.NewAssetDB()
)
