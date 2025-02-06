//wails:inject console.log("Hello everywhere!");
//wails:inject **:console.log("Hello everywhere again!");
//wails:inject *c:console.log("Hello Classes!");
//wails:inject *i:console.log("Hello Interfaces!");
//wails:inject j*:console.log("Hello JS!");
//wails:inject jc:console.log("Hello JS Classes!");
//wails:inject ji:console.log("Hello JS Interfaces!");
//wails:inject t*:console.log("Hello TS!");
//wails:inject tc:console.log("Hello TS Classes!");
//wails:inject ti:console.log("Hello TS Interfaces!");
package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/internal/generator/testcases/directives/otherpackage"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type IgnoredType struct {
	Field int
}

//wails:inject j*:/**
//wails:inject j*: * @param {string} arg
//wails:inject j*: * @returns {Promise<void>}
//wails:inject j*: */
//wails:inject j*:export async function CustomMethod(arg) {
//wails:inject t*:export async function CustomMethod(arg: string): Promise<void> {
//wails:inject     await InternalMethod("Hello " + arg + "!");
//wails:inject }
type Service struct{}

func (*Service) VisibleMethod(otherpackage.Dummy) {}

//wails:ignore
func (*Service) IgnoredMethod(IgnoredType) {}

//wails:internal
func (*Service) InternalMethod(string) {}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&Service{}),
			application.NewService(&unexportedService{}),
			application.NewService(&InternalService{}),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
