---
title: "Bonnes pratiques de sécurité"
description: "Sécurisez votre application Wails"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Vue d’ensemble

La sécurité est essentielle pour les applications de bureau. Suivez ces pratiques pour sécuriser votre application.

## Validation des entrées

### Validez systématiquement

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

### Assainissez le HTML

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

### Chiffrez les données sensibles

Utilisez un nonce inédit pour chaque message chiffré avec une clé donnée. Avec la chaîne d’outils Go requise par Wails v3, `NewGCMWithRandomNonce` génère automatiquement le nonce et l’ajoute en préfixe. Stockez tels quels les octets renvoyés ; `decrypt` lit le nonce dans ce préfixe. Conservez la clé séparément du texte chiffré.

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

Utilisez une clé AES appropriée générée aléatoirement et respectez la limite de messages par clé de l’AEAD. Consultez la [référence sur le chiffrement en Go](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

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

### Utilisez HTTPS

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

### Vérifiez les certificats

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

## Opérations sur les fichiers

### Restreignez l’accès aux fichiers

Acceptez un chemin relatif à un répertoire contrôlé par l’application. `os.OpenRoot` restreint à ce répertoire la résolution des chemins du système de fichiers, y compris celle des liens symboliques. Rejeter `..` ou vérifier uniquement le préfixe d’une chaîne ne permet pas d’établir cette limite.

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

Par exemple, appelez `readFile(appDataDir, "notes/today.json")`. Choisissez `appDataDir` dans le code backend de confiance ; ne laissez pas le frontend fournir la racine autorisée. Une racine limite le parcours du système de fichiers, mais ne détermine pas quels fichiers l’utilisateur est autorisé à lire dans cette racine.

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

## Bonnes pratiques

### ✅ À faire

- Validez toutes les entrées
- Utilisez HTTPS pour les appels réseau
- Chiffrez les données sensibles
- Utilisez un hachage sécurisé des mots de passe
- Mettez en œuvre une limitation du débit
- Maintenez les dépendances à jour
- Journalisez les événements de sécurité
- Utilisez les trousseaux de clés du système d’exploitation

### ❌ À ne pas faire

- Ne faites pas confiance aux entrées utilisateur
- Ne stockez pas les mots de passe en texte clair
- Ne codez pas les secrets en dur
- N’omettez pas la vérification des certificats
- N’exposez pas de données sensibles dans les journaux
- N’utilisez pas de chiffrement faible
- N’ignorez pas les mises à jour de sécurité

## Liste de contrôle de sécurité

- [ ] Toutes les entrées utilisateur sont validées
- [ ] Les mots de passe sont hachés avec bcrypt
- [ ] Les données sensibles sont chiffrées
- [ ] HTTPS est utilisé pour tous les appels réseau
- [ ] La limitation du débit est mise en œuvre
- [ ] Les chemins de fichiers sont validés
- [ ] Les dépendances sont à jour
- [ ] La journalisation des événements de sécurité est activée
- [ ] Les messages d’erreur ne divulguent aucune information
- [ ] Le code a été examiné afin de détecter les vulnérabilités

## Étapes suivantes

- [Architecture](/guides/architecture/) - Modèles d’architecture d’application
- [Bonnes pratiques](/features/bindings/best-practices/) - Bonnes pratiques relatives aux liaisons
