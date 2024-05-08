package analyse

import (
	"go/token"
	"go/types"

	"github.com/pterm/pterm"
)

// Result is an alias for the type of a single output from the analyser.
type Result = *types.TypeName

// reportResult checks whether a service type is valid and newly discovered;
// if both checks succeed, the type is added to the result set
// and reported to the consumer callback.
func (analyser *Analyser) reportResult(pkgi int, pos token.Pos, typ types.Type) {
	var named *types.Named

	switch t := types.Unalias(typ).(type) {
	case *types.Named:
		named = t
	case *types.Pointer:
		if elem, ok := types.Unalias(t.Elem()).(*types.Named); ok {
			named = elem
		} else {
			pterm.Warning.Printfln(
				"%s: ignoring service expression with non-named element type %s",
				analyser.pkgs[pkgi].Fset.Position(pos),
				t.Elem(),
			)
		}
	default:
		pterm.Warning.Printfln(
			"%s: ignoring service expression with non-named type %s",
			analyser.pkgs[pkgi].Fset.Position(pos),
			typ,
		)
	}

	if named == nil {
		return
	}

	if named.TypeParams() != nil {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with generic named type %s",
			analyser.pkgs[pkgi].Fset.Position(pos),
			typ,
		)
		return
	}

	// Retrieve type object.
	// If original type was an alias, use its object.
	result := named.Obj()
	if alias, ok := typ.(*types.Alias); ok {
		result = alias.Obj()
	}

	if analyser.found.Add(result) {
		if analyser.yield != nil {
			analyser.yield(result)
		}
	}
}
