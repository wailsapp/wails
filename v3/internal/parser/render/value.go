package render

import (
	"fmt"
	"math/big"
	"strconv"
	"text/template"
)

// RenderValue renders a Go constant value to its Javascript representation.
func RenderValue(value any) string {
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
