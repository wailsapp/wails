package selfupdate

import (
	"fmt"
	"sync"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]ProviderFactory)
)

// ProviderFactory is a function that creates a new UpdateProvider instance.
type ProviderFactory func() UpdateProvider

// RegisterProvider registers a provider factory with the given name.
// This is typically called from init() functions in provider packages.
//
// Example:
//
//	func init() {
//	    selfupdate.RegisterProvider("github", func() selfupdate.UpdateProvider {
//	        return &GitHubProvider{}
//	    })
//	}
func RegisterProvider(name string, factory ProviderFactory) {
	providersMu.Lock()
	defer providersMu.Unlock()

	if factory == nil {
		panic("selfupdate: RegisterProvider factory is nil")
	}
	if _, exists := providers[name]; exists {
		panic(fmt.Sprintf("selfupdate: RegisterProvider called twice for provider %q", name))
	}
	providers[name] = factory
}

// GetProvider returns a new instance of the named provider.
// Returns an error if the provider is not registered.
func GetProvider(name string) (UpdateProvider, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()

	factory, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("selfupdate: unknown provider %q (available: %v)", name, AvailableProviders())
	}
	return factory(), nil
}

// AvailableProviders returns a list of all registered provider names.
func AvailableProviders() []string {
	providersMu.RLock()
	defer providersMu.RUnlock()

	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}

// HasProvider returns true if a provider with the given name is registered.
func HasProvider(name string) bool {
	providersMu.RLock()
	defer providersMu.RUnlock()

	_, ok := providers[name]
	return ok
}

// UnregisterProvider removes a provider from the registry.
// This is primarily useful for testing.
func UnregisterProvider(name string) {
	providersMu.Lock()
	defer providersMu.Unlock()

	delete(providers, name)
}

// UnregisterAllProviders removes all providers from the registry.
// This is primarily useful for testing.
func UnregisterAllProviders() {
	providersMu.Lock()
	defer providersMu.Unlock()

	providers = make(map[string]ProviderFactory)
}
