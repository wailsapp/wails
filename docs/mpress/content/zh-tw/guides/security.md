---
title: "安全性最佳實務"
description: "保護您的 Wails 應用程式"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概觀

安全性對桌面應用程式至關重要。請遵循這些實務，確保應用程式安全。

## 輸入驗證

### 一律進行驗證

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

### 淨化 HTML

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

使用指定金鑰加密每則訊息時，都應使用全新的 nonce。在 Wails v3所需的 Go 工具鏈中，`NewGCMWithRandomNonce`會自動產生 nonce 並將其加到資料前面。請完整儲存傳回的位元組；`decrypt`會從此前綴讀取 nonce。金鑰應與密文分開存放。

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

請使用適當且隨機產生的 AES 金鑰，並遵守 AEAD 對每把金鑰的訊息數量限制。請參閱[Go 密碼學參考資料](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce)。

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

### 驗證憑證

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

## 檔案操作

### 限制檔案存取範圍

僅接受相對於應用程式所控目錄的路徑。`os.OpenRoot`會將檔案系統的路徑解析限制在該目錄內，符號連結也包含在內。僅拒絕`..`或檢查字串前綴，無法建立這項邊界。

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

例如，呼叫`readFile(appDataDir, "notes/today.json")`。請在受信任的後端程式碼中選擇`appDataDir`；不要讓前端提供允許的根目錄。根目錄只能限制檔案系統的遍歷範圍，不能決定使用者獲授權讀取其中哪些檔案。

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

### ✅ 應該做

- 驗證所有輸入
- 網路呼叫使用 HTTPS
- 加密敏感資料
- 使用安全的密碼雜湊
- 實作速率限制
- 讓相依套件保持最新狀態
- 記錄安全性事件
- 使用作業系統的鑰匙圈

### ❌ 不應該做

- 不要信任使用者輸入
- 不要以純文字儲存密碼
- 不要將密鑰硬式編碼
- 不要略過憑證驗證
- 不要在日誌中揭露敏感資料
- 不要使用弱式加密
- 不要忽略安全性更新

## 安全性檢查清單

- [ ] 已驗證所有使用者輸入
- [ ] 已使用 bcrypt 雜湊密碼
- [ ] 已加密敏感資料
- [ ] 所有網路呼叫皆使用 HTTPS
- [ ] 已實作速率限制
- [ ] 已驗證檔案路徑
- [ ] 相依套件均為最新版本
- [ ] 已啟用安全性記錄
- [ ] 錯誤訊息不會洩漏資訊
- [ ] 已檢閱程式碼中的安全漏洞

## 後續步驟

- [架構](/guides/architecture/) - 應用程式架構模式
- [最佳實務](/features/bindings/best-practices/) - 繫結最佳實務
