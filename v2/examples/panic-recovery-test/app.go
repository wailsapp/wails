package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("------------------------------%#v\n", err)
			}
		}()
		time.Sleep(5 * time.Second)
		// Fix signal handlers right before potential panic using the Wails runtime
		runtime.ResetSignalHandlers()
		// Nil pointer dereference - causes SIGSEGV
		var t *time.Time
		fmt.Println(t.Unix())
	}()

	return fmt.Sprintf("Hello %s, It's show time!", name)
}
