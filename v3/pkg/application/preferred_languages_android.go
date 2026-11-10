//go:build android && !server

package application

import "strings"

// platformPreferredLanguages asks WailsBridge.getPreferredLanguages for the
// list in Settings › System › Languages, as comma-separated tags.
//
// An app whose WailsBridge.java was copied from a template older than that
// method gets an empty list rather than a crash: the bridge call clears the
// NoSuchMethodError and reports no result. Builds without cgo have no bridge
// and get an empty list too.
func platformPreferredLanguages() []string {
	tags, ok := androidBridgeString("getPreferredLanguages")
	if !ok || tags == "" {
		return nil
	}
	return strings.Split(tags, ",")
}
