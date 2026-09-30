---
title: "Práticas recomendadas de segurança"
description: "Proteja seu aplicativo Wails"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Visão geral

A segurança é essencial para aplicativos desktop. Siga estas práticas para manter seu aplicativo seguro.

## Validação de entradas

### Sempre valide

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

### Sanitize o HTML

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

## Autenticação

### Armazenamento seguro de senhas

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

### Gerenciamento de sessões

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

## Proteção de dados

### Criptografe dados confidenciais

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

### Armazenamento seguro

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

## Segurança de rede

### Use HTTPS

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### Verifique os certificados

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

## Operações com arquivos

### Valide os caminhos

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

## Limitação de taxa

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

## Práticas recomendadas

### ✅ Faça

- Valide todas as entradas
- Use HTTPS em chamadas de rede
- Criptografe dados confidenciais
- Use hashing seguro de senhas
- Implemente a limitação de taxa
- Mantenha as dependências atualizadas
- Registre eventos de segurança
- Use os cofres de credenciais do sistema operacional

### ❌ Não faça

- Não confie nas entradas fornecidas pelo usuário
- Não armazene senhas em texto simples
- Não insira segredos diretamente no código
- Não ignore a verificação de certificados
- Não exponha dados confidenciais nos logs
- Não use criptografia fraca
- Não ignore atualizações de segurança

## Lista de verificação de segurança

- [ ] Todas as entradas fornecidas pelo usuário estão validadas
- [ ] As senhas estão protegidas por hash com bcrypt
- [ ] Os dados confidenciais estão criptografados
- [ ] HTTPS é usado em todas as chamadas de rede
- [ ] A limitação de taxa está implementada
- [ ] Os caminhos de arquivos estão validados
- [ ] As dependências estão atualizadas
- [ ] O registro de eventos de segurança está habilitado
- [ ] As mensagens de erro não expõem informações
- [ ] O código foi revisado em busca de vulnerabilidades

## Próximas etapas

- [Arquitetura](/guides/architecture/) - Padrões de arquitetura de aplicativos
- [Práticas recomendadas](/features/bindings/best-practices/) - Práticas recomendadas para bindings
