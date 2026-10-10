//go:build ios

package application

// SystemLocale returns the device's configured locale as a BCP-47 language tag
// (e.g. "nb-NO", "en-US"). Delegates to the iOS mobile manager which calls
// ios_system_locale() via the CGO preamble in mobile_features_ios.go.
func SystemLocale() string {
	return IOS.SystemLocale()
}
