// Package mypackage does all the things a mypackage can do
package mypackage

type Address struct {
	Number   int
	Street   string
	Town     string
	Postcode string
}

// Person defines a Person in the application
type Person struct {
	// Name is a name
	Name    string
	Age     int
	Address *Address
}

// Manager is the Mr Manager
type Manager struct {
	Name  string
	TwoIC *Person
}

// Hire me some peoples!
func (m *Manager) Hire(name, test string, bob int) *Person {
	return &Person{Name: name}
}

// func NewManagerPointer() *Manager {
// 	return &Manager{}
// }

// func NewManager() Manager {
// 	return Manager{}
// }
