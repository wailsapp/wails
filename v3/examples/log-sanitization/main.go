package main

import (
	"fmt"
	"regexp"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name: "Log Sanitization Demo",

		// Configure log sanitization
		SanitizeOptions: &application.SanitizeOptions{
			// Add custom fields to redact (merged with defaults)
			RedactFields: []string{"cardNumber", "cvv", "ssn"},

			// Add custom patterns
			RedactPatterns: []*regexp.Regexp{
				regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
			},
		},
	})

	// Demonstrate the sanitizer API
	fmt.Println("\n=== Log Sanitization Demo ===\n")

	sanitizer := app.Sanitizer()

	// Test data with sensitive information
	testData := map[string]any{
		"username":   "john_doe",
		"password":   "super_secret_123",
		"token":      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig",
		"apiKey":     "sk_live_abcdefghij1234567890",
		"email":      "john@example.com",
		"cardNumber": "4111-1111-1111-1111", // custom field
		"cvv":        "123",                 // custom field
		"ssn":        "123-45-6789",         // custom field
	}

	fmt.Println("Original data:")
	for k, v := range testData {
		fmt.Printf("  %s: %v\n", k, v)
	}

	fmt.Println("\nSanitized data:")
	cleanData := sanitizer.SanitizeMap(testData)
	for k, v := range cleanData {
		fmt.Printf("  %s: %v\n", k, v)
	}

	// Demonstrate nested sanitization
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

	fmt.Println("\nNested original:")
	fmt.Printf("  %+v\n", nestedData)

	fmt.Println("\nNested sanitized:")
	cleanNested := sanitizer.SanitizeMap(nestedData)
	fmt.Printf("  %+v\n", cleanNested)
}
