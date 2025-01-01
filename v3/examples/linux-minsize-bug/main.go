package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

    app := application.New(application.Options{
	Name: "v3-alpha8.3-minsize-linux-bug",
    })

    app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
	Title:     "v3-alpha8.3-minsize-linux-bug",
	URL:       "/",
	Width:     1024,
	Height:    768,
	MinWidth:  1024,
	MinHeight: 768,
    })

    err := app.Run()
    if err != nil {
	log.Fatal(err)
    }
}
