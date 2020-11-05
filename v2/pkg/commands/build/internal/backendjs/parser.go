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

	// Internal methods (WailsInit/WailsShutdown)
	internalMethods *slicer.StringSlicer
}

// NewParser creates a new Wails Project parser
func NewParser() *Parser {
	return &Parser{
		Packages:        make(map[string]*Package),
		internalMethods: slicer.String([]string{"WailsInit", "WailsShutdown"}),
	}
}

func (p *Parser) resolve() error {

	// Resolve bound structs
	err := p.resolveBoundStructs()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) resolveBoundStructs() error {

	// Resolve Struct Literals
	p.resolveBoundStructLiterals()

	// Resolve Struct Pointer Literals
	p.resolveBoundStructPointerLiterals()

	// Resolve functions that were bound
	// EG: app.Bind( newBasic() )
	p.resolveBoundFunctions()

	// Resolve variables that were bound
	p.resolveBoundVariables()

	return nil
}
func (p *Parser) resolveBoundStructLiterals() {

	// Resolve struct literals in each package
	for _, pkg := range p.Packages {
		pkg.resolveBoundStructLiterals()
	}
}

func (p *Parser) resolveBoundStructPointerLiterals() {

	// Resolve struct pointer literals
	for _, pkg := range p.Packages {
		pkg.resolveBoundStructPointerLiterals()
	}
}

func (p *Parser) resolveBoundFunctions() {

	// Loop over packages
	for _, pkg := range p.Packages {

		// Iterate over the method names
		pkg.structMethodsThatWereBound.Each(func(functionName string) {
			// Resolve each method name
			structName := p.resolveFunctionReturnType(pkg, functionName)

			strct := pkg.Structs[structName]
			if strct == nil {
				println("WARNING: Unable to find definition for struct", structName)
				return
			}
			strct.IsBound = true
		})

	}
}

// resolveFunctionReturnType gets the return type for the given package/function name combination
func (p *Parser) resolveFunctionReturnType(pkg *Package, functionName string) string {
	structName := pkg.functionsThatReturnStructPointers[functionName]
	if structName == "" {
		structName = pkg.functionsThatReturnStructs[functionName]
	}
	if structName == "" {
		println("WARNING: Unable to resolve bound function", functionName, "in package", pkg.Name)
	}
	return structName
}

func (p *Parser) markStructAsBound(pkg *Package, structName string) {
	strct := pkg.Structs[structName]
	if strct == nil {
		println("WARNING: Unable to find definition for struct", structName)
	}
	println("Found bound struct:", strct.Name)
	strct.IsBound = true
}

func (p *Parser) resolveBoundVariables() {

	for _, pkg := range p.Packages {

		// Iterate over the method names
		pkg.variablesThatWereBound.Each(func(variableName string) {
			println("Resolving variable: ", variableName)

			var structName string

			// Resolve each method name
			funcName := pkg.variablesThatWereAssignedByFunctions[variableName]
			if funcName != "" {
				// Found function name - resolve Function return type
				structName = p.resolveFunctionReturnType(pkg, funcName)
			}

			// If we couldn't resolve to a function, then let's try struct literals
			if structName == "" {
				funcName = pkg.variablesThatWereAssignedByStructLiterals[variableName]
				if funcName != "" {
					// Found function name - resolve Function return type
					structName = p.resolveFunctionReturnType(pkg, funcName)
				}
			}

			// Look for the variable as an external literal reference
			if structName == "" {
				sn := pkg.variablesThatWereAssignedByExternalStructLiterals[variableName]
				if sn != nil {
					pkg := p.Packages[sn.Package]
					if pkg != nil {
						p.markStructAsBound(pkg, sn.Name)
						return
					}

				}
			}

			if structName == "" {
				println("WARNING: Unable to resolve bound variable", variableName, "in package", pkg.Name)
				return
			}

			p.markStructAsBound(pkg, structName)
		})
	}
}

func (p *Parser) bindStructByStructName(sn *StructName) {
	// Get package
	pkg := p.Packages[sn.Package]
	if pkg == nil {
		// Ignore, it will get picked up by the compiler
		return
	}

	strct := pkg.Structs[sn.Name]
	if strct == nil {
		// Ignore, it will get picked up by the compiler
		return
	}

	println("Found bound Struct:", sn.ToString())
	strct.IsBound = true
}

func (p *Parser) getOrCreatePackage(name string) *Package {
	result := p.Packages[name]
	if result == nil {
		result = newPackage(name)
		p.Packages[name] = result
	}
	return result
}
