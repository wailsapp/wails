package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name: "Log Sanitization Demo",

		SanitizeOptions: &application.SanitizeOptions{
			// RedactFields: additional field names to redact (merged with defaults)
			RedactFields: []string{"cardNumber", "cvv", "ssn"},

			// RedactPatterns: additional regex patterns to match values
			RedactPatterns: []*regexp.Regexp{
				regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`), // card numbers
			},

			// CustomSanitizeFunc: full control - return (value, true) to override
			CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
				// Custom handling for specific paths
				if strings.HasPrefix(path, "payment.") && key != "amount" {
					return "[PAYMENT_REDACTED]", true
				}
				return nil, false // fall through to default logic
			},

			// Replacement: custom replacement string (default: "***")
			Replacement: "[REDACTED]",

			// DisableDefaults: if true, only use explicitly specified fields/patterns
			// DisableDefaults: false,

			// Disabled: completely disable sanitization
			// Disabled: false,
		},
	})

	sanitizer := app.Sanitizer()

	fmt.Println("\n=== Log Sanitization Demo ===")

	// Test default field redaction
	fmt.Println("\n--- Default Fields ---")
	defaultData := map[string]any{
		"username": "john_doe",
		"password": "super_secret_123",
		"token":    "abc123",
		"apiKey":   "sk_live_xyz",
		"email":    "john@example.com",
	}
	fmt.Println("Original:", defaultData)
	fmt.Println("Sanitized:", sanitizer.SanitizeMap(defaultData))

	// Test custom field redaction
	fmt.Println("\n--- Custom Fields ---")
	customData := map[string]any{
		"cardNumber": "4111-1111-1111-1111",
		"cvv":        "123",
		"ssn":        "123-45-6789",
		"name":       "John Doe",
	}
	fmt.Println("Original:", customData)
	fmt.Println("Sanitized:", sanitizer.SanitizeMap(customData))

	// Test pattern matching
	fmt.Println("\n--- Pattern Matching ---")
	patternData := map[string]any{
		"data":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig", // JWT
		"header":  "Bearer abc123xyz",
		"message": "Hello world",
	}
	fmt.Println("Original:", patternData)
	fmt.Println("Sanitized:", sanitizer.SanitizeMap(patternData))

	// Test CustomSanitizeFunc
	fmt.Println("\n--- CustomSanitizeFunc ---")
	paymentData := map[string]any{
		"payment": map[string]any{
			"cardNumber": "4111111111111111",
			"amount":     99.99,
			"currency":   "USD",
		},
	}
	fmt.Println("Original:", paymentData)
	fmt.Println("Sanitized:", sanitizer.SanitizeMap(paymentData))

	// Test nested structures
	fmt.Println("\n--- Nested Structures ---")
	nestedData := map[string]any{
		"user": map[string]any{
			"name":     "Jane",
			"password": "secret123",
			"settings": map[string]any{
				"theme":      "dark",
				"auth_token": "bearer_xyz789",
			},
		},
	}
	fmt.Println("Original:", nestedData)
	fmt.Println("Sanitized:", sanitizer.SanitizeMap(nestedData))

	// Test JSON sanitization
	fmt.Println("\n--- JSON Sanitization ---")
	jsonData := []byte(`{"user":"john","password":"secret","token":"abc123"}`)
	fmt.Println("Original:", string(jsonData))
	fmt.Println("Sanitized:", string(sanitizer.SanitizeJSON(jsonData)))
}
