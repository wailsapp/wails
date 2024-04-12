package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Field struct {
	Letter rune
	Number uint
}

type App struct {
	Locale map[string]map[string]string
	Board  [][]Field
}

type OtherService struct {
	SomeVariable int
}

func (o *OtherService) NewApp(locale map[string]map[string]string, board [][]Field) App {
	return App{locale, board}
}

func (o *OtherService) GetLocale(app App) map[string]map[string]string {
	return app.Locale
}

func (o *OtherService) GetBoard(app App) [][]Field {
	return app.Board
}

func main() {
	app := application.New(application.Options{
		Bind: []interface{}{
			&OtherService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
