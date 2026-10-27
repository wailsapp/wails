---
title: "Bewährte Sicherheitspraktiken"
description: "Schützen Sie Ihre Wails-Anwendung"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Überblick

Sicherheit ist für Desktopanwendungen unverzichtbar. Befolgen Sie diese Praktiken, um Ihre Anwendung zu schützen.

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

Verwenden Sie für jede Nachricht, die mit einem bestimmten Schlüssel verschlüsselt wird, einen neuen Nonce. In der von Wails benötigten Go-Toolchain v3 erzeugt `NewGCMWithRandomNonce` den Nonce automatisch und stellt ihn den Daten voran. Speichern Sie die zurückgegebenen Bytes unverändert; `decrypt` liest den Nonce aus diesem Präfix. Bewahren Sie den Schlüssel getrennt vom Chiffretext auf.

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

Verwenden Sie einen geeigneten, zufällig erzeugten AES-Schlüssel und beachten Sie die maximale Anzahl von Nachrichten pro Schlüssel für das AEAD-Verfahren. Weitere Informationen finden Sie in der [Go-Referenz zu Chiffrierverfahren](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

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

### Zertifikate überprüfen

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

## Dateioperationen

### Dateizugriff beschränken

Akzeptieren Sie einen Pfad relativ zu einem von der Anwendung kontrollierten Verzeichnis. `os.OpenRoot` beschränkt die Auflösung von Dateisystempfaden einschließlich symbolischer Links auf dieses Verzeichnis. Allein das Zurückweisen von `..` oder das Prüfen eines Zeichenkettenpräfixes stellt diese Begrenzung nicht sicher.

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

Rufen Sie beispielsweise `readFile(appDataDir, "notes/today.json")` auf. Legen Sie `appDataDir` in vertrauenswürdigem Backend-Code fest; lassen Sie nicht zu, dass das Frontend das zulässige Stammverzeichnis vorgibt. Ein Stammverzeichnis beschränkt die Navigation im Dateisystem, nicht jedoch, welche darin enthaltenen Dateien der Benutzer lesen darf.

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

### ✅ Empfohlen

- Alle Eingaben validieren
- HTTPS für Netzwerkaufrufe verwenden
- Sensible Daten verschlüsseln
- Sichere Passwort-Hashverfahren verwenden
- Ratenbegrenzung implementieren
- Abhängigkeiten aktuell halten
- Sicherheitsereignisse protokollieren
- Schlüsselbunde des Betriebssystems verwenden

### ❌ Nicht empfohlen

- Benutzereingaben nicht vertrauen
- Passwörter nicht im Klartext speichern
- Geheimnisse nicht fest im Code hinterlegen
- Zertifikatsprüfung nicht überspringen
- Sensible Daten nicht in Protokollen offenlegen
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

- [Architektur](/guides/architecture/) – Architekturmuster für Anwendungen
- [Bewährte Praktiken](/features/bindings/best-practices/) – Bewährte Praktiken für Bindings
