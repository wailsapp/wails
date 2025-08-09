package nobindingshere

import "github.com/wailsapp/wails/v3/internal/generator/testcases/no_bindings_here/other"

// SomeMethods exports some methods.
type SomeMethods struct {
	other.OtherMethods
}

// LikeThisOne is an example method that does nothing.
func (SomeMethods) LikeThisOne() (_ Person, _ HowDifferent[bool], _ PrivatePerson) {
	return
}
