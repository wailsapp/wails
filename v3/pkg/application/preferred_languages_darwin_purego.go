//go:build darwin && purego && !ios && !server

package application

// platformPreferredLanguages returns +[NSLocale preferredLanguages], the list
// in System Settings › General › Language & Region, through the Objective-C
// runtime helpers in darwin_purego_cocoa.go. It mirrors the cgo version in
// preferred_languages_darwin.go.
func platformPreferredLanguages() []string {
	var langs []string
	withAutoreleasePool(func() {
		list := class("NSLocale").send("preferredLanguages")
		if list.isNil() {
			return
		}
		count := get[uint](list, "count")
		for i := uint(0); i < count; i++ {
			langs = append(langs, list.send("objectAtIndex:", i).string())
		}
	})
	return langs
}
