# HTTP-Only Bindings Configuration Examples

This directory contains examples showing how to configure HTTP-only bindings for different environments.

## Examples

### Development Configuration (`binding-configuration-dev.go`)

Shows configuration suitable for development:
- **Permissive CORS**: Allows all origins for easy frontend development
- **Shorter timeouts**: Quick feedback during development
- **Additional headers**: Includes development-specific headers
- **Streaming enabled**: For testing large payload handling

Key features:
```go
Bindings: application.BindingConfig{
    Timeout: 5 * time.Minute,
    CORS: application.CORSConfig{
        Enabled: true,
        AllowedOrigins: []string{}, // Empty = allow all in dev mode
        AllowedHeaders: []string{
            "Content-Type",
            "x-wails-client-id",
            "x-wails-window-name",
            "x-dev-token", // Development-specific
        },
        MaxAge: 1 * time.Hour,
    },
    EnableStreaming: true,
}
```

### Production Configuration (`binding-configuration-prod.go`)

Shows configuration suitable for production:
- **Strict CORS**: Only allows specific trusted origins
- **Longer timeouts**: Accommodates complex business logic
- **Minimal headers**: Only includes necessary headers
- **Security focused**: Includes authentication headers

Key features:
```go
Bindings: application.BindingConfig{
    Timeout: 10 * time.Minute,
    CORS: application.CORSConfig{
        Enabled: true,
        AllowedOrigins: []string{
            "https://myapp.com",
            "https://app.myapp.com",
            "https://*.myapp.com", // Wildcard for subdomains
        },
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{
            "Content-Type",
            "x-wails-client-id",
            "x-wails-window-name",
            "Authorization", // For authentication
        },
        MaxAge: 24 * time.Hour,
    },
    EnableStreaming: true,
}
```

## Configuration Options

### BindingConfig

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `Timeout` | `time.Duration` | Maximum time for binding execution | 10 minutes |
| `CORS` | `CORSConfig` | Cross-origin request configuration | See below |
| `EnableStreaming` | `bool` | Enable streaming for large responses | `false` |

### CORSConfig

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `Enabled` | `bool` | Enable CORS headers | `false` |
| `AllowedOrigins` | `[]string` | Allowed origins (supports wildcards) | `[]` |
| `AllowedMethods` | `[]string` | Allowed HTTP methods | `["GET", "POST", "OPTIONS"]` |
| `AllowedHeaders` | `[]string` | Allowed request headers | Standard Wails headers |
| `MaxAge` | `time.Duration` | Preflight cache duration | 24 hours |

## Security Considerations

### Development vs Production

- **Development**: Uses permissive settings for ease of development
- **Production**: Uses strict settings for security

### CORS Security

- In development with empty `AllowedOrigins`, all origins are allowed
- In production, explicitly specify trusted origins
- Use wildcards carefully (e.g., `https://*.myapp.com`)
- Never use `*` as an allowed origin in production

### Timeout Configuration

- Set appropriate timeouts based on your application's needs
- Longer timeouts for complex data processing
- Shorter timeouts for simple operations to prevent resource exhaustion

## Error Handling

The HTTP-only binding system returns proper HTTP status codes:

- `200 OK` - Successful binding call
- `400 Bad Request` - Invalid arguments (TypeError)
- `404 Not Found` - Unknown binding method (ReferenceError)
- `408 Request Timeout` - Binding execution timeout
- `500 Internal Server Error` - Runtime error in binding

## Migration from Hybrid System

When migrating from the hybrid callback system:

1. Update frontend to handle direct HTTP responses instead of callbacks
2. Configure CORS if using external URLs
3. Update error handling to work with HTTP status codes
4. Test timeout behavior with your application's workload

## Running the Examples

```bash
# Development example
go run binding-configuration-dev.go

# Production example  
go run binding-configuration-prod.go
```