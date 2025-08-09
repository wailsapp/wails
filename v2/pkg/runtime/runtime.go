package runtime

import (
	"context"
	"log"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
)

const contextError = `An invalid context was passed. This method requires the specific context given in the lifecycle hooks:
https://wails.io/docs/reference/runtime/intro`

func getFrontend(ctx context.Context) frontend.Frontend {
	if ctx == nil {
		pc, _, _, _ := goruntime.Caller(1)
		funcName := goruntime.FuncForPC(pc).Name()
		log.Fatalf("cannot call '%s': %s", funcName, contextError)
	}
	result := ctx.Value("frontend")
	if result != nil {
		return result.(frontend.Frontend)
	}
	pc, _, _, _ := goruntime.Caller(1)
	funcName := goruntime.FuncForPC(pc).Name()
	log.Fatalf("cannot call '%s': %s", funcName, contextError)
	return nil
}

func getLogger(ctx context.Context) *logger.Logger {
	if ctx == nil {
		pc, _, _, _ := goruntime.Caller(1)
		funcName := goruntime.FuncForPC(pc).Name()
		log.Fatalf("cannot call '%s': %s", funcName, contextError)
	}
	result := ctx.Value("logger")
	if result != nil {
		return result.(*logger.Logger)
	}
	pc, _, _, _ := goruntime.Caller(1)
	funcName := goruntime.FuncForPC(pc).Name()
	log.Fatalf("cannot call '%s': %s", funcName, contextError)
	return nil
}

func getEvents(ctx context.Context) frontend.Events {
	if ctx == nil {
		pc, _, _, _ := goruntime.Caller(1)
		funcName := goruntime.FuncForPC(pc).Name()
		log.Fatalf("cannot call '%s': %s", funcName, contextError)
	}
	result := ctx.Value("events")
	if result != nil {
		return result.(frontend.Events)
	}
	pc, _, _, _ := goruntime.Caller(1)
	funcName := goruntime.FuncForPC(pc).Name()
	log.Fatalf("cannot call '%s': %s", funcName, contextError)
	return nil
}

// Quit the application
func Quit(ctx context.Context) {
	if ctx == nil {
		log.Fatalf("Error calling 'runtime.Quit': %s", contextError)
	}
	appFrontend := getFrontend(ctx)
	appFrontend.Quit()
}

// Hide the application
func Hide(ctx context.Context) {
	if ctx == nil {
		log.Fatalf("Error calling 'runtime.Hide': %s", contextError)
	}
	appFrontend := getFrontend(ctx)
	appFrontend.Hide()
}

// Show the application if it is hidden
func Show(ctx context.Context) {
	if ctx == nil {
		log.Fatalf("Error calling 'runtime.Show': %s", contextError)
	}
	appFrontend := getFrontend(ctx)
	appFrontend.Show()
}

// EnvironmentInfo contains information about the environment
type EnvironmentInfo struct {
	BuildType string `json:"buildType"`
	Platform  string `json:"platform"`
	Arch      string `json:"arch"`
}

// Environment returns information about the environment
func Environment(ctx context.Context) EnvironmentInfo {
	var result EnvironmentInfo
	buildType := ctx.Value("buildtype")
	if buildType != nil {
		result.BuildType = buildType.(string)
	}
	result.Platform = goruntime.GOOS
	result.Arch = goruntime.GOARCH
	return result
}
