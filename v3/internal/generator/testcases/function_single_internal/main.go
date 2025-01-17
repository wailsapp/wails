package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
}

// Greet someone
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

// Debugging name
func (*GreetService) ServiceName() string {
	return "GreetService"
}

// Lifecycle
func (*GreetService) ServiceStartup(context.Context, application.ServiceOptions) error {
	return nil
}

// Lifecycle
func (*GreetService) ServiceShutdown() error {
	return nil
}

// Serve some routes
func (*GreetService) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func NewGreetService() application.Service {
	return application.NewService(&GreetService{})
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			NewGreetService(),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
