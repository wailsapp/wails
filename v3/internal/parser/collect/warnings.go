package collect

import "github.com/pterm/pterm"

// complexWarning warns about unsupported use of complex types.
func complexWarning() {
	pterm.Warning.Println("complex types are not supported by encoding/json")
}

// chanWarning warns about unsupported use of channel types.
func chanWarning() {
	pterm.Warning.Println("channel types in models or bound method signatures are not supported")
}

// funcWarning warns about unsupported use of function types.
func funcWarning() {
	pterm.Warning.Println("function types in models or bound method signatures are not supported")
}

// genericWarning warns about unsupported use of generics.
// TODO: implement support for generics.
func genericWarning() {
	pterm.Warning.Println("generic types in models or bound method signatures are not fully supported yet")
}
