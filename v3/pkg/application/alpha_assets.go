//go:build !wails_native

package application

import "embed"

//go:embed assets/*
var alphaAssets embed.FS

// AlphaAssets is the default asset set for applications with a web frontend.
var AlphaAssets = AssetOptions{
	Handler: BundledAssetFileServer(alphaAssets),
}
