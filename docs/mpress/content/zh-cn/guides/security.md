---
title: "安全最佳实践"
description: "保护 Wails 应用的安全"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概述

安全性对于桌面应用至关重要。请遵循这些实践，确保应用安全。

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

对于使用同一密钥加密的每条消息，都应使用新的 nonce。在 Wails v3 所需的 Go 工具链中，`NewGCMWithRandomNonce`会自动生成 nonce 并将其添加到数据前面。请完整存储返回的字节；`decrypt`会从该前缀中读取 nonce。密钥应与密文分开保存。

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

使用适当的随机生成 AES 密钥，并遵守 AEAD 对每个密钥的消息数量限制。请参阅[Go 密码学参考文档](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce)。

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

### 验证证书

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

## 文件操作

### 限制文件访问范围

只接受相对于应用所控制目录的路径。`os.OpenRoot`会将文件系统路径解析限制在该目录内，符号链接也不例外。仅拒绝`..`或检查字符串前缀无法建立此边界。

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

例如，调用`readFile(appDataDir, "notes/today.json")`。请在可信的后端代码中选择`appDataDir`；不要允许前端提供获准访问的根目录。根目录只能限制文件系统遍历范围，无法限制用户有权读取其中的哪些文件。

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

## 最佳实践

### ✅ 应该做

- 验证所有输入
- 网络调用使用 HTTPS
- 加密敏感数据
- 使用安全的密码哈希算法
- 实施速率限制
- 及时更新依赖项
- 记录安全事件
- 使用操作系统密钥链

### ❌ 不应该做

- 不要信任用户输入
- 不要以明文形式存储密码
- 不要硬编码机密信息
- 不要跳过证书验证
- 不要在日志中暴露敏感数据
- 不要使用弱加密算法
- 不要忽略安全更新

## 安全检查清单

- [ ] 已验证所有用户输入
- [ ] 已使用 bcrypt 对密码进行哈希处理
- [ ] 已加密敏感数据
- [ ] 所有网络调用均使用 HTTPS
- [ ] 已实施速率限制
- [ ] 已验证文件路径
- [ ] 依赖项均为最新版本
- [ ] 已启用安全日志记录
- [ ] 错误消息不会泄露信息
- [ ] 已审查代码中的漏洞

## 后续步骤

- [架构](/guides/architecture/) - 应用架构模式
- [最佳实践](/features/bindings/best-practices/) - 绑定最佳实践
