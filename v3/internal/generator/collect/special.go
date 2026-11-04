package collect

import (
	"go/types"
	"strings"
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

// IsJSONv2CompatAlias returns true when the given type is an alias
// declared in package encoding/json whose right-hand side
// is a named type from encoding/json/v2, encoding/json/jsontext
// or another encoding/json subpackage.
//
// Starting with Go 1.27 (GOEXPERIMENT=jsonv2, enabled by default),
// some encoding/json types such as Marshaler, Unmarshaler and RawMessage
// are declared as aliases of their encoding/json/v2 counterparts.
// To keep generated bindings stable across Go versions,
// such aliases are treated as if they were the named types
// they used to be: they are modelled after their underlying type
// and their targets are never emitted as separate models.
func IsJSONv2CompatAlias(typ types.Type) bool {
	alias, ok := typ.(*types.Alias)
	if !ok || alias.TypeArgs() != nil {
		return false
	}

	obj := alias.Obj()
	if obj.Pkg() == nil || obj.Pkg().Path() != "encoding/json" {
		return false
	}

	named, ok := types.Unalias(alias).(*types.Named)
	if !ok || named.Obj().Pkg() == nil {
		return false
	}

	return strings.HasPrefix(named.Obj().Pkg().Path(), "encoding/json/")
}

// UnaliasJSONv2Compat is like [types.Unalias],
// except that it stops at the first alias in the chain
// that satisfies [IsJSONv2CompatAlias].
func UnaliasJSONv2Compat(typ types.Type) types.Type {
	for {
		if IsJSONv2CompatAlias(typ) {
			return typ
		}

		alias, ok := typ.(*types.Alias)
		if !ok {
			return typ
		}

		if alias.TypeArgs() != nil {
			// Do not walk instantiated generic aliases:
			// their right-hand side is expressed in terms of type parameters.
			return types.Unalias(typ)
		}

		typ = alias.Rhs()
	}
}
