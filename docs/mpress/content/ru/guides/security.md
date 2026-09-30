---
title: "Рекомендации по обеспечению безопасности"
description: "Обеспечьте безопасность приложения Wails"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Обзор

Безопасность критически важна для настольных приложений. Следуйте этим рекомендациям, чтобы обеспечить безопасность своего приложения.

## Проверка входных данных

### Всегда проверяйте входные данные

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

### Очищайте HTML

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

## Аутентификация

### Безопасное хранение паролей

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

### Управление сеансами

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

## Защита данных

### Шифруйте конфиденциальные данные

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

### Безопасное хранение

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

## Сетевая безопасность

### Используйте HTTPS

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### Проверяйте сертификаты

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

## Операции с файлами

### Проверяйте пути

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

## Ограничение частоты запросов

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

## Рекомендации

### ✅ Делайте

- Проверяйте все входные данные
- Используйте HTTPS для сетевых запросов
- Шифруйте конфиденциальные данные
- Используйте безопасное хеширование паролей
- Реализуйте ограничение частоты запросов
- Своевременно обновляйте зависимости
- Регистрируйте события безопасности
- Используйте системные хранилища учётных данных

### ❌ Не делайте

- Не доверяйте пользовательскому вводу
- Не храните пароли в открытом виде
- Не встраивайте секреты непосредственно в код
- Не пропускайте проверку сертификатов
- Не раскрывайте конфиденциальные данные в журналах
- Не используйте ненадёжное шифрование
- Не игнорируйте обновления безопасности

## Контрольный список безопасности

- [ ] Все пользовательские входные данные проверены
- [ ] Пароли хешируются с помощью bcrypt
- [ ] Конфиденциальные данные зашифрованы
- [ ] Для всех сетевых запросов используется HTTPS
- [ ] Реализовано ограничение частоты запросов
- [ ] Пути к файлам проверены
- [ ] Зависимости обновлены до актуальных версий
- [ ] Включено журналирование событий безопасности
- [ ] Сообщения об ошибках не раскрывают информацию
- [ ] Код проверен на наличие уязвимостей

## Дальнейшие действия

- [Архитектура](/guides/architecture/) — шаблоны архитектуры приложений
- [Рекомендации](/features/bindings/best-practices/) — рекомендации по использованию привязок
