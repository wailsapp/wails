//+build experimental

package runtime

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"log"
	goruntime "runtime"
)

func getFrontend(ctx context.Context) frontend.Frontend {
	result := ctx.Value("frontend")
	if result != nil {
		return result.(frontend.Frontend)
	}
	pc, _, _, _ := goruntime.Caller(1)
	funcName := goruntime.FuncForPC(pc).Name()
	log.Fatalf("cannot call '%s': Application not initialised", funcName)
	return nil
}

// Quit the application
func Quit(ctx context.Context) {
	frontend := getFrontend(ctx)
	frontend.Quit()
}
