package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "macOS Materials Demo",
		Description: "A demo of macOS materials support",
	})

	// Create a window with sidebar material background
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Sidebar Material Demo",
		Width:  400,
		Height: 300,
		X:      100,
		Y:      100,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropSidebar,
			MaterialOptions: application.MacMaterialOptions{
				BlendingMode:         application.MacMaterialBlendingModeWithinWindow,
				State:               application.MacMaterialStateActive,
				EmphasizedAppearance: false,
			},
			TitleBar: application.MacTitleBar{
				AppearsTransparent: true,
				FullSizeContent:    true,
			},
		},
		HTML: "<html><head><style>body{background:rgba(255,255,255,0.1);font-family:-apple-system;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;color:#333;}.content{background:rgba(255,255,255,0.7);padding:20px;border-radius:10px;backdrop-filter:blur(10px);text-align:center;}h1{margin-top:0;color:#333;}p{color:#666;}</style></head><body><div class='content'><h1>Sidebar Material</h1><p>This window demonstrates the Sidebar material backdrop effect.</p><p>Move this window around to see the material effect in action!</p></div></body></html>",
	})
	
	// Mac options are now properly set during window creation
	
	window.Show()

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
} 