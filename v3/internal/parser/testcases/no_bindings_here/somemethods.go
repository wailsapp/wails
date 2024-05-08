package nobindingshere

import "github.com/wailsapp/wails/v3/internal/parser/testcases/no_bindings_here/other"

// SomeMethods exports some methods.
type SomeMethods struct {
	other.OtherMethods
}

// LikeThisOne is an example method that does nothing.
func (SomeMethods) LikeThisOne(Person, Impersonator, HowDifferent[bool]) PrivatePerson {
	return PrivatePerson{}
}
