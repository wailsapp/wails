package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

type TimeAlias = time.Time

type TimeStruct struct {
	time.Time
}

type TimeAliasStruct struct {
	TimeAlias
}

type TimeFieldStruct struct {
	T1 time.Time
	T2 TimeAlias
	T3 TimeStruct
	T4 TimeAliasStruct

	Q time.Time `json:",string"`
	O time.Time `json:",omitempty"`
	P *time.Time
	A [3]time.Time
	S []time.Time
	M map[time.Time]time.Time
	I struct{ T time.Time }
}

func (*Service) GetTime() (_ time.Time) {
	return
}

func (*Service) SetTime(t time.Time) (_ error) {
	return
}

func (*Service) DoVariousTimeThings() (_ TimeFieldStruct) {
	return
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&Service{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
