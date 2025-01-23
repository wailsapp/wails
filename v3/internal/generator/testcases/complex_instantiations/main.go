package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/internal/generator/testcases/complex_instantiations/other"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service1 struct{}
type Service2 struct{}
type Service3 struct{}
type Service4 struct{}
type Service5 struct{}
type Service6 struct{}
type Service7 struct{}
type Service8 struct{}
type Service9 struct{}
type Service10 struct{}
type Service11 struct{}
type Service12 struct{}
type Service13 struct{}

func main() {
	factory := NewFactory[Service1, Service2]()
	otherFactory := other.NewFactory[Service3, Service4]()

	app := application.New(application.Options{
		Services: append(append(
			[]application.Service{
				factory.Get(),
				factory.GetU(),
				otherFactory.Get(),
				otherFactory.GetU(),
				application.NewService(&Service5{}),
				ServiceInitialiser[Service6]()(&Service6{}),
				other.CustomNewService(Service7{}),
				other.ServiceInitialiser[Service8]()(&Service8{}),
				application.NewServiceWithOptions(&Service13{}, application.ServiceOptions{Name: "custom name"}),
				other.LocalService,
			},
			CustomNewServices[Service9, Service10]()...),
			other.CustomNewServices[Service11, Service12]()...),
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
