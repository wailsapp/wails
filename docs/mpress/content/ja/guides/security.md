---
title: "セキュリティのベストプラクティス"
description: "Wails アプリケーションを安全に保つ"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概要

デスクトップアプリケーションでは、セキュリティが極めて重要です。以下のプラクティスに従って、アプリケーションを安全に保ってください。

## 入力の検証

### 必ず検証する

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

### HTML をサニタイズする

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

## 認証

### パスワードを安全に保存する

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

### セッション管理

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

## データ保護

### 機密データを暗号化する

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

### 安全なストレージ

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

## ネットワークセキュリティ

### HTTPS を使用する

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### 証明書を検証する

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

## ファイル操作

### パスを検証する

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

## レート制限

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

## ベストプラクティス

### ✅ 推奨事項

- すべての入力を検証する
- ネットワーク通信には HTTPS を使用する
- 機密データを暗号化する
- 安全なパスワードハッシュを使用する
- レート制限を実装する
- 依存関係を最新の状態に保つ
- セキュリティイベントをログに記録する
- OS のキーチェーンを使用する

### ❌ 禁止事項

- ユーザー入力を信用しない
- パスワードを平文で保存しない
- シークレットをハードコードしない
- 証明書の検証を省略しない
- 機密データをログに出力しない
- 脆弱な暗号化を使用しない
- セキュリティ更新を無視しない

## セキュリティチェックリスト

- [ ] すべてのユーザー入力を検証済み
- [ ] bcrypt でパスワードをハッシュ化済み
- [ ] 機密データを暗号化済み
- [ ] すべてのネットワーク通信で HTTPS を使用
- [ ] レート制限を実装済み
- [ ] ファイルパスを検証済み
- [ ] 依存関係が最新
- [ ] セキュリティログを有効化済み
- [ ] エラーメッセージから情報が漏えいしない
- [ ] 脆弱性についてコードをレビュー済み

## 次のステップ

- [アーキテクチャ](/guides/architecture/) - アプリケーションのアーキテクチャパターン
- [ベストプラクティス](/features/bindings/best-practices/) - バインディングのベストプラクティス
