package collect

import (
	"go/types"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/internal/generator/config"
)

// SpecialTypeCache caches types that are handled as special cases in the generation process.
type SpecialTypeCache struct {
	// timeTimeType caches the time.Time type, which may be rendered as a string, Date or Timestamp.
	timeTimeType atomic.Value
	// appVoidType caches the application.Void named type that stands in for the void TS type.
	appVoidType atomic.Value

	systemPaths *config.SystemPaths
}

// NewSpecialTypeCache initialises a new SpecialTypeCache instance for the given system paths.
func NewSpecialTypeCache(systemPaths *config.SystemPaths) SpecialTypeCache {
	return SpecialTypeCache{
		systemPaths: systemPaths,
	}
}

// IsSpecialType returns true when the given type or object is handled as a special case in the generation process.
func (specialTypes *SpecialTypeCache) IsSpecialType(typOrObj any) bool {
	return specialTypes.IsStdTime(typOrObj) ||
		specialTypes.IsVoidAlias(typOrObj)
}

// IsStdTime returns true when the given type or object is the time.Time type.
func (specialTypes *SpecialTypeCache) IsStdTime(typOrObj any) bool {
	var obj types.Object
	switch to := typOrObj.(type) {
	case types.Object:
		obj = to
	case interface{ Obj() *types.TypeName }:
		obj = to.Obj()
	default:
		return false
	}

	if vt := specialTypes.timeTimeType.Load(); vt != nil && obj == vt {
		return true
	} else if vt == nil && obj.Name() == "Time" && obj.Pkg() != nil && obj.Pkg().Path() == specialTypes.systemPaths.TimePackage { // Check name before package to fail fast
		// Cache time.Time type for faster checking
		specialTypes.timeTimeType.Store(obj)
		return true
	}

	return false
}

// IsVoidAlias returns true when the given type or object is the application.Void named type that stands in for the void TS type.
func (specialTypes *SpecialTypeCache) IsVoidAlias(typOrObj any) bool {
	var obj types.Object
	switch to := typOrObj.(type) {
	case types.Object:
		obj = to
	case interface{ Obj() *types.TypeName }:
		obj = to.Obj()
	default:
		return false
	}

	if vt := specialTypes.appVoidType.Load(); vt != nil && obj == vt {
		return true
	} else if vt == nil && obj.Name() == "Void" && obj.Pkg() != nil && obj.Pkg().Path() == specialTypes.systemPaths.ApplicationPackage { // Check name before package to fail fast
		// Cache void alias for faster checking
		specialTypes.appVoidType.Store(obj)
		return true
	}

	return false
}
