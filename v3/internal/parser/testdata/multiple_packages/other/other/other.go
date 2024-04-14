package other

type Person struct {
	Name  string
	Title string
}

type OtherService struct{}

func (o *OtherService) Greet(person Person) string {
	return "Hello " + person.Title + " " + person.Name + ", my name is other.other.OtherService"
}
