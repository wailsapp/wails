package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/events"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Video Demo",
		Description: "A demo of HTML5 Video API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			WndProcInterceptor:            nil,
			DisableQuitOnLastWindowClosed: false,
			WebviewUserDataPath:           "",
			WebviewBrowserPath:            "",
		},
	})
	app.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
		log.Println("ApplicationDidFinishLaunching")
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		BackgroundColour: application.NewRGB(33, 37, 41),
		Mac: application.MacWindow{
			DisableShadow: true,
			WebviewPreferences: application.MacWebviewPreferences{
				FullscreenEnabled: application.Enabled,
			},
		},
		HTML: "<video controls width=\"500\" >\n        <source\n          src=\"https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.webm\"\n          type=\"video/webm\"\n        />\n        <source\n          src=\"https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4\"\n          type=\"video/mp4\"\n        />\n      </video>",
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
