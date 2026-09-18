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

## 日誌去識別化

Wails 會自動遮蔽 IPC 日誌中的敏感資料，防止意外洩漏機密資訊。此功能**預設啟用**，並提供合理的預設設定。

### 預設保護

下列欄位名稱會自動遮蔽（不區分大小寫，以子字串比對）：

- **身分驗證**: `password`, `passwd`, `pwd`, `token`, `bearer`, `jwt`, `access_token`, `refresh_token`, `secret`, `apikey`, `api_key`, `auth`, `authorization`, `credential`
- **密碼學**: `private`, `privatekey`, `private_key`, `signing`, `encryption_key`
- **工作階段**: `session`, `sessionid`, `session_id`, `cookie`, `csrf`, `xsrf`

此外，也會在值中偵測下列模式：

- JWT 權杖 (`eyJhbG...`)
- Bearer 權杖 (`Bearer xxx`)
- 常見 API 金鑰格式 (`sk_live_xxx`, `pk_test_xxx`)

### 自訂設定

使用所有可用選項設定去識別化：

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

### 設定選項

| 選項 | 說明 |
| --- | --- |
| `RedactFields` | 額外需要遮蔽的欄位名稱（與預設值合併） |
| `RedactPatterns` | 額外用於比對值的正規表示式 |
| `CustomSanitizeFunc` | 完全控制處理過程的函式；傳回 `(value, true)` 可覆寫結果 |
| `DisableDefaults` | 僅使用明確指定的欄位與模式 |
| `Replacement` | 自訂替代字串（預設：`***`） |
| `Disabled` | 完全停用去識別化（請謹慎使用） |

### 公開去識別化 API

使用去識別化工具處理自己的資料：

```go
// Get the application's sanitizer
sanitizer := app.Sanitizer()

// Sanitize a map
cleanData := sanitizer.SanitizeMap(sensitiveData)

// Sanitize JSON
cleanJSON := sanitizer.SanitizeJSON(jsonBytes)
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
- 請勿在記錄中洩露敏感資料 (Wails 預設會將 IPC 日誌去識別化)
- 請勿使用強度不足的加密
- 請勿忽略安全性更新
- 不要在正式環境中停用日誌去識別化

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
- [ ] 已針對應用程式特定欄位設定日誌去識別化
- [ ] 自訂日誌未暴露敏感欄位

## 後續步驟

- [架構](/guides/architecture/) - 應用程式架構模式
- [最佳實務](/features/bindings/best-practices/) - 繫結最佳實務
