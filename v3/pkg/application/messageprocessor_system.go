package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	SystemIsDarkMode = 0
	Environment      = 1
)

var systemMethodNames = map[int]string{
	SystemIsDarkMode: "IsDarkMode",
	Environment:      "Environment",
}

func (m *MessageProcessor) processSystemMethod(req *RuntimeRequest) (any, error) {
	switch req.Method {
	case SystemIsDarkMode:
		return globalApplication.Env.IsDarkMode(), nil
	case Environment:
		return globalApplication.Env.Info(), nil
	default:
		return nil, errs.NewInvalidSystemCallErrorf("unknown method: %d", req.Method)
	}
}
