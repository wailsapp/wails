package changelog

import "testing"

func TestPlaceholderCaseVariants(t *testing.T) {
	validator := NewValidator()
	for _, token := range []string{"TODO", "todo", "ToDo", "FIXME", "fixme", "FixMe", "TBD", "tbd", "TbD"} {
		t.Run(token, func(t *testing.T) {
			content := token + ": describe the actual change"
			results := []ValidationResult{
				validator.ValidateContent("## Added\n- " + content + "\n"),
				validator.ValidateEntry(&ChangelogEntry{Added: []string{content}}),
			}
			for i, result := range results {
				if result.Valid || len(result.Errors) != 1 {
					t.Errorf("validation path %d: got %+v, want one placeholder error", i, result)
					continue
				}
				if got := result.Errors[0].Message; got != "bullet point contains placeholder text: '"+content+"'" {
					t.Errorf("validation path %d: error = %q", i, got)
				}
			}
		})
	}
}
