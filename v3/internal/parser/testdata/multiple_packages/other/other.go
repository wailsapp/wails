package other

import "github.com/google/uuid"

type Person struct {
	UUID uuid.UUID
	Name string
}

type OtherService struct{}

func (o *OtherService) Greet(person Person) string {
	return "Hello " + person.Name + " (" + person.UUID.String() + "), my name is other.OtherService"
}
