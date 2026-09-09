package main

import (
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

type AcceptanceService struct{}

func (*AcceptanceService) Ping(value string) string { return "go-response:" + value }
func (*AcceptanceService) Confirm(value string) {
	log.Print("HCL_STDERR_LOG_PASS")
	fmt.Printf("HCL_BINDING_ROUNDTRIP %s platform=%s/%s\n", value, runtime.GOOS, runtime.GOARCH)
}
func main() {
	app := application.New(application.Options{Name: "Wails Dev Acceptance", Services: []application.Service{application.NewService(&AcceptanceService{})}, Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)}})
	app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Wails Dev Acceptance", URL: "/"})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
