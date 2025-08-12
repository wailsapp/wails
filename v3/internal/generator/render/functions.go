package render

import (
	"fmt"
	"go/types"
	"math/big"
	"strconv"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// tmplFunctions holds a map of utility functions
// that should be available in every template.
var tmplFunctions = template.FuncMap{
	"fixext":     fixext,
	"hasdoc":     hasdoc,
	"isjsdocid":  isjsdocid,
	"isjsdocobj": isjsdocobj,
	"istpalias":  istpalias,
	"jsdoc":      jsdoc,
	"jsdocline":  jsdocline,
	"jsid":       jsid,
	"jsimport":   jsimport,
	"jsparam":    jsparam,
	"jsvalue":    jsvalue,
	"modelinfo":  modelinfo,
	"typeparam":  typeparam,
	"unalias":    types.Unalias,
}

// fixext replaces a *.ts extension with *.js in the given string.
// This is necessary to allow emitting javascript with the Typescript compiler.
func fixext(path string) string {
	if strings.HasSuffix(path, ".ts") {
		return path[:len(path)-3] + ".js"
	} else {
		return path
	}
}

// jsimport formats an external import name
// by joining the name with its occurrence index.
// Names are modified even when the index is 0
// to avoid collisions with Go identifiers.
func jsimport(info collect.ImportInfo) string {
	return fmt.Sprintf("%s$%d", info.Name, info.Index)
}

// jsparam renders the JS name of a parameter.
// Blank parameters are replaced with a dollar sign followed by the given index.
// Non-blank parameters are escaped by [jsid].
func jsparam(index int, param *collect.ParamInfo) string {
	if param.Blank {
		return "$" + strconv.Itoa(index)
	} else {
		return jsid(param.Name)
	}
}

// typeparam renders the TS name of a type parameter.
// Blank parameters are replaced with a double dollar sign
// followed by the given index.
// Non-blank parameters are escaped with jsid.
func typeparam(index int, param string) string {
	if param == "" || param == "_" {
		return "$$" + strconv.Itoa(index)
	} else {
		return jsid(param)
	}
}

// jsvalue renders a Go constant value to its Javascript representation.
func jsvalue(value any) string {
	switch v := value.(type) {
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case string:
		return fmt.Sprintf(`"%s"`, template.JSEscapeString(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case *big.Int:
		return v.String()
	case *big.Float:
		return v.Text('e', -1)
	case *big.Rat:
		return v.RatString()
	}

	// Fall back to undefined.
	return "(void(0))"
}

// istpalias determines whether typ is an alias
// that when uninstantiated resolves to a typeparam.
func istpalias(typ types.Type) bool {
	if alias, ok := typ.(*types.Alias); ok {
		return collect.IsTypeParam(alias.Origin())
	}

	return false
}
