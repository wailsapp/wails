---
title: "セキュリティのベストプラクティス"
description: "Wails アプリケーションを安全に保つ"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概要

デスクトップアプリケーションでは、セキュリティが極めて重要です。アプリケーションを安全に保つため、以下のプラクティスに従ってください。

## 入力検証

### 常に検証する

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

### パスワードの安全な保存

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

指定した鍵で暗号化するメッセージごとに、新しい nonce を使用してください。Wails v3 が必要とする Go ツールチェーンでは、`NewGCMWithRandomNonce` が nonce を自動的に生成し、先頭に付加します。返されたバイト列はそのまま保存してください。`decrypt` は、その接頭部から nonce を読み取ります。鍵は暗号文とは別に保管してください。

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

適切なランダム生成 AES 鍵を使用し、AEAD の鍵ごとのメッセージ数上限を守ってください。[Go の cipher リファレンス](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce)を参照してください。

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

### 証明書を検証する

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

## ファイル操作

### ファイルアクセスを制限する

アプリケーションが管理するディレクトリからの相対パスを受け入れてください。`os.OpenRoot` は、シンボリックリンクを含め、ファイルシステム上のパス解決をそのディレクトリ内に制限します。`..` を拒否したり、文字列の接頭辞だけを確認したりしても、この境界は確保できません。

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

たとえば、`readFile(appDataDir, "notes/today.json")` を呼び出します。信頼できるバックエンドコード内で `appDataDir` を選択し、許可するルートをフロントエンドから指定できないようにしてください。ルートによって制限されるのはファイルシステムの探索範囲であり、その内部でユーザーがどのファイルの読み取りを許可されているかではありません。

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
- ネットワーク呼び出しに HTTPS を使用する
- 機密データを暗号化する
- 安全なパスワードハッシュを使用する
- レート制限を実装する
- 依存関係を最新の状態に保つ
- セキュリティイベントを記録する
- OS のキーチェーンを使用する

### ❌ 禁止事項

- ユーザー入力を信用しない
- パスワードを平文で保存しない
- シークレットをハードコードしない
- 証明書の検証を省略しない
- ログに機密データを出力しない
- 脆弱な暗号化方式を使用しない
- セキュリティ更新を無視しない

## セキュリティチェックリスト

- [ ] すべてのユーザー入力を検証済み
- [ ] bcrypt でパスワードをハッシュ化済み
- [ ] 機密データを暗号化済み
- [ ] すべてのネットワーク呼び出しで HTTPS を使用
- [ ] レート制限を実装済み
- [ ] ファイルパスを検証済み
- [ ] 依存関係を最新の状態に更新済み
- [ ] セキュリティログを有効化済み
- [ ] エラーメッセージから情報が漏えいしない
- [ ] 脆弱性についてコードをレビュー済み

## 次のステップ

- [アーキテクチャ](/guides/architecture/) - アプリケーションアーキテクチャのパターン
- [ベストプラクティス](/features/bindings/best-practices/) - バインディングのベストプラクティス
