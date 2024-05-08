package config

import "github.com/wailsapp/wails/v3/pkg/application"

type Service7 struct{}
type Service8 struct{}
type Service9 struct{}
type Service10 struct{}
type Service11 struct{}
type Service12 struct{}

func (*Service7) TestMethod() {}

func (*Service9) TestMethod2() {}

func NewService7() application.Service {
	return application.NewService(new(Service7))
}

func NewService8() (result application.Service) {
	result = application.NewService(&Service8{})
	return
}

type ServiceProvider struct {
	AService application.Service
	*ProviderWithMethod
	HeresAnotherOne application.Service
}

type ProviderWithMethod struct {
	OtherService any
}

func (pm *ProviderWithMethod) Init() {
	pm.OtherService = application.NewService(&Service10{})
}

var Services []application.Service

func init() {
	var ourServices = []application.Service{
		NewService7(),
		NewService8(),
	}

	Services = make([]application.Service, len(ourServices))

	for i, el := range ourServices {
		Services[len(ourServices)-i] = el
	}
}

func MoreServices() ServiceProvider {
	var provider ServiceProvider

	provider.AService = application.NewService(&Service9{})
	provider.ProviderWithMethod = new(ProviderWithMethod)

	return provider
}

type ProviderInitialiser interface {
	InitProvider(provider any)
}

type internalProviderInitialiser struct{}

func NewProviderInitialiser() ProviderInitialiser {
	return internalProviderInitialiser{}
}

func (internalProviderInitialiser) InitProvider(provider any) {
	switch p := provider.(type) {
	case *ServiceProvider:
		p.HeresAnotherOne = application.NewService(&Service11{})
	default:
		if anyp, ok := p.(*any); ok {
			*anyp = application.NewService(&Service12{})
		}
	}
}
