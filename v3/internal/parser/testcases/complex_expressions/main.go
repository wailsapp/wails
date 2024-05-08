package main

import (
	_ "embed"
	"log"
	"slices"

	"github.com/wailsapp/wails/v3/internal/parser/testcases/complex_expressions/config"
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

	services = append(services, &Service2{}, nil, nil)
	services[2] = &Service3{}

	var options = application.Options{
		Bind: config.Services,
	}

	provider := config.MoreServices()
	provider.Init()

	pinit := config.NewProviderInitialiser()
	pinit.InitProvider(&provider)
	// Method resolution should work here just like above.
	config.NewProviderInitialiser().InitProvider(&services[3])

	copy(options.Bind, []any{&Service4{}, provider.HeresAnotherOne, provider.OtherService})
	(options.Bind) = append(options.Bind, slices.Insert(services, 1, GlobalServices[2:]...)...)

	app := application.New(options)

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	var global = make([]any, 4)

	global[0] = &Service5{}
	global = slices.Replace(global, 1, 4, any(&Service6{}), any(config.MoreServices().TM2Service))

	GlobalServices = slices.Clip(global)
}
