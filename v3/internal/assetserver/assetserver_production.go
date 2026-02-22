//go:build production

package assetserver

func defaultIndexHTML(_ string) []byte {
	return []byte("index.html not found")
}

func (a *AssetServer) LogDetails() {}
