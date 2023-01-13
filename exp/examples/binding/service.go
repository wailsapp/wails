package main

import (
	"github.com/wailsapp/wails/exp/pkg/application"
	"github.com/wailsapp/wails/exp/pkg/options"
)

//func appOptions() options.Application {
//	return options.Application{
//		Bind: []interface{}{
//			&GreetService{},
//		},
//	}
//}

//func appOptions() options.Application {
//	result := options.Application{
//		Bind: []interface{}{
//			&GreetService{},
//		},
//	}
//	return result
//}

func appOptions() options.Application {
	result := options.Application{
		Bind: []interface{}{
			&GreetService{},
		},
	}
	return result
}

func run() error {
	//o := options.Application{
	//	Bind: []interface{}{
	//		&GreetService{},
	//	},
	//}
	//app := application.New(o)

	//app := application.New(options.Application{
	//	Bind: []interface{}{
	//		&GreetService{},
	//	},
	//})

	app := application.New(appOptions())

	return app.Run()
}
