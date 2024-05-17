package services

import "github.com/wailsapp/wails/v3/pkg/application"

// OtherService is a struct
// that does things
type OtherService struct {
	t int
}

type Address struct {
	Street  string
	State   string
	Country string
}

// Yay does this and that
func (o *OtherService) Yay() *Address {
	return &Address{
		Street:  "123 Pitt Street",
		State:   "New South Wales",
		Country: "Australia",
	}
}

func NewOtherService() application.Service {
	return application.NewService(&OtherService{})
}
