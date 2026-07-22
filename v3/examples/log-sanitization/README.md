# Log Sanitization Example

Demonstrates Wails' automatic log sanitization feature which redacts sensitive data to prevent accidental exposure of secrets.

## Running

```bash
go run .
```

Watch the terminal output - it shows original vs sanitized data, demonstrating automatic redaction.

## What Gets Redacted

**Default fields** (case-insensitive, substring match):
- `password`, `passwd`, `pwd`, `token`, `secret`, `apikey`, `api_key`
- `auth`, `authorization`, `credential`, `bearer`, `jwt`
- `private`, `privatekey`, `session`, `cookie`, `csrf`

**Default patterns** (detected in values):
- JWT tokens: `eyJhbG...`
- Bearer tokens: `Bearer xxx`
- API keys: `sk_live_xxx`, `pk_test_xxx`

## Configuration

```go
app := application.New(application.Options{
    SanitizeOptions: &application.SanitizeOptions{
        RedactFields:   []string{"cardNumber", "ssn"},
        RedactPatterns: []*regexp.Regexp{...},
        Replacement:    "[REDACTED]", // default: "***"
    },
})
```

## Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | Working |
| Linux    | Working |
