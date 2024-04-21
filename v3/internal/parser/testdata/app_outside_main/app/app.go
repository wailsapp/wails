package app

import (
	"github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/app/services"
	"github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/other"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name string `json:"name"`
}

type GreetService struct{}

func (*GreetService) Greet(person Person) string {
	return "Hello " + person.Name
}

func NewApp() *application.App {
	return application.New(application.Options{
		Bind: []interface{}{
			&GreetService{},
			&other.OtherService{},
			&services.OtherService{},
		},
	})
}
