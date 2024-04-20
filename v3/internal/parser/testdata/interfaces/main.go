// This test produces different results depending on whether GODEBUG=gotypesalias=1 is set
// https://pkg.go.dev/go/types#Alias
package main

import (
	"io"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name string `json:"name"`

	// Writer is not visible, because it is an interface
	Writer io.Writer
}

type GreetService struct{}

func (*GreetService) Greet(person Person) string {
	return "Hello " + person.Name
}

// Write will not be bound, because writer is an interface
func (*GreetService) Write(writer io.Writer) {
	writer.Write([]byte("Hello"))
}

// GetWriter will not be bound, because the return value is an interface
func (*GreetService) GetWriter() io.Writer {
	return os.Stdout
}

func main() {
	app := application.New(application.Options{
		Bind: []interface{}{
			&GreetService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
