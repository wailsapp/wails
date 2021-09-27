package servicebus

import (
	"context"
	"log"
	"runtime"
)

func ExtractBus(ctx context.Context) *ServiceBus {
	bus := ctx.Value("bus")
	if bus == nil {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		log.Fatalf("cannot call '%s': Application not initialised", funcName)
	}
	return bus.(*ServiceBus)
}
