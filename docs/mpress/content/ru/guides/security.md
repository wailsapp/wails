---
title: "Рекомендации по безопасности"
description: "Обеспечение безопасности приложения Wails"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Обзор

Безопасность критически важна для настольных приложений. Следуйте этим рекомендациям, чтобы защитить свое приложение.

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

Для каждого сообщения, зашифрованного данным ключом, используйте новый nonce. В инструментарии Go, который требуется для Wails v3, `NewGCMWithRandomNonce` автоматически генерирует nonce и добавляет его в начало данных. Храните возвращенные байты без изменений: `decrypt` считывает nonce из этого префикса. Храните ключ отдельно от шифротекста.

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

Используйте подходящий случайно сгенерированный ключ AES и соблюдайте установленное для AEAD ограничение на количество сообщений для одного ключа. См. [справочник по шифрованию в Go](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

### Безопасное хранение данных

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

### Проверяйте сертификаты

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

## Операции с файлами

### Ограничивайте доступ к файлам

Принимайте путь относительно каталога, контролируемого приложением. `os.OpenRoot` ограничивает разрешение путей файловой системы этим каталогом, в том числе при использовании символических ссылок. Отклонение `..` или одна лишь проверка строкового префикса не обеспечивают соблюдение этой границы.

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

Например, вызовите `readFile(appDataDir, "notes/today.json")`. Задавайте `appDataDir` в доверенном коде серверной части; не позволяйте фронтенду передавать разрешенный корневой каталог. Корневой каталог ограничивает обход файловой системы, но не определяет, какие файлы внутри него пользователю разрешено читать.

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
- Используйте системные хранилища ключей ОС

### ❌ Не делайте

- Не доверяйте пользовательскому вводу
- Не храните пароли в открытом виде
- Не указывайте секреты непосредственно в коде
- Не пропускайте проверку сертификатов
- Не раскрывайте конфиденциальные данные в журналах
- Не используйте слабое шифрование
- Не игнорируйте обновления безопасности

## Контрольный список безопасности

- [ ] Все пользовательские данные проверяются
- [ ] Пароли хешируются с помощью bcrypt
- [ ] Конфиденциальные данные шифруются
- [ ] Для всех сетевых запросов используется HTTPS
- [ ] Реализовано ограничение частоты запросов
- [ ] Пути к файлам проверяются
- [ ] Зависимости обновлены
- [ ] Включено журналирование событий безопасности
- [ ] Сообщения об ошибках не раскрывают информацию
- [ ] Код проверен на наличие уязвимостей

## Дальнейшие шаги

- [Архитектура](/guides/architecture/) — Шаблоны архитектуры приложений
- [Рекомендации](/features/bindings/best-practices/) — Рекомендации по привязкам
