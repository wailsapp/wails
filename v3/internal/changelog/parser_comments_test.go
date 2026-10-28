package changelog

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseContentIgnoresMultilineComments(t *testing.T) {
	content := `## Fixed
- Visible first fix
<!--
- Hidden fix
## Added
- Hidden feature
-->
- Visible second fix
<!-- - Hidden single-line fix -->
## Added
- Visible feature
`
	entry, err := NewParser().ParseContent(strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(entry.Fixed, []string{"Visible first fix", "Visible second fix"}) {
		t.Errorf("Fixed = %q", entry.Fixed)
	}
	if !reflect.DeepEqual(entry.Added, []string{"Visible feature"}) {
		t.Errorf("Added = %q", entry.Added)
	}
}

func TestValidateContentIgnoresMultilineComments(t *testing.T) {
	content := `## Fixed
- Visible fix
<!--
## Invalid section
-
-->
- Another visible fix
`
	result := NewValidator().ValidateContent(content)
	if !result.Valid || len(result.Errors) != 0 {
		t.Errorf("comment contents produced validation errors: %v", result.Errors)
	}
}
