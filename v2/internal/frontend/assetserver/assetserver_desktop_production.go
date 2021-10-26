//go:build production

package assetserver

func (a *DesktopAssetServer) ReadFile(filename string) ([]byte, error) {
	return a.assets.ReadFile(filename)
}
