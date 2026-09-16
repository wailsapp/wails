//go:build darwin && !ios && !server

#ifndef services_provider_manager_darwin_h
#define services_provider_manager_darwin_h

// wailsServicesProviderRegister makes `<name>:userData:error:` resolvable on
// NSApplication.servicesProvider, records the send types whose data is read
// from the pasteboard for that service, and calls NSUpdateDynamicServices.
// sendTypesJSON is a JSON array of type identifiers. Call on the main
// thread.
void wailsServicesProviderRegister(const char *name, const char *sendTypesJSON);

#endif /* services_provider_manager_darwin_h */
