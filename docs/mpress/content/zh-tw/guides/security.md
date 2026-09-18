---
title: "安全性最佳實務"
description: "保護您的 Wails 應用程式"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概觀

安全性對桌面應用程式至關重要。請遵循這些實務，確保應用程式安全。

## 輸入驗證

### 一律驗證

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

### 清理 HTML

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

## 身分驗證

### 安全儲存密碼

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

### 工作階段管理

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

## 資料保護

### 加密敏感資料

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

### 安全儲存

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

## 網路安全

### 使用 HTTPS

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### 驗證憑證

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

## 檔案操作

### 驗證路徑

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

## 速率限制

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

## 最佳實務

### ✅ 應做事項

- 驗證所有輸入
- 網路呼叫使用 HTTPS
- 加密敏感資料
- 使用安全的密碼雜湊
- 實作速率限制
- 持續更新相依套件
- 記錄安全性事件
- 使用作業系統的鑰匙圈

### ❌ 請勿

- 請勿信任使用者輸入
- 請勿以純文字儲存密碼
- 請勿將機密資訊寫死在程式碼中
- 請勿略過憑證驗證
- 請勿在記錄中洩露敏感資料
- 請勿使用強度不足的加密
- 請勿忽略安全性更新

## 安全性檢查清單

- [ ] 已驗證所有使用者輸入
- [ ] 已使用 bcrypt 對密碼進行雜湊處理
- [ ] 已加密敏感資料
- [ ] 所有網路呼叫均使用 HTTPS
- [ ] 已實作速率限制
- [ ] 已驗證檔案路徑
- [ ] 相依套件均為最新版本
- [ ] 已啟用安全性記錄
- [ ] 錯誤訊息不會洩露資訊
- [ ] 已審查程式碼是否存在漏洞

## 後續步驟

- [架構](/guides/architecture/) - 應用程式架構模式
- [最佳實務](/features/bindings/best-practices/) - 繫結最佳實務
