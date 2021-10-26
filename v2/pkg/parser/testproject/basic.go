package main

import (
	"fmt"

	"testproject/mypackage"

	"github.com/wailsapp/wails/v2"
)

// Basic application struct
type Basic struct {
	runtime *wails.Runtime
}

// // Another application struct
// type Another struct {
// 	runtime *wails.Runtime
// }

// func (a *Another) Doit() {

// }

// // newBasicPointer creates a new Basic application struct
// func newBasicPointer() *Basic {
// 	return &Basic{}
// }

// // newBasic creates a new Basic application struct
// func newBasic() Basic {
// 	return Basic{}
// }

// WailsInit is called at application startup
func (b *Basic) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	b.runtime = runtime
	runtime.Window.SetTitle("jsbundle")
	return nil
}

// WailsShutdown is called at application termination
func (b *Basic) WailsShutdown() {
	// Perform your teardown here
}

// NewPerson creates a new person
func (b *Basic) NewPerson(name string, age int) *mypackage.Person {
	return &mypackage.Person{Name: name, Age: age}
}

// Greet returns a greeting for the given name
func (b *Basic) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

// MultipleGreets returns greetings for the given name
func (b *Basic) MultipleGreets(_ string) []string {
	return []string{"hi", "hello", "croeso!"}
}

// RemovePerson Removes the given person
func (b *Basic) RemovePerson(_ *mypackage.Person) {
	// dummy
}
