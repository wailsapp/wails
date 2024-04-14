package other

type Person struct {
	Name string
}

type OtherService struct{}

func (o *OtherService) Greet(person Person) string {
	return "Hello " + person.Name + ", my name is other.OtherService"
}
