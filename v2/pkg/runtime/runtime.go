package runtime

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"log"
	goruntime "runtime"
)

const contextError = `An invalid context was passed. This method requires the specific context given in the lifecycle hooks:
https://wails.io/docs/reference/runtime/intro`

func getFrontend(ctx context.Context) frontend.Frontend {
	if ctx == nil {
		pc, _, _, _ := goruntime.Caller(1)
		funcName := goruntime.FuncForPC(pc).Name()
		log.Fatalf("cannot call '%s': context is nil", funcName)
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
		log.Fatalf("cannot call '%s': context is nil", funcName)
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
		log.Fatalf("cannot call '%s': context is nil", funcName)
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
		log.Fatalf("cannot call Quit: context is nil")
	}
	appFrontend := getFrontend(ctx)
	appFrontend.Quit()
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
