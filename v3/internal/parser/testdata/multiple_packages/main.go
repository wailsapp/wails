package main

import (
	"log"
	"runtime/debug"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/parser/testdata/multiple_packages/other"
	otherother "github.com/wailsapp/wails/v3/internal/parser/testdata/multiple_packages/other/other"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
	ServiceName string
}

type Greeter GreetService

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + ", my name is " + g.ServiceName
}

func (g *GreetService) BuildInfo() (*debug.BuildInfo, bool) {
	return debug.ReadBuildInfo()
}

func main() {
	tupel := lo.T2(0, 1)

	app := application.New(application.Options{
		Bind: []interface{}{
			log.Default(),
			&GreetService{
				ServiceName: "GreetService",
			},
			&Greeter{
				ServiceName: "GreetService",
			},
			&other.OtherService{},
			&otherother.OtherService{},
			&tupel,
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
