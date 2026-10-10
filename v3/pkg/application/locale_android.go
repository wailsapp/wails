//go:build android

package application

// SystemLocale returns the device's configured locale as a BCP-47 language tag
// (e.g. "nb-NO", "en-US"). On Android this calls Locale.getDefault().toLanguageTag()
// via the WailsBridge.
func SystemLocale() string {
	s, _ := androidBridgeString("getLocale")
	if s == "" {
		return "en"
	}
	return s
}
