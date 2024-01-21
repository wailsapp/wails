//go:build !production

package assetserver

func (a *AssetServer) LogDetails() {
	var info = []any{
		"assetsFS", a.options.Assets != nil,
		"middleware", a.options.Middleware != nil,
		"handler", a.options.Handler != nil,
	}
	if a.devServerURL != "" {
		info = append(info, "devServerURL", a.devServerURL)
	}
	a.options.Logger.Info("AssetServer Info:", info...)
}
