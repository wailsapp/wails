---
title: "Bewährte Sicherheitspraktiken"
description: "Sichern Sie Ihre Wails-Anwendung ab"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Übersicht

Sicherheit ist für Desktopanwendungen von entscheidender Bedeutung. Befolgen Sie diese Praktiken, um Ihre Anwendung zu schützen.

## Eingabevalidierung

### Immer validieren

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

### HTML bereinigen

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

## Authentifizierung

### Sichere Passwortspeicherung

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

### Sitzungsverwaltung

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

## Datenschutz

### Sensible Daten verschlüsseln

```go
import "crypto/aes"
import "crypto/cipher"

func encrypt(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    return gcm.Seal(nonce, nonce, data, nil), nil
}
```

### Sichere Speicherung

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

## Netzwerksicherheit

### HTTPS verwenden

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### Zertifikate überprüfen

```go
import "crypto/tls"

func secureClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
        },
    }
}
```

## Dateioperationen

### Pfade validieren

```go
func readFile(path string) ([]byte, error) {
    // Prevent path traversal
    if strings.Contains(path, "..") {
        return nil, errors.New("invalid path")
    }
    
    // Check file exists in allowed directory
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }
    
    if !strings.HasPrefix(absPath, allowedDir) {
        return nil, errors.New("access denied")
    }
    
    return os.ReadFile(absPath)
}
```

## Ratenbegrenzung

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

## Bewährte Praktiken

### ✅ Das tun

- Alle Eingaben validieren
- HTTPS für Netzwerkaufrufe verwenden
- Sensible Daten verschlüsseln
- Sicheres Passwort-Hashing verwenden
- Ratenbegrenzung implementieren
- Abhängigkeiten aktuell halten
- Sicherheitsereignisse protokollieren
- Schlüsselbünde des Betriebssystems verwenden

### ❌ Das nicht tun

- Benutzereingaben nicht vertrauen
- Passwörter nicht im Klartext speichern
- Geheimnisse nicht fest im Code hinterlegen
- Die Zertifikatsprüfung nicht überspringen
- Keine sensiblen Daten in Protokollen offenlegen
- Keine schwache Verschlüsselung verwenden
- Sicherheitsupdates nicht ignorieren

## Sicherheitscheckliste

- [ ] Alle Benutzereingaben validiert
- [ ] Passwörter mit bcrypt gehasht
- [ ] Sensible Daten verschlüsselt
- [ ] HTTPS für alle Netzwerkaufrufe verwendet
- [ ] Ratenbegrenzung implementiert
- [ ] Dateipfade validiert
- [ ] Abhängigkeiten auf dem neuesten Stand
- [ ] Sicherheitsprotokollierung aktiviert
- [ ] Fehlermeldungen geben keine Informationen preis
- [ ] Code auf Schwachstellen geprüft

## Nächste Schritte

- [Architektur](/guides/architecture/) – Anwendungarchitekturmuster
- [Bewährte Praktiken](/features/bindings/best-practices/) – Bewährte Praktiken für Bindings
