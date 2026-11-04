---
title: "Security Best Practices"
description: "Secure your Wails application"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Overview

Security is critical for desktop applications. Follow these practices to keep your application secure.

## Input Validation

### Always Validate

```go
func (s *UserService) CreateUser(email, password string) (*User, error) {
    // Validate email
    if !isValidEmail(email) {
        return nil, errors.New("invalid email")
    }
    
    // Validate password strength
    if len(password) < 8 {
        return nil, errors.New("password too short")
    }
    
    // Sanitise input
    email = strings.TrimSpace(email)
    email = html.EscapeString(email)
    
    // Continue...
}
```

### Sanitise HTML

```go
import "html"

func (s *Service) SaveComment(text string) error {
    // Escape HTML
    text = html.EscapeString(text)
    
    // Validate length
    if len(text) > 1000 {
        return errors.New("comment too long")
    }
    
    return s.db.Save(text)
}
```

## Authentication

### Secure Password Storage

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}

func verifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Session Management

```go
type Session struct {
    UserID    int
    Token     string
    ExpiresAt time.Time
}

func (a *AuthService) CreateSession(userID int) (*Session, error) {
    token := generateSecureToken()
    
    session := &Session{
        UserID:    userID,
        Token:     token,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    
    return session, a.saveSession(session)
}
```

## Data Protection

### Encrypt Sensitive Data

Use a fresh nonce for every message encrypted with a given key. On the Go toolchain required by Wails v3, `NewGCMWithRandomNonce` generates and prefixes the nonce automatically. Store the returned bytes intact; `decrypt` reads the nonce from that prefix. Keep the key separate from the ciphertext.

```go
import (
    "crypto/aes"
    "crypto/cipher"
)

func encrypt(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aead, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }
    return aead.Seal(nil, nil, data, nil), nil
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aead, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }
    return aead.Open(nil, nil, ciphertext, nil)
}
```

Use a suitable randomly generated AES key and respect the AEAD's per-key message limit. See the [Go cipher reference](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

### Secure Storage

```go
// Use OS keychain for sensitive data
import "github.com/zalando/go-keyring"

func saveAPIKey(key string) error {
    return keyring.Set("myapp", "api_key", key)
}

func getAPIKey() (string, error) {
    return keyring.Get("myapp", "api_key")
}
```

## Network Security

### Use HTTPS

```go
import (
    "errors"
    "net/http"
    "strings"
)

func makeAPICall(url string) (*http.Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return secureClient().Get(url)
}
```

### Verify Certificates

```go
import (
    "crypto/tls"
    "errors"
    "net/http"
    "time"
)

func secureClient() *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if req.URL.Scheme != "https" {
                return errors.New("redirect to non-HTTPS URL denied")
            }
            if len(via) >= 10 {
                return errors.New("too many redirects")
            }
            return nil
        },
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
        },
    }
}
```

## File Operations

### Confine File Access

Accept a path relative to an application-controlled directory. `os.OpenRoot` confines filesystem resolution to that directory, including symlinks. Rejecting `..` or checking a string prefix alone does not provide this boundary.

```go
import "os"

func readFile(allowedDir, name string) ([]byte, error) {
    root, err := os.OpenRoot(allowedDir)
    if err != nil {
        return nil, err
    }
    defer root.Close()
    return root.ReadFile(name)
}
```

For example, call `readFile(appDataDir, "notes/today.json")`. Choose `appDataDir` in trusted backend code; do not let the frontend supply the allowed root. A root limits filesystem traversal, not which files the user is authorised to read within it.

## Rate Limiting

```go
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func (r *RateLimiter) Allow(key string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    
    // Clean old requests
    var recent []time.Time
    for _, t := range r.requests[key] {
        if now.Sub(t) < r.window {
            recent = append(recent, t)
        }
    }
    
    if len(recent) >= r.limit {
        return false
    }
    
    r.requests[key] = append(recent, now)
    return true
}
```

## Best Practices

### ✅ Do

- Validate all input
- Use HTTPS for network calls
- Encrypt sensitive data
- Use secure password hashing
- Implement rate limiting
- Keep dependencies updated
- Log security events
- Use OS keychains

### ❌ Don't

- Don't trust user input
- Don't store passwords in plain text
- Don't hardcode secrets
- Don't skip certificate verification
- Don't expose sensitive data in logs
- Don't use weak encryption
- Don't ignore security updates

## Security Checklist

- [ ] All user input validated
- [ ] Passwords hashed with bcrypt
- [ ] Sensitive data encrypted
- [ ] HTTPS used for all network calls
- [ ] Rate limiting implemented
- [ ] File paths validated
- [ ] Dependencies up to date
- [ ] Security logging enabled
- [ ] Error messages don't leak info
- [ ] Code reviewed for vulnerabilities

## Next Steps

- [Architecture](/guides/architecture/) - Application architecture patterns
- [Best Practices](/features/bindings/best-practices/) - Bindings best practices
