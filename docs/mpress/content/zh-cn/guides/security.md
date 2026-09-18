---
title: "安全最佳实践"
description: "保护您的 Wails 应用程序"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概述

安全性对于桌面应用程序至关重要。请遵循这些实践，确保应用程序安全。

## 输入验证

### 始终进行验证

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

### 净化 HTML

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

## 身份验证

### 安全存储密码

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

### 会话管理

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

## 数据保护

### 加密敏感数据

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

### 安全存储

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

## 网络安全

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

### 验证证书

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

## 文件操作

### 验证路径

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

## 日志脱敏

Wails 会自动对 IPC 日志中的敏感数据进行脱敏，防止意外泄露秘密信息。此功能**默认启用**，并提供合理的默认设置。

### 默认保护

以下字段名称会被自动脱敏（不区分大小写，按子字符串匹配）：

- **身份验证**: `password`, `passwd`, `pwd`, `token`, `bearer`, `jwt`, `access_token`, `refresh_token`, `secret`, `apikey`, `api_key`, `auth`, `authorization`, `credential`
- **密码学**: `private`, `privatekey`, `private_key`, `signing`, `encryption_key`
- **会话**: `session`, `sessionid`, `session_id`, `cookie`, `csrf`, `xsrf`

此外，还会在值中检测以下模式：

- JWT 令牌 (`eyJhbG...`)
- Bearer 令牌 (`Bearer xxx`)
- 常见 API 密钥格式 (`sk_live_xxx`, `pk_test_xxx`)

### 自定义配置

使用全部可用选项配置脱敏：

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

### 配置选项

| 选项 | 描述 |
| --- | --- |
| `RedactFields` | 额外需要脱敏的字段名称（与默认值合并） |
| `RedactPatterns` | 额外用于匹配值的正则表达式 |
| `CustomSanitizeFunc` | 完全控制处理过程的函数；返回 `(value, true)` 可覆盖结果 |
| `DisableDefaults` | 仅使用明确指定的字段和模式 |
| `Replacement` | 自定义替换字符串（默认：`***`） |
| `Disabled` | 完全禁用脱敏（请谨慎使用） |

### 公共脱敏 API

使用脱敏器处理自己的数据：

```go
// Get the application's sanitizer
sanitizer := app.Sanitizer()

// Sanitize a map
cleanData := sanitizer.SanitizeMap(sensitiveData)

// Sanitize JSON
cleanJSON := sanitizer.SanitizeJSON(jsonBytes)
```

## 最佳实践

### ✅ 应该做

- 验证所有输入
- 使用 HTTPS 发起网络请求
- 加密敏感数据
- 使用安全的密码哈希算法
- 实施速率限制
- 及时更新依赖项
- 记录安全事件
- 使用操作系统密钥链

### ❌ 不应该做

- 不要信任用户输入
- 不要以明文存储密码
- 不要硬编码机密信息
- 不要跳过证书验证
- 不要在日志中暴露敏感数据 (Wails 默认对 IPC 日志进行脱敏)
- 不要使用安全性弱的加密算法
- 不要忽略安全更新
- 不要在生产环境中禁用日志脱敏

## 安全检查清单

- [ ] 已验证所有用户输入
- [ ] 已使用 bcrypt 对密码进行哈希处理
- [ ] 已加密敏感数据
- [ ] 所有网络请求均使用 HTTPS
- [ ] 已实施速率限制
- [ ] 已验证文件路径
- [ ] 依赖项均为最新版本
- [ ] 已启用安全日志记录
- [ ] 错误消息不会泄露信息
- [ ] 已审查代码中的漏洞
- [ ] 已为应用特定字段配置日志脱敏
- [ ] 自定义日志未暴露敏感字段

## 后续步骤

- [架构](/guides/architecture/)——应用程序架构模式
- [最佳实践](/features/bindings/best-practices/)——绑定最佳实践
