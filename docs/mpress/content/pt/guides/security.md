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

Use um nonce novo para cada mensagem criptografada com determinada chave. Na cadeia de ferramentas do Go exigida pelo Wails v3, `NewGCMWithRandomNonce` gera o nonce e o adiciona como prefixo automaticamente. Armazene os bytes retornados sem alterações; `decrypt` lê o nonce desse prefixo. Mantenha a chave separada do texto cifrado.

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

Use uma chave AES adequada, gerada aleatoriamente, e respeite o limite de mensagens por chave do AEAD. Consulte a [referência de cifras do Go](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

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

### Verifique os certificados

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

## Operações com arquivos

### Restrinja o acesso a arquivos

Aceite um caminho relativo a um diretório controlado pelo aplicativo. `os.OpenRoot` restringe a resolução do sistema de arquivos a esse diretório, inclusive para links simbólicos. Rejeitar `..` ou verificar apenas o prefixo de uma string não estabelece esse limite.

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

Por exemplo, chame `readFile(appDataDir, "notes/today.json")`. Escolha `appDataDir` em código de backend confiável; não permita que o frontend forneça a raiz autorizada. Uma raiz limita a navegação pelo sistema de arquivos, mas não determina quais arquivos dentro dela o usuário está autorizado a ler.

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
- Use HTTPS nas chamadas de rede
- Criptografe dados confidenciais
- Use hash seguro para senhas
- Implemente a limitação de taxa
- Mantenha as dependências atualizadas
- Registre eventos de segurança
- Use os cofres de chaves do sistema operacional

### ❌ Não faça

- Não confie nas entradas do usuário
- Não armazene senhas em texto simples
- Não codifique segredos diretamente no código
- Não ignore a verificação de certificados
- Não exponha dados confidenciais nos logs
- Não use criptografia fraca
- Não ignore atualizações de segurança

## Lista de verificação de segurança

- [ ] Todas as entradas do usuário foram validadas
- [ ] Senhas protegidas com hash bcrypt
- [ ] Dados confidenciais criptografados
- [ ] HTTPS usado em todas as chamadas de rede
- [ ] Limitação de taxa implementada
- [ ] Caminhos de arquivos validados
- [ ] Dependências atualizadas
- [ ] Registro de eventos de segurança habilitado
- [ ] As mensagens de erro não expõem informações
- [ ] Código revisado em busca de vulnerabilidades

## Próximas etapas

- [Arquitetura](/guides/architecture/) - Padrões de arquitetura de aplicativos
- [Práticas recomendadas](/features/bindings/best-practices/) - Práticas recomendadas para bindings
