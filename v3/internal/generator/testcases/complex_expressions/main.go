package main

import (
	_ "embed"
	"log"
	"slices"

	"github.com/wailsapp/wails/v3/internal/generator/testcases/complex_expressions/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service1 struct{}
type Service2 struct{}
type Service3 struct{}
type Service4 struct{}
type Service5 struct{}
type Service6 struct{}

var GlobalServices []application.Service

func main() {
	services := []application.Service{
		application.NewService(&Service1{}),
	}

	services = append(services, application.NewService(&Service2{}), application.Service{}, application.Service{})
	services[2] = application.NewService(&Service3{})

	var options = application.Options{
		Services: config.Services,
	}

	provider := config.MoreServices()
	provider.Init()

	pinit := config.NewProviderInitialiser()
	pinit.InitProvider(&provider)
	// Method resolution should work here just like above.
	config.NewProviderInitialiser().InitProvider(&services[3])

	copy(options.Services, []application.Service{application.NewService(&Service4{}), provider.HeresAnotherOne, provider.OtherService.(application.Service)})
	(options.Services) = append(options.Services, slices.Insert(services, 1, GlobalServices[2:]...)...)

	app := application.New(options)

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	var global = make([]application.Service, 4)

	global[0] = application.NewService(&Service5{})
	global = slices.Replace(global, 1, 4,
		application.NewService(&Service6{}),
		config.MoreServices().AService,
	)

	GlobalServices = slices.Clip(global)
}
