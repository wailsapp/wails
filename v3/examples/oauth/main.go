package main

import (
	"embed"
	_ "embed"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/oauth"
	"log"
	"os"
)

//go:embed assets
var assets embed.FS

func main() {

	oAuthPlugin := oauth.NewPlugin(oauth.Config{
		Providers: []goth.Provider{
			github.New(
				os.Getenv("clientkey"),
				os.Getenv("secret"),
				"http://localhost:9876/auth/github/callback",
				"email",
				"profile"),
		},
	})

	app := application.New(application.Options{
		Name:        "OAuth Demo",
		Description: "A demo of the oauth Plugin",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
		Plugins: map[string]application.Plugin{
			"github.com/wailsapp/wails/v3/plugins/oauth": oAuthPlugin,
		},
	})

	oAuthWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:         "Login",
		Width:         600,
		Height:        850,
		Hidden:        true,
		DisableResize: true,
		URL:           "http://localhost:9876/auth/github",
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:                  "OAuth Demo",
		DevToolsEnabled:        true,
		OpenInspectorOnStartup: true,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	// Custom event handling
	app.Events.On(oauth.Success, func(e *application.WailsEvent) {
		oAuthWindow.Hide()
	})

	app.Events.On(oauth.Error, func(e *application.WailsEvent) {
		oAuthWindow.Hide()
	})
	app.Events.On("github-login", func(e *application.WailsEvent) {
		oAuthPlugin.Start()
		oAuthWindow.Show()
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
