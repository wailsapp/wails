package config

type Service7 struct{}
type Service8 struct{}
type Service9 struct{}
type Service10 struct{}
type Service11 struct{}
type Service12 struct{}

func (*Service7) TestMethod() {}

func (*Service9) TestMethod2() {}

func NewService7() interface{ TestMethod() } {
	return new(Service7)
}

func NewService8() (result any) {
	result = &Service8{}
	return
}

type TM2 interface {
	TestMethod2()
}

type ServiceProvider struct {
	TM2Service TM2
	*ProviderWithMethod
	HeresAnotherOne any
}

type ProviderWithMethod struct {
	OtherService any
}

func (pm *ProviderWithMethod) Init() {
	pm.OtherService = &Service10{}
}

var Services []any

func init() {
	var ourServices = []any{
		NewService7(),
		NewService8(),
	}

	Services = make([]any, len(ourServices))

	for i, el := range ourServices {
		Services[len(ourServices)-i] = el
	}
}

func MoreServices() ServiceProvider {
	var provider ServiceProvider

	provider.TM2Service = &Service9{}
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
		p.HeresAnotherOne = &Service11{}
	default:
		if anyp, ok := p.(*any); ok {
			*anyp = &Service12{}
		}
	}
}
