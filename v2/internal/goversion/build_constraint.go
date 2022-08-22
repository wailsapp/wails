//go:build !go1.18
// +build !go1.18

package goversion

const MinGoVersionRequired = "You need Go " + MinRequirement + " or newer to compile this program"

func init() {
	MinGoVersionRequired
}
