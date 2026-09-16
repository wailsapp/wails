//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include <stdlib.h>
#include "services_provider_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

type macosServicesProvider struct {
	app *App
}

func newServicesProviderImpl(app *App) servicesProviderImpl {
	return &macosServicesProvider{app: app}
}

// servicesRegisterTask defers the native registration until App.Run so
// Register may be called from main before the application starts.
type servicesRegisterTask struct {
	provider *macosServicesProvider
	name     string
}

func (t servicesRegisterTask) Run() {
	t.provider.registerNow(t.name)
}

func (p *macosServicesProvider) register(name string) error {
	p.app.runOrDeferToAppRun(servicesRegisterTask{provider: p, name: name})
	return nil
}

func (p *macosServicesProvider) registerNow(name string) {
	var sendTypes []string
	p.app.ServicesProvider.lock.RLock()
	if def, ok := p.app.ServicesProvider.services[name]; ok {
		sendTypes = def.SendTypes
	}
	p.app.ServicesProvider.lock.RUnlock()
	if sendTypes == nil {
		sendTypes = []string{}
	}
	encoded, _ := json.Marshal(sendTypes)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cTypes := C.CString(string(encoded))
	defer C.free(unsafe.Pointer(cTypes))
	InvokeSync(func() {
		C.wailsServicesProviderRegister(cName, cTypes)
	})
}

// servicesProviderResponse is the JSON handed back to the Objective-C IMP.
type servicesProviderResponse struct {
	ServiceResponse
	Error string `json:"error,omitempty"`
}

//export servicesProviderHandle
func servicesProviderHandle(name *C.char, requestJSON *C.char) *C.char {
	var request ServiceRequest
	if err := json.Unmarshal([]byte(C.GoString(requestJSON)), &request); err != nil {
		return C.CString(`{"error":"Wails: cannot decode the service request"}`)
	}
	response := servicesProviderResponse{}
	if globalApplication == nil || globalApplication.ServicesProvider == nil {
		response.Error = "Wails: application is not running"
	} else {
		result, err := globalApplication.ServicesProvider.handle(C.GoString(name), request)
		response.ServiceResponse = result
		if err != nil {
			response.Error = err.Error()
			globalApplication.error("services: %s: %v", C.GoString(name), err)
		}
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		return C.CString(`{"error":"Wails: cannot encode the service response"}`)
	}
	return C.CString(string(encoded))
}
