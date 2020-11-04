package backendjs

import (
	"github.com/leaanthony/slicer"
)

// Parser is our Wails project parser
type Parser struct {
	Packages map[string]*Package

	// The variable used to store the Wails application
	// EG: app := wails.CreateApp()
	applicationVariable string

	// The name of the wails package that's imported
	// import "github.com/wailsapp/wails/v2"  -> wails
	// import mywails "github.com/wailsapp/wails/v2" -> mywails
	wailsPackageVariable string

	// A list of methods that returns structs to the Bind method
	// EG: app.Bind( newMyStruct() )
	structMethodsThatWereBound slicer.StringSlicer

	// A list of struct literals that were bound to the application
	// EG: app.Bind( &mystruct{} )
	structLiteralsThatWereBound slicer.StringSlicer

	// A list of struct pointer literals that were bound to the application
	// EG: app.Bind( &mystruct{} )
	structPointerLiteralsThatWereBound slicer.StringSlicer

	// A list of variables that were used for binding
	// Eg: myVar := &mystruct{}; app.Bind( myVar )
	variablesThatWereBound slicer.StringSlicer

	// A list of variables that were assigned using a function call
	// EG: myVar := newStruct()
	variablesThatWereAssignedByFunctions map[string]string

	// A map of variables that were assigned using a struct literal
	// EG: myVar := MyStruct{}
	variablesThatWereAssignedByStructLiterals map[string]string

	// Internal methods (WailsInit/WailsShutdown)
	internalMethods *slicer.StringSlicer

	// A list of functions that return struct pointers
	functionsThatReturnStructPointers map[string]string

	// A list of functions that return structs
	functionsThatReturnStructs map[string]string
}

// NewParser creates a new Wails Project parser
func NewParser() *Parser {
	return &Parser{
		Packages:                                  make(map[string]*Package),
		variablesThatWereAssignedByFunctions:      make(map[string]string),
		variablesThatWereAssignedByStructLiterals: make(map[string]string),
		functionsThatReturnStructPointers:         make(map[string]string),
		functionsThatReturnStructs:                make(map[string]string),
		internalMethods:                           slicer.String([]string{"WailsInit", "WailsShutdown"}),
	}
}

func (p *Parser) resolve() error {
	return nil
}
