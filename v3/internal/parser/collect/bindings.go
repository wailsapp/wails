package collect

import (
	"go/types"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/hash"
	"golang.org/x/tools/go/types/typeutil"
)

type (
	BoundTypeInfo struct {
		*TypeDefInfo
		Imports *ImportMap
		Methods []*BoundMethodInfo
	}

	BoundMethodInfo struct {
		*MethodInfo
		FQN     string
		ID      string
		Params  []*ParamInfo
		Results []types.Type
	}

	ParamInfo struct {
		Name     string
		Type     types.Type
		Blank    bool
		Variadic bool
	}
)

// BoundType collects and returns information about a type
// for which bindings have to be generated.
//
// If package loading or object lookup fails at any point, BoundType returns nil.
// Errors are printed directly to the pterm Error logger.
//
// BoundType is safe for concurrent use.
func (collector *Collector) BoundType(typ *types.TypeName) *BoundTypeInfo {
	// Collect package information.
	pkg := collector.Package(typ.Pkg().Path())
	if !pkg.Collect() {
		return nil
	}

	info := &BoundTypeInfo{
		TypeDefInfo: pkg.Types[typ.Name()],
		Imports:     NewImportMap(pkg),
	}

	// Check type def information.
	if info.TypeDefInfo == nil {
		pterm.Error.Printfln(
			"package %s: type %s not found; try cleaning the build cache (go clean -cache)",
			pkg.Path,
			typ.Name(),
		)
		return nil
	}

	// Recover real type behind alias:
	// this is required to obtain the right method set
	// and fully qualified method names.
	realType := typ
	if _, isAlias := typ.Type().(*types.Alias); isAlias {
		switch real := types.Unalias(typ.Type()).(type) {
		case *types.Named:
			realType = real.Obj()
		case *types.Pointer:
			realType = types.Unalias(real.Elem()).(*types.Named).Obj()
		}
	}

	// Compute intuitive method set (i.e. both pointer and non-pointer receiver).
	// Do not use a method set cache because
	//   - it would hurt concurrency (requires mutual exclusion),
	//   - it is only useful when the same type is queried many times;
	//     this may only happen here if some embedded types appear frequently,
	//     which should be far from average.
	mset := typeutil.IntuitiveMethodSet(realType.Type(), nil)

	info.Methods = make([]*BoundMethodInfo, 0, len(mset))
	for _, sel := range mset {
		if !sel.Obj().Exported() {
			// Ignore unexported methods
			continue
		}

		methodInfo := collector.BoundMethod(realType, info.Imports, sel.Obj().(*types.Func))
		if methodInfo == nil {
			return nil
		}

		info.Methods = append(info.Methods, methodInfo)
	}

	// Record generated bindings.
	if len(info.Methods) > 0 {
		pkg.AddBindings(typ)
	}

	return info
}

// typeError caches type corresponding to the Go error interface.
var typeError = types.Universe.Lookup("error").Type()

// BoundMethod collects and returns information about a method of a type
// for which bindings have to be generated.
//
// If package loading or object lookup fails at any point, BoundMethod returns nil.
// Errors are printed directly to the pterm Error logger.
//
// BoundMethod is safe for concurrent use.
func (collector *Collector) BoundMethod(typ *types.TypeName, imports *ImportMap, method *types.Func) *BoundMethodInfo {
	// Collect package information.
	pkg := collector.Package(method.Pkg().Path())
	if !pkg.Collect() {
		return nil
	}

	signature := method.Type().(*types.Signature)

	// Retrieve original receiver type.
	var recv *types.TypeName

	switch rtype := signature.Recv().Type().(type) {
	case *types.Named:
		recv = rtype.Obj()
	case *types.Pointer:
		recv = rtype.Elem().(*types.Named).Obj()
	}

	// Retrieve original receiver type information.
	recvInfo := pkg.Types[recv.Name()]
	if recvInfo == nil {
		pterm.Error.Printfln(
			"package %s: type %s not found; try cleaning the build cache (go clean -cache)",
			pkg.Path,
			recv.Name(),
		)
		return nil
	}

	// Compute fully qualified name.
	fqn := typ.Pkg().Path() + "." + typ.Name() + "." + method.Name()
	id, _ := hash.Fnv(fqn)

	info := &BoundMethodInfo{
		MethodInfo: recvInfo.Methods[method.Name()],
		FQN:        fqn,
		ID:         strconv.FormatUint(uint64(id), 10),
		Params:     make([]*ParamInfo, 0, signature.Params().Len()),
		Results:    make([]types.Type, 0, signature.Results().Len()),
	}

	// Check method information.
	if info.MethodInfo == nil {
		pterm.Error.Printfln(
			"package %s: method %s.%s not found; try cleaning the build cache (go clean -cache)",
			pkg.Path,
			recv.Name(),
			method.Name(),
		)
		return nil
	}

	// Collect parameters.
	for i := range signature.Params().Len() {
		param := signature.Params().At(i)

		if i == 0 {
			// Skip first parameter if it has context type.
			if named, ok := param.Type().(*types.Named); ok && named.Obj().Pkg().Path() == "context" && named.Obj().Name() == "Context" {
				continue
			}
		}

		// Record type dependencies.
		imports.AddType(collector, param.Type())

		// Record parameter.
		info.Params = append(info.Params, &ParamInfo{
			Name:  param.Name(),
			Type:  param.Type(),
			Blank: param.Name() == "" || param.Name() == "_",
		})
	}

	if signature.Variadic() {
		info.Params[len(info.Params)-1].Variadic = true
	}

	// Collect results.
	for i := range signature.Results().Len() {
		result := signature.Results().At(i)

		if types.Identical(result.Type(), typeError) {
			// Skip error results, they are thrown as exceptions
			continue
		}

		// Record type dependencies.
		imports.AddType(collector, result.Type())

		// Record result.
		info.Results = append(info.Results, result.Type())
	}

	return info
}
