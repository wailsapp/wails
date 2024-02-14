package main

import (
	"embed"
	_ "embed"
	"log"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/oauth"
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
			Handler: application.AssetFileServerFS(assets),
		},
		Plugins: map[string]application.Plugin{
			"github.com/wailsapp/wails/v3/plugins/oauth": oAuthPlugin,
		},
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

	app.Events.On("github-login", func(e *application.WailsEvent) {
		oAuthPlugin.Github()
	})
	app.Events.On("github-logout", func(e *application.WailsEvent) {
		oAuthPlugin.LogoutGithub()
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
