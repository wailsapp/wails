package main

import (
	_ "embed"
	"log"
	"slices"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service1 struct{}
type Service2 struct{}
type Service3 struct{}
type Service4 struct{}
type Service5 struct{}
type Service6 struct{}

var GlobalServices []any

func main() {
	services := []any{
		&Service1{},
	}

	services = append(services, &Service2{}, nil)
	services[2] = &Service3{}

	var options application.Options

	copy(options.Bind, []any{&Service4{}})
	(options.Bind) = append(options.Bind, slices.Insert(services, 1, GlobalServices[2:]...)...)

	app := application.New(options)

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	var global = make([]any, 3)

	global[0] = &Service5{}
	global = slices.Replace(global, 1, 3, any(&Service6{}))

	GlobalServices = slices.Clip(global)
}
