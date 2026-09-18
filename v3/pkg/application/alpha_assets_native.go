//go:build wails_native

package application

// AlphaAssets is retained for source compatibility in native builds. Its
// handler is intentionally nil because wails_native does not include an asset
// server or embedded frontend runtime.
var AlphaAssets = AssetOptions{}
