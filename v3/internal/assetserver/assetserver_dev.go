//go:build !production

package assetserver

func (a *AssetServer) LogDetails() {
	var info = []any{
		"middleware", a.options.Middleware != nil,
		"handler", a.options.Handler != nil,
	}
	if devServerURL := GetDevServerURL(); devServerURL != "" {
		info = append(info, "devServerURL", devServerURL)
	}
	a.options.Logger.Info("AssetServer Info:", info...)
}
