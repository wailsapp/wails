package services

import "github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/app/models"

type OtherService struct{}

func (*OtherService) Greet(person models.Person) string {
	return "Hello " + person.Name
}
