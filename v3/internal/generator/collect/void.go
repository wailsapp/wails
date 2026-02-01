package collect

import (
	"go/types"
)

// IsVoidAlias returns true when the given type or object is the application.Void named type that stands in for the void TS type.
func (collector *Collector) IsVoidAlias(typOrObj any) bool {
	var obj types.Object
	switch to := typOrObj.(type) {
	case types.Object:
		obj = to
	case interface{ Obj() *types.TypeName }:
		obj = to.Obj()
	default:
		return false
	}

	if vt := collector.appVoidType.Load(); vt != nil && obj == vt {
		return true
	} else if vt == nil && obj.Name() == "Void" && obj.Pkg().Path() == collector.systemPaths.ApplicationPackage { // Check name before package to fail fast
		// Cache void alias for faster checking
		collector.appVoidType.Store(obj)
		return true
	}

	return false
}
