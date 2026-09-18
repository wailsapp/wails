---
title: "Bonnes pratiques de sécurité"
description: "Sécurisez votre application Wails"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Vue d’ensemble

La sécurité est essentielle pour les applications de bureau. Suivez ces pratiques pour sécuriser votre application.

## Validation des entrées

### Toujours valider

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

### Nettoyer le HTML

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

## Authentification

### Stockage sécurisé des mots de passe

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

### Gestion des sessions

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

## Protection des données

### Chiffrer les données sensibles

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

### Stockage sécurisé

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

## Sécurité du réseau

### Utiliser HTTPS

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### Vérifier les certificats

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

## Opérations sur les fichiers

### Valider les chemins d’accès

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

## Limitation du débit

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

## Assainissement des journaux

Wails masque automatiquement les données sensibles dans les journaux IPC pour éviter la divulgation accidentelle de secrets. Cette fonctionnalité est **activée par défaut**, avec des paramètres par défaut adaptés.

### Protection par défaut

Les noms de champs suivants sont automatiquement masqués (recherche de sous-chaînes sans distinction de casse) :

- **Authentification**: `password`, `passwd`, `pwd`, `token`, `bearer`, `jwt`, `access_token`, `refresh_token`, `secret`, `apikey`, `api_key`, `auth`, `authorization`, `credential`
- **Cryptographie**: `private`, `privatekey`, `private_key`, `signing`, `encryption_key`
- **Session**: `session`, `sessionid`, `session_id`, `cookie`, `csrf`, `xsrf`

Les motifs suivants sont également détectés dans les valeurs :

- Jetons JWT (`eyJhbG...`)
- Jetons Bearer (`Bearer xxx`)
- Formats courants de clés API (`sk_live_xxx`, `pk_test_xxx`)

### Configuration personnalisée

Configurez l’assainissement avec toutes les options disponibles :

```go
app := application.New(application.Options{
    Name: "MyApp",
    SanitizeOptions: &application.SanitizeOptions{
        // RedactFields: additional field names to redact (merged with defaults)
        RedactFields: []string{"cardNumber", "cvv", "ssn"},

        // RedactPatterns: additional regex patterns to match values
        RedactPatterns: []*regexp.Regexp{
            regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`), // card numbers
            regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`), // SSN format
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
```

### Options de configuration

| Option | Description |
| --- | --- |
| `RedactFields` | Noms de champs supplémentaires à masquer (ajoutés aux valeurs par défaut) |
| `RedactPatterns` | Expressions régulières supplémentaires à rechercher dans les valeurs |
| `CustomSanitizeFunc` | Fonction de contrôle complet ; renvoyez `(value, true)` pour remplacer la valeur |
| `DisableDefaults` | Utiliser uniquement les champs et motifs explicitement spécifiés |
| `Replacement` | Texte de remplacement personnalisé (par défaut : `***`) |
| `Disabled` | Désactiver complètement l’assainissement (à utiliser avec prudence) |

### API publique d’assainissement

Utilisez l’assainisseur pour vos propres données :

```go
// Get the application's sanitizer
sanitizer := app.Sanitizer()

// Sanitize a map
cleanData := sanitizer.SanitizeMap(sensitiveData)

// Sanitize JSON
cleanJSON := sanitizer.SanitizeJSON(jsonBytes)
```

## Bonnes pratiques

### ✅ À faire

- Validez toutes les entrées
- Utilisez HTTPS pour les appels réseau
- Chiffrez les données sensibles
- Utilisez un hachage sécurisé des mots de passe
- Mettez en place une limitation du débit
- Maintenez les dépendances à jour
- Journalisez les événements de sécurité
- Utilisez les trousseaux de clés du système d’exploitation

### ❌ À ne pas faire

- Ne faites pas confiance aux entrées utilisateur
- Ne stockez pas les mots de passe en texte clair
- Ne codez pas les secrets en dur
- N’omettez pas la vérification des certificats
- N’exposez pas de données sensibles dans les journaux (Wails assainit les journaux IPC par défaut)
- N’utilisez pas de chiffrement faible
- N’ignorez pas les mises à jour de sécurité
- Ne désactivez pas l’assainissement des journaux en production

## Liste de contrôle de sécurité

- [ ] Toutes les entrées utilisateur sont validées
- [ ] Les mots de passe sont hachés avec bcrypt
- [ ] Les données sensibles sont chiffrées
- [ ] HTTPS est utilisé pour tous les appels réseau
- [ ] La limitation du débit est mise en place
- [ ] Les chemins d’accès aux fichiers sont validés
- [ ] Les dépendances sont à jour
- [ ] La journalisation des événements de sécurité est activée
- [ ] Les messages d’erreur ne divulguent aucune information
- [ ] Le code a fait l’objet d’une recherche de vulnérabilités
- [ ] Assainissement des journaux configuré pour les champs propres à l’application
- [ ] Champs sensibles non divulgués dans la journalisation personnalisée

## Étapes suivantes

- [Architecture](/guides/architecture/) - Modèles d’architecture d’application
- [Bonnes pratiques](/features/bindings/best-practices/) - Bonnes pratiques relatives aux bindings
