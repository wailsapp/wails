package manifest

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// FormatValidationDiagnostics renders manifest validation failures with their
// source line and caret. It returns false for errors outside the manifest
// validation domain so callers can retain their normal error presentation.
func FormatValidationDiagnostics(err error) (string, bool) {
	validations := collectValidationErrors(err)
	if len(validations) == 0 {
		return "", false
	}
	var output strings.Builder
	for index, validation := range validations {
		if index != 0 {
			output.WriteString("\n\n")
		}
		formatValidationDiagnostic(&output, validation)
	}
	return output.String(), true
}

func collectValidationErrors(err error) []*ValidationError {
	var result []*ValidationError
	var visit func(error)
	visit = func(current error) {
		if current == nil {
			return
		}
		if joined, ok := current.(interface{ Unwrap() []error }); ok {
			for _, child := range joined.Unwrap() {
				visit(child)
			}
			return
		}
		var validation *ValidationError
		if errors.As(current, &validation) {
			result = append(result, validation)
		}
	}
	visit(err)
	return result
}

func formatValidationDiagnostic(output *strings.Builder, validation *ValidationError) {
	location := validation.Range.Filename
	if location == "" {
		location = Filename
	}
	if validation.Range.StartLine > 0 {
		fmt.Fprintf(output, "%s:%d:%d\n", location, validation.Range.StartLine, validation.Range.StartColumn)
	} else {
		fmt.Fprintln(output, location)
	}
	if validation.Range.Filename != "" && validation.Range.StartLine > 0 {
		if raw, err := os.ReadFile(validation.Range.Filename); err == nil {
			lines := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
			lineIndex := validation.Range.StartLine - 1
			if lineIndex >= 0 && lineIndex < len(lines) {
				width := len(fmt.Sprint(validation.Range.StartLine))
				fmt.Fprintf(output, "%*d | %s\n", width, validation.Range.StartLine, lines[lineIndex])
				start := max(validation.Range.StartColumn-1, 0)
				end := validation.Range.EndColumn - 1
				if validation.Range.EndLine != validation.Range.StartLine || end <= start {
					end = start + 1
				}
				fmt.Fprintf(output, "%*s | %s%s\n", width, "", strings.Repeat(" ", start), strings.Repeat("^", max(end-start, 1)))
			}
		}
	}
	fmt.Fprintf(output, "%s: %s", validation.Field, validation.Detail)
	if suggestion := validationSuggestion(validation); suggestion != "" {
		fmt.Fprintf(output, "\nHint: %s", suggestion)
	}
}

func validationSuggestion(validation *ValidationError) string {
	if strings.Contains(validation.Detail, "Unsupported argument") || strings.Contains(validation.Detail, "Unsupported block") {
		return "remove the unsupported field or compare this section with `wails3 eject` output"
	}
	switch validation.Field {
	case "frontend.install":
		return "use npm, pnpm, yarn, or bun"
	case "frontend.bindings.time_type":
		return `use "string", or use "Date" only when interfaces = false`
	}
	if strings.Contains(validation.Field, ".toolchain") {
		return "use auto, native, zig, or docker"
	}
	return ""
}
